package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"exchange-go/internal/api"
	"exchange-go/internal/di"
	"exchange-go/internal/order"
	"exchange-go/internal/payment"
	"exchange-go/internal/response"
	"exchange-go/internal/transaction"
	"exchange-go/internal/userbalance"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AutoExchangeTests struct {
	*suite.Suite
	httpServer      http.Handler
	adminHTTPServer http.Handler
	db              *gorm.DB
	redisClient     *redis.Client
	userActor       *userActor
	adminUserActor  *userActor
}

func (t *AutoExchangeTests) SetupSuite() {
	container := getContainer()
	t.httpServer = container.Get(di.HTTPServer).(api.HTTPServer).GetEngine()
	t.adminHTTPServer = container.Get(di.HTTPServer).(api.HTTPServer).GetAdminEngine()
	t.db = getDb()
	t.redisClient = getRedis()
	t.userActor = getUserActor()
	t.adminUserActor = getAdminUserActor()
}

func (t *AutoExchangeTests) SetupTest() {
	err := t.db.Where("id > ?", 0).Delete(userbalance.UserBalance{}).Error
	if err != nil {
		t.Fail(err.Error())
	}
	err = t.db.Where("id > ?", 0).Delete(transaction.Transaction{}).Error
	if err != nil {
		t.Fail(err.Error())
	}
	err = t.db.Where("id > ?", 0).Delete(payment.ExtraInfo{}).Error
	if err != nil {
		t.Fail(err.Error())
	}
	err = t.db.Where("id > ?", 0).Delete(payment.Payment{}).Error
	if err != nil {
		t.Fail(err.Error())
	}
	err = t.db.Where("id > ?", 0).Delete(order.Order{}).Error
	if err != nil {
		t.Fail(err.Error())
	}
}

func (t *AutoExchangeTests) TearDownTest() {
}

func (t *AutoExchangeTests) TearDownSuite() {

}

func (t *AutoExchangeTests) TestAutoExchange_Successful_BuyOrder() {
	t.redisClient.HMSet(context.Background(), "live_data:pair_currency:BTC-USDT", "price", "50000")
	user := getNewUserActor()
	usdtUb := &userbalance.UserBalance{
		UserID:           user.ID,
		CoinID:           1, //for usdt from currency seed
		FrozenBalance:    "",
		BalanceCoin:      "USDT",
		Status:           userbalance.StatusEnabled,
		Amount:           "1000.00",
		FrozenAmount:     "500.00",
		Address:          sql.NullString{String: "user1USDTAddress", Valid: true},
		AutoExchangeCoin: sql.NullString{String: "BTC", Valid: true},
	}
	err := t.db.Create(usdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	p := &payment.Payment{
		UserID:            user.ID,
		CoinID:            1,
		Type:              "DEPOSIT",
		Status:            "CREATED",
		Code:              "USDT",
		FromAddress:       sql.NullString{String: "fromAddress", Valid: true},
		ToAddress:         sql.NullString{String: "user1USDTAddress", Valid: true},
		TxID:              sql.NullString{String: "txIdFromCallback", Valid: true},
		Amount:            sql.NullString{String: "100.00000000", Valid: true},
		BlockchainNetwork: sql.NullString{String: "ETH", Valid: true},
	}
	err = t.db.Create(p).Error
	if err != nil {
		t.Fail(err.Error())
	}
	extraInfo := &payment.ExtraInfo{
		PaymentID: p.ID,
	}
	err = t.db.Create(extraInfo).Error
	if err != nil {
		t.Fail(err.Error())
	}

	data := `{
		"code" : "USDT",
		"amount" : "100.00000000",
		"type" : "DEPOSIT",
		"status" : "COMPLETED",
		"from_address" : "fromAddress",
		"to_address" : "user1USDTAddress",
		"tx_id" : "txIdFromCallback",
		"network" : "ETH"
	}`
	body := []byte(data)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payment/callback", bytes.NewReader(body))
	token := "Bearer " + t.adminUserActor.Token
	req.Header.Set("Authorization", token)
	t.adminHTTPServer.ServeHTTP(res, req)

	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.T().Error(err)
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)

	updatedPayment := &payment.Payment{}
	err = t.db.Where(payment.Payment{ID: p.ID}).First(updatedPayment).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), payment.StatusCompleted, updatedPayment.Status)
	assert.Equal(t.T(), "100.00000000", updatedPayment.Amount.String)
	assert.Equal(t.T(), "fromAddress", updatedPayment.FromAddress.String)
	assert.Equal(t.T(), "user1USDTAddress", updatedPayment.ToAddress.String)
	assert.Equal(t.T(), "USDT", updatedPayment.Code)
	assert.Equal(t.T(), "txIdFromCallback", updatedPayment.TxID.String)
	assert.Equal(t.T(), "DEPOSIT", updatedPayment.Type)
	assert.Equal(t.T(), "ETH", updatedPayment.BlockchainNetwork.String)

	updatedUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "1100.00000000", updatedUsdtUb.Amount)
	assert.Equal(t.T(), "500.00", updatedUsdtUb.FrozenAmount)

	tx := &transaction.Transaction{}
	err = t.db.Where(transaction.Transaction{PaymentID: sql.NullInt64{Int64: p.ID, Valid: true}}).First(tx).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "100.00000000", tx.Amount.String)
	assert.Equal(t.T(), int64(1), tx.CoinID)
	assert.Equal(t.T(), "USDT", tx.CoinName)
	assert.Equal(t.T(), "DEPOSIT", tx.Type)
	assert.Equal(t.T(), user.ID, tx.UserID)

	//sleeping so the auto exchange could be run
	time.Sleep(50 * time.Millisecond)
	//checking if the order is created because of autoExchange
	o := &order.Order{}
	err = t.db.Where("creator_user_id = ?", user.ID).First(o).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "0.00200000", o.DemandedAmount.String)
	assert.Equal(t.T(), "100.00000000", o.PayedByAmount.String)
	assert.Equal(t.T(), order.TypeBuy, o.Type)
	assert.Equal(t.T(), order.ExchangeTypeMarket, o.ExchangeType)
	assert.Equal(t.T(), order.StatusOpen, o.Status)
	assert.Equal(t.T(), "", o.Price.String)
	assert.Equal(t.T(), false, o.IsFastExchange)
	ei := &order.ExtraInfo{}
	err = t.db.Where("id = ?", o.ExtraInfoID).First(ei).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), true, ei.AutoExchange.Valid)
	assert.Equal(t.T(), true, ei.AutoExchange.Bool)

	//checking payment extra info
	pei := &payment.ExtraInfo{}
	err = t.db.Where("crypto_payment_id = ?", p.ID).First(pei).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), o.ID, pei.AutoExchangeOrderID.Int64)
	assert.Equal(t.T(), false, pei.AutoExchangeFailureReason.Valid)
	assert.Equal(t.T(), false, pei.AutoExchangeFailureType.Valid)
}

func (t *AutoExchangeTests) TestAutoExchange_Successful_SellOrder() {
	t.redisClient.HMSet(context.Background(), "live_data:pair_currency:BTC-USDT", "price", "50000")
	user := getNewUserActor()
	usdtUb := &userbalance.UserBalance{
		UserID:           user.ID,
		CoinID:           2, //for BTC from currency seed
		FrozenBalance:    "",
		BalanceCoin:      "BTC",
		Status:           userbalance.StatusEnabled,
		Amount:           "1.00",
		FrozenAmount:     "0.00",
		Address:          sql.NullString{String: "user1BTCAddress", Valid: true},
		AutoExchangeCoin: sql.NullString{String: "USDT", Valid: true},
	}
	err := t.db.Create(usdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	data := `{
		"code" : "BTC",
		"amount" : "0.002",
		"type" : "DEPOSIT",
		"status" : "COMPLETED",
		"from_address" : "fromAddress",
		"to_address" : "user1BTCAddress",
		"tx_id" : "txIdFromCallback",
		"network" : "BTC"
	}`
	body := []byte(data)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payment/callback", bytes.NewReader(body))
	token := "Bearer " + t.adminUserActor.Token
	req.Header.Set("Authorization", token)
	t.adminHTTPServer.ServeHTTP(res, req)

	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.T().Error(err)
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)

	createdPayment := &payment.Payment{}
	err = t.db.Where(payment.Payment{UserID: user.ID}).First(createdPayment).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), payment.StatusCompleted, createdPayment.Status)
	assert.Equal(t.T(), "0.00200000", createdPayment.Amount.String)
	assert.Equal(t.T(), "fromAddress", createdPayment.FromAddress.String)
	assert.Equal(t.T(), "user1BTCAddress", createdPayment.ToAddress.String)
	assert.Equal(t.T(), "BTC", createdPayment.Code)
	assert.Equal(t.T(), "txIdFromCallback", createdPayment.TxID.String)
	assert.Equal(t.T(), "DEPOSIT", createdPayment.Type)
	assert.Equal(t.T(), "BTC", createdPayment.BlockchainNetwork.String)

	updatedBTCUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedBTCUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "1.00200000", updatedBTCUb.Amount)
	assert.Equal(t.T(), "0.00", updatedBTCUb.FrozenAmount)

	tx := &transaction.Transaction{}
	err = t.db.Where(transaction.Transaction{PaymentID: sql.NullInt64{Int64: createdPayment.ID, Valid: true}}).First(tx).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "0.00200000", tx.Amount.String)
	assert.Equal(t.T(), int64(2), tx.CoinID)
	assert.Equal(t.T(), "BTC", tx.CoinName)
	assert.Equal(t.T(), "DEPOSIT", tx.Type)
	assert.Equal(t.T(), user.ID, tx.UserID)

	//sleeping so the auto exchange could be run
	time.Sleep(50 * time.Millisecond)
	//checking if the order is created because of autoExchange
	o := &order.Order{}
	err = t.db.Where("creator_user_id = ?", user.ID).First(o).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "100.00000000", o.DemandedAmount.String)
	assert.Equal(t.T(), "0.00200000", o.PayedByAmount.String)
	assert.Equal(t.T(), order.TypeSell, o.Type)
	assert.Equal(t.T(), order.ExchangeTypeMarket, o.ExchangeType)
	assert.Equal(t.T(), order.StatusOpen, o.Status)
	assert.Equal(t.T(), "", o.Price.String)
	assert.Equal(t.T(), false, o.IsFastExchange)
	ei := &order.ExtraInfo{}
	err = t.db.Where("id = ?", o.ExtraInfoID).First(ei).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), true, ei.AutoExchange.Valid)
	assert.Equal(t.T(), true, ei.AutoExchange.Bool)

	//checking payment extra info
	pei := &payment.ExtraInfo{}
	err = t.db.Where("crypto_payment_id = ?", createdPayment.ID).First(pei).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), o.ID, pei.AutoExchangeOrderID.Int64)
	assert.Equal(t.T(), false, pei.AutoExchangeFailureReason.Valid)
	assert.Equal(t.T(), false, pei.AutoExchangeFailureType.Valid)
}

func (t *AutoExchangeTests) TestAutoExchange_Unsuccessful_LogicalError() {
	//our logical error would be less than the minimum order amount
	//minimum order is 10$ but the user deposited 5$
	t.redisClient.HMSet(context.Background(), "live_data:pair_currency:BTC-USDT", "price", "50000")
	user := getNewUserActor()
	usdtUb := &userbalance.UserBalance{
		UserID:           user.ID,
		CoinID:           1, //for usdt from currency seed
		FrozenBalance:    "",
		BalanceCoin:      "USDT",
		Status:           userbalance.StatusEnabled,
		Amount:           "1000.00",
		FrozenAmount:     "500.00",
		Address:          sql.NullString{String: "user1USDTAddress", Valid: true},
		AutoExchangeCoin: sql.NullString{String: "BTC", Valid: true},
	}
	err := t.db.Create(usdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	p := &payment.Payment{
		UserID:            user.ID,
		CoinID:            1,
		Type:              "DEPOSIT",
		Status:            "CREATED",
		Code:              "USDT",
		FromAddress:       sql.NullString{String: "fromAddress", Valid: true},
		ToAddress:         sql.NullString{String: "user1USDTAddress", Valid: true},
		TxID:              sql.NullString{String: "txIdFromCallback", Valid: true},
		Amount:            sql.NullString{String: "5.00000000", Valid: true},
		BlockchainNetwork: sql.NullString{String: "ETH", Valid: true},
	}
	err = t.db.Create(p).Error
	if err != nil {
		t.Fail(err.Error())
	}
	extraInfo := &payment.ExtraInfo{
		PaymentID: p.ID,
	}
	err = t.db.Create(extraInfo).Error
	if err != nil {
		t.Fail(err.Error())
	}

	data := `{
		"code" : "USDT",
		"amount" : "5.00000000",
		"type" : "DEPOSIT",
		"status" : "COMPLETED",
		"from_address" : "fromAddress",
		"to_address" : "user1USDTAddress",
		"tx_id" : "txIdFromCallback",
		"network" : "ETH"
	}`
	body := []byte(data)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payment/callback", bytes.NewReader(body))
	token := "Bearer " + t.adminUserActor.Token
	req.Header.Set("Authorization", token)
	t.adminHTTPServer.ServeHTTP(res, req)

	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.T().Error(err)
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)

	updatedPayment := &payment.Payment{}
	err = t.db.Where(payment.Payment{ID: p.ID}).First(updatedPayment).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), payment.StatusCompleted, updatedPayment.Status)
	assert.Equal(t.T(), "5.00000000", updatedPayment.Amount.String)
	assert.Equal(t.T(), "fromAddress", updatedPayment.FromAddress.String)
	assert.Equal(t.T(), "user1USDTAddress", updatedPayment.ToAddress.String)
	assert.Equal(t.T(), "USDT", updatedPayment.Code)
	assert.Equal(t.T(), "txIdFromCallback", updatedPayment.TxID.String)
	assert.Equal(t.T(), "DEPOSIT", updatedPayment.Type)
	assert.Equal(t.T(), "ETH", updatedPayment.BlockchainNetwork.String)

	updatedUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "1005.00000000", updatedUsdtUb.Amount)
	assert.Equal(t.T(), "500.00", updatedUsdtUb.FrozenAmount)

	tx := &transaction.Transaction{}
	err = t.db.Where(transaction.Transaction{PaymentID: sql.NullInt64{Int64: p.ID, Valid: true}}).First(tx).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "5.00000000", tx.Amount.String)
	assert.Equal(t.T(), int64(1), tx.CoinID)
	assert.Equal(t.T(), "USDT", tx.CoinName)
	assert.Equal(t.T(), "DEPOSIT", tx.Type)
	assert.Equal(t.T(), user.ID, tx.UserID)

	//sleeping so the auto exchange could be run
	time.Sleep(50 * time.Millisecond)
	//checking if the order is created because of autoExchange
	o := &order.Order{}
	err = t.db.Where("creator_user_id = ?", user.ID).First(o).Error
	assert.Equal(t.T(), gorm.ErrRecordNotFound, err)

	//checking payment extra info
	pei := &payment.ExtraInfo{}
	err = t.db.Where("crypto_payment_id = ?", p.ID).First(pei).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), false, pei.AutoExchangeOrderID.Valid)
	assert.Equal(t.T(), int64(0), pei.AutoExchangeOrderID.Int64)
	assert.Equal(t.T(), true, pei.AutoExchangeFailureReason.Valid)
	assert.Equal(t.T(), "the minimum order amount must be more than 10 USDT", pei.AutoExchangeFailureReason.String)
	assert.Equal(t.T(), true, pei.AutoExchangeFailureType.Valid)
	assert.Equal(t.T(), payment.FailureTypeLogical, pei.AutoExchangeFailureType.String)

}

func TestAutoExchange(t *testing.T) {
	suite.Run(t, &AutoExchangeTests{
		Suite: new(suite.Suite),
	})

}
