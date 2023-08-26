package repository

import "database/sql"

// nullString converts a string to sql.NullString for empty string check
func nullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}

	return sql.NullString{
		String: value,
		Valid:  true,
	}
}
