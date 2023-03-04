package bscex

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type BSCEx interface {
	Withdraw(ctx context.Context, receiver common.Address, paymentToken common.Address, amount *big.Int) (string, error)
	ChainID() int
	ExchangeClient() *ethclient.Client
	ABI() abi.ABI
	CurrentBlock() int64
	ExSCAddress() string
	IsMainNet() bool
}
