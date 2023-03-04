package authstorage

import (
	"context"
	"trading-service/common"
	authmodel "trading-service/modules/auth/model"
)

func (s *sqlStore) CreateUser(
	ctx context.Context,
	data *authmodel.AuthDataCreation,
) error {

	db := s.db

	if err := db.Table(data.TableName()).Create(data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
