module github.com/IBM/appconfiguration-go-sdk

go 1.25.0

require (
	github.com/IBM/go-sdk-core/v5 v5.23.3
	github.com/IBM/secrets-manager-go-sdk/v2 v2.0.22
	github.com/emirpasic/gods v1.18.1
	github.com/gorilla/websocket v1.5.3
	github.com/robfig/cron v1.2.0
	github.com/sirupsen/logrus v1.10.2
	github.com/spaolacci/murmur3 v1.1.0
	github.com/stretchr/testify v1.12.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/gabriel-vasile/mimetype v1.4.15 // indirect
	github.com/go-openapi/errors v0.22.8 // indirect
	github.com/go-openapi/strfmt v0.27.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.30.3 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.8 // indirect
	github.com/leodido/go-urn v1.5.0 // indirect
	github.com/oklog/ulid/v2 v2.1.2 // indirect
	go.yaml.in/yaml/v2 v2.4.4 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/crypto v0.55.0 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)

//Retract v1.x.x versions
retract [v1.0.0, v1.2.1]
