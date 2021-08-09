package subscriber

import (
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func (subscriber *SubscriberCmd) newListSubscriberCmd() *cobra.Command {
	var listAllFlag bool
	var listSubscriberCmd = &cobra.Command{
		Use:   "list",
		Short: "List Gravity Subscribers",
		Long:  `List Gravity Subscribers`,
		Run: func(cmd *cobra.Command, args []string) {

			//Get auth client
			smClient, err := subscriber.cli.GetSubscriberManagerClient()
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

			// get subscribers
			subscribers, err := smClient.GetSubscribers()
			if err != nil {
				log.Fatal(err)
			}

			if listAllFlag {
				//show all information
				table.SetHeader([]string{"ID", "Name", "Component", "AppID", "Type", "Roles"})
				for _, subscriber := range subscribers {
					rolesStr := strings.Join(subscriber.Permissions, ",")

					table.Append([]string{subscriber.ID, subscriber.Name, subscriber.Component, subscriber.AppID, subscriber.Type.String(), rolesStr})
				}
			} else {
				table.SetHeader([]string{"ID", "Name", "Component", "Type"})
				for _, subscriber := range subscribers {
					table.Append([]string{subscriber.ID, subscriber.Name, subscriber.Component, subscriber.Type.String()})
				}

			}

			table.Render()

		},
	}

	listSubscriberCmd.Flags().BoolVar(&listAllFlag, "all", false, "List all information")
	return listSubscriberCmd
}
