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
	"sync"
	"testing"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/IBM/appconfiguration-go-sdk/lib/internal/utils/log"
)

var testLogger, hook = test.NewNullLogger()

func mockLogger() {
	log.SetLogger(testLogger)
}

func TestMeteringInit(t *testing.T) {
	resetMeteringInstance()
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
func TestMeteringSingletonConcurrent(t *testing.T) {
	resetMeteringInstance()
	defer resetMeteringInstance()

	const numGoroutines = 50
	instances := make([]*Metering, numGoroutines)
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			instances[idx] = GetMeteringInstance()
		}(i)
	}
	wg.Wait()

	first := instances[0]
	assert.NotNil(t, first)
	for i := 1; i < numGoroutines; i++ {
		assert.Same(t, first, instances[i], "all goroutines should receive the same Metering instance")
	}
}

func TestMeteringConcurrentEvaluationAndRotation(t *testing.T) {
	resetMeteringInstance()
	defer resetMeteringInstance()

	m := GetMeteringInstance()
	m.Init("test-guid", "test-env", "test-col")

	const numWriters = 10
	const numIterations = 200
	var wg sync.WaitGroup
	wg.Add(numWriters)
	stop := make(chan struct{})

	var flusherWg sync.WaitGroup
	flusherWg.Add(1)
	go func() {
		defer flusherWg.Done()
		for {
			select {
			case <-stop:
				return
			case <-time.After(5 * time.Millisecond):
				m.sendMetering()
			}
		}
	}()

	for w := 0; w < numWriters; w++ {
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < numIterations; i++ {
				fID := fmt.Sprintf("feat-%d", i%5)
				pID := fmt.Sprintf("prop-%d", i%5)
				eID := fmt.Sprintf("entity-%d", workerID)
				sID := fmt.Sprintf("seg-%d", i%2)
				m.RecordEvaluation(fID, "", eID, sID)
				m.RecordEvaluation("", pID, eID, sID)
			}
		}(w)
	}

	wg.Wait()
	close(stop)
	flusherWg.Wait()
}

func TestMeteringSendMeteringSnapshot(t *testing.T) {
	resetMeteringInstance()
	defer resetMeteringInstance()

	m := GetMeteringInstance()
	m.Init("guid-1", "env-1", "col-1")
	m.RecordEvaluation("f1", "", "e1", "s1")
	assert.Equal(t, 1, len(m.meteringFeatureData))
	m.sendMetering()
	assert.Equal(t, 0, len(m.meteringFeatureData))
	assert.Equal(t, 0, len(m.meteringPropertyData))
}

// covered above, but testing again aggressively
func TestMeteringSingletonAndMapRotationRace(t *testing.T) {
	resetMeteringInstance()
	defer resetMeteringInstance()

	stop := make(chan struct{})
	var wg sync.WaitGroup
	for w := 0; w < 8; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mt := GetMeteringInstance()
			mt.Init("guid", "env", "collection")
			for {
				select {
				case <-stop:
					return
				default:
				}
				mt.RecordEvaluation("f1", "", "entity-1", "")
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		mt := GetMeteringInstance()
		for {
			select {
			case <-stop:
				return
			default:
			}
			mt.sendMetering()
		}
	}()

	time.Sleep(200 * time.Millisecond)
	close(stop)
	wg.Wait()
}

// covered above, but testing again aggressively
func TestConcurrentMapIterationCrash(t *testing.T) {
	resetMeteringInstance()
	defer resetMeteringInstance()

	mt := GetMeteringInstance()
	mt.Init("guid", "env", "collection")

	for i := 0; i < 50000; i++ {
		mt.RecordEvaluation(fmt.Sprintf("f%d", i), "", "entity-1", "")
	}

	stop := make(chan struct{})
	var wg sync.WaitGroup

	for w := 0; w < 16; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			i := 0
			for {
				select {
				case <-stop:
					return
				default:
				}
				i++
				mt.RecordEvaluation(fmt.Sprintf("new-f%d-%d", workerID, i), "", "entity-1", "")
			}
		}(w)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 500; i++ {
			select {
			case <-stop:
				return
			default:
			}
			mt.sendMetering()
		}
	}()

	time.Sleep(2 * time.Second)
	close(stop)
	wg.Wait()
	t.Log("test completed without the runtime firing 'concurrent map iteration and map write'")
}

func resetMeteringInstance() {
	ResetMeteringInstance()
	ResetURLBuilderInstance()
	log.SetLogLevel("info")
}
