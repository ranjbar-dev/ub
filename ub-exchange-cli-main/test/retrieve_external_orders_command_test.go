package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"exchange-go/internal/command"
	"exchange-go/internal/di"
	"exchange-go/internal/externalexchange"
	"exchange-go/internal/order"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type RetrieveExternalOrdersCmd struct {
	*suite.Suite
	retrieveExternalOrdersCmd command.ConsoleCommand
	db                        *gorm.DB
	redisClient               *redis.Client
	userActor                 *userActor
}

func (t *RetrieveExternalOrdersCmd) SetupSuite() {
	container := getContainer()
	t.retrieveExternalOrdersCmd = container.Get(di.RetrieveExternalOrdersToRedisCommand).(command.ConsoleCommand)
	t.db = getDb()
	t.redisClient = getRedis()
	t.userActor = getUserActor()
}

func (t *RetrieveExternalOrdersCmd) SetupTest() {
	t.db.Where("id > ?", int64(0)).Delete(externalexchange.Order{})
	t.db.Where("id > ?", int64(0)).Delete(order.Order{})
	btcUsdtQueue := "not-calculated:trades:1"
	_, err := t.redisClient.Del(context.Background(), btcUsdtQueue).Result()
	if err != nil {
		t.Fail(err.Error())
	}
}

func (t *RetrieveExternalOrdersCmd) TearDownTest() {
	t.db.Where("id > ?", int64(0)).Delete(externalexchange.Order{})
	t.db.Where("id > ?", int64(0)).Delete(order.Trade{})
	t.db.Where("id > ?", int64(0)).Delete(order.Order{})
}

func (t *RetrieveExternalOrdersCmd) TearDownSuite() {}

func (t *RetrieveExternalOrdersCmd) TestRun() {
	ctx := context.Background()

	o1 := &order.Order{
		ID:                 1,
		UserID:             t.userActor.ID,
		Type:               order.TypeBuy,
		ExchangeType:       order.ExchangeTypeLimit,
		Price:              sql.NullString{String: "50000", Valid: true},
		Status:             order.StatusFilled,
		DemandedAmount:     sql.NullString{String: "0.1", Valid: true},
		PayedByAmount:      sql.NullString{String: "5000", Valid: true},
		PairID:             1,
		Level:              sql.NullInt64{Int64: 1, Valid: true},
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		CurrentMarketPrice: sql.NullString{String: "50000", Valid: true},
	}

	o2 := &order.Order{
		ID:                 2,
		UserID:             t.userActor.ID,
		Type:               order.TypeSell,
		ExchangeType:       order.ExchangeTypeLimit,
		Price:              sql.NullString{String: "50000", Valid: true},
		Status:             order.StatusFilled,
		DemandedAmount:     sql.NullString{String: "0.1", Valid: true},
		PayedByAmount:      sql.NullString{String: "5000", Valid: true},
		PairID:             1,
		Level:              sql.NullInt64{Int64: 1, Valid: true},
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		CurrentMarketPrice: sql.NullString{String: "50000", Valid: true},
	}

	orders := []*order.Order{o1, o2}
	err := t.db.Create(orders).Error
	if err != nil {
		t.Fail(err.Error())
	}

	//insert data in db
	t1 := &order.Trade{
		ID:           2,
		Price:        sql.NullString{String: "50000", Valid: true},
		Amount:       sql.NullString{String: "0.1", Valid: true},
		PairID:       1,
		BuyOrderID:   sql.NullInt64{Int64: 1, Valid: true},
		SellOrderID:  sql.NullInt64{},
		BotOrderType: sql.NullString{String: externalexchange.TypeSell, Valid: true},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	t2 := &order.Trade{
		ID:           3,
		Price:        sql.NullString{String: "50000", Valid: true},
		Amount:       sql.NullString{String: "0.1", Valid: true},
		PairID:       1,
		BuyOrderID:   sql.NullInt64{},
		SellOrderID:  sql.NullInt64{Int64: 2, Valid: true},
		BotOrderType: sql.NullString{String: externalexchange.TypeBuy, Valid: true},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	trades := []*order.Trade{t1, t2}
	err = t.db.Create(trades).Error
	if err != nil {
		t.Fail(err.Error())
	}

	eo := &externalexchange.Order{
		PairID:      sql.NullInt64{Int64: 1, Valid: true},
		LastTradeID: sql.NullInt64{Int64: 1, Valid: true},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = t.db.Create(eo).Error
	if err != nil {
		t.Fail(err.Error())
	}

	//put one of the trade in redis queue
	btcUsdtQueue := "not-calculated:trades:1"
	data1 := `{"tradeId":2,"pairId":1,"robotType":"SELL","amount":"0.10000000","price":"50000.00000000","lastOrderId":1}`
	_, err = t.redisClient.LPush(ctx, btcUsdtQueue, data1).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	//call command
	var flags []string
	t.retrieveExternalOrdersCmd.Run(ctx, flags)

	externalOrders, err := t.redisClient.LRange(ctx, btcUsdtQueue, 0, -1).Result()
	assert.Equal(t.T(), 2, len(externalOrders))
	for _, eo := range externalOrders {
		botData := &order.BotAggregationData{}
		err := json.Unmarshal([]byte(eo), botData)
		if err != nil {
			t.Fail(err.Error())
		}

		switch botData.TradeID {
		case int64(2):
			assert.Equal(t.T(), int64(1), botData.LastOrderID)
		case int64(3):
			assert.Equal(t.T(), int64(2), botData.LastOrderID)
		default:
			t.Fail("we should not be in default case")
		}

	}

}

func TestRetrieveExternalOrdersCmd(t *testing.T) {
	suite.Run(t, &RetrieveExternalOrdersCmd{
		Suite: new(suite.Suite),
	})
}
