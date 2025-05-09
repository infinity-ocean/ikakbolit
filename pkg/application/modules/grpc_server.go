package modules

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/infinity-ocean/ikakbolit/internal/config"
	"github.com/infinity-ocean/ikakbolit/internal/server/grpc"
)

type GRPCServer struct {}

func (h GRPCServer) Run(
	grpc *grpc.GRPCServer,
	cfg config.Config,
	log *slog.Logger,
) {
	go func() {
		log.Info("Starting gRPC server", "address", strconv.Itoa(cfg.GRPC.ListenAddress))
		if err := grpc.Run(); err != nil {
			log.Error("gRPC server error:", "err", err)
			os.Exit(1)
	}
	}()
}
