package model

import (
	"strings"
	"trading-service/common"
)

const EntityName = "Bank Account"

var (
	ErrTitleCannotBeBlank         = common.ValidateError("title cannot be blank", "ErrTitleCannotBeBlank")
	ErrAccountNumberCannotBeBlank = common.ValidateError("account number cannot be blank", "ErrAccountNumberCannotBeBlank")
	ErrAccountNameCannotBeBlank   = common.ValidateError("account name cannot be blank", "ErrAccountNameCannotBeBlank")
	ErrBankNameCannotBeBlank      = common.ValidateError("bank name cannot be blank", "ErrBankNameCannotBeBlank")
)

type BankAccount struct {
	common.SQLModel
	Title         string `json:"title" gorm:"column:title;"`
	BankName      string `json:"bank_name" gorm:"column:bank_name;"`
	AccountNumber string `json:"account_number" gorm:"column:account_number;"`
	AccountName   string `json:"account_name" gorm:"column:account_name;"`
	UserId        int    `json:"-" gorm:"column:user_id;"`
}

func (BankAccount) TableName() string {
	return "user_bank_accounts"
}

type BankAccountCreate struct {
	common.SQLModel
	Title         string `json:"title" gorm:"column:title;"`
	BankName      string `json:"bank_name" gorm:"column:bank_name;"`
	AccountNumber string `json:"account_number" gorm:"column:account_number;"`
	AccountName   string `json:"account_name" gorm:"column:account_name;"`
	UserId        int    `json:"-" gorm:"column:user_id;"`
}

func (BankAccountCreate) TableName() string {
	return BankAccount{}.TableName()
}

func (data *BankAccountCreate) Validate() error {
	data.Title = strings.TrimSpace(data.Title)
	data.BankName = strings.TrimSpace(data.BankName)
	data.AccountNumber = strings.TrimSpace(data.AccountNumber)
	data.AccountName = strings.TrimSpace(data.AccountName)

	if data.Title == "" {
		return ErrTitleCannotBeBlank
	}

	if data.BankName == "" {
		return ErrBankNameCannotBeBlank
	}

	if data.AccountNumber == "" {
		return ErrAccountNumberCannotBeBlank
	}

	if data.AccountName == "" {
		return ErrAccountNameCannotBeBlank
	}

	return nil
}

type BankAccountUpdate struct {
	Title         *string `json:"title" gorm:"column:title;"`
	BankName      *string `json:"bank_name" gorm:"column:bank_name;"`
	AccountNumber *string `json:"account_number" gorm:"column:account_number;"`
	AccountName   *string `json:"account_name" gorm:"column:account_name;"`
}

func (BankAccountUpdate) TableName() string {
	return BankAccount{}.TableName()
}
