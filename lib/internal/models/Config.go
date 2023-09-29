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

// Configurations : all the configurations pulled out from bootstrapJson/persistentCache/GitConfig/APIResponse
// will be stored in this struct format.
type Configurations struct {
	Features   []Feature  `json:"features"`
	Properties []Property `json:"properties"`
	Segments   []Segment  `json:"segments"`
}

// BootstrapConfig : the format of bootstrap configurations.
type BootstrapConfig struct {
	Environments []EnvironmentC `json:"environments"`
	Collections  []Collection   `json:"collections"`
	Segments     []Segment      `json:"segments"`
}
type EnvironmentC struct {
	Name          string      `json:"name"`
	EnvironmentID string      `json:"environment_id"`
	Features      []FeatureC  `json:"features"`
	Properties    []PropertyC `json:"properties"`
}
type Collection struct {
	Name         string `json:"name"`
	CollectionID string `json:"collection_id"`
}
type FeatureC struct {
	Feature
	Collections []Collection `json:"collections"`
}
type PropertyC struct {
	Property
	Collections []Collection `json:"collections"`
}

// APIConfig : the format of configurations returned from the API: `GET /config?action=sdkConfig`
type APIConfig struct {
	Environments []Environment `json:"environments"`
	Segments     []Segment     `json:"segments"`
}
type Environment struct {
	Name          string     `json:"name"`
	EnvironmentID string     `json:"environment_id"`
	Features      []Feature  `json:"features"`
	Properties    []Property `json:"properties"`
}
