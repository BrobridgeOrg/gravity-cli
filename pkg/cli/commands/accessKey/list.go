package accessKey

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func (accessKey *AccessKeyCmd) newListAccessKeyCmd() *cobra.Command {
	var listAllFlag bool
	var listAccessKeyCmd = &cobra.Command{
		Use:   "list [AppID]",
		Short: "List Gravity Subscriber's Access Key",
		Long:  `List Gravity Subscriber's Access Key`,
		Run: func(cmd *cobra.Command, args []string) {

			//Get auth client
			authClient, err := accessKey.cli.GetAuthClient()
			if err != nil {
				log.Fatal(err)
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetRowLine(true)
			appID := ""
			tmpAppID := ""
			total := 0
			if len(args) == 0 {
				// get entities
			LOOP:
				entities, _, err := authClient.GetEntities(appID, 10)
				if err != nil {
					log.Fatal(err)
				}

				if listAllFlag {
					//show all information
					table.SetHeader([]string{"AppID", "AppName", "AccessKey"})
					for i, entity := range entities {
						appID = entity.AppID
						if tmpAppID == appID && i == 0 {
							continue
						}
						table.Append([]string{entity.AppID, entity.AppName, entity.AccessKey})
						total++
					}
				} else {
					table.SetHeader([]string{"AppID", "AppName"})
					for i, entity := range entities {
						appID = entity.AppID
						if tmpAppID == appID && i == 0 {
							continue
						}
						table.Append([]string{entity.AppID, entity.AppName})
						total++
					}

				}

				if appID != tmpAppID {
					tmpAppID = appID
					goto LOOP
				}

				caption := fmt.Sprintf("Total: %d", total)
				table.SetCaption(true, caption)

			} else {

				table.SetHeader([]string{"AppID", "AppName", "AccessKey"})
				for _, arg := range args {
					// get entity
					entity, err := authClient.GetEntity(arg)
					if err != nil {
						log.Fatal(err)
					}
					table.Append([]string{entity.AppID, entity.AppName, entity.AccessKey})
				}
			}

			table.Render()

		},
	}

	listAccessKeyCmd.Flags().BoolVar(&listAllFlag, "all", false, "List all information")
	return listAccessKeyCmd
}
