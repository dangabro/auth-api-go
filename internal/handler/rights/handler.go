package rights

import (
	"auth-api-go/internal/data"
	"auth-api-go/internal/handler"
	"auth-api-go/internal/service/business"
	"context"
	"database/sql"
	"net/http"
)

type rights struct {
	config data.Config
	db     *sql.DB
}

func New(config data.Config, db *sql.DB) http.Handler {
	return &rights{
		config: config,
		db:     db,
	}
}

func (h *rights) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	_, err := handler.CheckRequestToken(ctx, h.db, r, []string{data.AUTH_CONNECT}, true)
	if err != nil {
		handler.Error(err, writer, http.StatusBadRequest)
		return
	}

	rights, err := business.NewRights(h.db).GetData(ctx)

	if err != nil {
		handler.Error(err, writer, http.StatusBadRequest)
		return
	}

	_ = handler.WriteJsonResponse(writer, rights)
}
