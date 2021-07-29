package main

import (
	"fmt"
	"os"
	//	"path"
	//	"strings"

	log "github.com/sirupsen/logrus"
	//	"github.com/spf13/viper"

	app "github.com/BrobridgeOrg/gravity-cli/pkg/app/instance"
)

func init() {

	debugLevel := log.InfoLevel
	switch os.Getenv("GRAVITY_DEBUG") {
	case log.TraceLevel.String():
		debugLevel = log.TraceLevel
	case log.DebugLevel.String():
		debugLevel = log.DebugLevel
	case log.ErrorLevel.String():
		debugLevel = log.ErrorLevel
	case log.InfoLevel.String():
		debugLevel = log.InfoLevel
	}

	log.SetOutput(os.Stdout)
	log.SetLevel(debugLevel)

	debugMsg := fmt.Sprintf("Debug level is set to \"%s\"\n", debugLevel.String())
	log.Debug(debugMsg)

	/*
		// From the environment
		viper.SetEnvPrefix("GRAVITY_CLI")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()

		//Generate config path
		homeDir, _ := os.UserHomeDir()
		configPath := path.Join(homeDir, ".gravity/")

		// From config file
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath(configPath)

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// Create configPath
				err := os.MkdirAll(configPath, 0700)
				if err != nil {
					log.Fatal(err)
				}
				// Write Config file
				configFile := path.Join(configPath, "config.toml")
				viper.WriteConfigAs(configFile)

			} else {
				log.Error("Reading confing error: ", err)
			}
		}
	*/
}

func main() {

	// Initializing application
	a := app.NewAppInstance()

	err := a.Init()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Starting application
	err = a.Run()
	if err != nil {
		log.Fatal(err)
		return
	}
}
