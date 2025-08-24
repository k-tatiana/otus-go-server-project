package internal

import "github.com/kelseyhightower/envconfig"

type DBConfig struct {
	Host     string `envconfig:"HOST" default:"localhost"`
	Port     int    `envconfig:"PORT" default:"5432"`
	User     string `envconfig:"USER" default:"postgres"`
	Password string `envconfig:"PASSWORD" default:"password"`
	Database string `envconfig:"NAME" default:"postgres"`
	SSLMode  string `envconfig:"SSLMODE" default:"disable"`
}

type Config struct {
	DBMaster   DBConfig `envconfig:"DB_MASTER"`
	DBReplica1 DBConfig `envconfig:"DB_REPLICA1"`
	DBReplica2 DBConfig `envconfig:"DB_REPLICA2"`
	Secret     string   `envconfig:"SECRET"`
	Cache      Cache    `envconfig:"CACHE"`
}

type Cache struct {
	UseCache bool `envconfig:"ENABLE" default:"true"`
}

func EnvParse() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
