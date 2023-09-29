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
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/utils/log"
	"github.com/sirupsen/logrus/hooks/test"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var value []interface{}

func Int(v int) *int { return &v }

func Interface(v interface{}) *interface{} { return &v }

var rule = Rule{
	Operator:      "startsWith",
	AttributeName: "attribute_name",
	Values:        append(value, "first"),
}

var segment = Segment{
	Name:      "segmentName",
	SegmentID: "segmentID",
	Rules:     []Rule{rule},
}

var ruleElem = RuleElem{
	Segments: []string{"segmentID"},
}

var segmentRule = SegmentRule{
	RolloutPercentage: Interface(50.0),
	Order:             1,
	Value:             true,
	Rules:             []RuleElem{ruleElem},
}

var feature = Feature{
	Name:              "featureName",
	FeatureID:         "featureID",
	EnabledValue:      true,
	DisabledValue:     false,
	Enabled:           true,
	DataType:          "BOOLEAN",
	Format:            "",
	SegmentRules:      []SegmentRule{segmentRule},
	RolloutPercentage: Int(60),
}

var property = Property{
	DataType:     "BOOLEAN",
	Format:       "",
	Name:         "propertyName",
	PropertyID:   "propertyID",
	Value:        true,
	SegmentRules: []SegmentRule{segmentRule},
}

var secretPropertyValueJSON = map[string]interface{}{"secret_type": "username_password", "sm_instance_crn": "crn:v1:staging:public:secrets-manager:eu-gb:a/3268cfe9e25d41q1232132f9a731a:d614a8ba-a13a-41cc-9e18-82132133ad9845::"}
var secretPropertySegmentValueJSON = map[string]interface{}{"id": "1212433434353535"}

var secretRefSegmentRule = SegmentRule{
	Order: 1,
	Value: secretPropertySegmentValueJSON,
	Rules: []RuleElem{ruleElem},
}

var propertySecretRefData = Property{
	DataType:     "SECRETREF",
	Format:       "",
	Name:         "propertyName",
	PropertyID:   "propertySecretDataId",
	Value:        secretPropertyValueJSON,
	SegmentRules: []SegmentRule{secretRefSegmentRule},
}

var secretproperty = SecretProperty{
	PropertyID: "propertySecretDataId",
}

func TestCacheWithDebugMode(t *testing.T) {
	featureMap := make(map[string]Feature)
	featureMap["featureID"] = feature
	segmentMap := make(map[string]Segment)
	segmentMap["segmentID"] = segment
	propertyMap := make(map[string]Property)
	propertyMap["propertyID"] = property
	propertyMap["propertySecretDataId"] = propertySecretRefData
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
	if feature.GetRolloutPercentage() != 60 {
		t.Error("Expected TestFeatureGetRolloutPercentage test case to pass")
	}
	if feature.IsEnabled() != true {
		t.Error("Expected TestFeatureIsEnabled test case to pass")
	}
	if !reflect.DeepEqual(feature.GetSegmentRules()[0], segmentRule) {
		t.Error("Expected TestFeatureGetSegmentRules test case to pass")
	}
	entityMap := make(map[string]interface{})
	entityMap["attribute_name"] = "first"
	if feature.GetCurrentValue("entityID123") != true {
		t.Error("Expected TestFeatureGetCurrentValueBoolean test case to pass")
	}
	if feature.GetCurrentValue("entityID123", entityMap) != true {
		t.Error("Expected TestFeatureGetCurrentValueBoolean test case to pass")
	}
	if feature.GetCurrentValue("entityID456", entityMap) != false {
		t.Error("Expected TestFeatureGetCurrentValueBoolean test case to pass")
	}
	feature.DataType = "STRING"
	feature.Format = "TEXT"
	feature.EnabledValue = "EnabledValue"
	feature.DisabledValue = "DisabledValue"
	segmentRule.Value = "OverriddenValue"
	feature.SegmentRules = []SegmentRule{segmentRule}
	if feature.GetCurrentValue("entityID123") != "EnabledValue" {
		t.Error("Expected TestFeatureGetCurrentValueStringText test case to pass")
	}
	if feature.GetCurrentValue("entityID123", entityMap) != "OverriddenValue" {
		t.Error("Expected TestFeatureGetCurrentValueStringText test case to pass")
	}
	if feature.GetCurrentValue("entityID456", entityMap) != "DisabledValue" {
		t.Error("Expected TestFeatureGetCurrentValueStringText test case to pass")
	}
	entityMap["attribute_name"] = "second"
	if feature.GetCurrentValue("entityID123", entityMap) != "EnabledValue" {
		t.Error("Expected TestFeatureGetCurrentValueStringText test case to pass")
	}
	if feature.GetCurrentValue("entityID456", entityMap) != "DisabledValue" {
		t.Error("Expected TestFeatureGetCurrentValueStringText test case to pass")
	}
	entityMap["attribute_name"] = "first"
	feature.DataType = "STRING"
	feature.Format = "JSON"
	enabledJSON := make(map[string]interface{})
	enabledJSON["key"] = "enabled value"
	feature.EnabledValue = enabledJSON
	disabledJSON := make(map[string]interface{})
	disabledJSON["key"] = "disabled value"
	feature.DisabledValue = disabledJSON
	overriddenJSON := make(map[string]interface{})
	overriddenJSON["key"] = "overridden value"
	segmentRule.Value = overriddenJSON
	feature.SegmentRules = []SegmentRule{segmentRule}
	if !reflect.DeepEqual(feature.GetCurrentValue("entityID123", entityMap), overriddenJSON) {
		t.Error("Expected TestFeatureGetCurrentValueStringJSON test case to pass")
	}
	feature.DataType = "STRING"
	feature.Format = "YAML"
	feature.EnabledValue = "men:\n  - John Smith\n  - Bill Jones\nwomen:\n  - Mary Smith\n  - Susan Williams"
	feature.DisabledValue = "key:value"
	segmentRule.Value = "key1:value1"
	feature.SegmentRules = []SegmentRule{segmentRule}
	if !reflect.DeepEqual(feature.GetCurrentValue("entityID123", entityMap), "key1:value1") {
		t.Error("Expected TestFeatureGetCurrentValueStringYAML test case to pass")
	}
	feature.DataType = "NUMERIC"
	feature.Format = ""
	feature.EnabledValue = float64(1)
	feature.DisabledValue = float64(0)
	segmentRule.Value = float64(5)
	feature.SegmentRules = []SegmentRule{segmentRule}
	if feature.GetCurrentValue("entityID123") != float64(1) {
		t.Error("Expected TestFeatureGetCurrentValueNumeric test case to pass")
	}
	if feature.GetCurrentValue("entityID456") != float64(0) {
		t.Error("Expected TestFeatureGetCurrentValueNumeric test case to pass")
	}
	if feature.GetCurrentValue("entityID123", entityMap) != float64(5) {
		t.Error("Expected TestFeatureGetCurrentValueNumeric test case to pass")
	}
	if feature.GetCurrentValue("entityID456", entityMap) != float64(0) {
		t.Error("Expected TestFeatureGetCurrentValueNumeric test case to pass")
	}
	feature.DataType = "INVALID_DATATYPE"
	feature.Format = ""
	if feature.GetCurrentValue("entityID123", entityMap) != nil {
		t.Error("Expected TestFeatureGetCurrentValueWithInvalidFeatureDatatype test case to pass")
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
	segmentRule.Value = true
	property.SegmentRules = []SegmentRule{segmentRule}
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
	segmentRule.Value = "OverriddenValue"
	property.SegmentRules = []SegmentRule{segmentRule}
	if property.GetCurrentValue("entityID123") != "Value" {
		t.Error("Expected TestPropertyGetCurrentValueStringText test case to pass")
	}
	if property.GetCurrentValue("entityID123", entityMap) != "OverriddenValue" {
		t.Error("Expected TestPropertyGetCurrentValueStringText test case to pass")
	}
	property.DataType = "STRING"
	property.Format = "JSON"
	propertyValueJSON := make(map[string]interface{})
	propertyValueJSON["key"] = "property value"
	property.Value = propertyValueJSON
	overriddenJSON := make(map[string]interface{})
	overriddenJSON["key"] = "overridden value"
	segmentRule.Value = overriddenJSON
	property.SegmentRules = []SegmentRule{segmentRule}
	if !reflect.DeepEqual(property.GetCurrentValue("entityId123", entityMap), overriddenJSON) {
		t.Error("Expected TestPropertyGetCurrentValueStringJson test case to pass")
	}
	property.DataType = "STRING"
	property.Format = "YAML"
	property.Value = "men:\n  - John Smith\n  - Bill Jones\nwomen:\n  - Mary Smith\n  - Susan Williams"
	segmentRule.Value = "key1:value1"
	property.SegmentRules = []SegmentRule{segmentRule}
	if !reflect.DeepEqual(property.GetCurrentValue("entityId123", entityMap), "key1:value1") {
		t.Error("Expected TestPropertyGetCurrentValueStringYaml test case to pass")
	}
	property.DataType = "NUMERIC"
	property.Format = ""
	property.Value = float64(1)
	segmentRule.Value = float64(5)
	property.SegmentRules = []SegmentRule{segmentRule}
	if property.GetCurrentValue("entityID123", entityMap) != float64(5) {
		t.Error("Expected TestPropertyGetCurrentValueNumeric test case to pass")
	}
	if property.GetCurrentValue("entityID123") != float64(1) {
		t.Error("Expected TestPropertyGetCurrentValueNumeric test case to pass")
	}
	property.DataType = "INVALID_DATATYPE"
	property.Format = ""
	property.Value = float64(1)
	if property.GetCurrentValue("entityID123", entityMap) != nil {
		t.Error("Expected TestPropertyGetCurrentValueWithInvalidPropertyDatatype test case to pass")
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

func TestSecretProperty(t *testing.T) {

	entityMap := make(map[string]interface{})
	entityMap["email"] = "user@ibm.com"
	if secretproperty.PropertyID != "propertySecretDataId" {
		t.Error("Expected test case to pass but failed")
	}
	_, _, entityIDError := secretproperty.GetCurrentValue("", entityMap)
	if entityIDError == nil {
		t.Error("Expected TestSecretPropertyGetCurrentValueWithEmptyEntityID test case to pass")
	}
	_, _, multipleEntityMapError := secretproperty.GetCurrentValue("entityID123", entityMap, entityMap)
	if multipleEntityMapError == nil {
		t.Error("Expected TestSecretPropertyGetCurrentValue test case to pass")
	}
	_, _, secretIDError := secretproperty.GetCurrentValue("entityID123", entityMap)
	if secretIDError == nil {
		t.Error("Expected TestSecretPropertyGetCurrentValueWithNoSecretId test case to pass")
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
	segmentRule.Value = true
	_ = []SegmentRule{segmentRule}

	if segmentRule.GetValue() != true {
		t.Error("Expected TestSegmentRuleGetValue test case to pass")
	}
	if segmentRule.GetOrder() != 1 {
		t.Error("Expected TestSegmentRuleGetOrder test case to pass")
	}
	if segmentRule.GetRolloutPercentage() != 50.0 {
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
	if !reflect.DeepEqual(rule.GetValues(), rule.Values) {
		t.Error("Expected TestRuleGetValues test case to pass")
	}
	entityMap := make(map[string]interface{})
	entityMap["attribute_name"] = "first"
	if rule.EvaluateRule(entityMap) != true {
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

func TestFormatConfig(t *testing.T) {
	data := `{"features":[{"name":"Cycle Rentals8","feature_id":"cycle-rentals8","type":"BOOLEAN","enabled_value":true,"disabled_value":false,"segment_rules":[],"enabled":true,"rollout_percentage":75}],"properties":[{"name":"ShowAd","property_id":"show-ad","tags":"","type":"BOOLEAN","value":false,"segment_rules":[]}],"segments":[{"name":"beta-users","segment_id":"knliu818","rules":[{"values":["ibm.com"],"operator":"contains","attribute_name":"email"}]},{"name":"ibm employees","segment_id":"ka761hap","rules":[{"values":["ibm.com","in.ibm.com"],"operator":"endsWith","attribute_name":"email"}]}]}`
	expectedReformat := `{"environments":[{"name":"myEnvironment","environment_id":"myEnvironment","features":[{"name":"Cycle Rentals8","feature_id":"cycle-rentals8","type":"BOOLEAN","format":"","enabled_value":true,"disabled_value":false,"segment_rules":[],"enabled":true,"rollout_percentage":75}],"properties":[{"name":"ShowAd","property_id":"show-ad","type":"BOOLEAN","format":"","value":false,"segment_rules":[]}]}],"segments":[{"name":"beta-users","segment_id":"knliu818","rules":[{"values":["ibm.com"],"operator":"contains","attribute_name":"email"}]},{"name":"ibm employees","segment_id":"ka761hap","rules":[{"values":["ibm.com","in.ibm.com"],"operator":"endsWith","attribute_name":"email"}]}]}`
	reformattedData := FormatConfig([]byte(data), "myEnvironment")
	if !reflect.DeepEqual([]byte(expectedReformat), reformattedData) {
		t.Error("Expected TestFormatConfig test case to pass")
	}
}

var testLogger, hook = test.NewNullLogger()

func mockLogger() {
	log.SetLogger(testLogger)
}

func TestExtractConfigurationsFromBootstrapJson(t *testing.T) {
	mockLogger()
	bootstrapJsonData := `invalidJsonStr`
	output := ExtractConfigurationsFromBootstrapJson([]byte(bootstrapJsonData), "myCollection", "myEnvironment")
	assert.Nil(t, output)
	if hook.LastEntry().Message != "AppConfiguration - Error occurred while reading bootstrap configurations - invalid character 'i' looking for beginning of value" {
		t.Errorf("Test failed: Incorrect error message")
	}

	bootstrapJsonData = `{"environments":[{"name":"Dev","environment_id":"dev","description":"","tags":"","color_code":"#FDD13A","features":[],"properties":[]}],"collections":[],"segments":[]}`
	output = ExtractConfigurationsFromBootstrapJson([]byte(bootstrapJsonData), "myCollection", "myEnvironment")
	assert.Nil(t, output)
	if hook.LastEntry().Message != "AppConfiguration - Error occurred while reading bootstrap configurations - no data matching for environment id: myEnvironment" {
		t.Errorf("Test failed: Incorrect error message")
	}

	bootstrapJsonData = `{"environments":[{"name":"Dev","environment_id":"myEnvironment","description":"","tags":"","color_code":"#FDD13A","features":[],"properties":[]}],"collections":[],"segments":[]}`
	output = ExtractConfigurationsFromBootstrapJson([]byte(bootstrapJsonData), "myCollection", "myEnvironment")
	assert.Nil(t, output)
	if hook.LastEntry().Message != "AppConfiguration - Error occurred while reading bootstrap configurations - no data matching for collection id: myCollection" {
		t.Errorf("Test failed: Incorrect error message")
	}

	bootstrapJsonData = `{"environments":[{"name":"My Environment","environment_id":"myEnvironment","description":"Environment created on instance creation","tags":"","color_code":"#FDD13A","features":[{"name":"F1","feature_id":"f1","description":"","tags":"","type":"NUMERIC","enabled_value":5,"disabled_value":0,"segment_rules":[{"rules":[{"segments":["mysegment"]}],"value":40,"order":1}],"collections":[{"collection_id":"myCollection","name":"My Collection"}],"enabled":true,"isOverridden":true}],"properties":[]}],"collections":[{"name":"My Collection","collection_id":"myCollection","description":"","tags":""}],"segments":[{"name":"test","segment_id":"test","description":"","tags":"","rules":[{"values":["test"],"operator":"startsWith","attribute_name":"test"}]}]}`
	output = ExtractConfigurationsFromBootstrapJson([]byte(bootstrapJsonData), "myCollection", "myEnvironment")
	assert.Nil(t, output)
	if hook.LastEntry().Message != "AppConfiguration - Error occurred while reading bootstrap configurations - no data matching for segment id: mysegment" {
		t.Errorf("Test failed: Incorrect error message")
	}

	bootstrapJsonData = `{"environments":[{"name":"My Environment","environment_id":"myEnvironment","description":"Environment created on instance creation","tags":"","color_code":"#FDD13A","features":[{"name":"F1","feature_id":"f1","description":"","tags":"","type":"NUMERIC","enabled_value":5,"disabled_value":0,"segment_rules":[{"rules":[{"segments":["l2dfo8do"]}],"value":40,"order":1}],"collections":[{"collection_id":"myCollection","name":"My Collection"}],"enabled":true,"isOverridden":true}],"properties":[{"name":"p1","property_id":"p1","description":"","tags":"","type":"NUMERIC","value":5,"segment_rules":[{"rules":[{"segments":["l2dfo8do"]}],"value":40,"order":1}],"collections":[{"collection_id":"myCollection","name":"My Collection"}],"isOverridden":true}]}],"collections":[{"name":"My Collection","collection_id":"myCollection","description":"","tags":""}],"segments":[{"name":"test","segment_id":"l2dfo8do","description":"","tags":"","rules":[{"values":["test"],"operator":"startsWith","attribute_name":"test"}]}]}`
	output = ExtractConfigurationsFromBootstrapJson([]byte(bootstrapJsonData), "myCollection", "myEnvironment")
	assert.NotNil(t, output)
	expectedConfigurations := `{"features":[{"name":"F1","feature_id":"f1","type":"NUMERIC","format":"","enabled_value":5,"disabled_value":0,"segment_rules":[{"rules":[{"segments":["l2dfo8do"]}],"value":40,"order":1,"rollout_percentage":null}],"enabled":true,"rollout_percentage":null}],"properties":[{"name":"p1","property_id":"p1","type":"NUMERIC","format":"","value":5,"segment_rules":[{"rules":[{"segments":["l2dfo8do"]}],"value":40,"order":1,"rollout_percentage":null}]}],"segments":[{"name":"test","segment_id":"l2dfo8do","rules":[{"values":["test"],"operator":"startsWith","attribute_name":"test"}]}]}`
	if !reflect.DeepEqual([]byte(expectedConfigurations), output) {
		t.Error("Expected TestExtractConfigurationsFromBootstrapJson test case to pass")
	}

}
func TestExtractConfigurationsFromAPIResponse(t *testing.T) {
	mockLogger()
	configRes := `invalidJsonStr`
	output := ExtractConfigurationsFromAPIResponse([]byte(configRes))
	assert.Nil(t, output)
	if hook.LastEntry().Message != "AppConfiguration - Error occurred while reading fetched configurations - invalid character 'i' looking for beginning of value" {
		t.Errorf("Test failed: Incorrect error message")
	}

	configRes = `{"environments":[{"name":"Dev","environment_id":"dev","features":[{"name":"Flight Booking Discounts","feature_id":"discount-on-flight-booking","type":"NUMERIC","enabled_value":15,"disabled_value":0,"segment_rules":[],"enabled":true,"rollout_percentage":100}],"properties":[{"name":"CI Pipeline","property_id":"ci-pipeline","type":"BOOLEAN","value":false,"segment_rules":[{"rules":[{"segments":["kk488156"]}],"value":true,"order":1}]}]}],"segments":[{"name":"spartans","segment_id":"kk488156","rules":[{"values":["bob@bluecharge.com","alice@bluecharge.com"],"operator":"is","attribute_name":"email"}]}]}`
	output = ExtractConfigurationsFromAPIResponse([]byte(configRes))
	assert.NotNil(t, output)
	expectedConfigurations := `{"features":[{"name":"Flight Booking Discounts","feature_id":"discount-on-flight-booking","type":"NUMERIC","format":"","enabled_value":15,"disabled_value":0,"segment_rules":[],"enabled":true,"rollout_percentage":100}],"properties":[{"name":"CI Pipeline","property_id":"ci-pipeline","type":"BOOLEAN","format":"","value":false,"segment_rules":[{"rules":[{"segments":["kk488156"]}],"value":true,"order":1,"rollout_percentage":null}]}],"segments":[{"name":"spartans","segment_id":"kk488156","rules":[{"values":["bob@bluecharge.com","alice@bluecharge.com"],"operator":"is","attribute_name":"email"}]}]}`
	if !reflect.DeepEqual([]byte(expectedConfigurations), output) {
		t.Error("Expected TestExtractConfigurationsFromAPIResponse test case to pass")
	}
}
