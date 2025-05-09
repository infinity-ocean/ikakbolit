package config

type GRPC struct {
	ListenAddress   string        `env:"GRPC_LISTEN_ADDRESS,notEmpty"`
}