package tests

import (
	"context"
	"fmt"

	"net/http/httptest"
	"sync"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/samber/lo"

	"github.com/infinity-ocean/ikakbolit/internal/application"
	"github.com/infinity-ocean/ikakbolit/internal/config"
	"github.com/infinity-ocean/ikakbolit/internal/domain/service"
	"github.com/infinity-ocean/ikakbolit/internal/infrastructure/repository"
	"github.com/infinity-ocean/ikakbolit/internal/server/rest"
	"github.com/infinity-ocean/ikakbolit/pkg/application/connectors"
	"github.com/infinity-ocean/ikakbolit/pkg/dbtest"
	"github.com/infinity-ocean/ikakbolit/pkg/tests"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	wg        sync.WaitGroup
	cfg       config.Config
	apiClient tests.APIClient
	ts        *httptest.Server
	db        *sqlx.DB
}

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	suite.Run(t, &Suite{})
}

func (s *Suite) SetupSuite() {
	godotenv.Load("../.env")

	rq := s.Require()

	ctx := context.Background()

	app := application.New()

	s.cfg = app.Cfg
	fmt.Println("check dsn")
	fmt.Println(s.cfg.Postgres.DSN)

	var err error
	s.db, err = sqlx.ConnectContext(ctx, "pgx", s.cfg.Postgres.DSN)
	rq.NoError(err)

	pool := lo.Must(repository.MakePool(app.Cfg.Postgres.DSN))
	repo := repository.New(pool)

	logger := connectors.Slog{Debug: true}
	slog := logger.Logger(context.Background())

	svc := service.New(repo, slog, app.Cfg)

	router := rest.NewHTTPRouter(svc, app.Cfg.HTTP.Port, slog)
	s.ts = httptest.NewServer(router.GetRouter())

	s.apiClient = tests.NewAPIClient(
		s.ts.URL,
		s.ts.Client(),
	)
}

func (s *Suite) SetupTest() {
	rq := s.Require()

	err := dbtest.MigrateFromFile(s.db, "testdata/cleanup.sql")
	rq.NoError(err)
}

func (s *Suite) TearDownSuite() {
	rq := s.Require()

	s.wg.Wait()
	s.ts.Close()
	rq.NoError(s.db.Close())
}
