package checkToken

import (
	"auth-api-go/internal/data"
	"auth-api-go/internal/handler"
	"auth-api-go/internal/service/business"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
)

type checkToken struct {
	config data.Config
	db     *sql.DB
}

func New(config data.Config, db *sql.DB) http.Handler {
	return &checkToken{
		config: config,
		db:     db,
	}
}

func (h *checkToken) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "method", "checkToken")

	// read login data from the request
	tokenData, err := readTokenData(r)
	if err != nil {
		handler.Error(err, writer, http.StatusBadRequest)
		return
	}

	checkTokenData, err := business.CheckToken(ctx, h.db, tokenData.Token)

	if err != nil {
		handler.Error(err, writer, http.StatusBadRequest)
		return
	}

	_ = handler.WriteJsonResponse(writer, checkTokenData)
}

func readTokenData(r *http.Request) (data.SimpleToken, error) {
	tokenData := data.SimpleToken{}

	err := json.NewDecoder(r.Body).Decode(&tokenData)
	if err != nil {
		return data.SimpleToken{}, err
	}

	return tokenData, nil
}
