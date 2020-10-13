package cmd_host

import (
	"fmt"

	"github.com/spf13/cobra"
)

var HostFlag string

var HostCmd = &cobra.Command{
	Use:   "host",
	Short: "Add gravity host",
	Long:  `Add gravity host using endpoint style`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("=========Gravity Host=========")
		fmt.Println("Host: " + HostFlag)
		fmt.Println("=================================")
	},
}

func NewHostCmd() *cobra.Command {

	HostCmd.Flags().StringVarP(&HostFlag, "gravity-host", "H", "", "Gravity host")

	return HostCmd
}
