package cli

import (
	"fmt"
	host "gravity-cli/pkg/cli/cmd_host"
	login "gravity-cli/pkg/cli/cmd_login"
	root "gravity-cli/pkg/cli/cmd_root"
	store "gravity-cli/pkg/cli/cmd_store"
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

	// New Store Command
	store := store.NewStoreCmd()

	// New Login Command
	login := login.NewLoginCmd()

	// New Host Command
	host := host.NewHostCmd()

	// Root Combine
	root.RootCmd.AddCommand(store)
	root.RootCmd.AddCommand(login)
	root.RootCmd.AddCommand(host)
}

func RunGracityCli() {

	if err := root.RootCmd.Execute(); err != nil {
		fmt.Println(err)

		os.Exit(1)
	}

}
