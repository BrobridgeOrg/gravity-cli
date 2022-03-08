module github.com/BrobridgeOrg/gravity-cli

go 1.15

require (
	github.com/BrobridgeOrg/gravity-api v0.2.25
	github.com/BrobridgeOrg/gravity-dispatcher v0.0.0-20220206181110-2ea65aa048be
	github.com/BrobridgeOrg/gravity-sdk v0.0.50
	github.com/bketelsen/crypt v0.0.4 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.1.2
	github.com/manifoldco/promptui v0.8.0
	github.com/nats-io/nats.go v1.13.1-0.20220121202836-972a071d373d
	github.com/olekukonko/tablewriter v0.0.5
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.1
	go.uber.org/fx v1.16.0
	go.uber.org/zap v1.17.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/BrobridgeOrg/gravity-sdk => ../gravity-sdk
