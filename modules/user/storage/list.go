package userstorage

import (
	"context"
	"trading-service/common"
	usermodel "trading-service/modules/user/model"
)

func (s *sqlStore) ListUsers(ctx context.Context, filter *usermodel.Filter, paging *common.Paging, moreInfo ...string) ([]usermodel.User, error) {
	var result []usermodel.User

	db := s.db.Table(usermodel.User{}.TableName())

	if f := filter; f != nil {
	}

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
