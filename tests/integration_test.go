package tests

import (
	"context"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/infinity-ocean/ikakbolit/internal/application"
	"github.com/infinity-ocean/ikakbolit/internal/config"
	"github.com/infinity-ocean/ikakbolit/pkg/dbtest"
	"github.com/infinity-ocean/ikakbolit/pkg/tests"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite

	wg sync.WaitGroup

	cfg config.Config

	apiClient tests.APIClient

	db *sqlx.DB
}

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	suite.Run(t, &Suite{})
}

func (s *Suite) SetupSuite() {
	var app *application.App

	var err error
	rq := s.Require()

	ctx := context.Background()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		app = application.New()
		app.Run("v0.0.0")
	}()

	time.Sleep(time.Second)

	s.cfg = app.Cfg

	s.db, err = sqlx.ConnectContext(ctx, "pgx", s.cfg.Postgres.DSN)
	rq.NoError(err)


	s.apiClient = tests.NewAPIClient(
		"http://"+s.cfg.HTTP.ListenAddress+s.cfg.HTTP.Port,
		http.DefaultClient,
	)
}

func (s *Suite) SetupTest() {
	rq := s.Require()

	err := dbtest.MigrateFromFile(s.db, "testdata/cleanup.sql")
	rq.NoError(err)
}

func (s *Suite) TearDownSuite() {
	rq := s.Require()

	p, err := os.FindProcess(os.Getpid())
	rq.NoError(err)
	rq.NoError(p.Signal(os.Interrupt))

	s.wg.Wait()

	rq.NoError(s.db.Close())
}
