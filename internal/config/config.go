package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Postgres Postgres
	HTTP     HTTP
	GRPC     GRPC
	Log      Log
	Options  Options
	Debug    bool `env:"DEBUG" envDefault:"false"`
}

func Load() (Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			if err := godotenv.Load("../../.env"); err != nil {
			log.Fatalf("Failed to load .env: %v", err)
		}
		}
	}
	var config Config

	if err := env.Parse(&config); err != nil {
		return Config{}, fmt.Errorf("env.Parse: %w", err)
	}

	return config, nil
}
