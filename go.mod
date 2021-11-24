module github.com/IBM/appconfiguration-go-sdk

go 1.16

replace github.com/gobuffalo/packr/v2 => github.com/gobuffalo/packr/v2 v2.3.2

require (
	github.com/IBM/go-sdk-core/v5 v5.7.0
	github.com/google/uuid v1.2.0 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/robfig/cron v1.2.0
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	golang.org/x/text v0.3.6 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200605160147-a5ece683394c
)

//Retract v1.x.x versions
retract [v1.0.0, v1.2.1]
