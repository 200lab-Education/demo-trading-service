package txstorage

import (
	"context"
	"trading-service/common"
	txmodel "trading-service/modules/transaction/model"
)

func (s *sqlStore) CreateData(
	ctx context.Context,
	tx *txmodel.BSCTransaction,
) error {
	db := s.db

	if err := db.Create(tx).Error; err != nil {

		return common.ErrDB(err)
	}

	return nil
}
