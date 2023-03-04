package userstorage

import (
	"context"
	"trading-service/common"
	usermodel "trading-service/modules/user/model"
)

func (s *sqlStore) UpdateUser(ctx context.Context, condition map[string]interface{}, data *usermodel.UserUpdate) error {
	db := s.db.Table(data.TableName())

	if err := db.Where(condition).Updates(data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
