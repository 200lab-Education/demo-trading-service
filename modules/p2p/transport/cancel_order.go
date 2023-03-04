package transport

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"trading-service/common"
	assetstore "trading-service/modules/assets/storage"
	bankStore "trading-service/modules/banking/storage"
	mstTxStore "trading-service/modules/mastertx/storage"
	"trading-service/modules/p2p/biz"
	"trading-service/modules/p2p/model"
	"trading-service/modules/p2p/storage"
	"trading-service/plugin/locker"
)

func CancelOrder(sc goservice.ServiceContext, status string) func(*gin.Context) {
	return func(c *gin.Context) {
		//db := sc.MustGet(common.PluginMainDB).(*gorm.DB)
		lck := sc.MustGet(common.PluginMutexLock).(locker.Locker)
		requester := c.MustGet(common.CurrentUser).(common.Requester)

		id, err := common.FromBase58(c.Param("id"))

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var updateData model.P2pOrderUpdate

		updateData.Status = status

		store := storage.NewSQLStore(sc.MustGet(common.PluginMainDB).(*gorm.DB))
		assetStore := assetstore.NewSQLStore(sc.MustGet(common.PluginMainDB).(*gorm.DB))
		masterTxStore := mstTxStore.NewSQLStore(sc.MustGet(common.PluginMainDB).(*gorm.DB))
		bankAccStore := bankStore.NewSQLStore(sc.MustGet(common.PluginMainDB).(*gorm.DB))
		business := biz.NewUpdateStBiz(store, lck, assetStore, bankAccStore, masterTxStore, requester)

		if err := business.CancelOrDeleteOrder(c.Request.Context(), int(id.GetLocalID()), &updateData); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
