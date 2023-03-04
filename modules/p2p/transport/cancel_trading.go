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

func CancelTrading(sc goservice.ServiceContext, isAdminReject bool) func(*gin.Context) {
	return func(c *gin.Context) {
		//db := sc.MustGet(common.PluginMainDB).(*gorm.DB)
		lck := sc.MustGet(common.PluginMutexLock).(locker.Locker)
		requester := c.MustGet(common.CurrentUser).(common.Requester)

		id, err := common.FromBase58(c.Param("trade-id"))

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var data model.P2pTradingUpdate

		if err := c.ShouldBind(&data); err != nil {
			panic(err)
		}

		data.IsAdminReject = isAdminReject
		data.Status = model.TradeStCancelled.String()

		if isAdminReject {
			data.Status = model.TradeStRejected.String()
		}

		store := storage.NewSQLStore(sc.MustGet(common.PluginMainDB).(*gorm.DB))
		assetStore := assetstore.NewSQLStore(sc.MustGet(common.PluginMainDB).(*gorm.DB))
		masterTxStore := mstTxStore.NewSQLStore(sc.MustGet(common.PluginMainDB).(*gorm.DB))
		bankAccStore := bankStore.NewSQLStore(sc.MustGet(common.PluginMainDB).(*gorm.DB))
		business := biz.NewUpdateStBiz(store, lck, assetStore, bankAccStore, masterTxStore, requester)

		if err := business.CancelTrading(c.Request.Context(), int(id.GetLocalID()), &data); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
