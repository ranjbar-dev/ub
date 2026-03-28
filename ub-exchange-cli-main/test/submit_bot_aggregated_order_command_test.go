package test

import (
	"context"
	"encoding/json"
	"exchange-go/internal/command"
	"exchange-go/internal/di"
	"exchange-go/internal/externalexchange"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type SubmitBotOrderCmd struct {
	*suite.Suite
	submitBotOrderCmd command.ConsoleCommand
	db                *gorm.DB
	redisClient       *redis.Client
	userActor         *userActor
}

func (t *SubmitBotOrderCmd) SetupSuite() {
	container := getContainer()
	t.submitBotOrderCmd = container.Get(di.SubmitBotAggregatedOrderCommand).(command.ConsoleCommand)
	t.db = getDb()
	t.redisClient = getRedis()
	t.userActor = getUserActor()
}

func (t *SubmitBotOrderCmd) SetupTest() {
	btcUsdtQueue := "not-calculated:trades:1"
	_, err := t.redisClient.Del(context.Background(), btcUsdtQueue).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	t.db.Where("id > ?", int64(0)).Delete(externalexchange.Order{})

}

func (t *SubmitBotOrderCmd) TearDownTest() {}

func (t *SubmitBotOrderCmd) TearDownSuite() {}

func (t *SubmitBotOrderCmd) TestRun() {
	ctx := context.Background()
	// first put data in redis queue
	type botAggregationData struct {
		TradeID     int64  `json:"tradeId"`
		PairID      int64  `json:"pairId"`
		RobotType   string `json:"robotType"`
		Amount      string `json:"amount"`
		Price       string `json:"price"`
		LastOrderID int64  `json:"lastOrderId"`
		UserID      int64  `json:"userId"`
	}

	key := "not-calculated:trades:1"

	d1 := botAggregationData{
		TradeID:     1,
		PairID:      1,
		RobotType:   "BUY",
		Amount:      "0.1",
		Price:       "40000",
		LastOrderID: 1,
		UserID:      1,
	}

	data, err := json.Marshal(d1)
	finalData := string(data)
	if err != nil {
		t.Fail(err.Error())
	}
	_, err = t.redisClient.LPush(ctx, key, finalData).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	d2 := botAggregationData{
		TradeID:     2,
		PairID:      1,
		RobotType:   "SELL",
		Amount:      "0.05",
		Price:       "30000",
		LastOrderID: 2,
		UserID:      2,
	}

	data, err = json.Marshal(d2)
	finalData = string(data)
	if err != nil {
		t.Fail(err.Error())
	}
	_, err = t.redisClient.LPush(ctx, key, finalData).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	var flags []string
	t.submitBotOrderCmd.Run(ctx, flags)

	externalExchangeOrder := &externalexchange.Order{}

	err = t.db.First(&externalExchangeOrder).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Greater(t.T(), externalExchangeOrder.ID, int64(0))
	assert.Equal(t.T(), "0.05000000", externalExchangeOrder.Amount.String)
	assert.Equal(t.T(), "40000.00000000", externalExchangeOrder.Price.String)
	assert.Equal(t.T(), "0.05000000", externalExchangeOrder.BuyAmount.String)
	assert.Equal(t.T(), "30000.00000000", externalExchangeOrder.BuyPrice.String)
	assert.Equal(t.T(), "0.10000000", externalExchangeOrder.SellAmount.String)
	assert.Equal(t.T(), "40000.00000000", externalExchangeOrder.SellPrice.String)
	assert.Equal(t.T(), "2,1", externalExchangeOrder.OrderIds.String)
	assert.Equal(t.T(), int64(2), externalExchangeOrder.LastTradeID.Int64)
	assert.Equal(t.T(), int64(1), externalExchangeOrder.PairID.Int64)
	assert.Equal(t.T(), "SELL", externalExchangeOrder.Type.String)
	assert.Equal(t.T(), "MARKET", externalExchangeOrder.ExchangeType.String)
	assert.Equal(t.T(), "BOT", externalExchangeOrder.Source.String)
	assert.Equal(t.T(), "COMPLETED", externalExchangeOrder.Status.String)

	count, err := t.redisClient.Exists(ctx, key).Result()
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), int64(0), count)

	timestamp, err := t.redisClient.HGet(ctx, "live_data:pair_currency:BTC-USDT", "last_aggregation_time").Result()
	assert.Nil(t.T(), err)
	assert.NotEqual(t.T(), "", timestamp)

}

func TestSubmitBotOrderCmd(t *testing.T) {
	suite.Run(t, &SubmitBotOrderCmd{
		Suite: new(suite.Suite),
	})
}
