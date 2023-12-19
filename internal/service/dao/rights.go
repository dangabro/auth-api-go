package dao

import (
	"auth-api-go/internal/data"
	"context"
	"database/sql"
)

func GetRightDataList(ctx context.Context, db *sql.DB) ([]data.RightData, error) {
	sql := "select right_cd, name from `right` order by right_cd"
	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}

	defer closeRows(ctx, rows)

	var res = make([]data.RightData, 0)
	for rows.Next() {
		var cd string
		var name string

		err = rows.Scan(&cd, &name)
		if err != nil {
			return nil, err
		}

		res = append(res, data.RightData{
			Cd:   cd,
			Name: name,
		})
	}

	return res, nil
}
