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

func CreateBankAccount(sc goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)
		requester := c.MustGet(common.CurrentUser).(common.Requester)

		var data model.BankAccountCreate

		if err := c.ShouldBind(&data); err != nil {
			panic(err)
		}

		data.UserId = requester.GetUserId()

		store := storage.NewSQLStore(db)
		business := biz.NewCreateBiz(store)

		if err := business.CreateBankAccount(c.Request.Context(), &data); err != nil {
			panic(err)
		}

		data.Mask(common.DbTypeBankAcc)

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.FakeId.String()))
	}
}
