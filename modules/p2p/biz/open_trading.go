package biz

import (
	"context"
	"fmt"
	"time"
	"trading-service/common"
	"trading-service/modules/p2p/model"
	"trading-service/plugin/locker"
)

type UpdateStatusStore interface {
	Find(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.P2pOrder, error)
	FindTrade(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.P2pTrading, error)
	IsTrading(ctx context.Context, id int) (bool, error)
	CancelOrder(ctx context.Context, userId, id int, data *model.P2pOrderUpdate) error
	OpenTrading(ctx context.Context, data *model.P2pTradingCreate) error
	CancelTrading(ctx context.Context, userId, id int, data *model.P2pTradingUpdate) error
	UpdateTrading(ctx context.Context, userId, id int, data *model.P2pTradingUpdate) error
}

type MasterTxStore interface {
	Create(ctx context.Context, data common.ExtendMasterTxData) error
}

type updateStatusBiz struct {
	locker      locker.Locker
	store       UpdateStatusStore
	assetStore  AssetStore
	bankStore   BankStore
	masterStore MasterTxStore
	requester   common.Requester
}

func NewUpdateStBiz(store UpdateStatusStore, locker locker.Locker, assetStore AssetStore,
	bankStore BankStore, masterStore MasterTxStore, requester common.Requester) *updateStatusBiz {
	return &updateStatusBiz{store: store, locker: locker, assetStore: assetStore, bankStore: bankStore,
		masterStore: masterStore, requester: requester}
}

func (biz updateStatusBiz) OpenTrading(ctx context.Context, id int, tradingData *model.P2pTradingCreate) error {
	if err := tradingData.Validate(); err != nil {
		return err
	}

	lockOrderKey := fmt.Sprintf(model.KeyLockOrderFmt, id)

	_ = biz.locker.Lock(lockOrderKey, time.Minute*2)
	defer biz.locker.Unlock(lockOrderKey)

	data, err := biz.store.Find(ctx, map[string]interface{}{"id": id})

	if err != nil {
		return common.ErrCannotGetEntity(model.OrderEntityName, err)
	}

	if data.Type == model.P2pTypeSell {
		tradingData.Type = model.P2pTypeBuy
	} else {
		tradingData.Type = model.P2pTypeSell
	}

	tradingData.OrderUserId = data.UserId
	tradingData.UserId = biz.requester.GetUserId()

	if biz.requester.GetUserId() == data.UserId {
		return model.ErrOpenYourselfOffer
	}

	if data.Status != model.OrdStActive.String() {
		return model.ErrOrderStoppedFinished
	}

	minAmount := data.MinTradeAmount.Decimal
	maxAmount := data.MaxTradeAmount.Decimal
	amount := tradingData.Quantity

	if (!minAmount.IsZero() && minAmount.GreaterThan(amount)) || (!maxAmount.IsZero() && maxAmount.LessThan(amount)) {
		return model.ErrAmountOutMinMax
	}

	if amount.GreaterThan(data.AvailableQuantity.Decimal) {
		return model.ErrOrderNotEnoughAmount
	}

	if data.Type == model.P2pTypeBuy {
		tradingData.Type = model.P2pTypeSell
		if data.PaymentType == model.PaymentBankTransfer {
			bankAcc, err := biz.bankStore.Find(ctx, map[string]interface{}{"id": tradingData.BankAccId.GetLocalID()})

			if err != nil || bankAcc.Status == 0 || bankAcc.UserId == 0 || bankAcc.UserId != tradingData.UserId {
				return model.ErrBankAccNotFound(err)
			}
		}

		assetId := data.OfferAssetId
		userId := biz.requester.GetUserId()

		lockUserWalletKey := common.GetUserAssetWalletKey(userId, assetId, common.WalletSPOT)
		_ = biz.locker.Lock(lockUserWalletKey, time.Minute*2)
		defer biz.locker.Unlock(lockUserWalletKey)

		assetData, err := biz.assetStore.GetAssetInWallet(ctx, userId, assetId, common.WalletSPOT)

		if err != nil {
			return model.ErrAssetNotFound(err)
		}

		payableAmount := tradingData.Quantity

		if assetData.Amount.Decimal.LessThan(payableAmount) {
			return model.ErrNotEnoughAssetBalance
		}

		ctx = context.WithValue(ctx, common.MasterTxData, common.NewExtTxData(
			0,
			tradingData.TableName(),
			common.WithMoneyData(
				userId, assetId, common.WalletSPOT, payableAmount, true,
			),
		))
	} else {
		tradingData.BankAccId = nil
		tradingData.Type = model.P2pTypeBuy
	}

	tradingData.OrderId = id
	tradingData.PaymentType = data.PaymentType
	tradingData.UserId = biz.requester.GetUserId()
	tradingData.OfferAssetId = data.OfferAssetId
	tradingData.PayableAssetId = data.PayableAssetId
	tradingData.Price = data.Price.Decimal
	tradingData.Status = model.TradeStOpening.String()
	bankAccId := common.NewUID(uint32(data.BankAccId), int(common.DbTypeBankAcc), 1)
	tradingData.BankAccId = &bankAccId

	if tradingData.Quantity.Equals(data.AvailableQuantity.Decimal) {
		tradingData.SellFee = data.AvailableFee.Decimal
		tradingData.BuyFee = data.AvailableFee.Decimal
	} else {
		feeAmount := tradingData.Quantity.Mul(data.TotalFee.Decimal).Div(data.TotalQuantity.Decimal)
		tradingData.SellFee = feeAmount
		tradingData.BuyFee = feeAmount
	}

	if err := biz.store.OpenTrading(ctx, tradingData); err != nil {
		return model.ErrCannotOpen(err)
	}

	return nil
}
