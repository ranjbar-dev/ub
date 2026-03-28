package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"exchange-go/internal/api"
	"exchange-go/internal/di"
	"exchange-go/internal/engine"
	"exchange-go/internal/order"
	"exchange-go/internal/response"
	"exchange-go/internal/userbalance"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type OrderCancelTests struct {
	*suite.Suite
	httpServer  http.Handler
	db          *gorm.DB
	redisClient *redis.Client
	userActor   *userActor
}

func (t *OrderCancelTests) SetupSuite() {
	container := getContainer()
	t.httpServer = container.Get(di.HTTPServer).(api.HTTPServer).GetEngine()
	t.db = getDb()
	t.redisClient = getRedis()
	t.userActor = getUserActor()
}

func (t *OrderCancelTests) TearDownSuite() {
}

func (t *OrderCancelTests) SetupTest() {

}

func (t *OrderCancelTests) TearDownTest() {
	t.db.Where("user_id = ?", t.userActor.ID).Delete(userbalance.UserBalance{})
}

func (t *OrderCancelTests) TestOrderCancel_Fail_OrderDoesNotExist() {
	data := `{"order_id":12231111}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/order/cancel", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := response.APIResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.T().Error(err)
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "order not found", result.Message)

}

func (t *OrderCancelTests) TestOrderCancel_Fail_StatusIsNotOpen() {
	now := time.Now()
	o := &order.Order{
		UserID:         t.userActor.ID,
		Type:           order.TypeBuy,
		ExchangeType:   order.ExchangeTypeLimit,
		Price:          sql.NullString{String: "50000", Valid: true},
		Status:         order.StatusFilled,
		DemandedAmount: sql.NullString{String: "0.1", Valid: true},
		PayedByAmount:  sql.NullString{String: "5000", Valid: true},
		PairID:         1,
		Level:          sql.NullInt64{Int64: 1, Valid: true},
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	err := t.db.Create(o).Error
	if err != nil {
		t.Fail(err.Error())
	}
	orderIDString := strconv.FormatInt(o.ID, 10)

	data := "{\"order_id\":" + orderIDString + "}"

	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/order/cancel", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "order status is not open", result.Message)
}

func (t *OrderCancelTests) TestOrderCancel_Buy_Successful_StopOrder() {
	//empty  redis
	queue := order.StopOrderQueuePrefix + "BUY:" + "BTC-USDT"
	_, err := t.redisClient.Del(context.Background(), queue).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	o := &order.Order{
		UserID:         t.userActor.ID,
		Type:           order.TypeBuy,
		ExchangeType:   order.ExchangeTypeLimit,
		Price:          sql.NullString{String: "51000", Valid: true},
		Status:         order.StatusOpen,
		DemandedAmount: sql.NullString{String: "0.1", Valid: true},
		PayedByAmount:  sql.NullString{String: "5100", Valid: true},
		PairID:         1,
		Level:          sql.NullInt64{Int64: 1, Valid: true},
		StopPointPrice: sql.NullString{String: "50000", Valid: true},
	}
	//insert order in db
	err = t.db.Create(o).Error
	if err != nil {
		t.Fail(err.Error())
	}
	orderIDString := strconv.FormatInt(o.ID, 10)

	//insert order in redis
	data := &redis.Z{
		Score:  50000, //is the stop point price
		Member: orderIDString,
	}

	t.redisClient.ZAdd(context.Background(), queue, data)

	//insert userBalance for user
	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "10000.00",
		FrozenAmount:  "6000.00",
	}

	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.01",
		FrozenAmount:  "0",
	}

	t.db.Create(usdtUb)
	t.db.Create(btcUb)

	params := "{\"order_id\":" + orderIDString + "}"

	res := httptest.NewRecorder()
	body := []byte(params)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/order/cancel", bytes.NewReader(body))

	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)

	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	updatedOrder := &order.Order{}
	err = t.db.Where(order.Order{ID: o.ID}).First(updatedOrder).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), order.StatusCanceled, updatedOrder.Status)

	_, err = t.redisClient.ZRank(context.Background(), queue, orderIDString).Result()
	assert.Equal(t.T(), redis.Nil, err)

	updatedUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "900.00000000", updatedUsdtUb.FrozenAmount)
}

func (t *OrderCancelTests) TestOrderCancel_Sell_Successful_StopOrder() {
	//empty  redis
	queue := order.StopOrderQueuePrefix + "SELL:" + "BTC-USDT"
	_, err := t.redisClient.Del(context.Background(), queue).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	o := &order.Order{
		UserID:         t.userActor.ID,
		Type:           order.TypeSell,
		ExchangeType:   order.ExchangeTypeLimit,
		Price:          sql.NullString{String: "49000", Valid: true},
		Status:         order.StatusOpen,
		DemandedAmount: sql.NullString{String: "49000", Valid: true},
		PayedByAmount:  sql.NullString{String: "0.1", Valid: true},
		PairID:         1,
		Level:          sql.NullInt64{Int64: 1, Valid: true},
		StopPointPrice: sql.NullString{String: "50000", Valid: true},
	}
	//insert order in db
	err = t.db.Create(o).Error
	if err != nil {
		t.T().Error("can not insert order in database in order cancel test")
	}
	orderIDString := strconv.FormatInt(o.ID, 10)

	//insert order in redis
	data := &redis.Z{
		Score:  50000, //is the stop point price
		Member: orderIDString,
	}

	t.redisClient.ZAdd(context.Background(), queue, data)

	//insert userBalance for user
	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "10000.00",
		FrozenAmount:  "6000.00",
	}

	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.3",
		FrozenAmount:  "0.2",
	}

	t.db.Create(usdtUb)
	t.db.Create(btcUb)

	params := "{\"order_id\":" + orderIDString + "}"

	res := httptest.NewRecorder()
	body := []byte(params)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/order/cancel", bytes.NewReader(body))

	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)

	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	updatedOrder := &order.Order{}
	err = t.db.Where(order.Order{ID: o.ID}).First(updatedOrder).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), order.StatusCanceled, updatedOrder.Status)

	_, err = t.redisClient.ZRank(context.Background(), queue, orderIDString).Result()
	assert.Equal(t.T(), redis.Nil, err)

	updatedBtcUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: btcUb.ID}).First(updatedBtcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "0.10000000", updatedBtcUb.FrozenAmount)
}

func (t *OrderCancelTests) TestOrderCancel_Buy_Successful_LimitOrder() {
	//we empty the redis queues first
	allOrdersQueue := engine.QueueName
	zQueue := engine.OrderBookBidsPrefix + "BTC-USDT"
	_, err := t.redisClient.Del(context.Background(), allOrdersQueue).Result()
	if err != nil {
		t.Fail(err.Error())
	}
	_, err = t.redisClient.Del(context.Background(), zQueue).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	//insert data in db
	now := time.Now()
	o := &order.Order{
		UserID:             t.userActor.ID,
		Type:               order.TypeBuy,
		ExchangeType:       order.ExchangeTypeLimit,
		Price:              sql.NullString{String: "50000", Valid: true},
		Status:             order.StatusOpen,
		DemandedAmount:     sql.NullString{String: "0.1", Valid: true},
		PayedByAmount:      sql.NullString{String: "5000", Valid: true},
		PairID:             1,
		Level:              sql.NullInt64{Int64: 1, Valid: true},
		CreatedAt:          now,
		UpdatedAt:          now,
		CurrentMarketPrice: sql.NullString{String: "50000", Valid: true},
	}
	err = t.db.Create(o).Error
	if err != nil {
		t.Fail(err.Error())
	}
	orderIDString := strconv.FormatInt(o.ID, 10)
	engineOrderIDString := fmt.Sprintf("%011d", o.ID)

	//insert order in redis
	eo := engine.Order{
		Pair:              "BTC-USDT",
		ID:                engineOrderIDString,
		Side:              "bid",
		Quantity:          "0.10000000",
		Price:             "50000",
		Timestamp:         now.Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MarketPrice:       "50000.00000000",
		MinThresholdPrice: "48000.00000000",
		MaxThresholdPrice: "52000.00000000",
	}

	engineOrderData, err := json.Marshal(&eo)
	if err != nil {
		t.Fail(err.Error())
	}
	t.redisClient.LPush(context.Background(), allOrdersQueue, string(engineOrderData))

	engineOrderBookData, err := eo.MarshalForOrderbook()
	if err != nil {
		t.Fail(err.Error())
	}
	zData := &redis.Z{
		Score:  50000,
		Member: string(engineOrderBookData),
	}

	t.redisClient.ZAdd(context.Background(), zQueue, zData)

	//insert userBalance for user
	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "10000.00",
		FrozenAmount:  "6000.00",
	}

	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.01",
		FrozenAmount:  "0",
	}

	t.db.Create(usdtUb)
	t.db.Create(btcUb)

	params := "{\"order_id\":" + orderIDString + "}"

	res := httptest.NewRecorder()
	body := []byte(params)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/order/cancel", bytes.NewReader(body))

	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)

	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	updatedOrder := &order.Order{}
	err = t.db.Where(userbalance.UserBalance{ID: o.ID}).First(updatedOrder).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), order.StatusCanceled, updatedOrder.Status)

	//checking redis
	_, err = t.redisClient.ZRank(context.Background(), zQueue, string(engineOrderBookData)).Result()
	assert.Equal(t.T(), redis.Nil, err)

	pos, err := t.redisClient.LPos(context.Background(), allOrdersQueue, string(engineOrderData), redis.LPosArgs{}).Result()
	assert.Equal(t.T(), redis.Nil, err)
	assert.Equal(t.T(), int64(0), pos)

	updatedUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "1000.00000000", updatedUsdtUb.FrozenAmount)

}

func (t *OrderCancelTests) TestOrderCancel_Sell_Successful_LimitOrder() {
	//we empty the redis queues first
	allOrdersQueue := engine.QueueName
	zQueue := engine.OrderBookAsksPrefix + "BTC-USDT"
	_, err := t.redisClient.Del(context.Background(), allOrdersQueue).Result()
	if err != nil {
		t.Fail(err.Error())
	}
	_, err = t.redisClient.Del(context.Background(), zQueue).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	//insert data in db
	now := time.Now()
	o := &order.Order{
		UserID:             t.userActor.ID,
		Type:               order.TypeSell,
		ExchangeType:       order.ExchangeTypeLimit,
		Price:              sql.NullString{String: "50000", Valid: true},
		Status:             order.StatusOpen,
		DemandedAmount:     sql.NullString{String: "5000", Valid: true},
		PayedByAmount:      sql.NullString{String: "0.1", Valid: true},
		PairID:             1,
		Level:              sql.NullInt64{Int64: 1, Valid: true},
		CreatedAt:          now,
		UpdatedAt:          now,
		CurrentMarketPrice: sql.NullString{String: "50000", Valid: true},
	}
	err = t.db.Create(o).Error
	if err != nil {
		t.Fail(err.Error())
	}
	orderIDString := strconv.FormatInt(o.ID, 10)
	engineOrderIDString := fmt.Sprintf("%011d", o.ID)
	//insert order in redis
	eo := engine.Order{
		Pair:              "BTC-USDT",
		ID:                engineOrderIDString,
		Side:              "ask",
		Quantity:          "0.10000000",
		Price:             "50000",
		Timestamp:         now.Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MarketPrice:       "50000.00000000",
		MinThresholdPrice: "48000.00000000",
		MaxThresholdPrice: "52000.00000000",
	}

	engineOrderData, err := json.Marshal(&eo)
	if err != nil {
		t.Fail(err.Error())
	}
	t.redisClient.LPush(context.Background(), allOrdersQueue, string(engineOrderData))

	engineOrderBookData, err := eo.MarshalForOrderbook()
	if err != nil {
		t.Fail(err.Error())
	}
	zData := &redis.Z{
		Score:  50000,
		Member: string(engineOrderBookData),
	}

	t.redisClient.ZAdd(context.Background(), zQueue, zData)

	//insert userBalance for user
	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "10000.00",
		FrozenAmount:  "6000.00",
	}

	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.3",
		FrozenAmount:  "0.2",
	}

	t.db.Create(usdtUb)
	t.db.Create(btcUb)

	params := "{\"order_id\":" + orderIDString + "}"

	res := httptest.NewRecorder()
	body := []byte(params)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/order/cancel", bytes.NewReader(body))

	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)

	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	updatedOrder := &order.Order{}
	err = t.db.Where(order.Order{ID: o.ID}).First(updatedOrder).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), order.StatusCanceled, updatedOrder.Status)

	//checking redis

	_, err = t.redisClient.ZRank(context.Background(), zQueue, string(engineOrderBookData)).Result()
	assert.Equal(t.T(), redis.Nil, err)

	pos, err := t.redisClient.LPos(context.Background(), allOrdersQueue, string(engineOrderData), redis.LPosArgs{}).Result()
	assert.Equal(t.T(), redis.Nil, err)
	assert.Equal(t.T(), int64(0), pos)

	updatedBtcUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: btcUb.ID}).First(updatedBtcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "0.10000000", updatedBtcUb.FrozenAmount)

}

func (t *OrderCancelTests) TestOrderCreateAndCancelFrequently_Limit_100Buy_And_100Sell() {
	//setting price in redis
	t.redisClient.HMSet(context.Background(), "live_data:pair_currency:BTC-USDT", "price", "50000")

	//insert user balance
	usdtUb := &userbalance.UserBalance{
		UserID:       t.userActor.ID,
		CoinID:       1, //for usdt from currency seed
		Status:       userbalance.StatusEnabled,
		Amount:       "10000.00000000",
		FrozenAmount: "0.00000000",
	}

	btcUb := &userbalance.UserBalance{
		UserID: t.userActor.ID,
		CoinID: 2, //for btc from currency seed

		Status:       userbalance.StatusEnabled,
		Amount:       "1.00000000",
		FrozenAmount: "0.00000000",
	}
	t.db.Create(usdtUb)
	t.db.Create(btcUb)

	token := "Bearer " + t.userActor.Token
	wg := sync.WaitGroup{}
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			//create and cancel sleep randomly between
			//call order api
			data := `{"type":"buy","amount":"0.01","exchange_type":"limit","pair_currency_id":1,"price":"10000"}`
			res := httptest.NewRecorder()
			body := []byte(data)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/order/create", bytes.NewReader(body))
			req.Header.Set("Authorization", token)
			t.httpServer.ServeHTTP(res, req)
			result := struct {
				Status  bool
				Message string
				Data    order.CreateOrderResponse
			}{}
			err := json.Unmarshal(res.Body.Bytes(), &result)
			if err != nil {
				t.Fail(err.Error())
			}
			orderID := result.Data.ID

			//there with this id we cancel the order
			orderIDString := strconv.FormatInt(orderID, 10)
			params := "{\"orderId\":" + orderIDString + "}"
			cancelRes := httptest.NewRecorder()
			cancelBody := []byte(params)
			cancelReq := httptest.NewRequest(http.MethodPost, "/api/v1/order/cancel", bytes.NewReader(cancelBody))

			cancelReq.Header.Set("Authorization", token)
			t.httpServer.ServeHTTP(cancelRes, cancelReq)
		}()
	}

	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			//create and cancel sleep randomly between
			//call order api
			data := `{"type":"sell","amount":"0.01","exchange_type":"limit","pair_currency_id":1,"price":"50000"}`
			res := httptest.NewRecorder()
			body := []byte(data)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/order/create", bytes.NewReader(body))
			req.Header.Set("Authorization", token)
			t.httpServer.ServeHTTP(res, req)
			result := struct {
				Status  bool
				Message string
				Data    order.CreateOrderResponse
			}{}
			err := json.Unmarshal(res.Body.Bytes(), &result)
			if err != nil {
				t.Fail(err.Error())
			}

			orderID := result.Data.ID

			//there with this id we cancel the order
			orderIDString := strconv.FormatInt(orderID, 10)
			params := "{\"orderId\":" + orderIDString + "}"
			cancelRes := httptest.NewRecorder()
			cancelBody := []byte(params)
			cancelReq := httptest.NewRequest(http.MethodPost, "/api/v1/order/cancel", bytes.NewReader(cancelBody))

			cancelReq.Header.Set("Authorization", token)
			t.httpServer.ServeHTTP(cancelRes, cancelReq)
		}()
	}

	wg.Wait()
	updatedUsdtUb := &userbalance.UserBalance{}
	err := t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "10000.00000000", updatedUsdtUb.Amount)
	assert.Equal(t.T(), "0.00000000", updatedUsdtUb.FrozenAmount)

	updatedBtcUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: btcUb.ID}).First(updatedBtcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "1.00000000", updatedBtcUb.Amount)
	assert.Equal(t.T(), "0.00000000", updatedBtcUb.FrozenAmount)
}

func TestOrderCancel(t *testing.T) {
	suite.Run(t, &OrderCancelTests{
		Suite: new(suite.Suite),
	})
}
