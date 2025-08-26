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

type DB struct {
	UseReplicas bool     `envconfig:"USE_REPLICAS" default:"false"`
	Master      DBConfig `envconfig:"MASTER"`
	Replica1    DBConfig `envconfig:"REPLICA1"`
	Replica2    DBConfig `envconfig:"REPLICA2"`
}

type Config struct {
	DB     DB     `envconfig:"DB"`
	Secret string `envconfig:"SECRET"`
	Cache  Cache  `envconfig:"CACHE"`
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
