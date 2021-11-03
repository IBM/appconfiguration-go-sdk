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
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var rule = Rule{
	Operator:      "startsWith",
	AttributeName: "attribute_name",
	Values:        values,
}

var strs = []string{"first"}
var values = make([]interface{}, len(strs))

var segment = Segment{
	Name:      "segmentName",
	SegmentID: "segmentID",
	Rules:     []Rule{rule},
}

var ruleElem = RuleElem{
	Segments: []string{"segmentID"},
}

var segmentRule = SegmentRule{
	Order: 1,
	Value: true,
	Rules: []RuleElem{ruleElem},
}

var feature = Feature{
	Name:          "featureName",
	FeatureID:     "featureID",
	EnabledValue:  true,
	DisabledValue: false,
	Enabled:       true,
	DataType:      "BOOLEAN",
	Format:        "",
	SegmentRules:  []SegmentRule{segmentRule},
}

var property = Property{
	DataType:     "BOOLEAN",
	Format:       "",
	Name:         "propertyName",
	PropertyID:   "propertyID",
	Value:        true,
	SegmentRules: []SegmentRule{segmentRule},
}

func TestCacheWithDebugMode(t *testing.T) {
	os.Setenv("ENABLE_DEBUG", "true")
	featureMap := make(map[string]Feature)
	featureMap["featureID"] = feature
	segmentMap := make(map[string]Segment)
	segmentMap["segmentID"] = segment
	propertyMap := make(map[string]Property)
	propertyMap["propertyID"] = property
	SetCache(featureMap, propertyMap, segmentMap)
	cacheInstance := GetCacheInstance()
	if !reflect.DeepEqual(cacheInstance.FeatureMap, featureMap) {
		t.Error("Expected TestCacheFeatureMap test case to pass")
	}
	if !reflect.DeepEqual(cacheInstance.SegmentMap, segmentMap) {
		t.Error("Expected TestCacheSegmentMap test case to pass")
	}
	if !reflect.DeepEqual(cacheInstance.PropertyMap, propertyMap) {
		t.Error("Expected TestCachePropertyMap test case to pass")
	}
}

func TestFeature(t *testing.T) {
	if feature.GetFeatureID() != "featureID" {
		t.Error("Expected TestFeatureGetFeatureID test case to pass")
	}
	if feature.GetFeatureName() != "featureName" {
		t.Error("Expected TestFeatureGetFeatureName test case to pass")
	}
	if feature.GetFeatureDataType() != "BOOLEAN" {
		t.Error("Expected TestFeatureGetFeatureDataType test case to pass")
	}
	if feature.GetFeatureDataFormat() != "" {
		t.Error("Expected TestFeatureGetFeatureDataFormat test case to pass")
	}
	if feature.GetEnabledValue() != true {
		t.Error("Expected TestFeatureGetEnabledValue test case to pass")
	}
	if feature.GetDisabledValue() != false {
		t.Error("Expected TestFeatureGetDisabledValue test case to pass")
	}
	if feature.IsEnabled() != true {
		t.Error("Expected TestFeatureIsEnabled test case to pass")
	}
	if !reflect.DeepEqual(feature.GetSegmentRules()[0], segmentRule) {
		t.Error("Expected TestFeatureGetSegmentRules test case to pass")
	}
	entityMap := make(map[string]interface{})
	entityMap["attribute_name"] = "first"
	if feature.GetCurrentValue("entityID123", entityMap) != true {
		t.Error("Expected TestFeatureGetCurrentValueBoolean test case to pass")
	}
	feature.DataType = "STRING"
	feature.Format = "TEXT"
	feature.EnabledValue = "EnabledValue"
	feature.DisabledValue = "DisabledValue"
	if feature.GetCurrentValue("entityID123", entityMap) != "EnabledValue" {
		t.Error("Expected TestFeatureGetCurrentValueStringText test case to pass")
	}
	feature.DataType = "STRING"
	feature.Format = "JSON"
	enabledJSON := make(map[string]interface{})
	enabledJSON["key"] = "enabled value"
	feature.EnabledValue = enabledJSON
	disabledJSON := make(map[string]interface{})
	disabledJSON["key"] = "disabled value"
	feature.DisabledValue = disabledJSON
	if !reflect.DeepEqual(feature.GetCurrentValue("entityId123", entityMap), enabledJSON) {
		t.Error("Expected TestFeatureGetCurrentValueStringJSON test case to pass")
	}
	feature.DataType = "STRING"
	feature.Format = "YAML"
	feature.EnabledValue = "men:\n  - John Smith\n  - Bill Jones\nwomen:\n  - Mary Smith\n  - Susan Williams"
	feature.DisabledValue = "key:value"
	if !reflect.DeepEqual(feature.GetCurrentValue("entityId123", entityMap), feature.GetEnabledValue()) {
		t.Error("Expected TestFeatureGetCurrentValueStringYAML test case to pass")
	}
	feature.DataType = "NUMERIC"
	feature.Format = ""
	feature.EnabledValue = float64(1)
	feature.DisabledValue = float64(0)
	if feature.GetCurrentValue("entityID123", entityMap) != float64(1) {
		t.Error("Expected TestFeatureGetCurrentValueNumeric test case to pass")
	}
	feature.DataType = "BOOLEAN"
	feature.EnabledValue = true
	feature.DisabledValue = false
	feature.Enabled = false
	if feature.GetCurrentValue("entityID123", entityMap) != false {
		t.Error("Expected TestFeatureGetCurrentValueDisabledFeature test case to pass")
	}
	feature.Enabled = true

	if feature.GetCurrentValue("", entityMap) != nil {
		t.Error("Expected TestFeatureGetCurrentValueWithEmptyEntityID test case to pass")
	}
	feature.FeatureID = ""
	if feature.GetCurrentValue("entityID123", entityMap) != nil {
		t.Error("Expected TestFeatureGetCurrentValueWithEmptyFeatureID test case to pass")
	}
	feature.FeatureID = "featureID"

	feature.SegmentRules = []SegmentRule{}
	if feature.GetCurrentValue("entityID123", entityMap) != true {
		t.Error("Expected TestFeatureGetCurrentValueWithEmptySegmentRules test case to pass")
	}
	feature.SegmentRules = []SegmentRule{segmentRule}

	entityMap = make(map[string]interface{})
	entityMap["attributeName"] = "FirstLast"
	if feature.GetCurrentValue("entityID123", entityMap) != true {
		t.Error("Expected TestFeatureGetCurrentValueWrongAttribute test case to pass")
	}
}

func TestProperty(t *testing.T) {
	if property.GetPropertyID() != "propertyID" {
		t.Error("Expected TestPropertyGetPropertyID test case to pass")
	}
	if property.GetPropertyName() != "propertyName" {
		t.Error("Expected TestPropertyGetPropertyName test case to pass")
	}
	if property.GetPropertyDataType() != "BOOLEAN" {
		t.Error("Expected TestPropertyGetPropertyDataType test case to pass")
	}
	if property.GetPropertyDataFormat() != "" {
		t.Error("Expected TestPropertyGetPropertyDataFormat test case to pass")
	}
	if property.GetValue() != true {
		t.Error("Expected TestPropertyGetValue test case to pass")
	}
	if !reflect.DeepEqual(property.GetSegmentRules()[0], segmentRule) {
		t.Error("Expected TestPropertyGetSegmentRules test case to pass")
	}
	entityMap := make(map[string]interface{})
	entityMap["attribute_name"] = "first"
	if property.GetCurrentValue("entityID123", entityMap) != true {
		t.Error("Expected TestPropertyGetCurrentValueBoolean test case to pass")
	}
	property.DataType = "STRING"
	property.Format = "TEXT"
	property.Value = "Value"
	if property.GetCurrentValue("entityID123", entityMap) != "Value" {
		t.Error("Expected TestPropertyGetCurrentValueStringText test case to pass")
	}
	property.DataType = "STRING"
	property.Format = "JSON"
	propertyValueJSON := make(map[string]interface{})
	propertyValueJSON["key"] = "property value"
	property.Value = propertyValueJSON
	if !reflect.DeepEqual(property.GetCurrentValue("entityId123", entityMap), propertyValueJSON) {
		t.Error("Expected TestPropertyGetCurrentValueStringJson test case to pass")
	}
	property.DataType = "STRING"
	property.Format = "YAML"
	property.Value = "men:\n  - John Smith\n  - Bill Jones\nwomen:\n  - Mary Smith\n  - Susan Williams"
	if !reflect.DeepEqual(property.GetCurrentValue("entityId123", entityMap), property.GetValue()) {
		t.Error("Expected TestPropertyGetCurrentValueStringYaml test case to pass")
	}
	property.DataType = "NUMERIC"
	property.Format = ""
	property.Value = float64(1)
	if property.GetCurrentValue("entityID123", entityMap) != float64(1) {
		t.Error("Expected TestPropertyGetCurrentValueNumeric test case to pass")
	}

	property.DataType = "BOOLEAN"
	property.Value = true

	if property.GetCurrentValue("", entityMap) != nil {
		t.Error("Expected TestPropertyGetCurrentValueWithEmptyEntityID test case to pass")
	}
	property.PropertyID = ""
	if property.GetCurrentValue("entityID123", entityMap) != nil {
		t.Error("Expected TestPropertyGetCurrentValueWithEmptyPropertyID test case to pass")
	}
	property.PropertyID = "propertyID"

	property.SegmentRules = []SegmentRule{}
	if property.GetCurrentValue("entityID123", entityMap) != true {
		t.Error("Expected TestPropertyGetCurrentValueWithEmptySegmentRules test case to pass")
	}
	property.SegmentRules = []SegmentRule{segmentRule}

	entityMap = make(map[string]interface{})
	entityMap["attributeName"] = "FirstLast"
	if property.GetCurrentValue("entityID123", entityMap) != true {
		t.Error("Expected TestPropertyGetCurrentValueWrongAttribute test case to pass")
	}
}

func TestSegment(t *testing.T) {
	if segment.GetName() != "segmentName" {
		t.Error("Expected TestSegmentGetName test case to pass")
	}
	if segment.GetSegmentID() != "segmentID" {
		t.Error("Expected TestSegmentGetSegmentID test case to pass")
	}
	if !reflect.DeepEqual(segment.GetRules(), []Rule{rule}) {
		t.Error("Expected TestSegmentGetRules test case to pass")
	}
	entityMap := make(map[string]interface{})
	entityMap["k1"] = 7
	if segment.EvaluateRule(entityMap) != false {
		t.Error("Expected TestSegmentEvaluateRule test case to pass")
	}
}

func TestSegmentRule(t *testing.T) {
	if segmentRule.GetValue() != true {
		t.Error("Expected TestSegmentRuleGetValue test case to pass")
	}
	if segmentRule.GetOrder() != 1 {
		t.Error("Expected TestSegmentRuleGetOrder test case to pass")
	}
	if !reflect.DeepEqual(segmentRule.GetRules()[0].Segments, ruleElem.Segments) {
		t.Error("Expected TestSegmentRuleGetRules test case to pass")
	}
	segmentRule.GetRules()
}

func TestRule(t *testing.T) {
	if rule.GetOperator() != "startsWith" {
		t.Error("Expected TestRuleGetOperator test case to pass")
	}
	if rule.GetAttributeName() != "attribute_name" {
		t.Error("Expected TestRuleGetAttributeName test case to pass")
	}
	if !reflect.DeepEqual(rule.GetValues(), values) {
		t.Error("Expected TestRuleGetValues test case to pass")
	}
	entityMap := make(map[string]interface{})
	entityMap["attribute_name"] = "first"
	if rule.EvaluateRule(entityMap) != false {
		t.Error("Expected TestRuleEvaluateRule test case to pass")
	}
	entityMap["attribute_name"] = "last"
	if rule.EvaluateRule(entityMap) != false {
		t.Error("Expected TestRuleEvaluateRule test case to pass")
	}

	//
	if isNumber(1) != true {
		t.Error("Expected TestIsNumber test case to pass when input provided is a number.")
	}
	if isNumber("a") != false {
		t.Error("Expected TestIsNumber test case to pass when input provided is a string.")
	}
	//
	if isBool("a") != false {
		t.Error("Expected TestIsBool test case to pass when input provided is a string.")
	}
	if isBool(true) != true {
		t.Error("Expected TestIsBool test case to passwhen input provided is a string.")
	}
	//
	if isString("a") != true {
		t.Error("Expected TestIsString test case to pass when input provided is a string.")
	}
	if isString(1) != false {
		t.Error("Expected TestIsString test case to pass when input provided is a number.")
	}

	//
	if val, _ := formatBool(true); val != "true" {
		t.Error("Expected TestFormatBool test case to pass when input provided is boolean true.")
	}
	if val, _ := formatBool(false); val != "false" {
		t.Error("Expected TestFormatBool test case to pass when input provided is boolean false.")
	}

	val := rule.operatorCheck("ibm.com", "ibm")
	assert.Equal(t, true, val)

	//

	rule = Rule{
		Operator: "endsWith",
	}
	val = rule.operatorCheck("ibm.com", "com")
	assert.Equal(t, true, val)

	//

	rule = Rule{
		Operator: "contains",
	}
	val = rule.operatorCheck("ibm.com", "ibm")
	assert.Equal(t, true, val)

	rule = Rule{
		Operator: "is",
	}
	val = rule.operatorCheck("ibm.com", "ibm.com")
	assert.Equal(t, true, val)

	val = rule.operatorCheck(1.5, "1.5")
	assert.Equal(t, true, val)

	val = rule.operatorCheck(true, "true")
	assert.Equal(t, true, val)

	rule = Rule{
		Operator: "greaterThan",
	}
	val = rule.operatorCheck(1.5, "1")
	assert.Equal(t, true, val)

	val = rule.operatorCheck("1.5", "1")
	assert.Equal(t, true, val)

	rule = Rule{
		Operator: "greaterThanEquals",
	}
	val = rule.operatorCheck(1.5, "1.5")
	assert.Equal(t, true, val)

	val = rule.operatorCheck("1.5", "1.5")
	assert.Equal(t, true, val)

	rule = Rule{
		Operator: "lesserThan",
	}
	val = rule.operatorCheck(0.5, "1")
	assert.Equal(t, true, val)

	val = rule.operatorCheck("0.5", "1")
	assert.Equal(t, true, val)

	rule = Rule{
		Operator: "lesserThanEquals",
	}
	val = rule.operatorCheck(0.5, "0.5")
	assert.Equal(t, true, val)

	val = rule.operatorCheck("0.5", "0.5")
	assert.Equal(t, true, val)

}
