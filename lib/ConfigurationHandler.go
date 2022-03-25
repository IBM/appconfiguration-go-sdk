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
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/constants"
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/messages"
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/models"
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/utils"
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/utils/log"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/gorilla/websocket"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

type configurationUpdateListenerFunc func()

// ConfigurationHandler : Configuration Handler
type ConfigurationHandler struct {
	isInitialized               bool
	collectionID                string
	environmentID               string
	apikey                      string
	guid                        string
	region                      string
	urlBuilder                  *utils.URLBuilder
	appConfig                   *AppConfiguration
	cache                       *models.Cache
	configurationUpdateListener configurationUpdateListenerFunc
	persistentCacheDirectory    string
	bootstrapFile               string
	liveConfigUpdateEnabled     bool
	persistentData              []byte
	retryInterval               int64
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
func (ch *ConfigurationHandler) Init(region, guid, apikey string) {
	ch.region = region
	ch.guid = guid
	ch.apikey = apikey
}

// SetContext : Set Context
func (ch *ConfigurationHandler) SetContext(collectionID, environmentID string, options ContextOptions) {
	ch.collectionID = collectionID
	ch.environmentID = environmentID
	ch.urlBuilder = utils.GetInstance()
	ch.urlBuilder.Init(ch.collectionID, ch.environmentID, ch.region, ch.guid, ch.apikey, OverrideServerHost)
	utils.GetMeteringInstance().Init(ch.guid, environmentID, collectionID)
	ch.persistentCacheDirectory = options.PersistentCacheDirectory
	ch.bootstrapFile = options.BootstrapFile
	ch.liveConfigUpdateEnabled = options.LiveConfigUpdateEnabled
	ch.isInitialized = true
	ch.retryInterval = 600
}
func (ch *ConfigurationHandler) loadData() {
	if !ch.isInitialized {
		log.Error(messages.ConfigurationHandlerInitError)
	}
	if len(ch.persistentCacheDirectory) > 0 {
		ch.persistentData = utils.ReadFiles(filepath.Join(utils.SanitizePath(ch.persistentCacheDirectory), constants.ConfigurationFile))
		if !bytes.Equal(ch.persistentData, []byte(`{}`)) {
			// no updating the listener here. Only updating cache is enough
			ch.saveInCache(ch.persistentData)
		}
	}
	if len(ch.bootstrapFile) > 0 {
		log.Debug(messages.BootstrapFileProvided)
		if len(ch.persistentCacheDirectory) > 0 {
			if bytes.Equal(ch.persistentData, []byte(`{}`)) {
				bootstrapFileData := utils.ReadFiles(utils.SanitizePath(ch.bootstrapFile))
				go utils.StoreFiles(string(bootstrapFileData), ch.persistentCacheDirectory)
				ch.updateCacheAndListener(bootstrapFileData)
			} else {
				// update the only listener here. Because, cache is already updated above (line 100)
				if ch.configurationUpdateListener != nil {
					ch.configurationUpdateListener()
				}
			}
		} else {
			bootstrapFileData := utils.ReadFiles(utils.SanitizePath(ch.bootstrapFile))
			ch.updateCacheAndListener(bootstrapFileData)
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
	configResponse := models.ConfigResponse{}
	err := json.Unmarshal(data, &configResponse)
	if err != nil {
		log.Error(messages.UnmarshalJSONErr, err)
		return
	}
	log.Debug(configResponse)
	featureMap := make(map[string]models.Feature)
	for _, feature := range configResponse.Features {
		featureMap[feature.GetFeatureID()] = feature
	}

	propertyMap := make(map[string]models.Property)
	for _, property := range configResponse.Properties {
		propertyMap[property.GetPropertyID()] = property
	}

	segmentMap := make(map[string]models.Segment)
	for _, segment := range configResponse.Segments {
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
func (ch *ConfigurationHandler) fetchFromAPI() {
	if ch.isInitialized {
		builder := core.NewRequestBuilder(core.GET)
		builder.AddQuery("environment_id", ch.environmentID)
		pathParamsMap := map[string]string{
			"guid":          ch.guid,
			"collection_id": ch.collectionID,
		}
		_, err := builder.ResolveRequestURL(ch.urlBuilder.GetBaseServiceURL(), `/apprapp/feature/v1/instances/{guid}/collections/{collection_id}/config`, pathParamsMap)
		if err != nil {
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
		// Both the cases [429 & 5xx] we schedule a retry after 10 minutes.

		response := utils.GetAPIManagerInstance().Request(builder)
		if response != nil && response.StatusCode == constants.StatusCodeGET {
			log.Debug(messages.FetchAPISuccessful)
			jsonData, _ := json.Marshal(response.Result)
			// asynchronously write the response to persistent volume, if enabled
			if len(ch.persistentCacheDirectory) > 0 {
				go utils.StoreFiles(string(jsonData), ch.persistentCacheDirectory)
			}
			// load the configurations in the response to cache maps
			ch.updateCacheAndListener(jsonData)
		} else {
			if response != nil {
				if response.Result != nil {
					log.Error(response.Result)
				} else {
					log.Error(string(response.RawResult))
				}
				if response.StatusCode == constants.StatusCodeTooManyRequests || (response.StatusCode >= constants.StatusCodeServerErrorBegin && response.StatusCode <= constants.StatusCodeServerErrorEnd) {
					time.AfterFunc(time.Second*time.Duration(ch.retryInterval), func() {
						ch.fetchFromAPI()
					})
					log.Info(messages.RetryScheduledMessage)
				}
			} else {
				log.Error(messages.ConfigAPIError)
			}
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
	h := http.Header{"Authorization": []string{authToken}}
	var err error
	if ch.socketConnection != nil {
		ch.socketConnection.Close()
	}
	ch.socketConnection, ch.socketConnectionResponse, err = websocket.DefaultDialer.Dial(ch.urlBuilder.GetWebSocketURL(), h)
	if err != nil {
		if ch.socketConnectionResponse != nil {
			log.Error(messages.WebSocketConnectErr, err, ch.socketConnectionResponse.StatusCode)
		}
		go ch.startWebSocket()
		return
	}
	// defer c.Close()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if ch.socketConnection != nil {
				_, message, err := ch.socketConnection.ReadMessage()
				log.Debug(string(message))
				if err != nil {
					log.Error(messages.WebsocketErrorReadingMessage, err.Error())
					go ch.startWebSocket()
					return
				}
				if string(message) != "test message" {
					log.Debug(messages.WebsocketReceivingMessage + string(message))
					ch.fetchFromAPI()
				}
			} else {
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
