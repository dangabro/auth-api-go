package changePassword

import (
	"auth-api-go/internal/data"
	"auth-api-go/internal/handler"
	"auth-api-go/internal/service/dao"
	"auth-api-go/internal/service/process"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

type changePassword struct {
	config data.Config
	db     *sql.DB
}

func New(config data.Config, db *sql.DB) http.Handler {
	return &changePassword{
		config: config,
		db:     db,
	}
}

func (h *changePassword) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	res, err := h.process(ctx, writer, r)

	if err != nil {
		handler.Error(err, writer, http.StatusBadRequest)
		log.Err(err).Msg("error changePassword")
		return
	}

	log.Info().Msgf("success changePassword")

	_ = handler.WriteJsonResponse(writer, res)
}

func (h *changePassword) process(ctx context.Context, writer http.ResponseWriter, r *http.Request) (any, error) {
	// check the token and ensure at least auth_connect
	loggedIn, err := handler.CheckRequestToken(ctx, h.db, r, []string{data.AUTH_CONNECT}, true)
	if err != nil {
		return nil, err
	}

	// get the user id associated with the token
	userId := loggedIn.Id

	// get and check the change password data
	pwdData, err := h.getAndCheckPasswordData(ctx, r)
	if err != nil {
		return nil, err
	}
	// user must exist
	err = h.confirmUserExists(ctx, pwdData.Id)
	if err != nil {
		return nil, err
	}
	// password should not be empty
	err = handler.CheckPassword(ctx, pwdData.Password)
	if err != nil {
		return nil, err
	}

	// if id different than the logged id, then check for auth_admin
	// to change the pwd for another person, you need admin rights
	if pwdData.Id != userId {
		// check the account admin
		err = loggedIn.CheckEither([]string{data.AUTH_ADMIN})

		if err != nil {
			return nil, err
		}
	}

	hashedPassword := process.HashPassword(pwdData.Password)
	err = dao.ChangeUserPassword(ctx, h.db, pwdData.Id, hashedPassword)
	if err != nil {
		return nil, err
	}

	return data.NewSuccessData(), nil
}

func (h *changePassword) confirmUserExists(ctx context.Context, id string) error {
	_, err := dao.GetUser(ctx, h.db, id)
	if err != nil {
		return err
	}

	return nil
}

func (h *changePassword) getAndCheckPasswordData(_ context.Context, r *http.Request) (data.ChangePasswordData, error) {
	pwdData := data.ChangePasswordData{}
	err := json.NewDecoder(r.Body).Decode(&pwdData)

	if err != nil {
		return pwdData, err
	}

	return pwdData, nil
}
