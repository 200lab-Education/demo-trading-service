package txbiz

import (
	"context"
	"trading-service/common"
	txmodel "trading-service/modules/transaction/model"
)

type ListTxStore interface {
	ListDataWithCondition(
		ctx context.Context,
		filter *txmodel.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]txmodel.BSCTransaction, error)
}

type listTxBiz struct {
	store ListTxStore
}

func NewListTxBiz(store ListTxStore) *listTxBiz {
	return &listTxBiz{store: store}
}

func (biz *listTxBiz) ListTx(
	ctx context.Context,
	filter *txmodel.Filter,
	paging *common.Paging,
) ([]txmodel.BSCTransaction, error) {
	result, err := biz.store.ListDataWithCondition(ctx, filter, paging)

	if err != nil {
		return nil, common.ErrCannotListEntity(txmodel.EntityName, err)
	}

	return result, nil
}
