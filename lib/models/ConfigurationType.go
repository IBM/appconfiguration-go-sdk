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
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/messages"
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/utils/log"
	"gopkg.in/yaml.v3"
)

func IsValidDataType(category string) bool {
	switch category {
	case
		"NUMERIC",
		"BOOLEAN",
		"STRING":
		return true
	}
	return false
}

func getTypeCastedValue(val interface{}, valType string, valFormat string) interface{} {

	if valType == "NUMERIC" && isNumber(val) {
		return val.(float64)
	} else if valType == "BOOLEAN" && isBool(val) {
		return val.(bool)
	} else if valType == "STRING" {
		if valFormat == "TEXT" && isString(val) {
			return val.(string)
		} else if valFormat == "JSON" {
			return val
		} else if valFormat == "YAML" {
			// isString() is added to avoid multiple parsing of yaml value
			// if it is string, then only parse it to map. Else, it would have already parsed.
			if isString(val) {
				var result interface{}
				// TODO: support for multi-document yaml
				if err := yaml.Unmarshal([]byte(val.(string)), &result); err != nil {
					log.Error(messages.UnmarshalYAMLErr, err)
					return nil
				}
				return result
			}
			return val
		} else {
			log.Error(messages.InvalidDataFormat)
			return nil
		}
	} else {
		if !IsValidDataType(valType) {
			log.Error(messages.InvalidDataType, valType)
			return nil
		}
		log.Error(messages.TypeCastingError)
		return nil
	}
}
