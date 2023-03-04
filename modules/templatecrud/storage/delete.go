package storage

import (
	"context"
	"trading-service/common"
	"trading-service/modules/banking/model"
)

func (s *sqlStore) Delete(ctx context.Context, condition map[string]interface{}) error {
	db := s.db.Table(model.BankAccount{}.TableName())

	if err := db.Where(condition).Updates(map[string]interface{}{"status": 0}).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
