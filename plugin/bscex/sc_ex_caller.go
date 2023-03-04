package bscex

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"flag"
	"github.com/200Lab-Education/go-sdk/logger"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"os"
	"strconv"
	"time"
	common2 "trading-service/common"
	"trading-service/component/asyncjob"
)

const (
	BlockStart         int64 = 16061292
	DefaultRPCURL            = "https://data-seed-prebsc-1-s1.binance.org:8545"
	ContractAddress          = "0x28Fb673E647EE8485a2B7683833034582CdA41F9"
	DefaultHotWalletPK       = "f03730fde4b1671da17b4113932d31400e451bebbae3a3f6e2d5bc245131b3a1" // testnet
	DefaultABIFilePath       = "./withdraw_deposit.abi"
	DefaultNetworkType       = "test"
)

var ErrorCastingPublicKey = errors.New("error casting public key to ECDSA")

type scConfig struct {
	rpcURL              string
	exSCAddress         string
	hotWalletPrivateKey string
	currentBlock        int64
	networkType         string
}

type scExchangeCaller struct {
	name    string
	abiObj  abi.ABI
	client  *ethclient.Client
	config  scConfig
	logger  logger.Logger
	chainId int
}

func NewSCCaller(name string) *scExchangeCaller {
	return &scExchangeCaller{
		name:   name,
		config: scConfig{},
	}
}

func (t *scExchangeCaller) GetPrefix() string {
	return t.name
}

func (t *scExchangeCaller) Get() interface{} {
	return t
}

func (t *scExchangeCaller) Name() string {
	return t.name
}

func (t *scExchangeCaller) InitFlags() {
	flag.Int64Var(&t.config.currentBlock, "bsc-block-number", BlockStart, "Current block number to start the scanner")
	flag.StringVar(&t.config.rpcURL, "bsc-rpc-url", DefaultRPCURL, "RPC URL of BSC Node")
	flag.StringVar(&t.config.hotWalletPrivateKey, "bsc-hot-wallet-pk", DefaultHotWalletPK, "Hot wallet private key")
	flag.StringVar(&t.config.exSCAddress, "bsc-exchange-smart-contract", ContractAddress, "Address of exchange smart contract")
	flag.StringVar(&t.config.networkType, "bsc-network-type", DefaultNetworkType, "main or test net. Default: test")
}

func (t *scExchangeCaller) Configure() error {
	t.logger = logger.GetCurrent().GetLogger(t.name)

	client, err := ethclient.Dial(t.config.rpcURL)
	if err != nil {
		t.logger.Errorln(err.Error())
		return err
	}

	t.client = client

	chainId, err := t.getChainId()

	if err != nil {
		t.logger.Errorln(err.Error())
		return err
	}

	t.chainId = chainId

	f, err := os.Open(DefaultABIFilePath)

	mkpABI, err := abi.JSON(f)

	if err != nil {
		t.logger.Errorln(err.Error())
		return err
	}

	t.abiObj = mkpABI

	return nil
}

func (t *scExchangeCaller) Run() error {
	return t.Configure()
}

func (t *scExchangeCaller) Stop() <-chan bool {
	c := make(chan bool)
	go func() {
		t.client.Close()
		c <- true
	}()
	return c
}

func (t *scExchangeCaller) Withdraw(ctx context.Context, receiver common.Address, paymentToken common.Address, amount *big.Int) (string, error) {
	data, err := t.abiObj.Pack("withdrawToken", receiver, paymentToken, amount)

	//

	privateKey, err := crypto.HexToECDSA(t.config.hotWalletPrivateKey)
	if err != nil {
		return "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", ErrorCastingPublicKey
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := t.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}

	value := new(big.Int)
	value.SetString("0", 10)
	gasLimit := common2.GasLimit

	gasPrice, err := t.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	to := common.HexToAddress(t.config.exSCAddress)

	newLegacyTx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &to,
		Value:    value,
		Data:     data,
	})

	chainID, err := t.client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}

	signedTx, err := types.SignTx(newLegacyTx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}

	err = t.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	//txHash := signedTx.Hash().Hex()

	return signedTx.Hash().Hex(), err
}

func (t *scExchangeCaller) getChainId() (int, error) {
	var result int = 1

	job := asyncjob.NewJob(func(ctx context.Context) error {
		chainId, err := t.client.ChainID(context.Background())

		if err != nil {
			return err
		}

		rs, err := strconv.ParseInt(chainId.String(), 10, 64)

		if err != nil {
			return err
		}

		result = int(rs)

		return nil
	})

	job.SetRetryDurations(time.Second, time.Second*2, time.Second*3)

	if err := asyncjob.NewGroup(false, job).Run(context.Background()); err != nil {
		return 0, err
	}

	return result, nil
}

func (t *scExchangeCaller) ChainID() int                      { return t.chainId }
func (t *scExchangeCaller) ExchangeClient() *ethclient.Client { return t.client }
func (t *scExchangeCaller) ABI() abi.ABI                      { return t.abiObj }
func (t *scExchangeCaller) CurrentBlock() int64               { return t.config.currentBlock }
func (t *scExchangeCaller) ExSCAddress() string               { return t.config.exSCAddress }
func (t *scExchangeCaller) IsMainNet() bool                   { return t.config.networkType == "main" }
