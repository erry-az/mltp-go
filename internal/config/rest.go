package config

import "time"

// RestConfig represents the Rest server configuration.
type RestConfig struct {
	Host         string        `mapstructure:"host"`
	Port         uint          `mapstructure:"port"`
	ReadTimeOut  time.Duration `mapstructure:"read_timeout"`
	WriteTimeOut time.Duration `mapstructure:"write_timeout"`
}
