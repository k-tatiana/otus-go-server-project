package internal

import "github.com/kelseyhightower/envconfig"

type DBConfig struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     int    `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER" default:"postgres"`
	Password string `envconfig:"DB_PASSWORD" default:"password"`
	Database string `envconfig:"DB_NAME" default:"postgres"`
	SSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`
}

type Config struct {
	DB     DBConfig
	Secret string `envconfig:"SECRET"`
}

func EnvParse() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
