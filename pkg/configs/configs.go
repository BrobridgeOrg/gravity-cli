package configs

import (
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Host   string
	Domain string
}

func GetConfig() *Config {

	// From the environment
	viper.SetEnvPrefix("GRAVITY_CLI")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// From config file
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./configs")

	viper.ReadInConfig()

	runtime.GOMAXPROCS(8)

	config := &Config{}

	// Specify events from environment variable for watching
	config.Host = viper.GetString("host")
	config.Domain = viper.GetString("domain")

	return config
}

func (config *Config) SetHost(host string) {
	config.Host = host
}

func (config *Config) SetDomain(domain string) {
	config.Domain = domain
}
