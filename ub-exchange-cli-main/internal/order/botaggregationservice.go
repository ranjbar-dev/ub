package order

import (
	"context"
	"encoding/json"
	"exchange-go/internal/platform"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
)

const botAggregationListPrefix = "not-calculated:trades"

type BotAggregationData struct {
	TradeID     int64  `json:"tradeId"`
	PairID      int64  `json:"pairId"`
	RobotType   string `json:"robotType"`
	Amount      string `json:"amount"`
	Price       string `json:"price"`
	LastOrderID int64  `json:"lastOrderId"`
	//UserId      int64  `json:"userId"`
}

type FinalAggregatedOrderData struct {
	Type         string
	Amount       string
	ExchangeType string
	Price        string
	BuyAmount    string
	BuyPrice     string
	SellAmount   string
	SellPrice    string
}
type AggregationResult struct {
	LastTradeID         int64
	OrderIds            []string
	AggregatedOrderData FinalAggregatedOrderData
}

// BotAggregationService aggregates multiple small bot orders into a single
// consolidated order for efficient submission to an external exchange.
type BotAggregationService interface {
	// AddToList appends a bot trade's aggregation data to the Redis list for the trade's pair.
	AddToList(data BotAggregationData) error
	// DeleteList removes the entire aggregation list for the specified currency pair.
	DeleteList(pairID int64) error
	// GetAggregationResultForPair computes the net aggregated order from all pending
	// bot trades for the given pair, returning the consolidated amount, price, and direction.
	GetAggregationResultForPair(pairID int64) (AggregationResult, error)
	// GetListForPair returns all pending bot aggregation entries for the specified pair.
	GetListForPair(pairID int64) ([]BotAggregationData, error)
}

type botAggregationService struct {
	rc platform.RedisClient
}

func (b *botAggregationService) AddToList(botAggregationData BotAggregationData) error {
	ctx := context.Background()
	key := b.getListName(botAggregationData.PairID)
	data, err := json.Marshal(botAggregationData)
	finalData := string(data)
	if err != nil {
		return fmt.Errorf("AddToList: marshal data: %w", err)
	}
	_, err = b.rc.LRem(ctx, key, 0, finalData)
	if err != nil {
		return fmt.Errorf("AddToList: remove existing entry: %w", err)
	}
	_, err = b.rc.LPush(ctx, key, finalData)
	if err != nil {
		return fmt.Errorf("AddToList: push to list: %w", err)
	}
	return nil
}

func (b *botAggregationService) GetAggregationResultForPair(pairID int64) (result AggregationResult, err error) {
	finalOrderData := FinalAggregatedOrderData{
		ExchangeType: ExchangeTypeMarket,
	}

	data, err := b.GetListForPair(pairID)
	if err != nil {
		return result, nil
	}

	buyAmountDecimal := decimal.NewFromFloat(0)
	buyPriceDecimal := decimal.NewFromFloat(0)

	sellAmountDecimal := decimal.NewFromFloat(0)
	sellPriceDecimal := decimal.NewFromFloat(0)

	buySumDecimal := decimal.NewFromFloat(0)
	sellSumDecimal := decimal.NewFromFloat(0)

	lastTradeID := int64(0)
	var orderIds = make([]string, 0)

	for _, botData := range data {
		if botData.TradeID > lastTradeID {
			lastTradeID = botData.TradeID
		}
		OrderIDString := strconv.FormatInt(botData.LastOrderID, 10)
		orderIds = append(orderIds, OrderIDString)

		amountDecimal, _ := decimal.NewFromString(botData.Amount)
		priceDecimal, _ := decimal.NewFromString(botData.Price)
		totalDecimal := amountDecimal.Mul(priceDecimal)

		// if robotType is "buy" we must submit a sell order and if it is "sell" we must submit a buy order
		if strings.ToLower(botData.RobotType) == strings.ToLower(TypeBuy) {
			sellAmountDecimal = sellAmountDecimal.Add(amountDecimal)
			sellSumDecimal = sellSumDecimal.Add(totalDecimal)
		} else {
			buyAmountDecimal = buyAmountDecimal.Add(amountDecimal)
			buySumDecimal = buySumDecimal.Add(totalDecimal)
		}
	}

	if sellAmountDecimal.IsPositive() {
		sellPriceDecimal = sellSumDecimal.Div(sellAmountDecimal)
	}

	if buyAmountDecimal.IsPositive() {
		buyPriceDecimal = buySumDecimal.Div(buyAmountDecimal)
	}

	finalOrderData.BuyAmount = buyAmountDecimal.StringFixed(8)
	finalOrderData.BuyPrice = buyPriceDecimal.StringFixed(8)
	finalOrderData.SellAmount = sellAmountDecimal.StringFixed(8)
	finalOrderData.SellPrice = sellPriceDecimal.StringFixed(8)

	if buyAmountDecimal.GreaterThan(sellAmountDecimal) {
		finalOrderData.Type = TypeBuy
		finalOrderData.Amount = buyAmountDecimal.Sub(sellAmountDecimal).StringFixed(8)
		finalOrderData.Price = buyPriceDecimal.StringFixed(8)
	} else {
		finalOrderData.Type = TypeSell
		finalOrderData.Amount = sellAmountDecimal.Sub(buyAmountDecimal).StringFixed(8)
		finalOrderData.Price = sellPriceDecimal.StringFixed(8)
	}

	result = AggregationResult{
		LastTradeID:         lastTradeID,
		OrderIds:            orderIds,
		AggregatedOrderData: finalOrderData,
	}
	return result, nil

}

func (b *botAggregationService) DeleteList(pairID int64) error {
	ctx := context.Background()
	key := b.getListName(pairID)
	_, err := b.rc.Del(ctx, key)
	if err != nil && err != redis.Nil {
		return fmt.Errorf("DeleteList: delete key: %w", err)
	}

	return nil
}

func (b *botAggregationService) GetListForPair(pairID int64) ([]BotAggregationData, error) {
	result := make([]BotAggregationData, 0)

	ctx := context.Background()
	key := b.getListName(pairID)
	externalOrders, err := b.rc.LRange(ctx, key, 0, -1)
	if err != nil && err != redis.Nil {
		return result, fmt.Errorf("GetListForPair: lrange: %w", err)
	}

	if err == redis.Nil {
		return result, nil
	}

	if len(externalOrders) < 1 {
		return result, nil
	}

	for _, eo := range externalOrders {
		botData := &BotAggregationData{}
		err := json.Unmarshal([]byte(eo), botData)
		if err != nil && err != redis.Nil {
			return result, fmt.Errorf("GetListForPair: unmarshal entry: %w", err)
		}

		result = append(result, *botData)

	}
	return result, nil
}

func (b *botAggregationService) getListName(pairID int64) string {

	return botAggregationListPrefix + ":" + strconv.FormatInt(pairID, 10)
}

func NewBotAggregationService(rc platform.RedisClient) BotAggregationService {
	return &botAggregationService{
		rc: rc,
	}
}
