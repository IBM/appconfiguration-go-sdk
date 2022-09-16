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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLBuilder(t *testing.T) {

	// test websocket url
	urlBuilder := GetInstance()
	urlBuilder.SetWebSocketURL("wss://test-service.com/apprapp/wsfeature?instance_id=guid&collection_id=CollectionID&environment_id=EnvironmentID")
	assert.Equal(t, "wss://test-service.com/apprapp/wsfeature?instance_id=guid&collection_id=CollectionID&environment_id=EnvironmentID", urlBuilder.GetWebSocketURL())
	resetURLBuilderInstance()

	// test base service url
	urlBuilder = GetInstance()
	urlBuilder.SetBaseServiceURL("https://region.apprapp.cloud.ibm.com")
	assert.Equal(t, "https://region.apprapp.cloud.ibm.com", urlBuilder.GetBaseServiceURL())
	resetURLBuilderInstance()

	// test private endpoint url (prod endpoints)
	urlBuilder = GetInstance()
	urlBuilder.Init("collectionId", "environmentId", "region", "guid", "apikey", "", true)
	assert.Equal(t, "https://private.region.apprapp.cloud.ibm.com", urlBuilder.GetBaseServiceURL())
	assert.Equal(t, "wss://private.region.apprapp.cloud.ibm.com/apprapp/wsfeature?instance_id=guid&collection_id=collectionId&environment_id=environmentId", urlBuilder.GetWebSocketURL())
	resetURLBuilderInstance()

	// test private endpoint url (dev & stage endpoints)
	urlBuilder = GetInstance()
	urlBuilder.Init("collectionId", "environmentId", "region", "guid", "apikey", "https://region.apprapp.test.cloud.ibm.com", true)
	assert.Equal(t, "https://private.region.apprapp.test.cloud.ibm.com", urlBuilder.GetBaseServiceURL())
	assert.Equal(t, "wss://private.region.apprapp.test.cloud.ibm.com/apprapp/wsfeature?instance_id=guid&collection_id=collectionId&environment_id=environmentId", urlBuilder.GetWebSocketURL())
	resetURLBuilderInstance()

	// test when get token encounters an error while retrieving token and returns an token of size 0
	urlBuilder = GetInstance()
	urlBuilder.SetAuthenticator(&core.NoAuthAuthenticator{})
	token := urlBuilder.GetToken()
	assert.Equal(t, 0, len(token))
	resetURLBuilderInstance()

}

func resetURLBuilderInstance() {
	urlBuilderInstance = nil
}
