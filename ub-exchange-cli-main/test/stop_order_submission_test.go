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
	"exchange-go/internal/transaction"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
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

type StopOrderSubmissionTests struct {
	*suite.Suite
	httpServer                 http.Handler
	db                         *gorm.DB
	redisClient                *redis.Client
	userActor                  *userActor
	engine                     engine.Engine
	stopOrderSubmissionManager order.StopOrderSubmissionManager
}

func (t *StopOrderSubmissionTests) SetupSuite() {
	container := getContainer()
	t.httpServer = container.Get(di.HTTPServer).(api.HTTPServer).GetEngine()
	t.db = getDb()
	t.redisClient = getRedis()
	t.userActor = getUserActor()

	//set engine
	rc := container.Get(di.RedisClient).(platform.RedisClient)
	logger := container.Get(di.LoggerService).(platform.Logger)
	obp := engine.NewRedisOrderBookProvider(rc, logger)
	rh := container.Get(di.EngineResultHandler).(order.EngineResultHandler)
	e := engine.NewEngine(rc, obp, rh, logger, platform.EnvTest)
	t.engine = e
	t.stopOrderSubmissionManager = container.Get(di.StopOrderSubmissionManager).(order.StopOrderSubmissionManager)

}

func (t *StopOrderSubmissionTests) TearDownSuite() {
	t.db.Where("user_id = ?", t.userActor.ID).Delete(userbalance.UserBalance{})
}

func (t *StopOrderSubmissionTests) SetupTest() {
	up := user.UsersPermissions{
		UserID:           t.userActor.ID,
		UserPermissionID: 3, //see the userPermissionSeed id for exchange is 3
	}
	t.db.Create(&up)
}

func (t *StopOrderSubmissionTests) TearDownTest() {
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

func (t *StopOrderSubmissionTests) TestStopOrderSubmission_Buy_Filled() {
	//set pair price in redis
	t.redisClient.HMSet(context.Background(), "live_data:pair_currency:BTC-USDT", "price", "50000")

	//empty redis queue
	stopOrderQueue := order.StopOrderQueuePrefix + "BUY:" + "BTC-USDT"
	_, err := t.redisClient.Del(context.Background(), stopOrderQueue).Result()
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
	zMembers, err := t.redisClient.ZRangeByScoreWithScores(context.Background(), stopOrderQueue, &redis.ZRangeBy{
		Min:    "49000",
		Max:    "51000",
		Offset: 0,
		Count:  1,
	}).Result()
	assert.Equal(t.T(), nil, err)
	assert.Equal(t.T(), float64(49000), zMembers[0].Score)
	assert.Equal(t.T(), orderIDString, zMembers[0].Member.(string))

	//here we are going to simulate price update and check if the stop order is changed properly
	t.engine.Run(1, false)
	time.Sleep(200 * time.Millisecond)
	t.stopOrderSubmissionManager.Submit(context.Background(), "BTC-USDT", "49000", "48000")
	time.Sleep(200 * time.Millisecond)
	t.engine.DispatchManually()
	time.Sleep(200 * time.Millisecond)
	t.engine.Stop()
	time.Sleep(200 * time.Millisecond)

	//checking the order status again
	updatedOrder := &order.Order{}
	err = t.db.Where(order.Order{ID: orderID}).First(updatedOrder).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), true, updatedOrder.IsSubmitted.Bool)
	assert.Equal(t.T(), "49000", updatedOrder.CurrentMarketPrice.String)
	assert.Equal(t.T(), order.StatusFilled, updatedOrder.Status)
	assert.Equal(t.T(), true, updatedOrder.IsMaker.Valid)
	assert.Equal(t.T(), false, updatedOrder.IsMaker.Bool)
	assert.Equal(t.T(), "0.10000000", updatedOrder.FinalDemandedAmount.String)
	assert.Equal(t.T(), "5000.00000000", updatedOrder.FinalPayedByAmount.String)
	assert.Equal(t.T(), 0.3, updatedOrder.FeePercentage.Float64)
	assert.Equal(t.T(), "50000.00000000", updatedOrder.TradePrice.String)

	//check redis again
	zMembers, err = t.redisClient.ZRangeByScoreWithScores(context.Background(), stopOrderQueue, &redis.ZRangeBy{
		Min:    "49000",
		Max:    "51000",
		Offset: 0,
		Count:  1,
	}).Result()
	assert.Equal(t.T(), nil, err)
	assert.Equal(t.T(), 0, len(zMembers))

	//check the user balance again
	updatedUsdtUb = &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "1000.00000000", updatedUsdtUb.FrozenAmount)
	assert.Equal(t.T(), "5000.00000000", updatedUsdtUb.Amount)

	updatedBtcUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: btcUb.ID}).First(updatedBtcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "0.17000000", updatedBtcUb.Amount)

	//check trade table
	trade := &order.Trade{}
	err = t.db.Where("buy_order_id = ?", orderID).First(trade).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "0.10000000", trade.Amount.String)
	assert.Equal(t.T(), "50000.00000000", trade.Price.String)
	assert.Equal(t.T(), false, trade.SellOrderID.Valid)
	assert.Equal(t.T(), order.TypeSell, trade.BotOrderType.String)

	//check transaction table
	var transactions []transaction.Transaction
	err = t.db.Where("user_id = ?", t.userActor.ID).Find(&transactions).Error
	if err != nil {
		t.Fail(err.Error())
	}

	for _, tx := range transactions {
		switch tx.Type {
		case transaction.TypeDemanded:
			assert.Equal(t.T(), "0.10000000", tx.Amount.String)
			assert.Equal(t.T(), orderID, tx.OrderID.Int64)
		case transaction.TypePayedBy:
			assert.Equal(t.T(), "5000.00000000", tx.Amount.String)
			assert.Equal(t.T(), orderID, tx.OrderID.Int64)
		case transaction.TypeTakerFee:
			assert.Equal(t.T(), "0.03", tx.Amount.String)
			assert.Equal(t.T(), orderID, tx.OrderID.Int64)
		default:
			t.Fail("we should not be in default case")
		}
	}

	// check redis queue again
	allOrdersQueue := engine.QueueName

	eo := engine.Order{
		Pair:              "BTC-USDT",
		ID:                orderIDString,
		Side:              "bid",
		Quantity:          "0.10000000",
		Price:             "50000.00000000",
		Timestamp:         o.CreatedAt.Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49950.00000000",
		MaxThresholdPrice: "50050.00000000",
	}

	engineOrderData, err := json.Marshal(&eo)

	_, err = t.redisClient.LPos(context.Background(), allOrdersQueue, string(engineOrderData), redis.LPosArgs{}).Result()
	assert.Equal(t.T(), redis.Nil, err)

	count, err := t.redisClient.ZCount(context.Background(), orderBookBidQueue, "-inf", "+inf").Result()
	assert.Equal(t.T(), int64(0), count)

}

func (t *StopOrderSubmissionTests) TestStopOrderSubmission_Buy_NotFilled() {
	//set pair price in redis
	t.redisClient.HMSet(context.Background(), "live_data:pair_currency:BTC-USDT", "price", "50000")

	//empty redis queue
	stopOrderQueue := order.StopOrderQueuePrefix + "BUY:" + "BTC-USDT"
	_, err := t.redisClient.Del(context.Background(), stopOrderQueue).Result()
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
	data := `{"type":"buy","amount":"0.1","exchange_type":"limit","pair_currency_id":1,"price":"47000","stop_point_price":"49000"}`
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

	assert.Equal(t.T(), "47000.00000000", result.Data.Price)
	orderID := result.Data.ID
	assert.Greater(t.T(), orderID, int64(0))

	////check db for order
	o := &order.Order{}
	err = t.db.Where(order.Order{ID: orderID}).First(o).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "0.10000000", o.DemandedAmount.String)
	assert.Equal(t.T(), "4700.00000000", o.PayedByAmount.String)
	assert.Equal(t.T(), order.TypeBuy, o.Type)
	assert.Equal(t.T(), order.ExchangeTypeLimit, o.ExchangeType)
	assert.Equal(t.T(), order.StatusOpen, o.Status)
	assert.Equal(t.T(), "49000.00000000", o.StopPointPrice.String)
	assert.Equal(t.T(), "47000.00000000", o.Price.String)

	//check db for user balance
	updatedUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "5700.00000000", updatedUsdtUb.FrozenAmount)

	orderIDString := strconv.FormatInt(orderID, 10)
	time.Sleep(20 * time.Millisecond)
	//check if it is in redis order book queue
	zMembers, err := t.redisClient.ZRangeByScoreWithScores(context.Background(), stopOrderQueue, &redis.ZRangeBy{
		Min:    "49000",
		Max:    "51000",
		Offset: 0,
		Count:  1,
	}).Result()
	assert.Equal(t.T(), nil, err)
	assert.Equal(t.T(), float64(49000), zMembers[0].Score)
	assert.Equal(t.T(), orderIDString, zMembers[0].Member.(string))

	//here we are going to simulate price update and check if the stop order is changed properly
	t.engine.Run(1, false)
	time.Sleep(200 * time.Millisecond)
	t.stopOrderSubmissionManager.Submit(context.Background(), "BTC-USDT", "49000", "48000")
	time.Sleep(200 * time.Millisecond)
	t.engine.DispatchManually()
	time.Sleep(200 * time.Millisecond)
	t.engine.Stop()
	time.Sleep(200 * time.Millisecond)

	//checking the order status again
	updatedOrder := &order.Order{}
	err = t.db.Where(order.Order{ID: orderID}).First(updatedOrder).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), true, updatedOrder.IsSubmitted.Bool)
	assert.Equal(t.T(), "49000", updatedOrder.CurrentMarketPrice.String)
	assert.Equal(t.T(), order.StatusOpen, updatedOrder.Status)

	//check the user balance again
	updatedUsdtUb = &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "5700.00000000", updatedUsdtUb.FrozenAmount)
	assert.Equal(t.T(), "10000.00000000", updatedUsdtUb.Amount)

	updatedBtcUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: btcUb.ID}).First(updatedBtcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "0.10000000", updatedBtcUb.Amount)

	//check order book
	zMembers, err = t.redisClient.ZRangeByScoreWithScores(context.Background(), orderBookBidQueue, &redis.ZRangeBy{
		Min:    "47000",
		Max:    "48000",
		Offset: 0,
		Count:  1,
	}).Result()
	assert.Equal(t.T(), nil, err)
	assert.Equal(t.T(), float64(47000), zMembers[0].Score)

}

func (t *StopOrderSubmissionTests) TestStopOrderSubmission_Sell_Filled() {
	//set pair price in redis
	t.redisClient.HMSet(context.Background(), "live_data:pair_currency:BTC-USDT", "price", "50000")

	//empty redis queue
	stopOrderQueue := order.StopOrderQueuePrefix + "SELL:" + "BTC-USDT"
	_, err := t.redisClient.Del(context.Background(), stopOrderQueue).Result()
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

	orderIDString := strconv.FormatInt(orderID, 10)
	//check if it is in redis order book queue
	zMembers, err := t.redisClient.ZRangeByScoreWithScores(context.Background(), stopOrderQueue, &redis.ZRangeBy{
		Min:    "51000",
		Max:    "52000",
		Offset: 0,
		Count:  1,
	}).Result()
	assert.Equal(t.T(), nil, err)
	assert.Equal(t.T(), float64(51000), zMembers[0].Score)
	assert.Equal(t.T(), orderIDString, zMembers[0].Member.(string))

	//here we are going to simulate price update and check if the stop order is changed properly
	t.engine.Run(1, false)
	time.Sleep(200 * time.Millisecond)
	t.stopOrderSubmissionManager.Submit(context.Background(), "BTC-USDT", "51000", "49000")
	time.Sleep(200 * time.Millisecond)
	t.engine.DispatchManually()
	time.Sleep(200 * time.Millisecond)
	t.engine.Stop()
	time.Sleep(200 * time.Millisecond)

	//checking the order status again
	updatedOrder := &order.Order{}
	err = t.db.Where(order.Order{ID: orderID}).First(updatedOrder).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), true, updatedOrder.IsSubmitted.Bool)
	assert.Equal(t.T(), "51000", updatedOrder.CurrentMarketPrice.String)
	assert.Equal(t.T(), order.StatusFilled, updatedOrder.Status)
	assert.Equal(t.T(), true, updatedOrder.IsMaker.Valid)
	assert.Equal(t.T(), false, updatedOrder.IsMaker.Bool)
	assert.Equal(t.T(), "5000.00000000", updatedOrder.FinalDemandedAmount.String)
	assert.Equal(t.T(), "0.10000000", updatedOrder.FinalPayedByAmount.String)
	assert.Equal(t.T(), 0.3, updatedOrder.FeePercentage.Float64)
	assert.Equal(t.T(), "50000.00000000", updatedOrder.TradePrice.String)

	//check redis again
	zMembers, err = t.redisClient.ZRangeByScoreWithScores(context.Background(), stopOrderQueue, &redis.ZRangeBy{
		Min:    "49000",
		Max:    "51000",
		Offset: 0,
		Count:  1,
	}).Result()
	assert.Equal(t.T(), nil, err)
	assert.Equal(t.T(), 0, len(zMembers))

	//check the user balance again
	updatedUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "1000.00", updatedUsdtUb.FrozenAmount)
	assert.Equal(t.T(), "13500.00000000", updatedUsdtUb.Amount)

	updatedBtcUb = &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: btcUb.ID}).First(updatedBtcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "0.20000000", updatedBtcUb.Amount)

	//check trade table
	trade := &order.Trade{}
	err = t.db.Where("sell_order_id = ?", orderID).First(trade).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "0.10000000", trade.Amount.String)
	assert.Equal(t.T(), "50000.00000000", trade.Price.String)
	assert.Equal(t.T(), false, trade.BuyOrderID.Valid)
	assert.Equal(t.T(), order.TypeBuy, trade.BotOrderType.String)

	//check transaction table
	var transactions []transaction.Transaction
	err = t.db.Where("user_id = ?", t.userActor.ID).Find(&transactions).Error
	if err != nil {
		t.Fail(err.Error())
	}

	for _, tx := range transactions {
		switch tx.Type {
		case transaction.TypeDemanded:
			assert.Equal(t.T(), "5000.00000000", tx.Amount.String)
			assert.Equal(t.T(), orderID, tx.OrderID.Int64)
		case transaction.TypePayedBy:
			assert.Equal(t.T(), "0.10000000", tx.Amount.String)
			assert.Equal(t.T(), orderID, tx.OrderID.Int64)
		case transaction.TypeTakerFee:
			assert.Equal(t.T(), "1500", tx.Amount.String)
			assert.Equal(t.T(), orderID, tx.OrderID.Int64)
		default:
			t.Fail("we should not be in default case")
		}
	}

	// check redis queue again
	allOrdersQueue := engine.QueueName

	eo := engine.Order{
		Pair:              "BTC-USDT",
		ID:                orderIDString,
		Side:              "bid",
		Quantity:          "0.10000000",
		Price:             "50000.00000000",
		Timestamp:         o.CreatedAt.Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49950.00000000",
		MaxThresholdPrice: "50050.00000000",
	}

	engineOrderData, err := json.Marshal(&eo)

	_, err = t.redisClient.LPos(context.Background(), allOrdersQueue, string(engineOrderData), redis.LPosArgs{}).Result()
	assert.Equal(t.T(), redis.Nil, err)

	count, err := t.redisClient.ZCount(context.Background(), orderBookAskQueue, "-inf", "+inf").Result()
	assert.Equal(t.T(), int64(0), count)
}

func (t *StopOrderSubmissionTests) TestStopOrderSubmission_Sell_NotFilled() {
	//set pair price in redis
	t.redisClient.HMSet(context.Background(), "live_data:pair_currency:BTC-USDT", "price", "50000")

	//empty redis queue
	stopOrderQueue := order.StopOrderQueuePrefix + "SELL:" + "BTC-USDT"
	_, err := t.redisClient.Del(context.Background(), stopOrderQueue).Result()
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
		Amount:        "0.3",
		FrozenAmount:  "0.1",
	}
	t.db.Create(usdtUb)
	t.db.Create(btcUb)

	//call order api
	data := `{"type":"sell","amount":"0.1","exchange_type":"limit","pair_currency_id":1,"price":"53000","stop_point_price":"51000"}`
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

	assert.Equal(t.T(), "53000.00000000", result.Data.Price)
	orderID := result.Data.ID
	assert.Greater(t.T(), orderID, int64(0))

	////check db for order
	o := &order.Order{}
	err = t.db.Where(order.Order{ID: orderID}).First(o).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "0.10000000", o.PayedByAmount.String)
	assert.Equal(t.T(), "5300.00000000", o.DemandedAmount.String)
	assert.Equal(t.T(), order.TypeSell, o.Type)
	assert.Equal(t.T(), order.ExchangeTypeLimit, o.ExchangeType)
	assert.Equal(t.T(), order.StatusOpen, o.Status)
	assert.Equal(t.T(), "51000.00000000", o.StopPointPrice.String)
	assert.Equal(t.T(), "53000.00000000", o.Price.String)

	//check db for user balance
	updatedBtcUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: btcUb.ID}).First(updatedBtcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "0.20000000", updatedBtcUb.FrozenAmount)

	orderIDString := strconv.FormatInt(orderID, 10)
	//check if it is in redis order book queue
	zMembers, err := t.redisClient.ZRangeByScoreWithScores(context.Background(), stopOrderQueue, &redis.ZRangeBy{
		Min:    "51000",
		Max:    "52000",
		Offset: 0,
		Count:  1,
	}).Result()
	assert.Equal(t.T(), nil, err)
	assert.Equal(t.T(), float64(51000), zMembers[0].Score)
	assert.Equal(t.T(), orderIDString, zMembers[0].Member.(string))

	//here we are going to simulate price update and check if the stop order is changed properly
	t.engine.Run(1, false)
	time.Sleep(200 * time.Millisecond)
	t.stopOrderSubmissionManager.Submit(context.Background(), "BTC-USDT", "51000", "49000")
	time.Sleep(200 * time.Millisecond)
	t.engine.DispatchManually()
	time.Sleep(200 * time.Millisecond)
	t.engine.Stop()
	time.Sleep(200 * time.Millisecond)

	//checking the order status again
	updatedOrder := &order.Order{}
	err = t.db.Where(order.Order{ID: orderID}).First(updatedOrder).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), true, updatedOrder.IsSubmitted.Bool)
	assert.Equal(t.T(), "51000", updatedOrder.CurrentMarketPrice.String)
	assert.Equal(t.T(), order.StatusOpen, updatedOrder.Status)

	//check redis again
	zMembers, err = t.redisClient.ZRangeByScoreWithScores(context.Background(), stopOrderQueue, &redis.ZRangeBy{
		Min:    "49000",
		Max:    "51000",
		Offset: 0,
		Count:  1,
	}).Result()
	assert.Equal(t.T(), nil, err)
	assert.Equal(t.T(), 0, len(zMembers))

	//check the user balance again
	updatedUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "1000.00", updatedUsdtUb.FrozenAmount)
	assert.Equal(t.T(), "10000.00000000", updatedUsdtUb.Amount)

	updatedBtcUb = &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: btcUb.ID}).First(updatedBtcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "0.30000000", updatedBtcUb.Amount)
	assert.Equal(t.T(), "0.20000000", updatedBtcUb.FrozenAmount)

	// check redis queue again
	allOrdersQueue := engine.QueueName

	eo := engine.Order{
		Pair:              "BTC-USDT",
		ID:                orderIDString,
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

	_, err = t.redisClient.LPos(context.Background(), allOrdersQueue, string(engineOrderData), redis.LPosArgs{}).Result()
	assert.Equal(t.T(), redis.Nil, err)

	zMembers, err = t.redisClient.ZRangeByScoreWithScores(context.Background(), orderBookAskQueue, &redis.ZRangeBy{
		Min:    "52000",
		Max:    "53000",
		Offset: 0,
		Count:  1,
	}).Result()
	assert.Equal(t.T(), nil, err)
	assert.Equal(t.T(), float64(53000), zMembers[0].Score)
}

func TestStopOrderSubmission(t *testing.T) {
	suite.Run(t, &StopOrderSubmissionTests{
		Suite: new(suite.Suite),
	})

}
