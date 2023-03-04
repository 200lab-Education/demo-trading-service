package userbiz

import (
	"context"
	"trading-service/common"
	usermodel "trading-service/modules/user/model"
)

type ListUserStore interface {
	ListUsers(ctx context.Context, filter *usermodel.Filter, paging *common.Paging, moreInfo ...string) ([]usermodel.User, error)
}

type listUserBiz struct {
	store ListUserStore
}

func NewListUserBiz(store ListUserStore) *listUserBiz {
	return &listUserBiz{store: store}
}

func (biz *listUserBiz) ListUsers(
	ctx context.Context,
	filter *usermodel.Filter,
	paging *common.Paging,
) ([]usermodel.User, error) {
	result, err := biz.store.ListUsers(ctx, filter, paging)

	if err != nil {
		return nil, common.ErrCannotListEntity(usermodel.EntityName, err)
	}

	return result, nil
}
