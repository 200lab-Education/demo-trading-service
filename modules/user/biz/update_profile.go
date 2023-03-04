package userbiz

import (
	"context"
	"trading-service/common"
	usermodel "trading-service/modules/user/model"
)

type UpdateUserProfileStore interface {
	UpdateUser(ctx context.Context, condition map[string]interface{}, data *usermodel.UserUpdate) error
}

type updateUserProfileBiz struct {
	store     UpdateUserProfileStore
	requester common.Requester
}

func NewUpdateUserProfileBiz(store UpdateUserProfileStore, requester common.Requester) *updateUserProfileBiz {
	return &updateUserProfileBiz{store: store, requester: requester}
}

func (biz *updateUserProfileBiz) UpdateProfile(
	ctx context.Context,
	data *usermodel.UserUpdate,
) error {
	if err := biz.store.UpdateUser(ctx, map[string]interface{}{"id": biz.requester.GetUserId()}, data); err != nil {
		return common.ErrCannotUpdateEntity(usermodel.EntityName, err)
	}

	return nil
}

func (biz *updateUserProfileBiz) UpdateUser(
	ctx context.Context,
	id int,
	data *usermodel.UserUpdate,
) error {
	if err := biz.store.UpdateUser(ctx, map[string]interface{}{"id": id}, data); err != nil {
		return common.ErrCannotUpdateEntity(usermodel.EntityName, err)
	}

	return nil
}
