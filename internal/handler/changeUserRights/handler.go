package changeUserRights

import (
	"auth-api-go/internal/data"
	"auth-api-go/internal/handler"
	"auth-api-go/internal/service/dao"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

type changeUserRights struct {
	config data.Config
	db     *sql.DB
}

func New(config data.Config, db *sql.DB) http.Handler {
	return &changeUserRights{
		config: config,
		db:     db,
	}
}

func (h *changeUserRights) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	res, err := h.process(ctx, writer, r)

	if err != nil {
		handler.Error(err, writer, http.StatusBadRequest)
		log.Err(err).Msg("error changeUserRights")

		return
	}

	log.Info().Msg("success changeUserRights")

	_ = handler.WriteJsonResponse(writer, res)
}

func (h *changeUserRights) process(ctx context.Context, writer http.ResponseWriter, r *http.Request) (any, error) {
	// check token that has connect and admin both
	userLogin, err := handler.CheckRequestToken(ctx, h.db, r, []string{data.AUTH_CONNECT, data.AUTH_ADMIN}, false)
	if err != nil {
		return nil, err
	}

	// parse the payload
	var userRights data.ChangeUserRights
	err = json.NewDecoder(r.Body).Decode(&userRights)
	if err != nil {
		return nil, err
	}

	// check id exists
	_, err = dao.GetUserById(ctx, h.db, userRights.Id)
	if err != nil {
		return nil, err
	}

	// check rights mentioned exist in the database
	allRights, err := dao.GetRightDataList(ctx, h.db)
	if err != nil {
		return nil, err
	}

	err = data.CheckValidRights(userRights.Rights, allRights)
	if err != nil {
		return nil, err
	}

	// if id is the logged in id, check it does not remove auth_connect and auth_admin
	if userRights.Id == userLogin.Id {
		err = data.CheckAuthFullConnect(userRights.Rights)
		if err != nil {
			return nil, err
		}
	}

	// operate the rights in the database
	err = dao.SaveUserRights(ctx, h.db, userRights.Id, userRights.Rights)
	if err != nil {
		return nil, err
	}

	return data.NewSuccessData(), nil
}
