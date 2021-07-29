package adapter

import (
	"fmt"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/BrobridgeOrg/gravity-cli/pkg/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	hostFlag      string
	domainFlag    string
	accessKeyFlag string
	appIDFlag     string
)

type SetConfigCmd struct {
	cli          cli.Cli
	setConfigCmd *cobra.Command
}

func NewSetConfigCmd(c cli.Cli) *SetConfigCmd {
	setConfig := &SetConfigCmd{
		cli: c,
	}

	return setConfig
}

func (setConfig *SetConfigCmd) Init() *cobra.Command {
	setConfig.setConfigCmd = &cobra.Command{
		Use:   "setConfig",
		Short: "Set Gravity CLI configuration",
		Long:  `Set Gravity CLI configuration and write config file to Gravity CLI config file`,
		Run: func(cmd *cobra.Command, args []string) {

			configFileFlag := setConfig.cli.GetConfigFile()

			// Check configPath is exist
			configPath := path.Dir(configFileFlag)
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				// Create configPath
				err := os.MkdirAll(configPath, 0700)
				if err != nil {
					log.Fatal(err)
				}
			}

			viper.SetConfigType("toml")
			viper.SetConfigFile(configFileFlag)

			// Set viper value
			tableName := "gravity"
			appIDKey := fmt.Sprintf("%s.%s", tableName, "appID")
			hostKey := fmt.Sprintf("%s.%s", tableName, "host")
			domainKey := fmt.Sprintf("%s.%s", tableName, "domain")
			accKey := fmt.Sprintf("%s.%s", tableName, "accessKey")

			viper.Set(appIDKey, appIDFlag)
			viper.Set(hostKey, hostFlag)
			viper.Set(domainKey, domainFlag)
			viper.Set(accKey, accessKeyFlag)

			// Write to config file
			err := viper.WriteConfigAs(configFileFlag)
			if err != nil {
				log.Fatal(err)
			}
			return
		},
	}

	setConfig.setConfigCmd.Flags().StringVarP(&hostFlag, "host", "H", "", "Specify Gravity host:port (example \"127.0.0.1:4222\")")
	setConfig.setConfigCmd.Flags().StringVarP(&domainFlag, "domain", "d", "gravity", "Specify Gravity Domain")
	setConfig.setConfigCmd.Flags().StringVarP(&accessKeyFlag, "accessKey", "k", "", "Specify Gravity AccessKey.")
	setConfig.setConfigCmd.Flags().StringVarP(&appIDFlag, "appID", "i", "", "Specify Application ID")
	setConfig.setConfigCmd.MarkFlagRequired("host")
	//setConfig.setConfigCmd.MarkFlagRequired("accessKey")
	//setConfig.setConfigCmd.MarkFlagRequired("appID")

	return setConfig.setConfigCmd
}
