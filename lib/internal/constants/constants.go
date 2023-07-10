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

package constants

// DefaultSegmentID : Default Segment ID
const DefaultSegmentID = "$$null$$"

// DefaultEntityID : Default Entity ID
const DefaultEntityID = "$$null$$"

// DefaultUsageLimit : Default Usage Limit
const DefaultUsageLimit = 10

// UserAgent specifies the user agent name
const UserAgent = "appconfiguration-go-sdk/0.4.1"

// ConfigurationFile : Name of file to which configurations will be written
const ConfigurationFile = "appconfiguration.json"

// MaxNumberOfRetries : Maximum number of retries
const MaxNumberOfRetries = 3

// MaxRetryInterval : Maximum duration between successive retries (in seconds)
const MaxRetryInterval = 30

// StatusCodeOK : Http status code for successful GET call
const StatusCodeOK = 200

// StatusCodeAccepted : Http status code for successful POST call
const StatusCodeAccepted = 202

// StatusCodeClientErrorBegin : Beginning status code for client related errors
const StatusCodeClientErrorBegin = 400

// StatusCodeTooManyRequests : Http status code for API call exceeding rate limit
const StatusCodeTooManyRequests = 429

// StatusCodeClientErrorEnd : End status code for client related errors
const StatusCodeClientErrorEnd = 499

// StatusCodeServerErrorBegin : Beginning status code for server errors
const StatusCodeServerErrorBegin = 500

// StatusCodeServerErrorEnd : End status code for server errors
const StatusCodeServerErrorEnd = 599
