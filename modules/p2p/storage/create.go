package storage

import (
	"context"
	"gorm.io/gorm"
	"trading-service/common"
	assetmodel "trading-service/modules/assets/model"
	"trading-service/modules/p2p/model"
)

func (s *sqlStore) Create(ctx context.Context, data *model.P2pOrderCreate) error {
	db := s.db.Begin()

	data.SQLModel = common.NewSQLModel()

	if err := db.Table(data.TableName()).Create(data).Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	if data.Type == model.P2pTypeSell {
		// lock asset
		if err := db.Table(assetmodel.UserAsset{}.TableName()).
			Where("user_id = ? and asset_id = ? and wallet_id =?", data.UserId, data.OfferAssetId.GetLocalID(), common.WalletSPOT).
			Updates(map[string]interface{}{"amount": gorm.Expr("amount - ?", data.TotalQuantity.String())}).Error; err != nil {
			db.Rollback()
			return common.ErrDB(err)
		}
	}

	//// Insert logs
	//log := model.P2pTradingLog{
	//	SQLModel: common.NewSQLModel(),
	//	TxId:     data.Id,
	//	UserId:   data.OfferUserId,
	//	Action:   model.ActionCreated,
	//}
	//
	//if err := db.Table(log.TableName()).Create(&log).Error; err != nil {
	//	db.Rollback()
	//	return common.ErrDB(err)
	//}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	return nil
}

func (s *sqlStore) OpenTrading(ctx context.Context, data *model.P2pTradingCreate) error {
	data.SQLModel = common.NewSQLModel()

	db := s.db.Begin()

	if v := ctx.Value(common.MasterTxData); v != nil {
		dataMoneyOut := v.(common.ExtendMasterTxData).MoneyOutData()

		if err := db.Table(assetmodel.UserAsset{}.TableName()).
			Where("user_id = ? and asset_id = ? and wallet_id =?", dataMoneyOut.UserId(), dataMoneyOut.AssetId(), dataMoneyOut.WalletId()).
			Updates(map[string]interface{}{"amount": gorm.Expr("amount - ?", dataMoneyOut.Amount().String())}).Error; err != nil {
			db.Rollback()
			return common.ErrDB(err)
		}
	}

	if err := db.Table(data.TableName()).Create(data).Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	{
		if err := db.Table(model.P2pOrder{}.TableName()).
			Where("id = ?", data.OrderId).
			Updates(map[string]interface{}{
				"available_quantity": gorm.Expr("available_quantity - ?", data.Quantity.String()),
				"available_fee":      gorm.Expr("available_fee - ?", data.BuyFee.String()),
			}).
			Error; err != nil {
			db.Rollback()
			return common.ErrDB(err)
		}
	}

	{
		// Insert logs
		log := model.P2pTradingLog{
			SQLModel: common.NewSQLModel(),
			TxId:     data.OrderId,
			UserId:   data.UserId,
			Action:   model.ActionOpenTrading,
		}

		if err := db.Table(log.TableName()).Create(&log).Error; err != nil {
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
