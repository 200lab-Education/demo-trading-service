package storage

import (
	"context"
	"trading-service/common"
	"trading-service/modules/banking/model"
)

func (s *sqlStore) Create(ctx context.Context, data *model.BankAccountCreate) error {
	db := s.db

	data.SQLModel = common.NewSQLModel()

	if err := db.Table(data.TableName()).Create(data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
