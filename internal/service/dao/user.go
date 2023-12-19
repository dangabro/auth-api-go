package dao

import (
	"auth-api-go/internal/data"
	"auth-api-go/internal/service/process"
	"context"
	"database/sql"
	"fmt"
)

func GetUserById(ctx context.Context, db *sql.DB, id string) (data.CompleteUserData, error) {
	res := data.CompleteUserData{}

	strSql := `select 
		u.user_id, u.name, u.login, u.password, r.right_cd 
			from 
		user u left join user_right r 
			on u.user_id = r.user_id where u.user_id = ?`

	rows, err := db.QueryContext(ctx, strSql, id)
	if err != nil {
		return res, err
	}

	defer closeRows(ctx, rows)

	found := false
	rights := make([]string, 0)

	for rows.Next() {
		var userId string
		var name string
		var login string
		var password sql.NullString

		var rightCd sql.NullString
		if rightCd.Valid {
			rights = append(rights, rightCd.String)
		}

		if !found {
			pass := ""
			if password.Valid {
				pass = password.String
			}

			res = data.CompleteUserData{
				Id:       userId,
				Name:     name,
				Login:    login,
				Password: pass,
				Rights:   nil,
			}
		}

		found = true

		err = rows.Scan(&userId, &name, &login, &password, &rightCd)
		if err != nil {
			return res, err
		}
	}

	// stick the rights to the result
	res.Rights = rights

	if !found {
		return res, fmt.Errorf("cannot find user with the mentioned id %s", id)
	}

	return res, nil
}

func SearchUser(ctx context.Context, db *sql.DB, search string) ([]data.CompleteUserData, error) {
	strSql := `select 
		u.user_id, u.name, u.login, u.password, r.right_cd 
			from 
		user u left join user_right r 
			on u.user_id = r.user_id where u.login like ? or u.name like ? order by login`

	rows, err := db.QueryContext(ctx, strSql, search, search)
	if err != nil {
		return nil, err
	}

	defer closeRows(ctx, rows)

	var userList []*data.CompleteUserData
	userMap := make(map[string]*data.CompleteUserData)

	for rows.Next() {
		var id string
		var name string
		var login string
		var password sql.NullString
		var rightCd sql.NullString

		err = rows.Scan(&id, &name, &login, &password, &rightCd)
		if err != nil {
			return nil, err
		}

		// if id does not exist in the map, then proceed to add it
		completeUserData, ok := userMap[id]
		if !ok {
			completeUserData = &data.CompleteUserData{}
			completeUserData.Rights = make([]string, 0)
			userMap[id] = completeUserData

			userList = append(userList, completeUserData)

			// fill out the values
			completeUserData.Id = id
			completeUserData.Name = name
			completeUserData.Login = login
			if password.Valid {
				completeUserData.Password = password.String
			}
		}

		// if rightCd is good, then add it to the list
		if rightCd.Valid {
			completeUserData.Rights = append(completeUserData.Rights, rightCd.String)
		}
	}

	users := make([]data.CompleteUserData, 0)
	for _, user := range userList {
		users = append(users, *user)
	}
	return users, nil
}

func UpdateUser(ctx context.Context, db *sql.DB, user data.UpdateUserData) error {
	strSql := "update user set name = ?, login = ? where user_id = ?"
	_, err := db.ExecContext(ctx, strSql, user.Name, user.Login, user.Id)

	return err
}

func AddUser(ctx context.Context, db *sql.DB, user data.UpdateUserData) error {
	strSql := "insert into user (user_id, name, login) values (?, ?, ?)"

	id := user.Id
	name := user.Name
	login := user.Login

	_, err := db.ExecContext(ctx, strSql, id, name, login)

	return err
}

// CheckLoginNotDuplicate check the login is not being employed for a user with an id different from the current user id
func CheckLoginNotDuplicate(ctx context.Context, db *sql.DB, userId string, login string) error {
	strSql := "select user_id from user where login = ? and user_id <> ?"
	rows, err := db.QueryContext(ctx, strSql, login, userId)
	if err != nil {
		return err
	}

	defer closeRows(ctx, rows)

	if rows.Next() {
		return fmt.Errorf("there is already another user using the login %s", login)
	}

	return nil
}

func SaveUserRights(ctx context.Context, db *sql.DB, id string, rights []string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer rollback(tx)

	err = deleteUserRights(ctx, tx, id)
	if err != nil {
		return err
	}

	err = insertUserRights(ctx, tx, id, rights)
	if err != nil {
		return err
	}

	// all right, commit
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func deleteUserRights(ctx context.Context, tx *sql.Tx, id string) error {
	strSql := "delete from user_right where user_id = ?"
	_, err := tx.ExecContext(ctx, strSql, id)
	if err != nil {
		return err
	}

	return nil
}

func insertUserRights(ctx context.Context, tx *sql.Tx, id string, rights []string) error {
	strSql := "insert into user_right (user_right_id, user_id, right_cd) values (?, ?, ?)"
	for _, right := range rights {
		_, err := tx.ExecContext(ctx, strSql, process.GetNewUUID(), id, right)

		if err != nil {
			return err
		}
	}

	return nil
}

func rollback(transaction *sql.Tx) {
	_ = transaction.Rollback()
}
