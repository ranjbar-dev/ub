package order

import (
	"exchange-go/internal/currency"
	"exchange-go/internal/platform"

	"go.uber.org/zap"
)

// TradeEventsHandler publishes trade events to downstream services such as
// bot aggregation and Centrifugo client notifications.
type TradeEventsHandler interface {
	// HandleTradesCreation processes newly created trades by forwarding bot trades
	// to the aggregation service for consolidation and external exchange submission.
	HandleTradesCreation(tradesDate []TradeData, pair currency.Pair)
}

type tradeEventsHandler struct {
	BotAggregationService BotAggregationService
	configs               platform.Configs
	logger                platform.Logger
}
type TradeData struct {
	Trade     Trade
	UserEmail string
	UserID    int
}

func (s *tradeEventsHandler) HandleTradesCreation(tradesData []TradeData, pair currency.Pair) {
	if pair.AggregationStatus == currency.AggregationStatusStop {
		return
	}
	for _, t := range tradesData {
		//todo this should not be hard coded here better be in config.yml
		if t.UserEmail == "rafsanjan@gmail.com" && s.configs.GetEnv() == platform.EnvProd {
			continue
		}
		if t.Trade.BotOrderType.Valid {
			var lastOrderID int64
			if t.Trade.SellOrderID.Valid {
				lastOrderID = t.Trade.SellOrderID.Int64
			}
			if t.Trade.BuyOrderID.Valid {
				lastOrderID = t.Trade.BuyOrderID.Int64
			}
			bat := BotAggregationData{
				TradeID:     t.Trade.ID,
				PairID:      t.Trade.PairID,
				RobotType:   t.Trade.BotOrderType.String,
				Amount:      t.Trade.Amount.String,
				Price:       t.Trade.Price.String,
				LastOrderID: lastOrderID,
				//UserId:      userId,
			}
			err := s.BotAggregationService.AddToList(bat)
			if err != nil {
				s.logger.Error2("error in adding to list", err,
					zap.String("service", "TradeEventsHandler"),
					zap.String("method", "HandleTradesCreation"),
					zap.Int64("tradeID", t.Trade.ID),
				)
			}
		}
	}
}

func NewTradeEventsHandler(botAggregationService BotAggregationService, configs platform.Configs, logger platform.Logger) TradeEventsHandler {
	return &tradeEventsHandler{
		BotAggregationService: botAggregationService,
		configs:               configs,
		logger:                logger,
	}
}
