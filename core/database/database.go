package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/nineteenseventy/minichat/core/logging"
)

var globalDatabase *pgx.Conn

func InitDatabase(ctx context.Context, config DatabaseConfig) error {
	logger := logging.GetLogger("database")

	configStruct, err := ParseConfig(config)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to parse database config")
		return err
	}

	configStruct.AfterConnect = ValidateConnect(logger)
	configStruct.OnPgError = PgError(logger)
	configStruct.OnNotice = Notice(logger)
	configStruct.Tracer = NewDatabaseTracer(logger)
	globalDatabase, err = pgx.ConnectConfig(ctx, configStruct)
	return err
}

func GetDatabase() *pgx.Conn {
	if globalDatabase == nil {
		panic("Database not initialized")
	}
	return globalDatabase
}
