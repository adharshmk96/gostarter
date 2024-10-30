package config

import "github.com/spf13/viper"

type Config struct {
	Version    string           `mapstructure:"version"`
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Encryption EncryptionConfig `mapstructure:"encryption"`
}

type ServerConfig struct {
	Port         string `mapstructure:"port"`
	AllowOrigins string `mapstructure:"allow_origins"`
}

type DatabaseConfig struct {
	MigrationFiles string         `mapstructure:"migration_files"`
	Postgres       PostgresConfig `mapstructure:"postgres"`
}

type PostgresConfig struct {
	Connection string `mapstructure:"connection"`
}

type JWTConfig struct {
	PrivateKeyPath  string `mapstructure:"jwt_private_key_path"`
	PublicKeyPath   string `mapstructure:"jwt_public_key_path"`
	ExpirationHours int    `mapstructure:"jwt_expiration_hours"`
}

type ObservabilityConfig struct {
	MeterName      string `mapstructure:"meter_name"`
	MeterExporter  string `mapstructure:"meter_exporter"`
	TracerName     string `mapstructure:"tracer_name"`
	TracerExporter string `mapstructure:"tracer_exporter"`
}

type EncryptionConfig struct {
	Key string `mapstructure:"key"`
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

func GetConfig() *Config {
	if config == nil {
		config = NewConfig()
	}
	return config
}
