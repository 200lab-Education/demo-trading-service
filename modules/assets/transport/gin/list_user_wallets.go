package assetgin

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"trading-service/common"
	assetbiz "trading-service/modules/assets/biz"
	assetmodel "trading-service/modules/assets/model"
	assetstorage "trading-service/modules/assets/storage"
	walletstorage "trading-service/modules/wallet/storage"
)

func ListUserWallets(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)

		var filter assetmodel.Filter

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		filter.OwnerId = requester.GetUserId()

		walletStore := walletstorage.NewLocalStore()
		assetStore := assetstorage.NewSQLStore(db)

		biz := assetbiz.NewListUserWalletBiz(requester, assetStore, walletStore)

		result, err := biz.ListUserWallets(c.Request.Context(), &filter)

		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask(false)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
