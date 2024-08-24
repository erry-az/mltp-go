package config

import "time"

type DbConfig struct {
	Dsn             string        `mapstructure:"dsn"`
	MaxConn         int           `mapstructure:"max_conn"`
	MinConn         int           `mapstructure:"min_conn"`
	MaxConnLifetime time.Duration `mapstructure:"max_conn_lifetime"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time"`
}
