package config

import "github.com/spf13/viper"

type Config struct {
	Version       string              `mapstructure:"version"`
	Server        ServerConfig        `mapstructure:"server"`
	Database      DatabaseConfig      `mapstructure:"database"`
	JWT           JWTConfig           `mapstructure:"jwt"`
	Encryption    EncryptionConfig    `mapstructure:"encryption"`
	Observability ObservabilityConfig `mapstructure:"observability"`
	Vault         VaultConfig         `mapstructure:"vault"`
	Consul        ConsulConfig        `mapstructure:"consul"`
}

var config *Config

func NewConfig() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigName(".gostarter")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	return config
}
