package util

import "database/sql"

func ParseSqlString(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}
