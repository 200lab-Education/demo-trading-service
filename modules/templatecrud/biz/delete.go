package biz

import (
	"context"
	"trading-service/common"
	"trading-service/modules/banking/model"
)

type DeleteStore interface {
	Find(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.BankAccount, error)
	Delete(ctx context.Context, condition map[string]interface{}) error
}

type deleteBiz struct {
	store     DeleteStore
	requester common.Requester
}

func NewDeleteBiz(store DeleteStore, requester common.Requester) *deleteBiz {
	return &deleteBiz{store: store, requester: requester}
}

func (biz deleteBiz) DeleteBankAccount(ctx context.Context, id int) error {
	oldData, err := biz.store.Find(ctx, map[string]interface{}{"id": id})

	if err != nil {
		return common.ErrCannotGetEntity(model.EntityName, err)
	}

	if oldData.UserId != biz.requester.GetUserId() {
		return common.ErrNoPermission(nil)
	}

	if err := biz.store.Delete(ctx, map[string]interface{}{"id": id}); err != nil {
		return common.ErrCannotDeleteEntity(model.EntityName, err)
	}

	return nil
}
