package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	ProdENV  = "prod"
	LocalENV = "local"
)

type Config struct {
	ENV      string `env:"ENV" env-required:"true"`
	Database struct {
		Host     string `env:"DATABASE_HOST" env-required:"true"`
		Port     int    `env:"DATABASE_PORT" env-required:"true"`
		User     string `env:"DATABASE_USER" env-required:"true"`
		Password string `env:"DATABASE_PASSWORD" env-required:"true"`
		Name     string `env:"DATABASE_NAME" env-required:"true"`
	}

	Server struct {
		Host string `env:"SERVER_HOST" env-default:"127.0.0.1"`
		Port int    `env:"SERVER_PORT" env-default:"8080"`
	}

	Memcached struct {
		Host string `env:"MEMCACHED_HOST" env-required:"true"`
		Port int    `env:"MEMCACHED_PORT" env-required:"true"`
	}
}

func MustInit() *Config {

	cfg := Config{}

	err := cleanenv.ReadEnv(&cfg)

	if err != nil {
		err = cleanenv.ReadConfig(".env", &cfg)
		if err != nil {
			panic(err)
		}
	}

	return &cfg
}
