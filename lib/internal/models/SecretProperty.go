/**
 * (C) Copyright IBM Corp. 2022.
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
	"errors"

	"github.com/IBM/appconfiguration-go-sdk/lib/internal/messages"
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/utils/log"
	"github.com/IBM/go-sdk-core/v5/core"
	sm "github.com/IBM/secrets-manager-go-sdk/secretsmanagerv1"
)

// SecretProperty : SecretProperty struct
type SecretProperty struct {
	PropertyID string
}

// GetCurrentValue returns the actual secret value(default or overridden) based on the evaluation.
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
func (sp *SecretProperty) GetCurrentValue(entityID string, entityAttributes ...map[string]interface{}) (result *sm.GetSecret, response *core.DetailedResponse, err error) {

	if len(entityID) <= 0 {
		log.Error("SecretProperty evaluation: ", messages.InvalidEntityId, "GetCurrentValue")
		return nil, nil, errors.New("error: " + messages.InvalidEntityId + "GetCurrentValue")
	}

	if len(entityAttributes) > 1 {
		log.Error("SecretProperty evaluation: ", messages.IncorrectUsageOfEntityAttributes, "GetCurrentValue")
		return nil, nil, errors.New("error: " + messages.IncorrectUsageOfEntityAttributes + "SecretProperty GetCurrentValue")
	}

	propertyObject := GetCacheInstance().PropertyMap[sp.PropertyID]
	propertySecretValue := propertyObject.GetValue().(map[string]interface{})
	var propertySecretType string
	if secretTypeValue, secretTypeExist := propertySecretValue["secret_type"]; secretTypeExist {
		propertySecretType = secretTypeValue.(string)
	} else {
		//secret_type does not exist then throw error
		log.Error("SecretProperty evaluation: secret_type is missing from the Property value of:", propertyObject.GetPropertyName())
		return nil, nil, errors.New("error: secret_type is missing from the Property value of :" + propertyObject.GetPropertyName())
	}

	var propertyCurrentVal interface{}
	if entityAttributes == nil {
		propertyCurrentVal = propertyObject.GetCurrentValue(entityID)
	} else {
		propertyCurrentVal = propertyObject.GetCurrentValue(entityID, entityAttributes[0])
	}

	if propertyCurrentVal == nil {
		log.Error(messages.InvalidPropertyValueMessage)
		return nil, nil, errors.New("error: " + messages.InvalidPropertyValueMessage)
	}

	valMap, isTypeMap := propertyCurrentVal.(map[string]interface{})
	if !isTypeMap {
		return nil, nil, errors.New("error: " + messages.InvalidPropertyValueMessage)
	}
	if secretID, secretIDExist := valMap["id"]; secretIDExist {
		id := secretID.(string)
		//sm sdk call
		secretManager := GetCacheInstance().SecretManagerMap[sp.PropertyID].(*sm.SecretsManagerV1)
		secretData, detailedResp, err := secretManager.GetSecret(&sm.GetSecretOptions{
			SecretType: core.StringPtr(propertySecretType),
			ID:         &id,
		})

		if err != nil {
			return nil, nil, err
		}
		return secretData, detailedResp, err
	}
	log.Error(messages.InvalidSecretID)
	return nil, nil, errors.New("error: " + messages.InvalidSecretID)
}
