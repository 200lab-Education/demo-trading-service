package assetmodel

import (
	"github.com/shopspring/decimal"
	walletmodel "trading-service/modules/wallet/model"
)

type UserAsset struct {
	UserId   int                 `json:"-" gorm:"column:user_id;"`
	AssetId  int                 `json:"-" gorm:"column:asset_id;"`
	WalletId int                 `json:"-" gorm:"column:wallet_id;"`
	Amount   decimal.NullDecimal `json:"amount" gorm:"column:amount;"`
	Asset    *Asset              `json:"asset" gorm:"preload:false;foreignKey:AssetId;"`
	Wallet   *walletmodel.Wallet `json:"wallet" gorm:"-"`
}

func (UserAsset) TableName() string {
	return "user_assets"
}

func (data *UserAsset) Mask(isAdminOrOwner bool) {
	if v := data.Asset; v != nil {
		v.Mask(isAdminOrOwner)
	}

	if v := data.Wallet; v != nil {
		v.Mask(isAdminOrOwner)
	}
}
