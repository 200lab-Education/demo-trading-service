package storage

import (
	"context"
	"gorm.io/gorm"
	"trading-service/common"
	assetmodel "trading-service/modules/assets/model"
	"trading-service/modules/p2p/model"
)

func (s *sqlStore) UpdateTrading(ctx context.Context, userId, id int, data *model.P2pTradingUpdate) error {
	db := s.db.Begin()

	if v := ctx.Value(common.MasterTxData); v != nil {
		dataMoneyIn := v.(common.ExtendMasterTxData).MoneyInData()

		if err := db.Table(assetmodel.UserAsset{}.TableName()).
			Where("user_id = ? and asset_id = ? and wallet_id =?", dataMoneyIn.UserId(), dataMoneyIn.AssetId(), dataMoneyIn.WalletId()).
			Updates(map[string]interface{}{"amount": gorm.Expr("amount + ?", dataMoneyIn.Amount().String())}).Error; err != nil {
			db.Rollback()
			return common.ErrDB(err)
		}
	}

	if err := db.Table(data.TableName()).Where("id = ?", id).Updates(data).Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	if data.Status == model.TradeStRejected.String() || data.Status == model.TradeStVerified.String() {
		// Insert logs
		log := model.P2pTradingLog{
			SQLModel: common.NewSQLModel(),
			TxId:     id,
			UserId:   userId,
			Action:   model.ActionApproved,
		}

		if data.Status == model.TradeStRejected.String() {
			log.Action = model.ActionRejected
		}

		if err := db.Table(log.TableName()).Create(&log).Error; err != nil {
			db.Rollback()
			return common.ErrDB(err)
		}
	}

	if data.IsOrderFilled {
		if err := db.Table(model.P2pOrder{}.TableName()).Where("id = ?", data.OrderId).Updates(map[string]interface{}{
			"status": model.OrdStFilled.String(),
		}).Error; err != nil {
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
