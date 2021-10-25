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
	"os"

	"github.com/IBM/appconfiguration-go-sdk/lib/internal/utils/log"
)

// Cache : Cache struct
type Cache struct {
	FeatureMap  map[string]Feature
	PropertyMap map[string]Property
	SegmentMap  map[string]Segment
}

// CacheInstance : Cache Instance
var CacheInstance *Cache

// Init : init
func init() {
	if os.Getenv("ENABLE_DEBUG") == "true" {
		log.SetLogLevel("debug")
	} else {
		log.SetLogLevel("info")
	}
}

// SetCache : Set Cache
func SetCache(featureMap map[string]Feature, propertyMap map[string]Property, segmentMap map[string]Segment) {
	CacheInstance = new(Cache)
	CacheInstance.FeatureMap = featureMap
	CacheInstance.PropertyMap = propertyMap
	CacheInstance.SegmentMap = segmentMap
	log.Debug(CacheInstance)
}

// GetCacheInstance : Get Cache Instance
func GetCacheInstance() *Cache {
	return CacheInstance
}
