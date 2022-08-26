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

// GetCurrentValue returns the default property value or its overridden value based on the evaluation.
//
// The function takes in entityId & entityAttributes parameters.
//
// entityId is a string identifier related to the Entity against which the property will be evaluated.
// For example, an entity might be an instance of an app that runs on a mobile device, a microservice that runs on the cloud, or a component of infrastructure that runs that microservice.
// For any entity to interact with App Configuration, it must provide a unique entity ID.
//
// entityAttributes is a map of type `map[string]interface{}` consisting of the attribute name and their values that defines the specified entity.
// This is an optional parameter if the property is not configured with any targeting definition.
// If the targeting is configured, then entityAttributes should be provided for the rule evaluation.
// An attribute is a parameter that is used to define a segment. The SDK uses the attribute values to determine if the
// specified entity satisfies the targeting rules, and returns the appropriate property value.
func (p *Property) GetCurrentValue(entityID string, entityAttributes ...map[string]interface{}) interface{} {
	log.Debug(messages.RetrievingProperty)
	if len(entityID) <= 0 {
		log.Error("Property evaluation: ", messages.InvalidEntityId, "GetCurrentValue")
		return nil
	}
	var temp map[string]interface{}
	switch len(entityAttributes) {
	case 0: // Do Nothing
	case 1:
		temp = entityAttributes[0]
	default:
		log.Error("Property evaluation: ", messages.IncorrectUsageOfEntityAttributes, "GetCurrentValue")
		return nil
	}

	if p.isPropertyValid() {
		val := p.propertyEvaluation(entityID, temp)
		return getTypeCastedValue(val, p.GetPropertyDataType(), p.GetPropertyDataFormat())
	}
	log.Error("Invalid property")
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

	if len(p.GetSegmentRules()) > 0 && len(entityAttributes) > 0 {
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
