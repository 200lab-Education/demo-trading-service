package storage

import (
	"context"
	"gorm.io/gorm"
	"trading-service/common"
	"trading-service/modules/location/model"
)

func (s *sqlStore) ListCountries(ctx context.Context, filter *model.Filter, paging *common.Paging, moreInfo ...string) ([]model.Country, error) {
	var result []model.Country

	db := s.db.Table(model.Country{}.TableName())

	//for i := range moreInfo {
	//	db = db.Preload(moreInfo[i])
	//}

	db = db.Preload("Cities", func(db *gorm.DB) *gorm.DB {
		return db.Order("cities.name ASC")
	})

	if f := filter; f != nil {
		if v := f.CountryId; v != "" {
			if uid, err := common.FromBase58(v); err == nil {
				db = db.Where("country_id = ?", uid.GetLocalID())
			}
		}
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
		Order("name asc").
		Find(&result).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	return result, nil
}
