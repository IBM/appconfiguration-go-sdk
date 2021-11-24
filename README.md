# IBM Cloud App Configuration Go server SDK 0.2.2

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
appConfiguration := AppConfiguration.GetInstance()
appConfiguration.Init("region", "guid", "apikey")

collectionId := "airlines-webapp"
environmentId := "dev"
appConfiguration.SetContext(collectionId, environmentId)
```

- region : Region name where the App Configuration service instance is created. Use
    - `AppConfiguration.REGION_US_SOUTH` for Dallas
    - `AppConfiguration.REGION_EU_GB` for London
    - `AppConfiguration.REGION_AU_SYD` for Sydney
- guid : Instance Id of the App Configuration service. Obtain it from the service credentials section of the App
  Configuration dashboard.
- apikey : ApiKey of the App Configuration service. Obtain it from the service credentials section of the App
  Configuration dashboard.
* collectionId: Id of the collection created in App Configuration service instance under the **Collections** section.
* environmentId: Id of the environment created in App Configuration service instance under the **Environments** section.

### (Optional)

In order for your application and SDK to continue its operations even during the unlikely scenario of App Configuration
service across your application restarts, you can configure the SDK to work using a persistent cache. The SDK uses the
persistent cache to store the App Configuration data that will be available across your application restarts.

```go
// 1. default (without persistent cache)
appConfiguration.SetContext(collectionId, environmentId)

// 2. with persistent cache
appConfiguration.SetContext(collectionId, environmentId, AppConfiguration.ContextOptions{
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
appConfiguration.SetContext(collectionId, environmentId, AppConfiguration.ContextOptions{
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
feature, err := appConfiguration.GetFeature("online-check-in")
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
features, err := appConfiguration.GetFeatures()
if err == nil {
    feature := features["online-check-in"]
    
    fmt.Println("Feature Name", feature.GetFeatureName())
    fmt.Println("Feature Id", feature.GetFeatureID())
    fmt.Println("Feature Type", feature.GetFeatureDataType())
    fmt.Println("Feature is enabled", feature.IsEnabled())
}
```

## Evaluate a feature

You can use the ` feature.GetCurrentValue(entityId, entityAttributes)` method to evaluate the value of the feature
flag. You should pass an unique entityId as the parameter to perform the feature flag evaluation. If the feature flag
is configured with segments in the App Configuration service, you can set the attributes values as a map.

```go
entityId := "john_doe"
entityAttributes := make(map[string]interface{})
entityAttributes["city"] = "Bangalore"
entityAttributes["country"] = "India"

featureVal := feature.GetCurrentValue(entityId, entityAttributes)
```

## Get single property

```go
property, err := appConfiguration.GetProperty("check-in-charges")
if err == nil {
    fmt.Println("Property Name", property.GetPropertyName())
    fmt.Println("Property Id", property.GetPropertyID())
    fmt.Println("Property Type", property.GetPropertyDataType())
}
```

## Get all properties

```go
properties, err := appConfiguration.GetProperties()
if err == nil {
    property := properties["check-in-charges"]
    
    fmt.Println("Property Name", property.GetPropertyName())
    fmt.Println("Property Id", property.GetPropertyID())
    fmt.Println("Property Type", property.GetPropertyDataType())
}
```

## Evaluate a property

You can use the ` property.GetCurrentValue(entityId, entityAttributes)` method to evaluate the value of the
property. You should pass an unique entityId as the parameter to perform the property evaluation. If the property is
configured with segments in the App Configuration service, you can set the attributes values as a map.

```go
entityId := "john_doe"
entityAttributes := make(map[string]interface{})
entityAttributes["city"] = "Bangalore"
entityAttributes["country"] = "India"

propertyVal := property.GetCurrentValue(entityId, entityAttributes)
```

## Supported Data types

App Configuration service allows to configure the feature flag and properties in the following data types : Boolean,
Numeric, String. The String data type can be of the format of a text string , JSON or YAML. The SDK processes each
format accordingly as shown in the below table.
<details><summary>View Table</summary>

| **Feature or Property value**                                                                                      | **DataType** | **DataFormat** | **Type of data returned <br> by `GetCurrentValue()`** | **Example output**                                                   |
| ------------------------------------------------------------------------------------------------------------------ | ------------ | -------------- | ----------------------------------------------------- | -------------------------------------------------------------------- |
| `true`                                                                                                             | BOOLEAN      | not applicable | `bool`                                                | `true`                                                               |
| `25`                                                                                                               | NUMERIC      | not applicable | `float64`                                             | `25`                                                                 |
| "a string text"                                                                                                    | STRING       | TEXT           | `string`                                              | `a string text`                                                      |
| <pre>{<br>  "firefox": {<br>    "name": "Firefox",<br>    "pref_url": "about:config"<br>  }<br>}</pre> | STRING       | JSON           | `map[string]interface{}`                              | `map[browsers:map[firefox:map[name:Firefox pref_url:about:config]]]` |
| <pre>men:<br>  - John Smith<br>  - Bill Jones<br>women:<br>  - Mary Smith<br>  - Susan Williams</pre>  | STRING       | YAML           | `map[string]interface{}`                              | `map[men:[John Smith Bill Jones] women:[Mary Smith Susan Williams]]` |
</details>

<details><summary>Feature flag</summary>

  ```go
feature, err := appConfiguration.GetFeature("json-feature")
if err == nil {
    feature.GetFeatureDataType() // STRING
    feature.GetFeatureDataFormat() // JSON
    
    // Example (traversing the returned map)
    result := feature.GetCurrentValue(entityID, entityAttributes) // JSON value is returned as a Map
    result.(map[string]interface{})["key"] // returns the value of the key
}

feature, err := appConfiguration.GetFeature("yaml-feature")
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
property, err := appConfiguration.GetProperty("json-property")
if err == nil {
    property.GetPropertyDataType() // STRING
    property.GetPropertyDataFormat() // JSON

    // Example (traversing the returned map)
    result := property.GetCurrentValue(entityID, entityAttributes) // JSON value is returned as a Map
    result.(map[string]interface{})["key"] // returns the value of the key
}

property, err := appConfiguration.GetProperty("yaml-property")
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

To listen to the configurations changes in your App Configuration service instance, implement the `RegisterConfigurationUpdateListener` event listener as mentioned below 

```go
appConfiguration.RegisterConfigurationUpdateListener(func () {
    fmt.Println("Get your feature/property value now ")
})
```

## Fetch latest data

```go
appConfiguration.FetchConfigurations()
```

## Enable debugger (Optional)

```go
appConfiguration.EnableDebug(true)
```

## Examples

Try [this](https://github.com/IBM/appconfiguration-go-sdk/tree/master/examples) sample application in the examples
folder to learn more about feature and property evaluation.

## License

This project is released under the Apache 2.0 license. The license's full text can be found
in [LICENSE](https://github.com/IBM/appconfiguration-go-sdk/blob/master/LICENSE)
