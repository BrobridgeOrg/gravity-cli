module github.com/BrobridgeOrg/gravity-cli

go 1.15

require (
	github.com/BrobridgeOrg/compton v0.0.0-20220617174904-7083c8a5102d
	github.com/BrobridgeOrg/gravity-sdk/v2 v2.0.2
	github.com/docker/go-units v0.5.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/nats-io/nats.go v1.16.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/spf13/afero v1.8.1 // indirect
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.1
	github.com/stretchr/testify v1.8.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/fx v1.17.0
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.0.0-20220307211146-efcb8507fb70 // indirect
	golang.org/x/sys v0.0.0-20220307203707-22a9840ba4d7 // indirect
	google.golang.org/protobuf v1.27.1
	gopkg.in/ini.v1 v1.66.4 // indirect
)

// replace github.com/BrobridgeOrg/compton => ../../compton

//replace github.com/BrobridgeOrg/gravity-sdk/v2 => ../gravity-sdk
