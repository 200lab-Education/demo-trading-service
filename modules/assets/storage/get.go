package assetstorage

import (
	"context"
	"gorm.io/gorm"
	"trading-service/common"
	assetmodel "trading-service/modules/assets/model"
)

func (s *sqlStore) GetDataWithCondition(
	ctx context.Context,
	condition map[string]interface{},
	moreKeys ...string,
) (*assetmodel.Asset, error) {
	db := s.db

	var result assetmodel.Asset

	for i := range moreKeys {
		db = db.Preload(moreKeys[i]) // for auto preload
	}

	if err := db.
		Where(condition).
		First(&result).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrRecordNotFound
		}
		return nil, common.ErrDB(err)
	}

	return &result, nil
}
