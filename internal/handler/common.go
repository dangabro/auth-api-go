package handler

import (
	"auth-api-go/internal/data"
	"auth-api-go/internal/service/business"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteJsonResponse(writer http.ResponseWriter, payload any) error {
	status := http.StatusOK

	header := writer.Header()
	header.Set("Content-Type", "application/json")
	header.Set("Cache-Control", "no-cache")

	writer.WriteHeader(status)

	// write the actual object here
	bts, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = writer.Write(bts)
	if err != nil {
		return err
	}

	return nil
}

func WriteHtmlResponse(writer http.ResponseWriter, payload any) error {
	status := http.StatusOK
	header := writer.Header()
	header.Set("Content-Type", "text/html")
	header.Set("Cache-Control", "no-cache")

	writer.WriteHeader(status)

	// write the actual object here
	strPayload := fmt.Sprintf("%v", payload)
	bts := []byte(strPayload)

	_, err := writer.Write(bts)
	if err != nil {
		return err
	}

	return nil
}

func Error(err error, writer http.ResponseWriter, httpCode int) {
	message := "Error: " + err.Error()

	header := writer.Header()
	header.Set("Content-Type", "text/plain")
	header.Set("Cache-Control", "no-cache")
	writer.WriteHeader(httpCode)

	bts := []byte(message)

	_, _ = writer.Write(bts)
}

func CheckRequestToken(ctx context.Context, db *sql.DB, r *http.Request, rights []string, anyRight bool) (data.RightsBased, error) {
	if len(rights) == 0 {
		// there are no rights to check, getting out
		return data.RightsBased{}, nil
	}

	// retrieve token from the request
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return data.RightsBased{}, errors.New("cannot find authorization header token")
	}

	token := authHeader
	if strings.HasPrefix(authHeader, "bearer ") {
		token = authHeader[7:]
	}

	rBased, err := business.RightsBasedCheckToken(ctx, db, token)
	if err != nil {
		return rBased, err
	}

	if anyRight {
		return rBased, rBased.CheckEither(rights)
	} else {
		return rBased, rBased.CheckAll(rights)
	}

}

func ReadRequestLoginData(r *http.Request) (data.LoginData, error) {
	loginData := data.LoginData{}

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		return data.LoginData{}, err
	}

	return loginData, nil
}

func CheckPassword(ctx context.Context, password string) error {
	strPass := strings.TrimSpace(password)
	if len(strPass) == 0 {
		return errors.New("we don't accept empty passwords at this time")
	}

	return nil
}
