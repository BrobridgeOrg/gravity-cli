package adapter

import (
	"github.com/BrobridgeOrg/gravity-cli/pkg/cli"
	"github.com/spf13/cobra"
)

type AdapterCmd struct {
	cli        cli.Cli
	adapterCmd *cobra.Command
}

func NewAdapterCmd(c cli.Cli) *AdapterCmd {
	adapter := &AdapterCmd{
		cli: c,
	}

	return adapter
}

func (adapter *AdapterCmd) Init() *cobra.Command {
	adapter.adapterCmd = &cobra.Command{
		Use:   "adapter",
		Short: "Gravity Adapter Manager",
		Long:  `Gravity Adapter Manager`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	list := adapter.newListAdapterCmd()
	adapter.adapterCmd.AddCommand(list)

	return adapter.adapterCmd
}
