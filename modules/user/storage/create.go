package userstorage

import (
	"context"
	"trading-service/common"
	usermodel "trading-service/modules/user/model"
)

func (s *sqlStore) CreateUser(ctx context.Context, data *usermodel.UserCreate) error {
	db := s.db.Begin()
	data.SQLModel = common.NewSQLModel()

	if err := db.Table(data.TableName()).Create(data).Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	return nil
}
