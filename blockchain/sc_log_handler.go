package blockchain

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
	"trading-service/component/asyncjob"
	assetmodel "trading-service/modules/assets/model"
	authmodel "trading-service/modules/auth/model"
	txmodel "trading-service/modules/transaction/model"
	walletmodel "trading-service/modules/wallet/model"
	"trading-service/plugin/bscex"
)

const (
	EvtMoneyIn  = "0xff4970bbc238020f1e1ca43cf8ca919450c525dada25f3ef91f1baf7c4d5fbe7"
	EvtMoneyOut = "0x6cf264d20c0889cb2e2f1f620c16e43005791c83fb04a6cccc0b3173c18991db"
)

type TxStorage interface {
	GetDataWithCondition(
		ctx context.Context,
		condition map[string]interface{},
		moreKeys ...string,
	) (*txmodel.BSCTransaction, error)

	CreateData(
		ctx context.Context,
		tx *txmodel.BSCTransaction,
	) error

	Update(ctx context.Context, condition map[string]interface{}, data *txmodel.BSCTransactionUpdate) error
}

type UserStore interface {
	CreateUser(
		ctx context.Context,
		data *authmodel.AuthDataCreation,
	) error
	FindUserWithCondition(
		ctx context.Context,
		condition map[string]interface{},
	) (*authmodel.AuthData, error)
}

type WalletStore interface {
	ListDataWithCondition(ctx context.Context) ([]walletmodel.Wallet, error)
	WalletMap(ctx context.Context) (map[int]walletmodel.Wallet, error)
}

type AssetStore interface {
	GetDataWithCondition(
		ctx context.Context,
		condition map[string]interface{},
		moreKeys ...string,
	) (*assetmodel.Asset, error)

	CreateData(
		ctx context.Context,
		data *assetmodel.Asset,
	) error

	GetAssetInWallet(ctx context.Context, userId, assetId, walletId int) (*assetmodel.UserAsset, error)
	CreateAssetInWallet(ctx context.Context, userId, assetId, walletId int, amount decimal.NullDecimal) error
	IncreaseAmountAssetInWallet(ctx context.Context, userId, assetId, walletId int, amount decimal.NullDecimal) error
	FinishLockUserAsset(ctx context.Context, lockId int) error
}

type scLogHdl struct {
	bscExStore  bscex.BSCEx
	txStore     TxStorage
	userStore   UserStore
	walletStore WalletStore
	assetStore  AssetStore
}

func NewSCLogHdl(
	bscExStore bscex.BSCEx,
	txStore TxStorage,
	userStore UserStore,
	walletStore WalletStore,
	assetStore AssetStore,
) *scLogHdl {
	return &scLogHdl{
		bscExStore:  bscExStore,
		txStore:     txStore,
		userStore:   userStore,
		walletStore: walletStore,
		assetStore:  assetStore,
	}
}

func (p *scLogHdl) Run(ctx context.Context, queue <-chan types.Log) {
	for l := range queue {
		functionHash := strings.ToLower(l.Topics[0].Hex())
		var job asyncjob.Job

		switch functionHash {
		case EvtMoneyIn:
			job = asyncjob.NewJob(func(ctx context.Context) error {
				return p.handleMoneyIn(ctx, l)
			})
		case EvtMoneyOut:
			job = asyncjob.NewJob(func(ctx context.Context) error {
				return p.handleMoneyOut(ctx, l)
			})

		default:
			log.Printf("Block %d - Tx %s - Function Hash %s \n", l.BlockNumber, l.TxHash.Hex(), functionHash)
			continue
		}

		job.SetRetryDurations(time.Second, time.Second, time.Second, time.Second) // 4 times (1s each)

		if err := job.Execute(ctx); err != nil {
			log.Errorln(err)
		}
	}
}

func (p *scLogHdl) findUser(ctx context.Context, walletAddress string) (*authmodel.AuthData, error) {
	user, err := p.userStore.FindUserWithCondition(ctx, map[string]interface{}{"eth_address": walletAddress})

	if err != nil {
		return nil, err
	}

	return user, nil

	//if err != appCommon.ErrRecordNotFound {
	//	return nil, err
	//}

	//if err == appCommon.ErrRecordNotFound {
	//	s1 := rand.NewSource(time.Now().UnixNano())
	//	r1 := rand.New(s1)
	//	nonce := r1.Intn(9999) + 10000
	//
	//	newUser := &authmodel.AuthDataCreation{
	//		WalletAddress: walletAddress,
	//		Nonce:         nonce,
	//	}
	//	newUser.PrepareForCreating()
	//	//newUser.Status = "not_verified"
	//
	//	if err := p.userStore.CreateUser(ctx, newUser); err != nil {
	//		return nil, err
	//	}
	//
	//	user = &authmodel.AuthData{
	//		Id:            newUser.Id,
	//		WalletAddress: walletAddress,
	//	}
	//}

	return user, nil
}

// BE: Todo service
// User authen (JWT)
// Comments, Likes todo
// Blockchain (web3)
