package blockchain

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"math/big"
	"strings"
	"time"
	appCommon "trading-service/common"
	"trading-service/component/asyncjob"
	assetmodel "trading-service/modules/assets/model"
	txmodel "trading-service/modules/transaction/model"
)

const (
	WalletSPOT    = 1
	WalletFunding = 2
)

func (p *scLogHdl) handleMoneyIn(ctx context.Context, l types.Log) error {
	log.Infof("Block %d - Tx %s - Event MoneyIn \n", l.BlockNumber, l.TxHash.Hex())

	tx, err := p.txStore.GetDataWithCondition(ctx, map[string]interface{}{"tx_hash": l.TxHash.Hex()})

	if tx != nil {
		return nil
	}

	if err != appCommon.ErrRecordNotFound {
		return err
	}

	// No tx in DB
	event := struct {
		PaymentToken common.Address `json:"payment_token"`
		Amount       *big.Int       `json:"amount"`
	}{}

	err = p.bscExStore.ABI().UnpackIntoInterface(&event, "MoneyIn", l.Data)

	if err != nil {
		log.Errorln(err)
		return err
	}

	sender := strings.ToLower(common.HexToAddress(l.Topics[1].Hex()).Hex())
	paymentToken := strings.ToLower(event.PaymentToken.Hex())
	amount := decimal.NewNullDecimal(decimal.NewFromBigInt(event.Amount, 0))
	fee := decimal.NewNullDecimal(decimal.NewFromInt(0))

	newTx := txmodel.BSCTransaction{
		SQLModel:            appCommon.NewSQLModel(),
		Type:                "deposit",
		TxHash:              l.TxHash.Hex(),
		BlockNumber:         l.BlockNumber,
		EventName:           "MoneyIn",
		SenderAddress:       &sender,
		PayableTokenAddress: paymentToken,
		PayableTokenName:    "USDT",
		Amount:              amount,
		Fee:                 fee,
		Status:              "verified",
	}

	u, err := p.findUser(ctx, sender)

	if err == nil {
		newTx.UserId = u.Id
	}

	if err := p.txStore.CreateData(ctx, &newTx); err != nil {
		return err
	}

	// Increase user asset, need to improve to be a message async call
	{
		job := asyncjob.NewJob(func(ctx context.Context) error {
			if u == nil {
				return nil
			}

			// 1. Find the asset record, insert if not existed
			asset, err := p.assetStore.GetDataWithCondition(ctx, map[string]interface{}{"token_address": paymentToken, "chain_id": p.bscExStore.ChainID()})

			if err != nil && err != appCommon.ErrRecordNotFound {
				return err
			}

			if err == appCommon.ErrRecordNotFound {
				asset = &assetmodel.Asset{
					SQLModel:          appCommon.SQLModel{},
					Type:              "crypto",
					Name:              "USDT Token",
					Symbol:            "USDT",
					ChainName:         "BSC",
					ChainId:           p.bscExStore.ChainID(),
					TokenAddress:      paymentToken,
					Decimal:           18,
					MinWithdrawAmount: 0,
					WithdrawFee:       0,
					AllowWithdraw:     true,
					AllowDeposit:      true,
					DisplayFormat:     "%symbol %amount",
					Status:            "active",
				}

				if err := p.assetStore.CreateData(ctx, asset); err != nil {
					return err
				}
			}

			// 2. Find user asset in SPOT Wallet
			// 2.1 Increase amount asset in wallet
			// 2.2 Else create amount asset in wallet
			userAsset, err := p.assetStore.GetAssetInWallet(ctx, u.Id, asset.Id, WalletSPOT)

			if err != nil && err != appCommon.ErrRecordNotFound {
				return err
			}

			if userAsset != nil {
				if err := p.assetStore.IncreaseAmountAssetInWallet(ctx, u.Id, asset.Id, WalletSPOT, amount); err != nil {
					return err
				}
				return nil
			}

			if err == appCommon.ErrRecordNotFound {
				if err := p.assetStore.CreateAssetInWallet(ctx, u.Id, asset.Id, WalletSPOT, amount); err != nil {
					return err
				}
			}

			return nil
		})

		job.SetRetryDurations(time.Second, time.Second*2, time.Second*3)

		_ = asyncjob.NewGroup(false, job).Run(context.Background())
	}

	return nil
}
