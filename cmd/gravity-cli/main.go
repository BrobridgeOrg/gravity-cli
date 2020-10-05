package main

import (
	"context"
	"gravity-cli/pkg/cli"

	"go.uber.org/fx"
)

func main() {

	app := fx.New(
		cli.Modual,
		fx.NopLogger,
	)

	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
	defer app.Stop(context.Background())
}
