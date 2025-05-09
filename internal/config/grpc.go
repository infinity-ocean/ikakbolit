package config

type GRPC struct {
	Port string `env:"GRPC_PORT,notEmpty"`
}
