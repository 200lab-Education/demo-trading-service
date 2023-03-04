package assetmodel

import (
	"fmt"
	"trading-service/common"
)

func ErrCannotListUserWallet(err error) *common.AppError {
	return common.NewCustomError(
		err,
		fmt.Sprintf("Cannot list user wallets"),
		fmt.Sprintf("ErrCannotListUserWallet"),
	)
}
