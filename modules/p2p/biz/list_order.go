package biz

import (
	"context"
	"trading-service/common"
	"trading-service/modules/p2p/model"
)

type ListStore interface {
	List(ctx context.Context, filter *model.Filter, paging *common.Paging, moreInfo ...string) ([]model.P2pOrder, error)
	ListTrades(ctx context.Context, filter *model.Filter, paging *common.Paging, moreInfo ...string) ([]model.P2pTrading, error)
	Find(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.P2pOrder, error)
}

type listBiz struct {
	store     ListStore
	requester common.Requester
}

func NewListBiz(store ListStore, requester common.Requester) *listBiz {
	return &listBiz{store: store, requester: requester}
}

func (biz *listBiz) ListOrder(
	ctx context.Context,
	filter *model.Filter,
	paging *common.Paging,
) ([]model.P2pOrder, error) {
	result, err := biz.store.List(ctx, filter, paging, "OfferAsset", "PayableAsset", "User")

	if err != nil {
		return nil, common.ErrCannotListEntity(model.OrderEntityName, err)
	}

	return result, nil
}

func (biz *listBiz) ListTrade(
	ctx context.Context,
	filter *model.Filter,
	paging *common.Paging,
) ([]model.P2pTrading, error) {
	if filter.OrderId > 0 {
		data, err := biz.store.Find(ctx, map[string]interface{}{"id": filter.OrderId})

		if err != nil {
			return nil, common.ErrCannotGetEntity(model.OrderEntityName, err)
		}

		mine := data.UserId == biz.requester.GetUserId()

		if !common.IsAdmin(biz.requester) && !mine {
			return nil, common.ErrNoPermission(nil)
		}
	}

	result, err := biz.store.ListTrades(ctx, filter, paging, "OfferAsset", "PayableAsset", "User", "OrderOwner")

	if err != nil {
		return nil, common.ErrCannotListEntity(model.TradeEntityName, err)
	}

	return result, nil
}
