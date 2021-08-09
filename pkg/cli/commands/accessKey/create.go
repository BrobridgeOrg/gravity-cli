package accessKey

import (
	"strings"

	auth "github.com/BrobridgeOrg/gravity-sdk/authenticator"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func (accessKey *AccessKeyCmd) newCreateAccessKeyCmd() *cobra.Command {
	var appNameFlag string
	var appIDFlag string
	var accessKeyFlag string
	var roleFlag string
	var collectionFlag string
	var createAccessKeyCmd = &cobra.Command{
		Use:   "create",
		Short: "Create Gravity Access Key",
		Long:  `Create Gravity Access Key`,
		Run: func(cmd *cobra.Command, args []string) {

			//Get auth client
			authClient, err := accessKey.cli.GetAuthClient()
			if err != nil {
				log.Fatal(err)
			}

			//process access key role
			roles := []string{}
			if len(roleFlag) > 0 {
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
			}

			//process access key collections
			collections := []string{}
			if len(collectionFlag) > 0 {
				rs := strings.Split(collectionFlag, ",")
				for _, r := range rs {
					r = strings.TrimSpace(r)

					if len(r) == 0 {
						continue
					}

					appendCol := true
					for _, collection := range collections {
						if collection == r {
							appendCol = false
							break
						}
					}

					if appendCol {
						collections = append(collections, r)
					}
				}
			}

			entity := auth.NewEntity()
			entity.AppID = appIDFlag
			entity.AccessKey = accessKeyFlag
			entity.AppName = appNameFlag
			entity.Properties["permissions"] = roles
			entity.Properties["collections"] = collections

			if err := authClient.CreateEntity(entity); err != nil {
				log.Fatal(err)
			}

		},
	}

	createAccessKeyCmd.Flags().StringVarP(&appNameFlag, "name", "n", "", "Specify client's accessKey name")
	createAccessKeyCmd.Flags().StringVarP(&appIDFlag, "appID", "i", "", "Specify client's appID")
	createAccessKeyCmd.Flags().StringVarP(&accessKeyFlag, "accessKey", "k", "", "Specify client's accessKey")
	createAccessKeyCmd.Flags().StringVarP(&roleFlag, "roles", "r", "", "Specify accessKey's roles [ SYSTEM | ADAPTER | SUBSCRBIER ], This flag can using \",\" to  specified multiple roles.")
	createAccessKeyCmd.Flags().StringVarP(&collectionFlag, "collections", "c", "", "Assign accessKey can subscribe's collections, This flag can using \",\" to  specified multiple collections. (Default can subscribe any collections.)")

	createAccessKeyCmd.MarkFlagRequired("name")
	createAccessKeyCmd.MarkFlagRequired("appID")
	createAccessKeyCmd.MarkFlagRequired("accessKey")
	createAccessKeyCmd.MarkFlagRequired("roles")

	return createAccessKeyCmd
}
