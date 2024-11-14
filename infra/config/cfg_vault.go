package config

type VaultConfig struct {
	Url   string `mapstructure:"url"`
	Token string `mapstructure:"token"`
}
