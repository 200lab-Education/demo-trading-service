package transport

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"trading-service/common"
	assetstore "trading-service/modules/assets/storage"
	bankStore "trading-service/modules/banking/storage"
	"trading-service/modules/p2p/biz"
	"trading-service/modules/p2p/model"
	"trading-service/modules/p2p/storage"
	"trading-service/plugin/locker"
)

func UserCreateOrder(sc goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		lck := sc.MustGet(common.PluginMutexLock).(locker.Locker)
		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)
		requester := c.MustGet(common.CurrentUser).(common.Requester)

		var data model.P2pOrderCreate

		if err := c.ShouldBind(&data); err != nil {
			panic(err)
		}

		data.Type = c.DefaultQuery("type", model.P2pTypeSell)

		//data.Type = model.P2pTypeSell // hard code because buy feature have not done yet
		data.Status = model.OrdStActive.String()
		data.UserId = requester.GetUserId()

		store := storage.NewSQLStore(db)
		assetStore := assetstore.NewSQLStore(sc.MustGet(common.PluginMainDB).(*gorm.DB))
		bankAccStore := bankStore.NewSQLStore(sc.MustGet(common.PluginMainDB).(*gorm.DB))

		business := biz.NewCreateOfferBiz(store, assetStore, bankAccStore, lck)

		if err := business.CreateSellOffer(c.Request.Context(), &data); err != nil {
			panic(err)
		}

		data.Mask(common.DbTypeP2pOrder)

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.FakeId.String()))
	}
}
