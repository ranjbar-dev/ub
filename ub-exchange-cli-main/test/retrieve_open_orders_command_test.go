package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"exchange-go/internal/command"
	"exchange-go/internal/di"
	"exchange-go/internal/engine"
	"exchange-go/internal/externalexchange"
	"exchange-go/internal/order"
	"fmt"
	"strconv"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type RetrieveOpenOrdersCmd struct {
	*suite.Suite
	retrieveOpenOrdersCmd command.ConsoleCommand
	db                    *gorm.DB
	redisClient           *redis.Client
	userActor             *userActor
}

func (t *RetrieveOpenOrdersCmd) SetupSuite() {
	container := getContainer()
	t.retrieveOpenOrdersCmd = container.Get(di.RetrieveOpenOrdersToRedisCommand).(command.ConsoleCommand)
	t.db = getDb()
	t.redisClient = getRedis()
	t.userActor = getUserActor()
}

func (t *RetrieveOpenOrdersCmd) SetupTest() {}

func (t *RetrieveOpenOrdersCmd) TearDownTest() {}

func (t *RetrieveOpenOrdersCmd) TearDownSuite() {
	t.db.Where("id > ?", int64(0)).Delete(externalexchange.Order{})
	t.db.Where("id > ?", int64(0)).Delete(order.Order{})
	//empty redis queue
	allOrdersQueue := engine.QueueName
	_, err := t.redisClient.Del(context.Background(), allOrdersQueue).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	orderBookBidQueue := engine.OrderBookBidsPrefix + "BTC-USDT"
	orderBookAskQueue := engine.OrderBookAsksPrefix + "BTC-USDT"
	_, err = t.redisClient.Del(context.Background(), orderBookBidQueue).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	_, err = t.redisClient.Del(context.Background(), orderBookAskQueue).Result()
	if err != nil {
		t.Fail(err.Error())
	}
}

func (t *RetrieveOpenOrdersCmd) TestRun() {
	ctx := context.Background()

	//insert data in db
	//limit order
	o1 := &order.Order{
		ID:                 1,
		UserID:             t.userActor.ID,
		Type:               order.TypeBuy,
		ExchangeType:       order.ExchangeTypeLimit,
		Price:              sql.NullString{String: "50000.00000000", Valid: true},
		Status:             order.StatusOpen,
		DemandedAmount:     sql.NullString{String: "0.10000000", Valid: true},
		DemandedCoin:       "BTC",
		PayedByAmount:      sql.NullString{String: "5000.00000000", Valid: true},
		PayedByCoin:        "USDT",
		PairID:             1,
		CurrentMarketPrice: sql.NullString{String: "50000.00000000", Valid: true},
	}

	//market order
	o2 := &order.Order{
		ID:                 2,
		UserID:             t.userActor.ID,
		Type:               order.TypeSell,
		ExchangeType:       order.ExchangeTypeMarket,
		Status:             order.StatusOpen,
		DemandedAmount:     sql.NullString{String: "5000.00000000", Valid: true},
		DemandedCoin:       "USDT",
		PayedByAmount:      sql.NullString{String: "0.10000000", Valid: true},
		PayedByCoin:        "BTC",
		PairID:             1,
		CurrentMarketPrice: sql.NullString{String: "50000.00000000", Valid: true},
	}

	//stop order
	o3 := &order.Order{
		ID:             3,
		UserID:         t.userActor.ID,
		Type:           order.TypeBuy,
		ExchangeType:   order.ExchangeTypeLimit,
		Price:          sql.NullString{String: "51000.00000000", Valid: true},
		StopPointPrice: sql.NullString{String: "50000.00000000", Valid: true},
		Status:         order.StatusOpen,
		DemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		DemandedCoin:   "BTC",
		PayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		PayedByCoin:    "USDT",
		PairID:         1,
	}

	//limit order which already is in redis orderbook
	o4 := &order.Order{
		ID:                 4,
		UserID:             t.userActor.ID,
		Type:               order.TypeBuy,
		ExchangeType:       order.ExchangeTypeLimit,
		Price:              sql.NullString{String: "50000.00000000", Valid: true},
		Status:             order.StatusOpen,
		DemandedAmount:     sql.NullString{String: "0.10000000", Valid: true},
		DemandedCoin:       "BTC",
		PayedByAmount:      sql.NullString{String: "5000.00000000", Valid: true},
		PayedByCoin:        "USDT",
		PairID:             1,
		CurrentMarketPrice: sql.NullString{String: "50000.00000000", Valid: true},
	}

	orders := []*order.Order{o1, o2, o3, o4}
	err := t.db.Create(orders).Error
	if err != nil {
		t.Fail(err.Error())
	}

	//empty redis queue
	allOrdersQueue := engine.QueueName
	_, err = t.redisClient.Del(context.Background(), allOrdersQueue).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	orderBookBidQueue := engine.OrderBookBidsPrefix + "BTC-USDT"
	orderBookAskQueue := engine.OrderBookAsksPrefix + "BTC-USDT"
	_, err = t.redisClient.Del(context.Background(), orderBookBidQueue).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	_, err = t.redisClient.Del(context.Background(), orderBookAskQueue).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	//put the o4 in orderbook
	engineOrder4IDString := fmt.Sprintf("%011d", o4.ID)
	eo4 := engine.Order{
		Pair:              "BTC-USDT",
		ID:                engineOrder4IDString,
		Side:              "bid",
		Quantity:          "0.10000000",
		Price:             "50000.00000000",
		Timestamp:         o1.CreatedAt.Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MarketPrice:       "50000.00000000",
		MinThresholdPrice: "48000.00000000",
		MaxThresholdPrice: "52000.00000000",
	}
	engineOrderData4, err := json.Marshal(&eo4)
	if err != nil {
		t.Fail(err.Error())
	}
	engineOrderBookData4, err := eo4.MarshalForOrderbook()
	if err != nil {
		t.Fail(err.Error())
	}
	member := &redis.Z{
		Score:  50000,
		Member: string(engineOrderBookData4),
	}

	_, err = t.redisClient.ZAdd(context.Background(), orderBookBidQueue, member).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	//call command
	var flags []string
	t.retrieveOpenOrdersCmd.Run(ctx, flags)

	engineOrder1IDString := fmt.Sprintf("%011d", o1.ID)
	eo1 := engine.Order{
		Pair:              "BTC-USDT",
		ID:                engineOrder1IDString,
		Side:              "bid",
		Quantity:          "0.10000000",
		Price:             "50000.00000000",
		Timestamp:         o1.CreatedAt.Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MarketPrice:       "50000.00000000",
		MinThresholdPrice: "48000.00000000",
		MaxThresholdPrice: "52000.00000000",
	}

	engineOrderData1, err := json.Marshal(&eo1)
	if err != nil {
		t.Fail(err.Error())
	}

	pos, err := t.redisClient.LPos(context.Background(), allOrdersQueue, string(engineOrderData1), redis.LPosArgs{}).Result()
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), pos, int64(1)) //will be shifted one because of the other order, so the pos is second index which is 1

	engineOrder2IDString := fmt.Sprintf("%011d", o2.ID)
	eo2 := engine.Order{
		Pair:              "BTC-USDT",
		ID:                engineOrder2IDString,
		Side:              "ask",
		Quantity:          "0.10000000",
		Price:             "",
		Timestamp:         o2.CreatedAt.Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MarketPrice:       "50000.00000000",
		MinThresholdPrice: "48000.00000000",
		MaxThresholdPrice: "52000.00000000",
	}

	engineOrderData2, err := json.Marshal(&eo2)
	if err != nil {
		t.Fail(err.Error())
	}
	pos, err = t.redisClient.LPos(context.Background(), allOrdersQueue, string(engineOrderData2), redis.LPosArgs{}).Result()
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), pos, int64(0))

	pos, err = t.redisClient.LPos(context.Background(), allOrdersQueue, string(engineOrderData4), redis.LPosArgs{}).Result()
	assert.Equal(t.T(), redis.Nil, err)

	stopOrderQueue := order.StopOrderQueuePrefix + "BUY:" + "BTC-USDT"
	stopOrderID := strconv.FormatInt(o3.ID, 10)
	rank, err := t.redisClient.ZRank(context.Background(), stopOrderQueue, stopOrderID).Result()
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), rank, int64(0)) //because it is the only one the rank index is zero
}

func TestRetrieveOpenOrdersCmd(t *testing.T) {
	suite.Run(t, &RetrieveOpenOrdersCmd{
		Suite: new(suite.Suite),
	})
}
