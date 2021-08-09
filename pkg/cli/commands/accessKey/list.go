package accessKey

import (
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func (accessKey *AccessKeyCmd) newListAccessKeyCmd() *cobra.Command {
	var listAllFlag bool
	var listAccessKeyCmd = &cobra.Command{
		Use:   "list [AppID]",
		Short: "List Gravity Access Key",
		Long:  `List Gravity Access Key`,
		Run: func(cmd *cobra.Command, args []string) {

			//Get auth client
			authClient, err := accessKey.cli.GetAuthClient()
			if err != nil {
				log.Fatal(err)
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetAutoWrapText(false)
			table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetCenterSeparator("")
			table.SetColumnSeparator("")
			table.SetRowSeparator("")
			table.SetHeaderLine(false)
			table.SetBorder(false)
			table.SetTablePadding("\t") // pad with tabs
			table.SetNoWhiteSpace(true)

			appID := ""
			tmpAppID := ""
			if len(args) == 0 {
				// get entities
			LOOP:
				entities, _, err := authClient.GetEntities(appID, 10)
				if err != nil {
					log.Fatal(err)
				}

				if listAllFlag {
					//show all information
					table.SetHeader([]string{"AppID", "AppName", "AccessKey", "Roles", "Collections"})
					for i, entity := range entities {
						appID = entity.AppID
						if tmpAppID == appID && i == 0 {
							continue
						}

						roles := []string{}
						permissions := entity.Properties["permissions"].([]interface{})
						for _, permission := range permissions {
							roles = append(roles, permission.(string))

						}
						rolesStr := strings.Join(roles, ",")

						collections := []string{}
						cols := entity.Properties["collections"].([]interface{})
						for _, col := range cols {
							collections = append(collections, col.(string))

						}
						collectionsStr := strings.Join(collections, ",")

						table.Append([]string{entity.AppID, entity.AppName, entity.AccessKey, rolesStr, collectionsStr})
					}
				} else {
					table.SetHeader([]string{"AppID", "AppName"})
					for i, entity := range entities {
						appID = entity.AppID
						if tmpAppID == appID && i == 0 {
							continue
						}
						table.Append([]string{entity.AppID, entity.AppName})
					}

				}

				if appID != tmpAppID {
					tmpAppID = appID
					goto LOOP
				}

			} else {

				table.SetHeader([]string{"AppID", "AppName", "AccessKey", "Roles", "Collections"})
				for _, arg := range args {
					// get entity
					entity, err := authClient.GetEntity(arg)
					if err != nil {
						log.Fatal(err)
					}

					roles := []string{}
					permissions := entity.Properties["permissions"].([]interface{})
					for _, permission := range permissions {
						roles = append(roles, permission.(string))

					}
					rolesStr := strings.Join(roles, ",")

					collections := []string{}
					cols := entity.Properties["collections"].([]interface{})
					for _, col := range cols {
						collections = append(collections, col.(string))

					}
					collectionsStr := strings.Join(collections, ",")

					table.Append([]string{entity.AppID, entity.AppName, entity.AccessKey, rolesStr, collectionsStr})
				}
			}

			table.Render()

		},
	}

	listAccessKeyCmd.Flags().BoolVar(&listAllFlag, "all", false, "List all information")
	return listAccessKeyCmd
}
