package data

import (
	"time"
)

type LoginData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserData struct {
	Id       string
	Name     string
	Login    string
	Password string
}

type CompleteUserData struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Login    string   `json:"login"`
	Password string   `json:"-"`
	Rights   []string `json:"rights"`
}

type UserSearchResult struct {
	Users []CompleteUserData `json:"users"`
}

type TokenData struct {
	Id       string    `json:"id"`
	UserId   string    `json:"-"`
	Token    string    `json:"token"`
	Added    time.Time `json:"createdDate"`
	Expired  bool      `json:"expired"`
	ExpiryDt time.Time `json:"expireDate"`
}

type TokenResponse struct {
	Tokens []TokenData `json:"tokens"`
}

type LoginResponse struct {
	Id     string   `json:"id"`
	Login  string   `json:"login"`
	Name   string   `json:"name"`
	Rights []string `json:"rights"`
	Token  string   `json:"token"`
}

type SimpleToken struct {
	Token string `json:"token"`
}

type RightData struct {
	Cd   string `json:"cd"`
	Name string `json:"name"`
}

type Rights struct {
	List []RightData `json:"rights"`
}

type IdsHolder struct {
	Ids []string `json:"ids"`
}

type IdHolder struct {
	Id string `json:"id"`
}

type SearchUserData struct {
	Search string `json:"searchString"`
}

type UpdateUserData struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
}

type ChangeUserRights struct {
	Id     string   `json:"id"`
	Rights []string `json:"rights"'`
}
