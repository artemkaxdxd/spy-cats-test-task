package config

import "github.com/kelseyhightower/envconfig"

type (
	Config struct {
		Server   Server
		Postgres Postgres
	}

	Server struct {
		Port  string `envconfig:"SERVER_PORT"`
		IsDev bool   `envconfig:"SERVER_IS_DEV"`
	}

	Postgres struct {
		User     string `envconfig:"POSTGRES_USER"`
		Password string `envconfig:"POSTGRES_PASSWORD"`
		Host     string `envconfig:"POSTGRES_HOST"`
		Port     string `envconfig:"POSTGRES_PORT"`
		DBName   string `envconfig:"POSTGRES_DB_NAME"`
	}
)

func New() (config Config, err error) {
	err = envconfig.Process("", &config)
	return
}
