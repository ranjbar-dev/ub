package command

import (
	"exchange-go/internal/currency"
	"flag"
	"fmt"
	"strings"
	"time"
)

func (cmd *syncKlineCmd) setNeededData(flags []string) (shouldConsiderFlags bool, err error) {
	if len(flags) < 1 {
		return false, nil
	}

	pairName := flag.String("pair", "", "")
	timeFrame := flag.String("frame", "", "")
	start := flag.String("start", "", "")
	end := flag.String("end", "", "")
	withUpdate := flag.Bool("update", false, "")
	err = flag.CommandLine.Parse(flags)
	if err != nil {
		return false, err
	}

	pair, err := cmd.currencyService.GetPairByName(strings.ToUpper(*pairName))
	if err != nil {
		return false, err
	}
	cmd.pair = pair

	activeTimeFrames := cmd.klineService.GetActiveTimeFrames()
	isTimeFrameCorrect := false
	for _, t := range activeTimeFrames {
		if t == *timeFrame {
			isTimeFrameCorrect = true
			break
		}
	}
	if isTimeFrameCorrect {
		cmd.timeFrame = *timeFrame
	} else {
		return false, fmt.Errorf("time frame is not correct")
	}

	startTime, err := time.Parse("2006-01-02 15:04:05", *start)
	if err != nil {
		return false, err
	} else {
		cmd.startTime = startTime
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", *end)
	if err != nil {
		return false, err
	} else {
		cmd.endTime = endTime
	}

	cmd.withUpdate = *withUpdate

	return true, nil
}

func (cmd *syncKlineCmd) getParametersFromKlineSync() (klineSync *currency.KlineSync, params neededParams, err error) {
	processingKlineSyncs := cmd.klineService.GetKlineSyncsByStatusAndLimit(currency.SyncStatusProcessing, 10)
	if len(processingKlineSyncs) > 0 {
		//if the oldest running ohlc sync is older than 10 hour we considered it at not running and try to run that
		now := time.Now()
		oldestProcessingKlineSync := processingKlineSyncs[0]
		last10Hour := now.Add(-10 * time.Hour)
		if oldestProcessingKlineSync.UpdatedAt.Sub(last10Hour) > 0 {
			return nil, params, nil
		}
		klineSync = &oldestProcessingKlineSync
	} else {
		createdKlineSyncs := cmd.klineService.GetKlineSyncsByStatusAndLimit(currency.SyncStatusCreated, 1)
		if len(createdKlineSyncs) > 0 {
			klineSync = &createdKlineSyncs[0]
		}
	}

	if klineSync != nil {
		pair, err := cmd.currencyService.GetPairByID(klineSync.PairID)
		if err != nil {
			return klineSync, params, err
		}

		activeTimeFrames := cmd.klineService.GetActiveTimeFrames()
		isTimeFrameCorrect := false
		for _, t := range activeTimeFrames {
			if t == klineSync.TimeFrame.String {
				isTimeFrameCorrect = true
				break
			}
		}

		if !isTimeFrameCorrect {
			return klineSync, params, fmt.Errorf("timeFrame is not correct in klineSync table")
		}

		params = neededParams{
			pair:       pair,
			timeFrame:  klineSync.TimeFrame.String,
			startTime:  klineSync.StartTime.Time,
			endTime:    klineSync.EndTime.Time,
			withUpdate: klineSync.WithUpdate,
		}

		return klineSync, params, nil
	}

	return nil, params, nil
}
