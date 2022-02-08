package product

import (
	"github.com/BrobridgeOrg/gravity-cli/pkg/configs"
	"github.com/BrobridgeOrg/gravity-cli/pkg/connector"
	"github.com/BrobridgeOrg/gravity-sdk/product"
	"go.uber.org/zap"
)

var logger *zap.Logger

type Product struct {
	config        *configs.Config
	connector     *connector.Connector
	productClient *product.ProductClient
}

func New(config *configs.Config, l *zap.Logger, c *connector.Connector) *Product {

	logger = l

	p := &Product{
		config:    config,
		connector: c,
	}

	pcOpts := product.NewOptions()
	pcOpts.Domain = c.GetDomain()
	productClient := product.NewProductClient(
		c.GetClient(),
		pcOpts,
	)

	if productClient == nil {
		return nil
	}

	p.productClient = productClient

	return p
}

func (p *Product) GetClient() *product.ProductClient {
	return p.productClient
}
