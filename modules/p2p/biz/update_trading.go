package biz

import (
	"context"
	"fmt"
	"time"
	"trading-service/common"
	"trading-service/modules/p2p/model"
)

func (biz updateStatusBiz) UpdateTrading(ctx context.Context, id int, updateData *model.P2pTradingUpdate) error {
	lockTradeKey := fmt.Sprintf(model.KeyLockTradeFmt, id)

	_ = biz.locker.Lock(lockTradeKey, time.Minute*2)
	defer biz.locker.Unlock(lockTradeKey)

	tradeData, err := biz.store.FindTrade(ctx, map[string]interface{}{"id": id})

	if err != nil {
		return common.ErrCannotGetEntity(model.TradeEntityName, err)
	}

	lockOrderKey := fmt.Sprintf(model.KeyLockOrderFmt, tradeData.OrderId)
	_ = biz.locker.Lock(lockOrderKey, time.Minute*2)
	defer biz.locker.Unlock(lockOrderKey)

	updateData.OrderId = tradeData.OrderId

	orderData, err := biz.store.Find(ctx, map[string]interface{}{"id": tradeData.OrderId})

	if err != nil {
		return common.ErrCannotGetEntity(model.OrderEntityName, err)
	}

	if !common.IsAdmin(biz.requester) {
		if tradeData.UserId != biz.requester.GetUserId() && orderData.UserId != biz.requester.GetUserId() {
			return common.ErrNoPermission(nil)
		}
	}

	if !model.StatusCanUpdate(tradeData.Status) {
		return model.ErrInvalidStatus
	}

	isOrderOwner := orderData.UserId == biz.requester.GetUserId()
	ordOwnerAllowStatuses := []string{model.TradeStVerified.String(), model.TradeStRejected.String()}
	ordTraderAllowStatuses := []string{model.TradeStWaitingPay.String(), model.TradeStWaitVerify.String(), model.TradeStRejected.String()}

	if orderData.Type == model.P2pTypeBuy {
		ordOwnerAllowStatuses, ordTraderAllowStatuses = ordTraderAllowStatuses, ordOwnerAllowStatuses
	}

	if !common.IsAdmin(biz.requester) {
		if (isOrderOwner && !common.HasString(ordOwnerAllowStatuses, updateData.Status)) ||
			(!isOrderOwner && !common.HasString(ordTraderAllowStatuses, updateData.Status)) {
			return model.ErrInvalidStatus
		}
	}

	now := time.Now().UTC()
	if updateData.Status == model.TradeStVerified.String() {
		updateData.VerifiedAt = &now
	}

	if updateData.Status == model.TradeStWaitingPay.String() || updateData.Status == model.TradeStWaitVerify.String() {
		updateData.WaitedAt = &now
	}

	if updateData.Status == model.TradeStVerified.String() {
		benefitUserId := tradeData.UserId
		benefitAssetIs := tradeData.OfferAssetId
		benefitQuantity := tradeData.Quantity.Sub(tradeData.BuyFee)

		if orderData.Type == model.P2pTypeBuy {
			benefitUserId = orderData.UserId
			benefitAssetIs = tradeData.PayableAssetId
			benefitQuantity = tradeData.Quantity.Sub(tradeData.SellFee)
		}

		ctx = context.WithValue(ctx, common.MasterTxData, common.NewExtTxData(
			id,
			updateData.TableName(),
			common.WithMoneyData(
				benefitUserId, benefitAssetIs, common.WalletSPOT, benefitQuantity, false,
			),
		))

		updateData.IsOrderFilled = orderData.AvailableQuantity.Decimal.IsZero()
	}

	if err := biz.store.UpdateTrading(ctx, tradeData.UserId, id, updateData); err != nil {
		return common.ErrCannotUpdateEntity(model.TradeEntityName, err)
	}

	return nil
}
