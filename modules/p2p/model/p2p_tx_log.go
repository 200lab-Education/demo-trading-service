package model

import "trading-service/common"

const (
	ActionOpenTrading = "opened"
	ActionPaid        = "paid"
	ActionCancelled   = "cancelled"
	ActionApproved    = "approved"
	ActionRejected    = "rejected"
	ActionDeleted     = "deleted"
)

type P2pTradingLog struct {
	common.SQLModel
	TxId   int                `json:"-" gorm:"column:tx_id;"`
	UserId int                `json:"-" gorm:"column:user_id;"`
	Action string             `json:"action" gorm:"column:action;"`
	User   *common.SimpleUser `json:"user" gorm:"foreignKey:UserId;"`
}

func (P2pTradingLog) TableName() string {
	return "p2p_trading_logs"
}

func (l *P2pTradingLog) Mask() {
	l.SQLModel.Mask(common.DbTypeP2pTradeLog)

	if u := l.User; u != nil {
		u.Mask(common.DbTypeUser)
	}
}
