package command

import (
	"context"
	"exchange-go/internal/currency"
	"exchange-go/internal/platform"
	"time"

	"go.uber.org/zap"
)

type generateKlineSyncCmd struct {
	currencyService currency.Service
	klineService    currency.KlineService
	logger          platform.Logger
}

func (cmd *generateKlineSyncCmd) Run(ctx context.Context, flags []string) {
	cmd.logger.Info("start of generate kline sync command")
	now := time.Now()
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	lastDay := now.Add(-24 * time.Hour)
	startTime := time.Date(lastDay.Year(), lastDay.Month(), lastDay.Day(), 0, 0, 0, 0, lastDay.Location())

	activeTimeFrames := cmd.klineService.GetActiveTimeFrames()
	pairs := cmd.currencyService.GetActivePairCurrenciesList()

	for _, pair := range pairs {
		for _, frame := range activeTimeFrames {
			params := currency.CreateKlineSyncParams{
				StartTime:  startTime.Format("2006-01-02 15:04:05"),
				EndTime:    endTime.Format("2006-01-02 15:04:05"),
				PairID:     pair.ID,
				TimeFrame:  frame,
				WithUpdate: false,
				Type:       currency.SyncTypeAuto,
			}

			err := cmd.klineService.CreateKlineSync(params)
			if err != nil {
				cmd.logger.Error2("error creating klineSync", err,
					zap.String("service", "generateKlineSyncCmd"),
					zap.String("method", "Run"),
					zap.String("pairName", pair.Name),
				)
				continue
			}
		}

	}

	cmd.logger.Info("end of generate kline sync command")
}

func NewGenerateKlineSyncCmd(currencyService currency.Service, klineService currency.KlineService, logger platform.Logger) ConsoleCommand {
	return &generateKlineSyncCmd{
		currencyService: currencyService,
		klineService:    klineService,
		logger:          logger,
	}
}
