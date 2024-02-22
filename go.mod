module github.com/IBM/appconfiguration-go-sdk

go 1.16

require (
	github.com/IBM/go-sdk-core/v5 v5.15.0
	github.com/IBM/secrets-manager-go-sdk/v2 v2.0.2
	github.com/go-openapi/strfmt v0.22.0 // indirect
	github.com/gorilla/websocket v1.5.1
	github.com/robfig/cron v1.2.0
	github.com/sirupsen/logrus v1.9.3
	github.com/spaolacci/murmur3 v1.1.0
	github.com/stretchr/testify v1.8.4
	golang.org/x/crypto v0.19.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	gopkg.in/yaml.v3 v3.0.1
)

//Retract v1.x.x versions
retract [v1.0.0, v1.2.1]
