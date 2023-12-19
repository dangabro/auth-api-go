package updateUser

import (
	"auth-api-go/internal/data"
	"auth-api-go/internal/handler"
	"auth-api-go/internal/service/dao"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"net/http"
)

type updateUser struct {
	config data.Config
	db     *sql.DB
}

func New(config data.Config, db *sql.DB) http.Handler {
	return &updateUser{
		config: config,
		db:     db,
	}
}

func (h *updateUser) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	res, err := h.process(ctx, writer, r)

	if err != nil {
		handler.Error(err, writer, http.StatusBadRequest)
		log.Err(err).Msgf("error updateUser")

		return
	}

	log.Info().Msgf("update user, success")

	_ = handler.WriteJsonResponse(writer, res)
}

func (h *updateUser) process(ctx context.Context, writer http.ResponseWriter, r *http.Request) (any, error) {
	userLogin, err := handler.CheckRequestToken(ctx, h.db, r, []string{data.AUTH_CONNECT}, true)
	if err != nil {
		return nil, err
	}

	loginId := userLogin.Id

	// get the payload
	var updateUserData data.UpdateUserData
	err = json.NewDecoder(r.Body).Decode(&updateUserData)
	if err != nil {
		return nil, err
	}

	// check user is all right
	err = data.CheckUserData(updateUserData)
	if err != nil {
		return nil, err
	}

	// well, now check the login is not a duplicate situation
	// get the user id and double check the user exists
	updateUserId := updateUserData.Id
	_, err = dao.GetUserById(ctx, h.db, updateUserId)
	adding := false
	if err != nil {
		adding = true
	}
	if adding || loginId != updateUserId {
		// different user, let's see if I am admin
		err = userLogin.CheckEither([]string{data.AUTH_ADMIN})

		if err != nil {
			return nil, errors.New("to update other user or to add user admin rights are required")
		}
	}

	// now check the login id
	err = dao.CheckLoginNotDuplicate(ctx, h.db, updateUserId, updateUserData.Login)
	if err != nil {
		return nil, err
	}

	// all good, update
	if adding {
		err = dao.AddUser(ctx, h.db, updateUserData)

		if err != nil {
			return nil, err
		}
	} else {
		err = dao.UpdateUser(ctx, h.db, updateUserData)

		if err != nil {
			return nil, err
		}
	}

	return data.NewSuccessData(), nil
}
