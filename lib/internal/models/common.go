/**
 * (C) Copyright IBM Corp. 2023.
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
	"encoding/json"
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/utils"
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/utils/log"
	"github.com/spaolacci/murmur3"
	"math"
)

func computeHash(str string) float64 {
	seed := 0
	hasher := murmur3.New32WithSeed(uint32(seed))
	hasher.Write([]byte(str))
	return float64(hasher.Sum32())
}

func GetNormalizedValue(str string) int {
	maxHashValue := math.Pow(2, 32)
	normalizer := 100
	return int((computeHash(str) / maxHashValue) * float64(normalizer))
}

// ExtractConfigurationsFromBootstrapJson will
// 1. Parse the bootstrap configuration into models.BootstrapConfig struct
// 2. Extract all the Features, Properties & Segments that are under environmentId and assigned to collectionId into models.Configurations struct
// 3. And Marshals the models.Configurations
// Currently, for any type of errors that will occur, this method will not return error instead returns a nil.
// In the future, we will consider returning the error.
func ExtractConfigurationsFromBootstrapJson(bootstrapFileData []byte, collectionId, environmentId string) []byte {
	errMessage := `Error occurred while reading bootstrap configurations - `

	bootstrapConfig := BootstrapConfig{}
	err := json.Unmarshal(bootstrapFileData, &bootstrapConfig)
	if err != nil {
		log.Error(errMessage + err.Error())
		return nil
	}

	// validate environmentId & pick the environment
	var matchingEnv *EnvironmentC
	for _, env := range bootstrapConfig.Environments {
		if env.EnvironmentID == environmentId {
			matchingEnv = &env
			break
		}
	}
	if matchingEnv == nil {
		log.Error(errMessage+"no data matching for environment id: ", environmentId)
		return nil
	}

	// validate collectionId
	var matchingCol *Collection
	for _, col := range bootstrapConfig.Collections {
		if col.CollectionID == collectionId {
			matchingCol = &col
			break
		}
	}
	if matchingCol == nil {
		log.Error(errMessage+"no data matching for collection id: ", collectionId)
		return nil
	}
	// slice's to store the extracted features, properties & segments
	var features []Feature
	var properties []Property
	var segments []Segment

	var segmentIds []string
	uniqueSegmentIdsMap := make(map[string]bool)

	// loop through features in the matching environment, and pick all the features matching the collectionId
	for _, feature := range matchingEnv.Features {
		match := false
		for _, collection := range feature.Collections {
			if collection.CollectionID == collectionId {
				match = true
				break
			}
		}
		if match {
			features = append(features, feature.Feature)
			// get the segmentIds from the extracted feature. Use the segmentId to extract segments.
			for _, segmentRule := range feature.SegmentRules {
				for _, rules := range segmentRule.Rules {
					for _, segment := range rules.Segments {
						segmentIds = append(segmentIds, segment)
					}
				}
			}
		}
	}

	// loop through properties in the matching environment, and pick all the properties matching the collectionId
	for _, property := range matchingEnv.Properties {
		match := false
		for _, collection := range property.Collections {
			if collection.CollectionID == collectionId {
				match = true
				break
			}
		}
		if match {
			properties = append(properties, property.Property)
			// get the segmentIds from the extracted feature. Use the segmentId to extract segments.
			for _, segmentRule := range property.SegmentRules {
				for _, rules := range segmentRule.Rules {
					for _, segment := range rules.Segments {
						segmentIds = append(segmentIds, segment)
					}
				}
			}
		}
	}

	// segmentIds can contain duplicates. Filter them out using the map
	for _, value := range segmentIds {
		uniqueSegmentIdsMap[value] = true
	}
	uniqueSegmentIdsSlice := make([]string, 0, len(uniqueSegmentIdsMap))
	for value := range uniqueSegmentIdsMap {
		uniqueSegmentIdsSlice = append(uniqueSegmentIdsSlice, value)
	}

	// For all the unique segmentIds extract their segment
	for _, segmentId := range uniqueSegmentIdsSlice {
		var matchingSeg *Segment
		for _, segment := range bootstrapConfig.Segments {
			if segment.SegmentID == segmentId {
				matchingSeg = &segment
				break
			}
		}
		if matchingSeg == nil {
			log.Error(errMessage+"no data matching for segment id: ", segmentId)
			return nil
		}
		segments = append(segments, *matchingSeg)
	}

	c, err := json.Marshal(Configurations{
		Features:   features,
		Properties: properties,
		Segments:   segments,
	})
	// highly unlikely that err will not be nil.
	if err != nil {
		log.Error(errMessage + err.Error())
		return nil
	}
	return c
}

// ExtractConfigurationsFromAPIResponse will
// 1. Parse the configuration into models.APIConfig struct
// 2. Extract all the Features, Properties & Segments into models.Configurations struct
func ExtractConfigurationsFromAPIResponse(res []byte) []byte {
	defer utils.GracefullyHandleError()
	errMessage := `Error occurred while reading fetched configurations - `

	apiConfig := APIConfig{}
	err := json.Unmarshal(res, &apiConfig)
	if err != nil {
		log.Error(errMessage + err.Error())
		return nil
	}

	features := apiConfig.Environments[0].Features
	properties := apiConfig.Environments[0].Properties
	segments := apiConfig.Segments

	c, err := json.Marshal(Configurations{
		Features:   features,
		Properties: properties,
		Segments:   segments,
	})
	// highly unlikely that err will not be nil.
	if err != nil {
		log.Error(errMessage + err.Error())
		return nil
	}
	return c
}

// AliasFunction : Only for readability purpose.
// The configurations stored in Persistent cache & configuration fetched from API both are in same format
var ExtractConfigurationsFromPersistentCache = ExtractConfigurationsFromAPIResponse

// FormatConfig : will reformat the configurations from type Configurations to type APIConfig
func FormatConfig(data []byte, environmentId string) []byte {
	configurations := Configurations{}
	_ = json.Unmarshal(data, &configurations)

	reformatted := APIConfig{}
	reformatted.Environments = append(reformatted.Environments, Environment{})

	reformatted.Environments[0].Name = environmentId
	reformatted.Environments[0].EnvironmentID = environmentId
	reformatted.Environments[0].Features = configurations.Features
	reformatted.Environments[0].Properties = configurations.Properties
	reformatted.Segments = configurations.Segments

	c, _ := json.Marshal(reformatted)
	return c
}
