package transport

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"trading-service/common"
	"trading-service/modules/banking/biz"
	"trading-service/modules/banking/model"
	"trading-service/modules/banking/storage"
)

func UpdateBankAccount(sc goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)

		userId, err := common.FromBase58(c.Param("id"))
		requester := c.MustGet(common.CurrentUser).(common.Requester)

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var data model.BankAccountUpdate

		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := storage.NewSQLStore(db)
		business := biz.NewUpdateBiz(store, requester)

		if err := business.UpdateBankAccount(c.Request.Context(), int(userId.GetLocalID()), &data); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
