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
	grpcServer   *grpc.GRPCServer
	Service      *service.Service
	Repo         *repository.Repo
	postgresPool *pgxpool.Pool
}

func New(appVersion string) App { //nolint:funlen
	const appName = "ikakbolit"

	cfg := lo.Must(config.Load())
	log := &connectors.Slog{Name: appName, Version: appVersion, Debug: cfg.Debug}

	pool := lo.Must(repository.MakePool(cfg.Postgres.DSN))

	repo := repository.New(pool)

	return App{
		cfg: cfg,
		log: log,
		httpServer: &modules.HTTPServer{
			Port: cfg.HTTP.Port,
		},
		grpcServer: &grpc.GRPCServer{
			Port: cfg.GRPC.Port,
		},
		Repo:         repo,
		postgresPool: pool,
	}
}

func (app App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	defer stop()

	log := app.log.Logger(ctx)

	log.Info("Program is starting...")

	g, ctx := errgroup.WithContext(ctx)

	app.httpServer.Log = log
	app.httpServer.Router = rest.NewHTTPRouter(app.Service, app.cfg.HTTP.Port, log)
	grpc := grpc.NewGRPCServer(app.Service, app.cfg.GRPC.Port, log)

	app.httpServer.Run(ctx, g)
	modules.RunGRPC(grpc, app.cfg, log)
	log.Info("Servers are started")

	<-ctx.Done()
	log.Info("Shutdown signal received, shutting down servers...")
	return nil
}
