package cli

import (
	auth "github.com/BrobridgeOrg/gravity-sdk/authenticator"
)

type Cli interface {
	GetAuthClient() (*auth.Authenticator, error)
	GetConfigFile() string
	SetConfigFile(string)
}
