package transport

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"trading-service/common"
	"trading-service/modules/p2p/biz"
	"trading-service/modules/p2p/storage"
)

func GetOrder(sc goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)
		requester := c.MustGet(common.CurrentUser).(common.Requester)

		id, err := common.FromBase58(c.Param("id"))

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := storage.NewSQLStore(db)
		business := biz.NewGetBiz(store, requester)

		data, err := business.GetOrder(c.Request.Context(), int(id.GetLocalID()))

		if err != nil {
			panic(err)
		}

		data.Mask(common.IsAdmin(requester))

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}
