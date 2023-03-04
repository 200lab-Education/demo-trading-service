package storage

import (
	"context"
	"gorm.io/gorm"
	"trading-service/common"
	"trading-service/modules/fiattx/model"
)

func (s *sqlStore) Find(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.FiatDW, error) {
	db := s.db.Table(model.FiatDW{}.TableName())

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	var data model.FiatDW

	if err := db.Where(conditions).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrRecordNotFound
		}

		return nil, common.ErrDB(err)
	}

	return &data, nil
}
