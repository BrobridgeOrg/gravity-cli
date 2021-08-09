package adapter

import (
	"os"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func (adapter *AdapterCmd) newListAdapterCmd() *cobra.Command {
	//var listAllFlag bool
	var listSubscriberCmd = &cobra.Command{
		Use:   "list",
		Short: "List Gravity Adapters",
		Long:  `List Gravity Adapters`,
		Run: func(cmd *cobra.Command, args []string) {

			//Get auth client
			client, err := adapter.cli.GetAdapterManagerClient()
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
			adapters, err := client.GetAdapters()
			if err != nil {
				log.Fatal(err)
			}

			table.SetHeader([]string{"ID", "Name", "Component"})
			for _, adapter := range adapters {
				table.Append([]string{adapter.ID, adapter.Name, adapter.Component})
			}

			table.Render()

		},
	}

	return listSubscriberCmd
}
