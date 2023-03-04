package usergin

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"trading-service/common"
	"trading-service/component/hasher"
	userbiz "trading-service/modules/user/biz"
	usermodel "trading-service/modules/user/model"
	userstorage "trading-service/modules/user/storage"
	"trading-service/plugin/tokenprovider"
)

func Login(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginUserData usermodel.UserLogin

		if err := c.ShouldBind(&loginUserData); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)
		tokenProvider := sc.MustGet(common.PluginJWT).(tokenprovider.Provider)

		store := userstorage.NewSQLStore(db)
		md5 := hasher.NewMd5Hash()

		business := userbiz.NewLoginBusiness(store, tokenProvider, md5, 60*60*24*30)
		account, err := business.Login(c.Request.Context(), &loginUserData)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(account))
	}
}
