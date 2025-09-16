/**
 * (C) Copyright IBM Corp. 2021.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/IBM/appconfiguration-go-sdk/lib/internal/constants"
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/messages"
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/models"
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/utils"
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/utils/log"
	"github.com/IBM/go-sdk-core/v5/core"
	sm "github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
	"github.com/gorilla/websocket"
)

// variable to keep track of status of server-client connection
var isAlive = false

type configurationUpdateListenerFunc func()

// ConfigurationHandler : Configuration Handler
type ConfigurationHandler struct {
	isInitialized               bool
	collectionID                string
	environmentID               string
	apikey                      string
	guid                        string
	region                      string
	usePrivateEndpoint          bool
	urlBuilder                  *utils.URLBuilder
	appConfig                   *AppConfiguration
	cache                       *models.Cache
	configurationUpdateListener configurationUpdateListenerFunc
	persistentCacheDirectory    string
	bootstrapFile               string
	liveConfigUpdateEnabled     bool
	persistentData              []byte
	retryInterval               int64
	scheduledRetry              *time.Timer
	socketConnection            *websocket.Conn
	socketConnectionResponse    *http.Response
	mu                          sync.Mutex
}

var configurationHandlerInstance *ConfigurationHandler

// GetConfigurationHandlerInstance : Get Configuration Handler Instance
func GetConfigurationHandlerInstance() *ConfigurationHandler {
	if configurationHandlerInstance == nil {
		configurationHandlerInstance = new(ConfigurationHandler)
	}
	return configurationHandlerInstance
}

// Init : Init App Configuration Instance
func (ch *ConfigurationHandler) Init(region, guid, apikey string, usePrivateEndpoint bool) {
	ch.region = region
	ch.guid = guid
	ch.apikey = apikey
	ch.usePrivateEndpoint = usePrivateEndpoint
}

// SetContext : Set Context
func (ch *ConfigurationHandler) SetContext(collectionID, environmentID string, options ContextOptions) {
	ch.collectionID = collectionID
	ch.environmentID = environmentID
	ch.urlBuilder = utils.GetInstance()
	ch.urlBuilder.Init(ch.collectionID, ch.environmentID, ch.region, ch.guid, ch.apikey, overrideServiceUrl, ch.usePrivateEndpoint)
	utils.GetMeteringInstance().Init(ch.guid, environmentID, collectionID)
	ch.persistentCacheDirectory = options.PersistentCacheDirectory
	ch.bootstrapFile = options.BootstrapFile
	ch.liveConfigUpdateEnabled = options.LiveConfigUpdateEnabled
	ch.isInitialized = true
	ch.retryInterval = 2 // two minutes
}
func (ch *ConfigurationHandler) loadData() {
	persistentCacheRead := false

	if len(ch.persistentCacheDirectory) > 0 {
		path := filepath.Join(utils.SanitizePath(ch.persistentCacheDirectory), constants.ConfigurationFile)
		log.Info(messages.ReadPersistentCache, path)
		ch.persistentData = utils.ReadFiles(path)
		if !bytes.Equal(ch.persistentData, []byte(`{}`)) {
			configurations, err := models.ExtractConfigurations(ch.persistentData, ch.environmentID, ch.collectionID)
			if err != nil {
				log.Error("Error occurred while reading persistent cache configurations - ", err.Error())
			} else {
				ch.saveInCache(configurations)
				persistentCacheRead = true
			}
		}
	}
	if len(ch.bootstrapFile) > 0 {
		path := utils.SanitizePath(ch.bootstrapFile)
		if len(ch.persistentCacheDirectory) > 0 {
			if !persistentCacheRead {
				bootstrapFileData := utils.ReadFiles(path)
				bootstrapConfigurations, err := models.ExtractConfigurations(bootstrapFileData, ch.environmentID, ch.collectionID)
				if err != nil {
					log.Error("Error occurred while reading bootstrap configurations - ", err.Error())
				} else {
					ch.saveInCache(bootstrapConfigurations)
					go utils.StoreFiles(string(models.FormatConfig(bootstrapConfigurations, ch.environmentID, ch.collectionID)), ch.persistentCacheDirectory)

				}
			}
		} else {
			log.Info(messages.ReadBootstrapConfigurations, path)
			bootstrapFileData := utils.ReadFiles(path)
			bootstrapConfigurations, err := models.ExtractConfigurations(bootstrapFileData, ch.environmentID, ch.collectionID)
			if err != nil {
				log.Error("Error occurred while reading bootstrap configurations - ", err.Error())
			} else {
				ch.saveInCache(bootstrapConfigurations)
			}
		}
	}
	if ch.liveConfigUpdateEnabled {
		ch.FetchConfigurationData()
	}
}

// FetchConfigurationData : Fetch Configuration Data
func (ch *ConfigurationHandler) FetchConfigurationData() {
	log.Debug(messages.FetchConfigurationData)
	if ch.isInitialized {
		ch.fetchFromAPI()
		go ch.startWebSocket()
	}
}
func (ch *ConfigurationHandler) saveInCache(data []byte) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	configurations := models.CacheConfig{}
	err := json.Unmarshal(data, &configurations)
	if err != nil {
		log.Error(messages.UnmarshalJSONErr, err)
		return
	}
	log.Debug(configurations)
	featureMap := make(map[string]models.Feature)
	for _, feature := range configurations.Features {
		featureMap[feature.GetFeatureID()] = feature.Feature
	}

	propertyMap := make(map[string]models.Property)
	for _, property := range configurations.Properties {
		propertyMap[property.GetPropertyID()] = property.Property
	}

	segmentMap := make(map[string]models.Segment)
	for _, segment := range configurations.Segments {
		segmentMap[segment.GetSegmentID()] = segment
	}
	log.Debug(messages.SetInMemoryCache)
	models.SetCache(featureMap, propertyMap, segmentMap)
	ch.cache = models.GetCacheInstance()
}
func (ch *ConfigurationHandler) updateCacheAndListener(data []byte) {
	ch.saveInCache(data)
	if ch.configurationUpdateListener != nil {
		ch.configurationUpdateListener()
	}
}
func extractErrorMessage(err error, response *core.DetailedResponse) string {
	if err != nil {
		return err.Error()
	}
	if response == nil {
		return ""
	}
	if response.Result != nil {
		resultStr, ok := response.Result.(string)
		if ok {
			return resultStr
		}
	}
	return string(response.RawResult)
}
func (ch *ConfigurationHandler) fetchFromAPI() {
	if ch.isInitialized {
		builder := core.NewRequestBuilder(core.GET)
		builder.AddQuery("action", "sdkConfig")
		builder.AddQuery("collection_id", ch.collectionID)
		builder.AddQuery("environment_id", ch.environmentID)
		pathParamsMap := map[string]string{
			"guid": ch.guid,
		}
		_, err := builder.ResolveRequestURL(ch.urlBuilder.GetBaseServiceURL(), `/apprapp/feature/v1/instances/{guid}/config`, pathParamsMap)
		if err != nil {
			log.Error(err.Error())
			return
		}
		builder.AddHeader("Accept", "application/json")
		builder.AddHeader("User-Agent", constants.UserAgent)

		// 2xx - Do not retry (Success)
		// 3xx - Do not retry (Redirect)
		// 4xx - Do not retry (Client errors)
		// 429 - Retry ("Too Many Requests")
		// 5xx - Retry (Server errors)

		// The Request() below is itself an retryableRequest. Hence, we don't need to write the retry logic again.
		//
		// The API call gets retried within Request() for 3 times in an exponential interval(0.5s, 1s, 1.5s) between each retry.
		// If all the 3 retries fails, the call is returned and execution is given back to us to take the response object ahead.
		//
		// For 429 error code - The Request() will retry the request 3 times in an interval of time mentioned in ["Retry-after"] header.
		// If all the 3 retries exhausts the call is returned and execution is given back to us to take the response object ahead.
		//
		// Both the cases [429, 499 & 5xx] we schedule a retry after 2 minutes.

		response, err := utils.GetAPIManagerInstance().Request(builder)
		if response != nil && response.StatusCode == 200 {
			log.Info(messages.FetchAPISuccessful)
			jsonData, _ := json.Marshal(response.Result)
			configurations, err := models.ExtractConfigurations(jsonData, ch.environmentID, ch.collectionID)
			if err != nil {
				log.Error("Error occurred while reading fetched configurations - ", err.Error())
				return
			}
			// asynchronously write the response to persistent volume, if enabled
			if len(ch.persistentCacheDirectory) > 0 {
				go utils.StoreFiles(string(jsonData), ch.persistentCacheDirectory)
			}
			// load the configurations in the response to cache maps
			ch.updateCacheAndListener(configurations)
		} else {
			if response != nil && response.StatusCode >= 400 && response.StatusCode < 499 && response.StatusCode != 429 {
				// Do Nothing! GET "/config" failed due to a client-side error.
				// Print the error message and return.
				log.Error(extractErrorMessage(err, response))
				return
			}
			errMessage := extractErrorMessage(err, response)
			log.Error(messages.ConfigAPIError, errMessage, fmt.Sprintf(messages.RetryScheduledMessage, ch.retryInterval))
			if ch.scheduledRetry != nil {
				ch.scheduledRetry.Stop()
			}
			ch.scheduledRetry = time.AfterFunc(time.Minute*time.Duration(ch.retryInterval), func() {
				ch.fetchFromAPI()
			})
		}
	} else {
		log.Debug(messages.FetchFromAPISdkInitError)
	}
}

func (ch *ConfigurationHandler) startWebSocket() {
	defer utils.GracefullyHandleError()
	log.Debug(messages.StartWebSocket)
	authToken := ch.urlBuilder.GetToken()
	if len(authToken) == 0 {
		log.Error(messages.WebSocketConnectFailed, messages.AuthTokenError)
		return
	}
	h := make(http.Header)
	h.Add("Authorization", authToken)
	h.Add("User-Agent", constants.UserAgent)
	var err error
	if ch.socketConnection != nil {
		isAlive = false
		ch.socketConnection.Close()
	}
	ch.socketConnection, ch.socketConnectionResponse, err = websocket.DefaultDialer.Dial(ch.urlBuilder.GetWebSocketURL(), h)
	if err != nil {
		isAlive = false
		if ch.socketConnectionResponse != nil {
			statusCode := ch.socketConnectionResponse.StatusCode
			if statusCode >= 400 && statusCode < 499 && statusCode != 429 {
				// websocket dial that fails with response status code in between 400-499, except 429 & 499, are not retried.
				// Do Nothing! Since websocket connect failed due to a client-side error.
				log.Error(messages.WebSocketConnectErr+err.Error(), " ", statusCode)
				return
			}
		}
		log.Error(messages.WebSocketConnectErr, err.Error(), " Retrying websocket connect in 15 seconds...")
		time.Sleep(15 * time.Second)
		go ch.startWebSocket()
		return
	}
	log.Debug(messages.WebSocketConnectSuccess)
	// defer c.Close()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if ch.socketConnection != nil {
				isAlive = true
				_, message, err := ch.socketConnection.ReadMessage()
				log.Debug(string(message))
				if err != nil {
					isAlive = false
					log.Error(messages.WebsocketErrorReadingMessage, err.Error(), " Retrying websocket connect in 15 seconds...")
					time.Sleep(15 * time.Second)
					go ch.startWebSocket()
					return
				}
				if string(message) != "test message" {
					log.Debug(messages.WebsocketReceivingMessage + string(message))
					ch.fetchFromAPI()
				}
			} else {
				isAlive = false
				time.Sleep(15 * time.Second)
				go ch.startWebSocket()
				return
			}
		}
	}()
}
func (ch *ConfigurationHandler) getFeatures() (map[string]models.Feature, error) {
	if ch.cache == nil {
		return nil, errors.New(messages.InitError)
	}
	return ch.cache.FeatureMap, nil
}
func (ch *ConfigurationHandler) getFeature(featureID string) (models.Feature, error) {
	if ch.cache != nil && len(ch.cache.FeatureMap) > 0 {
		if val, ok := ch.cache.FeatureMap[featureID]; ok {
			return val, nil
		}
	}
	log.Error(messages.InvalidFeatureID, featureID)
	return models.Feature{}, errors.New(messages.ErrorInvalidFeatureID + featureID)

}
func (ch *ConfigurationHandler) getProperties() (map[string]models.Property, error) {
	if ch.cache == nil {
		return nil, errors.New(messages.InitError)
	}
	return ch.cache.PropertyMap, nil
}
func (ch *ConfigurationHandler) getProperty(propertyID string) (models.Property, error) {
	if ch.cache != nil && len(ch.cache.PropertyMap) > 0 {
		if val, ok := ch.cache.PropertyMap[propertyID]; ok {
			return val, nil
		}
	}
	log.Error(messages.InvalidPropertyID, propertyID)
	return models.Property{}, errors.New(messages.ErrorInvalidPropertyID + propertyID)
}
func (ch *ConfigurationHandler) getSecret(propertyID string, secretsManagerService *sm.SecretsManagerV2) (models.SecretProperty, error) {
	property, err := ch.getProperty(propertyID)
	if err != nil {
		return models.SecretProperty{}, err
	}
	if property.GetPropertyDataType() == "SECRETREF" {
		ch.cache.SecretManagerMap[propertyID] = secretsManagerService
		return models.SecretProperty{PropertyID: propertyID}, nil
	}
	log.Error("Invalid operation: GetSecret() cannot be called on a ", property.GetPropertyDataType(), " property.")
	return models.SecretProperty{}, errors.New("error: GetSecret() cannot be called on a " + property.GetPropertyDataType() + " property.")
}

func (ch *ConfigurationHandler) registerConfigurationUpdateListener(chl configurationUpdateListenerFunc) {
	defer func() {
		if r := recover(); r != nil {
			log.Error(messages.ConfigurationUpdateListenerMethodError)
		}
	}()
	if ch.isInitialized {
		ch.configurationUpdateListener = chl
	} else {
		log.Error(messages.CollectionIDError)
	}
}
