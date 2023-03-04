package handlers

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"trading-service/common"
	"trading-service/middleware"
	assetgin "trading-service/modules/assets/transport/gin"
	authgin "trading-service/modules/auth/transport/gin"
	bankAccTprt "trading-service/modules/banking/transport"
	fiatTxTprt "trading-service/modules/fiattx/transport"
	locationgin "trading-service/modules/location/transport/gin"
	p2pmodel "trading-service/modules/p2p/model"
	p2pTransport "trading-service/modules/p2p/transport"
	txgin "trading-service/modules/transaction/transport/gin"
	userstorage "trading-service/modules/user/storage"
	usergin "trading-service/modules/user/transport/gin"
)

func MainRoute(v1 *gin.RouterGroup, sc goservice.ServiceContext) {
	dbConn := sc.MustGet(common.PluginMainDB).(*gorm.DB)

	v1.POST("/register", usergin.Register(sc))
	v1.POST("/login", usergin.Login(sc))

	auth := v1.Group("/auth")
	{
		auth.GET("/nonce", authgin.RequestNonce(sc))
		auth.POST("/verify_signature", authgin.VerifySignature(sc))
	}

	authStore := userstorage.NewSQLStore(dbConn)

	v1.GET("/profile", middleware.RequiredAuth(sc, authStore), usergin.GetProfile(sc))
	v1.PUT("/profile", middleware.RequiredAuth(sc, authStore), usergin.UpdateProfile(sc))
	v1.GET("/my-wallets", middleware.RequiredAuth(sc, authStore), assetgin.ListUserWallets(sc))

	v1.POST("/my-wallets/:wallet_id/withdraw/:asset_id", middleware.RequiredAuth(sc, authStore), txgin.BscWithdrawMoney(sc))

	v1.GET("/profile/nonce", middleware.RequiredAuth(sc, authStore), authgin.RequestNonceOldUser(sc))

	v1.GET("/locations", locationgin.ListCountries(sc))

	bankAccounts := v1.Group("/bank-accounts", middleware.RequiredAuth(sc, authStore))
	{
		bankAccounts.GET("", bankAccTprt.ListBankAccount(sc))
		bankAccounts.GET("/:id", bankAccTprt.GetBankAccount(sc))
		bankAccounts.POST("", bankAccTprt.CreateBankAccount(sc))
		bankAccounts.PUT("/:id", bankAccTprt.UpdateBankAccount(sc))
		bankAccounts.DELETE("/:id", bankAccTprt.DeleteBankAccount(sc))
	}

	fiatTransactions := v1.Group("/fiat-dw-transactions", middleware.RequiredAuth(sc, authStore))
	{
		fiatTransactions.GET("", fiatTxTprt.ListFiatTx(sc, true))
		fiatTransactions.POST("", fiatTxTprt.UserCreateFiatTx(sc))
		fiatTransactions.GET("/:id", fiatTxTprt.GetTx(sc))
		fiatTransactions.PUT("/:id/st/:status", fiatTxTprt.UserUpdateStatusTx(sc))
	}

	p2pOrders := v1.Group("/p2p-orders", middleware.RequiredAuth(sc, authStore))
	{
		p2pOrders.GET("", p2pTransport.ListP2pOrder(sc, true))
		p2pOrders.POST("", p2pTransport.UserCreateOrder(sc))
		p2pOrders.GET("/:id", p2pTransport.GetOrder(sc))
		p2pOrders.GET("/:id/trades", p2pTransport.ListP2pTrade(sc, true))
		p2pOrders.PUT("/:id/cancel", p2pTransport.CancelOrder(sc, "cancelled"))
		p2pOrders.POST("/:id/open-trade", p2pTransport.UserOpenTrade(sc))

		p2pTrades := v1.Group("/p2p-trades", middleware.RequiredAuth(sc, authStore))
		{
			p2pTrades.GET("", p2pTransport.ListP2pTrade(sc, true))
			p2pTrades.GET("/:trade-id", p2pTransport.GetTrade(sc))
			p2pTrades.PUT("/:trade-id/cancel-trade", p2pTransport.CancelTrading(sc, false))
			p2pTrades.PUT("/:trade-id/confirm", p2pTransport.UpdateTrading(sc, p2pmodel.TradeStWaitingPay.String()))
			p2pTrades.PUT("/:trade-id/transferred", p2pTransport.UpdateTrading(sc, p2pmodel.TradeStWaitVerify.String()))
			p2pTrades.PUT("/:trade-id/verify", p2pTransport.UpdateTrading(sc, p2pmodel.TradeStVerified.String()))
		}
	}

	metadata := v1.Group("/metadata")
	{
		metadata.GET("/deposit-accounts", bankAccTprt.SystemBankAccount(sc))
		metadata.GET("/assets", assetgin.ListAssets(sc))
	}
}
