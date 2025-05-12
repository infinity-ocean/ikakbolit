package config

type HTTP struct {
	Port          string `env:"HTTP_PORT,notEmpty"`
	ListenAddress string `env:"HTTP_ADDRESS,notEmpty"`
}
