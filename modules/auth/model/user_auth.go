package authmodel

import (
	"errors"
	common2 "github.com/ethereum/go-ethereum/common"
	"strings"
	"time"
	"trading-service/common"
)

var (
	ErrWalletAddressInvalid       = errors.New("wallet address is invalid")
	ErrWalletAddressHasBinded     = errors.New("wallet address has connected to another wallet")
	ErrNonceMustBeGreaterThanZero = errors.New("nonce must be greater than zero")
	ErrSignatureInvalid           = errors.New("signature is invalid")
)

const (
	EntityName = "User"
	TbName     = "app_users"
)

type AuthData struct {
	Id            int    `json:"-" gorm:"column:id;"`
	WalletAddress string `json:"eth_address" gorm:"column:eth_address;"`
	Nonce         int    `json:"nonce" gorm:"column:nonce;"`
	Email         string `json:"-" gorm:"column:email;"`
}

type AuthVerifyData struct {
	AuthData
	Signature string `json:"signature"`
}

func (AuthData) TableName() string { return TbName }

type AuthDataCreation struct {
	common.SQLModel
	WalletAddress string `json:"eth_address" gorm:"column:eth_address;"`
	Nonce         int    `json:"nonce" gorm:"column:nonce;"`
	LastName      string `json:"last_name" gorm:"column:last_name;"`
	FirstName     string `json:"first_name" gorm:"column:first_name;"`
	IsSetPassword bool   `json:"-" gorm:"column:is_set_password;"`
}

func (AuthDataCreation) TableName() string { return TbName }

func (data *AuthDataCreation) PrepareForCreating() {
	data.SQLModel = common.NewSQLModel()
	data.WalletAddress = strings.ToLower(strings.TrimSpace(data.WalletAddress))
}

func (data *AuthVerifyData) Validate() error {
	data.WalletAddress = strings.ToLower(strings.TrimSpace(data.WalletAddress))
	data.Signature = strings.ToLower(strings.TrimSpace(data.Signature))

	if data.WalletAddress == "" || !common2.IsHexAddress(data.WalletAddress) {
		return ErrWalletAddressInvalid
	}

	if data.Signature == "" {
		return ErrSignatureInvalid
	}

	if data.Nonce <= 0 {
		return ErrNonceMustBeGreaterThanZero
	}

	return nil
}

type AuthDataUpdating struct {
	Nonce                *int       `json:"-" gorm:"column:nonce;"`
	DisplayName          *string    `json:"-" gorm:"column:name;"`
	WalletAddress        *string    `json:"eth_address" gorm:"column:eth_address;"`
	EthAddressVerified   *bool      `json:"-" gorm:"column:eth_address_verified;"`
	EthAddressVerifiedAt *time.Time `json:"-" gorm:"column:eth_address_verified_at;"`
	Active               *int       `json:"-" gorm:"column:active;"`
}

func (AuthDataUpdating) TableName() string { return TbName }
