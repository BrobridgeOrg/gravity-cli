package cli

import (
	"fmt"
	"gravity-cli/pkg/cli/cmds"
	"os"

	"go.uber.org/fx"
)

var Modual = fx.Options(
	fx.Invoke(
		NewGravityCli,
		RunGracityCli,
	),
)

func NewGravityCli() {
	cmds.ListCmd.Flags().StringVarP(&cmds.DbFlag, "dbinfo", "d", "", "Database information")
	cmds.ListCmd.Flags().BoolVarP(&cmds.AllFlag, "alldbinfo", "a", false, "All Database information")

	cmds.RecoverCmd.Flags().StringVarP(&cmds.DbFlag, "dbinfo", "d", "", "Recover Database")
	cmds.RecoverCmd.Flags().BoolVarP(&cmds.AllFlag, "alldbinfo", "a", false, "Recover all Database")

	cmds.RootCmd.AddCommand(cmds.StoreCmd)
	cmds.StoreCmd.AddCommand(cmds.ListCmd)
	cmds.StoreCmd.AddCommand(cmds.RecoverCmd)
}

func RunGracityCli() {

	if err := cmds.RootCmd.Execute(); err != nil {
		fmt.Println(err)

		os.Exit(1)
	}

}
