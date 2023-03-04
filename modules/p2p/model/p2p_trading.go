package model

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
	"time"
	"trading-service/common"
	assetmodel "trading-service/modules/assets/model"
	bankmodel "trading-service/modules/banking/model"
)

type TradingStatus int

const (
	TradeStOpening TradingStatus = iota
	TradeStWaitingPay
	TradeStWaitVerify
	TradeStCancelled
	TradeStRejected
	TradeStDeleted
	TradeStVerified
)

var allTradeStatuses = []string{"opening", "wait_payment",
	"wait_verify", "cancelled", "rejected", "deleted", "verified"}

func StatusCanUpdate(s string) bool {
	return s != TradeStCancelled.String() && s != TradeStRejected.String() && s != TradeStDeleted.String() &&
		s != TradeStVerified.String()
}

func (s TradingStatus) String() string {
	return allTradeStatuses[s]
}

func ParseTradeStatus(s string) (TradingStatus, error) {
	for i := range allTradeStatuses {
		if allTradeStatuses[i] == s {
			return TradingStatus(i), nil
		}
	}

	return TradingStatus(0), errors.New("status not found")
}

type P2pTrading struct {
	common.SQLModel
	OrderId        int                    `json:"-" gorm:"column:order_id;"`
	OrderUserId    int                    `json:"-" gorm:"column:order_user_id;"`
	UserId         int                    `json:"-" gorm:"column:user_id;"`
	OfferAssetId   int                    `json:"-" gorm:"column:offer_asset_id;"`
	PayableAssetId int                    `json:"-" gorm:"column:payable_asset_id;"`
	Type           string                 `json:"type" gorm:"column:type;"`
	PaymentType    string                 `json:"payment_type" gorm:"column:payment_type;"`
	BankAccId      int                    `json:"-" gorm:"column:bank_acc_id;"`
	RefCode        string                 `json:"ref_code" gorm:"column:ref_code;"`
	Price          decimal.Decimal        `json:"price" gorm:"column:price;"`
	Quantity       decimal.Decimal        `json:"quantity" gorm:"column:quantity;"`
	BuyFee         decimal.Decimal        `json:"buy_fee" gorm:"column:buy_fee;"`
	SellFee        decimal.Decimal        `json:"sell_fee" gorm:"column:sell_fee;"`
	AuthByUserId   int                    `json:"auth_by_user_id" gorm:"column:auth_by_user_id;"`
	Status         string                 `json:"status" gorm:"column:status;"`
	FailedReason   string                 `json:"failed_reason" gorm:"column:failed_reason;"`
	WaitedAt       *time.Time             `json:"waited_at" gorm:"column:waited_at;"`
	VerifiedAt     *time.Time             `json:"verified_at" gorm:"column:verified_at;"`
	OfferAsset     *assetmodel.Asset      `json:"offer_asset" gorm:"foreignKey:OfferAssetId;"`
	PayableAsset   *assetmodel.Asset      `json:"payable_asset" gorm:"foreignKey:PayableAssetId;"`
	BankAccount    *bankmodel.BankAccount `json:"bank_account" gorm:"foreignKey:BankAccId;"`
	User           *common.SimpleUser     `json:"user" gorm:"foreignKey:UserId;"`
	OrderOwner     *common.SimpleUser     `json:"order_owner" gorm:"foreignKey:order_user_id;"`
	FakeOrderId    *common.UID            `json:"order_id" gorm:"-"`
	Logs           []P2pTradingLog        `json:"logs,omitempty" gorm:"foreignKey:TxId;references:OrderId;"`
}

func (P2pTrading) TableName() string {
	return "p2p_trades"
}

func (c *P2pTrading) Mask(isAdmin bool) {
	c.SQLModel.Mask(common.DbTypeP2pTrade)

	if v := c.User; v != nil {
		v.Mask(common.DbTypeUser)
	}

	if v := c.OrderOwner; v != nil {
		v.Mask(common.DbTypeUser)
	}

	if v := c.BankAccount; v != nil {
		v.Mask(common.DbTypeBankAcc)
	}

	if v := c.OfferAsset; v != nil {
		v.Mask(isAdmin)
	}

	if v := c.PayableAsset; v != nil {
		v.Mask(isAdmin)
	}

	orderId := common.NewUID(uint32(c.OrderId), int(common.DbTypeP2pOrder), 1)
	c.FakeOrderId = &orderId

	for i := range c.Logs {
		c.Logs[i].Mask()
	}
}

type P2pTradingCreate struct {
	common.SQLModel
	OrderId        int             `json:"-" gorm:"column:order_id;"`
	OrderUserId    int             `json:"-" gorm:"column:order_user_id;"`
	UserId         int             `json:"-" gorm:"column:user_id;"`
	OfferAssetId   int             `json:"-" gorm:"column:offer_asset_id;"`
	PayableAssetId int             `json:"-" gorm:"column:payable_asset_id;"`
	Type           string          `json:"-" gorm:"column:type;"`
	PaymentType    string          `json:"-" gorm:"column:payment_type;"`
	BankAccId      *common.UID     `json:"bank_acc_id" gorm:"column:bank_acc_id;"`
	RefCode        string          `json:"ref_code" gorm:"column:ref_code;"`
	Price          decimal.Decimal `json:"-" gorm:"column:price;"`
	Quantity       decimal.Decimal `json:"quantity" gorm:"column:quantity;"`
	BuyFee         decimal.Decimal `json:"-" gorm:"column:buy_fee;"`
	SellFee        decimal.Decimal `json:"-" gorm:"column:sell_fee;"`
	Status         string          `json:"-" gorm:"column:status;"`
}

func (P2pTradingCreate) TableName() string {
	return P2pTrading{}.TableName()
}

func (c *P2pTradingCreate) Validate() error {
	c.Type = strings.TrimSpace(strings.ToLower(c.Type))

	if c.Type != P2pTypeSell && c.Type != P2pTypeBuy {
		c.Type = P2pTypeSell
	}

	if c.Quantity.IsNegative() || c.Quantity.IsZero() {
		return ErrQuantity
	}

	c.RefCode = strings.TrimSpace(c.RefCode)

	if c.RefCode == "" {
		c.RefCode = fmt.Sprintf("%d", time.Now().UTC().UnixNano())
	}

	c.PaymentType = strings.TrimSpace(c.PaymentType)

	return nil
}

type P2pTradingUpdate struct {
	AuthByUserId int        `json:"-" gorm:"column:auth_by_user_id;"`
	Status       string     `json:"-" gorm:"column:status;"`
	FailedReason string     `json:"failed_reason" gorm:"column:failed_reason;"`
	WaitedAt     *time.Time `json:"-" gorm:"column:waited_at;"`
	VerifiedAt   *time.Time `json:"-" gorm:"column:verified_at;"`
	// More info
	IsAdminReject bool            `json:"-" gorm:"-"`
	OrderId       int             `json:"-" gorm:"-"`
	Quantity      decimal.Decimal `json:"-" gorm:"-"`
	Fee           decimal.Decimal `json:"-" gorm:"-"`
	IsOrderFilled bool            `json:"-" gorm:"-"`
}

func (P2pTradingUpdate) TableName() string {
	return P2pTrading{}.TableName()
}

func (c *P2pTradingUpdate) Validate() error {

	return nil
}
