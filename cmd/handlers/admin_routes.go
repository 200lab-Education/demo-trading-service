package handlers

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"trading-service/common"
	"trading-service/middleware"
	fiatTxTprt "trading-service/modules/fiattx/transport"
	p2pmodel "trading-service/modules/p2p/model"
	p2pTransport "trading-service/modules/p2p/transport"
	userstorage "trading-service/modules/user/storage"
	usergin "trading-service/modules/user/transport/gin"
)

func AdminRoutes(v1 *gin.RouterGroup, sc goservice.ServiceContext) {
	dbConn := sc.MustGet(common.PluginMainDB).(*gorm.DB)
	authStore := userstorage.NewSQLStore(dbConn)

	admin := v1.Group("/admin", middleware.RequiredAuth(sc, authStore), middleware.RequiredRoles(sc, common.RoleAdmin, common.RoleMod))

	users := admin.Group("/users")
	{
		users.GET("", usergin.ListUsers(sc))
		users.GET("/:user_id", usergin.GetUser(sc))
		users.PUT("/:user_id", usergin.UpdateUser(sc))
	}

	fiatTransactions := admin.Group("/fiat-dw-transactions")
	{
		fiatTransactions.GET("", fiatTxTprt.ListFiatTx(sc, false))
		fiatTransactions.GET("/:id", fiatTxTprt.GetTx(sc))
		fiatTransactions.PUT("/:id/st/:status", fiatTxTprt.UserUpdateStatusTx(sc))
	}

	p2pOrders := admin.Group("/p2p-orders")
	{
		p2pOrders.GET("", p2pTransport.ListP2pOrder(sc, false))
		p2pOrders.GET("/:id", p2pTransport.GetOrder(sc))
		p2pOrders.GET("/:id/trades", p2pTransport.ListP2pTrade(sc, false))
		p2pOrders.PUT("/:id/cancel", p2pTransport.CancelOrder(sc, "cancelled"))

		p2pTrades := admin.Group("/p2p-trades")
		{
			p2pTrades.GET("", p2pTransport.ListP2pTrade(sc, false))
			p2pTrades.GET("/:trade-id", p2pTransport.GetTrade(sc))
			p2pTrades.PUT("/:trade-id/cancel-trade", p2pTransport.CancelTrading(sc, false))
			p2pTrades.PUT("/:trade-id/verify", p2pTransport.UpdateTrading(sc, p2pmodel.TradeStVerified.String()))
		}
	}

}
