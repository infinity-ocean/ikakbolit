package config

type HTTP struct {
	ListenAddress int `env:"HTTP_LISTEN_ADDRESS,notEmpty"`
}
