package storage

import (
	"context"
	"trading-service/common"
	"trading-service/modules/banking/model"
)

func (s *sqlStore) List(ctx context.Context, filter *model.Filter, paging *common.Paging, moreInfo ...string) ([]model.BankAccount, error) {
	var result []model.BankAccount

	db := s.db.Table(model.BankAccount{}.TableName())

	if f := filter; f != nil {
		if v := f.UserId; v > 0 {
			db = db.Where("user_id = ?", v)
		}

		if f.IsSystemAccount {
			db = db.Where("user_id = 0")
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
