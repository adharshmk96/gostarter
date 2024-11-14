package config

type DatabaseConfig struct {
	MigrationFiles string         `mapstructure:"migration_files"`
	Postgres       PostgresConfig `mapstructure:"postgres"`
}

type PostgresConfig struct {
	Connection string `mapstructure:"connection"`
}
