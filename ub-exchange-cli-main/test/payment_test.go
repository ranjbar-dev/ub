package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"exchange-go/internal/api"
	"exchange-go/internal/di"
	"exchange-go/internal/payment"
	"exchange-go/internal/response"
	"exchange-go/internal/transaction"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type PaymentTests struct {
	*suite.Suite
	httpServer      http.Handler
	adminHTTPServer http.Handler
	db              *gorm.DB
	redisClient     *redis.Client
	userActor       *userActor
	adminUserActor  *userActor
}

func (t *PaymentTests) SetupSuite() {
	container := getContainer()
	t.httpServer = container.Get(di.HTTPServer).(api.HTTPServer).GetEngine()
	t.adminHTTPServer = container.Get(di.HTTPServer).(api.HTTPServer).GetAdminEngine()
	t.db = getDb()
	t.redisClient = getRedis()
	t.userActor = getUserActor()
	t.adminUserActor = getAdminUserActor()
}

func (t *PaymentTests) SetupTest() {
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
	err = t.db.Where("user_id = ?", t.userActor.ID).Delete(user.Config{}).Error
	if err != nil {
		t.Fail(err.Error())
	}
	t.db.Where("user_id = ?", t.userActor.ID).Delete(userbalance.UserBalance{})
	up := user.UsersPermissions{
		UserID:           t.userActor.ID,
		UserPermissionID: 2, //see the userPermissionSeed id for withdraw is 2
	}
	t.db.Create(&up)

	updatingUser := &user.User{ID: t.userActor.ID}
	t.db.Model(updatingUser).Updates(map[string]interface{}{"status": "VERIFIED", "google2fa_disabled_at": nil, "is_two_fa_enabled": false})

	//delete former withdraw-confirmation for user
	key := "withdraw-confirmation:" + strconv.Itoa(t.userActor.ID)
	t.redisClient.Del(context.Background(), key)
}

func (t *PaymentTests) TearDownTest() {
	t.db.Where("user_id = ?", t.userActor.ID).Delete(user.UsersPermissions{})
}

func (t *PaymentTests) TearDownSuite() {
	t.db.Where("id > ?", 0).Delete(payment.Payment{})
}

func (t *PaymentTests) TestGetPayments() {
	p1 := &payment.Payment{
		ID:                1,
		UserID:            t.userActor.ID,
		CoinID:            1,
		Type:              payment.TypeWithdraw,
		Status:            payment.StatusCreated,
		Code:              "USDT",
		ToAddress:         sql.NullString{String: "usdtToAddress1", Valid: true},
		BlockchainNetwork: sql.NullString{String: "ETH", Valid: true},
		Amount:            sql.NullString{String: "100.00000000", Valid: true},
		FeeAmount:         sql.NullString{String: "10.00000000", Valid: true},
		UpdatedAt:         time.Now().Add(-2 * time.Second),
	}

	p2 := &payment.Payment{
		ID:                2,
		UserID:            t.userActor.ID,
		CoinID:            1,
		Type:              payment.TypeDeposit,
		Status:            payment.StatusInProgress,
		Code:              "USDT",
		ToAddress:         sql.NullString{String: "usdtToAddress2", Valid: true},
		BlockchainNetwork: sql.NullString{String: "ETH", Valid: true},
		Amount:            sql.NullString{String: "100.00000000", Valid: true},
		FeeAmount:         sql.NullString{String: "", Valid: false},
		UpdatedAt:         time.Now().Add(-3 * time.Second),
	}

	p3 := &payment.Payment{
		ID:        3,
		UserID:    t.userActor.ID,
		CoinID:    2,
		Type:      payment.TypeWithdraw,
		Status:    payment.StatusCompleted,
		Code:      "BTC",
		ToAddress: sql.NullString{String: "btcToAddress1", Valid: true},
		Amount:    sql.NullString{String: "0.01000000", Valid: true},
		FeeAmount: sql.NullString{String: "0.00100000", Valid: true},
		TxID:      sql.NullString{String: "btcTxId1", Valid: true},
		UpdatedAt: time.Now().Add(-4 * time.Second),
	}

	p4 := &payment.Payment{
		ID:        4,
		UserID:    t.userActor.ID,
		CoinID:    3,
		Type:      payment.TypeWithdraw,
		Status:    payment.StatusFailed,
		Code:      "ETH",
		ToAddress: sql.NullString{String: "ethToAddress1", Valid: true},
		Amount:    sql.NullString{String: "0.10000000", Valid: true},
		FeeAmount: sql.NullString{String: "0.00100000", Valid: true},
		UpdatedAt: time.Now().Add(-5 * time.Second),
	}

	payments := []*payment.Payment{p1, p2, p3, p4}
	err := t.db.Create(payments).Error
	if err != nil {
		t.Fail(err.Error())
	}

	queryParams := url.Values{}
	queryParams.Set("page_size", "2")
	paramsString := queryParams.Encode()

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/crypto-payment?"+paramsString, nil)
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    map[string][]payment.GetPaymentsResponse
	}{}
	//result := response.ApiResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)
	finalData, ok := result.Data["payments"]
	if !ok {
		t.Error(fmt.Errorf("payments key does not exists in response"))
	}

	assert.Equal(t.T(), 2, len(finalData))

	for _, p := range finalData {
		switch p.ID {
		case 1:
			assert.Equal(t.T(), "USDT", p.Coin)
			assert.Equal(t.T(), "pending", p.Status)
			assert.Equal(t.T(), "usdtToAddress1", p.Address)
			assert.Equal(t.T(), "100.00000000", p.Amount)
			assert.Equal(t.T(), "", p.TxID)
			assert.Equal(t.T(), "withdraw", p.Type)
		case 2:
			assert.Equal(t.T(), "USDT", p.Coin)
			assert.Equal(t.T(), "in progress", p.Status)
			assert.Equal(t.T(), "usdtToAddress2", p.Address)
			assert.Equal(t.T(), "100.00000000", p.Amount)
			assert.Equal(t.T(), "", p.TxID)
			assert.Equal(t.T(), "deposit", p.Type)
		default:
			t.Fail("we should not be in default case")
		}

	}

	//test with pagination
	queryParams = url.Values{}
	queryParams.Set("page", "1")
	queryParams.Set("page_size", "2")
	paramsString = queryParams.Encode()
	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/crypto-payment?"+paramsString, nil)
	token = "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result = struct {
		Status  bool
		Message string
		Data    map[string][]payment.GetPaymentsResponse
	}{}
	//result := response.ApiResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	finalData, ok = result.Data["payments"]
	if !ok {
		t.Error(fmt.Errorf("payments key does not exists in response"))
	}

	assert.Equal(t.T(), 2, len(finalData))

	for _, p := range finalData {
		switch p.ID {
		case 3:
			assert.Equal(t.T(), "BTC", p.Coin)
			assert.Equal(t.T(), "completed", p.Status)
			assert.Equal(t.T(), "btcToAddress1", p.Address)
			assert.Equal(t.T(), "0.01000000", p.Amount)
			assert.Equal(t.T(), "btcTxId1", p.TxID)
			assert.Equal(t.T(), "withdraw", p.Type)
		case 4:
			assert.Equal(t.T(), "ETH", p.Coin)
			assert.Equal(t.T(), "failed", p.Status)
			assert.Equal(t.T(), "ethToAddress1", p.Address)
			assert.Equal(t.T(), "0.10000000", p.Amount)
			assert.Equal(t.T(), "", p.TxID)
			assert.Equal(t.T(), "withdraw", p.Type)
		default:
			t.Fail("we should not be in default case")
		}

	}

	//test with pagination
	queryParams = url.Values{}
	queryParams.Set("page", "0")
	queryParams.Set("page_size", "2")
	queryParams.Set("type", "withdraw")
	queryParams.Set("code", "btc")
	paramsString = queryParams.Encode()
	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/crypto-payment?"+paramsString, nil)
	token = "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result = struct {
		Status  bool
		Message string
		Data    map[string][]payment.GetPaymentsResponse
	}{}
	//result := response.ApiResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)
	finalData, ok = result.Data["payments"]
	if !ok {
		t.Error(fmt.Errorf("payments key does not exists in response"))
	}

	assert.Equal(t.T(), 1, len(finalData))

	for _, p := range finalData {
		switch p.ID {
		case 3:
			assert.Equal(t.T(), "BTC", p.Coin)
			assert.Equal(t.T(), "completed", p.Status)
			assert.Equal(t.T(), "btcToAddress1", p.Address)
			assert.Equal(t.T(), "0.01000000", p.Amount)
			assert.Equal(t.T(), "btcTxId1", p.TxID)
			assert.Equal(t.T(), "withdraw", p.Type)
		default:
			t.Fail("we should not be in default case")
		}

	}

}

func (t *PaymentTests) TestGetPaymentDetail() {
	p1 := &payment.Payment{
		ID:        1,
		UserID:    t.userActor.ID,
		CoinID:    2,
		Type:      payment.TypeWithdraw,
		Status:    payment.StatusCompleted,
		Code:      "BTC",
		ToAddress: sql.NullString{String: "btcToAddress1", Valid: true},
		Amount:    sql.NullString{String: "0.01000000", Valid: true},
		FeeAmount: sql.NullString{String: "0.00100000", Valid: true},
		TxID:      sql.NullString{String: "btcTxId1", Valid: true},
		UpdatedAt: time.Now().Add(-2 * time.Second),
	}

	err := t.db.Create(p1).Error
	if err != nil {
		t.Fail(err.Error())
	}

	queryParams := url.Values{}
	queryParams.Set("id", "1")
	paramsString := queryParams.Encode()
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/crypto-payment/detail?"+paramsString, nil)
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    payment.DetailResponse
	}{}
	//result := response.ApiResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	p := result.Data
	assert.Equal(t.T(), "btcToAddress1", p.Address)
	assert.Equal(t.T(), "", p.RejectionReason)
	assert.Equal(t.T(), "btcTxId1", p.TxID)

	//scenario when id not found
	queryParams = url.Values{}
	queryParams.Set("id", "3")
	paramsString = queryParams.Encode()

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/crypto-payment/detail?"+paramsString, nil)
	token = "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	failResult := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &failResult)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), false, failResult.Status)
	assert.Equal(t.T(), "payment not found", failResult.Message)

}

type withdrawValidationFailedScenarios struct {
	data         string
	reason       string
	errorMessage string
}

func (t *PaymentTests) TestPreWithdraw_ValidationFail() {
	failedScenarios := []withdrawValidationFailedScenarios{
		{
			data:         `{"code":"","amount":"0.01","address":"btcAddress1"}`,
			reason:       "code not provided",
			errorMessage: "code is required",
		},
		{
			data:         `{"code":"BTC","amount":"","address":"btcAddress1"}`,
			reason:       "amount not provided",
			errorMessage: "amount is required",
		},
		{
			data:         `{"code":"BTC","amount":"-0.1","address":"btcAddress1"}`,
			reason:       "amount is negative",
			errorMessage: "amount is not correct",
		},
		{
			data:         `{"code":"BTC","amount":"0.01","address":""}`,
			reason:       "address not provided",
			errorMessage: "address is required",
		},
		{
			data:         `{"code":"btcq","amount":"0.01","address":"btcAddress1"}`,
			reason:       "wrong code",
			errorMessage: "coin not found",
		},
		{
			data:         `{"code":"usdt","amount":"0.01","address":"usdtAddress1","network":"ethq"}`,
			reason:       "wrong network code",
			errorMessage: "network not found",
		},
	}

	for _, item := range failedScenarios {
		res := httptest.NewRecorder()
		body := []byte(item.data)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/pre-withdraw", bytes.NewReader(body))
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

func (t *PaymentTests) TestPreWithdraw_UserStatusIsNotVerified() {
	//update user status to not verified
	updatingUser := &user.User{
		ID:     t.userActor.ID,
		Status: user.StatusRegistered,
	}

	err := t.db.Model(updatingUser).Updates(updatingUser).Error
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()
	data := `{"code":"btc","amount":"0.001","address":"btcAddress1"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/pre-withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "user account is not verified", result.Message)

}

func (t *PaymentTests) TestPreWithdraw_GoogleTwoDisableIn24Hour() {
	//update user status to not verified
	updatingUser := &user.User{
		ID:                  t.userActor.ID,
		Google2faDisabledAt: sql.NullTime{Time: time.Now().Add(-2 * time.Hour), Valid: true},
	}

	err := t.db.Model(updatingUser).Updates(updatingUser).Error
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()
	data := `{"code":"btc","amount":"0.001","address":"btcAddress1"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/pre-withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "for the security reasons, after disabling/enabling 2fa the withdraw request is not allowed for 24 hours", result.Message)

}

func (t *PaymentTests) TestPreWithdraw_NoPermission() {
	err := t.db.Where("user_id = ?  and user_permission_id = ?", t.userActor.ID, 2).Delete(user.UsersPermissions{}).Error
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()
	data := `{"code":"btc","amount":"0.001","address":"btcAddress1"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/pre-withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "withdraw permission is not granted", result.Message)

}

func (t *PaymentTests) TestPreWithdraw_LessThanMinimumAmount() {
	res := httptest.NewRecorder()
	data := `{"code":"btc","amount":"0.0001","address":"btcAddress1"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/pre-withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "minimum withdraw is: 0.001", result.Message)

}

func (t *PaymentTests) TestPreWithdraw_MoreThanMaximum() {
	res := httptest.NewRecorder()
	data := `{"code":"btc","amount":"10.5","address":"btcAddress1"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/pre-withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "maximum withdraw is: 5.0", result.Message)

}

func (t *PaymentTests) TestPreWithdraw_AccountIsInReadOnlyMode() {
	//insert user config in database
	userConfig := user.Config{
		ID:         1,
		UserID:     t.userActor.ID,
		IsReadOnly: true,
	}
	t.db.Create(&userConfig)

	res := httptest.NewRecorder()
	data := `{"code":"btc","amount":"0.01","address":"btcAddress1"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/pre-withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "this account is in read only mode", result.Message)

}

func (t *PaymentTests) TestPreWithdraw_AddressIsNotInWhitelist() {
	//insert user config in database
	userConfig := user.Config{
		ID:                 1,
		UserID:             t.userActor.ID,
		IsWhiteListEnabled: true,
	}
	t.db.Create(&userConfig)

	res := httptest.NewRecorder()
	data := `{"code":"btc","amount":"0.01","address":"btcAddress1"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/pre-withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "this address is not in white list", result.Message)

}

func (t *PaymentTests) TestPreWithdraw_UserBalanceIsNotEnough() {
	//insert user balance
	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.009",
		FrozenAmount:  "0",
	}
	t.db.Create(btcUb)

	//insert user config in database
	userConfig := user.Config{
		ID:     1,
		UserID: t.userActor.ID,
	}
	t.db.Create(&userConfig)

	res := httptest.NewRecorder()
	data := `{"code":"btc","amount":"0.01","address":"btcAddress1"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/pre-withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "user balance is not enough to withdraw this much", result.Message)

}

func (t *PaymentTests) TestPreWithdraw_Successful_TwoFaAndEmailCodeNotNeeded_UserConfigExists() {
	//insert user balance
	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.1",
		FrozenAmount:  "0",
	}
	t.db.Create(btcUb)

	//insert user config in database
	userConfig := user.Config{
		ID:                                    1,
		UserID:                                t.userActor.ID,
		IsTwoFaVerificationForWithdrawEnabled: false,
		IsEmailVerificationForWithdrawEnabled: false,
	}
	t.db.Create(&userConfig)

	res := httptest.NewRecorder()
	data := `{"code":"btc","amount":"0.01","address":"btcAddress1"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/pre-withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    payment.PreWithdrawResponse
	}{}
	//result := response.ApiResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), false, result.Data.Need2fa)
	assert.Equal(t.T(), true, result.Data.NeedEmailCode)

}

func (t *PaymentTests) TestPreWithdraw_Successful_TwoFaAndEmailCodeBothNeeded_UserConfigExists() {
	//insert user balance
	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.1",
		FrozenAmount:  "0",
	}
	t.db.Create(btcUb)

	//insert user config in database
	userConfig := user.Config{
		ID:                                    1,
		UserID:                                t.userActor.ID,
		IsTwoFaVerificationForWithdrawEnabled: true,
		IsEmailVerificationForWithdrawEnabled: true,
	}
	t.db.Create(&userConfig)

	res := httptest.NewRecorder()
	data := `{"code":"btc","amount":"0.01","address":"btcAddress1"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/pre-withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    payment.PreWithdrawResponse
	}{}
	//result := response.ApiResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), true, result.Data.Need2fa)
	assert.Equal(t.T(), true, result.Data.NeedEmailCode)

	//check redis if code generated for user
	userIDString := strconv.Itoa(t.userActor.ID)
	key := "withdraw-confirmation:" + userIDString
	redisData, err := t.redisClient.HGetAll(context.Background(), key).Result()
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), userIDString, redisData["userId"])
	assert.Equal(t.T(), "0.01", redisData["amount"])
	assert.Equal(t.T(), "BTC", redisData["coin"])
	assert.Equal(t.T(), "btcAddress1", redisData["address"])

	expiredAtString := redisData["expiredAt"]
	expiredAt, err := strconv.ParseInt(expiredAtString, 10, 64)
	assert.Nil(t.T(), err)

	assert.GreaterOrEqual(t.T(), expiredAt, time.Now().Unix()+(30*60)-10) //30 minutes minus 10 seconds to be sure our assertion would be true
}

func (t *PaymentTests) TestPreWithdraw_Successful_TwoFaNeeded_UserConfigDoesNotExists() {
	//insert user balance
	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.1",
		FrozenAmount:  "0",
	}
	t.db.Create(btcUb)

	updatingUser := &user.User{
		ID:             t.userActor.ID,
		IsTwoFaEnabled: true,
	}
	err := t.db.Model(updatingUser).Updates(updatingUser).Error
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()
	data := `{"code":"btc","amount":"0.01","address":"btcAddress1"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/pre-withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    payment.PreWithdrawResponse
	}{}
	//result := response.ApiResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), true, result.Data.Need2fa)
	assert.Equal(t.T(), true, result.Data.NeedEmailCode)
}

func (t *PaymentTests) TestWithdraw_ValidationFail() {
	failedScenarios := []withdrawValidationFailedScenarios{
		{
			data:         `{"code":"","amount":"0.01","address":"btcAddress1"}`,
			reason:       "code not provided",
			errorMessage: "code is required",
		},
		{
			data:         `{"code":"BTC","amount":"","address":"btcAddress1"}`,
			reason:       "amount not provided",
			errorMessage: "amount is required",
		},
		{
			data:         `{"code":"BTC","amount":"-0.1","address":"btcAddress1"}`,
			reason:       "amount is negative",
			errorMessage: "amount is not correct",
		},
		{
			data:         `{"code":"BTC","amount":"0.01","address":""}`,
			reason:       "address not provided",
			errorMessage: "address is required",
		},
		{
			data:         `{"code":"btcq","amount":"0.01","address":"btcAddress1"}`,
			reason:       "wrong code",
			errorMessage: "coin not found",
		},
		{
			data:         `{"code":"usdt","amount":"0.01","address":"usdtAddress1","network":"ethq"}`,
			reason:       "wrong network code",
			errorMessage: "network not found",
		},
	}

	for _, item := range failedScenarios {
		res := httptest.NewRecorder()
		body := []byte(item.data)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/withdraw", bytes.NewReader(body))
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

func (t *PaymentTests) TestWithdraw_UserStatusIsNotVerified() {
	//update user status to not verified
	updatingUser := &user.User{
		ID:     t.userActor.ID,
		Status: user.StatusRegistered,
	}

	err := t.db.Model(updatingUser).Updates(updatingUser).Error
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()
	data := `{"code":"btc","amount":"0.001","address":"btcAddress1"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "user account is not verified", result.Message)

}

func (t *PaymentTests) TestWithdraw_GoogleTwoDisableIn24Hour() {
	//update user status to not verified
	updatingUser := &user.User{
		ID:                  t.userActor.ID,
		Google2faDisabledAt: sql.NullTime{Time: time.Now().Add(-2 * time.Hour), Valid: true},
	}

	err := t.db.Model(updatingUser).Updates(updatingUser).Error
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()
	data := `{"code":"btc","amount":"0.001","address":"btcAddress1"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "for the security reasons, after disabling/enabling 2fa the withdraw request is not allowed for 24 hours", result.Message)

}

func (t *PaymentTests) TestWithdraw_NoPermission() {
	err := t.db.Where("user_id = ?  and user_permission_id = ?", t.userActor.ID, 2).Delete(user.UsersPermissions{}).Error
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()
	data := `{"code":"btc","amount":"0.001","address":"btcAddress1"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "withdraw permission is not granted", result.Message)

}

func (t *PaymentTests) TestWithdraw_LessThanMinimumAmount() {
	res := httptest.NewRecorder()
	data := `{"code":"btc","amount":"0.0001","address":"btcAddress1"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "minimum withdraw is: 0.001", result.Message)

}

func (t *PaymentTests) TestWithdraw_MoreThanMaximum() {
	res := httptest.NewRecorder()
	data := `{"code":"btc","amount":"10.5","address":"btcAddress1"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "maximum withdraw is: 5.0", result.Message)

}

func (t *PaymentTests) TestWithdraw_AccountIsInReadOnlyMode() {
	//insert user config in database
	userConfig := user.Config{
		ID:         1,
		UserID:     t.userActor.ID,
		IsReadOnly: true,
	}
	t.db.Create(&userConfig)

	res := httptest.NewRecorder()
	data := `{"code":"btc","amount":"0.01","address":"btcAddress1"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "this account is in read only mode", result.Message)

}

func (t *PaymentTests) TestWithdraw_AddressIsNotInWhitelist() {
	//insert user config in database
	userConfig := user.Config{
		ID:                 1,
		UserID:             t.userActor.ID,
		IsWhiteListEnabled: true,
	}
	t.db.Create(&userConfig)

	res := httptest.NewRecorder()
	data := `{"code":"btc","amount":"0.01","address":"btcAddress1"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "this address is not in white list", result.Message)

}

//func (t *PaymentTests) TestWithdraw_UserBalanceIsNotEnough() {
//	//insert user balance
//	btcUb := &userbalance.UserBalance{
//		UserId:        t.userActor.ID,
//		CoinId:        2, //for btc from currency seed
//		FrozenBalance: "",
//		BalanceCoin:   "BTC",
//		Status:        userbalance.StatusEnabled,
//		Amount:        "0.009",
//		FrozenAmount:  "0",
//	}
//	t.db.Create(btcUb)
//
//	//insert user config in database
//	userConfig := user.Config{
//		ID:     1,
//		UserId: t.userActor.ID,
//	}
//	t.db.Create(&userConfig)
//
//
//
//
//	res := httptest.NewRecorder()
//	data := `{"code":"btc","amount":"0.01","address":"btcAddress1"}`
//	body := []byte(data)
//	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/withdraw", bytes.NewReader(body))
//	token := "Bearer " + t.userActor.Token
//	req.Header.Set("Authorization", token)
//	t.httpServer.ServeHTTP(res, req)
//	result := response.ApiResponse{}
//	err := json.Unmarshal(res.Body.Bytes(), &result)
//	if err != nil {
//		t.Fail(err.Error())
//	}
//
//	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
//	assert.Equal(t.T(), "user balance is not enough to withdraw this much", result.Message)
//
//}

func (t *PaymentTests) TestWithdraw_TwoFaCodeIsNotCorrect() {
	//insert user balance
	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.1",
		FrozenAmount:  "0",
	}
	t.db.Create(btcUb)

	//insert user config in database
	userConfig := user.Config{
		ID:                                    1,
		UserID:                                t.userActor.ID,
		IsTwoFaVerificationForWithdrawEnabled: true,
	}
	t.db.Create(&userConfig)

	res := httptest.NewRecorder()
	data := `{"code":"btc","amount":"0.01","address":"btcAddress1","2fa_code":"111111"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "2fa code is not correct", result.Message)
}

func (t *PaymentTests) TestWithdraw_EmailCodeIsNotCorrect() {
	//insert user balance
	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.1",
		FrozenAmount:  "0",
	}
	t.db.Create(btcUb)

	//insert user config in database
	userConfig := user.Config{
		ID:                                    1,
		UserID:                                t.userActor.ID,
		IsEmailVerificationForWithdrawEnabled: true,
	}
	t.db.Create(&userConfig)

	res := httptest.NewRecorder()
	data := `{"code":"btc","amount":"0.01","address":"btcAddress1","email_code":"111111"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "email confirmation code is not correct", result.Message)

}

//func (t *PaymentTests) TestWithdraw_Successful_NoTwoFaAndEmailCodeNeeded() {
//	//insert user balance
//	btcUb := &userbalance.UserBalance{
//		UserId:        t.userActor.ID,
//		CoinId:        2, //for btc from currency seed
//		FrozenBalance: "",
//		BalanceCoin:   "BTC",
//		Status:        userbalance.StatusEnabled,
//		Amount:        "0.1",
//		FrozenAmount:  "0",
//	}
//	t.db.Create(btcUb)
//
//	//insert user config in database
//	userConfig := user.Config{
//		ID:     1,
//		UserId: t.userActor.ID,
//	}
//
//	t.db.Create(&userConfig)
//
//	res := httptest.NewRecorder()
//	data := `{"code":"btc","amount":"0.01","address":"btcAddress1"}`
//	body := []byte(data)
//	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/withdraw", bytes.NewReader(body))
//	token := "Bearer " + t.userActor.Token
//	req.Header.Set("Authorization", token)
//	t.httpServer.ServeHTTP(res, req)
//
//	result := struct {
//		Status  bool
//		Message string
//		Data    map[string][]payment.GetPaymentsResponse
//	}{}
//	//result := response.ApiResponse{}
//	err := json.Unmarshal(res.Body.Bytes(), &result)
//	if err != nil {
//		t.Fail(err.Error())
//	}
//
//	assert.Equal(t.T(), http.StatusOK, res.Code)
//	finalData, ok := result.Data["payments"]
//	if !ok {
//		t.Error(fmt.Errorf("payments key does not exists in response"))
//	}
//
//	assert.Equal(t.T(), 1, len(finalData))
//
//	p := finalData[0]
//	assert.Equal(t.T(), "BTC", p.Coin)
//	assert.Equal(t.T(), "in progress", p.Status)
//	assert.Equal(t.T(), "btcAddress1", p.Address)
//	assert.Equal(t.T(), "0.01000000", p.Amount)
//	assert.Equal(t.T(), "", p.TxId)
//	assert.Equal(t.T(), "withdraw", p.Type)
//
//	// check payment exists in database
//	id := p.Id
//	createdPayment := &payment.Payment{}
//	err = t.db.Where(payment.Payment{ID: id}).First(createdPayment).Error
//
//	assert.Equal(t.T(), t.userActor.ID, createdPayment.UserId)
//	assert.Equal(t.T(), int64(2), createdPayment.CoinId)
//	assert.Equal(t.T(), "WITHDRAW", createdPayment.Type)
//	assert.Equal(t.T(), "CREATED", createdPayment.Status)
//	assert.Equal(t.T(), "BTC", createdPayment.Code)
//	assert.Equal(t.T(), "", createdPayment.FromAddress.String)
//	assert.Equal(t.T(), "btcAddress1", createdPayment.ToAddress.String)
//	assert.Equal(t.T(), "", createdPayment.TxId.String)
//	assert.Equal(t.T(), "", createdPayment.BlockchainNetwork.String)
//	assert.Equal(t.T(), "PENDING", createdPayment.AdminStatus.String)
//	assert.Equal(t.T(), "0.01000000", createdPayment.Amount.String)
//	assert.Equal(t.T(), "0.00010000", createdPayment.FeeAmount.String)
//
//	// check extra payment info exists in database
//	extraInfo := &payment.ExtraInfo{}
//	err = t.db.Where(payment.ExtraInfo{PaymentId: id}).First(extraInfo).Error
//	if err != nil {
//		t.Fail(err.Error())
//	}
//	assert.Greater(t.T(), extraInfo.ID, int64(0))
//
//	//check user balance
//	updatedBtcUb := &userbalance.UserBalance{}
//	err = t.db.Where(userbalance.UserBalance{ID: btcUb.ID}).First(updatedBtcUb).Error
//	if err != nil {
//		t.Fail(err.Error())
//	}
//	assert.Equal(t.T(), "0.01000000", updatedBtcUb.FrozenAmount)
//
//}

func (t *PaymentTests) TestWithdraw_Successful_TwoFaAndEmailCodeNeeded() {
	//insert user balance
	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.1",
		FrozenAmount:  "0",
	}
	t.db.Create(btcUb)

	//insert user config in database
	userConfig := user.Config{
		ID:                                    1,
		UserID:                                t.userActor.ID,
		IsTwoFaVerificationForWithdrawEnabled: true,
		IsEmailVerificationForWithdrawEnabled: true,
	}

	t.db.Create(&userConfig)

	//first we call pre-withdraw
	res := httptest.NewRecorder()
	data := `{"code":"btc","amount":"0.01","address":"btcAddress1"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/pre-withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	preWithdrawResult := struct {
		Status  bool
		Message string
		Data    payment.PreWithdrawResponse
	}{}
	//result := response.ApiResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &preWithdrawResult)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), true, preWithdrawResult.Data.Need2fa)
	assert.Equal(t.T(), true, preWithdrawResult.Data.NeedEmailCode)

	//get code from redis
	userIDString := strconv.Itoa(t.userActor.ID)
	key := "withdraw-confirmation:" + userIDString
	redisData, err := t.redisClient.HGetAll(context.Background(), key).Result()
	assert.Nil(t.T(), err)
	emailCode := redisData["code"]

	//get 2Fa code
	twoFaCode, err := totp.GenerateCode("HWOAQZBGXCKJZQVH", time.Now()) //this secret is for the userActor set in main_test.go
	if err != nil {
		t.Fail(err.Error())
	}
	res = httptest.NewRecorder()
	data = `{"code":"btc","amount":"0.01","address":"btcAddress1","email_code":"` + emailCode + `","2fa_code":"` + twoFaCode + `"}`
	body = []byte(data)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/withdraw", bytes.NewReader(body))
	token = "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    map[string][]payment.GetPaymentsResponse
	}{}
	//result := response.ApiResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)

	finalData, ok := result.Data["payments"]
	if !ok {
		t.Error(fmt.Errorf("payments key does not exists in response"))
	}

	assert.Equal(t.T(), 1, len(finalData))

	p := finalData[0]
	assert.Equal(t.T(), "BTC", p.Coin)
	assert.Equal(t.T(), "pending", p.Status)
	assert.Equal(t.T(), "btcAddress1", p.Address)
	assert.Equal(t.T(), "0.01000000", p.Amount)
	assert.Equal(t.T(), "", p.TxID)
	assert.Equal(t.T(), "withdraw", p.Type)

	// check payment exists in database
	id := p.ID
	createdPayment := &payment.Payment{}
	err = t.db.Where(payment.Payment{ID: id}).First(createdPayment).Error

	assert.Equal(t.T(), t.userActor.ID, createdPayment.UserID)
	assert.Equal(t.T(), int64(2), createdPayment.CoinID)
	assert.Equal(t.T(), "WITHDRAW", createdPayment.Type)
	assert.Equal(t.T(), "CREATED", createdPayment.Status)
	assert.Equal(t.T(), "BTC", createdPayment.Code)
	assert.Equal(t.T(), "", createdPayment.FromAddress.String)
	assert.Equal(t.T(), "btcAddress1", createdPayment.ToAddress.String)
	assert.Equal(t.T(), "", createdPayment.TxID.String)
	assert.Equal(t.T(), "", createdPayment.BlockchainNetwork.String)
	assert.Equal(t.T(), "PENDING", createdPayment.AdminStatus.String)
	assert.Equal(t.T(), "0.01000000", createdPayment.Amount.String)
	assert.Equal(t.T(), "0.00010000", createdPayment.FeeAmount.String)

	// check extra payment info exists in database
	extraInfo := &payment.ExtraInfo{}
	err = t.db.Where(payment.ExtraInfo{PaymentID: id}).First(extraInfo).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Greater(t.T(), extraInfo.ID, int64(0))

	//check user balance
	updatedBtcUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: btcUb.ID}).First(updatedBtcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "0.01000000", updatedBtcUb.FrozenAmount)

}

func (t *PaymentTests) TestWithdraw_Successful_USDT_TRC_NETWORK_DIFFERENT_FEE() {
	//insert user balance
	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "110",
		FrozenAmount:  "0",
	}
	t.db.Create(usdtUb)

	//insert user config in database
	userConfig := user.Config{
		ID:                                    1,
		UserID:                                t.userActor.ID,
		IsTwoFaVerificationForWithdrawEnabled: false,
		IsEmailVerificationForWithdrawEnabled: true,
	}

	t.db.Create(&userConfig)

	//first we call pre-withdraw
	res := httptest.NewRecorder()
	data := `{"code":"USDT","amount":"100","address":"usdtAddress1","network":"TRX"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/pre-withdraw", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	preWithdrawResult := struct {
		Status  bool
		Message string
		Data    payment.PreWithdrawResponse
	}{}
	//result := response.ApiResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &preWithdrawResult)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), false, preWithdrawResult.Data.Need2fa)
	assert.Equal(t.T(), true, preWithdrawResult.Data.NeedEmailCode)

	//get code from redis
	userIDString := strconv.Itoa(t.userActor.ID)
	key := "withdraw-confirmation:" + userIDString
	redisData, err := t.redisClient.HGetAll(context.Background(), key).Result()
	assert.Nil(t.T(), err)
	emailCode := redisData["code"]

	res = httptest.NewRecorder()
	data = `{"code":"USDT","amount":"100","address":"usdtAddress1","network":"TRX","email_code":"` + emailCode + `"}`
	body = []byte(data)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/withdraw", bytes.NewReader(body))
	token = "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    map[string][]payment.GetPaymentsResponse
	}{}
	//result := response.ApiResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)

	finalData, ok := result.Data["payments"]
	if !ok {
		t.Error(fmt.Errorf("payments key does not exists in response"))
	}

	assert.Equal(t.T(), 1, len(finalData))

	p := finalData[0]
	assert.Equal(t.T(), "USDT", p.Coin)
	assert.Equal(t.T(), "pending", p.Status)
	assert.Equal(t.T(), "usdtAddress1", p.Address)
	assert.Equal(t.T(), "100.00000000", p.Amount)
	assert.Equal(t.T(), "", p.TxID)
	assert.Equal(t.T(), "withdraw", p.Type)

	// check payment exists in database
	id := p.ID
	createdPayment := &payment.Payment{}
	err = t.db.Where(payment.Payment{ID: id}).First(createdPayment).Error

	assert.Equal(t.T(), t.userActor.ID, createdPayment.UserID)
	assert.Equal(t.T(), int64(1), createdPayment.CoinID)
	assert.Equal(t.T(), "WITHDRAW", createdPayment.Type)
	assert.Equal(t.T(), "CREATED", createdPayment.Status)
	assert.Equal(t.T(), "USDT", createdPayment.Code)
	assert.Equal(t.T(), "", createdPayment.FromAddress.String)
	assert.Equal(t.T(), "usdtAddress1", createdPayment.ToAddress.String)
	assert.Equal(t.T(), "", createdPayment.TxID.String)
	assert.Equal(t.T(), "TRX", createdPayment.BlockchainNetwork.String)
	assert.Equal(t.T(), "PENDING", createdPayment.AdminStatus.String)
	assert.Equal(t.T(), "100.00000000", createdPayment.Amount.String)
	assert.Equal(t.T(), "2.5", createdPayment.FeeAmount.String)

	// check extra payment info exists in database
	extraInfo := &payment.ExtraInfo{}
	err = t.db.Where(payment.ExtraInfo{PaymentID: id}).First(extraInfo).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Greater(t.T(), extraInfo.ID, int64(0))

	//check user balance
	updatedUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "100.00000000", updatedUsdtUb.FrozenAmount)
}

func (t *PaymentTests) TestCancelWithdraw_Fail_PaymentDoesNotExist() {
	data := `{"id":2}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/cancel", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := response.APIResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.T().Error(err)
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "withdraw not found", result.Message)

}

func (t *PaymentTests) TestCancelWithdraw_Fail_PaymentIsNotWithdraw() {
	res := httptest.NewRecorder()
	p := &payment.Payment{
		UserID:            t.userActor.ID,
		CoinID:            1,
		Type:              payment.TypeDeposit,
		Status:            payment.StatusCreated,
		Code:              "USDT",
		ToAddress:         sql.NullString{String: "usdtToAddress1", Valid: true},
		BlockchainNetwork: sql.NullString{String: "ETH", Valid: true},
		Amount:            sql.NullString{String: "100.00000000", Valid: true},
		FeeAmount:         sql.NullString{String: "10.00000000", Valid: true},
		UpdatedAt:         time.Now().Add(-2 * time.Second),
	}

	err := t.db.Create(p).Error
	if err != nil {
		t.Fail(err.Error())
	}

	idString := strconv.FormatInt(p.ID, 10)
	data := `{"id":` + idString + `}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/cancel", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.T().Error(err)
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "withdraw not found", result.Message)
}

func (t *PaymentTests) TestCancelWithdraw_PaymentStatusIsNotCreated() {
	res := httptest.NewRecorder()
	p := &payment.Payment{
		UserID:            t.userActor.ID,
		CoinID:            1,
		Type:              payment.TypeWithdraw,
		Status:            payment.StatusInProgress,
		Code:              "USDT",
		ToAddress:         sql.NullString{String: "usdtToAddress1", Valid: true},
		BlockchainNetwork: sql.NullString{String: "ETH", Valid: true},
		Amount:            sql.NullString{String: "100.00000000", Valid: true},
		FeeAmount:         sql.NullString{String: "10.00000000", Valid: true},
		UpdatedAt:         time.Now().Add(-2 * time.Second),
	}

	err := t.db.Create(p).Error
	if err != nil {
		t.Fail(err.Error())
	}

	idString := strconv.FormatInt(p.ID, 10)
	data := `{"id":` + idString + `}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/cancel", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.T().Error(err)
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "withdraw can't be cancelled now", result.Message)
}

func (t *PaymentTests) TestCancelWithdraw_Successful() {
	res := httptest.NewRecorder()

	//insert user balance in db
	//insert userBalance for user
	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "1000.00",
		FrozenAmount:  "500.00",
	}

	err := t.db.Create(usdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	p := &payment.Payment{
		UserID:            t.userActor.ID,
		CoinID:            1,
		Type:              payment.TypeWithdraw,
		Status:            payment.StatusCreated,
		Code:              "USDT",
		ToAddress:         sql.NullString{String: "usdtToAddress1", Valid: true},
		BlockchainNetwork: sql.NullString{String: "ETH", Valid: true},
		Amount:            sql.NullString{String: "100.00000000", Valid: true},
		FeeAmount:         sql.NullString{String: "10.00000000", Valid: true},
		UpdatedAt:         time.Now().Add(-2 * time.Second),
	}

	err = t.db.Create(p).Error
	if err != nil {
		t.Fail(err.Error())
	}

	idString := strconv.FormatInt(p.ID, 10)
	data := `{"id":` + idString + `}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crypto-payment/cancel", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

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

	assert.Equal(t.T(), payment.StatusUserCanceled, updatedPayment.Status)

	updatedUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "400.00000000", updatedUsdtUb.FrozenAmount)

}

func (t *PaymentTests) TestHandleWalletCallback_InternalTransfer() {
	res := httptest.NewRecorder()
	cr := &CryptoBalance{
		Type: "EXTERNAL",
	}
	err := t.db.Create(cr).Error
	if err != nil {
		t.Fail(err.Error())
	}
	internalTransfer := &payment.InternalTransfer{
		Amount:        "1.0",
		Status:        payment.StatusCreated,
		Network:       "ETH",
		FromBalanceID: cr.ID,
	}
	err = t.db.Create(internalTransfer).Error
	if err != nil {
		t.Fail(err.Error())
	}
	internalTransferIDString := strconv.FormatInt(internalTransfer.ID, 10)
	data := fmt.Sprintf(`{
		"code" : "USDT",
		"amount" : "100.0",
		"type" : "DEPOSIT",
		"status" : "COMPLETED",
		"from_address" : "fromAddress",
		"to_address" : "user1USDTAddress",
		"tx_id" : "txIdFromCallback",
		"meta" : "{\"internal_transfer_id\":\"%s\"}",
		"network" : "ETH"
	}`, internalTransferIDString)
	body := []byte(data)
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

	updatedInternalTransfer := &payment.InternalTransfer{}
	err = t.db.Where(payment.InternalTransfer{ID: internalTransfer.ID}).First(updatedInternalTransfer).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), payment.StatusCompleted, updatedInternalTransfer.Status)
}

func (t *PaymentTests) TestHandleWalletCallback_Deposit_AlreadyDoesNotExists_StatusCreated() {
	res := httptest.NewRecorder()
	//first we insert userbalance
	user := getNewUserActor()
	usdtUb := &userbalance.UserBalance{
		UserID:        user.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "1000.00",
		FrozenAmount:  "500.00",
		Address:       sql.NullString{String: "user1USDTAddress", Valid: true},
	}
	err := t.db.Create(usdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	data := `{
		"code" : "USDT",
		"amount" : "100.0",
		"type" : "DEPOSIT",
		"status" : "CREATED",
		"from_address" : "fromAddress",
		"to_address" : "user1USDTAddress",
		"tx_id" : "txIdFromCallback",
		"network" : "ETH"
	}`
	body := []byte(data)
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

	p := &payment.Payment{}
	err = t.db.Where(payment.Payment{UserID: user.ID}).First(p).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), payment.StatusCreated, p.Status)
	assert.Equal(t.T(), "100.00000000", p.Amount.String)
	assert.Equal(t.T(), "fromAddress", p.FromAddress.String)
	assert.Equal(t.T(), "user1USDTAddress", p.ToAddress.String)
	assert.Equal(t.T(), "USDT", p.Code)
	assert.Equal(t.T(), int64(1), p.CoinID)
	assert.Equal(t.T(), "", p.FeeAmount.String)
	assert.Equal(t.T(), "txIdFromCallback", p.TxID.String)
	assert.Equal(t.T(), "DEPOSIT", p.Type)
	assert.Equal(t.T(), "ETH", p.BlockchainNetwork.String)

	ei := &payment.ExtraInfo{}
	err = t.db.Where(payment.ExtraInfo{PaymentID: p.ID}).First(ei).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.NotEmpty(t.T(), ei.Price.String)
	assert.NotEmpty(t.T(), ei.BtcPrice.String)

	nonUpdatedUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(nonUpdatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "1000.00", nonUpdatedUsdtUb.Amount)
	assert.Equal(t.T(), "500.00", nonUpdatedUsdtUb.FrozenAmount)
}

// instead of erc20 this test simulate trc20
func (t *PaymentTests) TestHandleWalletCallback_Deposit_AlreadyDoesNotExists_StatusCreated_OtherNetwork() {
	res := httptest.NewRecorder()
	//first we insert userbalance
	user := getNewUserActor()
	usdtUb := &userbalance.UserBalance{
		UserID:         user.ID,
		CoinID:         1, //for usdt from currency seed
		FrozenBalance:  "",
		BalanceCoin:    "USDT",
		Status:         userbalance.StatusEnabled,
		Amount:         "1000.00",
		FrozenAmount:   "500.00",
		Address:        sql.NullString{String: "someOtheruserUSDTAddress", Valid: true},
		OtherAddresses: sql.NullString{String: `[{"code":"TRX","address":"user2USDTAddress"}]`, Valid: true},
	}

	err := t.db.Create(usdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	data := `{
		"code" : "USDT",
		"amount" : "100.0",
		"type" : "DEPOSIT",
		"status" : "CREATED",
		"from_address" : "fromAddress",
		"to_address" : "user2USDTAddress",
		"tx_id" : "txIdFromCallback",
		"network" : "TRX"
	}`
	body := []byte(data)
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

	p := &payment.Payment{}
	err = t.db.Where(payment.Payment{UserID: user.ID}).First(p).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), payment.StatusCreated, p.Status)
	assert.Equal(t.T(), "100.00000000", p.Amount.String)
	assert.Equal(t.T(), "fromAddress", p.FromAddress.String)
	assert.Equal(t.T(), "user2USDTAddress", p.ToAddress.String)
	assert.Equal(t.T(), "USDT", p.Code)
	assert.Equal(t.T(), int64(1), p.CoinID)
	assert.Equal(t.T(), "", p.FeeAmount.String)
	assert.Equal(t.T(), "txIdFromCallback", p.TxID.String)
	assert.Equal(t.T(), "DEPOSIT", p.Type)
	assert.Equal(t.T(), "TRX", p.BlockchainNetwork.String)

	ei := &payment.ExtraInfo{}
	err = t.db.Where(payment.ExtraInfo{PaymentID: p.ID}).First(ei).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.NotEmpty(t.T(), ei.Price.String)
	assert.NotEmpty(t.T(), ei.BtcPrice.String)

	nonUpdatedUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(nonUpdatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "1000.00", nonUpdatedUsdtUb.Amount)
	assert.Equal(t.T(), "500.00", nonUpdatedUsdtUb.FrozenAmount)
}

func (t *PaymentTests) TestHandleWalletCallback_Deposit_AlreadyDoesNotExists_StatusCompleted() {
	res := httptest.NewRecorder()
	//first we insert userbalance
	user := getNewUserActor()
	usdtUb := &userbalance.UserBalance{
		UserID:        user.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "1000.00",
		FrozenAmount:  "500.00",
		Address:       sql.NullString{String: "user1USDTAddress", Valid: true},
	}
	err := t.db.Create(usdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	data := `{
		"code" : "USDT",
		"amount" : "100.0",
		"type" : "DEPOSIT",
		"status" : "COMPLETED",
		"from_address" : "fromAddress",
		"to_address" : "user1USDTAddress",
		"tx_id" : "txIdFromCallback",
		"network" : "ETH"
	}`
	body := []byte(data)
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

	p := &payment.Payment{}
	err = t.db.Where(payment.Payment{UserID: user.ID}).First(p).Error
	if err != nil {
		t.Fail(err.Error())
	}
	ei := &payment.ExtraInfo{}
	err = t.db.Where(payment.ExtraInfo{PaymentID: p.ID}).First(ei).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.NotEmpty(t.T(), ei.Price.String)
	assert.NotEmpty(t.T(), ei.BtcPrice.String)

	assert.Equal(t.T(), payment.StatusCompleted, p.Status)
	assert.Equal(t.T(), "100.00000000", p.Amount.String)
	assert.Equal(t.T(), "fromAddress", p.FromAddress.String)
	assert.Equal(t.T(), "user1USDTAddress", p.ToAddress.String)
	assert.Equal(t.T(), "USDT", p.Code)
	assert.Equal(t.T(), "txIdFromCallback", p.TxID.String)
	assert.Equal(t.T(), "DEPOSIT", p.Type)
	assert.Equal(t.T(), "ETH", p.BlockchainNetwork.String)

	updatedUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "1100.00000000", updatedUsdtUb.Amount)
	assert.Equal(t.T(), "500.00", updatedUsdtUb.FrozenAmount)

	//checking transaction
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
}

func (t *PaymentTests) TestHandleWalletCallback_Deposit_AlreadyExists_StatusCompleted() {
	res := httptest.NewRecorder()
	//first we insert userbalance
	user := getNewUserActor()
	usdtUb := &userbalance.UserBalance{
		UserID:        user.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "1000.00",
		FrozenAmount:  "500.00",
		Address:       sql.NullString{String: "user1USDTAddress", Valid: true},
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
}

func (t *PaymentTests) TestHandleWalletCallback_Withdraw_StatusCompleted() {
	res := httptest.NewRecorder()
	//first we insert userbalance
	user := getNewUserActor()
	usdtUb := &userbalance.UserBalance{
		UserID:        user.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "1000.00",
		FrozenAmount:  "500.00",
		Address:       sql.NullString{String: "user1USDTAddress", Valid: true},
	}
	err := t.db.Create(usdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	p := &payment.Payment{
		UserID:            user.ID,
		CoinID:            1,
		Type:              "WITHDRAW",
		Status:            "CREATED",
		Code:              "USDT",
		FromAddress:       sql.NullString{String: "fromAddress", Valid: true},
		ToAddress:         sql.NullString{String: "user1USDTAddress", Valid: true},
		TxID:              sql.NullString{String: "txIdFromCallback", Valid: true},
		Amount:            sql.NullString{String: "100.00000000", Valid: true},
		BlockchainNetwork: sql.NullString{String: "ETH", Valid: true},
		FeeAmount:         sql.NullString{String: "5.00000000", Valid: true},
	}

	err = t.db.Create(p).Error
	if err != nil {
		t.Fail(err.Error())
	}

	data := `{
		"code" : "USDT",
		"amount" : "100.0",
		"type" : "WITHDRAW",
		"status" : "COMPLETED",
		"from_address" : "fromAddress",
		"to_address" : "user1USDTAddress",
		"tx_id" : "txIdFromCallback",
		"network" : "ETH"
	}`
	body := []byte(data)
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
	assert.Equal(t.T(), "WITHDRAW", updatedPayment.Type)
	assert.Equal(t.T(), "ETH", updatedPayment.BlockchainNetwork.String)

	updatedUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "900.00000000", updatedUsdtUb.Amount)
	assert.Equal(t.T(), "400.00000000", updatedUsdtUb.FrozenAmount)

	tx := &transaction.Transaction{}
	err = t.db.Where(transaction.Transaction{Type: "WITHDRAW", PaymentID: sql.NullInt64{Int64: p.ID, Valid: true}}).First(tx).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "100.00000000", tx.Amount.String)
	assert.Equal(t.T(), int64(1), tx.CoinID)
	assert.Equal(t.T(), "USDT", tx.CoinName)
	assert.Equal(t.T(), "WITHDRAW", tx.Type)
	assert.Equal(t.T(), user.ID, tx.UserID)

	feeTx := &transaction.Transaction{}
	err = t.db.Where(transaction.Transaction{Type: "WITHDRAW_FEE", PaymentID: sql.NullInt64{Int64: p.ID, Valid: true}}).First(feeTx).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "5.00000000", feeTx.Amount.String)
	assert.Equal(t.T(), int64(1), feeTx.CoinID)
	assert.Equal(t.T(), "USDT", feeTx.CoinName)
	assert.Equal(t.T(), "WITHDRAW_FEE", feeTx.Type)
	assert.Equal(t.T(), user.ID, feeTx.UserID)
}

func (t *PaymentTests) TestHandleWalletCallback_Withdraw_StatusFailed() {
	res := httptest.NewRecorder()
	//first we insert userbalance
	user := getNewUserActor()
	usdtUb := &userbalance.UserBalance{
		UserID:        user.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "1000.00000000",
		FrozenAmount:  "500.00",
		Address:       sql.NullString{String: "user1USDTAddress", Valid: true},
	}
	err := t.db.Create(usdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	p := &payment.Payment{
		UserID:            user.ID,
		CoinID:            1,
		Type:              "WITHDRAW",
		Status:            "CREATED",
		Code:              "USDT",
		FromAddress:       sql.NullString{String: "fromAddress", Valid: true},
		ToAddress:         sql.NullString{String: "user1USDTAddress", Valid: true},
		TxID:              sql.NullString{String: "txIdFromCallback", Valid: true},
		Amount:            sql.NullString{String: "100.00000000", Valid: true},
		BlockchainNetwork: sql.NullString{String: "ETH", Valid: true},
		FeeAmount:         sql.NullString{String: "5.00000000", Valid: true},
	}

	err = t.db.Create(p).Error
	if err != nil {
		t.Fail(err.Error())
	}

	data := `{
		"code" : "USDT",
		"amount" : "100.0",
		"type" : "WITHDRAW",
		"status" : "FAILED",
		"from_address" : "fromAddress",
		"to_address" : "user1USDTAddress",
		"tx_id" : "txIdFromCallback",
		"network" : "ETH"
	}`
	body := []byte(data)
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

	assert.Equal(t.T(), payment.StatusFailed, updatedPayment.Status)
	assert.Equal(t.T(), "100.00000000", updatedPayment.Amount.String)
	assert.Equal(t.T(), "fromAddress", updatedPayment.FromAddress.String)
	assert.Equal(t.T(), "user1USDTAddress", updatedPayment.ToAddress.String)
	assert.Equal(t.T(), "USDT", updatedPayment.Code)
	assert.Equal(t.T(), "txIdFromCallback", updatedPayment.TxID.String)
	assert.Equal(t.T(), "WITHDRAW", updatedPayment.Type)
	assert.Equal(t.T(), "ETH", updatedPayment.BlockchainNetwork.String)

	updatedUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "1000.00000000", updatedUsdtUb.Amount)
	assert.Equal(t.T(), "400.00000000", updatedUsdtUb.FrozenAmount)

}

func (t *PaymentTests) TestUpdateWithdraw_AdminStatus_Recheck() {
	user := getNewUserActor()
	p := &payment.Payment{
		UserID:            user.ID,
		CoinID:            1,
		Type:              "WITHDRAW",
		Status:            "CREATED",
		Code:              "USDT",
		FromAddress:       sql.NullString{String: "fromAddress", Valid: true},
		ToAddress:         sql.NullString{String: "user1USDTAddress", Valid: true},
		TxID:              sql.NullString{String: "txIdFromCallback", Valid: true},
		Amount:            sql.NullString{String: "100.00000000", Valid: true},
		BlockchainNetwork: sql.NullString{String: "ETH", Valid: true},
		FeeAmount:         sql.NullString{String: "5.00000000", Valid: true},
	}

	err := t.db.Create(p).Error
	if err != nil {
		t.Fail(err.Error())
	}

	ei := &payment.ExtraInfo{
		PaymentID: p.ID,
	}

	err = t.db.Create(ei).Error
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()
	data := fmt.Sprintf(`{
		"id" : %d,
		"admin_status" : "recheck",
		"fee" : "10",
		"network_fee" : "0.001"
	}`, p.ID)
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payment/update-withdraw", bytes.NewReader(body))
	token := "Bearer " + t.adminUserActor.Token
	req.Header.Set("Authorization", token)
	t.adminHTTPServer.ServeHTTP(res, req)

	assert.Equal(t.T(), http.StatusOK, res.Code)

	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.T().Error(err)
	}
	updatedPayment := &payment.Payment{}
	err = t.db.Where(payment.Payment{ID: p.ID}).First(updatedPayment).Error
	if err != nil {
		t.Fail(err.Error())
	}
	updatedExtraInfo := &payment.ExtraInfo{}
	err = t.db.Where(payment.ExtraInfo{PaymentID: p.ID}).First(updatedExtraInfo).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "RECHECK", updatedPayment.AdminStatus.String)
	assert.Equal(t.T(), true, updatedPayment.FeeAmount.Valid)
	assert.Equal(t.T(), "10.00000000", updatedPayment.FeeAmount.String)
	assert.Equal(t.T(), "0.001", updatedExtraInfo.NetworkFee.String)
	assert.Equal(t.T(), int64(t.adminUserActor.ID), updatedExtraInfo.LastHandledID.Int64)
}

func (t *PaymentTests) TestUpdateWithdraw_Status_InProgress_And_AutoTransfer() {
	user := getNewUserActor()
	usdtUb := &userbalance.UserBalance{
		UserID:        user.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "1000.00",
		FrozenAmount:  "500.00",
		Address:       sql.NullString{String: "user1USDTAddress", Valid: true},
	}
	err := t.db.Create(usdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	p := &payment.Payment{
		UserID:            user.ID,
		CoinID:            1,
		Type:              "WITHDRAW",
		Status:            "CREATED",
		Code:              "USDT",
		FromAddress:       sql.NullString{String: "fromAddress", Valid: true},
		ToAddress:         sql.NullString{String: "user1USDTAddress", Valid: true},
		TxID:              sql.NullString{String: "txIdFromCallback", Valid: true},
		Amount:            sql.NullString{String: "100.00000000", Valid: true},
		BlockchainNetwork: sql.NullString{String: "ETH", Valid: true},
		FeeAmount:         sql.NullString{String: "5.00000000", Valid: true},
	}

	err = t.db.Create(p).Error
	if err != nil {
		t.Fail(err.Error())
	}

	ei := &payment.ExtraInfo{
		PaymentID: p.ID,
	}

	err = t.db.Create(ei).Error
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()
	data := fmt.Sprintf(`{
		"id" : %d,
		"status" : "in_progress",
		"auto_transfer" : true,
		"fee" : "10",
		"network_fee" : "0.001"
	}`, p.ID)
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payment/update-withdraw", bytes.NewReader(body))
	token := "Bearer " + t.adminUserActor.Token
	req.Header.Set("Authorization", token)
	t.adminHTTPServer.ServeHTTP(res, req)

	assert.Equal(t.T(), http.StatusOK, res.Code)

	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.T().Error(err)
	}
	updatedPayment := &payment.Payment{}
	err = t.db.Where(payment.Payment{ID: p.ID}).First(updatedPayment).Error
	if err != nil {
		t.Fail(err.Error())
	}
	updatedExtraInfo := &payment.ExtraInfo{}
	err = t.db.Where(payment.ExtraInfo{PaymentID: p.ID}).First(updatedExtraInfo).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "IN_PROGRESS", updatedPayment.Status)
	assert.Equal(t.T(), true, updatedPayment.TxID.Valid)
	assert.Equal(t.T(), "USDTTxId", updatedPayment.TxID.String)
	assert.Equal(t.T(), true, updatedPayment.FeeAmount.Valid)
	assert.Equal(t.T(), "10.00000000", updatedPayment.FeeAmount.String)
	assert.Equal(t.T(), "0.001", updatedExtraInfo.NetworkFee.String)
	assert.Equal(t.T(), int64(t.adminUserActor.ID), updatedExtraInfo.LastHandledID.Int64)
	assert.Equal(t.T(), true, updatedExtraInfo.AutoTransfer.Valid)
	assert.Equal(t.T(), true, updatedExtraInfo.AutoTransfer.Bool)

	//just to be sure that no balance and frozen get updated
	notUpdatedUserBalance := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(notUpdatedUserBalance).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), usdtUb.Amount, notUpdatedUserBalance.Amount)
	assert.Equal(t.T(), usdtUb.FrozenAmount, notUpdatedUserBalance.FrozenAmount)
}

func (t *PaymentTests) TestUpdateWithdraw_Status_InProgress_And_NotAutoTransfer() {
	user := getNewUserActor()
	usdtUb := &userbalance.UserBalance{
		UserID:        user.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "1000.00",
		FrozenAmount:  "500.00",
		Address:       sql.NullString{String: "user1USDTAddress", Valid: true},
	}
	err := t.db.Create(usdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	p := &payment.Payment{
		UserID:            user.ID,
		CoinID:            1,
		Type:              "WITHDRAW",
		Status:            "CREATED",
		Code:              "USDT",
		FromAddress:       sql.NullString{String: "fromAddress", Valid: true},
		ToAddress:         sql.NullString{String: "user1USDTAddress", Valid: true},
		TxID:              sql.NullString{String: "txIdFromCallback", Valid: true},
		Amount:            sql.NullString{String: "100.00000000", Valid: true},
		BlockchainNetwork: sql.NullString{String: "ETH", Valid: true},
		FeeAmount:         sql.NullString{String: "5.00000000", Valid: true},
	}

	err = t.db.Create(p).Error
	if err != nil {
		t.Fail(err.Error())
	}

	ei := &payment.ExtraInfo{
		PaymentID: p.ID,
	}

	err = t.db.Create(ei).Error
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()
	data := fmt.Sprintf(`{
		"id" : %d,
		"status" : "in_progress",
		"auto_transfer" : false,
		"fee" : "10",
		"network_fee" : "0.001"
	}`, p.ID)
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payment/update-withdraw", bytes.NewReader(body))
	token := "Bearer " + t.adminUserActor.Token
	req.Header.Set("Authorization", token)
	t.adminHTTPServer.ServeHTTP(res, req)

	assert.Equal(t.T(), http.StatusOK, res.Code)

	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.T().Error(err)
	}
	updatedPayment := &payment.Payment{}
	err = t.db.Where(payment.Payment{ID: p.ID}).First(updatedPayment).Error
	if err != nil {
		t.Fail(err.Error())
	}
	updatedExtraInfo := &payment.ExtraInfo{}
	err = t.db.Where(payment.ExtraInfo{PaymentID: p.ID}).First(updatedExtraInfo).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "COMPLETED", updatedPayment.Status)
	assert.Equal(t.T(), true, updatedPayment.FeeAmount.Valid)
	assert.Equal(t.T(), "10.00000000", updatedPayment.FeeAmount.String)
	assert.Equal(t.T(), "0.001", updatedExtraInfo.NetworkFee.String)
	assert.Equal(t.T(), int64(t.adminUserActor.ID), updatedExtraInfo.LastHandledID.Int64)
	assert.Equal(t.T(), true, updatedExtraInfo.AutoTransfer.Valid)
	assert.Equal(t.T(), false, updatedExtraInfo.AutoTransfer.Bool)

	//just to be sure that no balance and frozen get updated
	updatedUserBalance := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUserBalance).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "900.00000000", updatedUserBalance.Amount)
	assert.Equal(t.T(), "400.00000000", updatedUserBalance.FrozenAmount)

	//checking transactions
	tx := &transaction.Transaction{}
	err = t.db.Where(transaction.Transaction{Type: "WITHDRAW", PaymentID: sql.NullInt64{Int64: p.ID, Valid: true}}).First(tx).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "100.00000000", tx.Amount.String)
	assert.Equal(t.T(), int64(1), tx.CoinID)
	assert.Equal(t.T(), "USDT", tx.CoinName)
	assert.Equal(t.T(), "WITHDRAW", tx.Type)
	assert.Equal(t.T(), user.ID, tx.UserID)

	feeTx := &transaction.Transaction{}
	err = t.db.Where(transaction.Transaction{Type: "WITHDRAW_FEE", PaymentID: sql.NullInt64{Int64: p.ID, Valid: true}}).First(feeTx).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "10.00000000", feeTx.Amount.String)
	assert.Equal(t.T(), int64(1), feeTx.CoinID)
	assert.Equal(t.T(), "USDT", feeTx.CoinName)
	assert.Equal(t.T(), "WITHDRAW_FEE", feeTx.Type)
	assert.Equal(t.T(), user.ID, feeTx.UserID)
}

func (t *PaymentTests) TestUpdateWithdraw_Status_Rejected() {
	user := getNewUserActor()
	usdtUb := &userbalance.UserBalance{
		UserID:        user.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "1000.00000000",
		FrozenAmount:  "500.00000000",
		Address:       sql.NullString{String: "user1USDTAddress", Valid: true},
	}
	err := t.db.Create(usdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	p := &payment.Payment{
		UserID:            user.ID,
		CoinID:            1,
		Type:              "WITHDRAW",
		Status:            "CREATED",
		Code:              "USDT",
		FromAddress:       sql.NullString{String: "fromAddress", Valid: true},
		ToAddress:         sql.NullString{String: "user1USDTAddress", Valid: true},
		TxID:              sql.NullString{String: "txIdFromCallback", Valid: true},
		Amount:            sql.NullString{String: "100.00000000", Valid: true},
		BlockchainNetwork: sql.NullString{String: "ETH", Valid: true},
		FeeAmount:         sql.NullString{String: "5.00000000", Valid: true},
	}

	err = t.db.Create(p).Error
	if err != nil {
		t.Fail(err.Error())
	}

	ei := &payment.ExtraInfo{
		PaymentID: p.ID,
	}

	err = t.db.Create(ei).Error
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()
	data := fmt.Sprintf(`{
		"id" : %d,
		"status" : "rejected",
		"rejection_reason" : "rejected"
	}`, p.ID)
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payment/update-withdraw", bytes.NewReader(body))
	token := "Bearer " + t.adminUserActor.Token
	req.Header.Set("Authorization", token)
	t.adminHTTPServer.ServeHTTP(res, req)

	assert.Equal(t.T(), http.StatusOK, res.Code)

	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.T().Error(err)
	}
	updatedPayment := &payment.Payment{}
	err = t.db.Where(payment.Payment{ID: p.ID}).First(updatedPayment).Error
	if err != nil {
		t.Fail(err.Error())
	}
	updatedExtraInfo := &payment.ExtraInfo{}
	err = t.db.Where(payment.ExtraInfo{PaymentID: p.ID}).First(updatedExtraInfo).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "REJECTED", updatedPayment.Status)
	assert.Equal(t.T(), int64(t.adminUserActor.ID), updatedExtraInfo.LastHandledID.Int64)
	assert.Equal(t.T(), "rejected", updatedExtraInfo.RejectionReason.String)
	assert.Equal(t.T(), int64(t.adminUserActor.ID), updatedExtraInfo.LastHandledID.Int64)

	updatedUserBalance := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUserBalance).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "1000.00000000", updatedUserBalance.Amount)
	assert.Equal(t.T(), "400.00000000", updatedUserBalance.FrozenAmount)
}

func (t *PaymentTests) TestUpdateDeposit_ShouldNotDeposit() {
	user := getNewUserActor()
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
		FeeAmount:         sql.NullString{String: "5.00000000", Valid: true},
	}

	err := t.db.Create(p).Error
	if err != nil {
		t.Fail(err.Error())
	}

	ei := &payment.ExtraInfo{
		PaymentID: p.ID,
	}

	err = t.db.Create(ei).Error
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()
	data := fmt.Sprintf(`{
		"id" : %d,
		"amount" : "0.2",
		"from_address" : "fromAddressUpdated",
		"to_address" : "toAddressUpdated",
		"tx_id" : "txIdUpdated"
	}`, p.ID)
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payment/update-deposit", bytes.NewReader(body))
	token := "Bearer " + t.adminUserActor.Token
	req.Header.Set("Authorization", token)
	t.adminHTTPServer.ServeHTTP(res, req)

	assert.Equal(t.T(), http.StatusOK, res.Code)

	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.T().Error(err)
	}
	updatedPayment := &payment.Payment{}
	err = t.db.Where(payment.Payment{ID: p.ID}).First(updatedPayment).Error
	if err != nil {
		t.Fail(err.Error())
	}
	updatedExtraInfo := &payment.ExtraInfo{}
	err = t.db.Where(payment.ExtraInfo{PaymentID: p.ID}).First(updatedExtraInfo).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "CREATED", updatedPayment.Status)
	assert.Equal(t.T(), "txIdUpdated", updatedPayment.TxID.String)
	assert.Equal(t.T(), "fromAddressUpdated", updatedPayment.FromAddress.String)
	assert.Equal(t.T(), "toAddressUpdated", updatedPayment.ToAddress.String)
	assert.Equal(t.T(), "0.20000000", updatedPayment.Amount.String)

	assert.Equal(t.T(), int64(t.adminUserActor.ID), updatedExtraInfo.LastHandledID.Int64)
}

func (t *PaymentTests) TestUpdateDeposit_ShouldDeposit() {
	user := getNewUserActor()
	usdtUb := &userbalance.UserBalance{
		UserID:        user.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "1000.00000000",
		FrozenAmount:  "500.00000000",
		Address:       sql.NullString{String: "user1USDTAddress", Valid: true},
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
		FeeAmount:         sql.NullString{String: "5.00000000", Valid: true},
	}

	err = t.db.Create(p).Error
	if err != nil {
		t.Fail(err.Error())
	}

	ei := &payment.ExtraInfo{
		PaymentID: p.ID,
	}

	err = t.db.Create(ei).Error
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()
	data := fmt.Sprintf(`{
		"id" : %d,
		"amount" : "110",
		"from_address" : "fromAddressUpdated",
		"to_address" : "toAddressUpdated",
		"tx_id" : "txIdUpdated",
		"should_deposit" : true
	}`, p.ID)
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payment/update-deposit", bytes.NewReader(body))
	token := "Bearer " + t.adminUserActor.Token
	req.Header.Set("Authorization", token)
	t.adminHTTPServer.ServeHTTP(res, req)

	assert.Equal(t.T(), http.StatusOK, res.Code)

	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.T().Error(err)
	}
	updatedPayment := &payment.Payment{}
	err = t.db.Where(payment.Payment{ID: p.ID}).First(updatedPayment).Error
	if err != nil {
		t.Fail(err.Error())
	}
	updatedExtraInfo := &payment.ExtraInfo{}
	err = t.db.Where(payment.ExtraInfo{PaymentID: p.ID}).First(updatedExtraInfo).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "COMPLETED", updatedPayment.Status)
	assert.Equal(t.T(), "txIdUpdated", updatedPayment.TxID.String)
	assert.Equal(t.T(), "fromAddressUpdated", updatedPayment.FromAddress.String)
	assert.Equal(t.T(), "toAddressUpdated", updatedPayment.ToAddress.String)
	assert.Equal(t.T(), "110.00000000", updatedPayment.Amount.String)

	assert.Equal(t.T(), int64(t.adminUserActor.ID), updatedExtraInfo.LastHandledID.Int64)

	updatedUserBalance := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUserBalance).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "1110.00000000", updatedUserBalance.Amount)
	assert.Equal(t.T(), "500.00000000", updatedUserBalance.FrozenAmount)

	//checking transactions
	tx := &transaction.Transaction{}
	err = t.db.Where(transaction.Transaction{Type: "DEPOSIT", PaymentID: sql.NullInt64{Int64: p.ID, Valid: true}}).First(tx).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "110.00000000", tx.Amount.String)
	assert.Equal(t.T(), int64(1), tx.CoinID)
	assert.Equal(t.T(), "USDT", tx.CoinName)
	assert.Equal(t.T(), "DEPOSIT", tx.Type)
	assert.Equal(t.T(), user.ID, tx.UserID)
}

func TestPayment(t *testing.T) {
	suite.Run(t, &PaymentTests{
		Suite: new(suite.Suite),
	})
}
