package txbiz

import (
	"context"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"math/big"
	"strings"
	"time"
	"trading-service/common"
	assetmodel "trading-service/modules/assets/model"
	txmodel "trading-service/modules/transaction/model"
	"trading-service/plugin/locker"
)

type UserAssetStore interface {
	GetDataWithCondition(
		ctx context.Context,
		condition map[string]interface{},
		moreKeys ...string,
	) (*assetmodel.Asset, error)

	GetAssetInWallet(ctx context.Context, userId, assetId, walletId int) (*assetmodel.UserAsset, error)

	LockUserAsset(
		ctx context.Context,
		userId,
		assetId,
		walletId int,
		lockType string,
		amount decimal.NullDecimal,
	) (int, error)

	RefundUserAsset(
		ctx context.Context,
		lockId int,
		reason string,
	) error

	UpdateLockUserAsset(
		ctx context.Context,
		lockId int,
		updateData map[string]interface{},
	) error
}

type BSCTxStore interface {
	CreateData(
		ctx context.Context,
		tx *txmodel.BSCTransaction,
	) error

	Update(ctx context.Context, condition map[string]interface{}, data *txmodel.BSCTransactionUpdate) error
}

type SCCaller interface {
	Withdraw(ctx context.Context, receiver common2.Address, paymentToken common2.Address, amount *big.Int) (string, error)
}

type bscWithdrawBiz struct {
	requester      common.Requester
	userAssetStore UserAssetStore
	bscTxStore     BSCTxStore
	locker         locker.Locker
	scCaller       SCCaller
}

func NewBSCWithdrawBiz(
	requester common.Requester,
	userAssetStore UserAssetStore,
	bscTxStore BSCTxStore,
	locker locker.Locker,
	scCaller SCCaller,
) *bscWithdrawBiz {
	return &bscWithdrawBiz{
		requester:      requester,
		userAssetStore: userAssetStore,
		bscTxStore:     bscTxStore,
		locker:         locker,
		scCaller:       scCaller,
	}
}

func (biz *bscWithdrawBiz) WithdrawAsset(ctx context.Context, receiver common2.Address, assetId, walletId int, amount decimal.NullDecimal) (string, error) {
	lockKey := common.GetUserAssetWalletKey(biz.requester.GetUserId(), assetId, walletId)
	// Lock first for safe thread running
	biz.locker.Lock(lockKey, time.Second*10)
	defer biz.locker.Unlock(lockKey)

	// 1. Check user asset if enough balance
	asset, err := biz.userAssetStore.GetDataWithCondition(ctx, map[string]interface{}{"id": assetId})

	if err != nil {
		return "", txmodel.ErrCannotWithdraw(err)
	}

	if asset.Type != "crypto" {
		return "", txmodel.ErrInvalidAsset
	}

	userAsset, err := biz.userAssetStore.GetAssetInWallet(ctx, biz.requester.GetUserId(), assetId, walletId)

	if err != nil {
		return "", txmodel.ErrCannotWithdraw(err)
	}

	if userAsset.Amount.Decimal.LessThan(amount.Decimal) {
		return "", txmodel.ErrNotEnoughAsset
	}

	// 2. Lock amount of user
	lockId, err := biz.userAssetStore.LockUserAsset(ctx, biz.requester.GetUserId(), assetId, walletId, "withdraw_bsc", amount)

	if err != nil {
		return "", txmodel.ErrCannotWithdraw(err)
	}

	// 3. Call smart contract to withdraw asset

	txHash, err := biz.scCaller.Withdraw(ctx, receiver, common2.HexToAddress(asset.TokenAddress), amount.Decimal.BigInt())

	if err != nil {
		_ = biz.userAssetStore.RefundUserAsset(ctx, lockId, err.Error())
		return "", txmodel.ErrCannotWithdraw(err)
	}

	// 4. Create log transaction
	paymentToken := strings.ToLower(asset.TokenAddress)
	fee := decimal.NewNullDecimal(decimal.NewFromInt(0))
	receiverStr := receiver.String()

	newTx := txmodel.BSCTransaction{
		SQLModel:            common.NewSQLModel(),
		Type:                "withdraw",
		TxHash:              txHash,
		ReceiverAddress:     &receiverStr,
		BlockNumber:         0,
		EventName:           "MoneyOut",
		PayableTokenAddress: paymentToken,
		PayableTokenName:    "USDT", // fixed first
		Amount:              amount,
		Fee:                 fee,
		Status:              "pending",
		LockId:              &lockId,
	}

	_ = biz.userAssetStore.UpdateLockUserAsset(ctx, lockId, map[string]interface{}{"tx_hash": txHash})

	if err := biz.bscTxStore.CreateData(ctx, &newTx); err != nil {
		return "", txmodel.ErrCannotWithdraw(err)
	}

	// 5. Finish
	//_ = biz.userAssetStore.FinishLockUserAsset(ctx, lockId)

	return txHash, nil
}
