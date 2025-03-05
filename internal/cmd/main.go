package main

import (
	"log"

	"github.com/infinity-ocean/ikakbolit/internal/config"
	"github.com/infinity-ocean/ikakbolit/internal/repo"
	"github.com/infinity-ocean/ikakbolit/internal/service"
	"github.com/infinity-ocean/ikakbolit/internal/controller"
)

	
func main() {
	log.Println("program is started")
	
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