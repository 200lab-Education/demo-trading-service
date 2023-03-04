package assetgin

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"trading-service/common"
	assetmodel "trading-service/modules/assets/model"
	assetstorage "trading-service/modules/assets/storage"
)

func ListAssets(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)

		var filter assetmodel.Filter

		if err := c.ShouldBind(&filter); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		_ = paging.Process()

		assetStore := assetstorage.NewSQLStore(db)

		result, err := assetStore.List(c.Request.Context(), &filter, &paging)

		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask(false)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
