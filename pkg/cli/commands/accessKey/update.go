package accessKey

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func (accessKey *AccessKeyCmd) newUpdateAccessKeyCmd() *cobra.Command {
	var accessKeyFlag string
	var appNameFlag string
	var updateAccessKeyCmd = &cobra.Command{
		Use:   "update [AppID]",
		Short: "Update Gravity Subscriber's Access Key",
		Long:  `Update Gravity Subscriber's Access Key`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			//Get auth client
			authClient, err := accessKey.cli.GetAuthClient()
			if err != nil {
				log.Fatal(err)
			}

			if len(accessKeyFlag) > 0 {
				if err := authClient.UpdateEntityKey(args[0], accessKeyFlag); err != nil {
					log.Fatal(err)
				}
			}
			if len(appNameFlag) > 0 {
				entity, err := authClient.GetEntity(args[0])
				if err != nil {
					log.Fatal(err)
				}
				entity.AppName = appNameFlag
				if err := authClient.UpdateEntity(entity); err != nil {
					log.Fatal(err)
				}

			}
		},
	}

	updateAccessKeyCmd.Flags().StringVarP(&accessKeyFlag, "accessKey", "k", "", "Specify new accessKey")
	updateAccessKeyCmd.Flags().StringVarP(&appNameFlag, "name", "n", "", "Specify new appName")

	return updateAccessKeyCmd
}
