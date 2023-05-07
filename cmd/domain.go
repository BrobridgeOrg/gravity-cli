package cmd

import (
	"context"
	"fmt"

	"github.com/BrobridgeOrg/gravity-cli/pkg/configs"
	"github.com/BrobridgeOrg/gravity-cli/pkg/connector"
	"github.com/BrobridgeOrg/gravity-cli/pkg/logger"
	"github.com/BrobridgeOrg/gravity-cli/pkg/product"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	domainEventStream  = "GVT_%s"
	domainEventSubject = "$GVT.%s.EVENT.%s"
)

type DomainCommandContext struct {
	Config    *configs.Config
	Logger    *zap.Logger
	Connector *connector.Connector
	Publisher *connector.Connector
	Product   *product.Product
	Cmd       *cobra.Command
	Args      []string
}

type domainCmdFunc func(*DomainCommandContext) error

func init() {

	RootCmd.AddCommand(domainPurgeCmd)
}

var domainPurgeCmd = &cobra.Command{
	Use:   "purge event",
	Short: "Purge domain event",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runDomainCmd(runDomainPurgeCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runDomainCmd(fn domainCmdFunc, cmd *cobra.Command, args []string) error {

	var cctx *DomainCommandContext

	config.SetHost(host)
	config.SetDomain(domain)
	config.SetAccessToken(accessToken)

	app := fx.New(
		fx.Supply(config),
		fx.Provide(
			logger.GetLogger,
			connector.New,
			product.New,
		),
		fx.Supply(cmd),
		fx.Supply(args),
		fx.Provide(func(
			config *configs.Config,
			l *zap.Logger,
			c *connector.Connector,
			publisher *connector.Connector,
			cmd *cobra.Command,
			args []string,
		) *DomainCommandContext {
			return &DomainCommandContext{
				Config:    config,
				Logger:    l,
				Connector: c,
				Publisher: publisher,
				Cmd:       cmd,
				Args:      args,
			}
		}),
		fx.Populate(&cctx),
		fx.NopLogger,
	)

	if err := app.Start(context.Background()); err != nil {
		return err
	}

	return fn(cctx)
}

func runDomainPurgeCmd(cctx *DomainCommandContext) error {

	js, err := cctx.Connector.GetClient().GetJetStream()
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return err
	}

	streamName := fmt.Sprintf(domainEventStream, cctx.Connector.GetDomain())

	err = js.PurgeStream(streamName)
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return err
	}

	return nil
}
