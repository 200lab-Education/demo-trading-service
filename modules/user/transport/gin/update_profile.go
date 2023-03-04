package usergin

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"trading-service/common"
	userbiz "trading-service/modules/user/biz"
	usermodel "trading-service/modules/user/model"
	userstorage "trading-service/modules/user/storage"
)

func UpdateProfile(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)

		requester := c.MustGet(common.CurrentUser).(common.Requester)

		var data usermodel.UserUpdate

		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := userstorage.NewSQLStore(db)
		biz := userbiz.NewUpdateUserProfileBiz(store, requester)

		if err := biz.UpdateProfile(c.Request.Context(), &data); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
