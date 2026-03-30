package command

import (
	"context"
	"exchange-go/internal/externalexchange"
	"exchange-go/internal/order"
	"exchange-go/internal/platform"
	"time"

	"go.uber.org/zap"
)

type retrieveExternalOrdersToRedisCmd struct {
	externalExchangeOrderService externalexchange.OrderService
	botAggregationService        order.BotAggregationService
	tradeService                 order.TradeService
	logger                       platform.Logger
}

func (cmd *retrieveExternalOrdersToRedisCmd) Run(ctx context.Context, flags []string) {
	cmd.logger.Info("start of retrieve external orders command")

	pairsLastTradeIds := cmd.externalExchangeOrderService.GetExternalExchangeOrdersLastTradeIds()

	for _, item := range pairsLastTradeIds {
		if item.TradeID == 0 {
			continue
		}

		existingOrders, err := cmd.botAggregationService.GetListForPair(item.PairID)
		if err != nil {
			cmd.logger.Error2("error getting order list for pair", err,
				zap.String("service", "retrieveExternalOrdersToRedisCmd"),
				zap.String("method", "Run"),
				zap.Int64("pairID", item.PairID),
			)
			continue
		}
		//twenty minutes because this cron is run every 20 minute
		twentyMinAgo := time.Now().Add(-21 * time.Minute)
		trades := cmd.tradeService.GetBotTradesByIDAndCreatedAtGreaterThan(item.PairID, item.TradeID, twentyMinAgo)
		for _, t := range trades {
			//check with existing orders
			alreadyExistsInRedisQueue := false
			for _, o := range existingOrders {
				if t.ID == o.TradeID {
					alreadyExistsInRedisQueue = true
					break
				}
			}

			if alreadyExistsInRedisQueue {
				//doing nothing since the order is already in redis
				continue
			}

			lastOrderID := int64(0)
			if t.BuyOrderID.Valid {
				lastOrderID = t.BuyOrderID.Int64
			} else {
				lastOrderID = t.SellOrderID.Int64
			}
			data := order.BotAggregationData{
				TradeID:     t.ID,
				PairID:      t.PairID,
				RobotType:   t.BotOrderType.String,
				Amount:      t.Amount.String,
				Price:       t.Price.String,
				LastOrderID: lastOrderID,
				//UserId:      userId,
			}
			err := cmd.botAggregationService.AddToList(data)
			if err != nil {
				cmd.logger.Error2("error in adding to list", err,
					zap.String("service", "retrieveExternalOrdersToRedisCmd"),
					zap.String("method", "Run"),
					zap.Int64("pairID", data.PairID),
				)
				continue
			}

		}

	}

	cmd.logger.Info("end of retrieve external orders command")
}

func NewRetrieveExternalOrdersToRedisCmd(externalExchangeOrderService externalexchange.OrderService,
	botAggregationService order.BotAggregationService, tradeService order.TradeService,
	logger platform.Logger) ConsoleCommand {
	return &retrieveExternalOrdersToRedisCmd{
		externalExchangeOrderService: externalExchangeOrderService,
		botAggregationService:        botAggregationService,
		tradeService:                 tradeService,
		logger:                       logger,
	}

}
