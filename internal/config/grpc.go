package config

type GRPC struct {
	ListenAddress int `env:"GRPC_LISTEN_ADDRESS,notEmpty"`
}
