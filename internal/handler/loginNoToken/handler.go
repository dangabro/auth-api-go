package loginNoToken

import (
	"auth-api-go/internal/data"
	"auth-api-go/internal/handler"
	"auth-api-go/internal/service/dao"
	"auth-api-go/internal/service/process"
	"context"
	"database/sql"
	"errors"
	"github.com/rs/zerolog/log"
	"net/http"
)

type loginNt struct {
	config data.Config
	db     *sql.DB
}

func New(config data.Config, db *sql.DB) http.Handler {
	return &loginNt{
		config: config,
		db:     db,
	}
}

func (h *loginNt) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	res, err := h.process(ctx, writer, r)

	if err != nil {
		handler.Error(err, writer, http.StatusBadRequest)
		log.Err(err).Msg("error loginNoToken")

		return
	}

	log.Info().Msg("success loginNoToken")

	_ = handler.WriteJsonResponse(writer, res)
}

func (h *loginNt) process(ctx context.Context, writer http.ResponseWriter, r *http.Request) (any, error) {
	// take login and password from the request
	loginData, err := handler.ReadRequestLoginData(r)
	if err != nil {
		return nil, err
	}

	// check values; password is not accepted to be empty string under any circumstance
	login := loginData.Login
	password := loginData.Password
	err = handler.CheckPassword(ctx, password)
	if err != nil {
		return nil, err
	}

	// hash the password
	hashedPassword := process.HashPassword(password)

	// load the complete data for the user
	userData, err := dao.GetUserDataByLogin(ctx, h.db, login)
	if err != nil {
		return nil, err
	}

	// compare passwords
	if hashedPassword != userData.Password {
		return nil, errors.New("Invalid credentials")
	}

	// return the built value object
	return userData, nil
}
