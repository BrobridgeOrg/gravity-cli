package subscriber

import (
	"github.com/BrobridgeOrg/gravity-cli/pkg/cli"
	"github.com/spf13/cobra"
)

type SubscriberCmd struct {
	cli           cli.Cli
	subscriberCmd *cobra.Command
}

func NewSubscriberCmd(c cli.Cli) *SubscriberCmd {
	subscriber := &SubscriberCmd{
		cli: c,
	}
	return subscriber
}

func (subscriber *SubscriberCmd) Init() *cobra.Command {
	subscriber.subscriberCmd = &cobra.Command{
		Use:   "subscriber",
		Short: "Gravity Subscriber Manager",
		Long:  `Gravity Subscriber Manager`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	list := subscriber.newListSubscriberCmd()
	subscriber.subscriberCmd.AddCommand(list)

	return subscriber.subscriberCmd
}
