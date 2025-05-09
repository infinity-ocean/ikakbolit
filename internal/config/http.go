package config

type HTTP struct {
	Port string `env:"HTTP_PORT,notEmpty"`
}
