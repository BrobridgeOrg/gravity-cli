package main

import (
	"context"
	"fmt"
	"gravity-cli/pkg/cli"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func init() {

	// From the environment
	viper.SetEnvPrefix("GRAVITY_CLI")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// From config file
	viper.SetConfigName("configs")
	viper.AddConfigPath("../../configs")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No configuration file was loaded")
	}
}

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
