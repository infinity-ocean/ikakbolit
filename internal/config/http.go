package config

type HTTP struct {
	ListenAddress   string        `env:"HTTP_LISTEN_ADDRESS,notEmpty"`
}