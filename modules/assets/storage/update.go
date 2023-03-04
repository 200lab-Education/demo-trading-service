package assetstorage

import (
	"context"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"trading-service/common"
	assetmodel "trading-service/modules/assets/model"
)

func (s *sqlStore) IncreaseAmountAssetInWallet(ctx context.Context, userId, assetId, walletId int, amount decimal.NullDecimal) error {
	db := s.db.Begin()

	//ctxData := ctx.Value("data").(common.ContextData)
	//
	//masterTx := model.MasterTx{
	//	SQLModel:     common.NewSQLModel(),
	//	UserId:       userId,
	//	Amount:       amount,
	//	Type:         "in",
	//	WalletId:     walletId,
	//	AssetId:      assetId,
	//	RelatedId:    ctxData.GetRelatedId(),
	//	RelatedTable: ctxData.GetRelatedTable(),
	//}
	//
	//if err := db.Table(masterTx.TableName()).Create(&masterTx).Error; err != nil {
	//	db.Rollback()
	//	return common.ErrDB(err)
	//}

	var data assetmodel.UserAsset

	if err := db.Table(data.TableName()).
		Where("user_id = ? and asset_id = ? and wallet_id =?", userId, assetId, walletId).
		Updates(map[string]interface{}{"amount": gorm.Expr("amount + ?", amount.Decimal.String())}).Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	if err := db.Commit().Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}

func (s *sqlStore) CreateAssetInWallet(ctx context.Context, userId, assetId, walletId int, amount decimal.NullDecimal) error {
	db := s.db.Begin()

	//ctxData := ctx.Value("data").(common.ContextData)
	//
	//masterTx := model.MasterTx{
	//	SQLModel:     common.NewSQLModel(),
	//	UserId:       userId,
	//	Amount:       amount,
	//	Type:         "in",
	//	WalletId:     walletId,
	//	AssetId:      assetId,
	//	RelatedId:    ctxData.GetRelatedId(),
	//	RelatedTable: ctxData.GetRelatedTable(),
	//}
	//
	//if err := db.Table(masterTx.TableName()).Create(&masterTx).Error; err != nil {
	//	db.Rollback()
	//	return common.ErrDB(err)
	//}

	data := assetmodel.UserAsset{
		UserId:   userId,
		AssetId:  assetId,
		WalletId: walletId,
		Amount:   amount,
	}

	if err := db.Table(data.TableName()).Create(&data).Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	if err := db.Commit().Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}

func (s *sqlStore) GetAssetInWallet(ctx context.Context, userId, assetId, walletId int) (*assetmodel.UserAsset, error) {
	var data assetmodel.UserAsset

	if err := s.db.Table(data.TableName()).
		Where("user_id = ? and asset_id = ? and wallet_id =?", userId, assetId, walletId).
		First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrRecordNotFound
		}
		return nil, common.ErrDB(err)
	}

	return &data, nil
}
