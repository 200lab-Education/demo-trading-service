package biz

import (
	"context"
	"trading-service/common"
	"trading-service/modules/fiattx/model"
)

type ListStore interface {
	List(ctx context.Context, filter *model.Filter, paging *common.Paging, moreInfo ...string) ([]model.FiatDW, error)
}

type listBiz struct {
	store ListStore
}

func NewListBiz(store ListStore) *listBiz {
	return &listBiz{store: store}
}

func (biz *listBiz) ListTxs(
	ctx context.Context,
	filter *model.Filter,
	paging *common.Paging,
) ([]model.FiatDW, error) {
	result, err := biz.store.List(ctx, filter, paging, "Asset", "BankAccount", "User")

	if err != nil {
		return nil, common.ErrCannotListEntity(model.EntityName, err)
	}

	return result, nil
}
