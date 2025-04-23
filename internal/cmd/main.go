package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"log"
	"net"
	"google.golang.org/grpc"
	pb "github.com/infinity-ocean/ikakbolit/3-api-grpc-Homework/grpc/ikakbolit"

	"github.com/infinity-ocean/ikakbolit/internal/config"
	"github.com/infinity-ocean/ikakbolit/internal/controller"
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
	log.Println("program is started")

	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../../.env"); err != nil {
			log.Fatal(err)
		}
	}

	conf := config.Config{}
	if err := conf.Parse(); err != nil {
		log.Println(err)
	}

	pool, err := repo.MakePool(conf)
	if err != nil {
		log.Println(err)
	}

	repo := repo.New(pool)
	svc := service.New(repo)
	ctrl := controller.New(svc, ":8080")

	grpcSrv := grpc.NewServer()
	pb.RegisterIkakbolitServiceServer(grpcSrv, controller.NewGRPCServer(svc))

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Println("Starting gRPC server on :50051")
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	go func() {
		if err := ctrl.Run(); err != nil {
			log.Fatalf("failed to start REST server: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("Shutting down servers...")
}