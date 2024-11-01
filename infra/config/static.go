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

type ServerConfig struct {
	Port         string `mapstructure:"port"`
	AllowOrigins string `mapstructure:"allow_origins"`
	BaseURL      string `mapstructure:"base_url"`
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
	MeterName     string `mapstructure:"meter_name"`
	TracerName    string `mapstructure:"tracer_name"`
	MeterExporter string `mapstructure:"meter_exporter"`
	TraceExporter string `mapstructure:"trace_exporter"`
	LogExporter   string `mapstructure:"logger_exporter"`
}

type EncryptionConfig struct {
	Key string `mapstructure:"key"`
}

type VaultConfig struct {
	Url   string `mapstructure:"url"`
	Token string `mapstructure:"token"`
}

type ConsulConfig struct {
	Url string `mapstructure:"url"`
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
