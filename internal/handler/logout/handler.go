package logout

import (
	"auth-api-go/internal/data"
	"auth-api-go/internal/handler"
	"auth-api-go/internal/service/dao"
	"context"
	"database/sql"
	"errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

type logout struct {
	config data.Config
	db     *sql.DB
}

func New(config data.Config, db *sql.DB) http.Handler {
	return &logout{
		config: config,
		db:     db,
	}
}

func (h *logout) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	res, err := h.process(ctx, writer, r)

	if err != nil {
		handler.Error(err, writer, http.StatusBadRequest)
		log.Err(err).Msg("error logout")

		return
	}

	log.Info().Msg("success logout")

	_ = handler.WriteJsonResponse(writer, res)
}

func (h *logout) process(ctx context.Context, writer http.ResponseWriter, r *http.Request) (any, error) {
	// get the token
	token := r.Header.Get("Authorization")
	token = strings.TrimSpace(token)
	if len(token) == 0 {
		return nil, errors.New("there is no authorization header")
	}

	if strings.HasPrefix(token, "bearer") {
		token = token[7:]
	}

	// see if the token exists
	_, err := dao.GetTokenData(ctx, h.db, token)
	if err != nil {
		return nil, err
	}

	// if exists, just cancel it - no on shall have the token other than the person for whom it was created
	err = dao.CancelToken(ctx, h.db, token)
	if err != nil {
		return nil, err
	}

	// return success
	return data.NewSuccessData(), nil
}
