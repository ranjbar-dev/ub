package test

import (
	"bytes"
	"context"
	"encoding/json"
	"exchange-go/internal/api"
	"exchange-go/internal/di"
	"exchange-go/internal/engine"
	"exchange-go/internal/order"
	"exchange-go/internal/platform"
	"exchange-go/internal/response"
	"exchange-go/internal/transaction"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type OrderCreateTests struct {
	*suite.Suite
	httpServer  http.Handler
	db          *gorm.DB
	redisClient *redis.Client
	userActor   *userActor
	engine      engine.Engine
}

func (t *OrderCreateTests) SetupSuite() {
	container := getContainer()
	t.httpServer = container.Get(di.HTTPServer).(api.HTTPServer).GetEngine()
	t.db = getDb()
	t.redisClient = getRedis()
	t.userActor = getUserActor()

	//set engine
	rc := container.Get(di.RedisClient).(platform.RedisClient)
	obp := engine.NewRedisOrderBookProvider(rc)
	rh := container.Get(di.EngineResultHandler).(order.EngineResultHandler)
	logger := container.Get(di.LoggerService).(platform.Logger)
	e := engine.NewEngine(rc, obp, rh, logger, platform.EnvTest)
	t.engine = e

}

func (t *OrderCreateTests) TearDownSuite() {
	t.db.Where("user_id = ?", t.userActor.ID).Delete(userbalance.UserBalance{})
}

func (t *OrderCreateTests) SetupTest() {
	up := user.UsersPermissions{
		UserID:           t.userActor.ID,
		UserPermissionID: 3, //see the userPermissionSeed id for exchange is 3
	}
	t.db.Create(&up)
}

func (t *OrderCreateTests) TearDownTest() {
	t.db.Where("user_id = ?  and user_permission_id = ?", t.userActor.ID, 3).Delete(user.UsersPermissions{})
	t.db.Where("user_id = ?", t.userActor.ID).Delete(userbalance.UserBalance{})
	t.db.Where("id > ?", 0).Delete(transaction.Transaction{})

	//empty redis queue
	allOrdersQueue := engine.QueueName
	_, err := t.redisClient.Del(context.Background(), allOrdersQueue).Result()
	if err != nil {
		t.Fail(err.Error())
	}
}

func (t *OrderCreateTests) TestOrderCreate_AA_Market_Buy_Successful() {
	//set pair price in redis
	t.redisClient.HMSet(context.Background(), "live_data:pair_currency:BTC-USDT", "price", "50000")

	//empty redis queue
	allOrdersQueue := engine.QueueName
	_, err := t.redisClient.Del(context.Background(), allOrdersQueue).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	//insert user balance
	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "10000.00",
		FrozenAmount:  "1000.00",
	}

	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.1",
		FrozenAmount:  "0",
	}
	t.db.Create(usdtUb)
	t.db.Create(btcUb)

	//call order api
	data := `{"type":"buy","amount":"100","exchange_type":"market","pair_currency_id":1,"price":""}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/order/create", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    order.CreateOrderResponse
	}{}
	//result := response.ApiResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	assert.Equal(t.T(), "", result.Data.Price)
	orderID := result.Data.ID
	assert.Greater(t.T(), orderID, int64(0))

	////check db for order
	o := &order.Order{}
	err = t.db.Where(order.Order{ID: orderID}).First(o).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "0.00200000", o.DemandedAmount.String)
	assert.Equal(t.T(), "100.00000000", o.PayedByAmount.String)
	assert.Equal(t.T(), order.TypeBuy, o.Type)
	assert.Equal(t.T(), order.ExchangeTypeMarket, o.ExchangeType)
	assert.Equal(t.T(), order.StatusOpen, o.Status)
	assert.Equal(t.T(), "", o.Price.String)

	//check db for user balance
	updatedUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "1100.00000000", updatedUsdtUb.FrozenAmount)

	//check redis orders queue only
	engineOrderIDString := fmt.Sprintf("%011d", o.ID)
	eo := engine.Order{
		Pair:              "BTC-USDT",
		ID:                engineOrderIDString,
		Side:              "bid",
		Quantity:          "0.00200000",
		Price:             "",
		Timestamp:         o.CreatedAt.Unix(),
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
	time.Sleep(50 * time.Millisecond)
	//pos, err := t.redisClient.LPos(context.Background(), allOrdersQueue, string(engineOrderData), redis.LPosArgs{}).Result()
	orderInRedis, err := t.redisClient.LPop(context.Background(), allOrdersQueue).Result()

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), string(engineOrderData), orderInRedis)
}

func (t *OrderCreateTests) TestOrderCreate_AA_Market_Buy_Successful_FastExchange() {
	//set pair price in redis
	t.redisClient.HMSet(context.Background(), "live_data:pair_currency:BTC-USDT", "price", "50000")

	//empty redis queue
	allOrdersQueue := engine.QueueName
	_, err := t.redisClient.Del(context.Background(), allOrdersQueue).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	//insert user balance
	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "10000.00",
		FrozenAmount:  "1000.00",
	}

	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.1",
		FrozenAmount:  "0",
	}
	t.db.Create(usdtUb)
	t.db.Create(btcUb)

	//call order api
	data := `{"type":"buy","amount":"100","exchange_type":"market","pair_currency_id":1,"price":"","is_fast_exchange":true}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/order/create", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    order.CreateOrderResponse
	}{}
	//result := response.ApiResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	assert.Equal(t.T(), "", result.Data.Price)
	orderID := result.Data.ID
	assert.Greater(t.T(), orderID, int64(0))

	////check db for order
	o := &order.Order{}
	err = t.db.Where(order.Order{ID: orderID}).First(o).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "0.00200000", o.DemandedAmount.String)
	assert.Equal(t.T(), "100.00000000", o.PayedByAmount.String)
	assert.Equal(t.T(), order.TypeBuy, o.Type)
	assert.Equal(t.T(), order.ExchangeTypeMarket, o.ExchangeType)
	assert.Equal(t.T(), order.StatusOpen, o.Status)
	assert.Equal(t.T(), "", o.Price.String)
	assert.Equal(t.T(), true, o.IsFastExchange)

	//check db for user balance
	updatedUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "1100.00000000", updatedUsdtUb.FrozenAmount)

	//check redis orders queue only
	engineOrderIDString := fmt.Sprintf("%011d", o.ID)
	eo := engine.Order{
		Pair:              "BTC-USDT",
		ID:                engineOrderIDString,
		Side:              "bid",
		Quantity:          "0.00200000",
		Price:             "",
		Timestamp:         o.CreatedAt.Unix(),
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

	time.Sleep(50 * time.Millisecond)
	//pos, err := t.redisClient.LPos(context.Background(), allOrdersQueue, string(engineOrderData), redis.LPosArgs{}).Result()
	orderInRedis, err := t.redisClient.LPop(context.Background(), allOrdersQueue).Result()

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), string(engineOrderData), orderInRedis)
}

func (t *OrderCreateTests) TestOrderCreate_AB_Market_Sell_Successful() {
	//set pair price in redis
	t.redisClient.HMSet(context.Background(), "live_data:pair_currency:BTC-USDT", "price", "50000")

	//empty redis queue
	allOrdersQueue := engine.QueueName
	_, err := t.redisClient.Del(context.Background(), allOrdersQueue).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	//insert user balance
	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "10000.00",
		FrozenAmount:  "1000.00",
	}

	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.3",
		FrozenAmount:  "0.1",
	}
	t.db.Create(usdtUb)
	t.db.Create(btcUb)

	//call order api
	data := `{"type":"sell","amount":"0.1","exchange_type":"market","pair_currency_id":1,"price":""}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/order/create", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    order.CreateOrderResponse
	}{}
	//result := response.ApiResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	assert.Equal(t.T(), "", result.Data.Price)
	orderID := result.Data.ID
	assert.Greater(t.T(), orderID, int64(0))

	////check db for order
	o := &order.Order{}
	err = t.db.Where(order.Order{ID: orderID}).First(o).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "5000.00000000", o.DemandedAmount.String)
	assert.Equal(t.T(), "0.10000000", o.PayedByAmount.String)
	assert.Equal(t.T(), order.TypeSell, o.Type)
	assert.Equal(t.T(), order.ExchangeTypeMarket, o.ExchangeType)
	assert.Equal(t.T(), order.StatusOpen, o.Status)
	assert.Equal(t.T(), "", o.Price.String)

	//check db for user balance
	updatedBtcUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: btcUb.ID}).First(updatedBtcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "0.20000000", updatedBtcUb.FrozenAmount)

	//check redis orders queue only
	engineOrderIDString := fmt.Sprintf("%011d", o.ID)
	eo := engine.Order{
		Pair:              "BTC-USDT",
		ID:                engineOrderIDString,
		Side:              "ask",
		Quantity:          "0.10000000",
		Price:             "",
		Timestamp:         o.CreatedAt.Unix(),
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

	time.Sleep(100 * time.Millisecond)
	//pos, err := t.redisClient.LPos(context.Background(), allOrdersQueue, string(engineOrderData), redis.LPosArgs{}).Result()
	orderInRedis, err := t.redisClient.LPop(context.Background(), allOrdersQueue).Result()

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), string(engineOrderData), orderInRedis)

}

func (t *OrderCreateTests) TestOrderCreate_AB_Market_Sell_Successful_FastExchange() {
	//set pair price in redis
	t.redisClient.HMSet(context.Background(), "live_data:pair_currency:BTC-USDT", "price", "50000")

	//empty redis queue
	allOrdersQueue := engine.QueueName
	_, err := t.redisClient.Del(context.Background(), allOrdersQueue).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	//insert user balance
	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "10000.00",
		FrozenAmount:  "1000.00",
	}

	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.3",
		FrozenAmount:  "0.1",
	}
	t.db.Create(usdtUb)
	t.db.Create(btcUb)

	//call order api
	data := `{"type":"sell","amount":"0.1","exchange_type":"market","pair_currency_id":1,"price":"","is_fast_exchange":true}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/order/create", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    order.CreateOrderResponse
	}{}
	//result := response.ApiResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	assert.Equal(t.T(), "", result.Data.Price)
	orderID := result.Data.ID
	assert.Greater(t.T(), orderID, int64(0))

	////check db for order
	o := &order.Order{}
	err = t.db.Where(order.Order{ID: orderID}).First(o).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "5000.00000000", o.DemandedAmount.String)
	assert.Equal(t.T(), "0.10000000", o.PayedByAmount.String)
	assert.Equal(t.T(), order.TypeSell, o.Type)
	assert.Equal(t.T(), order.ExchangeTypeMarket, o.ExchangeType)
	assert.Equal(t.T(), order.StatusOpen, o.Status)
	assert.Equal(t.T(), "", o.Price.String)
	assert.Equal(t.T(), true, o.IsFastExchange)

	//check db for user balance
	updatedBtcUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: btcUb.ID}).First(updatedBtcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "0.20000000", updatedBtcUb.FrozenAmount)

	//check redis orders queue only
	engineOrderIDString := fmt.Sprintf("%011d", o.ID)
	eo := engine.Order{
		Pair:              "BTC-USDT",
		ID:                engineOrderIDString,
		Side:              "ask",
		Quantity:          "0.10000000",
		Price:             "",
		Timestamp:         o.CreatedAt.Unix(),
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

	time.Sleep(100 * time.Millisecond)
	//pos, err := t.redisClient.LPos(context.Background(), allOrdersQueue, string(engineOrderData), redis.LPosArgs{}).Result()
	orderInRedis, err := t.redisClient.LPop(context.Background(), allOrdersQueue).Result()

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), string(engineOrderData), orderInRedis)

}

func (t *OrderCreateTests) TestOrderCreate_AC_Market_Buy_SendToExternalExchange_Successful() {
	//set pair price in redis
	t.redisClient.HMSet(context.Background(), "live_data:pair_currency:BTC-USDT", "price", "50000")

	// update user kyc to be able to place this order

	updatingUser := &user.User{
		ID:  t.userActor.ID,
		Kyc: user.KycLevel1Confirmation,
	}
	err := t.db.Model(updatingUser).Updates(updatingUser).Error
	if err != nil {
		t.Fail(err.Error())
	}

	//insert user balance
	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "200000.00",
		FrozenAmount:  "1000.00",
	}

	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.1",
		FrozenAmount:  "0",
	}
	t.db.Create(usdtUb)
	t.db.Create(btcUb)

	//call order api
	res := httptest.NewRecorder()
	data := `{"type":"buy","amount":"155000","exchange_type":"market","pair_currency_id":1,"price":""}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/order/create", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    order.CreateOrderResponse
	}{}

	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	assert.Equal(t.T(), "", result.Data.Price)
	orderID := result.Data.ID
	assert.Greater(t.T(), orderID, int64(0))

	////check db for order
	o := &order.Order{}
	err = t.db.Where(order.Order{ID: orderID}).First(o).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "3.10000000", o.DemandedAmount.String)
	assert.Equal(t.T(), "155000.00000000", o.PayedByAmount.String)
	assert.Equal(t.T(), order.TypeBuy, o.Type)
	assert.Equal(t.T(), order.ExchangeTypeMarket, o.ExchangeType)
	assert.Equal(t.T(), order.StatusOpen, o.Status)
	assert.Equal(t.T(), "", o.Price.String)
	//check db for user balance

	updatedUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "156000.00000000", updatedUsdtUb.FrozenAmount)

	//since the order could be updated after insertion we sleep here and then try to recheck its new state
	time.Sleep(2 * time.Second) //this high amount sleep is because we have http call in process

	updatedOrder := &order.Order{}
	err = t.db.Where(order.Order{ID: orderID}).First(updatedOrder).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "3.10000000", updatedOrder.FinalDemandedAmount.String)
	assert.Equal(t.T(), "155000.00000000", updatedOrder.FinalPayedByAmount.String)
	assert.Equal(t.T(), order.StatusFilled, updatedOrder.Status)
	assert.Equal(t.T(), "50000.00000000", updatedOrder.TradePrice.String)

	updatedUsdtUb = &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "1000.00000000", updatedUsdtUb.FrozenAmount)

}

func (t *OrderCreateTests) TestOrderCreate_AD_Market_Sell_SendToExternalExchange_Successful() {
	//set pair price in redis
	t.redisClient.HMSet(context.Background(), "live_data:pair_currency:BTC-USDT", "price", "50000")

	//empty redis queue
	allOrdersQueue := engine.QueueName
	_, err := t.redisClient.Del(context.Background(), allOrdersQueue).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	//insert user balance
	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "10000.00",
		FrozenAmount:  "1000.00",
	}

	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "3.3",
		FrozenAmount:  "0.1",
	}
	t.db.Create(usdtUb)
	t.db.Create(btcUb)

	//call order api
	data := `{"type":"sell","amount":"3.1","exchange_type":"market","pair_currency_id":1,"price":""}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/order/create", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    order.CreateOrderResponse
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	assert.Equal(t.T(), "", result.Data.Price)
	orderID := result.Data.ID
	assert.Greater(t.T(), orderID, int64(0))

	////check db for order
	o := &order.Order{}
	err = t.db.Where(order.Order{ID: orderID}).First(o).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "155000.00000000", o.DemandedAmount.String)
	assert.Equal(t.T(), "3.10000000", o.PayedByAmount.String)
	assert.Equal(t.T(), order.TypeSell, o.Type)
	assert.Equal(t.T(), order.ExchangeTypeMarket, o.ExchangeType)
	assert.Equal(t.T(), order.StatusOpen, o.Status)
	assert.Equal(t.T(), "", o.Price.String)

	//check db for user balance
	updatedBtcUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: btcUb.ID}).First(updatedBtcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "3.20000000", updatedBtcUb.FrozenAmount)

	//since the order could be updated after insertion we sleep here and then try to recheck its new state
	time.Sleep(200 * time.Millisecond)

	updatedOrder := &order.Order{}
	err = t.db.Where(order.Order{ID: orderID}).First(updatedOrder).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "155000.00000000", updatedOrder.FinalDemandedAmount.String)
	assert.Equal(t.T(), "3.10000000", updatedOrder.FinalPayedByAmount.String)
	assert.Equal(t.T(), order.StatusFilled, updatedOrder.Status)
	assert.Equal(t.T(), "50000.00000000", updatedOrder.TradePrice.String)

	updatedBtcUb = &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: btcUb.ID}).First(updatedBtcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "0.10000000", updatedBtcUb.FrozenAmount)

}

func (t *OrderCreateTests) TestOrderCreate_AE_Limit_Buy_Successful() {
	//set pair price in redis
	t.redisClient.HMSet(context.Background(), "live_data:pair_currency:BTC-USDT", "price", "50000")

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

	//insert user balance
	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "10000.00",
		FrozenAmount:  "1000.00",
	}

	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.1",
		FrozenAmount:  "0",
	}
	t.db.Create(usdtUb)
	t.db.Create(btcUb)

	//call order api
	data := `{"type":"buy","amount":"0.1","exchange_type":"limit","pair_currency_id":1,"price":"50000"}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/order/create", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    order.CreateOrderResponse
	}{}
	//result := response.ApiResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	assert.Equal(t.T(), "50000.00000000", result.Data.Price)
	orderID := result.Data.ID
	assert.Greater(t.T(), orderID, int64(0))

	////check db for order
	o := &order.Order{}
	err = t.db.Where(order.Order{ID: orderID}).First(o).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "0.10000000", o.DemandedAmount.String)
	assert.Equal(t.T(), "5000.00000000", o.PayedByAmount.String)
	assert.Equal(t.T(), order.TypeBuy, o.Type)
	assert.Equal(t.T(), order.ExchangeTypeLimit, o.ExchangeType)
	assert.Equal(t.T(), order.StatusOpen, o.Status)
	assert.Equal(t.T(), "50000.00000000", o.Price.String)

	//check db for user balance
	updatedUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "6000.00000000", updatedUsdtUb.FrozenAmount)

	//check redis orders queue only
	engineOrderIDString := fmt.Sprintf("%011d", o.ID)
	eo := engine.Order{
		Pair:              "BTC-USDT",
		ID:                engineOrderIDString,
		Side:              "bid",
		Quantity:          "0.10000000",
		Price:             "50000.00000000",
		Timestamp:         o.CreatedAt.Unix(),
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
	engineOrderBookData, err := eo.MarshalForOrderbook()
	if err != nil {
		t.Fail(err.Error())
	}

	time.Sleep(200 * time.Millisecond)

	pos, err := t.redisClient.LPos(context.Background(), allOrdersQueue, string(engineOrderData), redis.LPosArgs{}).Result()
	//orderInRedis, err := t.redisClient.LPop(context.Background(), allOrdersQueue).Result()
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), int64(pos), pos)

	_ = t.engine.SetPostOrderMatchingCall(false)
	t.engine.Run(1, false)
	time.Sleep(100 * time.Millisecond)
	t.engine.DispatchManually()
	time.Sleep(100 * time.Millisecond)
	t.engine.Stop()

	//check if it is in redis order book queue
	redisZList, err := t.redisClient.ZPopMin(context.Background(), orderBookBidQueue).Result()
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), 1, len(redisZList))
	assert.Equal(t.T(), string(engineOrderBookData), redisZList[0].Member)
}

func (t *OrderCreateTests) TestOrderCreate_AF_Limit_Sell_Successful() {
	//set pair price in redis
	t.redisClient.HMSet(context.Background(), "live_data:pair_currency:BTC-USDT", "price", "50000")

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

	//insert user balance
	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		Status:        userbalance.StatusEnabled,
		Amount:        "10000.00",
		FrozenAmount:  "1000.00",
	}

	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.3",
		FrozenAmount:  "0.1",
	}
	t.db.Create(usdtUb)
	t.db.Create(btcUb)

	//call order api
	data := `{"type":"sell","amount":"0.1","exchange_type":"limit","pair_currency_id":1,"price":"50000"}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/order/create", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    order.CreateOrderResponse
	}{}
	//result := response.ApiResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	assert.Equal(t.T(), "50000.00000000", result.Data.Price)
	orderID := result.Data.ID
	assert.Greater(t.T(), orderID, int64(0))

	////check db for order
	o := &order.Order{}
	err = t.db.Where(order.Order{ID: orderID}).First(o).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "0.10000000", o.PayedByAmount.String)
	assert.Equal(t.T(), "5000.00000000", o.DemandedAmount.String)
	assert.Equal(t.T(), order.TypeSell, o.Type)
	assert.Equal(t.T(), order.ExchangeTypeLimit, o.ExchangeType)
	assert.Equal(t.T(), order.StatusOpen, o.Status)
	assert.Equal(t.T(), "50000.00000000", o.Price.String)

	//check db for user balance
	updatedBtcUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: btcUb.ID}).First(updatedBtcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "0.20000000", updatedBtcUb.FrozenAmount)

	//check redis orders queue onl
	engineOrderIDString := fmt.Sprintf("%011d", o.ID)
	eo := engine.Order{
		Pair:              "BTC-USDT",
		ID:                engineOrderIDString,
		Side:              "ask",
		Quantity:          "0.10000000",
		Price:             "50000.00000000",
		Timestamp:         o.CreatedAt.Unix(),
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
	engineOrderBookData, err := eo.MarshalForOrderbook()
	if err != nil {
		t.Fail(err.Error())
	}

	time.Sleep(100 * time.Millisecond)
	pos, err := t.redisClient.LPos(context.Background(), allOrdersQueue, string(engineOrderData), redis.LPosArgs{}).Result()
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), int64(pos), pos)

	_ = t.engine.SetPostOrderMatchingCall(false)
	t.engine.Run(1, false)
	time.Sleep(100 * time.Millisecond)
	t.engine.DispatchManually()
	time.Sleep(100 * time.Millisecond)
	t.engine.Stop()

	//check if it is in redis order book queue
	redisZList, err := t.redisClient.ZPopMin(context.Background(), orderBookAskQueue).Result()
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), 1, len(redisZList))
	assert.Equal(t.T(), string(engineOrderBookData), redisZList[0].Member)
}

func (t *OrderCreateTests) TestOrderCreate_AG_StopLimit_Buy_Successful() {
	//set pair price in redis
	t.redisClient.HMSet(context.Background(), "live_data:pair_currency:BTC-USDT", "price", "50000")

	//empty redis queue
	stopOrderQueue := order.StopOrderQueuePrefix + "BUY:" + "BTC-USDT"
	_, err := t.redisClient.Del(context.Background(), stopOrderQueue).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	//insert user balance
	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "10000.00",
		FrozenAmount:  "1000.00",
	}

	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.1",
		FrozenAmount:  "0",
	}
	t.db.Create(usdtUb)
	t.db.Create(btcUb)

	//call order api
	data := `{"type":"buy","amount":"0.1","exchange_type":"limit","pair_currency_id":1,"price":"50000","stop_point_price":"49000"}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/order/create", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    order.CreateOrderResponse
	}{}
	//result := response.ApiResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	assert.Equal(t.T(), "50000.00000000", result.Data.Price)
	orderID := result.Data.ID
	assert.Greater(t.T(), orderID, int64(0))

	////check db for order
	o := &order.Order{}
	err = t.db.Where(order.Order{ID: orderID}).First(o).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "0.10000000", o.DemandedAmount.String)
	assert.Equal(t.T(), "5000.00000000", o.PayedByAmount.String)
	assert.Equal(t.T(), order.TypeBuy, o.Type)
	assert.Equal(t.T(), order.ExchangeTypeLimit, o.ExchangeType)
	assert.Equal(t.T(), order.StatusOpen, o.Status)
	assert.Equal(t.T(), "49000.00000000", o.StopPointPrice.String)
	assert.Equal(t.T(), "50000.00000000", o.Price.String)

	//check db for user balance
	updatedUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "6000.00000000", updatedUsdtUb.FrozenAmount)

	orderIDString := strconv.FormatInt(orderID, 10)
	time.Sleep(20 * time.Millisecond)
	//check if it is in redis order book queue
	redisZList, err := t.redisClient.ZPopMin(context.Background(), stopOrderQueue).Result()
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), 1, len(redisZList))
	assert.Equal(t.T(), orderIDString, redisZList[0].Member)
}

func (t *OrderCreateTests) TestOrderCreate_AH_StopLimit_Sell_Successful() {
	//set pair price in redis
	t.redisClient.HMSet(context.Background(), "live_data:pair_currency:BTC-USDT", "price", "50000")

	//empty redis queue
	stopOrderQueue := order.StopOrderQueuePrefix + "SELL:" + "BTC-USDT"
	_, err := t.redisClient.Del(context.Background(), stopOrderQueue).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	//insert user balance
	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "10000.00",
		FrozenAmount:  "1000.00",
	}

	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.3",
		FrozenAmount:  "0.1",
	}
	t.db.Create(usdtUb)
	t.db.Create(btcUb)

	//call order api
	data := `{"type":"sell","amount":"0.1","exchange_type":"limit","pair_currency_id":1,"price":"50000","stop_point_price":"51000"}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/order/create", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    order.CreateOrderResponse
	}{}
	//result := response.ApiResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	assert.Equal(t.T(), "50000.00000000", result.Data.Price)
	orderID := result.Data.ID
	assert.Greater(t.T(), orderID, int64(0))

	////check db for order
	o := &order.Order{}
	err = t.db.Where(order.Order{ID: orderID}).First(o).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "0.10000000", o.PayedByAmount.String)
	assert.Equal(t.T(), "5000.00000000", o.DemandedAmount.String)
	assert.Equal(t.T(), order.TypeSell, o.Type)
	assert.Equal(t.T(), order.ExchangeTypeLimit, o.ExchangeType)
	assert.Equal(t.T(), order.StatusOpen, o.Status)
	assert.Equal(t.T(), "51000.00000000", o.StopPointPrice.String)
	assert.Equal(t.T(), "50000.00000000", o.Price.String)

	//check db for user balance
	updatedBtcUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: btcUb.ID}).First(updatedBtcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "0.20000000", updatedBtcUb.FrozenAmount)

	//check redis orders queue only
	orderIDString := strconv.FormatInt(orderID, 10)

	//check if it is in redis order book queue
	redisZList, err := t.redisClient.ZPopMin(context.Background(), stopOrderQueue).Result()
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), 1, len(redisZList))
	assert.Equal(t.T(), orderIDString, redisZList[0].Member)

}

type orderCreateValidationFailedScenarios struct {
	data         string
	reason       string
	errorMessage string
}

func (t *OrderCreateTests) TestOrderCreate_AJ_ValidationFail() {
	//insert usersPermissions
	failedScenarios := []orderCreateValidationFailedScenarios{
		{
			data:         `{"type":"buy","amount":"0.1","exchange_type":"limit","pair_currency_id":1,"price":""}`,
			reason:       "no price for limit buy",
			errorMessage: "price is not valid",
		},
		{
			data:         `{"type":"sell","amount":"0.1","exchange_type":"limit","pair_currency_id":1,"price":""}`,
			reason:       "no price for limit sell",
			errorMessage: "price is not valid",
		},
		{
			data:         `{"type":"buy","amount":"0.1","exchange_type":"limit","pair_currency_id":1,"price":"-50000"}`,
			reason:       "negative price for limit buy",
			errorMessage: "price is not valid",
		},
		{
			data:         `{"type":"sell","amount":"0.1","exchange_type":"limit","pair_currency_id":1,"price":"-50000"}`,
			reason:       "negative price for limit sell",
			errorMessage: "price is not valid",
		},
		{
			data:         `{"type":"buy","amount":"-0.1","exchange_type":"limit","pair_currency_id":1,"price":"50000"}`,
			reason:       "negative amount for limit buy",
			errorMessage: "amount is not valid",
		},
		{
			data:         `{"type":"sell","amount":"-0.1","exchange_type":"limit","pair_currency_id":1,"price":"50000"}`,
			reason:       "negative amount for limit sell",
			errorMessage: "amount is not valid",
		},
		{
			data:         `{"type":"buy","amount":"0.1","exchange_type":"limit","pair_currency_id":1,"stop_point_price":"50000","price":""}`,
			reason:       "no price for stop order",
			errorMessage: "price is not valid",
		},
		{
			data:         `{"type":"buy","amount":"0.1","exchange_type":"limit","pair_currency_id":1,"stop_point_price":"-50000","price":"50000"}`,
			reason:       "negative stop point price for stop order",
			errorMessage: "stop point price is not valid",
		},
		{
			data:         `{"type":"buy","amount":"-10","exchange_type":"market","pair_currency_id":1,"price":""}`,
			reason:       "negative amount for market type",
			errorMessage: "amount is not valid",
		},
		{
			data:         `{"type":"sell","amount":"-10","exchange_type":"market","pair_currency_id":1,"price":""}`,
			reason:       "negative amount for market type",
			errorMessage: "amount is not valid",
		},
		{
			data:         `{"type":"sell","amount":"10","exchange_type":"market","pair_currency_id":1,"stop_point_price":"25000","price":""}`,
			reason:       "negative amount for market type",
			errorMessage: "stop order must be limit",
		},
	}

	for _, item := range failedScenarios {
		res := httptest.NewRecorder()
		body := []byte(item.data)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/order/create", bytes.NewReader(body))
		token := "Bearer " + t.userActor.Token

		req.Header.Set("Authorization", token)
		t.httpServer.ServeHTTP(res, req)
		result := response.APIResponse{}
		err := json.Unmarshal(res.Body.Bytes(), &result)
		if err != nil {
			t.Fail(err.Error())
		}

		assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
		assert.Equal(t.T(), item.errorMessage, result.Message)
	}
}

func (t *OrderCreateTests) TestOrderCreate_AK_NoPairCurrency() {
	data := `{"type":"sell","amount":"10","exchange_type":"market","pair_currency_id":1234,"price":""}` //1253 id  is not valid pair id
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/order/create", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token

	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "pair currency id is not valid", result.Message)

}

func (t *OrderCreateTests) TestOrderCreate_AL_LessThanMinimumAmount() {
	//set price in redis for pair
	t.redisClient.HMSet(context.Background(), "live_data:pair_currency:BTC-USDT", "price", "50000")

	scenarios := []string{
		`{"type":"buy","amount":"9","exchange_type":"market","pair_currency_id":1,"price":""}`,
		`{"type":"sell","amount":"0.0001","exchange_type":"market","pair_currency_id":1,"price":""}`,
		`{"type":"buy","amount":"0.0001","exchange_type":"limit","pair_currency_id":1,"price":"50000"}`,
		`{"type":"sell","amount":"0.0001","exchange_type":"limit","pair_currency_id":1,"price":"50000"}`,
		`{"type":"buy","amount":"0.0001","exchange_type":"limit","pair_currency_id":1,"stop_point_price":"49000","price":"50000"}`,
		`{"type":"sell","amount":"0.0001","exchange_type":"limit","pair_currency_id":1,"stop_point_price":"49000","price":"50000"}`,
	}

	for _, data := range scenarios {
		res := httptest.NewRecorder()
		body := []byte(data)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/order/create", bytes.NewReader(body))
		token := "Bearer " + t.userActor.Token

		req.Header.Set("Authorization", token)
		t.httpServer.ServeHTTP(res, req)
		result := response.APIResponse{}
		err := json.Unmarshal(res.Body.Bytes(), &result)
		if err != nil {
			t.Fail(err.Error())
		}

		assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
		assert.Equal(t.T(), "the minimum order amount must be more than 10 USDT", result.Message)
	}
}

func (t *OrderCreateTests) TestOrderCreate_AM_UserBalanceIsNotEnough() {
	//set price in redis for pair
	t.redisClient.HMSet(context.Background(), "live_data:pair_currency:BTC-USDT", "price", "50000")

	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "4000.00",
		FrozenAmount:  "0",
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

	scenarios := []string{
		`{"type":"buy","amount":"5000","exchange_type":"market","pair_currency_id":1,"price":""}`,
		`{"type":"sell","amount":"0.1","exchange_type":"market","pair_currency_id":1,"price":""}`,
		`{"type":"buy","amount":"0.1","exchange_type":"limit","pair_currency_id":1,"price":"50000"}`,
		`{"type":"sell","amount":"0.1","exchange_type":"limit","pair_currency_id":1,"price":"50000"}`,
		`{"type":"buy","amount":"0.1","exchange_type":"limit","pair_currency_id":1,"stop_point_price":"49000","price":"50000"}`,
		`{"type":"sell","amount":"0.1","exchange_type":"limit","pair_currency_id":1,"stop_point_price":"49000","price":"50000"}`,
	}

	for _, data := range scenarios {
		res := httptest.NewRecorder()
		body := []byte(data)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/order/create", bytes.NewReader(body))
		token := "Bearer " + t.userActor.Token

		req.Header.Set("Authorization", token)
		t.httpServer.ServeHTTP(res, req)
		result := response.APIResponse{}
		err := json.Unmarshal(res.Body.Bytes(), &result)
		if err != nil {
			t.Fail(err.Error())
		}

		assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
		assert.Equal(t.T(), "user balance is not enough", result.Message)
	}

}

func (t *OrderCreateTests) TestOrderCreate_AN_UserLevelDoesNotAllows() {
	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "10000.00",
		FrozenAmount:  "1000.00",
	}

	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.1",
		FrozenAmount:  "0",
	}
	t.db.Create(usdtUb)
	t.db.Create(btcUb)

	updatingUser := &user.User{
		ID:                   t.userActor.ID,
		ExchangeVolumeAmount: "99.999",
		ExchangeNumber:       100,
	}

	err := t.db.Model(updatingUser).Updates(updatingUser).Error
	if err != nil {
		t.Fail(err.Error())
	}
	//since the zero value field are not updated by gorm we do the kyc here
	err = t.db.Model(&user.User{}).Where("id = ?", t.userActor.ID).Update("kyc", 0).Error
	if err != nil {
		t.Fail(err.Error())
	}

	data := `{"type":"buy","amount":"0.01","exchange_type":"limit","pair_currency_id":1,"price":"50000"}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/order/create", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token

	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "your user level is low to place this order. please verify your identity to boost up your level", result.Message)

}

func (t *OrderCreateTests) TestOrderCreate_AO_NoPermission() {
	t.db.Where("user_id = ?  and user_permission_id = ?", t.userActor.ID, 3).Delete(user.UsersPermissions{})

	//t.db.Where("user_id = ?  and user_permission_id = ?", userActor.ID, 3).Delete(user.UsersPermissions{})
	data := `{"type":"sell","amount":"10","exchange_type":"market","pair_currency_id":1,"price":""}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/order/create", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token

	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "permission is not granted to place the order", result.Message)
}

func TestOrderCreate(t *testing.T) {
	suite.Run(t, &OrderCreateTests{
		Suite: new(suite.Suite),
	})

}
