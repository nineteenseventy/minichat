package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nineteenseventy/minichat/core/logging"
)

var globalDatabase *pgxpool.Pool

func InitDatabase(ctx context.Context, config DatabaseConfig) error {
	logger := logging.GetLogger("database")

	configStruct, err := ParseConfig(config)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to parse database config")
		return err
	}

	configStruct.ConnConfig.AfterConnect = ValidateConnect(logger)
	configStruct.ConnConfig.OnPgError = PgError(logger)
	configStruct.ConnConfig.OnNotice = Notice(logger)
	configStruct.ConnConfig.Tracer = NewDatabaseTracer(logger)
	globalDatabase, err = pgxpool.NewWithConfig(ctx, configStruct)
	return err
}

func GetDatabase() *pgxpool.Pool {
	if globalDatabase == nil {
		panic("Database not initialized")
	}
	return globalDatabase
}
