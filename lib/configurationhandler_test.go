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
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/IBM/appconfiguration-go-sdk/lib/internal/models"
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/utils"
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/utils/log"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

var testLogger, hook = test.NewNullLogger()

func mockLogger() {
	log.SetLogger(testLogger)
}

func TestInitConfigurationHandlerInstance(t *testing.T) {
	// test init of config handler instance done properly
	ch := GetConfigurationHandlerInstance()
	ch.Init("us-south", "abc", "abc", false)
	assert.Equal(t, "us-south", ch.region)
	assert.Equal(t, "abc", ch.apikey)
	assert.Equal(t, "abc", ch.guid)
	assert.Equal(t, false, ch.usePrivateEndpoint)
}
func TestConfigHandlerSetContext(t *testing.T) {
	// test set context when initialised properly
	ch := GetConfigurationHandlerInstance()
	ch.SetContext("c1", "dev", ContextOptions{
		BootstrapFile:           "flights.json",
		LiveConfigUpdateEnabled: false,
	})
	assert.Equal(t, "c1", ch.collectionID)
	assert.Equal(t, "dev", ch.environmentID)
	assert.Equal(t, "flights.json", ch.bootstrapFile)
	assert.Equal(t, false, ch.liveConfigUpdateEnabled)
}

func TestSaveCache(t *testing.T) {
	// test save feature when empty data is passed.
	ch := GetConfigurationHandlerInstance()
	data := `{"features":null,"properties":null,"segments":null}`
	ch.saveInCache([]byte(data))
	assert.Equal(t, 0, len(ch.cache.FeatureMap))
	assert.Equal(t, 0, len(ch.cache.PropertyMap))
	assert.Equal(t, 0, len(ch.cache.SegmentMap))

	// test save feature when non-empty data is passed.
	data = `{"features":[{"name":"Cycle Rentals8","feature_id":"cycle-rentals8","type":"BOOLEAN","enabled_value":true,"disabled_value":false,"segment_rules":[],"enabled":true, "rollout_percentage": 95}],"properties":[{"name":"p1","property_id":"p1","tags":"","type":"BOOLEAN","value":false,"segment_rules":[]}],"segments":[{"name":"beta-users","segment_id":"knliu818","rules":[{"values":["ibm.com"],"operator":"contains","attribute_name":"email"}]},{"name":"ibm employees","segment_id":"ka761hap","rules":[{"values":["ibm.com","in.ibm.com"],"operator":"endsWith","attribute_name":"email"}]}]}`
	ch.saveInCache([]byte(data))
	assert.Equal(t, 1, len(ch.cache.FeatureMap))
	assert.Equal(t, 1, len(ch.cache.PropertyMap))
	assert.Equal(t, 2, len(ch.cache.SegmentMap))
}

func TestFetchApi(t *testing.T) {

	// test fetch api when backend returns proper response
	// create a temp server which will act as our backend for the test
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(200)
			fmt.Fprintf(w, "%s", `{ "environments": [ { "name": "Dev", "environment_id": "dev", "description": "Environment created on instance creation", "tags": "", "color_code": "#FDD13A", "features": [ { "name": "Cycle Rentals", "feature_id": "cycle-rentals", "type": "BOOLEAN", "enabled_value": true, "disabled_value": false, "segment_rules": [], "enabled": true, "rollout_percentage": 95 } ], "properties": [ { "name": "Show Ad", "property_id": "show-ad", "tags": "", "type": "BOOLEAN", "value": false, "segment_rules": [] } ] } ], "segments": [ { "name": "beta-users", "segment_id": "knliu818", "rules": [ { "values": [ "ibm.com" ], "operator": "contains", "attribute_name": "email" } ] }, { "name": "ibm employees", "segment_id": "ka761hap", "rules": [ { "values": [ "ibm.com", "in.ibm.com" ], "operator": "endsWith", "attribute_name": "email" } ] } ] }`)
		}))

	ch := GetConfigurationHandlerInstance()
	ch.urlBuilder.SetBaseServiceURL(ts.URL)
	ch.urlBuilder.SetAuthenticator(&core.NoAuthAuthenticator{})
	ch.liveConfigUpdateEnabled = true
	ch.fetchFromAPI()
	assert.Equal(t, 1, len(ch.cache.FeatureMap))
	assert.Equal(t, 1, len(ch.cache.PropertyMap))
	assert.Equal(t, 2, len(ch.cache.SegmentMap))
	assert.Equal(t, "Cycle Rentals", ch.cache.FeatureMap["cycle-rentals"].Name)
	assert.Equal(t, "Show Ad", ch.cache.PropertyMap["show-ad"].Name)
	ts.Close()
	resetConfigurationHandler(ch)

	// test fetch api when backend returns 500 response
	// create a temp server which will act as our backend for the test
	ts1 := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(500)
		}))

	ch = GetConfigurationHandlerInstance()
	ch.urlBuilder.SetBaseServiceURL(ts1.URL)
	ch.urlBuilder.SetAuthenticator(&core.NoAuthAuthenticator{})
	ch.liveConfigUpdateEnabled = true
	ch.fetchFromAPI()
	assert.Equal(t, 0, len(ch.cache.FeatureMap))
	assert.Equal(t, 0, len(ch.cache.PropertyMap))
	assert.Equal(t, 0, len(ch.cache.SegmentMap))
	ts1.Close()
	resetConfigurationHandler(ch)

	// test fetch api when configuration handler instance is not initialized
	ch = GetConfigurationHandlerInstance()
	ch.isInitialized = false
	ch.fetchFromAPI()
	assert.Equal(t, 0, len(ch.cache.FeatureMap))
	assert.Equal(t, 0, len(ch.cache.PropertyMap))
	assert.Equal(t, 0, len(ch.cache.SegmentMap))
	resetConfigurationHandler(ch)
}

func TestUpdateCacheAndListener(t *testing.T) {
	mockLogger()
	// valid data but no listener method provided
	data := `{ "features": [ { "name": "Cycle Rentals", "feature_id": "cycle-rentals", "type": "BOOLEAN", "enabled_value": true, "disabled_value": false, "segment_rules": [], "enabled": true, "rollout_percentage": 95 } ], "properties": [ { "name": "Show Ad", "property_id": "show-ad", "tags": "", "type": "BOOLEAN", "value": false, "segment_rules": [] } ], "segments": [ { "name": "beta-users", "segment_id": "knliu818", "rules": [ { "values": [ "ibm.com" ], "operator": "contains", "attribute_name": "email" } ] }, { "name": "ibm employees", "segment_id": "ka761hap", "rules": [ { "values": [ "ibm.com", "in.ibm.com" ], "operator": "endsWith", "attribute_name": "email" } ] } ] }`
	ch := GetConfigurationHandlerInstance()
	ch.Init("us-south", "abc", "abc", false)
	ch.updateCacheAndListener([]byte(data))
	assert.Equal(t, 1, len(ch.cache.FeatureMap))
	assert.Equal(t, 1, len(ch.cache.PropertyMap))
	assert.Equal(t, 2, len(ch.cache.SegmentMap))
	assert.Equal(t, "Cycle Rentals", ch.cache.FeatureMap["cycle-rentals"].Name)
	assert.Equal(t, "Show Ad", ch.cache.PropertyMap["show-ad"].Name)
	resetConfigurationHandler(ch)

	// valid data and listener method provided
	ch = GetConfigurationHandlerInstance()
	ch.Init("us-south", "abc", "abc", false)
	msg := ""
	ch.configurationUpdateListener = func() {
		msg = "Latest evaluation done."
	}
	ch.updateCacheAndListener([]byte(data))
	assert.Equal(t, "Latest evaluation done.", msg)

	assert.Equal(t, 1, len(ch.cache.FeatureMap))
	assert.Equal(t, 1, len(ch.cache.PropertyMap))
	assert.Equal(t, 2, len(ch.cache.SegmentMap))
	assert.Equal(t, "Cycle Rentals", ch.cache.FeatureMap["cycle-rentals"].Name)
	assert.Equal(t, "Show Ad", ch.cache.PropertyMap["show-ad"].Name)
	resetConfigurationHandler(ch)

	// invalid data
	data = "<not a valid json>"
	ch = GetConfigurationHandlerInstance()
	ch.Init("us-south", "abc", "abc", false)
	ch.updateCacheAndListener([]byte(data))
	if hook.LastEntry().Message != "AppConfiguration - Error while unmarshalling JSON invalid character '<' looking for beginning of value" {
		t.Errorf("Test failed: Incorrect error message")
	}
	assert.Equal(t, 0, len(ch.cache.FeatureMap))
	assert.Equal(t, 0, len(ch.cache.PropertyMap))
	assert.Equal(t, 0, len(ch.cache.SegmentMap))
	resetConfigurationHandler(ch)

}

func TestRegisterConfigurationUpdateListener(t *testing.T) {
	mockLogger()
	// test register config update listener when config handler is initialized
	ch := GetConfigurationHandlerInstance()
	ch.Init("us-south", "abc", "abc", false)
	ch.isInitialized = true
	listenerBeforeRegisteration := ch.configurationUpdateListener
	var listener configurationUpdateListenerFunc = func() {
	}
	ch.registerConfigurationUpdateListener(listener)
	if reflect.ValueOf(listenerBeforeRegisteration).Pointer() == reflect.ValueOf(ch.configurationUpdateListener).Pointer() {
		t.Errorf("Test failed: configurationUpdateListenr not registered successfully.")
	}

	// test register config update listener when config handler is not initialized
	ch = GetConfigurationHandlerInstance()
	ch.Init("us-south", "abc", "abc", false)
	ch.isInitialized = false
	listenerBeforeRegisteration = ch.configurationUpdateListener
	listener = func() {
	}
	ch.registerConfigurationUpdateListener(listener)
	if hook.LastEntry().Message != "AppConfiguration - Invalid action. You can perform this action only after a successful initialization. Check the initialization section for errors." {
		t.Errorf("Test failed: Incorrect error message")
	}
	if reflect.ValueOf(listenerBeforeRegisteration).Pointer() != reflect.ValueOf(ch.configurationUpdateListener).Pointer() {
		t.Errorf("Test failed: configurationUpdateListenr shouldnt have registered since config handler is not initialized.")
	}

}

func TestStartWebSocket(t *testing.T) {

	// test start web socket when connection is done successfully
	mux := http.NewServeMux()
	mux.HandleFunc("/", wsEndpoint)

	server := httptest.NewServer(mux)
	server2 := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-type", "application/json")
		res.WriteHeader(200)
		fmt.Fprintf(res, "%s", `{ "environments": [ { "name": "Dev", "environment_id": "dev", "description": "Environment created on instance creation", "tags": "", "color_code": "#FDD13A", "features": [ { "name": "Cycle Rentals", "feature_id": "cycle-rentals", "type": "BOOLEAN", "enabled_value": true, "disabled_value": false, "segment_rules": [], "enabled": true, "rollout_percentage": 95 } ], "properties": [ { "name": "Show Ad", "property_id": "show-ad", "tags": "", "type": "BOOLEAN", "value": false, "segment_rules": [], } ] } ], "segments": [ { "name": "beta-users", "segment_id": "knliu818", "rules": [ { "values": [ "ibm.com" ], "operator": "contains", "attribute_name": "email" } ] }, { "name": "ibm employees", "segment_id": "ka761hap", "rules": [ { "values": [ "ibm.com", "in.ibm.com" ], "operator": "endsWith", "attribute_name": "email" } ] } ] }`)
	}))

	webSocketURL := "ws" + strings.TrimPrefix(server.URL, "http")

	defer server.Close()
	defer server2.Close()

	ch := GetConfigurationHandlerInstance()

	ch.urlBuilder = utils.GetInstance()

	ch.urlBuilder.SetBaseServiceURL(server2.URL)
	ch.urlBuilder.SetWebSocketURL(webSocketURL)
	ch.urlBuilder.SetAuthenticator(&core.NoAuthAuthenticator{})
	ch.guid = "guid"
	ch.collectionID = "collectionID"
	ch.liveConfigUpdateEnabled = true
	ch.isInitialized = true
	ch.startWebSocket()
	if hook.LastEntry().Message != "AppConfiguration - Socket connect failed. Could not generate Bearer token for the provided API Key." {
		t.Errorf("Test failed: Incorrect error message")
	}
	resetConfigurationHandler(ch)

}

func TestConfigHandlerGetProperty(t *testing.T) {
	// when property id exists in the cache
	ch := GetConfigurationHandlerInstance()
	data := `{"features":[{"name":"Cycle Rentals8","feature_id":"cycle-rentals8","type":"BOOLEAN","enabled_value":true,"disabled_value":false,"segment_rules":[],"enabled":true,"rollout_percentage":90}],"properties":[{"name":"ShowAd","property_id":"show-ad","tags":"","type":"BOOLEAN","value":false,"segment_rules":[]}],"segments":[{"name":"beta-users","segment_id":"knliu818","rules":[{"values":["ibm.com"],"operator":"contains","attribute_name":"email"}]},{"name":"ibm employees","segment_id":"ka761hap","rules":[{"values":["ibm.com","in.ibm.com"],"operator":"endsWith","attribute_name":"email"}]}]}`
	ch.saveInCache([]byte(data))
	val, _ := ch.getProperty("show-ad")
	assert.Equal(t, "ShowAd", val.Name)

	// when property id doesnt exists in the cache
	val, err := ch.getProperty("show-add")
	assert.Equal(t, "", val.Name)
	assert.Equal(t, "error : invalid property id show-add", fmt.Sprint(err))

	// when cache is empty
	data = `{"Features":null,"Properties":null,"Collection":{"name":"","collection_id":""},"Segments":null}`
	ch.saveInCache([]byte(data))
	val, err = ch.getProperty("show-ad")
	assert.Equal(t, "", val.Name)
	assert.Equal(t, "error : invalid property id show-ad", fmt.Sprint(err))

}

func TestConfigHandlerGetProperties(t *testing.T) {
	// when property id exists in the cache
	ch := GetConfigurationHandlerInstance()
	data := `{"features":[{"name":"Cycle Rentals8","feature_id":"cycle-rentals8","type":"BOOLEAN","enabled_value":true,"disabled_value":false,"segment_rules":[],"enabled":true,"rollout_percentage":85}],"properties":[{"name":"ShowAd","property_id":"show-ad","tags":"","type":"BOOLEAN","value":false,"segment_rules":[]}],"segments":[{"name":"beta-users","segment_id":"knliu818","rules":[{"values":["ibm.com"],"operator":"contains","attribute_name":"email"}]},{"name":"ibm employees","segment_id":"ka761hap","rules":[{"values":["ibm.com","in.ibm.com"],"operator":"endsWith","attribute_name":"email"}]}]}`
	ch.saveInCache([]byte(data))
	val, _ := ch.getProperties()
	assert.Equal(t, "ShowAd", val["show-ad"].Name)

	// when cache is
	ch.cache = nil
	val, _ = ch.getProperties()
	assert.Equal(t, 0, len(val))

}

func TestConfigHandlerGetSecret(t *testing.T) {
	// when property id exists in the cache
	ch := GetConfigurationHandlerInstance()
	data := `{"features":[{"name":"Cycle Rentals8","feature_id":"cycle-rentals8","type":"BOOLEAN","enabled_value":true,"disabled_value":false,"segment_rules":[],"enabled":true,"rollout_percentage":90}],"properties":[{"name":"ShowAd","property_id":"show-ad","tags":"","type":"SECRETREF","value":{"sm_instance_crn": "crn:v1:staging:public:secrets-manager:eu-gb:a/3268cfe9e25d4146a03b31f22f9a731a:d614a8ba-a13a-41cc-9e18-8fbff3ad9845::",
    "secret_type": "username_password"},"segment_rules":[]}],"segments":[{"name":"beta-users","segment_id":"knliu818","rules":[{"values":["ibm.com"],"operator":"contains","attribute_name":"email"}]},{"name":"ibm employees","segment_id":"ka761hap","rules":[{"values":["ibm.com","in.ibm.com"],"operator":"endsWith","attribute_name":"email"}]}]}`
	ch.saveInCache([]byte(data))
	_, err := ch.getProperty("show-ad1")
	if err == nil {
		t.Error("Expected getProperty to fail as the show-ad1 not found in the data")
	}
	// assert.Error(t, err, "Expected GetSecret to return error")

	//when the data type is invalid or different other than SECRETREF
	data = `{"features": [{"name": "Cycle Rentals8",
			"feature_id": "cycle-rentals8",
			"type": "BOOLEAN",
			"enabled_value": true,
			"disabled_value": false,
			"segment_rules": [],
			"enabled": true,
			"rollout_percentage": 90
		}],
		"properties": [{
			"name": "ShowAd2",
			"property_id": "show-ad2",
			"tags": "",
			"type": "BOOLEAN",
			"value": false,
			"segment_rules": []
		}],
		"segments": [{
			"name": "beta-users",
			"segment_id": "knliu818",
			"rules": [{
				"values": ["ibm.com"],
				"operator": "contains",
				"attribute_name": "email"
			}]
		}, {
			"name": "ibm employees",
			"segment_id": "ka761hap",
			"rules": [{
				"values": ["ibm.com", "in.ibm.com"],
				"operator": "endsWith",
				"attribute_name": "email"
			}]
		}]
	}`
	ch.saveInCache([]byte(data))
	_, err = ch.getSecret("show-ad2", nil)
	if err == nil {
		t.Error("Expected getProperty to fail as the show-ad2 property type is not SECRETREF")
	}

}

func TestConfigHandlerGetFeature(t *testing.T) {
	// when property id exists in the cache
	ch := GetConfigurationHandlerInstance()
	data := `{"features":[{"name":"Cycle Rentals8","feature_id":"cycle-rentals8","type":"BOOLEAN","enabled_value":true,"disabled_value":false,"segment_rules":[],"enabled":true,"rollout_percentage":70}],"properties":[{"name":"ShowAd","property_id":"show-ad","tags":"","type":"BOOLEAN","value":false,"segment_rules":[]}],"segments":[{"name":"beta-users","segment_id":"knliu818","rules":[{"values":["ibm.com"],"operator":"contains","attribute_name":"email"}]},{"name":"ibm employees","segment_id":"ka761hap","rules":[{"values":["ibm.com","in.ibm.com"],"operator":"endsWith","attribute_name":"email"}]}]}`
	ch.saveInCache([]byte(data))
	val, _ := ch.getFeature("cycle-rentals8")
	assert.Equal(t, "Cycle Rentals8", val.Name)

	// when property id doesnt exists in the cache
	val, err := ch.getFeature("cycle-rentals9")
	assert.Equal(t, "", val.Name)
	assert.Equal(t, "error : invalid feature id cycle-rentals9", fmt.Sprint(err))

	// when cache is empty
	data = `{"features":null,"properties":null,"segments":null}`
	ch.saveInCache([]byte(data))
	val, err = ch.getFeature("cycle-rentals8")
	assert.Equal(t, "", val.Name)
	assert.Equal(t, "error : invalid feature id cycle-rentals8", fmt.Sprint(err))

}
func TestConfigHandlerGetFeatures(t *testing.T) {
	// when property id exists in the cache
	ch := GetConfigurationHandlerInstance()
	data := `{"features":[{"name":"Cycle Rentals8","feature_id":"cycle-rentals8","type":"BOOLEAN","enabled_value":true,"disabled_value":false,"segment_rules":[],"enabled":true,"rollout_percentage":75}],"properties":[{"name":"ShowAd","property_id":"show-ad","tags":"","type":"BOOLEAN","value":false,"segment_rules":[]}],"segments":[{"name":"beta-users","segment_id":"knliu818","rules":[{"values":["ibm.com"],"operator":"contains","attribute_name":"email"}]},{"name":"ibm employees","segment_id":"ka761hap","rules":[{"values":["ibm.com","in.ibm.com"],"operator":"endsWith","attribute_name":"email"}]}]}`
	ch.saveInCache([]byte(data))
	val, _ := ch.getFeatures()
	assert.Equal(t, "Cycle Rentals8", val["cycle-rentals8"].Name)

	// when cache is nil
	ch.cache = nil
	val, _ = ch.getFeatures()
	assert.Equal(t, 0, len(val))

}
func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	if err := ws.WriteMessage(1, []byte("sending message")); err != nil {
		fmt.Println(err)
		return
	}

}
func resetConfigurationHandler(ch *ConfigurationHandler) {
	ch.cache = new(models.Cache)
}
