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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorHandling(t *testing.T) {
	// do a division with error handling support
	result := divideWithErrorHandling(1, 0)
	assert.Equal(t, 0, result)

	// do a division with error handling support
	assert.Panics(t, func() { divideWithoutErrorHandling(1, 0) }, "The code did not panic")
}

func divideWithErrorHandling(m int, n int) int {
	defer GracefullyHandleError()
	return m / n
}

func divideWithoutErrorHandling(m int, n int) int {
	return m / n
}
