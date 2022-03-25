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
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/constants"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestFileManager(t *testing.T) {

	mockLogger()
	assert.Equal(t, SanitizePath(""), "/")
	assert.Equal(t, SanitizePath("Users/home/Desktop"), "/Users/home/Desktop")
	assert.Equal(t, SanitizePath("/Users/home/Desktop"), "/Users/home/Desktop")
	assert.Equal(t, SanitizePath("./Users/home/Desktop"), "/Users/home/Desktop")
	assert.Equal(t, SanitizePath("../../../etc/abc.conf"), "/etc/abc.conf")
	assert.Equal(t, SanitizePath("////../../Users/home/Desktop"), "/Users/home/Desktop")
	assert.Equal(t, SanitizePath("./Users/home/Desktop/abc/../abc1"), "/Users/home/Desktop/abc1")

	// TestStoreFilesWithValidJSONContent
	dir, _ := os.Getwd()
	StoreFiles(`{"key":"value"}`, dir)
	assert.EqualValues(t,
		string(ReadFiles(filepath.Join(SanitizePath(dir), constants.ConfigurationFile))), "{\n\t\"key\": \"value\"\n}")
	os.Remove(constants.ConfigurationFile)

	// TestStoreFilesWithInvalidJSONContent
	StoreFiles("", dir)
	if hook.LastEntry().Message != "AppConfiguration - Error while encoding json json: error calling MarshalJSON for type json.RawMessage: unexpected end of JSON input" {
		t.Errorf("Test failed: StoreFiles for Invalid json")
	}

	// TestReadFilesWithNonExistingFile
	assert.Equal(t, ReadFiles(SanitizePath("non-existing-file.txt")), []byte(`{}`))
}
