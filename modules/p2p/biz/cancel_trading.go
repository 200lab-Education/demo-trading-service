package biz

import (
	"context"
	"fmt"
	"time"
	"trading-service/common"
	"trading-service/modules/p2p/model"
)

func (biz updateStatusBiz) CancelTrading(ctx context.Context, id int, updateData *model.P2pTradingUpdate) error {
	lockTradeKey := fmt.Sprintf(model.KeyLockTradeFmt, id)

	_ = biz.locker.Lock(lockTradeKey, time.Minute*2)
	defer biz.locker.Unlock(lockTradeKey)

	tradeData, err := biz.store.FindTrade(ctx, map[string]interface{}{"id": id})

	if err != nil {
		return common.ErrCannotGetEntity(model.TradeEntityName, err)
	}

	updateData.OrderId = tradeData.OrderId

	lockOrderKey := fmt.Sprintf(model.KeyLockOrderFmt, tradeData.OrderId)
	_ = biz.locker.Lock(lockOrderKey, time.Minute*2)
	defer biz.locker.Unlock(lockOrderKey)

	orderData, err := biz.store.Find(ctx, map[string]interface{}{"id": tradeData.OrderId})

	if err != nil {
		return common.ErrCannotGetEntity(model.OrderEntityName, err)
	}

	if !common.IsAdmin(biz.requester) {
		if tradeData.UserId != biz.requester.GetUserId() && orderData.UserId != biz.requester.GetUserId() {
			return common.ErrNoPermission(nil)
		}
	}

	if tradeData.Status != model.TradeStOpening.String() && tradeData.Status != model.TradeStWaitVerify.String() &&
		tradeData.Status != model.TradeStWaitingPay.String() {
		return model.ErrTradingCannotCancel
	}

	updateData.Quantity = tradeData.Quantity
	updateData.Fee = tradeData.SellFee

	if tradeData.Type == model.P2pTypeSell {
		ctx = context.WithValue(ctx, common.MasterTxData, common.NewExtTxData(
			0,
			updateData.TableName(),
			common.WithMoneyData(
				tradeData.UserId, tradeData.PayableAssetId, common.WalletSPOT, tradeData.Quantity, false,
			),
		))
	}

	if err := biz.store.CancelTrading(ctx, tradeData.UserId, id, updateData); err != nil {
		return common.ErrCannotUpdateEntity(model.TradeEntityName, err)
	}

	return nil
}
