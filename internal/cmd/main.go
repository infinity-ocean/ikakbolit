package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/infinity-ocean/ikakbolit/internal/controller"
	"github.com/infinity-ocean/ikakbolit/internal/logger"
	"github.com/infinity-ocean/ikakbolit/internal/repo"
	"github.com/infinity-ocean/ikakbolit/internal/service"
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
		log.Fatalf("Failed to load .env: %v", err)
	}

	log := logger.MustInitLogger()
	log.Info("Program is starting...")

	pool, err := repo.MakePool()
	if err != nil {
		log.Error("Failed to create database pool:", "err", err)
		os.Exit(1)
	}

	repository := repo.New(pool)
	svc := service.New(repository)

	grpcCtrl := controller.NewGRPCServer(svc, ":50051")
	restCtrl := controller.NewRestServer(svc, ":8080")

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