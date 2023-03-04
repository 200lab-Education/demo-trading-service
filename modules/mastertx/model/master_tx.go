package model

import (
	"github.com/shopspring/decimal"
	"trading-service/common"
)

type MasterTx struct {
	common.SQLModel
	UserId       int             `json:"-" gorm:"column:user_id;"`
	Amount       decimal.Decimal `json:"amount" gorm:"column:amount;"`
	Type         string          `json:"type" gorm:"column:type;"`
	WalletId     int             `json:"-" gorm:"column:wallet_id;"`
	AssetId      int             `json:"-" gorm:"column:asset_id;"`
	RelatedId    int             `json:"-" gorm:"column:related_id"`
	RelatedTable string          `json:"-" gorm:"column:related_table;"`
}

func (MasterTx) TableName() string {
	return "master_transactions"
}
