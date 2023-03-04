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
)

func Register(sc goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)
		var data usermodel.UserCreate

		if err := c.ShouldBind(&data); err != nil {
			panic(err)
		}

		store := userstorage.NewSQLStore(db)
		md5 := hasher.NewMd5Hash()
		biz := userbiz.NewRegisterBusiness(store, md5)

		if err := biz.Register(c.Request.Context(), &data); err != nil {
			panic(err)
		}

		data.Mask(common.DbTypeUser)

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.FakeId.String()))
	}
}
