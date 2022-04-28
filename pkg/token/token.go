package token

import (
	"context"
	"errors"

	"github.com/BrobridgeOrg/gravity-cli/pkg/configs"
	"github.com/BrobridgeOrg/gravity-cli/pkg/connector"
	"github.com/BrobridgeOrg/gravity-sdk/token"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var logger *zap.Logger

type Token struct {
	config      *configs.Config
	connector   *connector.Connector
	tokenClient *token.TokenClient
}

func New(lifecycle fx.Lifecycle, config *configs.Config, l *zap.Logger, c *connector.Connector) *Token {

	logger = l.Named("Token")

	p := &Token{
		config:    config,
		connector: c,
	}

	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				pcOpts := token.NewOptions()
				pcOpts.Domain = c.GetDomain()
				tokenClient := token.NewTokenClient(
					c.GetClient(),
					pcOpts,
				)

				if tokenClient == nil {
					return errors.New("Failed to create token client")
				}

				p.tokenClient = tokenClient

				return nil
			},
			OnStop: func(ctx context.Context) error {
				return nil
			},
		},
	)

	return p
}

func (p *Token) GetClient() *token.TokenClient {
	return p.tokenClient
}
