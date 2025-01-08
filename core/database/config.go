package database

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseConfig struct {
	Host     string
	Port     uint16
	Database string
	User     string
	Password string
	Tls      bool
}

func ParseConfig(config DatabaseConfig) (*pgxpool.Config, error) {
	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	pgConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	if !config.Tls {
		pgConfig.ConnConfig.TLSConfig = nil
	}

	return pgConfig, nil
}
