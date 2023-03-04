package storage

import (
	"context"
	"trading-service/common"
	"trading-service/modules/banking/model"
)

func (s *sqlStore) Update(ctx context.Context, condition map[string]interface{}, data *model.BankAccountUpdate) error {
	db := s.db.Table(data.TableName())

	if err := db.Where(condition).Updates(data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
