package business

import (
	"auth-api-go/internal/data"
	"auth-api-go/internal/service/dao"
	"context"
	"database/sql"
)

type rightsService struct {
	db *sql.DB
}

func NewRights(db *sql.DB) *rightsService {
	return &rightsService{
		db: db,
	}
}

func (service *rightsService) GetData(ctx context.Context) (data.Rights, error) {
	rights, err := dao.GetRightDataList(ctx, service.db)

	if err != nil {
		return data.Rights{}, err
	}

	return data.Rights{
		List: rights,
	}, nil
}
