package command

import (
	"context"
	"exchange-go/internal/currency"
	"exchange-go/internal/currency/candle"
	"exchange-go/internal/order"
	"exchange-go/internal/platform"
	"exchange-go/internal/user"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type setUserLevelCmd struct {
	db               *gorm.DB
	userService      user.Service
	tradeService     order.TradeService
	pairRepo         currency.PairRepository
	coinRepo         currency.Repository
	klineService     currency.KlineService
	userLevelService user.LevelService
	logger           platform.Logger

	pairAveragePrices []candle.AveragePairPrice
	allPairs          []currency.Pair
	btcCoin           currency.Coin
	startTime         time.Time
	endTime           time.Time
}

func (cmd *setUserLevelCmd) Run(ctx context.Context, flags []string) {
	cmd.logger.Info("start of set user level command")
	err := cmd.setNeededData(flags)
	if err != nil {
		cmd.logger.Error2("error setting needed  data", err,
			zap.String("service", "setUserLevel"),
			zap.String("method", "Run"),
		)
		return

	}
	page := 0
	pageSize := 100
	filters := map[string]interface{}{}
	startTimeString := cmd.startTime.Format("2006-01-02 15:04:05")
	endTimeString := cmd.endTime.Format("2006-01-02 15:04:05")
	for {
		users := cmd.userService.GetUsersByPagination(int64(page), pageSize, filters)
		if len(users) == 0 {
			break
		}
		for _, u := range users {
			exchangeNumber := int64(0)
			exchangeVolumeDecimal := decimal.NewFromFloat(0) //based on btc
			if u.IsLevelManuallySet.Valid && u.IsLevelManuallySet.Bool {
				continue
			}

			trades := cmd.tradeService.GetTradesOfUserBetweenTimes(u.ID, startTimeString, endTimeString)
			for _, t := range trades {
				exchangeNumber++
				amountDecimal := cmd.getAmountDecimalBasedOnBtc(ctx, t)
				exchangeVolumeDecimal = exchangeVolumeDecimal.Add(amountDecimal)
			}
			cmd.updateUserLevelAndExchangeData(u, exchangeVolumeDecimal.StringFixed(8), exchangeNumber)
		}

		page++
	}

	cmd.logger.Info("end of set user level command")
}

func (cmd *setUserLevelCmd) getAmountDecimalBasedOnBtc(ctx context.Context, t order.Trade) decimal.Decimal {
	tradeAmountDecimal, err := decimal.NewFromString(t.Amount.String)
	var tradePair currency.Pair
	for _, pair := range cmd.allPairs {
		if pair.ID == t.PairID {
			tradePair = pair
		}
	}

	tradedCoinID := tradePair.DependentCoin.ID

	if tradedCoinID == cmd.btcCoin.ID {
		return tradeAmountDecimal
	}
	var pairNeeded currency.Pair

	for _, p := range cmd.allPairs {
		if p.DependentCoin.ID == tradedCoinID && p.BasisCoinID == cmd.btcCoin.ID {
			pairNeeded = p
			break
		}
		if p.DependentCoin.ID == cmd.btcCoin.ID && p.BasisCoinID == tradedCoinID {
			pairNeeded = p
			break
		}
	}

	createdAtDay := t.CreatedAt.Format("2006-01-02")

	for _, item := range cmd.pairAveragePrices {
		if item.PairName == pairNeeded.Name && item.Day == createdAtDay {
			averagePrice := item.Price
			averagePriceDecimal, err := decimal.NewFromString(averagePrice)
			if err != nil {
				cmd.logger.Error2("error converting averagePrice string to Decimal", err,
					zap.String("service", "setUserLevel"),
					zap.String("method", "getAmountDecimalBasedOnBtc"),
					zap.String("averagePrice", averagePrice),
				)
			}

			if pairNeeded.BasisCoinID == cmd.btcCoin.ID {
				return tradeAmountDecimal.Mul(averagePriceDecimal)
			}
			return tradeAmountDecimal.Div(averagePriceDecimal)
		}
	}

	cmd.logger.Error2("no price for pair", err,
		zap.String("service", "setUserLevel"),
		zap.String("method", "getAmountDecimalBasedOnBtc"),
		zap.String("pairName", pairNeeded.Name),
		zap.String("createdAtDay", createdAtDay),
	)
	return tradeAmountDecimal

}

func (cmd *setUserLevelCmd) updateUserLevelAndExchangeData(u user.User, exchangeVolume string, exchangeNumber int64) {
	level := cmd.userLevelService.RecalculateUserLevel(u, exchangeVolume)
	updatingUser := &user.User{
		ID:                     u.ID,
		ExchangeNumber:         exchangeNumber,
		ExchangeVolumeAmount:   exchangeVolume,
		ExchangeVolumeCoinCode: "BTC",
		UserLevelID:            level.ID,
	}
	err := cmd.db.Model(updatingUser).Updates(updatingUser).Error
	if err != nil {
		cmd.logger.Error2("can not update user level", err,
			zap.String("service", "setUserLevel"),
			zap.String("method", "updateUserLevelAndExchangeData"),
			zap.Int("userID", u.ID),
			zap.Int64("levelID", level.ID),
		)
	}
}

func NewSetUserLevelCmd(db *gorm.DB, userService user.Service, tradeService order.TradeService, pairRepo currency.PairRepository, coinRepo currency.Repository, klineService currency.KlineService, userLevelService user.LevelService, logger platform.Logger) ConsoleCommand {
	cmd := &setUserLevelCmd{
		db:               db,
		userService:      userService,
		tradeService:     tradeService,
		pairRepo:         pairRepo,
		coinRepo:         coinRepo,
		klineService:     klineService,
		userLevelService: userLevelService,
		logger:           logger,
	}
	return cmd

}

func (cmd *setUserLevelCmd) setNeededData(flags []string) error {
	t := time.Now()
	t = t.Add(-1 * 31 * 24 * time.Hour) // 31 day ago
	startTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, t.Nanosecond(), t.Location())

	t2 := time.Now()
	t2 = t2.Add(-1 * 24 * time.Hour) //1 day ago
	endTime := time.Date(t2.Year(), t2.Month(), t2.Day(), 23, 59, 59, t.Nanosecond(), t.Location())
	cmd.startTime = startTime
	cmd.endTime = endTime

	cmd.allPairs = cmd.pairRepo.GetAllPairs()

	btcCoin := currency.Coin{}
	_ = cmd.coinRepo.GetCoinByCode("BTC", &btcCoin)
	cmd.btcCoin = btcCoin

	var pairNames []string
	for _, p := range cmd.allPairs {
		pairNames = append(pairNames, p.Name)
	}

	pairAveragePrices, err := cmd.klineService.GetAveragePriceOfPairsUsingTime(pairNames, startTime, endTime)
	if err != nil {
		return err
	}
	cmd.pairAveragePrices = pairAveragePrices
	return nil
	//startTimeString := flag.String("startTime", "", "")
	//endTimeString := flag.String("endTime", "", "")
	//err := flag.CommandLine.Parse(flags)
	//if err != nil {
	//	fmt.Println("errrrrrs", err)
	//}
	//
	//t := time.Now()
	//t = t.Add(-1 * 24 * time.Hour)
	//startTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, t.Nanosecond(), t.Location())
	////.Format("2006-01-02 15:04:05")
	//endTime := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, t.Nanosecond(), t.Location())
	//
	//if *startTimeString != "" {
	//	t, err = time.Parse("2006-01-02", *startTimeString)
	//	if err != nil {
	//		fmt.Println("startTime format is not correct")
	//		os.Exit(1)
	//	}
	//	startTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, t.Nanosecond(), t.Location())
	//}
	//
	//if *endTimeString != "" {
	//	t, err = time.Parse("2006-01-02", *endTimeString)
	//	if err != nil {
	//		fmt.Println("endTime format is not correct")
	//		os.Exit(1)
	//	}
	//	endTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, t.Nanosecond(), t.Location())
	//}
	//
	//if startTime.Sub(endTime).Seconds() > 0 {
	//	fmt.Println("endTime is smaller than start Date")
	//	os.Exit(1)
	//}

}
