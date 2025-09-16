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

// CacheConfig : all the configurations pulled out from bootstrapJson/persistentCache/GitConfig/APIResponse
// will be stored in this struct format.
type CacheConfig struct {
	Features   []FeatureC  `json:"features"`
	Properties []PropertyC `json:"properties"`
	Segments   []Segment   `json:"segments"`
}

// Config : the format of all configurations.
type Config struct {
	Environments []Environment `json:"environments"`
	Collections  []Collection  `json:"collections"`
	Segments     []Segment     `json:"segments"`
}
type Environment struct {
	Name          string      `json:"name"`
	EnvironmentID string      `json:"environment_id"`
	Features      []FeatureC  `json:"features"`
	Properties    []PropertyC `json:"properties"`
}
type Collection struct {
	Name         string `json:"name,omitempty"`
	CollectionID string `json:"collection_id"`
}
type FeatureC struct {
	Feature
	Collections []Collection `json:"collections,omitempty"`
}
type PropertyC struct {
	Property
	Collections []Collection `json:"collections,omitempty"`
}
