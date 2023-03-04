package biz

import (
	"context"
	"trading-service/common"
	"trading-service/modules/banking/model"
)

type CreateStore interface {
	Create(ctx context.Context, data *model.BankAccountCreate) error
}

type createBiz struct {
	store CreateStore
}

func NewCreateBiz(store CreateStore) *createBiz {
	return &createBiz{store: store}
}

func (biz createBiz) CreateBankAccount(ctx context.Context, data *model.BankAccountCreate) error {
	if err := data.Validate(); err != nil {
		return common.ErrCannotCreateEntity(model.EntityName, err)
	}

	if err := biz.store.Create(ctx, data); err != nil {
		return common.ErrCannotCreateEntity(model.EntityName, err)
	}

	return nil
}
