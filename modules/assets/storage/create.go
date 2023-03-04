package assetstorage

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"trading-service/common"
	assetmodel "trading-service/modules/assets/model"
)

func (s *sqlStore) CreateData(
	ctx context.Context,
	data *assetmodel.Asset,
) error {
	db := s.db

	if err := db.Create(data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}

func (s *sqlStore) LockUserAsset(
	ctx context.Context,
	userId,
	assetId,
	walletId int,
	lockType string,
	amount decimal.NullDecimal,
) (int, error) {
	db := s.db.Begin()

	data := assetmodel.UserAssetLock{
		SQLModel: common.NewSQLModel(),
		Type:     lockType,
		UserId:   userId,
		AssetId:  assetId,
		WalletId: walletId,
		Amount:   amount,
		Status:   "pending",
	}

	if err := db.Create(&data).Error; err != nil {
		db.Rollback()
		return 0, common.ErrDB(err)
	}

	if err := db.Table(assetmodel.UserAsset{}.TableName()).
		Where("user_id = ? and asset_id = ? and wallet_id =?", userId, assetId, walletId).
		Updates(map[string]interface{}{"amount": gorm.Expr("amount - ?", amount.Decimal.String())}).Error; err != nil {
		db.Rollback()
		return 0, common.ErrDB(err)
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return 0, common.ErrDB(err)
	}

	return data.Id, nil
}

func (s *sqlStore) RefundUserAsset(
	ctx context.Context,
	lockId int,
	reason string,
) error {
	// 1. Get lock details
	var lockData assetmodel.UserAssetLock
	tbUserAsset := assetmodel.UserAsset{}.TableName()
	tbUserAssetLock := assetmodel.UserAssetLock{}.TableName()

	if err := s.db.Where("id = ? and status = ?", lockId, "pending").First(&lockData).Error; err != nil {
		return common.ErrDB(err)
	}

	// 2. Begin db transaction
	dbTx := s.db.Begin()

	// 3. Refund user asset base on amount has been locked
	exprStr := fmt.Sprintf("amount + (select amount from %s where id = %d)", tbUserAssetLock, lockId)

	if err := dbTx.Table(tbUserAsset).
		Where("user_id = ? and asset_id = ? and wallet_id =?", lockData.UserId, lockData.AssetId, lockData.WalletId).
		Updates(map[string]interface{}{"amount": gorm.Expr(exprStr)}).Error; err != nil {
		dbTx.Rollback()
		return common.ErrDB(err)
	}

	// 4. Set status amount locked to 'refunded'
	if err := dbTx.Table(tbUserAssetLock).
		Where("id = ?", lockId).
		Updates(map[string]interface{}{"status": "refunded", "refund_reason": reason}).Error; err != nil {
		dbTx.Rollback()
		return common.ErrDB(err)
	}

	if err := dbTx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *sqlStore) FinishLockUserAsset(
	ctx context.Context,
	lockId int,
) error {
	tbUserAssetLock := assetmodel.UserAssetLock{}.TableName()

	if err := s.db.Table(tbUserAssetLock).
		Where("id = ?", lockId).
		Updates(map[string]interface{}{"status": "done"}).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}

func (s *sqlStore) UpdateLockUserAsset(
	ctx context.Context,
	lockId int,
	updateData map[string]interface{},
) error {
	tbUserAssetLock := assetmodel.UserAssetLock{}.TableName()

	if err := s.db.Table(tbUserAssetLock).
		Where("id = ?", lockId).
		Updates(updateData).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
