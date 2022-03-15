package cmd

import (
	"context"

	"github.com/BrobridgeOrg/gravity-cli/pkg/configs"
	"github.com/BrobridgeOrg/gravity-cli/pkg/connector"
	"github.com/BrobridgeOrg/gravity-cli/pkg/logger"
	"github.com/BrobridgeOrg/gravity-cli/pkg/product"
	adapter_sdk "github.com/BrobridgeOrg/gravity-sdk/adapter"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type PubCommandContext struct {
	Config    *configs.Config
	Logger    *zap.Logger
	Connector *connector.Connector
	Product   *product.Product
	Cmd       *cobra.Command
	Args      []string
}

type pubCmdFunc func(*PubCommandContext) error

var pubEvent string
var pubPayload string

func init() {

	RootCmd.AddCommand(pubCmd)
}

var pubCmd = &cobra.Command{
	Use:   "pub [event] [payload]",
	Short: "Publish domain event",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runPubCmd(runPublishCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runPubCmd(fn pubCmdFunc, cmd *cobra.Command, args []string) error {

	var cctx *PubCommandContext

	config.SetHost(host)

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
			cmd *cobra.Command,
			args []string,
		) *PubCommandContext {
			return &PubCommandContext{
				Config:    config,
				Logger:    l,
				Connector: c,
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

func runPublishCmd(cctx *PubCommandContext) error {

	pubEvent = cctx.Args[0]
	pubPayload = cctx.Args[1]

	// Initializing adapter connector
	opts := adapter_sdk.NewOptions()
	opts.Domain = cctx.Connector.GetDomain()

	ac := adapter_sdk.NewAdapterConnectorWithClient(cctx.Connector.GetClient(), opts)
	_, err := ac.Publish(pubEvent, []byte(pubPayload), nil)
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return err
	}

	return nil
}
