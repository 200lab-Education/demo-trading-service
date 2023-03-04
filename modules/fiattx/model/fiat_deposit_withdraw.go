package model

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
	"time"
	"trading-service/common"
	assetmodel "trading-service/modules/assets/model"
	"trading-service/modules/banking/model"
)

const (
	EntityName = "FiatDeposit"
)

type TxStatus int

const (
	StatusPending TxStatus = iota
	StatusWaitingPay
	StatusPaymentTimeout
	StatusWaitVerify
	StatusVerified
	StatusCancelled
	StatusRejected
	StatusDeleted
)

var allTxStatuses = []string{"pending", "waiting_payment", "payment_timeout",
	"waiting_verify", "verified", "cancelled", "rejected", "deleted"}

func (s TxStatus) String() string {
	return allTxStatuses[s]
}

func ParseTxStatus(s string) (TxStatus, error) {
	for i := range allTxStatuses {
		if allTxStatuses[i] == s {
			return TxStatus(i), nil
		}
	}

	return TxStatus(0), errors.New("status not found")
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
	PaymentTypes     = []string{"bank_transfer"}
	ReadOnlyStatuses = []string{StatusVerified.String(), StatusCancelled.String(),
		StatusDeleted.String(), StatusRejected.String()}
)

var (
	ErrAssetNotFound = func(root error) error {
		return common.NewCustomError(root, "asset not found", "ErrAssetNotFound")
	}

	ErrBankAccNotFound = func(root error) error {
		return common.NewCustomError(root, "bank account not found", "ErrBankAccNotFound")
	}

	ErrPaymentTypeInvalid = common.ValidateError("payment type is invalid", "ErrPaymentTypeInvalid")
	ErrAssetIdInvalid     = common.ValidateError("asset id is invalid", "ErrAssetIdInvalid")
	ErrBankAccIdInvalid   = common.ValidateError("bank acc is invalid", "ErrBankAccIdInvalid")
	ErrInvalidStatus      = common.ValidateError("invalid status stage, cannot update", "ErrInvalidStatus")
	ErrInvalidAmount      = common.ValidateError("invalid amount", "ErrInvalidAmount")
)

type FiatDW struct {
	common.SQLModel
	UserId       int                 `json:"-" gorm:"column:user_id;"`
	AssetId      int                 `json:"-" gorm:"column:asset_id;"`
	PaymentType  string              `json:"payment_type" gorm:"column:payment_type;"`
	Type         string              `json:"type" gorm:"column:type;"`
	BankAccId    int                 `json:"-" gorm:"column:bank_acc_id;"`
	RefCode      string              `json:"ref_code" gorm:"column:ref_code;"`
	Amount       decimal.NullDecimal `json:"amount" gorm:"column:amount;"`
	Fee          decimal.NullDecimal `json:"fee" gorm:"column:fee;"`
	AuthByUserId int                 `json:"auth_by_user_id" gorm:"column:auth_by_user_id;"`
	Status       string              `json:"status" gorm:"column:status;"`
	FailedReason string              `json:"failed_reason" gorm:"column:failed_reason;"`
	WaitedAt     *time.Time          `json:"waited_at" gorm:"column:waited_at;"`
	VerifiedAt   *time.Time          `json:"verified_at" gorm:"column:verified_at;"`
	Asset        *assetmodel.Asset   `json:"asset" gorm:"foreignKey:AssetId;"`
	BankAccount  *model.BankAccount  `json:"bank_account" gorm:"foreignKey:BankAccId;"`
	User         *common.SimpleUser  `json:"user" gorm:"foreignKey:UserId;"`
	Logs         []FiatDWLog         `json:"logs,omitempty" gorm:"foreignKey:TxId;references:Id;"`
}

func (FiatDW) TableName() string {
	return "fiat_deposit_withdraws"
}

func (c *FiatDW) Mask(isAdmin bool) {
	c.SQLModel.Mask(common.DbTypeFiatTx)

	if v := c.User; v != nil {
		v.Mask(common.DbTypeUser)
	}

	if v := c.BankAccount; v != nil {
		v.Mask(common.DbTypeBankAcc)
	}

	if v := c.Asset; v != nil {
		v.Mask(isAdmin)
	}

	for i := range c.Logs {
		c.Logs[i].Mask()
	}
}

type FiatDWCreate struct {
	common.SQLModel
	UserId      int                  `json:"-" gorm:"column:user_id;"`
	AssetId     *common.UID          `json:"asset_id" gorm:"column:asset_id;"`
	PaymentType string               `json:"payment_type" gorm:"column:payment_type;"`
	BankAccId   *common.UID          `json:"bank_acc_id" gorm:"column:bank_acc_id;"`
	RefCode     string               `json:"ref_code" gorm:"column:ref_code;"`
	Amount      *decimal.Decimal     `json:"amount" gorm:"column:amount;"`
	Fee         *decimal.NullDecimal `json:"fee" gorm:"column:fee;"`
	Type        string               `json:"type" gorm:"column:type;"`
}

func (FiatDWCreate) TableName() string {
	return FiatDW{}.TableName()
}

func (c *FiatDWCreate) Validate() error {
	if c.Amount == nil || c.Amount.IsZero() || c.Amount.IsNegative() {
		return ErrInvalidAmount
	}

	c.RefCode = strings.TrimSpace(c.RefCode)

	if c.RefCode == "" {
		c.RefCode = fmt.Sprintf("%d", time.Now().UTC().UnixNano())
	}

	c.PaymentType = strings.TrimSpace(c.PaymentType)

	if !common.HasString(PaymentTypes, c.PaymentType) {
		return ErrPaymentTypeInvalid
	}

	if v := c.AssetId; v == nil {
		return ErrAssetIdInvalid
	}

	if v := c.BankAccId; v == nil {
		return ErrBankAccIdInvalid
	}

	return nil
}

type FiatDWUpdate struct {
	AuthByUserId *int       `json:"-" gorm:"column:auth_by_user_id;"`
	Status       *string    `json:"status" gorm:"column:status;"`
	FailedReason *string    `json:"failed_reason" gorm:"column:failed_reason;"`
	WaitedAt     *time.Time `json:"-" gorm:"column:waited_at;"`
	VerifiedAt   *time.Time `json:"-" gorm:"column:verified_at;"`
}

func (FiatDWUpdate) TableName() string {
	return FiatDW{}.TableName()
}
