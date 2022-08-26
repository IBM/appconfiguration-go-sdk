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
	"testing"

	"github.com/IBM/appconfiguration-go-sdk/lib/internal/models"
	// "github.com/IBM/appconfiguration-go-sdk/lib/internal/utils/log"

	"github.com/stretchr/testify/assert"
)

// var testLogger, hook = test.NewNullLogger()

// func mockLogger() {
// 	log.SetLogger(testLogger)
// }

func TestInit(t *testing.T) {
	// test init when not initialised properly
	mockLogger()
	ac := GetInstance()
	ac.Init("", "", "")
	if hook.LastEntry().Message != "AppConfiguration - Provide a valid apiKey." {
		t.Errorf("Test failed: Incorrect error message")
	}
	reset(ac)

	// test init when initialised properly
	assert.Nil(t, ac.configurationHandlerInstance)
	ac.Init("a", "b", "c")
	assert.NotNil(t, ac.configurationHandlerInstance)

}

func TestSetContext(t *testing.T) {
	// test set context when is ac is not initialized properly
	mockLogger()
	ac := GetInstance()
	ac.isInitialized = false
	ac.SetContext("c1", "dev")
	if hook.LastEntry().Message != "AppConfiguration - Invalid action. You can perform this action only after a successful initialization. Check the initialization section for errors." {
		t.Errorf("Test failed: Incorrect error message")
	}
	reset(ac)
	// when no collection id is provided
	ac.isInitialized = true
	ac.SetContext("", "dev")
	if hook.LastEntry().Message != "AppConfiguration - Provide a valid collectionId." {
		t.Errorf("Test failed: Incorrect error message")
	}
	reset(ac)
	// when no environment id is provided
	ac.isInitialized = true
	ac.SetContext("c1", "")
	if hook.LastEntry().Message != "AppConfiguration - Provide a valid environmentId." {
		t.Errorf("Test failed: Incorrect error message")
	}
	reset(ac)

	// when only collection id and environment id is provided (in-memory cache)
	//ac.Init("a", "b", "c")
	//ac.isInitialized = true
	//assert.Equal(t, false, ac.isInitializedConfig)
	//ac.SetContext("c1", "dev")
	//assert.Equal(t, true, ac.isInitializedConfig)
	//reset(ac)

	// when collection id and environment id is provided and the number of context options is more than 1
	ac.Init("a", "b", "c")
	ac.isInitialized = true
	ac.SetContext("c1", "dev", ContextOptions{
		BootstrapFile:           "saflights/flights.json",
		LiveConfigUpdateEnabled: false,
	}, ContextOptions{
		BootstrapFile:           "saflights/flights.json",
		LiveConfigUpdateEnabled: false,
	})
	if hook.LastEntry().Message != "AppConfiguration - Incorrect usage of context options. At most of one ContextOptions struct should be passed." {
		t.Errorf("Test failed: Incorrect error message")
	}
	reset(ac)

	// when collection id and environment id is provided and the number of context options is 1. (Bootstrap file evaluation)
	ac.Init("a", "b", "c")
	ac.isInitialized = true
	assert.Equal(t, false, ac.isInitializedConfig)
	ac.SetContext("c1", "dev", ContextOptions{
		BootstrapFile:           "saflights/flights.json",
		LiveConfigUpdateEnabled: false,
	})
	assert.Equal(t, true, ac.isInitializedConfig)
	reset(ac)

	// when collection id and environment id is provided and the context options has no config file inspite of live update enabled set to false
	ac.Init("a", "b", "c")
	ac.isInitialized = true
	assert.Equal(t, false, ac.isInitializedConfig)
	ac.SetContext("c1", "dev", ContextOptions{
		BootstrapFile:           "",
		LiveConfigUpdateEnabled: false,
	})
	if hook.LastEntry().Message != "AppConfiguration - Provide bootstrap_file value when live_config_update_enabled is false." {
		t.Errorf("Test failed: Incorrect error message")
	}
	reset(ac)
}
func TestGetFeature(t *testing.T) {
	// test get feature when not initialised properly
	ac := GetInstance()
	_, err := ac.GetFeature("FID1")
	assert.Error(t, err, "Expected GetFeature to return error")
	reset(ac)

	// test get feature when config has been initialized properly and feature exists in the cache
	mockInit(ac)
	mockSetCache(ac)
	feature, err := ac.GetFeature("FID1")
	if assert.NotNil(t, feature) {
		assert.Equal(t, "discountOnBikes", feature.Name)
	}
	reset(ac)
}

func TestGetFeatures(t *testing.T) {
	// test get features when not initialised properly
	ac := GetInstance()
	_, err := ac.GetFeatures()
	assert.Error(t, err, "Expected GetFeatures to return error")
	reset(ac)

	// test get features when config has been initialized properly and feature exists in the cache
	mockInit(ac)
	mockSetCache(ac)
	features, err := ac.GetFeatures()
	if assert.NotNil(t, features) {
		assert.Equal(t, "discountOnBikes", features["FID1"].Name)
	}
	reset(ac)
}
func TestGetProperty(t *testing.T) {
	// test get feature when not initialised properly
	ac := GetInstance()
	_, err := ac.GetProperty("PID1")
	assert.Error(t, err, "Expected GetFeature to return error")
	reset(ac)

	// test get feature when config has been initialized properly and feature exists in the cache
	mockInit(ac)
	mockSetCache(ac)
	property, err := ac.GetProperty("PID1")
	if assert.NotNil(t, property) {
		assert.Equal(t, "nodeReplica", property.Name)
	}
	reset(ac)
}

func TestGetSecret(t *testing.T) {
	// test get secret property when not initialised properly
	ac := GetInstance()
	_, err := ac.GetSecret("propertySecretId", nil)
	if err == nil {
		t.Error("Expected GetSecret to return error when initialization of app config and cache has failed")
	}
	// assert.Error(t, err, "Expected GetSecret to return error when initialization of app config and cache has failed")
	reset(ac)
	mockInit(ac)
	mockSetCache(ac)
	_, err = ac.GetSecret("propertySecretId", nil)
	if err == nil {
		t.Error("Expected GetSecret to return error when secretManagerObject is failed")
	}
	reset(ac)
	// assert.Error(t, err, "Expected GetSecret to return error when secretManagerObject is failed")
}

func TestGetProperties(t *testing.T) {
	// test get properties when not initialised properly
	ac := GetInstance()
	_, err := ac.GetProperties()
	assert.Error(t, err, "Expected GetProperties to return error")
	reset(ac)
	// test get properties when config has been initialized properly and property exists in the cache
	mockInit(ac)
	mockSetCache(ac)
	properties, err := ac.GetProperties()
	if assert.NotNil(t, properties) {
		assert.Equal(t, "nodeReplica", properties["PID1"].Name)
	}
	reset(ac)
}

func TestFetchConfigurations(t *testing.T) {
	// test fetch configurations when sdk is not initialised properly
	ac := GetInstance()
	ac.FetchConfigurations()
	if hook.LastEntry().Message != "AppConfiguration - Invalid action. You can perform this action only after a successful initialization and setting the context. Check the Init and SetContext section for errors." {
		t.Errorf("Test failed: Incorrect error message")
	}
}

func TestRegisterConfigurationsUpdateListener(t *testing.T) {
	// test TestRegisterConfigurationsUpdateListener when sdk is not initialised properly
	ac := GetInstance()
	ac.RegisterConfigurationUpdateListener(func() {})
	if hook.LastEntry().Message != "AppConfiguration - Invalid action. You can perform this action only after a successful initialization and setting the context. Check the Init and SetContext section for errors." {
		t.Errorf("Test failed: Incorrect error message")
	}
}

func reset(ac *AppConfiguration) {
	ac.isInitializedConfig = false
	ac.configurationHandlerInstance = nil
}
func mockInit(ac *AppConfiguration) {
	ac.isInitializedConfig = true
	ac.configurationHandlerInstance = new(ConfigurationHandler)
}
func mockSetCache(ac *AppConfiguration) {
	featureMap := make(map[string]models.Feature)
	propertyMap := make(map[string]models.Property)

	var testFeature models.Feature
	testFeature.Name = "discountOnBikes"
	testFeature.FeatureID = "FID1"
	featureMap["FID1"] = testFeature
	var cacheInstance *models.Cache
	cacheInstance = new(models.Cache)
	cacheInstance.FeatureMap = featureMap
	ac.configurationHandlerInstance.cache = cacheInstance

	var testProperty models.Property
	testProperty.Name = "nodeReplica"
	testProperty.PropertyID = "PID1"
	propertyMap["PID1"] = testProperty
	cacheInstance.PropertyMap = propertyMap
	ac.configurationHandlerInstance.cache = cacheInstance
}
