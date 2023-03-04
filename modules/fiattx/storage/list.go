package storage

import (
	"context"
	"trading-service/common"
	"trading-service/modules/fiattx/model"
)

func (s *sqlStore) List(ctx context.Context, filter *model.Filter, paging *common.Paging, moreInfo ...string) ([]model.FiatDW, error) {
	var result []model.FiatDW

	db := s.db.Table(model.FiatDW{}.TableName())

	if f := filter; f != nil {
		if v := f.UserId; v > 0 {
			db = db.Where("user_id = ?", v)
		}

		if v := f.AssetId; v > 0 {
			db = db.Where("asset_id = ?", v)
		}

		if v := f.Type; v != "" {
			db = db.Where("type = ?", v)
		}

		if v := f.Status; v != "" {
			db = db.Where("status = ?", v)
		}

		if v := f.PaymentType; v != "" {
			db = db.Where("payment_type = ?", v)
		}
	}

	db = db.Where("status not in ('deleted')")

	if err := db.Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	if v := paging.FakeCursor; v != "" {
		uid, err := common.FromBase58(v)

		if err != nil {
			return nil, common.ErrDB(err)
		}

		db = db.Where("id < ?", uid.GetLocalID())
	} else {
		offset := (paging.Page - 1) * paging.Limit
		db = db.Offset(offset)
	}

	if err := db.
		Limit(paging.Limit).
		Order("id desc").
		Find(&result).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	return result, nil
}
