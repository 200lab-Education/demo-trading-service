package storage

import (
	"context"
	"gorm.io/gorm"
	"trading-service/common"
	"trading-service/modules/p2p/model"
)

func (s *sqlStore) Find(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.P2pOrder, error) {
	db := s.db.Table(model.P2pOrder{}.TableName())

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	var data model.P2pOrder

	if err := db.Where(conditions).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrRecordNotFound
		}

		return nil, common.ErrDB(err)
	}

	return &data, nil
}

func (s *sqlStore) FindTrade(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.P2pTrading, error) {
	db := s.db.Table(model.P2pTrading{}.TableName())

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	var data model.P2pTrading

	if err := db.Where(conditions).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrRecordNotFound
		}

		return nil, common.ErrDB(err)
	}

	return &data, nil
}

func (s *sqlStore) IsTrading(ctx context.Context, id int) (bool, error) {
	db := s.db.Table(model.P2pTrading{}.TableName())

	var count int64

	if err := db.
		Where(
			"order_id = ? and status in (?)", id,
			[]string{model.TradeStOpening.String(), model.TradeStWaitingPay.String(), model.TradeStWaitVerify.String()},
		).
		Count(&count).Error; err != nil {
		return false, common.ErrDB(err)
	}

	return count > 0, nil
}
