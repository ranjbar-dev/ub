package command

import (
	"encoding/json"
	"exchange-go/internal/externalexchange"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type queueNeededKline struct {
	Pair                string `json:"pair"`
	TimeFrame           string `json:"timeFrame"`
	KlineStartTime      string `json:"startTime"`
	KlineCloseTime      string `json:"closeTime"`
	OpenPrice           string `json:"openPrice"`
	ClosePrice          string `json:"closePrice"`
	HighPrice           string `json:"highPrice"`
	LowPrice            string `json:"lowPrice"`
	BaseVolume          string `json:"baseVolume"`
	QuoteVolume         string `json:"quoteVolume"`
	TakerBuyBaseVolume  string `json:"takerBuyBaseVolume"`
	TakerBuyQuoteVolume string `json:"takerBuyQuoteVolume"`
	Spread              string `json:"spread"`
	IsOld               bool   `json:"isOld"`
}

func (cmd *syncKlineCmd) fetchAndSaveKlines(params neededParams) error {
	start := params.startTime
	end := params.endTime
	for end.Sub(start) >= 0 {
		fetchParams := externalexchange.FetchKlinesParams{
			Pair:      params.pair.Name,
			TimeFrame: params.timeFrame,
			From:      start,
			To:        params.endTime,
		}

		externalExchangeKlines, err := cmd.externalExchangeService.FetchKlines(fetchParams)
		if err != nil {
			return err
		}

		for _, k := range externalExchangeKlines {
			if k.StartTime.Sub(end) >= 0 {
				break
			}
			k.StartTime.Format("2006-01-02 15:04:05")

			queueNeededKline := queueNeededKline{
				Pair:                params.pair.Name,
				TimeFrame:           params.timeFrame,
				KlineStartTime:      k.StartTime.Format("2006-01-02 15:04:05"),
				KlineCloseTime:      k.EndTime.Format("2006-01-02 15:04:05"),
				OpenPrice:           k.OpenPrice,
				ClosePrice:          k.ClosePrice,
				HighPrice:           k.HighPrice,
				LowPrice:            k.LowPrice,
				BaseVolume:          k.BaseVolume,
				QuoteVolume:         k.QuoteVolume,
				TakerBuyBaseVolume:  k.TakerBuyBaseVolume,
				TakerBuyQuoteVolume: k.TakerBuyQuoteVolume,
				Spread:              strconv.FormatFloat(float64(params.pair.Spread), 'f', 8, 64),
				IsOld:               true,
			}

			payload, err := json.Marshal(queueNeededKline)
			if err != nil {
				cmd.logger.Warn("can not marshal kline",
					zap.Error(err),
					zap.String("service", "processor"),
					zap.String("method", "ProcessKline"),
					zap.String("pairName", params.pair.Name),
				)
				continue
			}
			cmd.queueManager.PublishKline(payload)
		}

		if len(externalExchangeKlines) < 1000 {
			break
		}

		lastExternalKline := externalExchangeKlines[len(externalExchangeKlines)-1]
		start = lastExternalKline.EndTime
		time.Sleep(1500 * time.Millisecond)
	}
	return nil
}
