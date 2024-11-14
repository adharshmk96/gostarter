package config

type ServerConfig struct {
	Port         string `mapstructure:"port"`
	AllowOrigins string `mapstructure:"allow_origins"`
	BaseURL      string `mapstructure:"base_url"`
}
