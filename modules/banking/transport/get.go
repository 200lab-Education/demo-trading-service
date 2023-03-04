package transport

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"trading-service/common"
	"trading-service/modules/banking/biz"
	"trading-service/modules/banking/storage"
)

func GetBankAccount(sc goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)

		userId, err := common.FromBase58(c.Param("id"))

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := storage.NewSQLStore(db)
		business := biz.NewGetBiz(store)

		data, err := business.GetBankAccount(c.Request.Context(), int(userId.GetLocalID()))

		if err != nil {
			panic(err)
		}

		data.Mask(common.DbTypeBankAcc)

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}
