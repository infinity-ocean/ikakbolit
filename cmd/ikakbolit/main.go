package main

import (
	"fmt"
	"os"

	"github.com/infinity-ocean/ikakbolit/internal/application"
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

var appVersion = "v1.0.0" //nolint:gochecknoglobals

func main() {
	if err := application.New(appVersion).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
