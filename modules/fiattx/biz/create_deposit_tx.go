package biz

import (
	"context"
	"github.com/shopspring/decimal"
	"trading-service/common"
	assetmodel "trading-service/modules/assets/model"
	bankmdl "trading-service/modules/banking/model"
	"trading-service/modules/fiattx/model"
	"trading-service/plugin/locker"
)

type CreateStore interface {
	Create(ctx context.Context, data *model.FiatDWCreate) error
}

type AssetStore interface {
	GetDataWithCondition(
		ctx context.Context,
		condition map[string]interface{},
		moreKeys ...string,
	) (*assetmodel.Asset, error)

	CreateData(
		ctx context.Context,
		data *assetmodel.Asset,
	) error

	GetAssetInWallet(ctx context.Context, userId, assetId, walletId int) (*assetmodel.UserAsset, error)
	CreateAssetInWallet(ctx context.Context, userId, assetId, walletId int, amount decimal.NullDecimal) error
	IncreaseAmountAssetInWallet(ctx context.Context, userId, assetId, walletId int, amount decimal.NullDecimal) error
	FinishLockUserAsset(ctx context.Context, lockId int) error
}

type BankStore interface {
	Find(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*bankmdl.BankAccount, error)
}

type createBiz struct {
	store      CreateStore
	assetStore AssetStore
	bankStore  BankStore
	lck        locker.Locker
}

func NewCreateBiz(store CreateStore, assetStore AssetStore, bankStore BankStore, lck locker.Locker) *createBiz {
	return &createBiz{store: store, assetStore: assetStore, bankStore: bankStore, lck: lck}
}

func (biz createBiz) CreateDepositTx(ctx context.Context, data *model.FiatDWCreate) error {
	if err := data.Validate(); err != nil {
		return err
	}

	assetData, err := biz.assetStore.GetDataWithCondition(ctx, map[string]interface{}{"id": data.AssetId.GetLocalID()})

	if err != nil || assetData.Status == "inactive" || !assetData.AllowDeposit || assetData.Type != "fiat" {
		return model.ErrAssetNotFound(err)
	}

	bankAcc, err := biz.bankStore.Find(ctx, map[string]interface{}{"id": data.BankAccId.GetLocalID()})

	if err != nil || bankAcc.Status == 0 || bankAcc.UserId > 0 {
		return model.ErrBankAccNotFound(err)
	}

	data.Type = "deposit"

	if err := biz.store.Create(ctx, data); err != nil {
		return common.ErrCannotCreateEntity(model.EntityName, err)
	}

	return nil
}
