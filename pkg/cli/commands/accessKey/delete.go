package accessKey

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func (accessKey *AccessKeyCmd) newDeleteAccessKeyCmd() *cobra.Command {
	var deleteAccessKeyCmd = &cobra.Command{
		Use:   "delete [AppID]",
		Short: "Delete Gravity Subscriber's Access Key",
		Long:  `Delete Gravity Subscriber's Access Key`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			//Get auth client
			authClient, err := accessKey.cli.GetAuthClient()
			if err != nil {
				log.Fatal(err)
			}

			//Delete
			for _, arg := range args {
				if err := authClient.DeleteEntity(arg); err != nil {
					log.Error(err)
				}
			}

		},
	}

	return deleteAccessKeyCmd
}
