package config

type OtelConfig struct {
	Endpoint      string `mapstructure:"endpoint"`
	UsePrometheus bool   `mapstructure:"use_prometheus"`
}
