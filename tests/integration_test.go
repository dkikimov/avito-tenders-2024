package tests

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"

	"avito-tenders/config"
	"avito-tenders/internal/api"
	"avito-tenders/pkg/backend"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	fixturesPath = "fixtures/api"
)

type TestSuite struct {
	suite.Suite
	psqlContainer *PostgreSQLContainer
	server        *httptest.Server
	loader        *FixtureLoader
	back          backend.Backend
}

func (s *TestSuite) SetupSuite() {
	// create db container
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	psqlContainer, err := NewPostgreSQLContainer(ctx)
	s.Require().NoError(err)

	s.psqlContainer = psqlContainer

	// run migrations
	var m *migrate.Migrate
	for i := 0; i < 20; i++ {
		m, err = migrate.New("file://../migrations", psqlContainer.GetDSN())
		if err == nil {
			break
		}

		time.Sleep(1 * time.Second)
	}
	if m == nil {
		s.T().FailNow()
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Panicf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrate: no change")
		return
	}

	log.Printf("Migrate: up success")

	// mock client
	mockClient := &http.Client{}
	httpmock.ActivateNonDefault(mockClient)

	// init backend
	back, err := backend.NewForServer(&config.Config{
		ServerAddress:    "0.0.0.0:8080",
		PostgresConn:     psqlContainer.GetDSN(),
		PostgresUsername: psqlContainer.Config.User,
		PostgresPassword: psqlContainer.Config.Password,
		PostgresHost:     psqlContainer.Config.Host,
		PostgresPort:     psqlContainer.Config.MappedPort,
		PostgresDatabase: psqlContainer.Config.Database,
	})
	if err != nil {
		log.Panicf("Failed to initialize backend: %v", err)
	}

	s.back = back

	// insert default data
	queriesLoader := NewQueriesLoader(s.T(), Queries)
	query := queriesLoader.LoadString("queries/default_data.sql")
	_, err = s.back.DB.Exec(query)
	s.Require().NoError(err)
	log.Printf("Inserted default data")

	// init routes
	routes, err := api.InitAPIRoutes(back)
	s.Require().NoError(err)

	// use httptest
	s.server = httptest.NewServer(routes)

	// create fixture loader
	s.loader = NewFixtureLoader(s.T(), Fixtures)
}

func (s *TestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.server.Close()

	httpmock.DeactivateAndReset()

	s.Require().NoError(s.psqlContainer.Terminate(ctx))
	s.Require().NoError(s.back.DB.Close())
}

// create fixtures before each test.
func (s *TestSuite) SetupTest() {
	fixtures, err := testfixtures.New(
		testfixtures.Database(s.back.DB.DB),
		testfixtures.Dialect("postgres"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
}

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
