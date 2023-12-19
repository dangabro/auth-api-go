package dao

import (
	"context"
	"database/sql"
)

func ChangeUserPassword(ctx context.Context, db *sql.DB, id string, password string) error {
	strSql := "update user set password = ? where user_id = ?"
	params := []any{password, id}

	_, err := db.ExecContext(ctx, strSql, params...)

	return err
}
