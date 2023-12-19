package business

import (
	"auth-api-go/internal/data"
	"auth-api-go/internal/service/dao"
	"context"
	"database/sql"
	"fmt"
	"time"
)

func RightsBasedCheckToken(ctx context.Context, db *sql.DB, token string) (data.RightsBased, error) {
	response, err := CheckToken(ctx, db, token)
	if err != nil {
		return data.RightsBased{}, err
	}

	rights := make(map[string]bool)
	for _, val := range response.Rights {
		rights[val] = true
	}

	return data.RightsBased{
		Id:     response.Id,
		Login:  response.Login,
		Name:   response.Name,
		Rights: rights,
		Token:  response.Token,
	}, nil
}

func CheckToken(ctx context.Context, db *sql.DB, token string) (data.LoginResponse, error) {
	// dummy response if needed
	res := data.LoginResponse{}

	// get the token
	tokenData, err := dao.GetTokenData(ctx, db, token)
	if err != nil {
		return data.LoginResponse{}, err
	}

	// check expired flag or not
	if tokenData.Expired {
		return res, fmt.Errorf("the token %s is already expired", token)
	}

	// check the date is after the expiry date
	currentTime := time.Now()

	delta := currentTime.UnixMilli() - tokenData.ExpiryDt.UnixMilli()

	if delta > 0 {
		return res, fmt.Errorf("the token expired %d milliseconds ago, sorry about that", delta)
	}

	// get the user and the rights
	userId := tokenData.UserId
	res, err = dao.GetCheckTokenData(ctx, db, userId, token)
	if err != nil {
		return data.LoginResponse{}, err
	}

	// prepare and return the data
	return res, nil
}
