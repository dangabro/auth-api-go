package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

func closeRows(ctx context.Context, rows *sql.Rows) {
	if rows != nil {
		err := rows.Close()

		if err != nil {
			log.Println(ctx, err.Error())
		}
	}
}

func getBool(val string) bool {
	trimmed := strings.TrimSpace(val)
	upper := strings.ToUpper(trimmed)

	return upper == "Y" || upper == "1"
}

func checkNumberRowsChanged(val sql.Result, rows int) error {
	lines, err := val.RowsAffected()
	if err != nil {
		return err
	}

	if lines != int64(rows) {
		return fmt.Errorf("expected to change %d record(s) and changed only %d", rows, lines)
	}

	return nil
}

func CheckConnection(ctx context.Context, db *sql.DB) (string, error) {
	strSql := "select current_timestamp"
	var res string
	rows, err := db.QueryContext(ctx, strSql)
	if err != nil {
		return "", err
	}
	defer closeRows(ctx, rows)

	if rows.Next() {
		err = rows.Scan(&res)

		if err != nil {
			return "", err
		}
	} else {
		return "", errors.New("there is no record returned by the database upon SQL health check")
	}

	return res, nil
}
