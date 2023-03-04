package middleware

import (
	"context"
	"errors"
	"fmt"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"strings"
	"trading-service/common"
	usermodel "trading-service/modules/user/model"
	"trading-service/plugin/tokenprovider"
)

type AuthStore interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*usermodel.User, error)
}

func ErrWrongAuthHeader(err error) *common.AppError {
	return common.NewCustomError(
		err,
		fmt.Sprintf("wrong authen header"),
		fmt.Sprintf("ErrWrongAuthHeader"),
	)
}

func extractTokenFromHeaderString(s string) (string, error) {
	parts := strings.Split(s, " ")
	//"Authorization" : "Bearer {token}"

	if parts[0] != "Bearer" || len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return "", ErrWrongAuthHeader(nil)
	}

	return parts[1], nil
}

// RequiredAuth
// 1. Get token from header
// 2. Validate token and parse to payload
// 3. From the token payload, we use user_id to find from DB
func RequiredAuth(sc goservice.ServiceContext, authStore AuthStore) func(c *gin.Context) {
	tokenProvider := sc.MustGet(common.PluginJWT).(tokenprovider.Provider)

	return func(c *gin.Context) {
		token, err := extractTokenFromHeaderString(c.GetHeader("Authorization"))

		if err != nil {
			panic(err)
		}

		payload, err := tokenProvider.Validate(token)
		if err != nil {
			panic(err)
		}

		user, err := authStore.FindUser(c.Request.Context(), map[string]interface{}{"id": payload.UserId()})

		if err != nil {
			panic(err)
		}

		if user.Active != 1 {
			panic(common.ErrNoPermission(errors.New("user has been deleted or banned")))
		}

		user.Mask(common.DbTypeUser)

		c.Set(common.CurrentUser, user)
		c.Next()
	}
}
