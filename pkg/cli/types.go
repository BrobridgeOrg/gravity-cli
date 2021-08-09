package cli

import (
	adapter_manager "github.com/BrobridgeOrg/gravity-sdk/adapter_manager"
	auth "github.com/BrobridgeOrg/gravity-sdk/authenticator"
	subscriber_manager "github.com/BrobridgeOrg/gravity-sdk/subscriber_manager"
)

type Cli interface {
	GetAuthClient() (*auth.Authenticator, error)
	GetSubscriberManagerClient() (*subscriber_manager.SubscriberManager, error)
	GetAdapterManagerClient() (*adapter_manager.AdapterManager, error)
	GetConfigFile() string
	SetConfigFile(string)
}
