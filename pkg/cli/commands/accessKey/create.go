package accessKey

import (
	auth "github.com/BrobridgeOrg/gravity-sdk/authenticator"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func (accessKey *AccessKeyCmd) newCreateAccessKeyCmd() *cobra.Command {
	var appNameFlag string
	var appIDFlag string
	var accessKeyFlag string
	var createAccessKeyCmd = &cobra.Command{
		Use:   "create",
		Short: "Create Gravity Subscriber's Access Key",
		Long:  `Create Gravity Subscriber's Access Key`,
		Run: func(cmd *cobra.Command, args []string) {

			//Get auth client
			authClient, err := accessKey.cli.GetAuthClient()
			if err != nil {
				log.Fatal(err)
			}

			entity := auth.NewEntity()
			entity.AppID = appIDFlag
			entity.AccessKey = accessKeyFlag
			entity.AppName = appNameFlag

			if err := authClient.CreateEntity(entity); err != nil {
				log.Fatal(err)
			}

		},
	}

	createAccessKeyCmd.Flags().StringVarP(&appNameFlag, "name", "n", "", "Specify client's accessKey name")
	createAccessKeyCmd.Flags().StringVarP(&appIDFlag, "appID", "i", "", "Specify client's appID")
	createAccessKeyCmd.Flags().StringVarP(&accessKeyFlag, "accessKey", "k", "", "Specify client's accessKey")

	createAccessKeyCmd.MarkFlagRequired("name")
	createAccessKeyCmd.MarkFlagRequired("appID")
	createAccessKeyCmd.MarkFlagRequired("accessKey")

	return createAccessKeyCmd
}
