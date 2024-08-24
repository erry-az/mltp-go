package config

type JwtConfig struct {
	SigningKey    string `mapstructure:"signing_key"`
	SigningMethod string `mapstructure:"signing_method"`
	Issuer        string `mapstructure:"issuer"`
}
