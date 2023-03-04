package authstorage

import (
	"context"
	"gorm.io/gorm"
	"trading-service/common"
	authmodel "trading-service/modules/auth/model"
)

func (s *sqlStore) FindUserWithCondition(
	ctx context.Context,
	condition map[string]interface{},
) (*authmodel.AuthData, error) {
	db := s.db

	var result authmodel.AuthData

	db = db.Where(condition)

	if err := db.Table(result.TableName()).First(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrRecordNotFound
		}

		return nil, common.ErrDB(err)
	}

	return &result, nil
}
