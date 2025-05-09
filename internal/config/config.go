package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Postgres Postgres
	HTTP     HTTP
	GRPC     GRPC
	Log      Log
	Debug    bool `env:"DEBUG" envDefault:"false"`
}

func Load() (Config, error) {
	var config Config

	if err := env.Parse(&config); err != nil {
		return Config{}, fmt.Errorf("env.Parse: %w", err)
	}

	return config, nil
}
