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
	cons "github.com/IBM/appconfiguration-go-sdk/lib/internal/constants"
	"github.com/IBM/go-sdk-core/v5/core"
	"sync"
	"time"
)

// APIManager : wrapper struct over core base service.
type APIManager struct {
	baseService    *core.BaseService
	serviceOptions *core.ServiceOptions
}

var apiManagerInstance *APIManager
var once sync.Once

// GetAPIManagerInstance : returns APIManager instance.
func GetAPIManagerInstance() *APIManager {
	once.Do(func() {
		apiManagerInstance = &APIManager{}
		apiManagerInstance.serviceOptions = &core.ServiceOptions{
			URL:           urlBuilderInstance.GetBaseServiceURL(),
			Authenticator: urlBuilderInstance.GetAuthenticator(),
		}
		apiManagerInstance.baseService, _ = core.NewBaseService(apiManagerInstance.serviceOptions)
		apiManagerInstance.baseService.EnableRetries(cons.MaxNumberOfRetries, time.Second*time.Duration(cons.MaxRetryInterval))
	})
	return apiManagerInstance
}

// Request : wrapper over core base service request method.
func (ap *APIManager) Request(builder *core.RequestBuilder) *core.DetailedResponse {
	request, err := builder.Build()
	if err != nil {
		return nil
	}
	var rawResponse map[string]json.RawMessage
	response, _ := ap.baseService.Request(request, &rawResponse)
	return response
}
