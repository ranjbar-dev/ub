package command

import (
	"context"
	"exchange-go/internal/currency"
	"exchange-go/internal/externalexchange"
	"exchange-go/internal/livedata"
	"exchange-go/internal/order"
	"exchange-go/internal/platform"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type submitBotAggregatedOrderCmd struct {
	currencyService              currency.Service
	botAggregationService        order.BotAggregationService
	liveDataService              livedata.Service
	externalExchangeOrderService externalexchange.OrderService
	logger                       platform.Logger
}

func (cmd *submitBotAggregatedOrderCmd) Run(ctx context.Context, flags []string) {
	fmt.Println("start of submit bot orders command")
	pairs := cmd.currencyService.GetActivePairCurrenciesList()
	for _, pair := range pairs {
		if pair.AggregationStatus == currency.AggregationStatusPause || pair.AggregationStatus == currency.AggregationStatusStop {
			continue
		}
		shouldSend, err := cmd.shouldSendToExternalExchange(pair)
		if err != nil {
			cmd.logger.Error2("error checking if we should send order to external exchange", err,
				zap.String("service", "submitBotAggregatedOrderCmd"),
				zap.String("method", "Run"),
				zap.String("pairName", pair.Name),
			)
			continue
		}
		if !shouldSend {
			continue
		}

		result, err := cmd.botAggregationService.GetAggregationResultForPair(pair.ID)
		if err != nil {
			cmd.logger.Error2("error getting aggregation result for pair", err,
				zap.String("service", "submitBotAggregatedOrderCmd"),
				zap.String("method", "Run"),
				zap.String("pairName", pair.Name),
			)
			continue
		}

		if result.AggregatedOrderData.Amount == "" {
			continue
		}

		amountDecimal, err := decimal.NewFromString(result.AggregatedOrderData.Amount)
		if err != nil {
			cmd.logger.Error2("error converting amount to decimal", err,
				zap.String("service", "submitBotAggregatedOrderCmd"),
				zap.String("method", "Run"),
				zap.String("pairName", pair.Name),
				zap.String("amount", result.AggregatedOrderData.Amount),
			)
			continue
		}
		if !amountDecimal.IsPositive() {
			continue
		}

		params := externalexchange.BotOrderParams{
			PairID:       pair.ID,
			PairName:     pair.Name,
			Type:         result.AggregatedOrderData.Type,
			ExchangeType: result.AggregatedOrderData.ExchangeType,
			Amount:       result.AggregatedOrderData.Amount,
			Price:        result.AggregatedOrderData.Price,
			BuyAmount:    result.AggregatedOrderData.BuyAmount,
			BuyPrice:     result.AggregatedOrderData.BuyPrice,
			SellAmount:   result.AggregatedOrderData.SellAmount,
			SellPrice:    result.AggregatedOrderData.SellPrice,
			LastTradeID:  result.LastTradeID,
			OrderIds:     result.OrderIds,
		}

		_, err = cmd.externalExchangeOrderService.CreateExternalExchangeOrderForBot(params)
		if err != nil {
			cmd.logger.Error2("error creating external exchange order for bot", err,
				zap.String("service", "submitBotAggregatedOrderCmd"),
				zap.String("method", "Run"),
				zap.String("pairName", pair.Name),
				zap.Int64("lastTradeID", result.LastTradeID),
			)
		}

		//deleting list whether submiting order was successful or not
		err = cmd.botAggregationService.DeleteList(pair.ID)
		if err != nil {
			cmd.logger.Error2("error deleting list from redis", err,
				zap.String("service", "submitBotAggregatedOrderCmd"),
				zap.String("method", "Run"),
				zap.String("pairName", pair.Name),
				zap.Int64("lastTradeID", result.LastTradeID),
			)
			continue
		}

		now := time.Now().Unix()
		err = cmd.liveDataService.SetLastAggregationTime(ctx, pair.Name, now)
		if err != nil {
			cmd.logger.Error2("error setting last aggregation time in redis", err,
				zap.String("service", "submitBotAggregatedOrderCmd"),
				zap.String("method", "Run"),
				zap.String("pairName", pair.Name),
			)
		}

	}

	fmt.Println("end of submit bot orders command")
}

func (cmd *submitBotAggregatedOrderCmd) shouldSendToExternalExchange(pair currency.Pair) (bool, error) {
	ctx := context.Background()
	lastAggregationTimeString, err := cmd.liveDataService.GetLastAggregationTime(ctx, pair.Name)
	if err != nil && err != redis.Nil {
		return false, err
	}
	if err == redis.Nil {
		return true, nil
	}
	botOrdersAggregationTime := pair.BotOrdersAggregationTime
	if !botOrdersAggregationTime.Valid || botOrdersAggregationTime.Int64 == int64(0) {
		return true, nil
	}
	lastAggregationTime, err := strconv.ParseInt(lastAggregationTimeString, 10, 64)
	if err != nil {
		return false, err
	}
	if time.Now().Unix()-lastAggregationTime >= botOrdersAggregationTime.Int64 {
		return true, nil
	}
	return false, nil
}

func NewSubmitBotAggregatedCmd(currencyService currency.Service, botAggregationService order.BotAggregationService,
	liveDataService livedata.Service, externalExchangeOrderService externalexchange.OrderService,
	logger platform.Logger) ConsoleCommand {
	cmd := &submitBotAggregatedOrderCmd{
		currencyService:              currencyService,
		botAggregationService:        botAggregationService,
		liveDataService:              liveDataService,
		externalExchangeOrderService: externalExchangeOrderService,
		logger:                       logger,
	}
	return cmd

}
