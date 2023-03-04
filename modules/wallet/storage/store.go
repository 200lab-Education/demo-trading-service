package walletstorage

import (
	"context"
	"trading-service/common"
	walletmodel "trading-service/modules/wallet/model"
)

type localStore struct {
}

func NewLocalStore() *localStore {
	return &localStore{}
}

func (s *localStore) ListDataWithCondition(ctx context.Context) ([]walletmodel.Wallet, error) {
	return []walletmodel.Wallet{
		walletmodel.NewWallet(1, "Spot"),
		walletmodel.NewWallet(2, "Funding"),
		walletmodel.NewWallet(3, "Investment"),
		walletmodel.NewWallet(4, "Interest"),
		walletmodel.NewWallet(5, "Trading"),
	}, nil
}

func (s *localStore) WalletMap(ctx context.Context) (map[int]walletmodel.Wallet, error) {
	wallets, err := s.ListDataWithCondition(ctx)

	if err != nil {
		return nil, common.ErrDB(err)
	}

	result := make(map[int]walletmodel.Wallet)

	for i, item := range wallets {
		result[item.Id] = wallets[i]
	}

	return result, nil
}
