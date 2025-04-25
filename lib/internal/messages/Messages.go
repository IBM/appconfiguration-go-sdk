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

package messages

// RegionError : RegionError const
const RegionError = "Provide a valid region."

// GUIDError : GUIDError const
const GUIDError = "Provide a valid guid."

// ApikeyError : ApikeyError const
const ApikeyError = "Provide a valid apiKey."

// CollectionIDValueError : CollectionIDValueError const
const CollectionIDValueError = "Provide a valid collectionId."

// EnvironmentIDValueError : EnvironmentIDValueError const
const EnvironmentIDValueError = "Provide a valid environmentId."

// CollectionIDError : CollectionIDError const
const CollectionIDError = "Invalid action. You can perform this action only after a successful initialization. Check the initialization section for errors."

// CollectionInitError : CollectionInitError const
const CollectionInitError = "Invalid action. You can perform this action only after a successful initialization and setting the context. Check the Init and SetContext section for errors."

// BootstrapFileNotFoundError : BootstrapFileNotFoundError const
const BootstrapFileNotFoundError = "Provide bootstrap_file value when live_config_update_enabled is false."

// InvalidBootstrapFile : InvalidBootstrapFile const
const InvalidBootstrapFile = "Invalid value provided for BootstrapFile parameter"

// IncorrectUsageOfContextOptions : IncorrectUsageOfContextOptions const
const IncorrectUsageOfContextOptions = "Incorrect usage of context options. At most of one ContextOptions struct should be passed."

// IncorrectUsageOfEntityAttributes : IncorrectUsageOfEntityAttributes const
const IncorrectUsageOfEntityAttributes = "Incorrect usage of entityAttributes in "

// ConfigAPIError : ConfigAPIError const
const ConfigAPIError = "Failed to fetch the configurations. "

// FetchFromAPISdkInitError : FetchFromAPISdkInitError const
const FetchFromAPISdkInitError = "fetchFromAPI() - Configuration SDK not initialized with call to init."

// WebsocketErrorReadingMessage : WebsocketErrorReadingMessage const
const WebsocketErrorReadingMessage = "Error while reading message from the socket. "

// WebsocketReceivingMessage : WebsocketReceivingMessage const
const WebsocketReceivingMessage = "Message received from socket. "

// ConfigurationUpdateListenerMethodError : ConfigurationUpdateListenerMethodError const
const ConfigurationUpdateListenerMethodError = "Configuration update listener should me a method or a function."

// InvalidEntityId : InvalidEntityId const
const InvalidEntityId = "Invalid entityId passed to "

// ConfigurationHandlerInitError : ConfigurationHandlerInitError const
const ConfigurationHandlerInitError = "Invalid action in ConfigurationHandler. You can perform this action only after a successful initialization. Check the initialization section for errors."

// CreatingNewAPIManagerInstance : CreatingNewAPIManagerInstance const
const CreatingNewAPIManagerInstance = "Creating new API manager instance."

// ParsingFeatureRules : ParsingFeatureRules const
const ParsingFeatureRules = "Parsing feature rules."

// EvaluatingSegments : EvaluatingSegments const
const EvaluatingSegments = "Evaluating segments."

// FeatureValue : FeatureValue const
const FeatureValue = "Feature value is "

// EvaluatingFeature : EvaluatingFeature const
const EvaluatingFeature = "Evaluating feature."

// RetrievingFeature : RetrievingFeature const
const RetrievingFeature = "Retrieving feature current value."

// InvalidFeatureID : InvalidFeatureID const
const InvalidFeatureID = "Invalid feature id - "

// ErrorInvalidFeatureID : ErrorInvalidFeatureID const
const ErrorInvalidFeatureID = "error : invalid feature id "

// ErrorInvalidFeatureAction : ErrorInvalidFeatureAction const
const ErrorInvalidFeatureAction = "error : feature object not initialized"

// ParsingPropertyRules : ParsingPropertyRules const
const ParsingPropertyRules = "Parsing property rules."

// PropertyValue : PropertyValue const
const PropertyValue = "Property value is "

// EvaluatingProperty : EvaluatingProperty const
const EvaluatingProperty = "Evaluating property."

// RetrievingProperty : RetrievingProperty const
const RetrievingProperty = "Retrieving property current value."

// InvalidPropertyID : InvalidPropertyID const
const InvalidPropertyID = "Invalid property id - "

// ErrorInvalidPropertyID : ErrorInvalidPropertyID const
const ErrorInvalidPropertyID = "error : invalid property id "

// ErrorInvalidPropertyAction : ErrorInvalidPropertyAction const
const ErrorInvalidPropertyAction = "error : property object not initialized"

// EvalSegmentRule : EvalSegmentRule const
const EvalSegmentRule = "Evaluating segment rule."

// ExecAPICall : ExecAPICall const
const ExecAPICall = "Executing API call."

// APICallError : APICallError const
const APICallError = "Error executing API call "

// EncodeJSONErr : EncodeJSONErr const
const EncodeJSONErr = "Error while encoding json "

// WriteFileErr : WriteFileErr const
const WriteFileErr = "Error while writing to a file "

// ReadFileErr : ReadFileErr const
const ReadFileErr = "Error while reading file "

// StoreFile : StoreFile const
const StoreFile = "Storing file."

// ReadFile : ReadFile const
const ReadFile = "Reading file."

// RetrieveMeteringInstance : RetrieveMeteringInstance const
const RetrieveMeteringInstance = "Retrieving metering instance."

// StartSendingMeteringData : StartSendingMeteringData const
const StartSendingMeteringData = "Start sending metering data in the background."

// AddMetering : AddMetering const
const AddMetering = "Add metering."

// RecordEval : RecordEval const
const RecordEval = "Record evaluation."

// TenMinExpiry : TenMinExpiry const
const TenMinExpiry = "10 mins expired, sending metering data if any."

// SendMeteringServer : SendMeteringServer const
const SendMeteringServer = "Sending to metering server."

// SendMeteringSuccess : SendMeteringSuccess const
const SendMeteringSuccess = "Successfully sent metering data to server."

// SendMeteringServerErr : SendMeteringServerErr const
const SendMeteringServerErr = "Error while sending metering data to server. "

// RetrieveingAppConfig : RetrieveingAppConfig const
const RetrieveingAppConfig = "Retrieving App Configuration instance."

// CreatingNewAppConfig : CreatingNewAppConfig const
const CreatingNewAppConfig = "Creating new App Configuration instance."

// AppConfigAlreadyInstantiated : AppConfigAlreadyInstantiated const
const AppConfigAlreadyInstantiated = "App Configuration instance is already instantiated."

// SettingContext : SettingContext const
const SettingContext = "Setting context."

// LoadingData : LoadingData const
const LoadingData = "Loading data."

// ReadPersistentCache : ReadPersistentCache const
const ReadPersistentCache = "Reading configurations from persistent cache: "

// CheckBootstrapFileProvided : CheckBootstrapFileProvided const
const CheckBootstrapFileProvided = "Checking bootstrap file is provided by the user or not."

// ReadBootstrapConfigurations : ReadBootstrapConfigurations const
const ReadBootstrapConfigurations = "Reading configurations from bootstrap file: "

// LoadingConfigurations : LoadingConfigurations const
const LoadingConfigurations = "Loading configurations."

// LiveUpdateCheck : LiveUpdateCheck const
const LiveUpdateCheck = "Checking live configuration update is enabled or not."

// FetchFromConfigurationFile : FetchFromConfigurationFile const
const FetchFromConfigurationFile = "Fetching from configuration file."

// FetchConfigurationData : FetchConfigurationData const
const FetchConfigurationData = "Fetching configuration data."

// FetchFromAPI : FetchFromAPI const
const FetchFromAPI = "Fetching from API."

// StartWebSocket : StartWebSocket const
const StartWebSocket = "Starting web socket connection."

// WebSocketConnectErr : WebSocketConnectErr const
const WebSocketConnectErr = "Error occurred while establishing websocket connection with the App Configuration server. "

// RetryWebSocketConnect : RetryWebSocketConnect const
const RetryWebSocketConnect = "Trying web socket connection again."

// UnmarshalJSONErr : UnmarshalJSONErr const
const UnmarshalJSONErr = "Error while unmarshalling JSON "

// UnmarshalYAMLErr : UnmarshalYAMLErr const
const UnmarshalYAMLErr = "Error while unmarshalling YAML "

// MarshalJSONErr : MarshalJSONErr const
const MarshalJSONErr = "Error while marshalling JSON "

// SetInMemoryCache : SetInMemoryCache const
const SetInMemoryCache = "Setting memory cache."

// ConfigurationFileEmpty : ConfigurationFileEmpty const
const ConfigurationFileEmpty = " file is empty."

// InitError : Caused due to initialization error
const InitError = "error: configurations not fetched, check the init and setcontext section for errors"

// InvalidDataType : Invalid Datatype
const InvalidDataType = "Invalid datatype: "

// InvalidDataFormat : Invalid Data Format
const InvalidDataFormat = "Invalid data format: "

// TypeCastingError : Type Casting Error
const TypeCastingError = "Error Type casting. Check the feature or property values."

// FetchAPISuccessful : FetchAPISuccessful const
const FetchAPISuccessful = "Successfully fetched the configurations."

// WebSocketConnectFailed : WebSocketConnectFailed const
const WebSocketConnectFailed = "Socket connect failed. "

// WebSocketConnectSuccess: WebSocketConnectSuccess const
const WebSocketConnectSuccess = "Socket connect successful."

// AuthTokenError : AuthTokenError const
const AuthTokenError = "Could not generate Bearer token for the provided API Key."

// RetryScheduledMessage : RetryScheduledMessage const
const RetryScheduledMessage = " Scheduling the API request to retry after %d minutes."

// InvalidPropertyValueMessage : InvalidPropertyValueMessage const
const InvalidPropertyValueMessage = "Property Value is either invalid or empty."

// InvalidSecretManagerMessage : InvalidSecretManagerMessage const
const InvalidSecretManagerMessage = "Secret Manager object is either invalid or empty."

// InvalidSecretID : InvalidSecretIdMessage const
const InvalidSecretID = "Secret Id is either invalid or empty."
