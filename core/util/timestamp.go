package util

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ParseTimestamp(timestamp string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.000-0700", timestamp)
}

func FormatTimestamp(timestamp time.Time) string {
	return timestamp.Format("2006-01-02T15:04:05.000-0700")
}

func FormatTimestampz(timestamp pgtype.Timestamptz) *string {
	if timestamp.Valid {
		result := FormatTimestamp(timestamp.Time)
		return &result
	}
	return nil
}
