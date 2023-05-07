package product

import (
	"context"
	"errors"

	"github.com/BrobridgeOrg/gravity-cli/pkg/configs"
	"github.com/BrobridgeOrg/gravity-cli/pkg/connector"
	"github.com/BrobridgeOrg/gravity-sdk/v2/product"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var logger *zap.Logger

type Product struct {
	config        *configs.Config
	connector     *connector.Connector
	productClient *product.ProductClient
}

func New(lifecycle fx.Lifecycle, config *configs.Config, l *zap.Logger, c *connector.Connector) *Product {

	logger = l.Named("Product")

	p := &Product{
		config:    config,
		connector: c,
	}

	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				pcOpts := product.NewOptions()
				pcOpts.Domain = c.GetDomain()
				productClient := product.NewProductClient(
					c.GetClient(),
					pcOpts,
				)

				if productClient == nil {
					return errors.New("Failed to create product client")
				}

				p.productClient = productClient

				return nil
			},
			OnStop: func(ctx context.Context) error {
				return nil
			},
		},
	)

	return p
}

func (p *Product) GetClient() *product.ProductClient {
	return p.productClient
}
