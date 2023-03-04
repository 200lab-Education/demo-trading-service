package txmodel

import (
	"github.com/shopspring/decimal"
	"trading-service/common"
)

const (
	EntityName  = "BSC Transaction"
	TbBSCTxName = "bsc_deposit_withdraws"
)

type BSCTransaction struct {
	common.SQLModel
	Type                string              `json:"type" gorm:"column:type;"`
	UserId              int                 `json:"user_id" gorm:"column:user_id;"`
	TxHash              string              `json:"tx_hash" gorm:"column:tx_hash;"`
	BlockNumber         uint64              `json:"block_number" gorm:"column:block_number;"`
	EventName           string              `json:"event_name" gorm:"column:event_name;"`
	SenderAddress       *string             `json:"sender_address" gorm:"column:sender_address;"`
	ReceiverAddress     *string             `json:"receiver_address" gorm:"column:receiver_address;"`
	PayableTokenAddress string              `json:"payable_token_address" gorm:"column:payable_token_address;"`
	PayableTokenName    string              `json:"payable_token_name" gorm:"column:payable_token_name;"`
	Status              string              `json:"status" gorm:"column:status;"`
	FailedReason        string              `json:"failed_reason" gorm:"column:failed_reason;"`
	LockId              *int                `json:"lock_id" gorm:"column:lock_id;"`
	Amount              decimal.NullDecimal `json:"amount" gorm:"column:amount;"`
	Fee                 decimal.NullDecimal `json:"fee" gorm:"column:fee;"`
}

func (BSCTransaction) TableName() string { return TbBSCTxName }

type BSCTransactionUpdate struct {
	Status       *string              `json:"status" gorm:"column:status;"`
	FailedReason *string              `json:"failed_reason" gorm:"column:failed_reason;"`
	BlockNumber  *uint64              `json:"block_number" gorm:"column:block_number;"`
	Fee          *decimal.NullDecimal `json:"fee" gorm:"column:fee;"`
}

func (BSCTransactionUpdate) TableName() string { return TbBSCTxName }
