package main

import (
	"log"

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

	if err := ctrl.Run(); err != nil {
		log.Println(err)
	}
}
