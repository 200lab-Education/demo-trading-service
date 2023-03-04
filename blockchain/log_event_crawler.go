package blockchain

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"math/big"
	"strconv"
	"strings"
	"time"
	appCommon "trading-service/common"
	"trading-service/component/asyncjob"
	"trading-service/plugin/bscex"
)

type BlockTrackerStorage interface {
	GetLatestBlockNumber(ctx context.Context) (int64, error)
	UpdateLatestBlockNumber(ctx context.Context, num int64) error
}

type scLogCrawler struct {
	bscStore     bscex.BSCEx
	currentBlock int64
	latestBlock  int64
	storage      BlockTrackerStorage
	logChan      chan types.Log
}

func NewSCLogCrawler(
	bscStore bscex.BSCEx,
	store BlockTrackerStorage,
) *scLogCrawler {
	mkpCrawler := &scLogCrawler{
		bscStore:     bscStore,
		storage:      store,
		logChan:      make(chan types.Log, 100),
		currentBlock: bscStore.CurrentBlock(),
	}

	return mkpCrawler
}

func (parser *scLogCrawler) Start() error {
	latestBlockNumber, err := parser.latestBlockNumber()
	parser.latestBlock = latestBlockNumber

	if err != nil {
		return err
	}

	currentBlockNumber, err := parser.latestDbBlockNumber()

	if err != nil {
		return err
	}

	parser.currentBlock = max(currentBlockNumber, parser.currentBlock)

	go func() {
		var stepBlockFastScan int64 = 200
		scanDelay := time.Second * 3

		if parser.bscStore.IsMainNet() {
			stepBlockFastScan = 20
			//scanDelay = time.Second * 3
		}

		for {
			time.Sleep(scanDelay)

			if err := parser.storage.UpdateLatestBlockNumber(context.Background(), parser.currentBlock); err != nil {
				log.Errorln(err)
				continue
			}

			if latestBlockNumber > parser.currentBlock {
				if v := latestBlockNumber - parser.currentBlock; v < stepBlockFastScan {
					stepBlockFastScan = v
				}

				if err := parser.scanBlock(parser.currentBlock, stepBlockFastScan); err != nil {
					log.Errorln(err)
					continue
				}

				parser.currentBlock += stepBlockFastScan + 1
				continue
			}

			latestBlockNumber, err = parser.latestBlockNumber()

			if err != nil {
				continue
			}

			parser.latestBlock = latestBlockNumber

			if err != nil {
				log.Errorln("Get latestBlockNumber:", err)
				continue
			}

			if latestBlockNumber <= parser.currentBlock {
				log.Debugln("Still no block available")
				continue
			}

			if err := parser.scanBlock(parser.currentBlock, 0); err != nil {
				log.Errorln(err)
				continue
			}

			parser.currentBlock += 1
		}
	}()

	return nil
}

func (parser *scLogCrawler) latestBlockNumber() (int64, error) {
	var result int64 = 1

	job := asyncjob.NewJob(func(ctx context.Context) error {
		latestBlock, err := parser.bscStore.ExchangeClient().HeaderByNumber(context.Background(), nil)

		if err != nil {
			return err
		}

		result, err = strconv.ParseInt(latestBlock.Number.String(), 10, 64)

		if err != nil {
			return err
		}

		return nil
	})

	job.SetRetryDurations(time.Second, time.Second*2, time.Second*3)

	if err := asyncjob.NewGroup(false, job).Run(context.Background()); err != nil {
		return 0, err
	}

	return result, nil
}

func (parser *scLogCrawler) latestDbBlockNumber() (int64, error) {
	var result int64 = 1

	job := asyncjob.NewJob(func(ctx context.Context) error {
		rs, err := parser.storage.GetLatestBlockNumber(ctx)

		if err != nil {
			if err == appCommon.ErrRecordNotFound {
				result = 1
				return nil
			}

			return err
		}

		result = rs
		return nil
	})

	job.SetRetryDurations(time.Second, time.Second*2, time.Second*3)

	if err := asyncjob.NewGroup(false, job).Run(context.Background()); err != nil {
		return 0, err
	}

	if result < parser.currentBlock {
		result = parser.currentBlock
	}

	return result, nil
}

func (parser *scLogCrawler) scanBlock(from, step int64) error {
	//log.Printf("Starting scan block from %d to %d. Latest onchain: %d", from, from+step, parser.latestBlock)

	logs, err := parser.bscStore.ExchangeClient().FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: big.NewInt(from),
		ToBlock:   big.NewInt(from + step),
		Addresses: []common.Address{common.HexToAddress(parser.bscStore.ExSCAddress())},
	})

	if err != nil {
		return err
	}

	for i, l := range logs {
		functionHash := strings.ToLower(l.Topics[0].Hex())
		log.Printf("Block %d - Tx %s - Function Hash %s \n", l.BlockNumber, l.TxHash.Hex(), functionHash)

		parser.logChan <- logs[i]
	}

	return nil
}

func (parser *scLogCrawler) GetLogsChan() <-chan types.Log { return parser.logChan }

func stringOrDefault(s, d string) string {
	if s == "" {
		return d
	}

	return s
}

func max(n1, n2 int64) int64 {
	if n1 > n2 {
		return n1
	}

	return n2
}
