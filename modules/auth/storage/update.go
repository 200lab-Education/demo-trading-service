package authstorage

import (
	"context"
	"trading-service/common"
	authmodel "trading-service/modules/auth/model"
)

func (s *sqlStore) UpdateUserWithCondition(
	ctx context.Context,
	data *authmodel.AuthDataUpdating,
	condition map[string]interface{},
) error {
	db := s.db

	db = db.Where(condition)

	if err := db.Table(data.TableName()).Updates(&data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
