package transport

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"trading-service/common"
	"trading-service/modules/p2p/biz"
	"trading-service/modules/p2p/model"
	"trading-service/modules/p2p/storage"
)

func ListP2pTrade(sc goservice.ServiceContext, userOnly bool) func(*gin.Context) {
	return func(c *gin.Context) {
		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)
		requester := c.MustGet(common.CurrentUser).(common.Requester)

		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		paging.Process()

		var filter model.Filter

		if err := c.ShouldBind(&filter); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		if v := c.Param("id"); v != "" {
			id, _ := common.FromBase58(v)
			filter.OrderId = int(id.GetLocalID())
		}

		if userOnly {
			filter.UserAccess = true
			filter.RequesterId = requester.GetUserId()
		}

		store := storage.NewSQLStore(db)
		business := biz.NewListBiz(store, requester)

		result, err := business.ListTrade(c.Request.Context(), &filter, &paging)

		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask(common.IsAdmin(requester))

			if i == len(result)-1 {
				paging.NextCursor = result[i].FakeId.String()
			}
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, filter))
	}
}
