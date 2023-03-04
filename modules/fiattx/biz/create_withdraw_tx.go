package biz

import (
	"context"
	"time"
	"trading-service/common"
	"trading-service/modules/fiattx/model"
)

func (biz createBiz) CreateWithdrawTx(ctx context.Context, data *model.FiatDWCreate) error {
	if err := data.Validate(); err != nil {
		return err
	}

	assetId := int(data.AssetId.GetLocalID())

	assetData, err := biz.assetStore.GetDataWithCondition(ctx, map[string]interface{}{"id": assetId})

	if err != nil || assetData.Status == "inactive" || !assetData.AllowWithdraw || assetData.Type != "fiat" {
		return model.ErrAssetNotFound(err)
	}

	bankAcc, err := biz.bankStore.Find(ctx, map[string]interface{}{"id": data.BankAccId.GetLocalID()})

	if err != nil || bankAcc.Status == 0 || bankAcc.UserId != data.UserId {
		return model.ErrBankAccNotFound(err)
	}

	lockUserWalletKey := common.GetUserAssetWalletKey(data.UserId, assetId, common.WalletSPOT)
	_ = biz.lck.Lock(lockUserWalletKey, time.Minute*2)
	defer biz.lck.Unlock(lockUserWalletKey)

	userAsset, err := biz.assetStore.GetAssetInWallet(ctx, data.UserId, assetId, common.WalletSPOT)

	if err != nil {
		return model.ErrAssetNotFound(err)
	}

	if userAsset.Amount.Decimal.LessThan(*data.Amount) {
		return common.ErrEnoughAmount
	}

	ctx = context.WithValue(ctx, common.MasterTxData, common.NewExtTxData(
		0,
		data.TableName(),
		common.WithMoneyData(
			data.UserId, assetId, common.WalletSPOT, *data.Amount, true,
		),
	))

	data.Type = "withdraw"

	if err := biz.store.Create(ctx, data); err != nil {
		return common.ErrCannotCreateEntity(model.EntityName, err)
	}

	return nil
}
