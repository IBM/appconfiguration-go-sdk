# Example Application

> **DISCLAIMER**: This is a guideline example application and is used for demonstrative and illustrative purposes only. This is not a production ready code.

## Step 1: Create an instance of App Configuration service

- Log in to your IBM Cloud account.
- In the [IBM Cloud catalog](https://cloud.ibm.com/catalog#services), search **App Configuration** and
  select [App Configuration](https://cloud.ibm.com/catalog/services/app-configuration). The service configuration screen
  opens.
- **Select a region** - Currently, Dallas (us-south), London (eu-gb), Sydney (au-syd), Washington DC (us-east), Frankfurt (eu-de), Toronto (ca-tor), Tokyo (jp-tok) and Osaka (jp-osa) regions are supported.
- Select a pricing plan, resource group and configure your resource with a service name, or use the preset name.
- Click **Create**. A new service instance is created and the App Configuration console displayed.

## Step 2: Generate Service Credentials

- Go to dashboard page of App Configuration service instance in the IBM Cloud UI.
- Navigate to Service Credentials section and generate a new set of credentials. Provide these generated `region`
  , `guid` and `apikey` values in [environment file](.env).

## Step 3: Create a collection, segment, feature flag & add targeting to feature flag

- On to dashboard page of App Configuration service instance created, navigate to Collections section and create a
  collection by clicking on create button.

```
//Collection details

Name: Car Rentals
Collection ID: car-rentals
```

- On the same instance, navigate to Segments section and create a segment by clicking on create button.

```
//Segment details

Name: Users from Bangalore urban area
Add rule one:
            Attribute name: city
            Operator: is
            Values: Bangalore
Add rule two:
            Attribute name: radius
            Operator: less than and equals
            Values: 60
            
```

- Again on the same instance, navigate to Feature flags section and create a feature flag by clicking on create button.

```
//Feature flag details

Name: Weekend discount
Feature flag ID: weekend-discount
Flag type: Numeric
Enabled value: 10
Disabled value: 0
Add this feature flag to the above collection created "Car Rentals"
```

- Navigate to Properties section and create a property by clicking on create button.

```
//Property details

Name: User location
Property ID: user-location
Property type & format: String & Text
Default value: other
Add this property to the above collection created "Car Rentals"
```

- Click `Add targeting` on the feature flag created & target this feature flag to the segment.

```
Select the segment "Users from Bangalore urban area" from the dropdown.
Click on radio button "Override" in the Enabled value section, and give some value(say 25) in the below numeric input.
Next, click on "Save rule" & and then click on "Add targeting". This applies the targeting to the feature flag.
```

- Click `Add targeting` on the property created & target this property to the segment.

```
Select the segment "Users from Bangalore urban area" from the dropdown.
Click on radio button "Override" in the Property value section, and give the value(say "Bangalore") in the below text input.
Next, click on "Save rule" & and then click on "Add targeting". This applies the targeting to the property.
```

## Step 4: Run the app

```bash
$ cd ./examples
$ go mod tidy
$ go run main.go
```

Expected output

```
FEATURE FLAG OPERATIONS

Feature Name: Weekend discount
Feature Id: weekend-discount
Feature Data type: NUMERIC
Is Feature enabled? *true or false* based on the configuration
Feature evaluated value is: 25


PROPERTY OPERATIONS

Property Name: User location
Property Id: user-location
Property Data type: STRING
Property evaluated value is: Bangalore

```
