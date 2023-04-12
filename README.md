# IBM Cloud App Configuration Go server SDK 0.3.3

IBM Cloud App Configuration SDK is used to perform feature flag and property evaluation based on the configuration on
IBM Cloud App Configuration service.

## Table of Contents

- [Overview](#overview)
- [Installation](#installation)
    - [`go get` command](#go-get-command)
    - [Go modules](#go-modules)
- [Using the SDK](#using-the-sdk)
- [License](#license)

## Overview

IBM Cloud App Configuration is a centralized feature management and configuration service
on [IBM Cloud](https://www.cloud.ibm.com) for use with web and mobile applications, microservices, and distributed
environments.

Instrument your applications with App Configuration Go SDK, and use the App Configuration dashboard, API or CLI to
define feature flags or properties, organized into collections and targeted to segments. Change feature flag states in
the cloud to activate or deactivate features in your application or environment, when required. You can also manage the
properties for distributed applications centrally.

## Prerequisites
- Go version 1.16 or newer


## Installation

**Note: The v1.x.x versions of the App Configuration Go SDK have been retracted. Use the latest available version of the SDK.**

There are a few different ways to download and install the IBM App Configuration Go SDK project for use by your Go
application:

#### `go get` command

Use this command to download and install the SDK (along with its dependencies) to allow your Go application to use it:

```
go get -u github.com/IBM/appconfiguration-go-sdk@latest
```

#### Go modules

If your application is using Go modules, you can add a suitable import to your Go application, like this:

```go
import (
	AppConfiguration "github.com/IBM/appconfiguration-go-sdk/lib"
)
```

then run `go mod tidy` to download and install the new dependency and update your Go application's go.mod file.

## Using the SDK

Initialize the sdk to connect with your App Configuration service instance.

```go
collectionId := "airlines-webapp"
environmentId := "dev"

appConfigClient := AppConfiguration.GetInstance()
appConfigClient.Init("region", "guid", "apikey")
appConfigClient.SetContext(collectionId, environmentId)
```

:red_circle: **Important** :red_circle:

The **`Init()`** and **`SetContext()`** are the initialisation methods and should be invoked **only once** using
appConfigClient. The appConfigClient, once initialised, can be obtained across modules
using **`AppConfiguration.GetInstance()`**. [See this example below](#fetching-the-appconfigclient-across-other-modules).

- region : Region name where the App Configuration service instance is created. Use
    - `AppConfiguration.REGION_US_SOUTH` for Dallas
    - `AppConfiguration.REGION_EU_GB` for London
    - `AppConfiguration.REGION_AU_SYD` for Sydney
    - `AppConfiguration.REGION_US_EAST` for Washington DC
- guid : Instance Id of the App Configuration service. Obtain it from the service credentials section of the App
  Configuration dashboard.
- apikey : ApiKey of the App Configuration service. Obtain it from the service credentials section of the App
  Configuration dashboard.
* collectionId: Id of the collection created in App Configuration service instance under the **Collections** section.
* environmentId: Id of the environment created in App Configuration service instance under the **Environments** section.

### Connect using private network connection (optional)

Set the SDK to connect to App Configuration service by using a private endpoint that is accessible only through the IBM
Cloud private network.

```go
appConfigClient.UsePrivateEndpoint(true)
```

This must be done before calling the `Init` function on the SDK.

### (Optional)

In order for your application and SDK to continue its operations even during the unlikely scenario of App Configuration
service across your application restarts, you can configure the SDK to work using a persistent cache. The SDK uses the
persistent cache to store the App Configuration data that will be available across your application restarts.

```go
// 1. default (without persistent cache)
appConfigClient.SetContext(collectionId, environmentId)

// 2. optional (with persistent cache)
appConfigClient.SetContext(collectionId, environmentId, AppConfiguration.ContextOptions{
    PersistentCacheDirectory: "/var/lib/docker/volumes/",
})
```

* PersistentCacheDirectory: Absolute path to a directory which has read & write permission for the user. The SDK will
  create a file - `appconfiguration.json` in the specified directory, and it will be used as the persistent cache to
  store the App Configuration service information.

When persistent cache is enabled, the SDK will keep the last known good configuration at the persistent cache. In the
case of App Configuration server being unreachable, the latest configurations at the persistent cache is loaded to the
application to continue working.

Please ensure that the cache file is not lost or deleted in any case. For example, consider the case when a kubernetes pod is restarted and the cache file (appconfiguration.json) was stored in ephemeral volume of the pod. As pod gets restarted, kubernetes destroys the ephermal volume in the pod, as a result the cache file gets deleted. So, make sure that the cache file created by the SDK is always stored in persistent volume by providing the correct absolute path of the persistent directory.

### (Optional)

The SDK is also designed to serve configurations, perform feature flag & property evaluations without being connected to
App Configuration service.

```go
appConfigClient.SetContext(collectionId, environmentId, AppConfiguration.ContextOptions{
    BootstrapFile: "saflights/flights.json",
    LiveConfigUpdateEnabled: false,
})
```

* BootstrapFile: Absolute path of the JSON file, which contains configuration details. Make sure to provide a proper
  JSON file. You can generate this file using `ibmcloud ac config` command of the IBM Cloud App Configuration CLI.
* LiveConfigUpdateEnabled: Live configuration update from the server. Set this value to `false` if the new configuration
  values shouldn't be fetched from the server. By default, this value is enabled.

## Get single feature

```go
feature, err := appConfigClient.GetFeature("online-check-in")
if err == nil {
    fmt.Println("Feature Name", feature.GetFeatureName())
    fmt.Println("Feature Id", feature.GetFeatureID())
    fmt.Println("Feature Type", feature.GetFeatureDataType())

    if (feature.IsEnabled()) {
        // feature flag is enabled
    } else {
        // feature flag is disabled
    }
}
```

## Get all features

```go
features, err := appConfigClient.GetFeatures()
if err == nil {
    feature := features["online-check-in"]
    
    fmt.Println("Feature Name", feature.GetFeatureName())
    fmt.Println("Feature Id", feature.GetFeatureID())
    fmt.Println("Feature Type", feature.GetFeatureDataType())
    fmt.Println("Is feature enabled?", feature.IsEnabled())
}
```

## Evaluate a feature

Use the `feature.GetCurrentValue(entityId, entityAttributes)` method to evaluate the value of the feature flag.
GetCurrentValue returns one of the Enabled/Disabled/Overridden value based on the evaluation.

```go
entityId := "john_doe"
entityAttributes := make(map[string]interface{})
entityAttributes["city"] = "Bangalore"
entityAttributes["country"] = "India"

featureVal := feature.GetCurrentValue(entityId, entityAttributes)
```

* entityId: entityId is a string identifier related to the Entity against which the feature will be evaluated. For
  example, an entity might be an instance of an app that runs on a mobile device, a microservice that runs on the cloud,
  or a component of infrastructure that runs that microservice. For any entity to interact with App Configuration, it
  must provide a unique entity ID.

* entityAttributes: entityAttributes is a map of type `map[string]interface{}` consisting of the attribute name and
  their values that defines the specified entity. This is an optional parameter if the feature flag is not configured
  with any targeting definition. If the targeting is configured, then entityAttributes should be provided for the rule
  evaluation. An attribute is a parameter that is used to define a segment. The SDK uses the attribute values to
  determine if the specified entity satisfies the targeting rules, and returns the appropriate feature flag value.

## Get single property

```go
property, err := appConfigClient.GetProperty("check-in-charges")
if err == nil {
    fmt.Println("Property Name", property.GetPropertyName())
    fmt.Println("Property Id", property.GetPropertyID())
    fmt.Println("Property Type", property.GetPropertyDataType())
}
```

## Get all properties

```go
properties, err := appConfigClient.GetProperties()
if err == nil {
    property := properties["check-in-charges"]
    
    fmt.Println("Property Name", property.GetPropertyName())
    fmt.Println("Property Id", property.GetPropertyID())
    fmt.Println("Property Type", property.GetPropertyDataType())
}
```

## Evaluate a property

Use the `property.GetCurrentValue(entityId, entityAttributes)` method to evaluate the value of the property.
GetCurrentValue returns the default property value or its overridden value based on the evaluation.

```go
entityId := "john_doe"
entityAttributes := make(map[string]interface{})
entityAttributes["city"] = "Bangalore"
entityAttributes["country"] = "India"

propertyVal := property.GetCurrentValue(entityId, entityAttributes)
```

* entityId: entityId is a string identifier related to the Entity against which the property will be evaluated. For
  example, an entity might be an instance of an app that runs on a mobile device, a microservice that runs on the cloud,
  or a component of infrastructure that runs that microservice. For any entity to interact with App Configuration, it
  must provide a unique entity ID.
* entityAttributes: entityAttributes is a map of type `map[string]interface{}` consisting of the attribute name and
  their values that defines the specified entity. This is an optional parameter if the property is not configured with
  any targeting definition. If the targeting is configured, then entityAttributes should be provided for the rule
  evaluation. An attribute is a parameter that is used to define a segment. The SDK uses the attribute values to
  determine if the specified entity satisfies the targeting rules, and returns the appropriate property value.

## Get secret property

```go
secretPropertyObject, err := appConfiguration.GetSecret(propertyID, secretsManagerObject)
```

* propertyID: propertyID is the unique string identifier, using this we will be able to fetch the property which will provide the necessary data to fetch the secret.
* secretsManagerObject: secretsManagerObject is an Secret Manager variable or object which will be used for getting the secrets during the secret property evaluation. How to create a secret manager object refer the secret manager docs:https://cloud.ibm.com/apidocs/secrets-manager?code=go

## Evaluate a secret property

Use the `secretPropertyObject.GetCurrentValue(entityId, entityAttributes)` method to evaluate the value of the secret property.
GetCurrentValue returns the secret value based on the evaluation.

```go
entityId := "john_doe"
entityAttributes := make(map[string]interface{})
entityAttributes["city"] = "Bangalore"
entityAttributes["country"] = "India"

getSecretRes, detailedResponse, err := secretPropertyObject.GetCurrentValue(entityId, entityAttributes)
```

* entityId: entityId is a string identifier related to the Entity against which the property will be evaluated. For
  example, an entity might be an instance of an app that runs on a mobile device, a microservice that runs on the cloud,
  or a component of infrastructure that runs that microservice. For any entity to interact with App Configuration, it
  must provide a unique entity ID.
* entityAttributes: entityAttributes is a map of type `map[string]interface{}` consisting of the attribute name and
  their values that defines the specified entity. This is an optional parameter if the property is not configured with
  any targeting definition. If the targeting is configured, then entityAttributes should be provided for the rule
  evaluation. An attribute is a parameter that is used to define a segment. The SDK uses the attribute values to
  determine if the specified entity satisfies the targeting rules, and returns the appropriate value.

## How to access the payload secret data from the response
```go
//make sure this import statement is added
import (sm "github.com/IBM/secrets-manager-go-sdk/secretsmanagerv1")

secret := getSecretRes.Resources[0].(*sm.SecretResource)
secretData := secret.SecretData.(map[string]interface{})
payload := secretData["payload"]
```

The GetCurrentValue will be sending the 3 objects as part of response. 

* getSecretRes:  this will give the meta data and payload.
* detailedResponse: this will give entire data which includes the http response header data, meta data and payload.
* err: this will give the error response if the request is invalid or failed for some reason.
> Note: `secretData["payload"] will return interface{}` so based on the data we need to do the type casting.

## Fetching the appConfigClient across other modules

Once the SDK is initialized, the appConfigClient can be obtained across other modules as shown below:

```go
// **other modules**

import (
        AppConfiguration "github.com/IBM/appconfiguration-go-sdk/lib"
)

appConfigClient := AppConfiguration.GetInstance()
feature, err := appConfigClient.GetFeature("online-check-in")
if (err == nil) {
    enabled := feature.IsEnabled()
    featureValue := feature.GetCurrentValue(entityId, entityAttributes)
}
```

## Supported Data types

App Configuration service allows to configure the feature flag and properties in the following data types : Boolean,
Numeric, SecretRef, String. The String data type can be of the format of a text string , JSON or YAML. The SDK processes each
format accordingly as shown in the below table.
<details><summary>View Table</summary>

| **Feature value or Property value**                                                                                      | **DataType** | **DataFormat** | **Type of data returned <br> by `GetCurrentValue()`** | **Example output**                                                   |
| ------------------------------------------------------------------------------------------------------------------ | ------------ | -------------- | ----------------------------------------------------- | -------------------------------------------------------------------- |
| `true`                                                                                                             | BOOLEAN      | not applicable | `bool`                                                | `true`                                                               |
| `25`                                                                                                               | NUMERIC      | not applicable | `float64`                                             | `25`                                                                 |
| <pre>{<br> "secret_type": "kv",<br>"id": "secret_id_data_here",<br> "sm_instance_crn": "crn_data_added-here"<br>}</pre>                                                                                                               | SECRETREF(`this type is applicable only for Property`)     | not applicable | `map[string]interface{}`                                             | <pre>{<br>"metadata": {<br>"collection_type":"application/vnd.ibm.secrets-manager.secret+json",<br>"collection_total": 1<br>},<br>"resources": [{"created_by": "iam-ServiceId-e4a2f0a4-3c76-4bef-b1f2-fbeae11c0f21",<br>"creation_date": "2020-10-05T21:33:11Z",<br>"crn": "crn:v1:bluemix:public:secrets-manager:us-south:a/a5ebf2570dcaedf18d7ed78e216c263a:f1bc94a6-64aa-4c55-b00f-f6cd70e4b2ce:secret:cb7a2502-8ede-47d6-b5b6-1b7af6b6f563",<br>"description": "Extended description for this secret.",<br>"expiration_date": "2021-01-01T00:00:00Z",<br>"id": "cb7a2502-8ede-47d6-b5b6-1b7af6b6f563",<br>"labels": ["dev","us-south"],<br>"last_update_date": "2020-10-05T21:33:11Z",<br>"name": "example-arbitrary-secret",<br>"secret_data": {"payload": "secret-data"},<br>"secret_type": "arbitrary",<br>"state": 1,<br>"state_description": "Active",<br>"versions_total": 1,<br>"versions": [{"created_by": "iam-ServiceId-222b47ab-b08e-4619-b68f-8014a2c3acb8","creation_date": "2020-11-23T20:15:01Z","id": "50277266-d706-4b3e-badb-f07257f8f581","payload_available": true,"downloaded": true}],"locks_total": 2}]<br>}</pre>  `Note: Along with the above data we will also provide the detailedResponse and error data. For more info on the response data refer #how-to-access-the-payload-secret-data-from-the-response`                                                                |
| "a string text"                                                                                                    | STRING       | TEXT           | `string`                                              | `a string text`                                                      |
| <pre>{<br>  "firefox": {<br>    "name": "Firefox",<br>    "pref_url": "about:config"<br>  }<br>}</pre> | STRING       | JSON           | `map[string]interface{}`                              | `map[browsers:map[firefox:map[name:Firefox pref_url:about:config]]]` |
| <pre>men:<br>  - John Smith<br>  - Bill Jones<br>women:<br>  - Mary Smith<br>  - Susan Williams</pre>  | STRING       | YAML           | `map[string]interface{}`                              | `map[men:[John Smith Bill Jones] women:[Mary Smith Susan Williams]]` |
</details>

<details><summary>Feature flag</summary>

  ```go
feature, err := appConfigClient.GetFeature("json-feature")
if err == nil {
    feature.GetFeatureDataType() // STRING
    feature.GetFeatureDataFormat() // JSON
    
    // Example (traversing the returned map)
    result := feature.GetCurrentValue(entityID, entityAttributes) // JSON value is returned as a Map
    result.(map[string]interface{})["key"] // returns the value of the key
}

feature, err := appConfigClient.GetFeature("yaml-feature")
if err == nil {
    feature.GetFeatureDataType() // STRING
    feature.GetFeatureDataFormat() // YAML
    
    // Example (traversing the returned map)
    result := feature.GetCurrentValue(entityID, entityAttributes) // YAML value is returned as a Map
    result.(map[string]interface{})["key"] // returns the value of the key
}
  ```

</details>
<details><summary>Property</summary>

  ```go
property, err := appConfigClient.GetProperty("json-property")
if err == nil {
    property.GetPropertyDataType() // STRING
    property.GetPropertyDataFormat() // JSON

    // Example (traversing the returned map)
    result := property.GetCurrentValue(entityID, entityAttributes) // JSON value is returned as a Map
    result.(map[string]interface{})["key"] // returns the value of the key
}

property, err := appConfigClient.GetProperty("yaml-property")
if err == nil {
    property.GetPropertyDataType() // STRING
    property.GetPropertyDataFormat() // YAML

    // Example (traversing the returned map)
    result := property.GetCurrentValue(entityID, entityAttributes) // YAML value is returned as a Map
    result.(map[string]interface{})["key"] // returns the value of the key
}
  ```

</details>

## Set listener for feature or property data changes

The SDK provides mechanism to notify you in real-time when feature flag's or property's configuration changes.
You can subscribe to configuration changes using the same appConfigClient.

```go
appConfigClient.RegisterConfigurationUpdateListener(func () {
      // **add your code**
      // To find the effect of any configuration changes, you can call the feature or property related methods

      // feature, err := appConfigClient.GetFeature("json-feature")
      // newValue := feature.GetCurrentValue(entityID, entityAttributes)
})
```

## Fetch latest data

```go
appConfigClient.FetchConfigurations()
```

## Enable debugger (Optional)

```go
appConfigClient.EnableDebug(true)
```

## Examples

Try [this](https://github.com/IBM/appconfiguration-go-sdk/tree/master/examples) sample application in the examples
folder to learn more about feature and property evaluation.

## License

This project is released under the Apache 2.0 license. The license's full text can be found
in [LICENSE](https://github.com/IBM/appconfiguration-go-sdk/blob/master/LICENSE)
