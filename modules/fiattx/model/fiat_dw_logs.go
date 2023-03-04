package model

import "trading-service/common"

const (
	ActionCreated   = "created"
	ActionPaid      = "paid"
	ActionCancelled = "cancelled"
	ActionApproved  = "approved"
	ActionRejected  = "rejected"
	ActionDeleted   = "deleted"
)

type FiatDWLog struct {
	common.SQLModel
	TxId   int                `json:"-" gorm:"column:tx_id;"`
	UserId int                `json:"-" gorm:"column:user_id;"`
	Action string             `json:"action" gorm:"column:action;"`
	User   *common.SimpleUser `json:"user" gorm:"foreignKey:UserId;"`
}

func (FiatDWLog) TableName() string {
	return "fiat_dw_logs"
}

func (l *FiatDWLog) Mask() {
	l.SQLModel.Mask(common.DbTypeFiatTxLog)

	if u := l.User; u != nil {
		u.Mask(common.DbTypeUser)
	}
}
