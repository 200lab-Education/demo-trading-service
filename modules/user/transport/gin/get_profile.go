package usergin

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"net/http"
	"trading-service/common"
)

func GetProfile(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := c.MustGet(common.CurrentUser)
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(u))
	}
}
