module github.com/IBM/appconfiguration-go-sdk

go 1.16

require (
	github.com/IBM/go-sdk-core/v5 v5.14.1
	github.com/IBM/secrets-manager-go-sdk/v2 v2.0.1
	github.com/go-openapi/strfmt v0.21.7 // indirect
	github.com/gorilla/websocket v1.5.0
	github.com/robfig/cron v1.2.0
	github.com/sirupsen/logrus v1.9.0
	github.com/spaolacci/murmur3 v1.1.0
	github.com/stretchr/testify v1.8.2
	go.mongodb.org/mongo-driver v1.11.4 // indirect
	golang.org/x/net v0.17.0 // indirect
	gopkg.in/yaml.v3 v3.0.1
)

//Retract v1.x.x versions
retract [v1.0.0, v1.2.1]
