package authgin

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"trading-service/common"
	authbiz "trading-service/modules/auth/biz"
	authmodel "trading-service/modules/auth/model"
	authstorage "trading-service/modules/auth/storage"
	"trading-service/plugin/tokenprovider"
)

func VerifySignature(sc goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)
		tokenProvider := sc.MustGet(common.PluginJWT).(tokenprovider.Provider)

		var data authmodel.AuthVerifyData

		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := authstorage.NewSQLStore(db)

		biz := authbiz.NewVerifySignatureBiz(store, tokenProvider)

		accessToken, err := biz.VerifySignature(c.Request.Context(), &data)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(accessToken))
	}
}
