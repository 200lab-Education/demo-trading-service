package usermodel

type WithdrawParams struct {
	PaymentToken *string `json:"payment_token"`
	Amount       *int64  `json:"amount"`
}
