package biz

import (
	"context"
	"trading-service/common"
	"trading-service/modules/banking/model"
)

type UpdateStore interface {
	Find(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.BankAccount, error)
	Update(ctx context.Context, condition map[string]interface{}, data *model.BankAccountUpdate) error
}

type updateBiz struct {
	store     UpdateStore
	requester common.Requester
}

func NewUpdateBiz(store UpdateStore, requester common.Requester) *updateBiz {
	return &updateBiz{store: store, requester: requester}
}

func (biz updateBiz) UpdateBankAccount(ctx context.Context, id int, data *model.BankAccountUpdate) error {
	oldData, err := biz.store.Find(ctx, map[string]interface{}{"id": id})

	if err != nil {
		return common.ErrCannotGetEntity(model.EntityName, err)
	}

	if oldData.UserId != biz.requester.GetUserId() {
		return common.ErrNoPermission(nil)
	}

	if err := biz.store.Update(ctx, map[string]interface{}{"id": id}, data); err != nil {
		return common.ErrCannotUpdateEntity(model.EntityName, err)
	}

	return nil
}
