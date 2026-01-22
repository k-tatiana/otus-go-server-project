package config

import "github.com/kelseyhightower/envconfig"

type DBConfig struct {
	Host     string `envconfig:"HOST" default:"localhost"`
	Port     int    `envconfig:"PORT" default:"5432"`
	User     string `envconfig:"USER" default:"postgres"`
	Password string `envconfig:"PASSWORD" default:"password"`
	Database string `envconfig:"NAME" default:"postgres"`
	SSLMode  string `envconfig:"SSLMODE" default:"disable"`
}

type RabbitMQConfig struct {
	Host     string `envconfig:"HOST" default:"localhost"`
	Port     int    `envconfig:"PORT" default:"5672"`
	User     string `envconfig:"USER" default:"guest"`
	Password string `envconfig:"PASSWORD" default:"guest"`
	VHost    string `envconfig:"VHOST" default:"/"`
}

type Tarantool struct {
	Host     string `envconfig:"HOST" default:"localhost"`
	Port     string `envconfig:"PORT" default:"3301"`
	User     string `envconfig:"USER" default:"guest"`
	Password string `envconfig:"PASSWORD" default:""`
	Enabled  bool   `envconfig:"ENABLED" default:"false"`
}

type DB struct {
	UseReplicas bool `envconfig:"USE_REPLICAS" default:"false"`
}

type Config struct {
	DB         DB             `envconfig:"DB"`
	DBMaster   DBConfig       `envconfig:"DB_MASTER"`
	DBReplica1 DBConfig       `envconfig:"DB_REPLICA1"`
	DBReplica2 DBConfig       `envconfig:"DB_REPLICA2"`
	Secret     string         `envconfig:"SECRET"`
	Cache      Cache          `envconfig:"CACHE"`
	RabbitMQ   RabbitMQConfig `envconfig:"RABBITMQ"`
	Tarantool  Tarantool      `envconfig:"TARANTOOL"`
	DialogAPI  struct {
		Address string `envconfig:"ADDRESS" default:"http://localhost:8082"`
	} `envconfig:"DIALOG_API"`
	CounterAPI struct {
		Address string `envconfig:"ADDRESS" default:"http://localhost:8088"`
	} `envconfig:"COUNTER_API"`
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
