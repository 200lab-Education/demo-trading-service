package walletmodel

import "trading-service/common"

type Wallet struct {
	common.SQLModel `json:",inline"`
	Name            string `json:"name" gorm:"column:name;"`
}

func (Wallet) TableName() string { return "wallets" }

func NewWallet(id int, name string) Wallet {
	return Wallet{
		SQLModel: common.SQLModel{Id: id},
		Name:     name,
	}
}

func (data *Wallet) Mask(isAdminOrOwner bool) {
	data.SQLModel.Mask(common.DbTypeWallet)
}
