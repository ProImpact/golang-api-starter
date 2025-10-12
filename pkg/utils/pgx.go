package utils

import "github.com/jackc/pgx/v5/pgtype"

func NilOrPgText(value *string) pgtype.Text {
	if value != nil {
		return pgtype.Text{
			String: *value,
			Valid:  true,
		}
	}
	return pgtype.Text{
		Valid: false,
	}
}
