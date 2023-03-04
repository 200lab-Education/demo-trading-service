package common

import (
	"fmt"
	"time"
)

const (
	DbTypeUser        = 1
	DbTypeBSCTx       = 2
	DbTypeAsset       = 3
	DbTypeWallet      = 4
	DbTypeCity        = 4
	DbTypeCountry     = 5
	DbTypeBankAcc     = 6
	DbTypeFiatTx      = 7
	DbTypeFiatTxLog   = 8
	DbTypeP2pOrder    = 8
	DbTypeP2pTrade    = 8
	DbTypeP2pTradeLog = 9

	CurrentUser = "user"

	RoleAdmin = "admin"
	RoleUser  = "user"
	RoleMod   = "mod"

	WalletSPOT = 1
)

const DateTimeFmt = "2006-01-02 15:04:05.999999"

type Requester interface {
	GetUserId() int
	GetRole() string
	GetWalletAddress() string
}

const (
	KeyUserBalanceLockKey = "%d-lock-balance"
	keyUserAssetWallet    = "user-%d-asset-%d-wallet-%d"
	KeyUserWallet         = "user-%d-wallet-%d"
	StatusPending         = "pending"
	TimeLock              = 10 * time.Minute
)

const (
	GasLimit = uint64(124014)

	PluginMainDB    = "mysql"
	PluginJWT       = "jwt"
	PluginMutexLock = "locker"
	PluginBSCEx     = "bsc-exchange"

	MasterTxData = "MasterTxData"
)

func GetUserAssetWalletKey(userId, assetId, walletId int) string {
	return fmt.Sprintf(keyUserAssetWallet, userId, assetId, walletId)
}

type TokenPayload struct {
	UId   int    `json:"user_id"`
	URole string `json:"role"`
}

func (p TokenPayload) UserId() int {
	return p.UId
}

func (p TokenPayload) Role() string {
	return p.URole
}
