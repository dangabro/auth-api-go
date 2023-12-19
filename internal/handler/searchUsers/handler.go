package searchUsers

import (
	"auth-api-go/internal/data"
	"auth-api-go/internal/handler"
	"auth-api-go/internal/service/dao"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

type searchUsers struct {
	config data.Config
	db     *sql.DB
}

func New(config data.Config, db *sql.DB) http.Handler {
	return &searchUsers{
		config: config,
		db:     db,
	}
}

func (h *searchUsers) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	res, err := h.process(ctx, writer, r)

	if err != nil {
		handler.Error(err, writer, http.StatusBadRequest)
		log.Err(err).Msg("error searchUsers")

		return
	}

	log.Info().Msgf("success searchUsers")

	_ = handler.WriteJsonResponse(writer, res)
}

func (h *searchUsers) process(ctx context.Context, writer http.ResponseWriter, r *http.Request) (any, error) {
	// check token and get rights etc
	// check token from the header, must have at least auth_connect
	userLogin, err := handler.CheckRequestToken(ctx, h.db, r, []string{data.AUTH_CONNECT}, true)
	if err != nil {
		return nil, err
	}

	userId := userLogin.Id
	admin := true
	err = userLogin.CheckEither([]string{data.AUTH_ADMIN})
	if err != nil {
		// the token does not belong to user who is administrator
		admin = false
	}

	if admin {
		// parse the input
		var searchData data.SearchUserData
		err = json.NewDecoder(r.Body).Decode(&searchData)
		if err != nil {
			return nil, err
		}

		// process the search value
		search := processSearch(searchData.Search)

		// if admin, proceed with the search order by login
		arr, err := dao.SearchUser(ctx, h.db, search)
		if err != nil {
			return nil, err
		}

		if arr == nil {
			// return empty array rather than nil
			arr = []data.CompleteUserData{}
		}

		return data.UserSearchResult{Users: arr}, nil
	} else {
		// if not admin, load current user and send it back
		dt, err := dao.GetUserById(ctx, h.db, userId)
		if err != nil {
			return nil, err
		}

		usersArray := []data.CompleteUserData{dt}
		userSearchResult := data.UserSearchResult{Users: usersArray}
		return userSearchResult, nil
	}

}

func processSearch(search string) string {
	// add % at the beginning
	// add % at the end
	// replace * with %
	// replace ? with _

	val := strings.TrimSpace(search)

	val = strings.ReplaceAll(val, "*", "%")
	val = strings.ReplaceAll(val, "?", "_")

	if !strings.HasPrefix(val, "%") {
		val = "%" + val
	}

	if !strings.HasSuffix(val, "%") {
		val = val + "%"
	}

	return val
}
