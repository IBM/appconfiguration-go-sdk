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
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/constants"
	utils "github.com/IBM/appconfiguration-go-sdk/lib/internal/utils"
)

// RuleElem : RuleElem struct
type RuleElem struct {
	Segments []string `json:"segments"`
}

// SegmentRule : SegmentRule struct
type SegmentRule struct {
	Rules                []RuleElem                  `json:"rules"`
	Value                interface{}                 `json:"value"`
	Order                int                         `json:"order"`
	RolloutPercentage    *interface{}                `json:"rollout_percentage"`
	RuleID               *string                     `json:"rule_id,omitempty"`
	RolloutType          string                      `json:"rollout_type,omitempty"`
	RolloutConfiguration *utils.RolloutConfiguration `json:"rollout_configuration,omitempty"`
}

// GetRuleID : Get Rule ID
func (sr *SegmentRule) GetRuleID() string {
	if sr.RuleID == nil {
		return ""
	} else {
		return *sr.RuleID
	}
}

// GetRules : Get Rules
func (sr *SegmentRule) GetRules() []RuleElem {
	return sr.Rules
}

// GetValue : Get Value
func (sr *SegmentRule) GetValue() interface{} {
	return sr.Value
}

// GetOrder : Get Order
func (sr *SegmentRule) GetOrder() int {
	return sr.Order
}

// GetRolloutPercentage : Get the rollout percentage of the segment rule
func (sr *SegmentRule) GetRolloutPercentage() interface{} {
	if sr.RolloutPercentage == nil {
		var v interface{} = 100.0
		sr.RolloutPercentage = &v
	}
	return *sr.RolloutPercentage
}

// GetRolloutType : Get Rollout Type of Segment Rule
func (sr *SegmentRule) GetRolloutType() string {
	if sr.RolloutType == "" {
		return constants.Manual
	} else {
		return sr.RolloutType
	}
}
