package dao

import (
	"auth-api-go/internal/data"
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

func GetTokenData(ctx context.Context, db *sql.DB, token string) (data.TokenData, error) {
	strSql := "select token_id, user_id, token, added_dt, expired_ind, expiry_dt from token where token = ?"
	rows, err := db.QueryContext(ctx, strSql, token)
	if err != nil {
		return data.TokenData{}, err
	}

	defer closeRows(ctx, rows)

	res := data.TokenData{}
	if rows.Next() {
		var expired string

		err = rows.Scan(&res.Id, &res.UserId, &res.Token, &res.Added, &expired, &res.ExpiryDt)
		res.Expired = getBool(expired)

		if err != nil {
			return data.TokenData{}, err
		}
	} else {
		return res, errors.New("cannot find the token you mentioned")
	}

	return res, nil
}

func GetCheckTokenData(ctx context.Context, db *sql.DB, userId string, token string) (data.LoginResponse, error) {
	strSql := `select 
		u.user_id, u.login, u.name, r.right_cd
		from 
			user u 
		left outer join user_right r on u.user_id = r.user_id
		left outer join token t on u.user_id = t.user_id
		where u.user_id = ? and t.token = ?`
	rows, err := db.QueryContext(ctx, strSql, userId, token)

	if err != nil {
		return data.LoginResponse{}, nil
	}

	defer closeRows(ctx, rows)
	first := true
	var res data.LoginResponse
	rights := make([]string, 0)

	for rows.Next() {
		var rightCd sql.NullString
		var id string
		var login string
		var name string

		err = rows.Scan(&id, &login, &name, &rightCd)

		if err != nil {
			return data.LoginResponse{}, err
		}

		if first {
			first = false

			res = data.LoginResponse{
				Id:     id,
				Login:  login,
				Name:   name,
				Rights: nil,
				Token:  token,
			}
		}

		if rightCd.Valid {
			rights = append(rights, rightCd.String)
		}
	}

	res.Rights = rights

	return res, nil
}

func ExpireTokens(ctx context.Context, db *sql.DB, tokens []string) error {
	sqls := []string{"update token set expired_ind = ? where token_id in ("}
	params := []any{"Y"}
	first := true

	for _, tk := range tokens {
		if first {
			first = false
		} else {
			sqls = append(sqls, ",")
		}

		sqls = append(sqls, "?")
		params = append(params, tk)
	}

	sqls = append(sqls, ")")
	strSql := strings.Join(sqls, "")
	_, err := db.ExecContext(ctx, strSql, params...)
	if err != nil {
		return err
	}

	return nil
}

func GetDistinctUserIdsForTokens(ctx context.Context, db *sql.DB, tokens []string) (map[string]bool, error) {
	sqls := []string{"select distinct user_id from token where token in ("}
	var params []any
	first := true
	for _, tk := range tokens {
		params = append(params, tk)

		if first {
			first = false
		} else {
			sqls = append(sqls, ",")
		}

		sqls = append(sqls, "?")
	}

	sqls = append(sqls, ")")

	strSql := strings.Join(sqls, "")
	rows, err := db.QueryContext(ctx, strSql, params...)
	if err != nil {
		return nil, err
	}

	defer closeRows(ctx, rows)

	mapRes := make(map[string]bool)
	for rows.Next() {
		var id string
		err = rows.Scan(&id)

		if err != nil {
			return nil, err
		}

		mapRes[id] = true
	}

	return mapRes, nil
}

func GetTokenValues(ctx context.Context, db *sql.DB, ids []string) (map[string]bool, error) {
	sqls := []string{"select token_id from token where token_id in ("}
	var params []any
	first := true
	for _, tk := range ids {
		params = append(params, tk)

		if first {
			first = false
		} else {
			sqls = append(sqls, ",")
		}

		sqls = append(sqls, "?")
	}

	sqls = append(sqls, ")")

	strSql := strings.Join(sqls, "")
	rows, err := db.QueryContext(ctx, strSql, params...)
	if err != nil {
		return nil, err
	}
	defer closeRows(ctx, rows)

	mapRes := make(map[string]bool)
	for rows.Next() {
		var tkn string
		err = rows.Scan(&tkn)

		if err != nil {
			return nil, err
		}

		mapRes[tkn] = true
	}

	return mapRes, nil
}

func CancelToken(ctx context.Context, db *sql.DB, token string) error {
	strSql := "update token set expired_ind = 'Y' where token = ?"
	params := []any{token}
	_, err := db.ExecContext(ctx, strSql, params...)

	return err // hopefully nil
}

func SetTokenExpiryDate(ctx context.Context, db *sql.DB, token string, newExpiry time.Time) error {
	strSql := "update token set expiry_dt = ? where token = ?"
	params := []any{newExpiry, token}

	val, err := db.ExecContext(ctx, strSql, params...)
	if err != nil {
		return err
	}

	err = checkNumberRowsChanged(val, 1)

	return err // nil or not
}

func GetTokenByUserId(ctx context.Context, db *sql.DB, id string) ([]data.TokenData, error) {
	strSql := "select token_id, token, added_dt, expired_ind, expiry_dt from token where user_id = ?"
	rows, err := db.QueryContext(ctx, strSql, id)
	if err != nil {
		return nil, err
	}

	defer closeRows(ctx, rows)
	res := make([]data.TokenData, 0)

	for rows.Next() {
		var tokenId string
		var token string
		var addedDt time.Time
		var expiredInd string
		var expiryDt time.Time

		err = rows.Scan(&tokenId, &token, &addedDt, &expiredInd, &expiryDt)
		if err != nil {
			return nil, err
		}

		tokenData := data.TokenData{
			Id:       tokenId,
			UserId:   id,
			Token:    token,
			Added:    addedDt,
			Expired:  getBool(expiredInd),
			ExpiryDt: expiryDt,
		}

		res = append(res, tokenData)
	}

	return res, nil
}
