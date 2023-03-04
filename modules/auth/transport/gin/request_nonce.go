package authgin

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"trading-service/common"
	authbiz "trading-service/modules/auth/biz"
	authstorage "trading-service/modules/auth/storage"
)

func RequestNonce(sc goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		walletAddress := c.DefaultQuery("address", "")

		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)

		store := authstorage.NewSQLStore(db)
		biz := authbiz.NewRequestNonceBiz(store)

		result, err := biz.RequestNonce(c.Request.Context(), walletAddress)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}

func RequestNonceOldUser(sc goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)

		requester := c.MustGet(common.CurrentUser).(common.Requester)

		walletAddress := c.DefaultQuery("address", "")

		store := authstorage.NewSQLStore(db)
		biz := authbiz.NewRequestNonceBiz(store)

		result, err := biz.RequestNonceOldUser(c.Request.Context(), requester, walletAddress)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
