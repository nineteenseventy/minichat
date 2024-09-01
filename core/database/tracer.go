package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

type DatabaseTracer struct {
	logger zerolog.Logger
}

func (tracer *DatabaseTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	tracer.logger.Debug().
		Str("address", conn.PgConn().Conn().LocalAddr().String()).
		Msg(data.SQL)
	return ctx
}

func (tracer *DatabaseTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
}

func NewDatabaseTracer(logger zerolog.Logger) *DatabaseTracer {
	return &DatabaseTracer{logger: logger}
}
