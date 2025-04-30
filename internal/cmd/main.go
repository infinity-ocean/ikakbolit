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
		log.Fatal(err)
	}
	
	log := logger.MustInitLogger()

	log.Info("program is started")

	pool, err := repo.MakePool()
	if err != nil {
		log.Error("can't create pool", "err", err)
	}

	repo := repo.New(pool)
	svc := service.New(repo)
	grpcCtrl := controller.NewGRPCServer(svc, ":50051")

	go func() {
	    if err := grpcCtrl.Run(); err != nil {
	        log.Error("failed to start gRPC server: %v", "err", err)
	    }
	}()

	restCtrl := controller.New(svc, ":8080")

	go func() {
	    if err := restCtrl.Run(); err != nil {
			log.Error("failed to start REST server: %v", "err", err)

	    }
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Info("Shutting down servers...")
}