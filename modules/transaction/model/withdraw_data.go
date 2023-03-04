package txmodel

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"strings"
)

var (
	ErrWalletAddressInvalid        = errors.New("wallet address is invalid")
	ErrWalletAddressHasBinded      = errors.New("wallet address has connected to another wallet")
	ErrAmountMustBeGreaterThanZero = errors.New("amount must be greater than zero")
)

type WithdrawData struct {
	//WalletId int                 `json:"-"`
	//AssetId  int                 `json:"-"`
	Receiver common.Address      `json:"receiver_address"`
	Amount   decimal.NullDecimal `json:"amount"`
	Type     string              `json:"type" form:"type"`
}

func (data *WithdrawData) Validate() error {
	data.Receiver = common.HexToAddress(strings.ToLower(strings.TrimSpace(data.Receiver.Hex())))

	if data.Receiver.Hex() == "" || !common.IsHexAddress(data.Receiver.Hex()) {
		return ErrWalletAddressInvalid
	}

	if data.Amount.Decimal.LessThan(decimal.NewFromInt(1)) {
		return ErrAmountMustBeGreaterThanZero
	}

	return nil
}
