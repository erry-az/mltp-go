package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceName string     `mapstructure:"service_name"`
	RestServer  RestConfig `mapstructure:"rest_server"`
	DB          DbConfig   `mapstructure:"db"`
	Jwt         JwtConfig  `mapstructure:"jwt"`
	Otel        OtelConfig `mapstructure:"otel"`
	Log         LogConfig  `mapstructure:"log"`
}

// Load loads the config file into Config struct
func Load(env string) (Config, error) {
	var cfg Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// More option of config path can be added here
	viper.AddConfigPath("/etc/wallet/config/")                      // Staging, Production or Docker
	viper.AddConfigPath(fmt.Sprintf("files/config/%s/", env))       // Unix Local
	viper.AddConfigPath(fmt.Sprintf("../../files/config/%s/", env)) // Windows Local

	viper.AutomaticEnv()

	// Get the config file
	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}

	// Convert into struct
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
