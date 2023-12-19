package extendToken

import (
	"auth-api-go/internal/data"
	"auth-api-go/internal/handler"
	"auth-api-go/internal/service/dao"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type extendToken struct {
	config data.Config
	db     *sql.DB
}

func New(config data.Config, db *sql.DB) http.Handler {
	return &extendToken{
		config: config,
		db:     db,
	}
}

func (h *extendToken) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	res, err := h.process(ctx, writer, r)

	if err != nil {
		handler.Error(err, writer, http.StatusBadRequest)
		return
	}

	_ = handler.WriteJsonResponse(writer, res)
}

func (h *extendToken) process(ctx context.Context, writer http.ResponseWriter, r *http.Request) (any, error) {
	// check token from the header, must have at least auth_connect
	userLogin, err := handler.CheckRequestToken(ctx, h.db, r, []string{"auth_connect"}, true)
	if err != nil {
		return nil, err
	}

	loggedUserId := userLogin.Id

	// get payload
	simpleToken, err := h.getPayloadSimpleToken(r)
	if err != nil {
		return nil, err
	}

	// check token exists and it is not expired already
	tokenData, err := dao.GetTokenData(ctx, h.db, simpleToken.Token)
	if err != nil {
		return nil, err
	}

	err = h.checkExpiredToken(tokenData)
	if err != nil {
		return nil, err
	}

	// if does not belong to current user, first check the current user is auth_admin
	if tokenData.UserId != loggedUserId {
		err = userLogin.CheckEither([]string{data.AUTH_ADMIN})
		if err != nil {
			return nil, err
		}
	}

	// ok now everything checks, calculate current date and add the value
	tokenDuration := h.config.TokenDurationMs
	newExpiryDt := time.Now().Add(time.Millisecond * time.Duration(tokenDuration))

	dao.SetTokenExpiryDate(ctx, h.db, tokenData.Token, newExpiryDt)

	return nil, nil
}

func (h *extendToken) getPayloadSimpleToken(r *http.Request) (data.SimpleToken, error) {
	simpleToken := data.SimpleToken{}

	err := json.NewDecoder(r.Body).Decode(&simpleToken)
	if err != nil {
		return data.SimpleToken{}, err
	}

	return simpleToken, nil
}

func (h *extendToken) checkExpiredToken(tokenData data.TokenData) error {
	if tokenData.Expired {
		return errors.New("the expired flag of the provided  token is already set")
	}

	expiredDt := tokenData.ExpiryDt
	if time.Now().After(expiredDt) {
		return errors.New("the token already past due, cannot be extended")
	}

	return nil
}
