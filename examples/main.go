// main.go
package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	AppConfiguration "github.com/IBM/appconfiguration-go-sdk/lib"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Welcome to Sample App HomePage!")
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	region := os.Getenv("REGION")
	guid := os.Getenv("GUID")
	apikey := os.Getenv("APIKEY")
	collectionId := os.Getenv("COLLECTION_ID")
	environmentId := os.Getenv("ENVIRONMENT_ID")

	appConfiguration := AppConfiguration.GetInstance()
	appConfiguration.Init(region, guid, apikey)
	appConfiguration.SetContext(collectionId, environmentId)
	entityID := "user123"
	entityAttributes := make(map[string]interface{})
	entityAttributes["city"] = "Bangalore"
	entityAttributes["radius"] = 60

	fmt.Println("\n\nFEATURE FLAG OPERATIONS\n")
	feature, err := appConfiguration.GetFeature(os.Getenv("FEATURE_ID"))
	if err == nil {
		fmt.Println("Feature Name:", feature.GetFeatureName())
		fmt.Println("Feature Id:", feature.GetFeatureID())
		fmt.Println("Feature Data type:", feature.GetFeatureDataType())
		fmt.Println("Is Feature enabled?", feature.IsEnabled())
		fmt.Println("Feature evaluated value is:", feature.GetCurrentValue(entityID, entityAttributes))
	}

	fmt.Println("\n\nPROPERTY OPERATIONS\n")
	property, err := appConfiguration.GetProperty(os.Getenv("PROPERTY_ID"))
	if err == nil {
		fmt.Println("Property Name:", property.GetPropertyName())
		fmt.Println("Property Id:", property.GetPropertyID())
		fmt.Println("Property Data type:", property.GetPropertyDataType())
		fmt.Println("Property evaluated value is:", property.GetCurrentValue(entityID, entityAttributes))
	}
	//whenever the configurations get changed/updated on the app configuration service instance the function inside this listener is triggered.
	//So, to keep track of live changes to configurations use this listener.
	appConfiguration.RegisterConfigurationUpdateListener(func() {
		fmt.Println("configurations updated")
		// To find the effect of any configuration changes, you should call the feature or property related methods again

		// feature, err := appConfigClient.GetFeature("feature-id")
		// newValue := feature.GetCurrentValue(entityID, entityAttributes)
	})

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
