package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/infinity-ocean/ikakbolit/internal/config"
	"github.com/infinity-ocean/ikakbolit/internal/domain/service"
	"github.com/infinity-ocean/ikakbolit/internal/infrastructure/repository"
	"github.com/infinity-ocean/ikakbolit/internal/server/grpc"
	"github.com/infinity-ocean/ikakbolit/internal/server/rest"
	"github.com/infinity-ocean/ikakbolit/pkg/application/connectors"

	"github.com/samber/lo"
)

// @title ikakbolit API
// @version 1.0
// @description This is the main entry point for the Ikakbolit application, which sets up and runs the application server.
// @contact.name Константин Троицкий
// @contact.url https://t.me/debussy3
// @contact.telegram_username @debussy3
// @contact.email varrr7@gmail.com
// @host localhost:8080
// @BasePath /

func main() {
	cfg := lo.Must(config.Load())

	log := connectors.MustInitLogger()
	if cfg.Debug {
		log = slog.Default()
		log.Info("Running in DEBUG mode")
	}
	
	log.Info("Program is starting...")

	pool, err := repository.MakePool(cfg.Postgres.DSN)
	if err != nil {
		log.Error("Failed to create database pool:", "err", err)
		os.Exit(1)
	}

	repo := repository.New(pool)
	svc := service.New(repo, log)

	grpcCtrl := grpc.NewGRPCServer(svc, cfg.GRPC.ListenAddress, log)
	restCtrl := rest.NewHTTPServer(svc, cfg.HTTP.ListenAddress, log)

	go func() {
		log.Info("Starting gRPC server", "address", strconv.Itoa(cfg.GRPC.ListenAddress))
		if err := grpcCtrl.Run(); err != nil {
			log.Error("gRPC server error:", "err", err)
			os.Exit(1)
	}
	}()

	go func() {
		log.Info("Starting REST server", "address", strconv.Itoa(cfg.HTTP.ListenAddress))
		if err := restCtrl.Run(); err != nil {
			log.Error("REST server error:", "err", err)
			os.Exit(1)
	}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Info("Shutdown signal received, shutting down servers...")
}