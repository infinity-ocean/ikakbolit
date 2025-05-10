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

func (s HTTPServer) Run(
	_ context.Context,
	g *errgroup.Group,
) {
	g.Go(func() error {
		s.Log.Info("Starting HTTP server", "port", s.Port)
		if err := http.ListenAndServe(s.Port, s.Router.GetRouter()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http.ListenAndServe: %w", err)
		}
		s.Log.Info("http server stopped")

		return nil
	})
}
