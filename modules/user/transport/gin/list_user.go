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

func ListUsers(sc goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var filter usermodel.Filter

		if err := c.ShouldBind(&filter); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		if err := paging.Process(); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		//requester := c.MustGet(common.CurrentUser).(common.Requester)
		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)

		store := userstorage.NewSQLStore(db)
		biz := userbiz.NewListUserBiz(store)

		result, err := biz.ListUsers(c.Request.Context(), &filter, &paging)

		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask(common.DbTypeUser)

			if i == len(result)-1 {
				paging.NextCursor = result[i].FakeId.String()
			}
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, filter))
	}
}
