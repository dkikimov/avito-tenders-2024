package app

import (
	"errors"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"avito-tenders/pkg/backend"
)

const (
	defaultAttempts = 20
	defaultTimeout  = time.Second
)

func Migrate(backend backend.Backend) {
	var (
		attempts = defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	driver, err := postgres.WithInstance(backend.DB.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not create driver to database: %s", err)
	}

	for attempts > 0 {
		m, err = migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
		if err == nil {
			break
		}

		log.Printf("Migrate: could not connect to database: %s, attemps lefr: %d", err, attempts)
		time.Sleep(defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("Migrate: postgres connect error: %s", err)
	}

	if m == nil {
		log.Fatal("Migrate is nil")
	}

	_, _, err = m.Version()
	if err != nil {
		// Skip first migration
		if errors.Is(err, migrate.ErrNilVersion) {
			log.Printf("Migrate: no migration found. Set force version")
			if err := m.Force(20240909113647); err != nil {
				log.Fatalf("Migrate: could not force migration: %s", err)
			}
		}
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrate: no change")
		return
	}

	log.Printf("Migrate: up success")
}
