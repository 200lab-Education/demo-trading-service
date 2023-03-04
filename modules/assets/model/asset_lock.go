package assetmodel

import (
	"github.com/shopspring/decimal"
	"trading-service/common"
)

type UserAssetLock struct {
	common.SQLModel
	Type     string              `json:"type" gorm:"column:type;"`
	UserId   int                 `json:"user_id" gorm:"column:user_id;"`
	AssetId  int                 `json:"asset_id" gorm:"column:asset_id;"`
	WalletId int                 `json:"wallet_id" gorm:"column:wallet_id;"`
	Amount   decimal.NullDecimal `json:"amount" gorm:"column:amount;"`
	RefName  string              `json:"ref_name" gorm:"column:ref_name;"`
	RefId    int                 `json:"ref_id" gorm:"column:ref_id;"`
	Status   string              `json:"status" gorm:"column:status;"`
}

func (UserAssetLock) TableName() string {
	return "user_asset_locks"
}
