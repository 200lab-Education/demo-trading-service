package authbiz

import (
	"context"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"math/rand"
	"time"
	"trading-service/common"
	authmodel "trading-service/modules/auth/model"
	"trading-service/plugin/tokenprovider"
)

type VerifySignatureStore interface {
	FindUserWithCondition(ctx context.Context, condition map[string]interface{}) (*authmodel.AuthData, error)
	CreateUser(ctx context.Context, data *authmodel.AuthDataCreation) error
	UpdateUserWithCondition(ctx context.Context, data *authmodel.AuthDataUpdating, condition map[string]interface{}) error
}

type verifySignatureBiz struct {
	store         VerifySignatureStore
	tokenProvider tokenprovider.Provider
}

func NewVerifySignatureBiz(store VerifySignatureStore, tokenProvider tokenprovider.Provider) *verifySignatureBiz {
	return &verifySignatureBiz{store: store, tokenProvider: tokenProvider}
}

func (biz *verifySignatureBiz) VerifySignature(
	ctx context.Context,
	data *authmodel.AuthVerifyData,
) (tokenprovider.Token, error) {
	if err := data.Validate(); err != nil {
		return nil, common.ErrInvalidRequest(err)
	}

	result, err := biz.store.FindUserWithCondition(
		ctx, map[string]interface{}{"eth_address": data.WalletAddress, "nonce": data.Nonce},
	)

	if err != nil {
		return nil, common.ErrCannotGetEntity(authmodel.EntityName, err)
	}

	if !verifySig(
		data.WalletAddress,
		data.Signature,
		[]byte(fmt.Sprintf("%d", data.Nonce)),
	) {
		return nil, authmodel.ErrSignatureInvalid
	}

	// change nonce and update status
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	nonce := r1.Intn(9999) + 10000

	isVerified := true
	now := time.Now().UTC()
	active := 1

	if err := biz.store.UpdateUserWithCondition(ctx, &authmodel.AuthDataUpdating{
		Nonce:                &nonce,
		EthAddressVerified:   &isVerified,
		EthAddressVerifiedAt: &now,
		Active:               &active,
	}, map[string]interface{}{"id": result.Id}); err != nil {
		return nil, common.ErrCannotUpdateEntity("user", err)
	}

	payload := &common.TokenPayload{
		UId: result.Id,
	}

	accessToken, err := biz.tokenProvider.Generate(payload, 60*60*24*30)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	return accessToken, nil
}

func verifySig(from, sigHex string, msg []byte) bool {
	fromAddr := ethcommon.HexToAddress(from)

	sig := hexutil.MustDecode(sigHex)
	// https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L442
	if sig[64] != 27 && sig[64] != 28 {
		return false
	}
	sig[64] -= 27

	pubKey, err := crypto.SigToPub(signHash(msg), sig)
	if err != nil {
		return false
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	return fromAddr == recoveredAddr
}

// https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L404
// signHash is a helper function that calculates a hash for the given message that can be
// safely used to calculate a signature from.
//
// The hash is calculated as
//
//	keccak256("\x19Ethereum Signed Message:\n"${message length}${message}).
//
// This gives context to the signed message and prevents signing of transactions.
func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}
