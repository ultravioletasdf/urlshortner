package main

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	TursoDatabaseName string `env:"TURSO_DATABASE_NAME" envDefault:"./turso.db"`
	TursoUrl          string `env:"TURSO_URL,required"`
	TursoToken        string `env:"TURSO_TOKEN,required"`

	ListenAddress string `env:"LISTEN_ADDRESS" envDefault:"127.0.0.1:3005"`
}

func GetConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}
