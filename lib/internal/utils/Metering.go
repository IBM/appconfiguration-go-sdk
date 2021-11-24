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
	"sync"
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

type featureMetric struct {
	count          int64
	evaluationTime string
}

// Metering : Metering struct
type Metering struct {
	CollectionID         string
	EnvironmentID        string
	guid                 string
	mu                   sync.Mutex
	meteringFeatureData  map[string]map[string]map[string]map[string]map[string]map[string]featureMetric //guid->EnvironmentID->CollectionID->featureId->entityId->segmentId
	meteringPropertyData map[string]map[string]map[string]map[string]map[string]map[string]featureMetric //guid->EnvironmentID->CollectionID->propertyId->entityId->segmentId
}

// SendInterval : SendInterval struct
const SendInterval = "10m"

var meteringInstance *Metering

// GetMeteringInstance : Get Metering Instance
func GetMeteringInstance() *Metering {
	log.Debug(messages.RetrieveMeteringInstance)
	if meteringInstance == nil {
		meteringInstance = &Metering{}
		guidFeatureMap := make(map[string]map[string]map[string]map[string]map[string]map[string]featureMetric)
		guidPropertyMap := make(map[string]map[string]map[string]map[string]map[string]map[string]featureMetric)
		meteringInstance.meteringFeatureData = guidFeatureMap
		meteringInstance.meteringPropertyData = guidPropertyMap
		// start sending metering data in the background
		log.Debug(messages.StartSendingMeteringData)
		c := cron.New()
		c.AddFunc("@every "+SendInterval, meteringInstance.sendMetering)
		c.Start()

	}
	return meteringInstance
}

// Init : Init
func (mt *Metering) Init(guid string, environmentID string, collectionID string) {
	mt.guid = guid
	mt.EnvironmentID = environmentID
	mt.CollectionID = collectionID
}

func (mt *Metering) addMetering(guid string, environmentID string, collectionID string, entityID string, segmentID string, featureID string, propertyID string) {
	log.Debug(messages.AddMetering)
	defer GracefullyHandleError()
	mt.mu.Lock()
	t := time.Now().UTC()
	formattedTime := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	var fm featureMetric
	fm.evaluationTime = formattedTime
	fm.count = 1
	meteringData := make(map[string]map[string]map[string]map[string]map[string]map[string]featureMetric)
	var modifyKey string
	if featureID != "" {
		meteringData = meteringInstance.meteringFeatureData
		modifyKey = featureID
	} else {
		meteringData = meteringInstance.meteringPropertyData
		modifyKey = propertyID
	}
	if _, ok := meteringData[guid]; ok {
		guidVal := meteringData[guid]
		if _, ok := guidVal[environmentID]; ok {
			envIDVal := guidVal[environmentID]
			if _, ok := envIDVal[collectionID]; ok {
				collectionIDVal := envIDVal[collectionID]
				if _, ok := collectionIDVal[modifyKey]; ok {
					modifyKeyVal := collectionIDVal[modifyKey]
					if _, ok := modifyKeyVal[entityID]; ok {
						entityIDVal := modifyKeyVal[entityID]
						if _, ok := entityIDVal[segmentID]; ok {
							segmentIDVal := entityIDVal[segmentID]
							segmentIDVal.evaluationTime = formattedTime
							segmentIDVal.count = segmentIDVal.count + 1
							entityIDVal[segmentID] = segmentIDVal
						} else {
							entityIDVal[segmentID] = fm
						}
					} else {
						segmentMap := make(map[string]featureMetric)
						segmentMap[segmentID] = fm
						modifyKeyVal[entityID] = segmentMap
					}
				} else {
					segmentMap := make(map[string]featureMetric)
					entityMap := make(map[string]map[string]featureMetric)
					segmentMap[segmentID] = fm
					entityMap[entityID] = segmentMap
					collectionIDVal[modifyKey] = entityMap
				}
			} else {
				segmentMap := make(map[string]featureMetric)
				entityMap := make(map[string]map[string]featureMetric)
				modifyKeyMap := make(map[string]map[string]map[string]featureMetric)
				segmentMap[segmentID] = fm
				entityMap[entityID] = segmentMap
				modifyKeyMap[modifyKey] = entityMap
				envIDVal[collectionID] = modifyKeyMap
			}
		} else {
			segmentMap := make(map[string]featureMetric)
			entityMap := make(map[string]map[string]featureMetric)
			modifyKeyMap := make(map[string]map[string]map[string]featureMetric)
			collectionMap := make(map[string]map[string]map[string]map[string]featureMetric)
			segmentMap[segmentID] = fm
			entityMap[entityID] = segmentMap
			modifyKeyMap[modifyKey] = entityMap
			collectionMap[collectionID] = modifyKeyMap
			guidVal[environmentID] = collectionMap
		}
	} else {
		segmentMap := make(map[string]featureMetric)
		entityMap := make(map[string]map[string]featureMetric)
		modifyKeyMap := make(map[string]map[string]map[string]featureMetric)
		collectionMap := make(map[string]map[string]map[string]map[string]featureMetric)
		environmentMap := make(map[string]map[string]map[string]map[string]map[string]featureMetric)
		segmentMap[segmentID] = fm
		entityMap[entityID] = segmentMap
		modifyKeyMap[modifyKey] = entityMap
		collectionMap[collectionID] = modifyKeyMap
		environmentMap[environmentID] = collectionMap
		meteringData[guid] = environmentMap
	}
	mt.mu.Unlock()
}

// RecordEvaluation : Record Evaluation
func (mt *Metering) RecordEvaluation(featureID string, propertyID string, entityID string, segmentID string) {
	log.Debug(messages.RecordEval)
	mt.addMetering(mt.guid, mt.EnvironmentID, mt.CollectionID, entityID, segmentID, featureID, propertyID)
}
func (mt *Metering) buildRequestBody(sendMeteringData map[string]map[string]map[string]map[string]map[string]map[string]featureMetric, guidMap map[string][]CollectionUsages, key string) {

	for guid, environmentMap := range sendMeteringData {
		var collectionUsageArray []CollectionUsages
		if _, ok := guidMap[guid]; !ok {
			guidMap[guid] = collectionUsageArray
		}
		for environmentID, collectionMap := range environmentMap {
			for collectionID, featureMap := range collectionMap {
				var usagesArray []Usages
				for featureID, entityMap := range featureMap {
					for entityID, segmentMap := range entityMap {
						for segmentID, val := range segmentMap {
							var usages Usages
							if key == "feature_id" {
								usages.FeatureID = featureID
							} else {
								usages.PropertyID = featureID
							}
							if entityID == constants.DefaultEntityID {
								usages.EntityID = nil
							} else {
								usages.EntityID = entityID
							}
							if segmentID == constants.DefaultSegmentID {
								usages.SegmentID = nil
							} else {
								usages.SegmentID = segmentID
							}
							usages.EvaluationTime = val.evaluationTime
							usages.Count = val.count
							usagesArray = append(usagesArray, usages)
						}
					}
				}
				var collectionUsageElem CollectionUsages
				collectionUsageElem.CollectionID = collectionID
				collectionUsageElem.EnvironmentID = environmentID
				collectionUsageElem.Usages = usagesArray
				collectionUsageArray = append(collectionUsageArray, collectionUsageElem)
			}
		}
		guidMap[guid] = append(guidMap[guid], collectionUsageArray...)
	}
}
func (mt *Metering) sendMetering() {
	log.Debug(messages.TenMinExpiry)
	defer GracefullyHandleError()
	log.Debug(mt.meteringFeatureData)
	log.Debug(mt.meteringPropertyData)
	mt.mu.Lock()
	if len(mt.meteringFeatureData) <= 0 && len(mt.meteringPropertyData) <= 0 {
		mt.mu.Unlock()
		return
	}
	sendFeatureData := make(map[string]map[string]map[string]map[string]map[string]map[string]featureMetric)
	sendFeatureData = mt.meteringFeatureData
	meteringFeatureDataMap := make(map[string]map[string]map[string]map[string]map[string]map[string]featureMetric)
	mt.meteringFeatureData = meteringFeatureDataMap

	sendPropertyData := make(map[string]map[string]map[string]map[string]map[string]map[string]featureMetric)
	sendPropertyData = mt.meteringPropertyData
	meteringPropertyDataMap := make(map[string]map[string]map[string]map[string]map[string]map[string]featureMetric)
	mt.meteringPropertyData = meteringPropertyDataMap

	mt.mu.Unlock()

	guidMap := make(map[string][]CollectionUsages)

	if len(sendFeatureData) > 0 {
		mt.buildRequestBody(sendFeatureData, guidMap, "feature_id")
	}

	if len(sendPropertyData) > 0 {
		mt.buildRequestBody(sendPropertyData, guidMap, "property_id")
	}

	for guid, val := range guidMap {
		for _, collectionUsage := range val {
			var count int = len(collectionUsage.Usages)
			if count > constants.DefaultUsageLimit {
				mt.sendSplitMetering(guid, collectionUsage, count)
			} else {
				mt.sendToServer(guid, collectionUsage)
			}
		}
	}

}
func (mt *Metering) sendSplitMetering(guid string, collectionUsages CollectionUsages, count int) {
	var lim int = 0
	subUsages := collectionUsages.Usages
	for lim <= count {
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
		mt.sendToServer(guid, collectionUsageElem)
		lim = lim + constants.DefaultUsageLimit
	}
}
func (mt *Metering) sendToServer(guid string, collectionUsages CollectionUsages) {
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
	response := GetAPIManagerInstance().Request(builder)
	if response != nil && response.StatusCode >= 200 && response.StatusCode <= 299 {
		log.Debug(messages.SendMeteringSuccess)
	} else {
		log.Error(messages.SendMeteringServerErr, err)
		return
	}
}
