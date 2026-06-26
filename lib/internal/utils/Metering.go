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
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/IBM/appconfiguration-go-sdk/lib/internal/constants"
	messages "github.com/IBM/appconfiguration-go-sdk/lib/internal/messages"
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/utils/log"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/robfig/cron"
)

// Usages : Usages struct
type Usages struct {
	FeatureID      string      `json:"feature_id,omitempty"`
	PropertyID     string      `json:"property_id,omitempty"`
	EntityID       interface{} `json:"entity_id"`
	SegmentID      interface{} `json:"segment_id"`
	EvaluationTime string      `json:"evaluation_time"`
	Count          int64       `json:"count"`
}

// CollectionUsages : CollectionUsages struct
type CollectionUsages struct {
	CollectionID  string   `json:"collection_id"`
	EnvironmentID string   `json:"environment_id"`
	Usages        []Usages `json:"usages"`
}

type meteringRecord struct {
	count      int64
	latestTime atomic.Value
}

func newMeteringRecord(timestamp string) *meteringRecord {
	record := &meteringRecord{count: 1}
	record.latestTime.Store(timestamp)
	return record
}

func (mr *meteringRecord) increment(newTime string) {
	atomic.AddInt64(&mr.count, 1)
	for {
		currentTime := mr.latestTime.Load().(string)
		if newTime > currentTime {
			if mr.latestTime.CompareAndSwap(currentTime, newTime) {
				break
			}
		} else {
			break
		}
	}
}

func (mr *meteringRecord) getCount() int64 {
	return atomic.LoadInt64(&mr.count)
}

func (mr *meteringRecord) getLatestTime() string {
	return mr.latestTime.Load().(string)
}

type Metering struct {
	CollectionID         string
	EnvironmentID        string
	guid                 string
	mu                   sync.Mutex
	meteringFeatureData  map[string]*meteringRecord
	meteringPropertyData map[string]*meteringRecord
}

// Delimiter for building key by joining strings
const delimiter = "\u001F"

// SendInterval : SendInterval struct
const SendInterval = "10m"

var meteringInstance *Metering

func GetMeteringInstance() *Metering {
	log.Debug(messages.RetrieveMeteringInstance)
	if meteringInstance == nil {
		meteringInstance = &Metering{}
		meteringInstance.meteringFeatureData = make(map[string]*meteringRecord)
		meteringInstance.meteringPropertyData = make(map[string]*meteringRecord)
		log.Debug(messages.StartSendingMeteringData)
		c := cron.New()
		c.AddFunc("@every "+SendInterval, meteringInstance.sendMetering)
		c.Start()
	}
	return meteringInstance
}

func (mt *Metering) Init(guid string, environmentID string, collectionID string) {
	mt.guid = guid
	mt.EnvironmentID = environmentID
	mt.CollectionID = collectionID
}

func (mt *Metering) addMetering(entityID string, segmentID string, featureID string, propertyID string) {
	log.Debug(messages.AddMetering)
	defer GracefullyHandleError()

	t := time.Now().UTC()
	formattedTime := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	var meteringData map[string]*meteringRecord
	var modifyKey string
	if featureID != "" {
		meteringData = meteringInstance.meteringFeatureData
		modifyKey = featureID
	} else {
		meteringData = meteringInstance.meteringPropertyData
		modifyKey = propertyID
	}

	key := buildCompositeKey(modifyKey, entityID, segmentID)
	
	mt.mu.Lock()
	record, exists := meteringData[key]
	if exists {
		mt.mu.Unlock()
		record.increment(formattedTime)
	} else {
		meteringData[key] = newMeteringRecord(formattedTime)
		mt.mu.Unlock()
	}
}

func (mt *Metering) RecordEvaluation(featureID string, propertyID string, entityID string, segmentID string) {
	log.Debug(messages.RecordEval)
	mt.addMetering(entityID, segmentID, featureID, propertyID)
}

func (mt *Metering) buildRequestBody(sendMeteringData map[string]*meteringRecord, collectionUsages *CollectionUsages, resourceType string) {
	for key, record := range sendMeteringData {
		keyParts := parseCompositeKey(key)
		if len(keyParts) != 3 {
			log.Error("Invalid composite key format: expected 3 parts, got ", len(keyParts))
			continue
		}
		var entityId interface{} = nil
		if keyParts[1] != "" && keyParts[1] != constants.DefaultEntityID {
			entityId = keyParts[1]
		}
		var segmentId interface{} = nil
		if keyParts[2] != "" && keyParts[2] != constants.DefaultSegmentID {
			segmentId = keyParts[2]
		}
		usages := Usages{
			EntityID:       entityId,
			SegmentID:      segmentId,
			EvaluationTime: record.getLatestTime(),
			Count:          record.getCount(),
		}
		if resourceType == "feature_id" {
			usages.FeatureID = keyParts[0]
		} else {
			usages.PropertyID = keyParts[0]
		}
		collectionUsages.Usages = append(collectionUsages.Usages, usages)
	}
}

func (mt *Metering) sendMetering() {
	log.Debug(messages.TenMinExpiry)
	defer GracefullyHandleError()

	mt.mu.Lock()
	currentFeatureData := mt.meteringFeatureData
	currentPropertyData := mt.meteringPropertyData
	mt.meteringFeatureData = make(map[string]*meteringRecord)
	mt.meteringPropertyData = make(map[string]*meteringRecord)
	mt.mu.Unlock()

	log.Debug(currentFeatureData)
	log.Debug(currentPropertyData)
	
	if len(currentFeatureData) == 0 && len(currentPropertyData) == 0 {
		return
	}

	collectionUsages := CollectionUsages{
		CollectionID:  mt.CollectionID,
		EnvironmentID: mt.EnvironmentID,
		Usages:        []Usages{},
	}
	
	if len(currentFeatureData) > 0 {
		mt.buildRequestBody(currentFeatureData, &collectionUsages, "feature_id")
	}
	if len(currentPropertyData) > 0 {
		mt.buildRequestBody(currentPropertyData, &collectionUsages, "property_id")
	}

	count := len(collectionUsages.Usages)
	if count > constants.DefaultUsageLimit {
		mt.sendSplitMetering(collectionUsages, count)
	} else {
		mt.sendToServer(collectionUsages)
	}
}
func (mt *Metering) sendSplitMetering(collectionUsages CollectionUsages, count int) {
	var lim int = 0
	subUsages := collectionUsages.Usages
	for lim < count {
		var endIndex int
		if lim+constants.DefaultUsageLimit >= count {
			endIndex = count
		} else {
			endIndex = lim + constants.DefaultUsageLimit
		}
		var collectionUsageElem CollectionUsages
		collectionUsageElem.CollectionID = collectionUsages.CollectionID
		collectionUsageElem.EnvironmentID = collectionUsages.EnvironmentID
		for i := lim; i < endIndex; i++ {
			collectionUsageElem.Usages = append(collectionUsageElem.Usages, subUsages[i])
		}
		mt.sendToServer(collectionUsageElem)
		lim = lim + constants.DefaultUsageLimit
	}
}
func (mt *Metering) sendToServer(collectionUsages CollectionUsages) {
	log.Debug(messages.SendMeteringServer)
	log.Debug(collectionUsages)
	builder := core.NewRequestBuilder(core.POST)
	pathParamsMap := map[string]string{
		"guid": mt.guid,
	}
	_, err := builder.ResolveRequestURL(urlBuilderInstance.GetBaseServiceURL(), `/apprapp/events/v1/instances/{guid}/usage`, pathParamsMap)
	if err != nil {
		return
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	builder.AddHeader("User-Agent", constants.UserAgent)
	_, err = builder.SetBodyContentJSON(collectionUsages)
	if err != nil {
		return
	}
	response, err := GetAPIManagerInstance().Request(builder)
	if response != nil && response.StatusCode == 202 {
		log.Debug(messages.SendMeteringSuccess)
	} else {
		// [first] Log the accurate reason
		if err != nil {
			log.Error(messages.SendMeteringServerErr + err.Error())
		} else {
			log.Error(messages.SendMeteringServerErr)
		}
		statusCode := -1
		if response != nil {
			statusCode = response.StatusCode
		}
		// status code -1 represents response has been failed to form else check the retry cases
		if statusCode == -1 || statusCode == 429 || (statusCode >= 500 && statusCode <= 599) {
			// schedule a function to send the same payload after 10 minutes
			minutes, _ := time.ParseDuration(SendInterval)
			time.AfterFunc(time.Second*time.Duration(minutes.Seconds()), func() {
				mt.sendToServer(collectionUsages)
			})
		}
	}
}

func buildCompositeKey(modifyKey, entityID, segmentID string) string {
	return strings.Join([]string{modifyKey, entityID, segmentID}, delimiter)
}

func parseCompositeKey(compositeKey string) []string {
	return strings.SplitN(compositeKey, delimiter, 3)
}
