package usergin

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"trading-service/common"
	userbiz "trading-service/modules/user/biz"
	userstorage "trading-service/modules/user/storage"
)

func GetUser(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := common.FromBase58(c.Param("user_id"))

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)

		store := userstorage.NewSQLStore(db)
		biz := userbiz.NewGetUserBiz(store)

		data, err := biz.GetUser(c.Request.Context(), int(userId.GetLocalID()))

		if err != nil {
			panic(err)
		}

		data.Mask(common.DbTypeUser)

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}
