package configs

import (
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Host        string
	Domain      string
	AccessToken string
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
	config.SetHost(viper.GetString("host"))
	config.SetDomain(viper.GetString("domain"))
	config.SetAccessToken(viper.GetString("accessToken"))

	return config
}

func (config *Config) SetHost(host string) {

	if len(host) == 0 {
		return
	}

	config.Host = host

	parts := strings.Split(host, ":")
	viper.Set("gravity.host", parts[0])

	if len(parts) == 2 {
		port, err := strconv.Atoi(parts[1])
		if err == nil {
			viper.Set("gravity.port", port)
		}
	}

}

func (config *Config) SetDomain(domain string) {
	config.Domain = domain
	viper.Set("gravity.domain", domain)
}

func (config *Config) SetAccessToken(accessToken string) {
	config.AccessToken = accessToken
	viper.Set("gravity.accessToken", accessToken)
}
