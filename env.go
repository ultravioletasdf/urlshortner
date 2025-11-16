package main

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	TursoDatabaseName string `env:"TURSO_DATABASE_NAME" envDefault:"./turso.db"`
	TursoUrl          string `env:"TURSO_URL"`
	TursoToken        string `env:"TURSO_TOKEN"`

	ListenAddress string `env:"LISTEN_ADDRESS" envDefault:":3005"`
	FQDN          string `env:"FQDN,required"`
}

func GetConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}
