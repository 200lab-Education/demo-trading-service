package biz

import (
	"context"
	"github.com/shopspring/decimal"
	"time"
	"trading-service/common"
	assetmodel "trading-service/modules/assets/model"
	bankmdl "trading-service/modules/banking/model"
	"trading-service/modules/p2p/model"
	"trading-service/plugin/locker"
)

type CreateStore interface {
	Create(ctx context.Context, data *model.P2pOrderCreate) error
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
	locker     locker.Locker
}

func NewCreateOfferBiz(store CreateStore, assetStore AssetStore, bankStore BankStore, l locker.Locker) *createBiz {
	return &createBiz{store: store, assetStore: assetStore, bankStore: bankStore, locker: l}
}

func (biz createBiz) CreateSellOffer(ctx context.Context, data *model.P2pOrderCreate) error {
	if err := data.Validate(); err != nil {
		return err
	}

	assetData, err := biz.assetStore.GetDataWithCondition(ctx, map[string]interface{}{"id": data.OfferAssetId.GetLocalID()})

	if err != nil || assetData.Status == "inactive" || !assetData.AllowP2p {
		return model.ErrAssetNotFound(err)
	}

	data.AvailableQuantity = decimal.NewNullDecimal(data.TotalQuantity)
	data.TotalFee = decimal.NewNullDecimal(decimal.NewFromInt(0))
	data.AvailableFee = data.TotalFee
	data.BuyFeeRate = assetData.P2pFee
	data.SellFeeRate = assetData.P2pFee

	feeRate := int64(assetData.P2pFee * 1000)
	if feeRate > 0 {
		feeAmount := data.TotalQuantity.Mul(decimal.NewFromFloat32(assetData.P2pFee / 100))
		data.AvailableQuantity = decimal.NewNullDecimal(data.TotalQuantity.Sub(feeAmount))
		data.TotalFee = decimal.NewNullDecimal(feeAmount)
		data.AvailableFee = data.TotalFee
	}

	payAssetData, err := biz.assetStore.GetDataWithCondition(ctx, map[string]interface{}{"id": data.PayableAssetId.GetLocalID()})

	if err != nil || payAssetData.Status == "inactive" || !payAssetData.AllowP2p || payAssetData.Type != "fiat" {
		return model.ErrAssetNotFound(err)
	}

	if data.Type == model.P2pTypeSell {
		bankAcc, err := biz.bankStore.Find(ctx, map[string]interface{}{"id": data.BankAccId.GetLocalID()})

		if err != nil || bankAcc.Status == 0 || bankAcc.UserId == 0 || bankAcc.UserId != data.UserId {
			return model.ErrBankAccNotFound(err)
		}

		assetId := int(data.OfferAssetId.GetLocalID())

		lockUserWalletKey := common.GetUserAssetWalletKey(data.UserId, assetId, common.WalletSPOT)
		_ = biz.locker.Lock(lockUserWalletKey, time.Minute*2)
		defer biz.locker.Unlock(lockUserWalletKey)

		assetData, err := biz.assetStore.GetAssetInWallet(ctx, data.UserId, assetId, common.WalletSPOT)

		if err != nil {
			return model.ErrAssetNotFound(err)
		}

		if assetData.Amount.Decimal.LessThan(data.TotalQuantity) {
			return model.ErrNotEnoughAssetBalance
		}
	} else {
		data.BankAccId = nil
	}

	if err := biz.store.Create(ctx, data); err != nil {
		return common.ErrCannotCreateEntity(model.OrderEntityName, err)
	}

	return nil
}
