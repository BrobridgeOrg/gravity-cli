module github.com/BrobridgeOrg/gravity-cli

go 1.15

require (
	github.com/BrobridgeOrg/gravity-api v0.2.25
	github.com/BrobridgeOrg/gravity-dispatcher v0.0.0-20220206181110-2ea65aa048be
	github.com/BrobridgeOrg/gravity-sdk v0.0.50
	github.com/bketelsen/crypt v0.0.4 // indirect
	github.com/google/uuid v1.1.2
	github.com/manifoldco/promptui v0.8.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.10.1
	go.uber.org/fx v1.16.0
	go.uber.org/zap v1.17.0
)

replace github.com/BrobridgeOrg/gravity-sdk => ../gravity-sdk
