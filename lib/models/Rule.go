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
	"reflect"
	"strconv"
	"strings"

	utils "github.com/IBM/appconfiguration-go-sdk/lib/internal/utils"
)

// Rule : Rule struct
type Rule struct {
	Values        []interface{} `json:"values"`
	Operator      string
	AttributeName string `json:"Attribute_Name"`
}

// GetValues : Get Values
func (r *Rule) GetValues() []interface{} {
	return r.Values
}

// GetOperator : Get Operator
func (r *Rule) GetOperator() string {
	return r.Operator
}

// GetAttributeName : Get Attribute Name
func (r *Rule) GetAttributeName() string {
	return r.AttributeName
}

func (r *Rule) operatorCheck(key interface{}, value interface{}) bool {

	var result bool = false

	if key == nil || value == nil {
		return result
	}

	switch r.GetOperator() {
	case "endsWith":
		result = strings.HasSuffix(key.(string), value.(string))
		break
	case "startsWith":
		result = strings.HasPrefix(key.(string), value.(string))
		break
	case "contains":
		result = strings.Contains(key.(string), value.(string))
		break
	case "is":
		if isNumber(key) {
			// compare number
			key, _ = getFloat(key)
			value, _ = strconv.ParseFloat(value.(string), 64)
			result = (key.(float64) == value.(float64))
		} else if isBool(key) {
			// compare boolean
			key, _ = formatBool(key) //convert boolean true/false to string "true"/"false"
			result = (key == value.(string))
		} else {
			// compare string
			result = (key == value)
		}
		break
	case "greaterThan":
		if isNumber(key) {
			key, _ = getFloat(key)
			value, _ = strconv.ParseFloat(value.(string), 64)
			result = key.(float64) > value.(float64)
		} else if isString(key) {
			key, _ = strconv.ParseFloat(key.(string), 64)
			value, _ = strconv.ParseFloat(value.(string), 64)
			result = key.(float64) > value.(float64)
		}
		break
	case "lesserThan":
		if isNumber(key) {
			key, _ = getFloat(key)
			value, _ = strconv.ParseFloat(value.(string), 64)
			result = key.(float64) < value.(float64)
		} else if isString(key) {
			key, _ = strconv.ParseFloat(key.(string), 64)
			value, _ = strconv.ParseFloat(value.(string), 64)
			result = key.(float64) < value.(float64)
		}
		break
	case "greaterThanEquals":
		if isNumber(key) {
			key, _ = getFloat(key)
			value, _ = strconv.ParseFloat(value.(string), 64)
			result = key.(float64) >= value.(float64)
		} else if isString(key) {
			key, _ = strconv.ParseFloat(key.(string), 64)
			value, _ = strconv.ParseFloat(value.(string), 64)
			result = key.(float64) >= value.(float64)
		}
		break
	case "lesserThanEquals":
		if isNumber(key) {
			key, _ = getFloat(key)
			value, _ = strconv.ParseFloat(value.(string), 64)
			result = key.(float64) <= value.(float64)
		} else if isString(key) {
			key, _ = strconv.ParseFloat(key.(string), 64)
			value, _ = strconv.ParseFloat(value.(string), 64)
			result = key.(float64) <= value.(float64)
		}
		break
	default:
		result = false
	}
	return result
}
func isNumber(val interface{}) bool {
	switch val.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return true
	default:
		return false
	}
}
func isBool(val interface{}) bool {
	switch val.(type) {
	case bool:
		return true
	default:
		return false
	}
}
func formatBool(val interface{}) (string, error) {
	if val == true {
		return "true", nil
	}
	return "false", nil
}
func isString(val interface{}) bool {
	return reflect.TypeOf(val).String() == "string"
}
func getFloat(unk interface{}) (float64, error) {
	switch i := unk.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int16:
		return float64(i), nil
	case int8:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case uint16:
		return float64(i), nil
	case uint8:
		return float64(i), nil
	case int:
		return float64(i), nil
	case uint:
		return float64(i), nil
	default:
		return float64(0), nil
	}
}

// EvaluateRule : Evaluate Rule
func (r *Rule) EvaluateRule(entityAttributes map[string]interface{}) bool {
	defer utils.GracefullyHandleError()
	var result = false
	key, ok := entityAttributes[r.GetAttributeName()]
	if !ok {
		return false
	}
	for _, val := range r.GetValues() {
		if r.operatorCheck(key, val) {
			result = true
		}
	}
	return result
}
