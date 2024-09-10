package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config is the structure that represents configuration for webserver.
type Config struct {
	ServerAddress    string `env:"SERVER_ADDRESS,required"`
	PostgresConn     string `env:"POSTGRES_CONN,required"`
	PostgresUsername string `env:"POSTGRES_USERNAME,required"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	PostgresHost     string `env:"POSTGRES_HOST,required"`
	PostgresPort     int    `env:"POSTGRES_PORT,required"`
	PostgresDatabase string `env:"POSTGRES_DATABASE,required"`
}

func NewConfig() (*Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
