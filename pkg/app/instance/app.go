package instance

import (
	gravity_cli "github.com/BrobridgeOrg/gravity-cli/pkg/cli/instance"
	//log "github.com/sirupsen/logrus"
)

type AppInstance struct {
	done chan bool
	cli  *gravity_cli.Cli
}

func NewAppInstance() *AppInstance {

	a := &AppInstance{
		done: make(chan bool),
	}

	a.cli = gravity_cli.NewCli(a)

	return a
}

func (a *AppInstance) Init() error {

	err := a.cli.Init()
	if err != nil {
		return err
	}

	return nil
}

func (a *AppInstance) Uninit() {
}

func (a *AppInstance) Run() error {
	// Nothing to do
	//	a.cli.RunCli()
	//	<-a.done

	return nil
}
