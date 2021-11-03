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
	messages "github.com/IBM/appconfiguration-go-sdk/lib/internal/messages"
	utils "github.com/IBM/appconfiguration-go-sdk/lib/internal/utils"
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/utils/log"
)

// Segment : Segment struct
type Segment struct {
	Name      string `json:"name"`
	SegmentID string `json:"segment_id"`
	Rules     []Rule `json:"rules"`
}

// GetName : Get Name
func (s *Segment) GetName() string {
	return s.Name
}

// GetSegmentID : Get SegmentID
func (s *Segment) GetSegmentID() string {
	return s.SegmentID
}

// GetRules : Get Rules
func (s *Segment) GetRules() []Rule {
	return s.Rules
}

// EvaluateRule : Evaluate Rule
func (s *Segment) EvaluateRule(entityAttributes map[string]interface{}) bool {
	log.Debug(messages.EvalSegmentRule)
	defer utils.GracefullyHandleError()
	for _, rule := range s.GetRules() {
		if !rule.EvaluateRule(entityAttributes) {
			return false
		}
	}
	return true
}
