package txgin

import (
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"trading-service/common"
	assetstorage "trading-service/modules/assets/storage"
	txbiz "trading-service/modules/transaction/biz"
	txmodel "trading-service/modules/transaction/model"
	txstorage "trading-service/modules/transaction/storage"
	"trading-service/plugin/bscex"
	locker "trading-service/plugin/locker"
)

func BscWithdrawMoney(sc goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		db := sc.MustGet(common.PluginMainDB).(*gorm.DB)
		lck := sc.MustGet(common.PluginMutexLock).(locker.Locker)
		scCaller := sc.MustGet(common.PluginBSCEx).(bscex.BSCEx)

		var withdrawData txmodel.WithdrawData

		if err := c.ShouldBind(&withdrawData); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		if withdrawData.Type == "" {
			withdrawData.Type = "bsc"
		}

		if err := withdrawData.Validate(); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)

		walletUID, err := common.FromBase58(c.Param("wallet_id"))

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		assetUID, err := common.FromBase58(c.Param("asset_id"))

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		userAssetStore := assetstorage.NewSQLStore(db)
		bscTxStore := txstorage.NewSQLStore(db)

		biz := txbiz.NewBSCWithdrawBiz(requester, userAssetStore, bscTxStore, lck, scCaller)

		txHash, err := biz.WithdrawAsset(
			c.Request.Context(),
			withdrawData.Receiver,
			int(assetUID.GetLocalID()),
			int(walletUID.GetLocalID()),
			withdrawData.Amount,
		)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(txHash))
	}
}
