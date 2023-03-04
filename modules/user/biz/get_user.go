package userbiz

import (
	"context"
	"trading-service/common"
	usermodel "trading-service/modules/user/model"
)

type GetUserStore interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*usermodel.User, error)
}

type getUserBiz struct {
	store GetUserStore
}

func NewGetUserBiz(store GetUserStore) *getUserBiz {
	return &getUserBiz{store: store}
}

func (biz *getUserBiz) GetUser(
	ctx context.Context,
	id int,
) (*usermodel.User, error) {
	data, err := biz.store.FindUser(ctx, map[string]interface{}{"id": id})

	if err != nil {
		return nil, common.ErrCannotGetEntity(usermodel.EntityName, err)
	}

	return data, nil
}
