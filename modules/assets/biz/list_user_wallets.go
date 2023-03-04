package assetbiz

import (
	"context"
	"trading-service/common"
	assetmodel "trading-service/modules/assets/model"
	walletmodel "trading-service/modules/wallet/model"
)

type ListUserWalletStore interface {
	ListUserWallets(ctx context.Context, filter *assetmodel.Filter, moreInfo ...string) ([]assetmodel.UserAsset, error)
}

type WalletStore interface {
	WalletMap(ctx context.Context) (map[int]walletmodel.Wallet, error)
}

type listUserWalletBiz struct {
	requester common.Requester
	store     ListUserWalletStore
	wStore    WalletStore
}

func NewListUserWalletBiz(requester common.Requester, store ListUserWalletStore, wStore WalletStore) *listUserWalletBiz {
	return &listUserWalletBiz{
		requester: requester,
		store:     store,
		wStore:    wStore,
	}
}

func (biz *listUserWalletBiz) ListUserWallets(ctx context.Context, filter *assetmodel.Filter) ([]assetmodel.UserAsset, error) {
	result, err := biz.store.ListUserWallets(ctx, filter, "Asset")

	if err != nil {
		return nil, assetmodel.ErrCannotListUserWallet(err)
	}

	walletMap, err := biz.wStore.WalletMap(ctx)

	if err != nil {
		return nil, assetmodel.ErrCannotListUserWallet(err)
	}

	for i, item := range result {
		wallet := walletMap[item.WalletId]
		result[i].Wallet = &wallet
	}

	return result, nil
}
