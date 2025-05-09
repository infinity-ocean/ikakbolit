package modules

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/infinity-ocean/ikakbolit/internal/config"
	"github.com/infinity-ocean/ikakbolit/internal/server/rest"
)

type HTTPServer struct {}

func (h HTTPServer) Run(
	rest *rest.HTTPServer,
	cfg config.Config,
	log *slog.Logger,
) {
	go func() {
			log.Info("Starting REST server", "address", strconv.Itoa(cfg.HTTP.ListenAddress))
			if err := rest.Run(); err != nil {
				log.Error("REST server error:", "err", err)
				os.Exit(1)
		}
		}()
}
