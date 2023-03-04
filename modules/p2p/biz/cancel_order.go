package biz

import (
	"context"
	"fmt"
	"time"
	"trading-service/common"
	"trading-service/modules/p2p/model"
)

func (biz updateStatusBiz) CancelOrDeleteOrder(ctx context.Context, id int, updateData *model.P2pOrderUpdate) error {
	lockOrderKey := fmt.Sprintf(model.KeyLockOrderFmt, id)

	_ = biz.locker.Lock(lockOrderKey, time.Minute*2)
	defer biz.locker.Unlock(lockOrderKey)

	data, err := biz.store.Find(ctx, map[string]interface{}{"id": id})

	if err != nil {
		return common.ErrCannotGetEntity(model.OrderEntityName, err)
	}

	if biz.requester.GetUserId() != data.UserId && !common.IsAdmin(biz.requester) {
		return model.ErrInvalidOrder
	}

	if common.HasString(model.ReadOnlyStatuses, data.Status) {
		return model.ErrInvalidStatus
	}

	if isTrading, err := biz.store.IsTrading(ctx, id); isTrading || err != nil {
		return model.ErrOrderIsTrading
	}

	if data.Type == model.P2pTypeSell {
		totalRefundAmount := data.AvailableQuantity.Decimal.Add(data.AvailableFee.Decimal)
		ctx = context.WithValue(ctx, common.MasterTxData, common.NewExtTxData(
			id,
			updateData.TableName(),
			common.WithMoneyData(
				data.UserId, data.OfferAssetId, common.WalletSPOT, totalRefundAmount, false,
			),
		))

		lockUserWalletKey := common.GetUserAssetWalletKey(data.UserId, data.OfferAssetId, common.WalletSPOT)
		_ = biz.locker.Lock(lockUserWalletKey, time.Minute*2)
		defer biz.locker.Unlock(lockUserWalletKey)
	}

	if err := biz.store.CancelOrder(ctx, data.UserId, data.Id, updateData); err != nil {
		return common.ErrCannotUpdateEntity(model.OrderEntityName, err)
	}

	return nil
}
