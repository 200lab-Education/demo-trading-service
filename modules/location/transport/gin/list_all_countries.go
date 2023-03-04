package locationgin

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"trading-service/common"
	"trading-service/modules/location/model"
	"trading-service/modules/location/storage"
)

func ListCountries(sc goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var filter model.Filter

		if err := c.ShouldBind(&filter); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		if err := paging.Process(); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		paging.Limit = 200

		//requester := c.MustGet(common.CurrentUser).(common.Requester)
		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)

		store := storage.NewSQLStore(db)

		result, err := store.ListCountries(c.Request.Context(), &filter, &paging)

		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask()

			if i == len(result)-1 {
				paging.NextCursor = result[i].FakeId.String()
			}
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, filter))
	}
}
