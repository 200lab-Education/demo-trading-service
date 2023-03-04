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

func ListBankAccount(sc goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)
		requester := c.MustGet(common.CurrentUser).(common.Requester)

		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		_ = paging.Process()

		var filter model.Filter

		filter.UserId = requester.GetUserId()

		store := storage.NewSQLStore(db)
		business := biz.NewListBiz(store)

		result, err := business.ListBankAccounts(c.Request.Context(), &filter, &paging)

		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask(common.DbTypeBankAcc)

			if i == len(result)-1 {
				paging.NextCursor = result[i].FakeId.String()
			}
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}

func SystemBankAccount(sc goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)

		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var filter model.Filter

		filter.UserId = 0
		filter.IsSystemAccount = true

		store := storage.NewSQLStore(db)
		business := biz.NewListBiz(store)

		result, err := business.ListBankAccounts(c.Request.Context(), &filter, &paging)

		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask(common.DbTypeBankAcc)

			if i == len(result)-1 {
				paging.NextCursor = result[i].FakeId.String()
			}
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
