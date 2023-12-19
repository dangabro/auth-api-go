package userTokens

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

type userTokens struct {
	config data.Config
	db     *sql.DB
}

func New(config data.Config, db *sql.DB) http.Handler {
	return &userTokens{
		config: config,
		db:     db,
	}
}

func (h *userTokens) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	res, err := h.process(ctx, writer, r)

	if err != nil {
		handler.Error(err, writer, http.StatusBadRequest)
		log.Err(err).Msg("error userTokens")

		return
	}

	log.Info().Msg("success userTokens")

	_ = handler.WriteJsonResponse(writer, res)
}

func (h *userTokens) process(ctx context.Context, writer http.ResponseWriter, r *http.Request) (any, error) {
	// check token and auth_conect
	// check the token and ensure at least auth_connect
	loggedIn, err := handler.CheckRequestToken(ctx, h.db, r, []string{data.AUTH_CONNECT}, true)
	if err != nil {
		return nil, err
	}

	// get the user id associated with the token
	userId := loggedIn.Id

	// get id from the payload
	var id data.IdHolder

	// if id different than user id, then check for admin
	json.NewDecoder(r.Body).Decode(&id)
	payloadId := id.Id
	if payloadId != userId {
		err := loggedIn.CheckEither([]string{data.AUTH_ADMIN})
		if err != nil {
			return nil, err
		}
	}

	// get the tokens from the database and return
	tokens, err := dao.GetTokenByUserId(ctx, h.db, payloadId)
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("userTokens, id: %s", id.Id)
	return data.TokenResponse{Tokens: tokens}, nil
}
