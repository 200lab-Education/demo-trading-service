package biz

import (
	"context"
	"trading-service/common"
	"trading-service/modules/banking/model"
)

type ListStore interface {
	List(ctx context.Context, filter *model.Filter, paging *common.Paging, moreInfo ...string) ([]model.BankAccount, error)
}

type listBiz struct {
	store ListStore
}

func NewListBiz(store ListStore) *listBiz {
	return &listBiz{store: store}
}

func (biz *listBiz) ListBankAccounts(
	ctx context.Context,
	filter *model.Filter,
	paging *common.Paging,
) ([]model.BankAccount, error) {
	result, err := biz.store.List(ctx, filter, paging)

	if err != nil {
		return nil, common.ErrCannotListEntity(model.EntityName, err)
	}

	return result, nil
}
