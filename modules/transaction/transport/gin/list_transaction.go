package txgin

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"trading-service/common"
	txbiz "trading-service/modules/transaction/biz"
	txmodel "trading-service/modules/transaction/model"
	txstorage "trading-service/modules/transaction/storage"
)

func ListTransactions(sc goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var filter txmodel.Filter

		if err := c.ShouldBind(&filter); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		if err := paging.Process(); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		//requester := c.MustGet(common.CurrentUser).(common.Requester)
		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)

		store := txstorage.NewSQLStore(db)
		biz := txbiz.NewListTxBiz(store)

		result, err := biz.ListTx(c.Request.Context(), &filter, &paging)

		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask(common.DbTypeBSCTx)

			if i == len(result)-1 {
				paging.NextCursor = result[i].FakeId.String()
			}
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, filter))
	}
}
