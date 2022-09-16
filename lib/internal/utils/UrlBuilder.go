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
	"github.com/IBM/go-sdk-core/v5/core"
	"net/http"
	"strings"
)

// URLBuilder : URLBuilder struct
type URLBuilder struct {
	baseURL               string
	wsPath                string
	path                  string
	service               string
	httpBase              string
	webSocketURL          string
	events                string
	region                string
	guid                  string
	privateEndpointPrefix string
	iamURL                string
	authenticator         core.Authenticator
}

var urlBuilderInstance *URLBuilder

// GetInstance : Get Instance
func GetInstance() *URLBuilder {
	if urlBuilderInstance == nil {
		urlBuilderInstance = &URLBuilder{
			baseURL:               ".apprapp.cloud.ibm.com",
			privateEndpointPrefix: "private.",
			wsPath:                "/wsfeature",
			path:                  "/feature/v1/instances/",
			service:               "/apprapp",
			events:                "/events/v1/instances/",
			httpBase:              "",
			webSocketURL:          "",
			iamURL:                "",
			region:                "",
			guid:                  "",
		}
	}
	return urlBuilderInstance
}

// Init : Init
func (ub *URLBuilder) Init(collectionID string, environmentID string, region string, guid string, apikey string, overrideServiceUrl string, usePrivateEndpoint bool) {
	ub.region = region
	ub.guid = guid

	// for dev & stage
	if len(overrideServiceUrl) > 0 {
		temp := strings.Split(overrideServiceUrl, "://")
		if usePrivateEndpoint {
			ub.httpBase = temp[0] + "://" + ub.privateEndpointPrefix + temp[1]
			ub.iamURL = "https://private.iam.test.cloud.ibm.com"
			ub.webSocketURL = "wss://" + ub.privateEndpointPrefix + temp[1] + ub.service + ub.wsPath + "?instance_id=" + guid + "&collection_id=" + collectionID + "&environment_id=" + environmentID
		} else {
			ub.httpBase = overrideServiceUrl
			ub.iamURL = "https://iam.test.cloud.ibm.com"
			ub.webSocketURL = "wss://" + temp[1] + ub.service + ub.wsPath + "?instance_id=" + guid + "&collection_id=" + collectionID + "&environment_id=" + environmentID

		}
		// for prod
	} else {
		if usePrivateEndpoint {
			ub.httpBase = "https://" + ub.privateEndpointPrefix + region + ub.baseURL
			ub.iamURL = "https://private.iam.cloud.ibm.com"
			ub.webSocketURL = "wss://" + ub.privateEndpointPrefix + region + ub.baseURL + ub.service + ub.wsPath + "?instance_id=" + guid + "&collection_id=" + collectionID + "&environment_id=" + environmentID
		} else {
			ub.httpBase = "https://" + region + ub.baseURL
			ub.iamURL = "https://iam.cloud.ibm.com"
			ub.webSocketURL = "wss://" + region + ub.baseURL + ub.service + ub.wsPath + "?instance_id=" + guid + "&collection_id=" + collectionID + "&environment_id=" + environmentID
		}
	}

	// Create the authenticator.
	var err error
	ub.authenticator, err = core.NewIamAuthenticatorBuilder().
		SetApiKey(apikey).
		SetURL(ub.iamURL).
		Build()
	if err != nil {
		panic(err)
	}
}

// GetBaseServiceURL returns base service url
func (ub *URLBuilder) GetBaseServiceURL() string {
	return ub.httpBase
}

// SetBaseServiceURL overrides the base service url if set
func (ub *URLBuilder) SetBaseServiceURL(url string) {
	ub.httpBase = url
}

// GetAuthenticator returns iam authenticator
func (ub *URLBuilder) GetAuthenticator() core.Authenticator {
	return ub.authenticator
}

// GetWebSocketURL returns web socket url
func (ub *URLBuilder) GetWebSocketURL() string {
	return ub.webSocketURL
}

// GetToken returns the string "Bearer <token>"
func (ub *URLBuilder) GetToken() string {
	req, _ := http.NewRequest("GET", "https://localhost", nil)
	var err error
	err = ub.authenticator.Authenticate(req)
	if err != nil {
		return ""
	}
	return req.Header.Get("Authorization")
}

// SetWebSocketURL : sets web socket url
func (ub *URLBuilder) SetWebSocketURL(webSocketURL string) {
	ub.webSocketURL = webSocketURL
}

// SetAuthenticator : assigns an authenticator to the url builder instance authenticator member variable.
func (ub *URLBuilder) SetAuthenticator(authenticator core.Authenticator) {
	ub.authenticator = authenticator
}
