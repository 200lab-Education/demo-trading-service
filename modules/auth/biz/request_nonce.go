package authbiz

import (
	"context"
	"fmt"
	common2 "github.com/ethereum/go-ethereum/common"
	"math/rand"
	"strings"
	"time"
	"trading-service/common"
	authmodel "trading-service/modules/auth/model"
)

type RequestNonceStore interface {
	FindUserWithCondition(
		ctx context.Context,
		condition map[string]interface{},
	) (*authmodel.AuthData, error)
	CreateUser(ctx context.Context, data *authmodel.AuthDataCreation) error
	UpdateUserWithCondition(
		ctx context.Context,
		data *authmodel.AuthDataUpdating,
		condition map[string]interface{},
	) error
}

type requestNonceItemBiz struct {
	store RequestNonceStore
}

func NewRequestNonceBiz(store RequestNonceStore) *requestNonceItemBiz {
	return &requestNonceItemBiz{store: store}
}

func (biz *requestNonceItemBiz) RequestNonce(
	ctx context.Context,
	walletAddress string,
) (*authmodel.AuthData, error) {
	walletAddress = strings.ToLower(strings.TrimSpace(walletAddress))

	if walletAddress == "" {
		return nil, authmodel.ErrWalletAddressInvalid
	}

	result, err := biz.store.FindUserWithCondition(ctx, map[string]interface{}{"eth_address": walletAddress})

	if err == common.ErrRecordNotFound {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		nonce := r1.Intn(9999) + 10000

		newData := authmodel.AuthDataCreation{
			WalletAddress: walletAddress,
			Nonce:         nonce,
			FirstName:     "user",
			LastName:      fmt.Sprintf("%d", time.Now().UTC().UnixNano()),
		}

		newData.PrepareForCreating()

		if err := biz.store.CreateUser(ctx, &newData); err != nil {
			return nil, common.ErrCannotCreateEntity(authmodel.EntityName, err)
		}

		result = &authmodel.AuthData{
			WalletAddress: walletAddress,
			Nonce:         nonce,
		}

		return result, nil
	}

	if err != nil {
		return nil, common.ErrCannotGetEntity(authmodel.EntityName, err)
	}

	return result, nil
}

func (biz *requestNonceItemBiz) RequestNonceOldUser(
	ctx context.Context,
	requester common.Requester,
	walletAddress string,
) (*authmodel.AuthData, error) {
	walletAddress = strings.ToLower(strings.TrimSpace(walletAddress))

	if walletAddress == "" || !common2.IsHexAddress(walletAddress) {
		return nil, authmodel.ErrWalletAddressInvalid
	}

	result, err := biz.store.FindUserWithCondition(ctx, map[string]interface{}{"eth_address": walletAddress})

	if err != nil && err != common.ErrRecordNotFound {
		return nil, common.ErrCannotGetEntity(authmodel.EntityName, err)
	}

	if result.Id != requester.GetUserId() {
		return nil, authmodel.ErrWalletAddressHasBinded
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	nonce := r1.Intn(9999) + 10000
	verified := false

	data := authmodel.AuthDataUpdating{
		WalletAddress:      &walletAddress,
		Nonce:              &nonce,
		EthAddressVerified: &verified,
	}

	if err := biz.store.UpdateUserWithCondition(ctx, &data, map[string]interface{}{"id": requester.GetUserId()}); err != nil {
		return nil, common.ErrCannotUpdateEntity(authmodel.EntityName, err)
	}

	result = &authmodel.AuthData{
		WalletAddress: walletAddress,
		Nonce:         nonce,
	}

	return result, nil
}
