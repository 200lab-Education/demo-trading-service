package txstorage

import (
	"context"
	"trading-service/common"
	txmodel "trading-service/modules/transaction/model"
)

func (s *sqlStore) Update(ctx context.Context, condition map[string]interface{}, data *txmodel.BSCTransactionUpdate) error {
	db := s.db.Begin()

	if err := db.Table(data.TableName()).Where(condition).Updates(data).Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	return nil
}
