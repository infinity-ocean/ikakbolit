package modules

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/infinity-ocean/ikakbolit/internal/server/rest"
	"golang.org/x/sync/errgroup"
)

type HTTPServer struct {
	Port   string
	Router *rest.HTTPRouter
	Log    *slog.Logger
}

func (h HTTPServer) Run(
	_ context.Context,
	g *errgroup.Group,
) {
	g.Go(func() error {
		h.Log.Info("Starting HTTP server", "port", h.Port)
		if err := http.ListenAndServe(h.Port, h.Router.GetRouter()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http.ListenAndServe: %w", err)
		}
		h.Log.Info("http server stopped")

		return nil
	})
}
