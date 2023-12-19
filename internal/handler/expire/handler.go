package expire

import (
	"auth-api-go/internal/data"
	"auth-api-go/internal/handler"
	"auth-api-go/internal/service/dao"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type expireTokens struct {
	config data.Config
	db     *sql.DB
}

func New(config data.Config, db *sql.DB) http.Handler {
	return &expireTokens{
		config: config,
		db:     db,
	}
}

func (h *expireTokens) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	res, err := h.process(ctx, writer, r)

	if err != nil {
		handler.Error(err, writer, http.StatusBadRequest)
		return
	}

	_ = handler.WriteJsonResponse(writer, res)
}

func (h *expireTokens) process(ctx context.Context, writer http.ResponseWriter, r *http.Request) (any, error) {
	// get token
	// check token with admin_connect
	loggedIn, err := handler.CheckRequestToken(ctx, h.db, r, []string{data.AUTH_CONNECT}, true)
	if err != nil {
		return nil, err
	}

	// get the user id associated with the token
	userId := loggedIn.Id

	// read request holder
	idsHolder, err := readIdsHolder(r)
	if err != nil {
		return nil, err
	}

	// ensure at least one token is being passed
	ids := idsHolder.Ids
	err = checkAtLeastOneToken(ids)
	if err != nil {
		return nil, err
	}

	// ensure all the token exist
	err = ensureAllTokensExist(ctx, h.db, ids)
	if err != nil {
		return nil, err
	}

	// see who are the users token belong to
	userIds, err := getUserIdsForTokens(ctx, h.db, ids)
	if err != nil {
		return nil, err
	}

	// if at least one of the users is not the logged in user then check admin
	different := atLeastOneDifferentUser(userId, userIds)
	if different {
		err = loggedIn.CheckEither([]string{data.AUTH_ADMIN})
		if err != nil {
			return nil, err
		}
	}

	// go ahead and expire the tokens
	err = dao.ExpireTokens(ctx, h.db, ids)
	if err != nil {
		return nil, err
	}

	return data.NewSuccessData(), nil
}

func atLeastOneDifferentUser(id string, ids map[string]bool) bool {
	for currentId, _ := range ids {
		if id != currentId {
			// at least one user is not the logged in user
			return true
		}
	}
	return false
}

func getUserIdsForTokens(ctx context.Context, db *sql.DB, tokens []string) (map[string]bool, error) {
	return dao.GetDistinctUserIdsForTokens(ctx, db, tokens)
}

func checkAtLeastOneToken(tokens []string) error {
	if len(tokens) == 0 {
		return errors.New("please provide at least one token value in the request")
	}

	return nil
}

func ensureAllTokensExist(ctx context.Context, db *sql.DB, ids []string) error {
	var tokenMap map[string]bool
	var err error
	tokenMap, err = dao.GetTokenValues(ctx, db, ids)

	if err != nil {
		return err
	}

	for _, id := range ids {
		_, ok := tokenMap[id]
		if !ok {
			return fmt.Errorf("at least the token %s does not exist", id)
		}
	}

	return nil
}

func readIdsHolder(r *http.Request) (data.IdsHolder, error) {
	idsHolder := data.IdsHolder{}

	err := json.NewDecoder(r.Body).Decode(&idsHolder)
	if err != nil {
		return data.IdsHolder{}, err
	}

	return idsHolder, nil
}
