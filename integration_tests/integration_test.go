package integration_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"

	"avito-tenders/config"
	"avito-tenders/internal/api"
	"avito-tenders/internal/api/tenders/dtos"
	"avito-tenders/pkg/backend"
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
	m, err := migrate.New("file://../migrations", psqlContainer.GetDSN())
	s.Require().NoError(err)

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
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
		log.Fatalf("Failed to initialize backend: %v", err)
	}

	s.back = back

	// insert default data
	queriesLoader := NewQueriesLoader(s.T(), Queries)
	query := queriesLoader.LoadString("queries/default_data.sql")
	_, err = s.back.DB.Exec(query)
	s.Require().NoError(err)

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

// create fixtures before each test
func (s *TestSuite) SetupTest() {
	fixtures, err := testfixtures.New(
		testfixtures.Database(s.back.DB.DB),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("./fixtures/storage"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
}

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestCreateTender() {
	type want struct {
		StatusCode int
	}
	type args struct {
		inputFileName string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Create tender",
			args: args{
				inputFileName: "tenders/new/create_tender.json",
			},
			want: want{
				StatusCode: 200,
			},
		},
		{
			name: "Missing name",
			args: args{
				inputFileName: "tenders/new/missing_name.json",
			},
			want: want{
				StatusCode: 400,
			},
		},
		{
			name: "Missing description",
			args: args{
				inputFileName: "tenders/new/missing_description.json",
			},
			want: want{
				StatusCode: 400,
			},
		},
		{
			name: "Missing service type",
			args: args{
				inputFileName: "tenders/new/missing_service_type.json",
			},
			want: want{
				StatusCode: 400,
			},
		},
		{
			name: "Missing organization",
			args: args{
				inputFileName: "tenders/new/missing_organization_id.json",
			},
			want: want{
				StatusCode: 400,
			},
		},
		{
			name: "Missing creator username",
			args: args{
				inputFileName: "tenders/new/missing_creator_username.json",
			},
			want: want{
				StatusCode: 400,
			},
		},
		{
			name: "Not from organization",
			args: args{
				inputFileName: "tenders/new/not_from_organization.json",
			},
			want: want{
				StatusCode: 403,
			},
		},
		{
			name: "Unknown user",
			args: args{
				inputFileName: "tenders/new/unknown_user.json",
			},
			want: want{
				StatusCode: 401,
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			requestBody := s.loader.LoadString(fmt.Sprintf("%s/%s", fixturesPath, tt.args.inputFileName))

			res, err := s.server.Client().Post(fmt.Sprintf("%s/api/tenders/new", s.server.URL), "", bytes.NewBufferString(requestBody))
			s.Require().NoError(err)

			defer res.Body.Close()

			s.Require().Equal(tt.want.StatusCode, res.StatusCode)
			if tt.want.StatusCode != 200 {
				return
			}

			// check response
			var response dtos.TenderResponse
			err = json.NewDecoder(res.Body).Decode(&response)
			s.Require().NoError(err)

			expected := s.loader.LoadTemplate(fmt.Sprintf("%s/%s.result", fixturesPath, tt.args.inputFileName), map[string]interface{}{
				"id":        response.Id,
				"createdAt": response.CreatedAt,
			})

			JSONEq(s.T(), expected, response)
		})
	}
}

// func (s *TestSuite) TestGetUser() {
// 	res, err := s.server.Client().Get(s.server.URL + "/users/1")
// 	s.Require().NoError(err)
//
// 	defer res.Body.Close()
//
// 	s.Require().Equal(http.StatusOK, res.StatusCode)
//
// 	// check response
// 	expected := s.loader.LoadString("fixtures/api/get_user_response.json")
//
// 	JSONEq(s.T(), expected, res.Body)
// }
//
// func (s *TestSuite) TestDepositBalance() {
// 	// mock http call
// 	httpmock.RegisterResponder(
// 		http.MethodPost,
// 		addr+"/deposit",
// 		httpmock.NewStringResponder(http.StatusOK, ""),
// 	)
// 	//
//
// 	requestBody := s.loader.LoadString("fixtures/api/deposit_user_request.json")
//
// 	res, err := s.server.Client().Post(s.server.URL+"/users/deposit", "", bytes.NewBufferString(requestBody))
// 	s.Require().NoError(err)
//
// 	defer res.Body.Close()
//
// 	s.Require().Equal(http.StatusOK, res.StatusCode)
//
// 	// check response
// 	expected := s.loader.LoadString("fixtures/api/deposit_user_response.json")
//
// 	JSONEq(s.T(), expected, res.Body)
// }
