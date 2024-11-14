package config

type JWTConfig struct {
	PrivateKeyPath  string `mapstructure:"jwt_private_key_path"`
	PublicKeyPath   string `mapstructure:"jwt_public_key_path"`
	ExpirationHours int    `mapstructure:"jwt_expiration_hours"`
}
