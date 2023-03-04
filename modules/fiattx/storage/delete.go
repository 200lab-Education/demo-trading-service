package storage

import (
	"context"
	"trading-service/common"
	"trading-service/modules/fiattx/model"
)

func (s *sqlStore) Delete(ctx context.Context, condition map[string]interface{}) error {
	db := s.db.Table(model.FiatDW{}.TableName())

	if err := db.Where(condition).Updates(map[string]interface{}{"status": "deleted"}).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
