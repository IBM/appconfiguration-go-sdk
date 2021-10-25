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
	Name          string        `json:"name"`
	FeatureID     string        `json:"feature_id"`
	DataType      string        `json:"type"`
	Format        string        `json:"format"`
	EnabledValue  interface{}   `json:"enabled_value"`
	DisabledValue interface{}   `json:"disabled_value"`
	SegmentRules  []SegmentRule `json:"segment_rules"`
	Enabled       bool          `json:"enabled"`
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

// IsEnabled : Is Enabled
func (f *Feature) IsEnabled() bool {
	return f.Enabled
}

// GetSegmentRules : Get Segment Rules
func (f *Feature) GetSegmentRules() []SegmentRule {
	return f.SegmentRules
}

// GetCurrentValue : Get Current Value
func (f *Feature) GetCurrentValue(entityID string, entityAttributes map[string]interface{}) interface{} {
	log.Debug(messages.RetrievingFeature)
	if len(entityID) <= 0 {
		log.Error(messages.SetEntityObjectIDError)
		return nil
	}

	if f.isFeatureValid() {
		val := f.featureEvaluation(entityID, entityAttributes)
		return getTypeCastedValue(val, f.GetFeatureDataType(), f.GetFeatureDataFormat())
	}
	return nil
}

func (f *Feature) isFeatureValid() bool {
	return !(f.Name == "" || f.FeatureID == "" || f.DataType == "" || f.EnabledValue == nil || f.DisabledValue == nil)
}
func (f *Feature) featureEvaluation(entityID string, entityAttributes map[string]interface{}) interface{} {

	var evaluatedSegmentID string = constants.DefaultSegmentID
	defer func() {
		utils.GetMeteringInstance().RecordEvaluation(f.GetFeatureID(), "", entityID, evaluatedSegmentID)
	}()

	if f.IsEnabled() {
		log.Debug(messages.EvaluatingFeature)
		defer utils.GracefullyHandleError()

		if len(entityAttributes) < 0 {
			log.Debug(f.GetEnabledValue())
			return f.GetEnabledValue()
		}

		if len(f.GetSegmentRules()) > 0 {

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
							if segmentRule.GetValue() == "$default" {
								log.Debug(messages.FeatureValue)
								log.Debug(f.GetEnabledValue())
								return f.GetEnabledValue()
							}
							log.Debug(messages.FeatureValue)
							log.Debug(segmentRule.GetValue())
							return segmentRule.GetValue()
						}
					}
				}
			}
		} else {
			return f.GetEnabledValue()
		}
		return f.GetEnabledValue()
	}
	return f.GetDisabledValue()
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
