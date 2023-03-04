package storage

import (
	"context"
	"gorm.io/gorm"
	"trading-service/common"
	"trading-service/modules/banking/model"
)

func (s *sqlStore) Find(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.BankAccount, error) {
	db := s.db.Table(model.BankAccount{}.TableName())

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	var data model.BankAccount

	if err := db.Where(conditions).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrRecordNotFound
		}

		return nil, common.ErrDB(err)
	}

	return &data, nil
}
