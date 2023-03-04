package storage

import (
	"context"
	"gorm.io/gorm"
	"trading-service/common"
	assetmodel "trading-service/modules/assets/model"
	"trading-service/modules/fiattx/model"
)

func (s *sqlStore) Update(ctx context.Context, condition map[string]interface{}, data *model.FiatDWUpdate) error {
	db := s.db.Table(data.TableName())

	if err := db.Where(condition).Updates(data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}

func (s *sqlStore) UpdateStatus(ctx context.Context, userId, id int, data *model.FiatDWUpdate) error {

	db := s.db.Begin()

	if err := db.Table(data.TableName()).Where("id = ?", id).Updates(data).Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	// Insert logs
	log := model.FiatDWLog{
		SQLModel: common.NewSQLModel(),
		TxId:     id,
		UserId:   userId,
	}

	st := *data.Status

	switch st {
	case model.StatusWaitVerify.String():
		log.Action = model.ActionPaid
	case model.StatusVerified.String():
		log.Action = model.ActionApproved
	case model.StatusRejected.String():
		log.Action = model.ActionRejected
	case model.StatusDeleted.String():
		log.Action = model.ActionDeleted
	case model.StatusCancelled.String():
		log.Action = model.ActionCancelled
	}

	if log.Action != "" {
		if err := db.Table(log.TableName()).Create(&log).Error; err != nil {
			db.Rollback()
			return common.ErrDB(err)
		}
	}

	if v := ctx.Value(common.MasterTxData); v != nil {
		dataMoneyIn := v.(common.ExtendMasterTxData).MoneyInData()

		if err := db.Table(assetmodel.UserAsset{}.TableName()).
			Where("user_id = ? and asset_id = ? and wallet_id =?", dataMoneyIn.UserId(), dataMoneyIn.AssetId(), dataMoneyIn.WalletId()).
			Updates(map[string]interface{}{"amount": gorm.Expr("amount + ?", dataMoneyIn.Amount().String())}).Error; err != nil {
			db.Rollback()
			return common.ErrDB(err)
		}
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	return nil
}
