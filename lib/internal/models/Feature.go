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

package models

import (
	constants "github.com/IBM/appconfiguration-go-sdk/lib/internal/constants"
	messages "github.com/IBM/appconfiguration-go-sdk/lib/internal/messages"
	utils "github.com/IBM/appconfiguration-go-sdk/lib/internal/utils"

	"sort"

	"github.com/IBM/appconfiguration-go-sdk/lib/internal/utils/log"
)

// Feature : Feature struct
type Feature struct {
	Name              string        `json:"name"`
	FeatureID         string        `json:"feature_id"`
	DataType          string        `json:"type"`
	Format            string        `json:"format"`
	EnabledValue      interface{}   `json:"enabled_value"`
	DisabledValue     interface{}   `json:"disabled_value"`
	SegmentRules      []SegmentRule `json:"segment_rules"`
	Enabled           bool          `json:"enabled"`
	RolloutPercentage *int          `json:"rollout_percentage"`
}

// GetFeatureName : Get Feature Name
func (f *Feature) GetFeatureName() string {
	return f.Name
}

// GetDisabledValue : Get Disabled Value
func (f *Feature) GetDisabledValue() interface{} {
	if f.Format == "YAML" {
		return getTypeCastedValue(f.DisabledValue, f.GetFeatureDataType(), f.GetFeatureDataFormat())
	}
	return f.DisabledValue
}

// GetEnabledValue : Get Enabled Value
func (f *Feature) GetEnabledValue() interface{} {
	if f.Format == "YAML" {
		return getTypeCastedValue(f.EnabledValue, f.GetFeatureDataType(), f.GetFeatureDataFormat())
	}
	return f.EnabledValue
}

// GetFeatureID : Get Feature ID
func (f *Feature) GetFeatureID() string {
	return f.FeatureID
}

// GetFeatureDataType : Get Feature Data Type
func (f *Feature) GetFeatureDataType() string {
	return f.DataType
}

// GetFeatureDataFormat : Get Feature Data Format
func (f *Feature) GetFeatureDataFormat() string {
	// Format will be empty string ("") for Boolean & Numeric feature flags
	// If the Format is empty for a String type, we default it to TEXT
	if f.Format == "" && f.DataType == "STRING" {
		f.Format = "TEXT"
	}
	return f.Format
}

// GetRolloutPercentage : Get the Feature flag rollout percentage
func (f *Feature) GetRolloutPercentage() int {
	if f.RolloutPercentage == nil {
		var v int = 100
		f.RolloutPercentage = &v
	}
	return *f.RolloutPercentage
}

// GetSegmentRules : Get Segment Rules
func (f *Feature) GetSegmentRules() []SegmentRule {
	return f.SegmentRules
}

// IsEnabled returns the state of the feature flag.
// Returns true, if the feature flag is enabled, otherwise returns false.
func (f *Feature) IsEnabled() bool {
	return f.Enabled
}

// GetCurrentValue returns one of the Enabled/Disabled/Overridden value based on the evaluation.
//
// The function takes in entityId & entityAttributes parameters.
//
// entityId is a string identifier related to the Entity against which the feature will be evaluated.
// For example, an entity might be an instance of an app that runs on a mobile device, a microservice that runs on the cloud, or a component of infrastructure that runs that microservice.
// For any entity to interact with App Configuration, it must provide a unique entity ID.
//
// entityAttributes is a map of type `map[string]interface{}` consisting of the attribute name and their values that defines the specified entity.
// This is an optional parameter if the feature flag is not configured with any targeting definition.
// If the targeting is configured, then entityAttributes should be provided for the rule evaluation.
// An attribute is a parameter that is used to define a segment. The SDK uses the attribute values to determine if the
// specified entity satisfies the targeting rules, and returns the appropriate feature flag value.
func (f *Feature) GetCurrentValue(entityID string, entityAttributes ...map[string]interface{}) interface{} {
	log.Debug(messages.RetrievingFeature)
	if len(entityID) <= 0 {
		log.Error("Feature flag evaluation: ", messages.InvalidEntityId, "GetCurrentValue")
		return nil
	}
	var temp map[string]interface{}
	switch len(entityAttributes) {
	case 0: // Do Nothing
	case 1:
		temp = entityAttributes[0]
	default:
		log.Error("Feature flag evaluation: ", messages.IncorrectUsageOfEntityAttributes, "GetCurrentValue")
		return nil
	}
	if f.isFeatureValid() {
		val, _ := f.featureEvaluation(entityID, temp)
		return getTypeCastedValue(val, f.GetFeatureDataType(), f.GetFeatureDataFormat())
	}
	return nil
}

func (f *Feature) isFeatureValid() bool {
	return !(f.Name == "" || f.FeatureID == "" || f.DataType == "" || f.EnabledValue == nil || f.DisabledValue == nil)
}
func (f *Feature) featureEvaluation(entityID string, entityAttributes map[string]interface{}) (interface{}, bool) {

	var evaluatedSegmentID string = constants.DefaultSegmentID
	defer func() {
		utils.GetMeteringInstance().RecordEvaluation(f.GetFeatureID(), "", entityID, evaluatedSegmentID)
	}()

	if f.Enabled {
		log.Debug(messages.EvaluatingFeature)
		defer utils.GracefullyHandleError()

		if len(f.GetSegmentRules()) > 0 && len(entityAttributes) > 0 {
			var rulesMap map[int]SegmentRule
			rulesMap = f.parseRules(f.GetSegmentRules())

			// sort the map elements as per ascending order of keys
			var keys []int
			for k := range rulesMap {
				keys = append(keys, k)
			}
			sort.Ints(keys)

			// after sorting , pick up each map element as per keys order
			for _, k := range keys {
				segmentRule := rulesMap[k]
				for _, rule := range segmentRule.GetRules() {
					for _, segmentKey := range rule.Segments {
						if f.evaluateSegment(string(segmentKey), entityAttributes) {
							evaluatedSegmentID = segmentKey
							var segmentLevelRolloutPercentage int
							if segmentRule.GetRolloutPercentage() == "$default" {
								segmentLevelRolloutPercentage = f.GetRolloutPercentage()
							} else {
								segmentLevelRolloutPercentage = int(segmentRule.GetRolloutPercentage().(float64))
							}
							if segmentLevelRolloutPercentage == 100 || utils.GetNormalizedValue(entityID+":"+f.GetFeatureID()) < segmentLevelRolloutPercentage {
								if segmentRule.GetValue() == "$default" {
									return f.GetEnabledValue(), true
								} else {
									return segmentRule.GetValue(), true
								}
							} else {
								return f.GetDisabledValue(), false
							}
						}
					}
				}
			}
		}
		if f.GetRolloutPercentage() == 100 || utils.GetNormalizedValue(entityID+":"+f.GetFeatureID()) < f.GetRolloutPercentage() {
			return f.GetEnabledValue(), true
		}
		return f.GetDisabledValue(), false
	}
	return f.GetDisabledValue(), false
}
func (f *Feature) parseRules(segmentRules []SegmentRule) map[int]SegmentRule {
	log.Debug(messages.ParsingFeatureRules)
	defer utils.GracefullyHandleError()
	var rulesMap map[int]SegmentRule
	rulesMap = make(map[int]SegmentRule)
	for _, rule := range segmentRules {
		rulesMap[rule.GetOrder()] = rule
	}
	log.Debug(rulesMap)
	return rulesMap
}
func (f *Feature) evaluateSegment(segmentKey string, entityAttributes map[string]interface{}) bool {
	log.Debug(messages.EvaluatingSegments)
	segment, ok := GetCacheInstance().SegmentMap[segmentKey]
	if ok {
		return segment.EvaluateRule(entityAttributes)
	}
	return false
}
