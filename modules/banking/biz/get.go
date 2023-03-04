package biz

import (
	"context"
	"trading-service/common"
	"trading-service/modules/banking/model"
)

type GetStore interface {
	Find(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.BankAccount, error)
}

type getBiz struct {
	store GetStore
}

func NewGetBiz(store GetStore) *getBiz {
	return &getBiz{store: store}
}

func (biz getBiz) GetBankAccount(ctx context.Context, id int) (*model.BankAccount, error) {
	data, err := biz.store.Find(ctx, map[string]interface{}{"id": id})

	if err != nil {
		return nil, common.ErrCannotGetEntity(model.EntityName, err)
	}

	return data, nil
}
