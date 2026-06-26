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

const guid, env, col, ent, seg, feat, prop = "guid", "dev", "c1", "e1", "s1", "f1", "p1"

func TestAddMetering(t *testing.T) {
	// test add metering when the meteringFeatureData is empty and first recording of the evaluation is done
	m := GetMeteringInstance()
	m.Init(guid, env, col)
	assert.Equal(t, 0, len(m.meteringFeatureData))

	m.addMetering(ent, seg, feat, prop)
	assert.Equal(t, 1, len(m.meteringFeatureData))
	record := m.meteringFeatureData[buildCompositeKey(feat, ent, seg)]
	assert.Equal(t, int64(1), record.getCount())

	// when the evaluation is done for the second time for the same feature against the same entity and segment

	m.addMetering(ent, seg, feat, prop)
	record = m.meteringFeatureData[buildCompositeKey(feat, ent, seg)]
	assert.Equal(t, int64(2), record.getCount())

	// when the evaluation is done  for the same feature against the same entity but different segment

	m.addMetering(ent, "s2", feat, prop)
	record = m.meteringFeatureData[buildCompositeKey(feat, ent, "s2")]
	assert.Equal(t, int64(1), record.getCount())

	// when the evaluation is done  for the same feature against but different entity

	m.addMetering("e2", seg, feat, prop)
	record = m.meteringFeatureData[buildCompositeKey(feat, "e2", seg)]
	assert.Equal(t, int64(1), record.getCount())

	// when the evaluation is done  for different feature but same collection

	m.addMetering(ent, seg, "f2", prop)
	record = m.meteringFeatureData[buildCompositeKey("f2", ent, seg)]
	assert.Equal(t, int64(1), record.getCount())

	// when the evaluation is done  for different collection but same environment
	// Note: With simplified keys, collection changes don't affect the key

	m.addMetering("e2", seg, "f2", prop)
	record = m.meteringFeatureData[buildCompositeKey("f2", "e2", seg)]
	assert.Equal(t, int64(1), record.getCount())

	// when the evaluation is done  for different environment but same guid
	// Note: With simplified keys, environment changes don't affect the key

	m.addMetering("e2", seg, "f2", prop)
	record = m.meteringFeatureData[buildCompositeKey("f2", "e2", seg)]
	assert.Equal(t, int64(2), record.getCount()) // Should increment existing record

	resetMeteringInstance()
}

func TestBuildRequestBody(t *testing.T) {
	// when request body contains only features evaluations
	m := GetMeteringInstance()
	m.Init(guid, env, col)
	assert.Equal(t, 0, len(m.meteringFeatureData))
	m.addMetering(ent, seg, feat, prop)
	m.addMetering(ent, seg, feat, prop)

	assert.Equal(t, 1, len(m.meteringFeatureData))
	record := m.meteringFeatureData[buildCompositeKey(feat, ent, seg)]
	assert.Equal(t, int64(2), record.getCount())
	collectionsUsages := CollectionUsages{}
	assert.Equal(t, 0, len(collectionsUsages.Usages))

	m.buildRequestBody(m.meteringFeatureData, &collectionsUsages, "feature_id")
	assert.Equal(t, int64(2), collectionsUsages.Usages[0].Count)
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
	m.addMetering(ent, seg, feat, prop)
	m.addMetering(ent, seg, feat, prop)

	assert.Equal(t, 1, len(m.meteringFeatureData))
	record := m.meteringFeatureData[buildCompositeKey(feat, ent, seg)]
	assert.Equal(t, int64(2), record.getCount())
	collectionsUsages := CollectionUsages{}
	assert.Equal(t, 0, len(collectionsUsages.Usages))

	m.buildRequestBody(m.meteringFeatureData, &collectionsUsages, "feature_id")
	assert.Equal(t, int64(2), collectionsUsages.Usages[0].Count)
	m.sendToServer(collectionsUsages)
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
	m.sendToServer(collectionsUsages)
	if hook.LastEntry().Message != "AppConfiguration - Error while sending metering data to server. Internal Server Error" {
		t.Errorf("Test failed: Incorrect error message -->")
	}
	resetMeteringInstance()

}
func resetMeteringInstance() {
	meteringInstance = nil
	urlBuilderInstance = nil
	log.SetLogLevel("info")
}
