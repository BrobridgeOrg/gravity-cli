package cmd_root

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "gravity-cli",
	Short: "Gravity-cli tool",
	Long:  `Gravity-cli tool by daginwu in Brobridge`,
}
