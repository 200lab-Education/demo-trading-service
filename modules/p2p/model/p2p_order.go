package model

import (
	"errors"
	"github.com/shopspring/decimal"
	"strings"
	"trading-service/common"
	assetmodel "trading-service/modules/assets/model"
	"trading-service/modules/banking/model"
)

const (
	OrderEntityName = "P2pOrder"
	TradeEntityName = "P2pTrade"
)

type OrderStatus int

const (
	OrdStActive OrderStatus = iota
	OrdStStopped
	OrdStCancelled
	OrdStFilled
	OrdStDeleted
)

const (
	P2pTypeSell = "sell"
	P2pTypeBuy  = "buy"

	PaymentBankTransfer = "bank_transfer"

	KeyLockOrderFmt = "p2p.order/%d"
	KeyLockTradeFmt = "p2p.order.trade/%d"
)

var allOrdStatuses = []string{"active", "stopped", "cancelled",
	"filled", "deleted"}

var allP2pTypes = []string{P2pTypeBuy, P2pTypeSell}

func (s OrderStatus) String() string {
	return allOrdStatuses[s]
}

func ParseTxStatus(s string) (OrderStatus, error) {
	for i := range allOrdStatuses {
		if allOrdStatuses[i] == s {
			return OrderStatus(i), nil
		}
	}

	return OrderStatus(0), errors.New("status not found")
}

func AllowStatus(current, s string) bool {
	cst, err := ParseTxStatus(current)

	if err != nil {
		return false
	}

	sst, err := ParseTxStatus(s)

	if err != nil {
		return false
	}

	return cst < sst && !common.HasString(ReadOnlyStatuses, current)
}

var (
	PaymentTypes     = []string{PaymentBankTransfer}
	ReadOnlyStatuses = []string{OrdStStopped.String(), OrdStFilled.String(),
		OrdStCancelled.String(), OrdStDeleted.String()}
)

var (
	ErrAssetNotFound = func(root error) error {
		return common.NewCustomError(root, "asset not found", "ErrAssetNotFound")
	}

	ErrBankAccNotFound = func(root error) error {
		return common.NewCustomError(root, "bank account not found", "ErrBankAccNotFound")
	}

	ErrCannotOpen = func(root error) error {
		return common.NewCustomError(root, "cannot open trading", "ErrCannotOpen")
	}

	ErrPaymentTypeInvalid       = common.ValidateError("payment type is invalid", "ErrPaymentTypeInvalid")
	ErrOfferAssetIdInvalid      = common.ValidateError("offer asset id is invalid", "ErrOfferAssetIdInvalid")
	ErrPayableAssetIdInvalid    = common.ValidateError("payable asset id is invalid", "ErrPayableAssetIdInvalid")
	ErrBankAccIdInvalid         = common.ValidateError("bank acc is invalid", "ErrBankAccIdInvalid")
	ErrInvalidStatus            = common.ValidateError("invalid status stage, cannot update", "ErrInvalidStatus")
	ErrQuantity                 = common.ValidateError("quantity must be greater than zero", "ErrQuantity")
	ErrInvalidMinAmount         = common.ValidateError("invalid min amount", "ErrInvalidMinAmount")
	ErrInvalidMaxAmount         = common.ValidateError("invalid max amount", "ErrInvalidMaxAmount")
	ErrSameOfferPayable         = common.ValidateError("offer and payable id can not be the same", "ErrSameOfferPayable")
	ErrAmountOutMinMax          = common.ValidateError("amount must be in min-max range", "ErrAmountOutMinMax")
	ErrNotEnoughAssetBalance    = common.ValidateError("asset balance not enough", "ErrNotEnoughAssetBalance")
	ErrMinGreaterThanMaxAmount  = common.ValidateError("invalid min greater than max amount", "ErrMinGreaterThanMaxAmount")
	ErrMaxIsGreaterThanQuantity = common.ValidateError("max must be less then quantity", "ErrMaxIsGreaterThanQuantity")
	ErrInvalidPrice             = common.ValidateError("invalid price", "ErrInvalidPrice")
	ErrInvalidOrder             = common.ValidateError("order not found", "ErrInvalidOrder")
	ErrOpenYourselfOffer        = common.ValidateError("cannot open yourself transaction", "ErrOpenYourselfOffer")
	ErrOrderIsTrading           = common.ValidateError("order is trading", "ErrOrderIsTrading")
	ErrOrderNotEnoughAmount     = common.ValidateError("order not have enough amount", "ErrOrderNotEnoughAmount")
	ErrOrderStoppedFinished     = common.ValidateError("order has been stopped or finished", "ErrOrderStoppedFinished")
	ErrTradingCannotCancel      = common.ValidateError("trading cannot cancel or delete", "ErrTradingCannotCancel")
)

type P2pOrder struct {
	common.SQLModel
	UserId            int                 `json:"-" gorm:"column:user_id;"`
	OfferAssetId      int                 `json:"-" gorm:"column:offer_asset_id;"`
	PayableAssetId    int                 `json:"-" gorm:"column:payable_asset_id;"`
	Type              string              `json:"type" gorm:"column:type;"`
	PaymentType       string              `json:"payment_type" gorm:"column:payment_type;"`
	BankAccId         int                 `json:"-" gorm:"column:bank_acc_id;"`
	Price             decimal.NullDecimal `json:"price" gorm:"column:price;"`
	TotalQuantity     decimal.NullDecimal `json:"total_quantity" gorm:"column:total_quantity;"`
	AvailableQuantity decimal.NullDecimal `json:"available_quantity" gorm:"column:available_quantity;"`
	TotalFee          decimal.NullDecimal `json:"total_fee" gorm:"column:total_fee;"`
	AvailableFee      decimal.NullDecimal `json:"available_fee" gorm:"column:available_fee;"`
	SellFeeRate       float32             `json:"sell_fee_rate" gorm:"column:sell_fee_rate;"`
	BuyFeeRate        float32             `json:"buy_fee_rate" gorm:"column:buy_fee_rate;"`
	MinTradeAmount    decimal.NullDecimal `json:"min_trade_amount" gorm:"column:min_trade_amount;"`
	MaxTradeAmount    decimal.NullDecimal `json:"max_trade_amount" gorm:"column:max_trade_amount;"`
	Headline          string              `json:"headline" gorm:"column:headline;"`
	Term              string              `json:"term" gorm:"column:term;"`
	Status            string              `json:"status" gorm:"column:status;"`
	OfferAsset        *assetmodel.Asset   `json:"offer_asset" gorm:"foreignKey:OfferAssetId;"`
	PayableAsset      *assetmodel.Asset   `json:"payable_asset" gorm:"foreignKey:PayableAssetId;"`
	BankAccount       *model.BankAccount  `json:"bank_account" gorm:"foreignKey:BankAccId;"`
	User              *common.SimpleUser  `json:"user" gorm:"foreignKey:UserId;"`
}

func (P2pOrder) TableName() string {
	return "p2p_orders"
}

func (c *P2pOrder) Mask(isAdmin bool) {
	c.SQLModel.Mask(common.DbTypeP2pOrder)

	if v := c.User; v != nil {
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
}

type P2pOrderCreate struct {
	common.SQLModel
	UserId            int                 `json:"-" gorm:"column:user_id;"`
	OfferAssetId      *common.UID         `json:"offer_asset_id" gorm:"column:offer_asset_id;"`
	PayableAssetId    *common.UID         `json:"payable_asset_id" gorm:"column:payable_asset_id;"`
	Type              string              `json:"type" gorm:"column:type;"`
	PaymentType       string              `json:"payment_type" gorm:"column:payment_type;"`
	BankAccId         *common.UID         `json:"bank_acc_id" gorm:"column:bank_acc_id;"`
	Price             decimal.NullDecimal `json:"price" gorm:"column:price;"`
	TotalQuantity     decimal.Decimal     `json:"total_quantity" gorm:"column:total_quantity;"`
	AvailableQuantity decimal.NullDecimal `json:"-" gorm:"column:available_quantity;"`
	TotalFee          decimal.NullDecimal `json:"-" gorm:"column:total_fee;"`
	AvailableFee      decimal.NullDecimal `json:"-" gorm:"column:available_fee;"`
	SellFeeRate       float32             `json:"-" gorm:"column:sell_fee_rate;"`
	BuyFeeRate        float32             `json:"-" gorm:"column:buy_fee_rate;"`
	MinTradeAmount    decimal.NullDecimal `json:"min_trade_amount" gorm:"column:min_trade_amount;"`
	MaxTradeAmount    decimal.NullDecimal `json:"max_trade_amount" gorm:"column:max_trade_amount;"`
	Headline          string              `json:"headline" gorm:"column:headline;"`
	Term              string              `json:"term" gorm:"column:term;"`
	Status            string              `json:"-" gorm:"column:status;"`
}

func (P2pOrderCreate) TableName() string {
	return P2pOrder{}.TableName()
}

func (c *P2pOrderCreate) Validate() error {
	c.Type = strings.TrimSpace(strings.ToLower(c.Type))

	if c.Type != P2pTypeSell && c.Type != P2pTypeBuy {
		c.Type = P2pTypeSell
	}

	c.PaymentType = strings.TrimSpace(strings.ToLower(c.PaymentType))

	if !common.HasString(PaymentTypes, c.PaymentType) {
		c.PaymentType = PaymentTypes[0]
	}

	if c.TotalQuantity.IsNegative() || c.TotalQuantity.IsZero() {
		return ErrQuantity
	}

	if !c.MinTradeAmount.Valid || c.MinTradeAmount.Decimal.IsNegative() {
		return ErrInvalidMinAmount
	}

	if !c.MaxTradeAmount.Valid || c.MaxTradeAmount.Decimal.IsNegative() {
		return ErrInvalidMaxAmount
	}

	if !c.MaxTradeAmount.Decimal.GreaterThan(c.MinTradeAmount.Decimal) {
		return ErrMinGreaterThanMaxAmount
	}

	if c.MaxTradeAmount.Decimal.GreaterThan(c.TotalQuantity) {
		return ErrMaxIsGreaterThanQuantity
	}

	if !c.Price.Valid || c.Price.Decimal.IsZero() || c.Price.Decimal.IsNegative() {
		return ErrInvalidPrice
	}

	c.PaymentType = strings.TrimSpace(c.PaymentType)

	if !common.HasString(PaymentTypes, c.PaymentType) {
		return ErrPaymentTypeInvalid
	}

	if v := c.OfferAssetId; v == nil {
		return ErrOfferAssetIdInvalid
	}

	if v := c.PayableAssetId; v == nil {
		return ErrPayableAssetIdInvalid
	}

	if c.OfferAssetId.GetLocalID() == c.PayableAssetId.GetLocalID() {
		return ErrSameOfferPayable
	}

	if v := c.BankAccId; v == nil {
		return ErrBankAccIdInvalid
	}

	return nil
}

type P2pOrderUpdate struct {
	Price             decimal.NullDecimal `json:"price" gorm:"column:price;"`
	AvailableQuantity decimal.NullDecimal `json:"-" gorm:"column:available_quantity;"`
	AvailableFee      decimal.NullDecimal `json:"-" gorm:"column:available_fee;"`
	MinTradeAmount    decimal.NullDecimal `json:"min_trade_amount" gorm:"column:min_trade_amount;"`
	MaxTradeAmount    decimal.NullDecimal `json:"max_trade_amount" gorm:"column:max_trade_amount;"`
	Status            string              `json:"-" gorm:"column:status;"`
	IsAdminReject     bool                `json:"-" gorm:"-"`
}

func (P2pOrderUpdate) TableName() string {
	return P2pOrder{}.TableName()
}

func (c *P2pOrderUpdate) Validate() error {
	return nil
}
