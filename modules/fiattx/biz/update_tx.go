package biz

import (
	"context"
	"time"
	"trading-service/common"
	"trading-service/modules/fiattx/model"
	"trading-service/plugin/locker"
)

type UpdateStatusStore interface {
	Find(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.FiatDW, error)
	UpdateStatus(ctx context.Context, userId, id int, data *model.FiatDWUpdate) error
}

type MasterTxStore interface {
	Create(ctx context.Context, data common.ExtendMasterTxData) error
}

type updateStatusBiz struct {
	locker        locker.Locker
	store         UpdateStatusStore
	assetStore    AssetStore
	masterTxStore MasterTxStore
	requester     common.Requester
}

func NewUpdateStBiz(
	store UpdateStatusStore,
	locker locker.Locker,
	assetStore AssetStore,
	masterTxStore MasterTxStore,
	requester common.Requester,
) *updateStatusBiz {
	return &updateStatusBiz{store: store, locker: locker, assetStore: assetStore,
		masterTxStore: masterTxStore, requester: requester}
}

func (biz updateStatusBiz) UpdateStatus(ctx context.Context, id int, updateData *model.FiatDWUpdate) error {
	data, err := biz.store.Find(ctx, map[string]interface{}{"id": id})

	if err != nil {
		return common.ErrCannotGetEntity(model.EntityName, err)
	}

	lockKey := common.GetUserAssetWalletKey(data.UserId, data.AssetId, common.WalletSPOT)
	_ = biz.locker.Lock(lockKey, time.Minute*2)
	defer biz.locker.Unlock(lockKey)

	if !common.IsAdmin(biz.requester) && data.UserId != biz.requester.GetUserId() {
		return common.ErrNoPermission(nil)
	}

	if !model.AllowStatus(data.Status, *updateData.Status) {
		return model.ErrInvalidStatus
	}

	now := time.Now().UTC()
	status := *updateData.Status

	if status == model.StatusVerified.String() {
		updateData.VerifiedAt = &now
		adminId := biz.requester.GetUserId()
		updateData.AuthByUserId = &adminId
	}

	if status == model.StatusWaitingPay.String() || status == model.StatusWaitVerify.String() {
		updateData.WaitedAt = &now
	}

	if (status == model.StatusRejected.String() || status == model.StatusCancelled.String()) && data.Type == "withdraw" {
		ctx = context.WithValue(ctx, common.MasterTxData, common.NewExtTxData(
			0,
			data.TableName(),
			common.WithMoneyData(
				data.UserId, data.AssetId, common.WalletSPOT, data.Amount.Decimal, false,
			),
		))
	}

	if err := biz.store.UpdateStatus(ctx, biz.requester.GetUserId(), id, updateData); err != nil {
		return common.ErrCannotUpdateEntity(model.EntityName, err)
	}

	_ = biz.masterTxStore.Create(ctx, common.NewExtTxData(id, updateData.TableName(),
		common.WithMoneyData(data.UserId, data.AssetId, common.WalletSPOT, data.Amount.Decimal, false)))

	if *updateData.Status == model.StatusVerified.String() {
		_, err := biz.assetStore.GetAssetInWallet(ctx, data.UserId, data.AssetId, common.WalletSPOT)

		if err != nil && err != common.ErrRecordNotFound {
			return err
		}

		if err == common.ErrRecordNotFound {
			if err := biz.assetStore.CreateAssetInWallet(ctx, data.UserId, data.AssetId, common.WalletSPOT, data.Amount); err != nil {
				return err
			}

			return nil
		}

		if err := biz.assetStore.IncreaseAmountAssetInWallet(ctx, data.UserId, data.AssetId, common.WalletSPOT, data.Amount); err != nil {
			return err
		}
	}

	return nil
}
