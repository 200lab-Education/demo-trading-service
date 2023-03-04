package assetmodel

import "trading-service/common"

type Asset struct {
	common.SQLModel   `json:",inline"`
	Type              string  `json:"type" gorm:"column:type;"`
	Name              string  `json:"name" gorm:"column:name;"`
	ShortName         string  `json:"short_name" gorm:"column:short_name;"`
	Symbol            string  `json:"symbol" gorm:"column:symbol;"`
	ChainName         string  `json:"chain_name" gorm:"column:chain_name;"`
	ChainId           int     `json:"chain_id" gorm:"column:chain_id;"`
	TokenAddress      string  `json:"token_address" gorm:"token_address;"`
	Decimal           int     `json:"decimal" gorm:"column:decimal;"`
	MinWithdrawAmount float32 `json:"min_withdraw_amount" gorm:"column:min_withdraw_amount;"`
	WithdrawFee       float32 `json:"withdraw_fee_rate" gorm:"column:withdraw_fee_rate;"`
	P2pFee            float32 `json:"p2p_fee_rate" gorm:"column:p2p_fee_rate;"`
	AllowWithdraw     bool    `json:"allow_withdraw" gorm:"column:allow_withdraw;"`
	AllowDeposit      bool    `json:"allow_deposit" gorm:"column:allow_deposit;"`
	AllowP2p          bool    `json:"allow_p2p" gorm:"column:allow_p2p;"`
	DisplayFormat     string  `json:"display_format" gorm:"column:display_format;"`
	Status            string  `json:"status" gorm:"column:status;"`
}

func (Asset) TableName() string { return "assets" }

func (data *Asset) Mask(isAdminOrOwner bool) {
	data.SQLModel.Mask(common.DbTypeAsset)
}
