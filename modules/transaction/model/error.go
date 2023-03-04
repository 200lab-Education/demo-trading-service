package txmodel

import (
	"fmt"
	"trading-service/common"
)

var (
	ErrNotEnoughAsset = common.NewCustomError(nil, "not enough asset", "ErrNotEnoughAsset")
	ErrInvalidAsset   = common.NewCustomError(nil, "invalid asset", "ErrInvalidAsset")
)

func ErrCannotWithdraw(err error) *common.AppError {
	return common.NewCustomError(
		err,
		fmt.Sprintf("cannot withdraw asset"),
		fmt.Sprintf("ErrCannotWithdrawAsset"),
	)
}
