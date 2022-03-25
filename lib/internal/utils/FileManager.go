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
	"encoding/json"
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/constants"
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/messages"
	"github.com/IBM/appconfiguration-go-sdk/lib/internal/utils/log"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

var fileMutex sync.Mutex

// SanitizePath : Sanitises the string path and restrict users from providing manipulated path
// Example :
//		input: ../../../etc/abc.conf
//		output: /etc/abc.conf
func SanitizePath(_path string) string {
	return filepath.FromSlash(path.Clean("/" + strings.Trim(_path, "/")))
}

// StoreFiles : Store Files
func StoreFiles(content, basePath string) {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	log.Debug(messages.StoreFile)

	file, err := json.MarshalIndent(json.RawMessage(content), "", "\t")
	if err != nil {
		log.Error(messages.EncodeJSONErr, err)
		return
	}
	sanitizedFilePath := filepath.Join(SanitizePath(basePath), constants.ConfigurationFile)
	err = ioutil.WriteFile(sanitizedFilePath, file, 0644)
	if err != nil {
		log.Error(messages.WriteFileErr, err)
		return
	}
}

// ReadFiles reads file from the file path
func ReadFiles(filePath string) []byte {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	log.Debug(messages.ReadFile)
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Error(messages.ReadFileErr, err)
		return []byte(`{}`)
	}
	return file
}
