package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/infinity-ocean/ikakbolit/internal/server/grpc"
	"github.com/infinity-ocean/ikakbolit/internal/server/rest"
	"github.com/infinity-ocean/ikakbolit/pkg/application/connectors"
	"github.com/infinity-ocean/ikakbolit/internal/repository"
	"github.com/infinity-ocean/ikakbolit/internal/domain/service"
	"github.com/joho/godotenv"
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
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../../.env"); err != nil{
			log.Fatalf("Failed to load .env: %v", err)
		}
	}

	log := connectors.MustInitLogger()
	if os.Getenv("DEBUG") == "true" {
		log = slog.Default()
		log.Info("Running in DEBUG mode")
	}
	
	log.Info("Program is starting...")

	pool, err := repository.MakePool()
	if err != nil {
		log.Error("Failed to create database pool:", "err", err)
		os.Exit(1)
	}

	repo := repository.New(pool)
	svc := service.New(repo, log)

	grpcCtrl := grpc.NewGRPCServer(svc, ":50051", log)
	restCtrl := rest.NewHTTPServer(svc, ":8080", log)

	go func() {
		log.Info("Starting gRPC server on :50051")
		if err := grpcCtrl.Run(); err != nil {
			log.Error("gRPC server error:", "err", err)
			os.Exit(1)
	}
	}()

	go func() {
		log.Info("Starting REST server on :8080")
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