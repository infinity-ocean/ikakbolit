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

const shutdownTimeout = 5 * 1e9

type HTTPServer struct {
	Port   string
	Router *rest.HTTPRouter
	Log    *slog.Logger
}

func (s HTTPServer) Run(
	ctx context.Context,
	g *errgroup.Group,
) {
	server := &http.Server{
		Addr:    s.Port,
		Handler: s.Router.GetRouter(),
	}

	g.Go(func() error {
		s.Log.Info("Starting HTTP server", "port", s.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http.ListenAndServe: %w", err)
		}
		s.Log.Info("HTTP server stopped")
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		s.Log.Info("Shutting down HTTP server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	})
}
