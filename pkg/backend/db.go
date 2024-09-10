package backend

import (
	"fmt"
	"log/slog"
	"time"

	"avito-tenders/pkg/postgres"

	"github.com/jmoiron/sqlx"
)

// dbConnectionOpts represents options for a new DB connection.
type dbConnectionOpts struct {
	Host                       string
	Port                       int
	Name                       string
	Username                   string
	Password                   string
	NewConnectionRetryInterval int
	NewConnectionRetryAttempts int
}

// newDBConnection wraps DB connection setup.
func newDBConnection(opts *dbConnectionOpts) (*sqlx.DB, error) {
	connOpts := &postgres.Opts{
		Host:     opts.Host,
		Port:     opts.Port,
		Name:     opts.Name,
		Username: opts.Username,
		Password: opts.Password,
	}

	var (
		conn *sqlx.DB
		err  error
	)

	newConnectionRetryInterval := time.Second * time.Duration(opts.NewConnectionRetryInterval)

	for i := 0; i < opts.NewConnectionRetryAttempts; i++ {
		conn, err = postgres.New(connOpts)
		if err == nil {
			slog.Info(
				"connected to DB",
				"host", opts.Host,
				"port", opts.Port,
			)

			return conn, nil
		}
		slog.Warn(
			"unable to create a new DB connection",
			"error", err,
			"attempt", i,
			"attempt interval", newConnectionRetryInterval,
		)

		time.Sleep(newConnectionRetryInterval)
	}

	return nil, fmt.Errorf(
		"unable to open DB connection, after %d retries last error was: %w",
		opts.NewConnectionRetryAttempts,
		err,
	)
}
