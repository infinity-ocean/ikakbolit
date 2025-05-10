package application

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/infinity-ocean/ikakbolit/internal/config"
	"github.com/infinity-ocean/ikakbolit/internal/domain/service"
	"github.com/infinity-ocean/ikakbolit/internal/infrastructure/repository"
	"github.com/infinity-ocean/ikakbolit/internal/server/grpc"
	"github.com/infinity-ocean/ikakbolit/internal/server/rest"
	"github.com/infinity-ocean/ikakbolit/pkg/application/connectors"
	"github.com/infinity-ocean/ikakbolit/pkg/application/modules"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

type App struct {
	cfg          config.Config
	log          *connectors.Slog
	httpServer   *modules.HTTPServer
	Service      *service.Service
	Repo         *repository.Repo
	postgresPool *pgxpool.Pool
}

func New() *App {
	return &App{
		cfg:          config.Config{},
		log:          nil,
		httpServer:   nil,
		Service:      nil,
		Repo:         nil,
		postgresPool: nil,
	}
}

func (app *App) Init(appVersion string) error { //nolint:funlen
	const appName = "ikakbolit"

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)

	app.cfg = lo.Must(config.Load())
	app.log = &connectors.Slog{Name: appName, Version: appVersion, Debug: app.cfg.Debug}
	log := app.log.Logger(ctx)
	log.Info("Program is starting...")

	app.postgresPool = lo.Must(repository.MakePool(app.cfg.Postgres.DSN))
	app.Repo = repository.New(app.postgresPool)
	app.Service = service.New(app.Repo, log, app.cfg)

	app.httpServer = &modules.HTTPServer{
		Port:   app.cfg.HTTP.Port,
		Log:    log,
		Router: rest.NewHTTPRouter(app.Service, app.cfg.HTTP.Port, log),
	}
	grpcServer := grpc.NewGRPCServer(app.Service, app.cfg.GRPC.Port, log)

	app.httpServer.Run(ctx, g)
	g.Go(func() error {
		return grpcServer.Run(ctx)
	})

	log.Info("Servers are started")

	<-ctx.Done()
	log.Info("Shutdown signal received, shutting down servers...")

	return g.Wait()
}
