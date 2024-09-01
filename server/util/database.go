package util

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/nineteenseventy/minichat/server/util/database"
)

var globalDatabase *pgx.Conn

func InitDatabase(ctx context.Context, config database.DatabaseConfig) error {
	logger := GetLogger("database")

	configStruct, err := database.ParseConfig(config)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to parse database config")
		return err
	}

	configStruct.AfterConnect = database.ValidateConnect(logger)
	configStruct.OnPgError = database.PgError(logger)
	configStruct.OnNotice = database.Notice(logger)
	configStruct.Tracer = database.NewDatabaseTracer(logger)
	globalDatabase, err = pgx.ConnectConfig(ctx, configStruct)
	return err
}

func GetDatabase() *pgx.Conn {
	if globalDatabase == nil {
		panic("Database not initialized")
	}
	return globalDatabase
}
