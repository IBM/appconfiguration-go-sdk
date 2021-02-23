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

type RuleElem struct {
	Segments []string
}
type SegmentRule struct {
	rules []RuleElem
	value interface{}
	order int
}

func (sr *SegmentRule) GetRules() []RuleElem {
	return sr.rules
}
func (sr *SegmentRule) GetValue() interface{} {
	return sr.value
}
func (sr *SegmentRule) GetOrder() int {
	return sr.order
}
