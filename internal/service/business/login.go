package business

import (
	"auth-api-go/internal/data"
	"auth-api-go/internal/service/dao"
	"auth-api-go/internal/service/process"
	"context"
	"database/sql"
	"errors"
)

func Login(ctx context.Context, db *sql.DB, tokenDurationMillisec int64, loginData data.LoginData) (data.LoginResponse, error) {
	// check login is filled out
	err := checkLoginData(loginData)
	if err != nil {
		return data.LoginResponse{}, err
	}

	// get the user and check it does exist
	userData, err := dao.GetUserDataByLogin(ctx, db, loginData.Login)
	if err != nil {
		return data.LoginResponse{}, err
	}

	// check the password
	err = checkPassword(loginData.Password, userData.Password)
	if err != nil {
		return data.LoginResponse{}, err
	}

	// if OK, create token entry and proceed
	token, err := dao.CreateTokenEntry(ctx, db, tokenDurationMillisec, userData.Id)
	if err != nil {
		return data.LoginResponse{}, err
	}

	return data.LoginResponse{
		Id:     userData.Id,
		Login:  userData.Login,
		Name:   userData.Name,
		Rights: userData.Rights,
		Token:  token,
	}, nil
}

func checkPassword(providedPassword string, hashedPassword string) error {

	hashedHexa := process.HashPassword(providedPassword)

	if hashedHexa != hashedPassword {
		return errors.New("Invalid credentials")
	}

	return nil
}

func checkLoginData(loginData data.LoginData) error {
	// TODO
	return nil
}
