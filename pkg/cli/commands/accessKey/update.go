package accessKey

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

func (accessKey *AccessKeyCmd) newUpdateAccessKeyCmd() *cobra.Command {
	var accessKeyFlag string
	var appNameFlag string
	var roleFlag string
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
			//process access key role
			if len(roleFlag) > 0 {
				entity, err := authClient.GetEntity(args[0])
				if err != nil {
					log.Fatal(err)
				}

				roles := []string{}
				roleFlag = strings.ToUpper(roleFlag)
				rs := strings.Split(roleFlag, ",")
				for _, r := range rs {
					r = strings.TrimSpace(r)

					if len(r) == 0 {
						continue
					}

					if r != "SYSTEM" && r != "ADAPTER" && r != "SUBSCRIBER" {
						log.Error("Unkonw role: ", r)
						continue
					}

					appendRole := true
					for _, role := range roles {
						if role == r {
							appendRole = false
							break
						}
					}

					if appendRole {
						roles = append(roles, r)
					}
				}

				entity.Properties["permissions"] = roles
				if err := authClient.UpdateEntity(entity); err != nil {
					log.Fatal(err)
				}
			}
		},
	}

	updateAccessKeyCmd.Flags().StringVarP(&accessKeyFlag, "accessKey", "k", "", "Specify new accessKey")
	updateAccessKeyCmd.Flags().StringVarP(&appNameFlag, "name", "n", "", "Specify new appName")
	updateAccessKeyCmd.Flags().StringVarP(&roleFlag, "roles", "r", "", "Specify accessKey's roles [ SYSTEM | ADAPTER | SUBSCRBIER ], This flag can using \",\" to  specified multiple roles.")

	return updateAccessKeyCmd
}
