package postgres

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	// for the sake of pgx compatibility.
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	defaultMaxOpenConns    = 30
	defaultMaxIdleConns    = 15
	defaultMaxConnLifetime = 180
)

// Opts represents options to initializes new postgresql wrapper.
type Opts struct {
	Host     string
	Port     int
	Name     string
	Username string
	Password string
}

// New initializes a new postgresql wrapper and verifies that connection is stable.
func New(opts *Opts) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s timezone=UTC",
		opts.Host,
		opts.Port,
		opts.Name,
		opts.Username,
		opts.Password,
	)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to init db connection: %w", err)
	}

	db.SetMaxOpenConns(defaultMaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(defaultMaxConnLifetime) * time.Second)
	db.SetMaxIdleConns(defaultMaxIdleConns)

	return db, nil
}
