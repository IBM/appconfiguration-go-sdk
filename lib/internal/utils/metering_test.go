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

package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IBM/go-sdk-core/v5/core"

	"github.com/IBM/appconfiguration-go-sdk/lib/internal/utils/log"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

var testLogger, hook = test.NewNullLogger()

func mockLogger() {
	log.SetLogger(testLogger)
}

func TestMeteringInit(t *testing.T) {
	// test init
	m := GetMeteringInstance()
	assert.Equal(t, "", m.guid)
	assert.Equal(t, "", m.CollectionID)
	assert.Equal(t, "", m.EnvironmentID)
	m.Init("guid", "dev", "c1")
	assert.Equal(t, "guid", m.guid)
	assert.Equal(t, "c1", m.CollectionID)
	assert.Equal(t, "dev", m.EnvironmentID)
	resetMeteringInstance()

}

func TestAddMetering(t *testing.T) {
	// test add metering when the meteringFeatureData is empty and first recording of the evaluation is done
	m := GetMeteringInstance()
	m.Init("guid", "dev", "c1")
	assert.Equal(t, 0, len(m.meteringFeatureData))
	m.addMetering("guid", "dev", "c1", "e1", "s1", "f1", "p1")
	assert.Equal(t, 1, len(m.meteringFeatureData))
	guidVal := m.meteringFeatureData["guid"]

	envtVal := guidVal["dev"]

	collectionVal := envtVal["c1"]

	featureVal := collectionVal["f1"]

	entityVal := featureVal["e1"]

	segmentVal := entityVal["s1"]

	assert.Equal(t, int64(1), segmentVal.count)

	// when the evaluation is done for the second time for the same feature against the same entity and segment

	m.addMetering("guid", "dev", "c1", "e1", "s1", "f1", "p1")

	guidVal = m.meteringFeatureData["guid"]

	envtVal = guidVal["dev"]

	collectionVal = envtVal["c1"]

	featureVal = collectionVal["f1"]

	entityVal = featureVal["e1"]

	segmentVal = entityVal["s1"]
	assert.Equal(t, int64(2), segmentVal.count)

	// when the evaluation is done  for the same feature against the same entity but different segment

	m.addMetering("guid", "dev", "c1", "e1", "s2", "f1", "p1")

	guidVal = m.meteringFeatureData["guid"]

	envtVal = guidVal["dev"]

	collectionVal = envtVal["c1"]

	featureVal = collectionVal["f1"]

	entityVal = featureVal["e1"]

	segmentVal = entityVal["s2"]
	assert.Equal(t, int64(1), segmentVal.count)

	// when the evaluation is done  for the same feature against but different entity

	m.addMetering("guid", "dev", "c1", "e2", "s1", "f1", "p1")

	guidVal = m.meteringFeatureData["guid"]

	envtVal = guidVal["dev"]

	collectionVal = envtVal["c1"]

	featureVal = collectionVal["f1"]

	entityVal = featureVal["e2"]

	segmentVal = entityVal["s1"]
	assert.Equal(t, int64(1), segmentVal.count)

	// when the evaluation is done  for different feature but same collection

	m.addMetering("guid", "dev", "c1", "e2", "s1", "f2", "p1")

	guidVal = m.meteringFeatureData["guid"]

	envtVal = guidVal["dev"]

	collectionVal = envtVal["c1"]

	featureVal = collectionVal["f2"]

	entityVal = featureVal["e2"]

	segmentVal = entityVal["s1"]
	assert.Equal(t, int64(1), segmentVal.count)

	// when the evaluation is done  for different collection but same environment

	m.addMetering("guid", "dev", "c2", "e2", "s1", "f2", "p1")

	guidVal = m.meteringFeatureData["guid"]

	envtVal = guidVal["dev"]

	collectionVal = envtVal["c2"]

	featureVal = collectionVal["f2"]

	entityVal = featureVal["e2"]

	segmentVal = entityVal["s1"]
	assert.Equal(t, int64(1), segmentVal.count)

	// when the evaluation is done  for different environment but same guid

	m.addMetering("guid", "prod", "c2", "e2", "s1", "f2", "p1")

	guidVal = m.meteringFeatureData["guid"]

	envtVal = guidVal["prod"]

	collectionVal = envtVal["c2"]

	featureVal = collectionVal["f2"]

	entityVal = featureVal["e2"]

	segmentVal = entityVal["s1"]
	assert.Equal(t, int64(1), segmentVal.count)

	resetMeteringInstance()
}

func TestBuildRequestBody(t *testing.T) {
	// when request body contains only features evaluations
	m := GetMeteringInstance()
	m.Init("guid", "dev", "c1")
	assert.Equal(t, 0, len(m.meteringFeatureData))
	m.addMetering("guid", "dev", "c1", "e1", "s1", "f1", "p1")
	m.addMetering("guid", "dev", "c1", "e1", "s1", "f1", "p1")

	assert.Equal(t, 1, len(m.meteringFeatureData))
	guidVal := m.meteringFeatureData["guid"]

	envtVal := guidVal["dev"]

	collectionVal := envtVal["c1"]

	featureVal := collectionVal["f1"]

	entityVal := featureVal["e1"]

	segmentVal := entityVal["s1"]

	assert.Equal(t, int64(2), segmentVal.count)
	guidMap := make(map[string][]CollectionUsages)
	assert.Equal(t, 0, len(guidMap))

	m.buildRequestBody(m.meteringFeatureData, guidMap, "feature_id")
	assert.Equal(t, int64(2), guidMap["guid"][0].Usages[0].Count)
	resetMeteringInstance()

}

func TestSendToServer(t *testing.T) {

	// test send to server with backend returning success

	mockLogger()
	log.SetLogLevel("debug")
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(202)
			fmt.Fprintf(w, "%s", `Success`)
		}))

	m := GetMeteringInstance()
	m.Init("guid", "dev", "c1")
	urlBuilderInstance = &URLBuilder{

		httpBase: ts.URL,
	}
	urlBuilderInstance.SetAuthenticator(&core.NoAuthAuthenticator{})

	assert.Equal(t, 0, len(m.meteringFeatureData))
	m.addMetering("guid", "dev", "c1", "e1", "s1", "f1", "p1")
	m.addMetering("guid", "dev", "c1", "e1", "s1", "f1", "p1")

	assert.Equal(t, 1, len(m.meteringFeatureData))
	guidVal := m.meteringFeatureData["guid"]

	envtVal := guidVal["dev"]

	collectionVal := envtVal["c1"]

	featureVal := collectionVal["f1"]

	entityVal := featureVal["e1"]

	segmentVal := entityVal["s1"]

	assert.Equal(t, int64(2), segmentVal.count)
	guidMap := make(map[string][]CollectionUsages)
	assert.Equal(t, 0, len(guidMap))

	m.buildRequestBody(m.meteringFeatureData, guidMap, "feature_id")
	assert.Equal(t, int64(2), guidMap["guid"][0].Usages[0].Count)
	m.sendToServer("guid", guidMap["guid"][0])
	if hook.LastEntry().Message != "AppConfiguration - Successfully sent metering data to server." {
		t.Errorf("Test failed: Incorrect error message")
	}
	ts.Close()
	resetMeteringInstance()

	// test send to server with backend returning failure

	mockLogger()
	log.SetLogLevel("debug")
	ts = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)

		}))
	urlBuilderInstance = &URLBuilder{

		httpBase: ts.URL,
	}
	urlBuilderInstance.SetAuthenticator(&core.NoAuthAuthenticator{})
	m.sendToServer("guid", guidMap["guid"][0])
	if hook.LastEntry().Message != "AppConfiguration - Error while sending metering data to server <nil>" {
		t.Errorf("Test failed: Incorrect error message -->")
	}
	resetMeteringInstance()

}
func resetMeteringInstance() {
	meteringInstance = nil
	urlBuilderInstance = nil
	log.SetLogLevel("info")
}
