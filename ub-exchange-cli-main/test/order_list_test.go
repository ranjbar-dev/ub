package test

import (
	"database/sql"
	"encoding/json"
	"exchange-go/internal/api"
	"exchange-go/internal/di"
	"exchange-go/internal/externalexchange"
	"exchange-go/internal/order"
	"exchange-go/internal/transaction"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type OrderListTests struct {
	*suite.Suite
	httpServer  http.Handler
	db          *gorm.DB
	redisClient *redis.Client
	userActor   *userActor
}

func (t *OrderListTests) SetupSuite() {
	container := getContainer()
	t.httpServer = container.Get(di.HTTPServer).(api.HTTPServer).GetEngine()
	t.db = getDb()
	t.redisClient = getRedis()
	t.userActor = getUserActor()
}

func (t *OrderListTests) TearDownSuite() {
	t.db.Where("id > ?", int64(0)).Delete(externalexchange.Order{})

	var ids []int64
	t.db.Where("id > ?", int64(0)).Delete(externalexchange.Order{})
	t.db.Table("orders").Where("creator_user_id = ?", t.userActor.ID).Select("orders.id").
		Order("id desc").Scan(&ids)
	for _, id := range ids {
		t.db.Where("id = ?", id).Delete(order.Order{})
	}
}

func (t *OrderListTests) SetupTest() {
	t.db.Where("id > ?", 0).Delete(order.Trade{})
	t.db.Where("id > ?", 0).Delete(transaction.Transaction{})
	t.db.Where("id > ?", int64(0)).Delete(externalexchange.Order{})
	var ids []int64
	t.db.Table("orders").Where("creator_user_id = ?", t.userActor.ID).Select("orders.id").
		Order("id desc").Scan(&ids)
	for _, id := range ids {
		t.db.Where("id = ?", id).Delete(order.Order{})
	}
	//t.db.Where("creator_user_id = ?", t.userActor.ID).Delete(order.Order{})

}

func (t *OrderListTests) TearDownTest() {
}

func (t *OrderListTests) TestOpenOrders() {
	o1 := &order.Order{
		ID:             1,
		UserID:         t.userActor.ID,
		Type:           order.TypeBuy,
		ExchangeType:   order.ExchangeTypeLimit,
		Price:          sql.NullString{String: "50000.00000000", Valid: true},
		Status:         order.StatusOpen,
		DemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		DemandedCoin:   "BTC",
		PayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		PayedByCoin:    "USDT",
		PairID:         1,
	}

	o2 := &order.Order{
		ID:             2,
		UserID:         t.userActor.ID,
		Type:           order.TypeSell,
		ExchangeType:   order.ExchangeTypeLimit,
		Price:          sql.NullString{String: "50000.00000000", Valid: true},
		Status:         order.StatusOpen,
		DemandedAmount: sql.NullString{String: "5000.00000000", Valid: true},
		DemandedCoin:   "USDT",
		PayedByAmount:  sql.NullString{String: "0.10000000", Valid: true},
		PayedByCoin:    "BTC",
		PairID:         1,
	}

	o3 := &order.Order{
		ID:             3,
		UserID:         t.userActor.ID,
		Type:           order.TypeBuy,
		ExchangeType:   order.ExchangeTypeLimit,
		Price:          sql.NullString{String: "50000.00000000", Valid: true},
		Status:         order.StatusOpen,
		DemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		DemandedCoin:   "USDT",
		PayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		PayedByCoin:    "BTC",
		PairID:         1,
		StopPointPrice: sql.NullString{String: "49000.00000000", Valid: true},
	}

	o4 := &order.Order{
		ID:             4,
		UserID:         t.userActor.ID,
		Type:           order.TypeSell,
		ExchangeType:   order.ExchangeTypeLimit,
		Price:          sql.NullString{String: "50000.00000000", Valid: true},
		Status:         order.StatusOpen,
		DemandedAmount: sql.NullString{String: "5000.00000000", Valid: true},
		DemandedCoin:   "USDT",
		PayedByAmount:  sql.NullString{String: "0.10000000", Valid: true},
		PayedByCoin:    "BTC",
		PairID:         1,
		StopPointPrice: sql.NullString{String: "51000.00000000", Valid: true},
	}

	//orders 5,6,7,8 are one order divided into 4 with open status for last one
	o5 := &order.Order{
		ID:                  5,
		UserID:              t.userActor.ID,
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.40000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "20000.00000000", Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 1, Valid: true},
		Path:                sql.NullString{String: "5,", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
	}

	o6 := &order.Order{
		ID:                  6,
		UserID:              t.userActor.ID,
		ParentID:            sql.NullInt64{Int64: 5, Valid: true},
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.30000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "15000.00000000", Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 2, Valid: true},
		Path:                sql.NullString{String: "5,6,", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
	}

	o7 := &order.Order{
		ID:                  7,
		UserID:              t.userActor.ID,
		ParentID:            sql.NullInt64{Int64: 6, Valid: true},
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.20000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "10000.00000000", Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 3, Valid: true},
		Path:                sql.NullString{String: "5,6,7,", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
	}

	o8 := &order.Order{
		ID:             8,
		UserID:         t.userActor.ID,
		ParentID:       sql.NullInt64{Int64: 7, Valid: true},
		Type:           order.TypeBuy,
		ExchangeType:   order.ExchangeTypeLimit,
		Price:          sql.NullString{String: "50000.00000000", Valid: true},
		Status:         order.StatusOpen,
		DemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		DemandedCoin:   "BTC",
		PayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		PayedByCoin:    "USDT",
		PairID:         1,
		TradePrice:     sql.NullString{String: "", Valid: false},
		Level:          sql.NullInt64{Int64: 4, Valid: true},
		Path:           sql.NullString{String: "5,6,7,", Valid: true},
	}

	//orders 9 and 10 are one order divided into 2 with open status for last one
	o9 := &order.Order{
		ID:                  9,
		UserID:              t.userActor.ID,
		Type:                order.TypeSell,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "20000.00000000", Valid: true},
		DemandedCoin:        "USDT",
		PayedByAmount:       sql.NullString{String: "0.40000000", Valid: true},
		PayedByCoin:         "BTC",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 1, Valid: true},
		Path:                sql.NullString{String: "9,", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "10000.00000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "0.20000000", Valid: true},
	}

	o10 := &order.Order{
		ID:             10,
		UserID:         t.userActor.ID,
		ParentID:       sql.NullInt64{Int64: 9, Valid: true},
		Type:           order.TypeSell,
		ExchangeType:   order.ExchangeTypeLimit,
		Price:          sql.NullString{String: "50000.00000000", Valid: true},
		Status:         order.StatusOpen,
		DemandedAmount: sql.NullString{String: "10000.00000000", Valid: true},
		DemandedCoin:   "USDT",
		PayedByAmount:  sql.NullString{String: "0.20000000", Valid: true},
		PayedByCoin:    "BTC",
		PairID:         1,
		TradePrice:     sql.NullString{String: "", Valid: false},
		Level:          sql.NullInt64{Int64: 2, Valid: true},
		Path:           sql.NullString{String: "9,10,", Valid: true},
	}

	//orders 11 and 12 are one submitted stop order divided into 2 with open status for last one
	o11 := &order.Order{
		ID:                  11,
		UserID:              t.userActor.ID,
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.20000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "10000.00000000", Valid: true},
		StopPointPrice:      sql.NullString{String: "49000.00000000", Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 1, Valid: true},
		Path:                sql.NullString{String: "11,", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		IsSubmitted:         sql.NullBool{Bool: true, Valid: true},
	}

	o12 := &order.Order{
		ID:             12,
		UserID:         t.userActor.ID,
		ParentID:       sql.NullInt64{Int64: 11, Valid: true},
		Type:           order.TypeBuy,
		ExchangeType:   order.ExchangeTypeLimit,
		Price:          sql.NullString{String: "50000.00000000", Valid: true},
		Status:         order.StatusOpen,
		DemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		DemandedCoin:   "BTC",
		PayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		PayedByCoin:    "USDT",
		PairID:         1,
		TradePrice:     sql.NullString{String: "", Valid: false},
		Level:          sql.NullInt64{Int64: 2, Valid: true},
		Path:           sql.NullString{String: "11,12,", Valid: true},
	}

	orders := []*order.Order{o1, o2, o3, o4, o5, o6, o7, o8, o9, o10, o11, o12}
	err := t.db.Create(orders).Error
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/order/open-orders", nil)
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    []order.OpenOrdersResponse
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	for _, o := range result.Data {
		switch o.ID {
		case 1:
			assert.Equal(t.T(), order.MainTypeOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.10000000", o.Amount)
			assert.Equal(t.T(), "5000.00000000", o.Total)
			assert.Equal(t.T(), "0.00 %", o.Executed)
			assert.Equal(t.T(), "", o.TriggerCondition)
		case 2:
			assert.Equal(t.T(), order.MainTypeOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeSell), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.10000000", o.Amount)
			assert.Equal(t.T(), "5000.00000000", o.Total)
			assert.Equal(t.T(), "0.00 %", o.Executed)
			assert.Equal(t.T(), "", o.TriggerCondition)

		case 3:
			assert.Equal(t.T(), order.MainTypeStopOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.10000000", o.Amount)
			assert.Equal(t.T(), "5000.00000000", o.Total)
			assert.Equal(t.T(), "0.00 %", o.Executed)
			assert.Equal(t.T(), ">= 49000.00000000", o.TriggerCondition)

		case 4:
			assert.Equal(t.T(), order.MainTypeStopOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeSell), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.10000000", o.Amount)
			assert.Equal(t.T(), "5000.00000000", o.Total)
			assert.Equal(t.T(), "0.00 %", o.Executed)
			assert.Equal(t.T(), "<= 51000.00000000", o.TriggerCondition)

		case 8:
			assert.Equal(t.T(), order.MainTypeOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.40000000", o.Amount)
			assert.Equal(t.T(), "20000.00000000", o.Total)
			assert.Equal(t.T(), "75.00 %", o.Executed)
			assert.Equal(t.T(), "", o.TriggerCondition)

		case 10:
			assert.Equal(t.T(), order.MainTypeOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeSell), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.40000000", o.Amount)
			assert.Equal(t.T(), "20000.00000000", o.Total)
			assert.Equal(t.T(), "50.00 %", o.Executed)
			assert.Equal(t.T(), "", o.TriggerCondition)

		case 12:
			assert.Equal(t.T(), order.MainTypeStopOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.20000000", o.Amount)
			assert.Equal(t.T(), "10000.00000000", o.Total)
			assert.Equal(t.T(), "50.00 %", o.Executed)
			assert.Equal(t.T(), ">= 49000.00000000", o.TriggerCondition)

		default:
			t.Fail("we should not be in default case")
		}
	}
}

func (t *OrderListTests) TestOrderHistory() {
	o1 := &order.Order{
		ID:                  1,
		UserID:              t.userActor.ID,
		Type:                order.TypeBuy,
		Path:                sql.NullString{String: "1,", Valid: true},
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.10000000", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "5000.00000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
	}

	o2 := &order.Order{
		ID:                  2,
		UserID:              t.userActor.ID,
		Type:                order.TypeSell,
		ExchangeType:        order.ExchangeTypeLimit,
		Path:                sql.NullString{String: "2,", Valid: true},
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "5000.00000000", Valid: true},
		DemandedCoin:        "USDT",
		PayedByAmount:       sql.NullString{String: "0.10000000", Valid: true},
		PayedByCoin:         "BTC",
		PairID:              1,
		FinalDemandedAmount: sql.NullString{String: "5000.00000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "0.10000000", Valid: true},
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	o3 := &order.Order{
		ID:                  3,
		UserID:              t.userActor.ID,
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Path:                sql.NullString{String: "3,", Valid: true},
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.10000000", Valid: true},
		DemandedCoin:        "USDT",
		PayedByAmount:       sql.NullString{String: "5000.00000000", Valid: true},
		PayedByCoin:         "BTC",
		PairID:              1,
		StopPointPrice:      sql.NullString{String: "49000.00000000", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	o4 := &order.Order{
		ID:                  4,
		UserID:              t.userActor.ID,
		Type:                order.TypeSell,
		ExchangeType:        order.ExchangeTypeLimit,
		Path:                sql.NullString{String: "4,", Valid: true},
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "5000.00000000", Valid: true},
		DemandedCoin:        "USDT",
		PayedByAmount:       sql.NullString{String: "0.10000000", Valid: true},
		PayedByCoin:         "BTC",
		PairID:              1,
		StopPointPrice:      sql.NullString{String: "51000.00000000", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "5000.00000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "0.10000000", Valid: true},
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	//orders 5,6,7,8 are one order divided into 4
	o5 := &order.Order{
		ID:                  5,
		UserID:              t.userActor.ID,
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.40000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "20000.00000000", Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 1, Valid: true},
		Path:                sql.NullString{String: "5,", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	o6 := &order.Order{
		ID:                  6,
		UserID:              t.userActor.ID,
		ParentID:            sql.NullInt64{Int64: 5, Valid: true},
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.30000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "15000.00000000", Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 2, Valid: true},
		Path:                sql.NullString{String: "5,6,", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	o7 := &order.Order{
		ID:                  7,
		UserID:              t.userActor.ID,
		ParentID:            sql.NullInt64{Int64: 6, Valid: true},
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.20000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "10000.00000000", Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 3, Valid: true},
		Path:                sql.NullString{String: "5,6,7,", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	o8 := &order.Order{
		ID:                  8,
		UserID:              t.userActor.ID,
		ParentID:            sql.NullInt64{Int64: 7, Valid: true},
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.10000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "5000.00000000", Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 4, Valid: true},
		Path:                sql.NullString{String: "5,6,7,8,", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
	}

	//orders 9 and 10 are one order divided into 2 with open status for last one
	o9 := &order.Order{
		ID:                  9,
		UserID:              t.userActor.ID,
		Type:                order.TypeSell,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "20000.00000000", Valid: true},
		DemandedCoin:        "USDT",
		PayedByAmount:       sql.NullString{String: "0.40000000", Valid: true},
		PayedByCoin:         "BTC",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 1, Valid: true},
		Path:                sql.NullString{String: "9,", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "10000.00000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "0.20000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	o10 := &order.Order{
		ID:             10,
		UserID:         t.userActor.ID,
		ParentID:       sql.NullInt64{Int64: 9, Valid: true},
		Type:           order.TypeSell,
		ExchangeType:   order.ExchangeTypeLimit,
		Price:          sql.NullString{String: "50000.00000000", Valid: true},
		Status:         order.StatusCanceled,
		DemandedAmount: sql.NullString{String: "10000.00000000", Valid: true},
		DemandedCoin:   "USDT",
		PayedByAmount:  sql.NullString{String: "0.20000000", Valid: true},
		PayedByCoin:    "BTC",
		PairID:         1,
		TradePrice:     sql.NullString{String: "", Valid: false},
		Level:          sql.NullInt64{Int64: 2, Valid: true},
		Path:           sql.NullString{String: "9,10,", Valid: true},
		FeePercentage:  sql.NullFloat64{Float64: 0.3, Valid: true},
		//FinalDemandedAmount: sql.NullString{String: "10000.00000000", Valid: true},
		//FinalPayedByAmount:  sql.NullString{String: "0.20000000", Valid: true},
	}

	o11 := &order.Order{
		ID:                  11,
		UserID:              t.userActor.ID,
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeMarket,
		Path:                sql.NullString{String: "11,", Valid: true},
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.10000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "5000.00000000", Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
		IsFastExchange:      true,
	}

	o12 := &order.Order{
		ID:                  12,
		UserID:              t.userActor.ID,
		Type:                order.TypeSell,
		ExchangeType:        order.ExchangeTypeMarket,
		Path:                sql.NullString{String: "12,", Valid: true},
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "5000.00000000", Valid: true},
		DemandedCoin:        "USDT",
		PayedByAmount:       sql.NullString{String: "0.10000000", Valid: true},
		PayedByCoin:         "BTC",
		PairID:              1,
		FinalDemandedAmount: sql.NullString{String: "5000.00000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "0.10000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	//orders 13,14 are one order divided into 2
	o13 := &order.Order{
		ID:                  13,
		UserID:              t.userActor.ID,
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.40000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "20000.00000000", Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
		TradePrice:          sql.NullString{String: "49000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 1, Valid: true},
		Path:                sql.NullString{String: "13,", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.20000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "9800.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	o14 := &order.Order{
		ID:                  14,
		UserID:              t.userActor.ID,
		ParentID:            sql.NullInt64{Int64: 13, Valid: true},
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.20000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "10000.00000000", Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 2, Valid: true},
		Path:                sql.NullString{String: "13,14,", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.20000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "10000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	orders := []*order.Order{o1, o2, o3, o4, o5, o6, o7, o8, o9, o10, o11, o12, o13, o14}
	err := t.db.Create(orders).Error
	if err != nil {
		t.Fail(err.Error())
	}

	//testing with no filter
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/order/history", nil)
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    []order.HistoryResponse
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	for _, o := range result.Data {
		switch o.ID {
		case 1:
			assert.Equal(t.T(), order.MainTypeOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), "50000.00000000", o.AveragePrice)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), strings.ToLower(order.StatusFilled), o.Status)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "0.1 BTC", o.Total)
			assert.Equal(t.T(), "100.00 %", o.Executed)
			assert.Equal(t.T(), "", o.TriggerCondition)
			assert.NotEmpty(t.T(), o.CreatedAt)
			assert.NotEmpty(t.T(), o.UpdatedAt)

		case 2:
			assert.Equal(t.T(), order.MainTypeOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeSell), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), "50000.00000000", o.AveragePrice)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), strings.ToLower(order.StatusFilled), o.Status)
			assert.Equal(t.T(), "0.1 BTC", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)
			assert.Equal(t.T(), "100.00 %", o.Executed)
			assert.Equal(t.T(), "", o.TriggerCondition)
			assert.NotEmpty(t.T(), o.CreatedAt)
			assert.NotEmpty(t.T(), o.UpdatedAt)

		case 3:
			assert.Equal(t.T(), order.MainTypeStopOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), "50000.00000000", o.AveragePrice)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), strings.ToLower(order.StatusFilled), o.Status)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "0.1 BTC", o.Total)
			assert.Equal(t.T(), "100.00 %", o.Executed)
			assert.Equal(t.T(), ">= 49000.00000000", o.TriggerCondition)
			assert.NotEmpty(t.T(), o.CreatedAt)
			assert.NotEmpty(t.T(), o.UpdatedAt)

		case 4:
			assert.Equal(t.T(), order.MainTypeStopOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeSell), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), "50000.00000000", o.AveragePrice)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), strings.ToLower(order.StatusFilled), o.Status)
			assert.Equal(t.T(), "0.1 BTC", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)
			assert.Equal(t.T(), "100.00 %", o.Executed)
			assert.Equal(t.T(), "<= 51000.00000000", o.TriggerCondition)
			assert.NotEmpty(t.T(), o.CreatedAt)
			assert.NotEmpty(t.T(), o.UpdatedAt)

		case 5:
			assert.Equal(t.T(), order.MainTypeOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), "50000.00000000", o.AveragePrice)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), strings.ToLower(order.StatusFilled), o.Status)
			assert.Equal(t.T(), "20000 USDT", o.Amount)
			assert.Equal(t.T(), "0.4 BTC", o.Total)
			assert.Equal(t.T(), "100.00 %", o.Executed)
			assert.Equal(t.T(), "", o.TriggerCondition)
			assert.NotEmpty(t.T(), o.CreatedAt)
			assert.NotEmpty(t.T(), o.UpdatedAt)

		case 9:
			assert.Equal(t.T(), order.MainTypeOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeSell), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), "50000.00000000", o.AveragePrice)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), strings.ToLower(order.StatusCanceled), o.Status)
			assert.Equal(t.T(), "0.4 BTC", o.Amount)
			assert.Equal(t.T(), "10000 USDT", o.Total)
			assert.Equal(t.T(), "50.00 %", o.Executed)
			assert.Equal(t.T(), "", o.TriggerCondition)
			assert.NotEmpty(t.T(), o.CreatedAt)
			assert.NotEmpty(t.T(), o.UpdatedAt)

		case 11:
			assert.Equal(t.T(), order.MainTypeOrder, o.MainType)
			assert.Equal(t.T(), "fast exchange", o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.Side)
			assert.Equal(t.T(), "", o.Price)
			assert.Equal(t.T(), "50000.00000000", o.AveragePrice)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), strings.ToLower(order.StatusFilled), o.Status)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "0.1 BTC", o.Total)
			assert.Equal(t.T(), "100.00 %", o.Executed)
			assert.Equal(t.T(), "", o.TriggerCondition)
			assert.NotEmpty(t.T(), o.CreatedAt)
			assert.NotEmpty(t.T(), o.UpdatedAt)

		case 12:
			assert.Equal(t.T(), order.MainTypeOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeMarket), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeSell), o.Side)
			assert.Equal(t.T(), "", o.Price)
			assert.Equal(t.T(), "50000.00000000", o.AveragePrice)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), strings.ToLower(order.StatusFilled), o.Status)
			assert.Equal(t.T(), "0.1 BTC", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)
			assert.Equal(t.T(), "100.00 %", o.Executed)
			assert.Equal(t.T(), "", o.TriggerCondition)
			assert.NotEmpty(t.T(), o.CreatedAt)
			assert.NotEmpty(t.T(), o.UpdatedAt)

		case 13:
			assert.Equal(t.T(), order.MainTypeOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), "49500.00000000", o.AveragePrice)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), strings.ToLower(order.StatusFilled), o.Status)
			assert.Equal(t.T(), "20000 USDT", o.Amount)
			assert.Equal(t.T(), "0.4 BTC", o.Total)
			assert.Equal(t.T(), "100.00 %", o.Executed)
			assert.Equal(t.T(), "", o.TriggerCondition)
			assert.NotEmpty(t.T(), o.CreatedAt)
			assert.NotEmpty(t.T(), o.UpdatedAt)

		default:
			t.Fail("we should not be in default case")
		}
	}

	//testing with filters
	queryParams := url.Values{}
	queryParams.Set("type", "buy")
	queryParams.Set("pair_currency_name", "BTC-USDT")

	paramsString := queryParams.Encode()
	res2 := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/order/history?"+paramsString, nil)
	token = "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res2, req)

	filteredResult := struct {
		Status  bool
		Message string
		Data    []order.HistoryResponse
	}{}
	err = json.Unmarshal(res2.Body.Bytes(), &filteredResult)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res2.Code)
	assert.Equal(t.T(), 5, len(filteredResult.Data))

	for _, o := range filteredResult.Data {
		switch o.ID {
		case 1:
			assert.Equal(t.T(), order.MainTypeOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), "50000.00000000", o.AveragePrice)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), strings.ToLower(order.StatusFilled), o.Status)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "0.1 BTC", o.Total)
			assert.Equal(t.T(), "100.00 %", o.Executed)
			assert.Equal(t.T(), "", o.TriggerCondition)

		case 3:
			assert.Equal(t.T(), order.MainTypeStopOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), "50000.00000000", o.AveragePrice)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), strings.ToLower(order.StatusFilled), o.Status)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "0.1 BTC", o.Total)
			assert.Equal(t.T(), "100.00 %", o.Executed)
			assert.Equal(t.T(), ">= 49000.00000000", o.TriggerCondition)

		case 5:
			assert.Equal(t.T(), order.MainTypeOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), "50000.00000000", o.AveragePrice)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), strings.ToLower(order.StatusFilled), o.Status)
			assert.Equal(t.T(), "20000 USDT", o.Amount)
			assert.Equal(t.T(), "0.4 BTC", o.Total)
			assert.Equal(t.T(), "100.00 %", o.Executed)
			assert.Equal(t.T(), "", o.TriggerCondition)

		case 11:
			assert.Equal(t.T(), order.MainTypeOrder, o.MainType)
			assert.Equal(t.T(), "fast exchange", o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.Side)
			assert.Equal(t.T(), "", o.Price)
			assert.Equal(t.T(), "50000.00000000", o.AveragePrice)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), strings.ToLower(order.StatusFilled), o.Status)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "0.1 BTC", o.Total)
			assert.Equal(t.T(), "100.00 %", o.Executed)
			assert.Equal(t.T(), "", o.TriggerCondition)

		case 13:
			assert.Equal(t.T(), order.MainTypeOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), "49500.00000000", o.AveragePrice)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), strings.ToLower(order.StatusFilled), o.Status)
			assert.Equal(t.T(), "20000 USDT", o.Amount)
			assert.Equal(t.T(), "0.4 BTC", o.Total)
			assert.Equal(t.T(), "100.00 %", o.Executed)
			assert.Equal(t.T(), "", o.TriggerCondition)

		default:
			t.Fail("we should not be in default case")
		}
	}

	//testing full history
	res3 := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/order/full-history", nil)
	token = "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res3, req)

	fullResultPage1 := struct {
		Status  bool
		Message string
		Data    []order.HistoryResponse
	}{}
	err = json.Unmarshal(res3.Body.Bytes(), &fullResultPage1)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res3.Code)

	assert.Equal(t.T(), 3, len(fullResultPage1.Data))

	for _, o := range fullResultPage1.Data {
		switch o.ID {
		case 11:
			assert.Equal(t.T(), order.MainTypeOrder, o.MainType)
			assert.Equal(t.T(), "fast exchange", o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.Side)
			assert.Equal(t.T(), "", o.Price)
			assert.Equal(t.T(), "50000.00000000", o.AveragePrice)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), strings.ToLower(order.StatusFilled), o.Status)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "0.1 BTC", o.Total)
			assert.Equal(t.T(), "100.00 %", o.Executed)
			assert.Equal(t.T(), "", o.TriggerCondition)

		case 12:
			assert.Equal(t.T(), order.MainTypeOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeMarket), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeSell), o.Side)
			assert.Equal(t.T(), "", o.Price)
			assert.Equal(t.T(), "50000.00000000", o.AveragePrice)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), strings.ToLower(order.StatusFilled), o.Status)
			assert.Equal(t.T(), "0.1 BTC", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)
			assert.Equal(t.T(), "100.00 %", o.Executed)
			assert.Equal(t.T(), "", o.TriggerCondition)

		case 13:
			assert.Equal(t.T(), order.MainTypeOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), "49500.00000000", o.AveragePrice)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), strings.ToLower(order.StatusFilled), o.Status)
			assert.Equal(t.T(), "20000 USDT", o.Amount)
			assert.Equal(t.T(), "0.4 BTC", o.Total)
			assert.Equal(t.T(), "100.00 %", o.Executed)
			assert.Equal(t.T(), "", o.TriggerCondition)

		default:
			t.Fail("we should not be in default case")
		}
	}

	queryParams2 := url.Values{}
	queryParams2.Set("last_id", "11")

	paramsString2 := queryParams2.Encode()
	res4 := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/order/full-history?"+paramsString2, nil)
	token = "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res4, req)

	fullResultPage2 := struct {
		Status  bool
		Message string
		Data    []order.HistoryResponse
	}{}
	err = json.Unmarshal(res4.Body.Bytes(), &fullResultPage2)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res2.Code)

	assert.Equal(t.T(), 3, len(fullResultPage2.Data))
	for _, o := range fullResultPage2.Data {
		switch o.ID {
		case 4:
			assert.Equal(t.T(), order.MainTypeStopOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeSell), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), "50000.00000000", o.AveragePrice)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), strings.ToLower(order.StatusFilled), o.Status)
			assert.Equal(t.T(), "0.1 BTC", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)
			assert.Equal(t.T(), "100.00 %", o.Executed)
			assert.Equal(t.T(), "<= 51000.00000000", o.TriggerCondition)

		case 5:
			assert.Equal(t.T(), order.MainTypeOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), "50000.00000000", o.AveragePrice)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), strings.ToLower(order.StatusFilled), o.Status)
			assert.Equal(t.T(), "20000 USDT", o.Amount)
			assert.Equal(t.T(), "0.4 BTC", o.Total)
			assert.Equal(t.T(), "100.00 %", o.Executed)
			assert.Equal(t.T(), "", o.TriggerCondition)

		case 9:
			assert.Equal(t.T(), order.MainTypeOrder, o.MainType)
			assert.Equal(t.T(), strings.ToLower(order.ExchangeTypeLimit), o.OrderType)
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeSell), o.Side)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), "50000.00000000", o.AveragePrice)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), strings.ToLower(order.StatusCanceled), o.Status)
			assert.Equal(t.T(), "0.4 BTC", o.Amount)
			assert.Equal(t.T(), "10000 USDT", o.Total)
			assert.Equal(t.T(), "50.00 %", o.Executed)
			assert.Equal(t.T(), "", o.TriggerCondition)

		default:
			t.Fail("we should not be in default case")
		}
	}
}

func (t *OrderListTests) TestTradeHistory() {
	o1 := &order.Order{
		ID:                  1,
		UserID:              t.userActor.ID,
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.10000000", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "5000.00000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
	}

	o2 := &order.Order{
		ID:                  2,
		UserID:              t.userActor.ID,
		Type:                order.TypeSell,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "5000.00000000", Valid: true},
		DemandedCoin:        "USDT",
		PayedByAmount:       sql.NullString{String: "0.10000000", Valid: true},
		PayedByCoin:         "BTC",
		PairID:              1,
		FinalDemandedAmount: sql.NullString{String: "5000.00000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "0.10000000", Valid: true},
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	o3 := &order.Order{
		ID:                  3,
		UserID:              t.userActor.ID,
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.10000000", Valid: true},
		DemandedCoin:        "USDT",
		PayedByAmount:       sql.NullString{String: "5000.00000000", Valid: true},
		PayedByCoin:         "BTC",
		PairID:              1,
		StopPointPrice:      sql.NullString{String: "49000.00000000", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	o4 := &order.Order{
		ID:                  4,
		UserID:              t.userActor.ID,
		Type:                order.TypeSell,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "5000.00000000", Valid: true},
		DemandedCoin:        "USDT",
		PayedByAmount:       sql.NullString{String: "0.10000000", Valid: true},
		PayedByCoin:         "BTC",
		PairID:              1,
		StopPointPrice:      sql.NullString{String: "51000.00000000", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "5000.00000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "0.10000000", Valid: true},
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	//orders 5,6,7,8 are one order divided into 4 with open status for last one
	o5 := &order.Order{
		ID:                  5,
		UserID:              t.userActor.ID,
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.40000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "20000.00000000", Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 1, Valid: true},
		Path:                sql.NullString{String: "5,", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	o6 := &order.Order{
		ID:                  6,
		UserID:              t.userActor.ID,
		ParentID:            sql.NullInt64{Int64: 5, Valid: true},
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.30000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "15000.00000000", Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 2, Valid: true},
		Path:                sql.NullString{String: "5,6,", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	o7 := &order.Order{
		ID:                  7,
		UserID:              t.userActor.ID,
		ParentID:            sql.NullInt64{Int64: 6, Valid: true},
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.20000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "10000.00000000", Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 3, Valid: true},
		Path:                sql.NullString{String: "5,6,7,", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	o8 := &order.Order{
		ID:                  8,
		UserID:              t.userActor.ID,
		ParentID:            sql.NullInt64{Int64: 7, Valid: true},
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.10000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "5000.00000000", Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 4, Valid: true},
		Path:                sql.NullString{String: "5,6,7,", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
	}

	//orders 9 and 10 are one order divided into 2 with open status for last one
	o9 := &order.Order{
		ID:                  9,
		UserID:              t.userActor.ID,
		Type:                order.TypeSell,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "20000.00000000", Valid: true},
		DemandedCoin:        "USDT",
		PayedByAmount:       sql.NullString{String: "0.40000000", Valid: true},
		PayedByCoin:         "BTC",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 1, Valid: true},
		Path:                sql.NullString{String: "9,", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "10000.00000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "0.20000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	o10 := &order.Order{
		ID:                  10,
		UserID:              t.userActor.ID,
		ParentID:            sql.NullInt64{Int64: 9, Valid: true},
		Type:                order.TypeSell,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "10000.00000000", Valid: true},
		DemandedCoin:        "USDT",
		PayedByAmount:       sql.NullString{String: "0.20000000", Valid: true},
		PayedByCoin:         "BTC",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 2, Valid: true},
		Path:                sql.NullString{String: "9,10,", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
		FinalDemandedAmount: sql.NullString{String: "10000.00000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "0.20000000", Valid: true},
	}

	o11 := &order.Order{
		ID:                  11,
		UserID:              t.userActor.ID,
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeMarket,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.10000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "5000.00000000", Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	o12 := &order.Order{
		ID:                  12,
		UserID:              t.userActor.ID,
		Type:                order.TypeSell,
		ExchangeType:        order.ExchangeTypeMarket,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "5000.00000000", Valid: true},
		DemandedCoin:        "USDT",
		PayedByAmount:       sql.NullString{String: "0.10000000", Valid: true},
		PayedByCoin:         "BTC",
		PairID:              1,
		FinalDemandedAmount: sql.NullString{String: "5000.00000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "0.10000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	orders := []*order.Order{o1, o2, o3, o4, o5, o6, o7, o8, o9, o10, o11, o12}
	err := t.db.Create(orders).Error
	if err != nil {
		t.Fail(err.Error())
	}

	//testing with no filter
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/trade/history", nil)
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    []order.TradeHistoryResponse
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	for _, o := range result.Data {
		switch o.ID {
		case 1:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.1 BTC", o.Executed)
			assert.Equal(t.T(), "0.03 BTC", o.Fee)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)

		case 2:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeSell), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.1 BTC", o.Executed)
			assert.Equal(t.T(), "1500 USDT", o.Fee)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)

		case 3:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.1 BTC", o.Executed)
			assert.Equal(t.T(), "0.03 BTC", o.Fee)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)

		case 4:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeSell), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.1 BTC", o.Executed)
			assert.Equal(t.T(), "1500 USDT", o.Fee)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)

		case 5:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.1 BTC", o.Executed)
			assert.Equal(t.T(), "0.03 BTC", o.Fee)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)

		case 6:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.1 BTC", o.Executed)
			assert.Equal(t.T(), "0.03 BTC", o.Fee)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)

		case 7:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.1 BTC", o.Executed)
			assert.Equal(t.T(), "0.03 BTC", o.Fee)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)

		case 8:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.1 BTC", o.Executed)
			assert.Equal(t.T(), "0.03 BTC", o.Fee)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)

		case 9:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeSell), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.2 BTC", o.Executed)
			assert.Equal(t.T(), "3000 USDT", o.Fee)
			assert.Equal(t.T(), "10000 USDT", o.Amount)
			assert.Equal(t.T(), "10000 USDT", o.Total)

		case 10:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeSell), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.2 BTC", o.Executed)
			assert.Equal(t.T(), "3000 USDT", o.Fee)
			assert.Equal(t.T(), "10000 USDT", o.Amount)
			assert.Equal(t.T(), "10000 USDT", o.Total)

		case 11:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.1 BTC", o.Executed)
			assert.Equal(t.T(), "0.03 BTC", o.Fee)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)

		case 12:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeSell), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.1 BTC", o.Executed)
			assert.Equal(t.T(), "1500 USDT", o.Fee)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)

		default:
			t.Fail("we should not be in default case")
		}
	}

	//testing with filters
	queryParams := url.Values{}
	queryParams.Set("type", "buy")
	queryParams.Set("pair_currency_name", "BTC-USDT")

	paramsString := queryParams.Encode()
	res2 := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/trade/history?"+paramsString, nil)
	token = "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res2, req)

	filteredResult := struct {
		Status  bool
		Message string
		Data    []order.TradeHistoryResponse
	}{}
	err = json.Unmarshal(res2.Body.Bytes(), &filteredResult)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res2.Code)

	for _, o := range filteredResult.Data {
		switch o.ID {
		case 1:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.1 BTC", o.Executed)
			assert.Equal(t.T(), "0.03 BTC", o.Fee)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)

		case 3:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.1 BTC", o.Executed)
			assert.Equal(t.T(), "0.03 BTC", o.Fee)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)

		case 5:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.1 BTC", o.Executed)
			assert.Equal(t.T(), "0.03 BTC", o.Fee)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)

		case 6:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.1 BTC", o.Executed)
			assert.Equal(t.T(), "0.03 BTC", o.Fee)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)

		case 7:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.1 BTC", o.Executed)
			assert.Equal(t.T(), "0.03 BTC", o.Fee)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)

		case 8:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.1 BTC", o.Executed)
			assert.Equal(t.T(), "0.03 BTC", o.Fee)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)

		case 11:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.1 BTC", o.Executed)
			assert.Equal(t.T(), "0.03 BTC", o.Fee)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)

		default:
			t.Fail("we should not be in default case")
		}
	}

	//testing the full history with pagination
	//queryParams := url.Values{}
	//queryParams.Set("type", "buy")
	//queryParams.Set("pair_currency_name", "BTC-USDT")
	//
	//paramsString := queryParams.Encode()
	res3 := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/trade/full-history", nil)
	token = "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res3, req)

	fullResultPage1 := struct {
		Status  bool
		Message string
		Data    []order.TradeHistoryResponse
	}{}
	err = json.Unmarshal(res3.Body.Bytes(), &fullResultPage1)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res3.Code)

	assert.Equal(t.T(), 3, len(fullResultPage1.Data))

	for _, o := range fullResultPage1.Data {
		switch o.ID {
		case 10:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeSell), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.2 BTC", o.Executed)
			assert.Equal(t.T(), "3000 USDT", o.Fee)
			assert.Equal(t.T(), "10000 USDT", o.Amount)
			assert.Equal(t.T(), "10000 USDT", o.Total)

		case 11:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.1 BTC", o.Executed)
			assert.Equal(t.T(), "0.03 BTC", o.Fee)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)

		case 12:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeSell), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.1 BTC", o.Executed)
			assert.Equal(t.T(), "1500 USDT", o.Fee)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)

		default:
			t.Fail("we should not be in default case")
		}
	}

	queryParams2 := url.Values{}
	queryParams2.Set("last_id", "10")

	paramsString2 := queryParams2.Encode()
	res4 := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/trade/full-history?"+paramsString2, nil)
	token = "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res4, req)

	fullResultPage2 := struct {
		Status  bool
		Message string
		Data    []order.TradeHistoryResponse
	}{}
	err = json.Unmarshal(res4.Body.Bytes(), &fullResultPage2)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res2.Code)

	assert.Equal(t.T(), 3, len(fullResultPage2.Data))
	for _, o := range fullResultPage2.Data {
		switch o.ID {
		case 7:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.1 BTC", o.Executed)
			assert.Equal(t.T(), "0.03 BTC", o.Fee)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)

		case 8:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeBuy), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.1 BTC", o.Executed)
			assert.Equal(t.T(), "0.03 BTC", o.Fee)
			assert.Equal(t.T(), "5000 USDT", o.Amount)
			assert.Equal(t.T(), "5000 USDT", o.Total)

		case 9:
			assert.Equal(t.T(), "BTC-USDT", o.Pair)
			assert.Equal(t.T(), strings.ToLower(order.TypeSell), o.OrderType)
			assert.Equal(t.T(), "50000.00000000", o.Price)
			assert.Equal(t.T(), 8, o.SubUnit)
			assert.Equal(t.T(), "0.2 BTC", o.Executed)
			assert.Equal(t.T(), "3000 USDT", o.Fee)
			assert.Equal(t.T(), "10000 USDT", o.Amount)
			assert.Equal(t.T(), "10000 USDT", o.Total)

		default:
			t.Fail("we should not be in default case")
		}
	}

}

func (t *OrderListTests) TestOrderDetail() {
	o1 := &order.Order{
		ID:                  1,
		UserID:              t.userActor.ID,
		Type:                order.TypeBuy,
		Path:                sql.NullString{String: "1,", Valid: true},
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.10000000", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "5000.00000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
	}

	o2 := &order.Order{
		ID:                  2,
		UserID:              t.userActor.ID,
		Type:                order.TypeSell,
		ExchangeType:        order.ExchangeTypeLimit,
		Path:                sql.NullString{String: "2,", Valid: true},
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "5000.00000000", Valid: true},
		DemandedCoin:        "USDT",
		PayedByAmount:       sql.NullString{String: "0.10000000", Valid: true},
		PayedByCoin:         "BTC",
		PairID:              1,
		FinalDemandedAmount: sql.NullString{String: "5000.00000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "0.10000000", Valid: true},
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	o3 := &order.Order{
		ID:                  3,
		UserID:              t.userActor.ID,
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Path:                sql.NullString{String: "3,", Valid: true},
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.10000000", Valid: true},
		DemandedCoin:        "USDT",
		PayedByAmount:       sql.NullString{String: "5000.00000000", Valid: true},
		PayedByCoin:         "BTC",
		PairID:              1,
		StopPointPrice:      sql.NullString{String: "49000.00000000", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	o4 := &order.Order{
		ID:                  4,
		UserID:              t.userActor.ID,
		Type:                order.TypeSell,
		ExchangeType:        order.ExchangeTypeLimit,
		Path:                sql.NullString{String: "4,", Valid: true},
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "5000.00000000", Valid: true},
		DemandedCoin:        "USDT",
		PayedByAmount:       sql.NullString{String: "0.10000000", Valid: true},
		PayedByCoin:         "BTC",
		PairID:              1,
		StopPointPrice:      sql.NullString{String: "51000.00000000", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "5000.00000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "0.10000000", Valid: true},
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	//orders 5, 6, 7, 8 are one order divided into 4
	o5 := &order.Order{
		ID:                  5,
		UserID:              t.userActor.ID,
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.40000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "20000.00000000", Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 1, Valid: true},
		Path:                sql.NullString{String: "5,", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	o6 := &order.Order{
		ID:                  6,
		UserID:              t.userActor.ID,
		ParentID:            sql.NullInt64{Int64: 5, Valid: true},
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.30000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "15000.00000000", Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 2, Valid: true},
		Path:                sql.NullString{String: "5,6,", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	o7 := &order.Order{
		ID:                  7,
		UserID:              t.userActor.ID,
		ParentID:            sql.NullInt64{Int64: 6, Valid: true},
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.20000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "10000.00000000", Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 3, Valid: true},
		Path:                sql.NullString{String: "5,6,7,", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	o8 := &order.Order{
		ID:                  8,
		UserID:              t.userActor.ID,
		ParentID:            sql.NullInt64{Int64: 7, Valid: true},
		Type:                order.TypeBuy,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "0.10000000", Valid: true},
		DemandedCoin:        "BTC",
		PayedByAmount:       sql.NullString{String: "5000.00000000", Valid: true},
		PayedByCoin:         "USDT",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 4, Valid: true},
		Path:                sql.NullString{String: "5,6,7,8,", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
		FinalDemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
	}

	//orders 9 and 10 are one order divided into 2 with open status for last one
	o9 := &order.Order{
		ID:                  9,
		UserID:              t.userActor.ID,
		Type:                order.TypeSell,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "20000.00000000", Valid: true},
		DemandedCoin:        "USDT",
		PayedByAmount:       sql.NullString{String: "0.40000000", Valid: true},
		PayedByCoin:         "BTC",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 1, Valid: true},
		Path:                sql.NullString{String: "9,", Valid: true},
		FinalDemandedAmount: sql.NullString{String: "10000.00000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "0.20000000", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
	}

	o10 := &order.Order{
		ID:                  10,
		UserID:              t.userActor.ID,
		ParentID:            sql.NullInt64{Int64: 9, Valid: true},
		Type:                order.TypeSell,
		ExchangeType:        order.ExchangeTypeLimit,
		Price:               sql.NullString{String: "50000.00000000", Valid: true},
		Status:              order.StatusFilled,
		DemandedAmount:      sql.NullString{String: "10000.00000000", Valid: true},
		DemandedCoin:        "USDT",
		PayedByAmount:       sql.NullString{String: "0.20000000", Valid: true},
		PayedByCoin:         "BTC",
		PairID:              1,
		TradePrice:          sql.NullString{String: "50000.00000000", Valid: true},
		Level:               sql.NullInt64{Int64: 2, Valid: true},
		Path:                sql.NullString{String: "9,10,", Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: 0.3, Valid: true},
		FinalDemandedAmount: sql.NullString{String: "10000.00000000", Valid: true},
		FinalPayedByAmount:  sql.NullString{String: "0.20000000", Valid: true},
	}

	orders := []*order.Order{o1, o2, o3, o4, o5, o6, o7, o8, o9, o10}
	err := t.db.Create(orders).Error
	if err != nil {
		t.Fail(err.Error())
	}

	queryParams := url.Values{}
	queryParams.Set("order_id", "1")
	paramsString := queryParams.Encode()

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/order/detail?"+paramsString, nil)
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    []order.DetailResponse
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)
	detail1 := result.Data[0]
	assert.Equal(t.T(), "buy", detail1.Type)
	assert.Equal(t.T(), "BTC-USDT", detail1.Pair)
	assert.Equal(t.T(), "0.03 BTC", detail1.Fee)
	assert.Equal(t.T(), "50000.00000000", detail1.Price)
	assert.Equal(t.T(), "5000 USDT", detail1.Amount)
	assert.Equal(t.T(), "0.1 BTC", detail1.Executed)

	queryParams = url.Values{}
	queryParams.Set("order_id", "2")
	paramsString = queryParams.Encode()

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/order/detail?"+paramsString, nil)
	token = "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result = struct {
		Status  bool
		Message string
		Data    []order.DetailResponse
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)
	detail2 := result.Data[0]
	assert.Equal(t.T(), "sell", detail2.Type)
	assert.Equal(t.T(), "BTC-USDT", detail2.Pair)
	assert.Equal(t.T(), "1500 USDT", detail2.Fee)
	assert.Equal(t.T(), "50000.00000000", detail2.Price)
	assert.Equal(t.T(), "5000 USDT", detail2.Amount)
	assert.Equal(t.T(), "0.1 BTC", detail2.Executed)

	queryParams = url.Values{}
	queryParams.Set("order_id", "3")
	paramsString = queryParams.Encode()

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/order/detail?"+paramsString, nil)
	token = "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result = struct {
		Status  bool
		Message string
		Data    []order.DetailResponse
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)
	detail3 := result.Data[0]
	assert.Equal(t.T(), "buy", detail3.Type)
	assert.Equal(t.T(), "BTC-USDT", detail3.Pair)
	assert.Equal(t.T(), "0.03 BTC", detail3.Fee)
	assert.Equal(t.T(), "50000.00000000", detail3.Price)
	assert.Equal(t.T(), "5000 USDT", detail3.Amount)
	assert.Equal(t.T(), "0.1 BTC", detail3.Executed)

	queryParams = url.Values{}
	queryParams.Set("order_id", "4")
	paramsString = queryParams.Encode()

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/order/detail?"+paramsString, nil)
	token = "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result = struct {
		Status  bool
		Message string
		Data    []order.DetailResponse
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)
	detail4 := result.Data[0]
	assert.Equal(t.T(), "sell", detail4.Type)
	assert.Equal(t.T(), "BTC-USDT", detail4.Pair)
	assert.Equal(t.T(), "1500 USDT", detail4.Fee)
	assert.Equal(t.T(), "50000.00000000", detail4.Price)
	assert.Equal(t.T(), "5000 USDT", detail4.Amount)
	assert.Equal(t.T(), "0.1 BTC", detail4.Executed)

	queryParams = url.Values{}
	queryParams.Set("order_id", "5")
	paramsString = queryParams.Encode()

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/order/detail?"+paramsString, nil)
	token = "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result = struct {
		Status  bool
		Message string
		Data    []order.DetailResponse
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)

	detail5 := result.Data[0]
	assert.Equal(t.T(), "buy", detail5.Type)
	assert.Equal(t.T(), "BTC-USDT", detail5.Pair)
	assert.Equal(t.T(), "0.03 BTC", detail5.Fee)
	assert.Equal(t.T(), "50000.00000000", detail5.Price)
	assert.Equal(t.T(), "5000 USDT", detail5.Amount)
	assert.Equal(t.T(), "0.1 BTC", detail5.Executed)

	detail6 := result.Data[1]
	assert.Equal(t.T(), "buy", detail6.Type)
	assert.Equal(t.T(), "BTC-USDT", detail6.Pair)
	assert.Equal(t.T(), "0.03 BTC", detail6.Fee)
	assert.Equal(t.T(), "50000.00000000", detail6.Price)
	assert.Equal(t.T(), "5000 USDT", detail6.Amount)
	assert.Equal(t.T(), "0.1 BTC", detail6.Executed)

	detail7 := result.Data[2]
	assert.Equal(t.T(), "buy", detail7.Type)
	assert.Equal(t.T(), "BTC-USDT", detail7.Pair)
	assert.Equal(t.T(), "0.03 BTC", detail7.Fee)
	assert.Equal(t.T(), "50000.00000000", detail7.Price)
	assert.Equal(t.T(), "5000 USDT", detail7.Amount)
	assert.Equal(t.T(), "0.1 BTC", detail7.Executed)

	detail8 := result.Data[3]
	assert.Equal(t.T(), "buy", detail8.Type)
	assert.Equal(t.T(), "BTC-USDT", detail8.Pair)
	assert.Equal(t.T(), "0.03 BTC", detail8.Fee)
	assert.Equal(t.T(), "50000.00000000", detail8.Price)
	assert.Equal(t.T(), "5000 USDT", detail8.Amount)
	assert.Equal(t.T(), "0.1 BTC", detail8.Executed)

	queryParams = url.Values{}
	queryParams.Set("order_id", "9")
	paramsString = queryParams.Encode()

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/order/detail?"+paramsString, nil)
	token = "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result = struct {
		Status  bool
		Message string
		Data    []order.DetailResponse
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)

	detail9 := result.Data[0]
	assert.Equal(t.T(), "sell", detail9.Type)
	assert.Equal(t.T(), "BTC-USDT", detail9.Pair)
	assert.Equal(t.T(), "3000 USDT", detail9.Fee)
	assert.Equal(t.T(), "50000.00000000", detail9.Price)
	assert.Equal(t.T(), "10000 USDT", detail9.Amount)
	assert.Equal(t.T(), "0.2 BTC", detail9.Executed)

	detail10 := result.Data[1]
	assert.Equal(t.T(), "sell", detail10.Type)
	assert.Equal(t.T(), "BTC-USDT", detail10.Pair)
	assert.Equal(t.T(), "3000 USDT", detail10.Fee)
	assert.Equal(t.T(), "50000.00000000", detail10.Price)
	assert.Equal(t.T(), "10000 USDT", detail10.Amount)
	assert.Equal(t.T(), "0.2 BTC", detail10.Executed)

}

func TestOrderList(t *testing.T) {
	suite.Run(t, &OrderListTests{
		Suite: new(suite.Suite),
	})

}
