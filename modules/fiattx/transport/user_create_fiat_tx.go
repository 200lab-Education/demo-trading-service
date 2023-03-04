package transport

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"trading-service/common"
	assetstore "trading-service/modules/assets/storage"
	bankStore "trading-service/modules/banking/storage"
	"trading-service/modules/fiattx/biz"
	"trading-service/modules/fiattx/model"
	"trading-service/modules/fiattx/storage"
	"trading-service/plugin/locker"
)

func UserCreateFiatTx(sc goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)
		requester := c.MustGet(common.CurrentUser).(common.Requester)
		lck := sc.MustGet(common.PluginMutexLock).(locker.Locker)

		var data model.FiatDWCreate

		if err := c.ShouldBind(&data); err != nil {
			panic(err)
		}

		data.UserId = requester.GetUserId()

		store := storage.NewSQLStore(db)
		assetStore := assetstore.NewSQLStore(db)
		bankAccStore := bankStore.NewSQLStore(db)
		business := biz.NewCreateBiz(store, assetStore, bankAccStore, lck)

		if c.DefaultQuery("type", "deposit") == "deposit" {
			if err := business.CreateDepositTx(c.Request.Context(), &data); err != nil {
				panic(err)
			}

		} else {
			if err := business.CreateWithdrawTx(c.Request.Context(), &data); err != nil {
				panic(err)
			}
		}

		data.Mask(common.DbTypeFiatTx)
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.FakeId.String()))
	}
}
