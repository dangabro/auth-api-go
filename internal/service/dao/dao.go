package dao

import (
	"auth-api-go/internal/data"
	"auth-api-go/internal/service/process"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func GetUser(ctx context.Context, db *sql.DB, id string) (data.UserData, error) {
	strSql := "select user_id, name, login, password from user where user_id = ?"
	rows, err := db.QueryContext(ctx, strSql, id)
	if err != nil {
		return data.UserData{}, err
	}

	defer closeRows(ctx, rows)

	var userId string
	var name string
	var login string
	var password sql.NullString

	if rows.Next() {
		err = rows.Scan(&userId, &name, &login, &password)

		if err != nil {
			return data.UserData{}, err
		}

		pwd := ""
		if password.Valid {
			pwd = password.String
		}

		return data.UserData{
			Id:       userId,
			Name:     name,
			Login:    login,
			Password: pwd,
		}, nil
	}

	return data.UserData{}, fmt.Errorf("cannot find user with the id: %s", id)
}

func GetUserDataByLogin(ctx context.Context, db *sql.DB, login string) (data.CompleteUserData, error) {
	strSql := "select user_id, name, login, password from user where login = ?"
	rows, err := db.QueryContext(ctx, strSql, login)
	if err != nil {
		return data.CompleteUserData{}, err
	}
	defer closeRows(ctx, rows)

	if rows.Next() {
		var res data.CompleteUserData
		var pwd sql.NullString

		err = rows.Scan(&res.Id, &res.Name, &res.Login, &pwd)
		if err != nil {
			return data.CompleteUserData{}, err
		}

		if pwd.Valid {
			res.Password = pwd.String
		}

		// all good, get the rights
		rightsSql := "select right_cd from user_right where user_id = ?"
		rowsRights, err1 := db.QueryContext(ctx, rightsSql, res.Id)
		if err1 != nil {
			return data.CompleteUserData{}, err1
		}

		defer closeRows(ctx, rowsRights)

		// initialization of empty string slice with literal is not well regarded and I don't want to return nil as the rights
		// but an empty JSON array
		var rights []string = make([]string, 0)

		for rowsRights.Next() {
			var strRight string
			err2 := rowsRights.Scan(&strRight)
			if err2 != nil {
				return data.CompleteUserData{}, err2
			}

			rights = append(rights, strRight)
		}

		res.Rights = rights

		return res, nil
	}

	errorMessage := fmt.Sprintf("cannot find user with login: %s", login)
	return data.CompleteUserData{}, errors.New(errorMessage)
}

func CreateTokenEntry(ctx context.Context, db *sql.DB, tokenDurationMillisecond int64, userId string) (string, error) {
	token := process.GetNewUUID()
	tokenId := process.GetNewUUID()
	expiredInd := "N"
	currentDate := time.Now()
	expiryDate := currentDate.Add(time.Duration(tokenDurationMillisecond) * time.Millisecond)

	strSql := "insert into token (token_id, user_id, token, added_dt, expired_ind, expiry_dt) values (?, ?, ?, ?, ?, ?)"

	_, err := db.ExecContext(ctx, strSql, tokenId, userId, token, currentDate, expiredInd, expiryDate)
	if err != nil {
		return "", err
	}

	return token, nil
}
