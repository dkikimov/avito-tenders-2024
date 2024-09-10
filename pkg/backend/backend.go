package backend

import (
	"fmt"

	"avito-tenders/config"

	"github.com/jmoiron/sqlx"
)

// Backend contains application connections to different external services and additional parameters that should be
// passed to API middlewares.
type Backend struct {
	DB *sqlx.DB
}

func NewForServer(cfg *config.Config) (Backend, error) {
	dbConn, err := newDBConnection(&dbConnectionOpts{
		Host:                       cfg.PostgresHost,
		Port:                       cfg.PostgresPort,
		Name:                       cfg.PostgresDatabase,
		Username:                   cfg.PostgresUsername,
		Password:                   cfg.PostgresPassword,
		NewConnectionRetryInterval: 1,
		NewConnectionRetryAttempts: 3,
	})
	if err != nil {
		return Backend{}, fmt.Errorf("unable to setup DB connection for the new backend: %w", err)
	}

	return Backend{
		DB: dbConn,
	}, nil
}
