package root

import (
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/BrobridgeOrg/gravity-cli/pkg/cli"
)

var configFileFlag string

type RootCmd struct {
	cli     cli.Cli
	rootCmd *cobra.Command
}

func NewRootCmd(c cli.Cli) *RootCmd {

	rootCmd := &RootCmd{
		cli: c,
	}

	return rootCmd
}

func (root *RootCmd) Init() *cobra.Command {

	root.rootCmd = &cobra.Command{
		Use:   "gravity-cli",
		Short: "Gravity-cli tool",
		Long: `Gravity CLI tool

If $GRACONFIG environment variable is set, then that config file is loaded.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			envConfigFile := os.Getenv("GRACONFIG")
			if len(envConfigFile) != 0 {
				configFileFlag = envConfigFile
			}

			root.cli.SetConfigFile(configFileFlag)
		},
	}

	//Generate config path
	homeDir, _ := os.UserHomeDir()
	configPath := path.Join(homeDir, ".gravity/")

	// Generate Config file
	configFile := path.Join(configPath, "config.toml")

	root.rootCmd.PersistentFlags().StringVar(&configFileFlag, "config", configFile, "Specify Gravity CLI config file")
	return root.rootCmd
}
