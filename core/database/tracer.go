package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

type DatabaseTracer struct {
	logger zerolog.Logger
}

func argsToString(args []interface{}) string {
	var strArgs []string
	for _, arg := range args {
		strArgs = append(strArgs, fmt.Sprintf("%v", arg))
	}
	return strings.Join(strArgs, ", ")
}

func cleanupSql(sql string) string {
	rows := strings.Split(sql, "\n")
	var cleanedRows []string
	for _, row := range rows {
		cleanedRows = append(cleanedRows, strings.TrimSpace(row))
	}
	return strings.Join(cleanedRows, "\\n")
}

func (tracer *DatabaseTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	tracer.logger.Debug().
		Str("address", conn.PgConn().Conn().RemoteAddr().String()).
		Str("args", argsToString(data.Args)).
		Msg(cleanupSql(data.SQL))
	return ctx
}

func (tracer *DatabaseTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
}

func NewDatabaseTracer(logger zerolog.Logger) *DatabaseTracer {
	return &DatabaseTracer{logger: logger}
}
