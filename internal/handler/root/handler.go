package root

import (
	"auth-api-go/internal/data"
	"auth-api-go/internal/handler"
	"context"
	"database/sql"
	"github.com/rs/zerolog/log"
	"net/http"
)

type root struct {
	config data.Config
	db     *sql.DB
}

func New(config data.Config, db *sql.DB) http.Handler {
	return &root{
		config: config,
		db:     db,
	}
}

func (h *root) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	res, err := h.process(ctx, writer, r)

	log.Info().Msg("root invoked with url: " + r.URL.Path)

	if err != nil {
		handler.Error(err, writer, http.StatusBadRequest)
		log.Err(err).Msg("error rootHandler")

		return
	}

	log.Info().Msg("success root handler")

	_ = handler.WriteJsonResponse(writer, res)
}

func (h *root) process(ctx context.Context, writer http.ResponseWriter, r *http.Request) (any, error) {
	return "this is the response from root todo api info " + r.URL.Path, nil
}
