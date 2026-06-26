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

package utils

import (
	"fmt"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	godsutils "github.com/emirpasic/gods/utils"
)

type RolloutPhase struct {
	Percentage   int    `json:"percentage"`
	Duration     int    `json:"duration,omitempty"`
	DurationType string `json:"duration_type,omitempty"`
}

type RolloutConfiguration struct {
	StartAt string         `json:"start_at"`
	Phases  []RolloutPhase `json:"phases"`
}

// ParseRolloutConfigurationPhases parses progressive rollout phases into a TreeMap
// for efficient timestamp-to-percentage lookups
func ParseRolloutConfigurationPhases(configuration *RolloutConfiguration) (*treemap.Map, error) {
	if configuration == nil || configuration.StartAt == "" || len(configuration.Phases) == 0 {
		return nil, fmt.Errorf("invalid rollout configuration")
	}

	// Parse start timestamp
	startTime, err := time.Parse(time.RFC3339, configuration.StartAt)
	if err != nil {
		return nil, fmt.Errorf("invalid start_at: %s, error: %v", configuration.StartAt, err)
	}

	// Create TreeMap with int64 comparator for timestamps
	result := treemap.NewWith(godsutils.Int64Comparator)
	result.Put(int64(0), 0)

	transitionTime := startTime.UnixMilli()

	// Duration multipliers in milliseconds
	multipliers := map[string]int64{
		"days":    86400000, // days
		"hours":   3600000,  // hours
		"minutes": 60000,    // minutes
	}

	for _, phase := range configuration.Phases {
		// Add phase entry
		result.Put(transitionTime, phase.Percentage)

		// Calculate next transition time if duration is specified
		if phase.Duration != 0 && phase.DurationType != "" {
			transitionTime += int64(phase.Duration) * multipliers[phase.DurationType]
		}
	}

	return result, nil
}

// GetCurrentRolloutPercentage returns the current rollout percentage based on current time
func GetCurrentRolloutPercentage(rolloutMap *treemap.Map) int {
	if rolloutMap == nil || rolloutMap.Empty() {
		return 0
	}

	currentTime := time.Now().UnixMilli()

	// Find the entry with the largest timestamp that is <= currentTime
	floorKey, floorValue := rolloutMap.Floor(currentTime)
	if floorKey != nil {
		if percentage, ok := floorValue.(int); ok {
			return percentage
		}
	}

	return 0
}
