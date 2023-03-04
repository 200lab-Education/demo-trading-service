package common

import "github.com/shopspring/decimal"

type extMasterTx struct {
	relatedId    int
	relatedTable string
	moneyOutData *userWalletData
	moneyInData  *userWalletData
}

type ExtOption func(*extMasterTx)

func NewExtTxData(relativeId int, relatedTable string, opts ...ExtOption) *extMasterTx {
	d := &extMasterTx{
		relatedId:    relativeId,
		relatedTable: relatedTable,
	}

	for i := range opts {
		opts[i](d)
	}

	return d
}

func WithMoneyData(userId, assetId, walletId int, amount decimal.Decimal, isOut bool) ExtOption {
	return func(tx *extMasterTx) {
		mData := &userWalletData{
			userId:   userId,
			walletId: walletId,
			assetId:  assetId,
			amount:   amount,
		}

		if isOut {
			tx.moneyOutData = mData
			return
		}

		tx.moneyInData = mData
	}
}

func (c *extMasterTx) GetRelatedId() int {
	return c.relatedId
}

func (c *extMasterTx) GetRelatedTable() string {
	return c.relatedTable
}

func (c *extMasterTx) MoneyInData() UserWalletAsset {
	return c.moneyInData
}

func (c *extMasterTx) MoneyOutData() UserWalletAsset {
	return c.moneyOutData
}

type ExtendMasterTxData interface {
	MoneyInData() UserWalletAsset
	MoneyOutData() UserWalletAsset
	GetRelatedId() int
	GetRelatedTable() string
}

type userWalletData struct {
	userId   int
	walletId int
	assetId  int
	amount   decimal.Decimal
}

func (data *userWalletData) UserId() int             { return data.userId }
func (data *userWalletData) WalletId() int           { return data.walletId }
func (data *userWalletData) AssetId() int            { return data.assetId }
func (data *userWalletData) Amount() decimal.Decimal { return data.amount }

type UserWalletAsset interface {
	UserId() int
	WalletId() int
	AssetId() int
	Amount() decimal.Decimal
}
