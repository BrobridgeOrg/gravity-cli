package cmd

import (
	"github.com/BrobridgeOrg/gravity-cli/pkg/configs"
	"github.com/BrobridgeOrg/gravity-cli/pkg/connector"
	"github.com/BrobridgeOrg/gravity-cli/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var config *configs.Config
var host string
var domain string

var RootCmd = &cobra.Command{
	Use:   "gravity-cli",
	Short: "Gravity utility",
	Long:  `gravity-cli is a command line utility of Gravity.`,
	/*
		RunE: func(cmd *cobra.Command, args []string) error {

			if err := run(); err != nil {
				return err
			}
			return nil
		},
	*/
}

func init() {
	config = configs.GetConfig()

	RootCmd.PersistentFlags().StringVarP(&host, "host", "s", "0.0.0.0:32803", "Specify server address")
	RootCmd.PersistentFlags().StringVarP(&domain, "domain", "d", "default", "Specify data domain")
}

func run() error {

	config.SetHost(host)

	fx.New(
		fx.Supply(config),
		fx.Provide(
			logger.GetLogger,
			connector.New,
		),
		fx.NopLogger,
	).Run()

	return nil
}
