package biz

import (
	"context"
	"trading-service/common"
	"trading-service/modules/p2p/model"
)

type GetStore interface {
	Find(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.P2pOrder, error)
	FindTrade(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.P2pTrading, error)
}

type getBiz struct {
	store     GetStore
	requester common.Requester
}

func NewGetBiz(store GetStore, requester common.Requester) *getBiz {
	return &getBiz{store: store, requester: requester}
}

func (biz getBiz) GetOrder(ctx context.Context, id int) (*model.P2pOrder, error) {
	data, err := biz.store.Find(ctx, map[string]interface{}{"id": id}, "OfferAsset", "PayableAsset", "BankAccount", "User")

	if err != nil {
		return nil, common.ErrCannotGetEntity(model.OrderEntityName, err)
	}

	mine := data.UserId == biz.requester.GetUserId()
	if !common.IsAdmin(biz.requester) && !mine && data.Status != model.OrdStActive.String() {
		return nil, common.ErrNoPermission(nil)
	}

	return data, nil
}

func (biz getBiz) GetTrade(ctx context.Context, id int) (*model.P2pTrading, error) {
	data, err := biz.store.FindTrade(ctx, map[string]interface{}{"id": id}, "OfferAsset", "PayableAsset", "BankAccount", "OrderOwner", "User", "Logs", "Logs.User")

	if err != nil {
		return nil, common.ErrCannotGetEntity(model.TradeEntityName, err)
	}

	mine := data.UserId == biz.requester.GetUserId() || data.OrderUserId == biz.requester.GetUserId()
	if !common.IsAdmin(biz.requester) && !mine {
		return nil, common.ErrNoPermission(nil)
	}

	return data, nil
}
