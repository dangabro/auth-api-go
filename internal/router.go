package internal

import (
	"auth-api-go/internal/data"
	"auth-api-go/internal/handler/changePassword"
	"auth-api-go/internal/handler/changeUserRights"
	"auth-api-go/internal/handler/checkResources"
	"auth-api-go/internal/handler/checkToken"
	"auth-api-go/internal/handler/expire"
	"auth-api-go/internal/handler/extendToken"
	"auth-api-go/internal/handler/login"
	"auth-api-go/internal/handler/loginNoToken"
	"auth-api-go/internal/handler/logout"
	"auth-api-go/internal/handler/rights"
	"auth-api-go/internal/handler/root"
	"auth-api-go/internal/handler/searchUsers"
	"auth-api-go/internal/handler/updateUser"
	"auth-api-go/internal/handler/userTokens"
	"database/sql"
	"github.com/gorilla/mux"
)

func Start(config data.Config, db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	r.Handle("/", root.New(config, db)).Methods("GET")
	r.Handle("/authapi", root.New(config, db)).Methods("GET")
	r.Handle("/login", login.New(config, db)).Methods("POST")
	r.Handle("/checkToken", checkToken.New(config, db)).Methods("POST")
	r.Handle("/getAccessRights", rights.New(config, db)).Methods("GET")
	r.Handle("/expireTokens", expire.New(config, db)).Methods("POST")
	r.Handle("/changePassword", changePassword.New(config, db)).Methods("POST")
	r.Handle("/loginNoToken", loginNoToken.New(config, db)).Methods("POST")
	r.Handle("/logout", logout.New(config, db)).Methods("POST")
	r.Handle("/extendToken", extendToken.New(config, db)).Methods("POST")
	r.Handle("/checkResources", checkResources.New(config, db)).Methods("GET")
	r.Handle("/getUserTokens", userTokens.New(config, db)).Methods("POST")
	r.Handle("/searchUsers", searchUsers.New(config, db)).Methods("POST")
	r.Handle("/updateUser", updateUser.New(config, db)).Methods("POST")
	r.Handle("/changeUserRights", changeUserRights.New(config, db)).Methods("POST")

	return r
}
