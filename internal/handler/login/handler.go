package login

import (
	"auth-api-go/internal/data"
	"auth-api-go/internal/handler"
	"auth-api-go/internal/service/business"
	"context"
	"database/sql"
	"github.com/rs/zerolog/log"
	"net/http"
)

type loginHandler struct {
	config data.Config
	db     *sql.DB
}

func New(config data.Config, db *sql.DB) http.Handler {
	return &loginHandler{
		config: config,
		db:     db,
	}
}

func (h *loginHandler) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	log.Info().Msg("login")
	ctx := context.Background()
	ctx = context.WithValue(ctx, "method", "login")

	// read login data from the request
	loginData, err := handler.ReadRequestLoginData(r)
	if err != nil {
		log.Error().Msg(err.Error())
		handler.Error(err, writer, http.StatusBadRequest)
		return
	}

	tokenDurationMs := h.config.TokenDurationMs
	loginResponse, err := business.Login(ctx, h.db, tokenDurationMs, loginData)

	if err != nil {
		log.Error().Msg(err.Error())
		handler.Error(err, writer, http.StatusBadRequest)
		return
	}

	log.Info().Msgf("successful logged in: %s", loginData.Login)

	_ = handler.WriteJsonResponse(writer, loginResponse)
}
