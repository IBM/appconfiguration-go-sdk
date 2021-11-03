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

// Property : Property struct
type Property struct {
	Name         string        `json:"name"`
	PropertyID   string        `json:"property_id"`
	DataType     string        `json:"type"`
	Format       string        `json:"format"`
	Value        interface{}   `json:"value"`
	SegmentRules []SegmentRule `json:"segment_rules"`
}

// GetPropertyName : Get Property Name
func (p *Property) GetPropertyName() string {
	return p.Name
}

// GetPropertyID : Get Property Id
func (p *Property) GetPropertyID() string {
	return p.PropertyID
}

// GetPropertyDataType : Get Property Data Type
func (p *Property) GetPropertyDataType() string {
	return p.DataType
}

// GetPropertyDataFormat : Get Property Data Format
func (p *Property) GetPropertyDataFormat() string {
	// Format will be empty string ("") for Boolean & Numeric properties
	// If the Format is empty for a String type, we default it to TEXT
	if p.Format == "" && p.DataType == "STRING" {
		p.Format = "TEXT"
	}
	return p.Format
}

// GetValue : Get Value
func (p *Property) GetValue() interface{} {
	if p.Format == "YAML" {
		return getTypeCastedValue(p.Value, p.GetPropertyDataType(), p.GetPropertyDataFormat())
	}
	return p.Value
}

// GetSegmentRules : Get Segment Rules
func (p *Property) GetSegmentRules() []SegmentRule {
	return p.SegmentRules
}

// GetCurrentValue : Get Current Value
func (p *Property) GetCurrentValue(entityID string, entityAttributes map[string]interface{}) interface{} {
	log.Debug(messages.RetrievingProperty)
	if len(entityID) <= 0 {
		log.Error(messages.SetEntityObjectIDError)
		return nil
	}

	if p.isPropertyValid() {
		val := p.propertyEvaluation(entityID, entityAttributes)
		return getTypeCastedValue(val, p.GetPropertyDataType(), p.GetPropertyDataFormat())
	}
	return nil
}

func (p *Property) isPropertyValid() bool {
	return !(p.Name == "" || p.PropertyID == "" || p.DataType == "" || p.Value == nil)
}

func (p *Property) propertyEvaluation(entityID string, entityAttributes map[string]interface{}) interface{} {

	var evaluatedSegmentID string = constants.DefaultSegmentID
	defer func() {
		utils.GetMeteringInstance().RecordEvaluation("", p.GetPropertyID(), entityID, evaluatedSegmentID)
	}()

	log.Debug(messages.EvaluatingProperty)
	defer utils.GracefullyHandleError()

	if len(entityAttributes) < 0 {
		log.Debug(p.GetValue())
		return p.GetValue()
	}

	if len(p.GetSegmentRules()) > 0 {

		var rulesMap map[int]SegmentRule
		rulesMap = p.parseRules(p.GetSegmentRules())

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
					if p.evaluateSegment(string(segmentKey), entityAttributes) {
						evaluatedSegmentID = segmentKey
						if segmentRule.GetValue() == "$default" {
							log.Debug(messages.PropertyValue)
							log.Debug(p.GetValue())
							return p.GetValue()
						}
						log.Debug(messages.PropertyValue)
						log.Debug(segmentRule.GetValue())
						return segmentRule.GetValue()
					}
				}
			}
		}
	} else {
		return p.GetValue()
	}
	return p.GetValue()
}
func (p *Property) parseRules(segmentRules []SegmentRule) map[int]SegmentRule {
	log.Debug(messages.ParsingPropertyRules)
	defer utils.GracefullyHandleError()
	var rulesMap map[int]SegmentRule
	rulesMap = make(map[int]SegmentRule)
	for _, rule := range segmentRules {
		rulesMap[rule.GetOrder()] = rule
	}
	log.Debug(rulesMap)
	return rulesMap
}
func (p *Property) evaluateSegment(segmentKey string, entityAttributes map[string]interface{}) bool {
	log.Debug(messages.EvaluatingSegments)
	segment, ok := GetCacheInstance().SegmentMap[segmentKey]
	if ok {
		return segment.EvaluateRule(entityAttributes)
	}
	return false
}
