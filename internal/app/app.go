package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	"avito-tenders/config"
	"avito-tenders/internal/api"
	"avito-tenders/pkg/backend"
	"avito-tenders/pkg/httpserver"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	back, err := backend.NewForServer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize backend: %v", err)
	}

	Migrate(back)

	routes, err := api.InitAPIRoutes(back)
	if err != nil {
		log.Fatalf("Failed to initialize API routes: %v", err)
	}

	server := httpserver.New(routes, httpserver.Address(cfg.ServerAddress))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Printf("Received signal \"%v\", shutting down", s)
	case err = <-server.Notify():
		log.Printf("Received error, shutting down: %s", err)
	}

	// Shutdown
	err = server.Shutdown()
	if err != nil {
		log.Printf("app - Run - httpServer.Shutdown: %s", err)
	}
}
