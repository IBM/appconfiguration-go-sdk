module github.com/IBM/appconfiguration-go-sdk

go 1.16

require (
	github.com/IBM/go-sdk-core/v5 v5.10.1
	github.com/IBM/secrets-manager-go-sdk v1.0.45
	github.com/google/uuid v1.2.0 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/robfig/cron v1.2.0
	github.com/sirupsen/logrus v1.7.0
	github.com/spaolacci/murmur3 v1.1.0
	github.com/stretchr/testify v1.7.1
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	gopkg.in/yaml.v3 v3.0.0-20200605160147-a5ece683394c
)

//Retract v1.x.x versions
retract [v1.0.0, v1.2.1]
