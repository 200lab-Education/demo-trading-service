package storage

import (
	"context"
	"trading-service/common"
	"trading-service/modules/p2p/model"
)

func (s *sqlStore) List(ctx context.Context, filter *model.Filter, paging *common.Paging, moreInfo ...string) ([]model.P2pOrder, error) {
	var result []model.P2pOrder

	db := s.db.Table(model.P2pOrder{}.TableName())

	if f := filter; f != nil {

		if f.UserAccess {
			db = db.Where("(user_id = ? or status in ('active'))", f.RequesterId)
		} else {
			db = db.Where("status not in ('deleted')")

			if v := f.UserId; v != "" {
				db = db.Where("user_id = ?", v)
			}
		}

		if v := f.OfferAssetId; v != "" {
			uid, _ := common.FromBase58(v)
			db = db.Where("offer_asset_id = ?", uid.GetLocalID())
		}

		if v := f.Type; v != "" {
			db = db.Where("type = ?", v)
		}

		if v := f.Status; v != "" {
			db = db.Where("status = ?", v)
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

func (s *sqlStore) ListTrades(ctx context.Context, filter *model.Filter, paging *common.Paging, moreInfo ...string) ([]model.P2pTrading, error) {
	var result []model.P2pTrading

	db := s.db.Table(model.P2pTrading{}.TableName())

	if f := filter; f != nil {
		if f.UserAccess {
			db = db.Where("(user_id = ? or order_user_id = ?)", f.RequesterId, f.RequesterId)
		} else {
			db = db.Where("status not in ('deleted')")

			if v := f.UserId; v != "" {
				db = db.Where("user_id = ?", v)
			}
		}

		if v := f.OfferAssetId; v != "" {
			uid, _ := common.FromBase58(v)
			db = db.Where("offer_asset_id = ?", uid.GetLocalID())
		}

		if v := f.OrderId; v > 0 {
			db = db.Where("order_id = ?", v)
		}

		if v := f.UserId; v != "" {
			uid, _ := common.FromBase58(v)
			db = db.Where("user_id = ?", uid.GetLocalID())
		}

		if v := f.Type; v != "" {
			db = db.Where("type = ?", v)
		}

		if v := f.Status; v != "" {
			db = db.Where("status = ?", v)
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
