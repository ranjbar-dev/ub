package command

import (
	"context"
	"exchange-go/internal/communication"
	"exchange-go/internal/currency"
	"exchange-go/internal/externalexchange"
	"exchange-go/internal/platform"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type syncKlineCmd struct {
	currencyService         currency.Service
	klineService            currency.KlineService
	externalExchangeService externalexchange.Service
	queueManager            communication.QueueManager
	logger                  platform.Logger
	pair                    currency.Pair
	timeFrame               string
	startTime               time.Time
	endTime                 time.Time
	withUpdate              bool
}

type neededParams struct {
	pair       currency.Pair
	timeFrame  string
	startTime  time.Time
	endTime    time.Time
	withUpdate bool
}

func (cmd *syncKlineCmd) Run(ctx context.Context, flags []string) {
	fmt.Println("start of sync kline command")
	shouldConsiderFlags, err := cmd.setNeededData(flags)
	if err != nil {
		fmt.Println("parameters are not correct: " + err.Error())
		cmd.logger.Warn("parameters are not correct",
			zap.Error(err),
		)
		return
	}

	params := neededParams{}
	var klineSync *currency.KlineSync

	if shouldConsiderFlags {
		params = neededParams{
			pair:       cmd.pair,
			timeFrame:  cmd.timeFrame,
			startTime:  cmd.startTime,
			endTime:    cmd.endTime,
			withUpdate: cmd.withUpdate,
		}
	} else {
		klineSync, params, err = cmd.getParametersFromKlineSync()
		if err != nil {
			cmd.logger.Error2("error getting parameter from kline sync", err,
				zap.String("service", "syncKlineCmd"),
				zap.String("method", "Run"),
				zap.String("pairName", params.pair.Name),
			)
			return
		}
		if klineSync == nil {
			fmt.Println("there is an active klineSync or all of them are done")
			return
		}
	}

	if klineSync != nil {
		klineSync.Status = currency.SyncStatusProcessing
		err := cmd.klineService.UpdateKlineSync(klineSync)
		if err != nil {
			cmd.logger.Error2("error updating kline sync", err,
				zap.String("service", "syncKlineCmd"),
				zap.String("method", "Run"),
				zap.String("status", currency.SyncStatusProcessing),
				zap.Int64("klineSyncID", klineSync.ID),
			)
			return
		}

		err = cmd.fetchAndSaveKlines(params)

		if err != nil {
			cmd.logger.Error2("error in fetch and save kline", err,
				zap.String("service", "syncKlineCmd"),
				zap.String("method", "Run"),
				zap.String("pairName", params.pair.Name),
			)

			klineSync.Status = currency.SyncStatusError
			err := cmd.klineService.UpdateKlineSync(klineSync)
			if err != nil {
				cmd.logger.Error2("error updating kline sync", err,
					zap.String("service", "syncKlineCmd"),
					zap.String("method", "Run"),
					zap.String("status", currency.SyncStatusError),
					zap.Int64("klineSyncID", klineSync.ID),
					zap.Int64("pairID", klineSync.PairID),
				)
				return
			}
			return
		}

		klineSync.Status = currency.SyncStatusDone
		err = cmd.klineService.UpdateKlineSync(klineSync)
		if err != nil {
			cmd.logger.Error2("error updating kline sync", err,
				zap.String("service", "syncKlineCmd"),
				zap.String("method", "Run"),
				zap.String("status", currency.SyncStatusDone),
				zap.Int64("klineSyncID", klineSync.ID),
				zap.Int64("pairID", klineSync.PairID),
			)
		}
	} else {
		err = cmd.fetchAndSaveKlines(params)

		if err != nil {
			cmd.logger.Error2("can not fetch and save klines", err,
				zap.String("service", "syncKlineCmd"),
				zap.String("method", "Run"),
				zap.String("pairName", params.pair.Name),
			)
			return
		}
	}
	fmt.Println("end of sync kline command")
}

func NewSyncKlineCmd(currencyService currency.Service, klineService currency.KlineService,
	externalExchangeService externalexchange.Service, queueManager communication.QueueManager, logger platform.Logger) ConsoleCommand {
	return &syncKlineCmd{
		currencyService:         currencyService,
		klineService:            klineService,
		externalExchangeService: externalExchangeService,
		queueManager:            queueManager,
		logger:                  logger,
	}
}
