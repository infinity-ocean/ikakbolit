package modules

import (
	"log/slog"
	"os"

	"github.com/infinity-ocean/ikakbolit/internal/config"
	"github.com/infinity-ocean/ikakbolit/internal/server/grpc"
)

func RunGRPC(
	grpc *grpc.GRPCServer,
	cfg config.Config,
	log *slog.Logger,
) {
	go func() {
		log.Info("Starting gRPC server", "address", cfg.GRPC.Port)
		if err := grpc.Run(); err != nil {
			log.Error("gRPC server error:", "err", err)
			os.Exit(1)
	}
	}()
}
