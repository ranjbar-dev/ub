package test

import (
	"bytes"
	"context"
	"encoding/json"
	"exchange-go/internal/api"
	"exchange-go/internal/di"
	"exchange-go/internal/engine"
	"exchange-go/internal/externalexchange"
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

type UnmatchedOrdersCmd struct {
	*suite.Suite
	httpServer             http.Handler
	db                     *gorm.DB
	redisClient            *redis.Client
	userActor              *userActor
	engine                 engine.Engine
	configs                platform.Configs
	unmatchedOrdersHandler order.UnmatchedOrdersHandler
}

func (t *UnmatchedOrdersCmd) SetupSuite() {
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
	e := engine.NewEngine(rc, obp, rh, logger,platform.EnvTest)
	t.engine = e
	t.configs = container.Get(di.ConfigService).(platform.Configs)
	t.unmatchedOrdersHandler = container.Get(di.UnmatchedOrderHandler).(order.UnmatchedOrdersHandler)

	//to be sure kyc level is not limitation
	updatingUser := &user.User{
		ID:                   t.userActor.ID,
		ExchangeVolumeAmount: "0",
		ExchangeNumber:       1,
	}

	err := t.db.Model(updatingUser).Updates(updatingUser).Error
	if err != nil {
		t.Fail(err.Error())
	}

	t.db.Where("id > ?", 0).Delete(userbalance.UserBalance{})
}

func (t *UnmatchedOrdersCmd) SetupTest() {
	up := user.UsersPermissions{
		UserID:           t.userActor.ID,
		UserPermissionID: 3, //see the userPermissionSeed id for exchange is 3
	}
	t.db.Create(&up)
}

func (t *UnmatchedOrdersCmd) TearDownTest() {
	t.db.Where("user_id > ?", 0).Delete(userbalance.UserBalance{})
}

func (t *UnmatchedOrdersCmd) TearDownSuite() {
	t.db.Where("user_id = ?  and user_permission_id = ?", t.userActor.ID, 3).Delete(user.UsersPermissions{})
	t.db.Where("id > ?", 0).Delete(userbalance.UserBalance{})
	t.db.Where("id > ?", 0).Delete(transaction.Transaction{})
	t.db.Where("id > ?", 0).Delete(order.Trade{})
	t.db.Where("id > ?", int64(0)).Delete(externalexchange.Order{})
	var ids []int64
	t.db.Table("orders").Where("id > ?", 0).Order("id desc").Select("orders.id").
		Order("id desc").Scan(&ids)
	for _, id := range ids {
		t.db.Where("id = ?", id).Delete(order.Order{})
	}

	//empty redis queue
	_, err := t.redisClient.Del(context.Background(), engine.QueueName).Result()
	if err != nil {
		t.Fail(err.Error())
	}
	//empty redis queue
	_, err = t.redisClient.Del(context.Background(), order.UnmatchedOrdersList).Result()
	if err != nil {
		t.Fail(err.Error())
	}
}

func (t *UnmatchedOrdersCmd) TestMatch() {
	//this is for causing error in PostOrderMatchingservice
	t.configs.Set("commitError", true)

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
	orderIDString := strconv.FormatInt(orderID, 10)
	_ = t.engine.SetPostOrderMatchingCall(true)
	t.engine.Run(1,false)
	time.Sleep(100 * time.Millisecond)
	t.engine.DispatchManually()
	time.Sleep(100 * time.Millisecond)
	t.engine.Stop()
	time.Sleep(100 * time.Millisecond)

	updatedOrder := &order.Order{}
	err = t.db.Where(order.Order{ID: orderID}).First(updatedOrder).Error
	if err != nil {
		t.Fail(err.Error())
	}

	//check the order again
	assert.Equal(t.T(), "", updatedOrder.FinalDemandedAmount.String)
	assert.Equal(t.T(), "", updatedOrder.FinalPayedByAmount.String)
	assert.Equal(t.T(), 0.0, updatedOrder.FeePercentage.Float64)
	assert.Equal(t.T(), false, updatedOrder.IsMaker.Valid)
	assert.Equal(t.T(), order.StatusOpen, updatedOrder.Status)
	assert.Equal(t.T(), "", updatedOrder.TradePrice.String)

	//check for unatched redis queue
	position, err := t.redisClient.LPos(context.Background(), order.UnmatchedOrdersList, orderIDString, redis.LPosArgs{}).Result()
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), int64(0), position)

	//running command
	t.configs.Set("commitError", false)
	t.unmatchedOrdersHandler.Match()

	//running the engine again
	_ = t.engine.SetPostOrderMatchingCall(true)
	t.engine.Run(1,false)
	time.Sleep(100 * time.Millisecond)
	t.engine.DispatchManually()
	time.Sleep(100 * time.Millisecond)
	t.engine.Stop()
	time.Sleep(100 * time.Millisecond)

	//checking all data again
	updatedOrder = &order.Order{}
	err = t.db.Where(order.Order{ID: orderID}).First(updatedOrder).Error
	if err != nil {
		t.Fail(err.Error())
	}

	//check the order again
	assert.Equal(t.T(), "0.00200000", updatedOrder.FinalDemandedAmount.String)
	assert.Equal(t.T(), "100.00000000", updatedOrder.FinalPayedByAmount.String)
	assert.Equal(t.T(), 0.3, updatedOrder.FeePercentage.Float64)
	assert.Equal(t.T(), true, updatedOrder.IsMaker.Valid)
	assert.Equal(t.T(), false, updatedOrder.IsMaker.Bool)
	assert.Equal(t.T(), order.StatusFilled, updatedOrder.Status)
	assert.Equal(t.T(), "50000.00000000", updatedOrder.TradePrice.String)

	//check the user balance again
	updatedUsdtUb = &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "1000.00000000", updatedUsdtUb.FrozenAmount)
	assert.Equal(t.T(), "9900.00000000", updatedUsdtUb.Amount)

	updatedBtcUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: btcUb.ID}).First(updatedBtcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "0.10140000", updatedBtcUb.Amount)

	//check trade table
	trade := &order.Trade{}
	err = t.db.Where("buy_order_id = ?", orderID).First(trade).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "0.00200000", trade.Amount.String)
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
			assert.Equal(t.T(), "0.00200000", tx.Amount.String)
			assert.Equal(t.T(), orderID, tx.OrderID.Int64)
		case transaction.TypePayedBy:
			assert.Equal(t.T(), "100.00000000", tx.Amount.String)
			assert.Equal(t.T(), orderID, tx.OrderID.Int64)
		case transaction.TypeTakerFee:
			assert.Equal(t.T(), "0.0006", tx.Amount.String)
			assert.Equal(t.T(), orderID, tx.OrderID.Int64)
		default:
			t.Fail("we should not be in default case")
		}
	}

	_, err = t.redisClient.LPos(context.Background(), order.UnmatchedOrdersList, orderIDString, redis.LPosArgs{}).Result()
	assert.Equal(t.T(), redis.Nil, err)
	t.configs.Set("commitError", false)
}

func TestUnmatchedOrdersCmd(t *testing.T) {
	suite.Run(t, &UnmatchedOrdersCmd{
		Suite: new(suite.Suite),
	})
}
