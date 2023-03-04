package assetstorage

import (
	"context"
	"trading-service/common"
	assetmodel "trading-service/modules/assets/model"
)

func (s *sqlStore) ListUserWallets(ctx context.Context, filter *assetmodel.Filter, moreInfo ...string) ([]assetmodel.UserAsset, error) {
	var result []assetmodel.UserAsset

	db := s.db.Table(assetmodel.UserAsset{}.TableName())

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	if f := filter; f != nil {
		if f.OwnerId > 0 {
			db = db.Where("user_id = ?", f.OwnerId)
		}
	}

	if err := db.Order("amount DESC").Find(&result).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	return result, nil
}

func (s *sqlStore) List(ctx context.Context, filter *assetmodel.Filter, paging *common.Paging, moreInfo ...string) ([]assetmodel.Asset, error) {
	var result []assetmodel.Asset

	db := s.db.Table(assetmodel.Asset{}.TableName())

	if f := filter; f != nil {
		switch f.Type {
		case "fiat":
			db = db.Where("type = ?", "fiat")
		case "crypto":
			db = db.Where("type = ?", "crypto")
		default:
		}
	}

	db = db.Where("status = 1")

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
