package storage

import (
	"context"
	"gorm.io/gorm"
	"trading-service/common"
	assetmodel "trading-service/modules/assets/model"
	"trading-service/modules/fiattx/model"
)

func (s *sqlStore) Create(ctx context.Context, data *model.FiatDWCreate) error {
	db := s.db.Begin()

	data.SQLModel = common.NewSQLModel()

	if err := db.Table(data.TableName()).Create(data).Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	if v := ctx.Value(common.MasterTxData); v != nil {
		dataMoneyOut := v.(common.ExtendMasterTxData).MoneyOutData()

		if err := db.Table(assetmodel.UserAsset{}.TableName()).
			Where("user_id = ? and asset_id = ? and wallet_id =?", dataMoneyOut.UserId(), dataMoneyOut.AssetId(), dataMoneyOut.WalletId()).
			Updates(map[string]interface{}{"amount": gorm.Expr("amount - ?", dataMoneyOut.Amount().String())}).Error; err != nil {
			db.Rollback()
			return common.ErrDB(err)
		}
	}

	// Insert logs
	log := model.FiatDWLog{
		SQLModel: common.NewSQLModel(),
		TxId:     data.Id,
		UserId:   data.UserId,
		Action:   model.ActionCreated,
	}

	if err := db.Table(log.TableName()).Create(&log).Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	return nil
}
