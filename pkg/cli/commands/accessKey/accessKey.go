package accessKey

import (
	"github.com/BrobridgeOrg/gravity-cli/pkg/cli"
	"github.com/spf13/cobra"
)

type AccessKeyCmd struct {
	cli          cli.Cli
	accessKeyCmd *cobra.Command
}

func NewAccessKeyCmd(c cli.Cli) *AccessKeyCmd {
	accessKey := &AccessKeyCmd{
		cli: c,
	}
	return accessKey

}

func (accessKey *AccessKeyCmd) Init() *cobra.Command {

	accessKey.accessKeyCmd = &cobra.Command{
		Use:   "accessKey",
		Short: "Gravity Access Key Manager",
		Long:  `Gravity Access Key Manager`,
		Run: func(cmd *cobra.Command, args []string) {
			// Nothing to do
		},
	}

	create := accessKey.newCreateAccessKeyCmd()
	delete := accessKey.newDeleteAccessKeyCmd()
	update := accessKey.newUpdateAccessKeyCmd()
	list := accessKey.newListAccessKeyCmd()

	accessKey.accessKeyCmd.AddCommand(create)
	accessKey.accessKeyCmd.AddCommand(delete)
	accessKey.accessKeyCmd.AddCommand(update)
	accessKey.accessKeyCmd.AddCommand(list)

	return accessKey.accessKeyCmd
}
