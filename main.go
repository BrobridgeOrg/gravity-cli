package main

import (
	"os"

	"github.com/BrobridgeOrg/gravity-cli/cmd"
)

func main() {

	err := cmd.RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
