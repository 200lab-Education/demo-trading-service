package biz

import (
	"context"
	"trading-service/common"
	"trading-service/modules/fiattx/model"
)

type GetStore interface {
	Find(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.FiatDW, error)
}

type getBiz struct {
	store     GetStore
	requester common.Requester
}

func NewGetBiz(store GetStore, requester common.Requester) *getBiz {
	return &getBiz{store: store, requester: requester}
}

func (biz getBiz) GetTx(ctx context.Context, id int) (*model.FiatDW, error) {
	data, err := biz.store.Find(ctx, map[string]interface{}{"id": id}, "Asset", "BankAccount", "User", "Logs", "Logs.User")

	if err != nil {
		return nil, common.ErrCannotGetEntity(model.EntityName, err)
	}

	if !common.IsAdmin(biz.requester) && data.UserId != biz.requester.GetUserId() {
		return nil, common.ErrNoPermission(nil)
	}

	return data, nil
}
