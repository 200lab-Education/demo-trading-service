package blockchain

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"math/big"
	appCommon "trading-service/common"
	txmodel "trading-service/modules/transaction/model"
)

func (p *scLogHdl) handleMoneyOut(ctx context.Context, l types.Log) error {
	log.Infof("Block %d - Tx %s - Event MoneyOut \n", l.BlockNumber, l.TxHash.Hex())

	tx, err := p.txStore.GetDataWithCondition(ctx, map[string]interface{}{"tx_hash": l.TxHash.Hex()})

	if err != nil && err != appCommon.ErrRecordNotFound {
		log.Errorln(err)
		return err
	}

	// No tx in DB

	event := struct {
		PaymentToken common.Address `json:"payment_token"`
		Amount       *big.Int       `json:"amount"`
		Fee          *big.Int       `json:"fee"`
	}{}

	err = p.bscExStore.ABI().UnpackIntoInterface(&event, "MoneyOut", l.Data)

	if err != nil {
		log.Errorln(err)
		return err
	}

	fee := decimal.NewNullDecimal(decimal.NewFromBigInt(event.Fee, 0))

	if tx != nil {
		verified := "verified"

		if err := p.txStore.Update(ctx, map[string]interface{}{"tx_hash": l.TxHash.Hex()}, &txmodel.BSCTransactionUpdate{
			Status:      &verified,
			BlockNumber: &l.BlockNumber,
			Fee:         &fee,
		}); err != nil {
			log.Errorln(err)
			return nil
		}

		if tx.LockId != nil && *tx.LockId > 0 {
			if err := p.assetStore.FinishLockUserAsset(ctx, *tx.LockId); err != nil {
				log.Errorln(err)
				return nil
			}
		}

		return nil
	}

	//sender := strings.ToLower(common.HexToAddress(l.Topics[1].Hex()).Hex())
	//receiver := strings.ToLower(common.HexToAddress(l.Topics[2].Hex()).Hex())
	//paymentToken := strings.ToLower(event.PaymentToken.Hex())
	//amount := decimal.NewNullDecimal(decimal.NewFromBigInt(event.Amount, 0))

	//newTx := txmodel.BSCTransaction{
	//	SQLModel:    appCommon.NewSQLModel(),
	//	Type:        "withdraw",
	//	TxHash:      l.TxHash.Hex(),
	//	BlockNumber: l.BlockNumber,
	//	EventName:   "MoneyOut",
	//	//SenderAddress:       &sender,
	//	ReceiverAddress:     &receiver,
	//	PayableTokenAddress: paymentToken,
	//	PayableTokenName:    "USDT",
	//	Amount:              amount,
	//	Fee:                 fee,
	//	Status:              "verified",
	//}
	//
	//if err := p.txStore.CreateData(ctx, &newTx); err != nil {
	//	return err
	//}

	return nil
}
