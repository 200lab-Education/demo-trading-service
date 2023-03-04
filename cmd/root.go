package cmd

import (
	"context"
	"fmt"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/200Lab-Education/go-sdk/plugin/storage/sdkgorm"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"net/http"
	"os"
	"trading-service/blockchain"
	"trading-service/cmd/handlers"
	"trading-service/common"
	"trading-service/middleware"
	assetstorage "trading-service/modules/assets/storage"
	authstorage "trading-service/modules/auth/storage"
	txstorage "trading-service/modules/transaction/storage"
	walletstorage "trading-service/modules/wallet/storage"
	"trading-service/plugin/bscex"
	"trading-service/plugin/locker/local"
	"trading-service/plugin/tokenprovider/jwt"
)

func newService() goservice.Service {
	service := goservice.New(
		goservice.WithName("trading-service"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(sdkgorm.NewGormDB("main", common.PluginMainDB)),
		goservice.WithInitRunnable(jwt.NewTokenJWTProvider(common.PluginJWT)),
		goservice.WithInitRunnable(local.NewLocalLocker(common.PluginMutexLock)),
		goservice.WithInitRunnable(bscex.NewSCCaller(common.PluginBSCEx)),
	)

	return service
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start an REST service",
	Run: func(cmd *cobra.Command, args []string) {
		service := newService()

		serviceLogger := service.Logger("service")

		if err := service.Init(); err != nil {
			serviceLogger.Fatalln(err)
		}

		db := service.MustGet(common.PluginMainDB).(*gorm.DB)
		bscExStore := service.MustGet(common.PluginBSCEx).(bscex.BSCEx)

		txStore := txstorage.NewSQLStore(db)
		userStore := authstorage.NewSQLStore(db)
		assetStore := assetstorage.NewSQLStore(db)
		walletStore := walletstorage.NewLocalStore()

		logCrawler := blockchain.NewSCLogCrawler(bscExStore, txStore)
		logHandler := blockchain.NewSCLogHdl(bscExStore, txStore, userStore, walletStore, assetStore)

		if err := logCrawler.Start(); err != nil {
			log.Errorln(err)
		}

		go logHandler.Run(context.Background(), logCrawler.GetLogsChan())

		service.HTTPServer().AddHandler(func(engine *gin.Engine) {
			engine.Use(middleware.Recover(service))
			engine.Use(middleware.AllowCORS())

			engine.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"data": "pong"})
			})

			v1 := engine.Group("/v1")

			handlers.MainRoute(v1, service)
			handlers.AdminRoutes(v1, service)
		})

		if err := service.Start(); err != nil {
			serviceLogger.Fatalln(err)
		}
	},
}

func Execute() {
	// TransAddPoint outenv as a sub command
	rootCmd.AddCommand(outEnvCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
