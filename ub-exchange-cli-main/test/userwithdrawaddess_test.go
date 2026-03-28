package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"exchange-go/internal/api"
	"exchange-go/internal/di"
	"exchange-go/internal/response"
	"exchange-go/internal/userwithdrawaddress"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserWithdrawAddressTests struct {
	*suite.Suite
	httpServer http.Handler
	db         *gorm.DB
	userActor  *userActor
}

func (t *UserWithdrawAddressTests) SetupSuite() {
	container := getContainer()
	t.httpServer = container.Get(di.HTTPServer).(api.HTTPServer).GetEngine()
	t.db = getDb()
	t.userActor = getUserActor()

}

func (t *UserWithdrawAddressTests) SetupTest() {
	t.db.Where("id > ?", 0).Delete(userwithdrawaddress.UserWithdrawAddress{})

}

func (t *UserWithdrawAddressTests) TearDownTest() {

}

func (t *UserWithdrawAddressTests) TearDownSuite() {

}

func (t *UserWithdrawAddressTests) TestGetWithdrawAddresses() {

	uwa1 := &userwithdrawaddress.UserWithdrawAddress{
		ID:      1,
		UserID:  t.userActor.ID,
		CoinID:  1,
		Address: "usdtAddress1",
		Label:   sql.NullString{String: "usdt1", Valid: true},
		Network: sql.NullString{String: "ETH", Valid: true},
	}

	uwa2 := &userwithdrawaddress.UserWithdrawAddress{
		ID:      2,
		UserID:  t.userActor.ID,
		CoinID:  1,
		Address: "usdtAddress2",
		Label:   sql.NullString{String: "usdt2", Valid: true},
		Network: sql.NullString{String: "TRX", Valid: true},
	}
	uwa3 := &userwithdrawaddress.UserWithdrawAddress{
		ID:      3,
		UserID:  t.userActor.ID,
		CoinID:  2,
		Address: "btcAddress1",
		Label:   sql.NullString{String: "btc1", Valid: true},
	}

	uwa4 := &userwithdrawaddress.UserWithdrawAddress{
		ID:      4,
		UserID:  t.userActor.ID,
		CoinID:  3,
		Address: "ethAddress1",
		Label:   sql.NullString{String: "eth1", Valid: true},
	}

	userWithdrawAddresses := []*userwithdrawaddress.UserWithdrawAddress{uwa1, uwa2, uwa3, uwa4}
	err := t.db.Create(userWithdrawAddresses).Error
	if err != nil {
		t.Fail(err.Error())
	}

	queryParams := url.Values{}
	queryParams.Set("page_size", "2")
	paramsString := queryParams.Encode()

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/withdraw-address?"+paramsString, nil)
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    []userwithdrawaddress.GetWithdrawAddressesResponse
	}{}
	//result := response.ApiResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), 2, len(result.Data))

	for _, uwa := range result.Data {
		switch uwa.ID {
		case 3:
			assert.Equal(t.T(), "BTC", uwa.Coin)
			assert.Equal(t.T(), "btcAddress1", uwa.Address)
			assert.Equal(t.T(), "btc1", uwa.Label)
			assert.Equal(t.T(), "Bitcoin", uwa.Name)
			assert.Equal(t.T(), "", uwa.Network)
		case 4:
			assert.Equal(t.T(), "ETH", uwa.Coin)
			assert.Equal(t.T(), "ethAddress1", uwa.Address)
			assert.Equal(t.T(), "eth1", uwa.Label)
			assert.Equal(t.T(), "Ethereum", uwa.Name)
			assert.Equal(t.T(), "", uwa.Network)

		default:
			t.Fail("we should not be in default case")
		}
	}

	//testing pagination
	queryParams = url.Values{}
	queryParams.Set("page", "1")
	queryParams.Set("page_size", "2")
	paramsString = queryParams.Encode()

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/withdraw-address?"+paramsString, nil)
	token = "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result = struct {
		Status  bool
		Message string
		Data    []userwithdrawaddress.GetWithdrawAddressesResponse
	}{}
	//result := response.ApiResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), 2, len(result.Data))

	for _, uwa := range result.Data {
		switch uwa.ID {
		case 1:
			assert.Equal(t.T(), "USDT", uwa.Coin)
			assert.Equal(t.T(), "usdtAddress1", uwa.Address)
			assert.Equal(t.T(), "usdt1", uwa.Label)
			assert.Equal(t.T(), "Tether", uwa.Name)
			assert.Equal(t.T(), "ETH", uwa.Network)
		case 2:
			assert.Equal(t.T(), "USDT", uwa.Coin)
			assert.Equal(t.T(), "usdtAddress2", uwa.Address)
			assert.Equal(t.T(), "usdt2", uwa.Label)
			assert.Equal(t.T(), "Tether", uwa.Name)
			assert.Equal(t.T(), "TRX", uwa.Network)
		default:
			t.Fail("we should not be in default case")
		}
	}

	//testing filters
	queryParams = url.Values{}
	queryParams.Set("page", "0")
	queryParams.Set("page_size", "2")
	queryParams.Set("code", "USDT")
	queryParams.Set("label", "usdt1")
	queryParams.Set("address", "usdtAddress1")
	paramsString = queryParams.Encode()

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/withdraw-address?"+paramsString, nil)
	token = "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result = struct {
		Status  bool
		Message string
		Data    []userwithdrawaddress.GetWithdrawAddressesResponse
	}{}
	//result := response.ApiResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), 1, len(result.Data))

	for _, uwa := range result.Data {
		switch uwa.ID {
		case 1:
			assert.Equal(t.T(), "USDT", uwa.Coin)
			assert.Equal(t.T(), "usdtAddress1", uwa.Address)
			assert.Equal(t.T(), "usdt1", uwa.Label)
			assert.Equal(t.T(), "Tether", uwa.Name)
			assert.Equal(t.T(), "ETH", uwa.Network)
		default:
			t.Fail("we should not be in default case")
		}
	}
}

func (t *UserWithdrawAddressTests) TestGetFormerAddresses_Successful() {
	uwa1 := &userwithdrawaddress.UserWithdrawAddress{
		ID:      1,
		UserID:  t.userActor.ID,
		CoinID:  1,
		Address: "usdtAddress1",
		Label:   sql.NullString{String: "usdt1", Valid: true},
		Network: sql.NullString{String: "ETH", Valid: true},
	}

	uwa2 := &userwithdrawaddress.UserWithdrawAddress{
		ID:      2,
		UserID:  t.userActor.ID,
		CoinID:  1,
		Address: "usdtAddress2",
		Label:   sql.NullString{String: "usdt2", Valid: true},
		Network: sql.NullString{String: "TRX", Valid: true},
	}
	uwa3 := &userwithdrawaddress.UserWithdrawAddress{
		ID:      3,
		UserID:  t.userActor.ID,
		CoinID:  2,
		Address: "btcAddress1",
		Label:   sql.NullString{String: "btc1", Valid: true},
	}

	uwa4 := &userwithdrawaddress.UserWithdrawAddress{
		ID:      4,
		UserID:  t.userActor.ID,
		CoinID:  3,
		Address: "ethAddress1",
		Label:   sql.NullString{String: "eth1", Valid: true},
	}

	userWithdrawAddresses := []*userwithdrawaddress.UserWithdrawAddress{uwa1, uwa2, uwa3, uwa4}
	err := t.db.Create(userWithdrawAddresses).Error
	if err != nil {
		t.Fail(err.Error())
	}

	queryParams := url.Values{}
	queryParams.Set("code", "USDT")
	paramsString := queryParams.Encode()

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/withdraw-address/former-addresses?"+paramsString, nil)
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    []userwithdrawaddress.GetWithdrawAddressesResponse
	}{}
	//result := response.ApiResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), 2, len(result.Data))

	for _, uwa := range result.Data {
		switch uwa.ID {
		case 1:
			assert.Equal(t.T(), "USDT", uwa.Coin)
			assert.Equal(t.T(), "usdtAddress1", uwa.Address)
			assert.Equal(t.T(), "usdt1", uwa.Label)
			assert.Equal(t.T(), "Tether", uwa.Name)
			assert.Equal(t.T(), "ETH", uwa.Network)
		case 2:
			assert.Equal(t.T(), "USDT", uwa.Coin)
			assert.Equal(t.T(), "usdtAddress2", uwa.Address)
			assert.Equal(t.T(), "usdt2", uwa.Label)
			assert.Equal(t.T(), "Tether", uwa.Name)
			assert.Equal(t.T(), "TRX", uwa.Network)
		default:
			t.Fail("we should not be in default case")
		}
	}

}

func (t *UserWithdrawAddressTests) TestGetFormerAddresses_Fail_CoinNotFound() {
	uwa1 := &userwithdrawaddress.UserWithdrawAddress{
		ID:      1,
		UserID:  t.userActor.ID,
		CoinID:  1,
		Address: "usdtAddress1",
		Label:   sql.NullString{String: "usdt1", Valid: true},
		Network: sql.NullString{String: "ETH", Valid: true},
	}

	uwa2 := &userwithdrawaddress.UserWithdrawAddress{
		ID:      2,
		UserID:  t.userActor.ID,
		CoinID:  1,
		Address: "usdtAddress2",
		Label:   sql.NullString{String: "usdt2", Valid: true},
		Network: sql.NullString{String: "TRX", Valid: true},
	}
	uwa3 := &userwithdrawaddress.UserWithdrawAddress{
		ID:      3,
		UserID:  t.userActor.ID,
		CoinID:  2,
		Address: "btcAddress1",
		Label:   sql.NullString{String: "btc1", Valid: true},
	}

	uwa4 := &userwithdrawaddress.UserWithdrawAddress{
		ID:      4,
		UserID:  t.userActor.ID,
		CoinID:  3,
		Address: "ethAddress1",
		Label:   sql.NullString{String: "eth1", Valid: true},
	}

	userWithdrawAddresses := []*userwithdrawaddress.UserWithdrawAddress{uwa1, uwa2, uwa3, uwa4}
	err := t.db.Create(userWithdrawAddresses).Error
	if err != nil {
		t.Fail(err.Error())
	}

	queryParams := url.Values{}
	queryParams.Set("code", "btcq") // wrong code
	paramsString := queryParams.Encode()
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/withdraw-address/former-addresses?"+paramsString, nil)
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "coin not found", result.Message)

}

type newWithdrawAddressValidationFailedScenarios struct {
	data         string
	reason       string
	errorMessage string
}

func (t *OrderCreateTests) TestNewWithdrawAddress_ValidationFail() {
	//insert usersPermissions
	failedScenarios := []newWithdrawAddressValidationFailedScenarios{
		{
			data:         `{"code":"","label":"btc1","address":"btcAddress1","network":""}`,
			reason:       "code not provided",
			errorMessage: "code is required",
		},
		{
			data:         `{"code":"btc","label":"","address":"btcAddress1","network":""}`,
			reason:       "label not provided",
			errorMessage: "label is required",
		},
		{
			data:         `{"code":"btc","label":"btc1","address":"","network":""}`,
			reason:       "address not provided",
			errorMessage: "address is required",
		},
		{
			data:         `{"code":"btcq","label":"btc1","address":"btcAddress1","network":""}`,
			reason:       "wrong code",
			errorMessage: "coin not found",
		},
		{
			data:         `{"code":"usdt","label":"usdt11","address":"usdtAddress1","network":"ethq"}`,
			reason:       "wrong network code",
			errorMessage: "network not found",
		},
	}

	for _, item := range failedScenarios {
		res := httptest.NewRecorder()
		body := []byte(item.data)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/withdraw-address/new", bytes.NewReader(body))
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

func (t *UserWithdrawAddressTests) Test_NewWithDrawAddress_Successful() {
	data := `{"code":"usdt","label":"usdt1","address":"usdtAddress1","network":"ETH"}`
	body := []byte(data)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/withdraw-address/new", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token

	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    []userwithdrawaddress.CreateAddressResponse
	}{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)

	assert.Equal(t.T(), "USDT", result.Data[0].Coin)
	assert.Equal(t.T(), "Tether", result.Data[0].Name)
	assert.Equal(t.T(), "usdtAddress1", result.Data[0].Address)
	assert.Equal(t.T(), "usdt1", result.Data[0].Label)
	assert.Equal(t.T(), false, result.Data[0].IsFavorite)
	assert.Equal(t.T(), "ETH", result.Data[0].Network)
}

func (t *UserWithdrawAddressTests) Test_AddToFavorite() {
	uwa := &userwithdrawaddress.UserWithdrawAddress{
		ID:      1,
		UserID:  t.userActor.ID,
		CoinID:  1,
		Address: "usdtAddress1",
		Label:   sql.NullString{String: "eth1", Valid: true},
	}

	err := t.db.Create(uwa).Error
	if err != nil {
		t.Fail(err.Error())
	}

	data := `{"id":1,"action":"add"}`
	body := []byte(data)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/withdraw-address/favorite", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token

	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)
	//get from db and check
	uwaUpdated := &userwithdrawaddress.UserWithdrawAddress{}
	err = t.db.Where("id = ?", int64(1)).First(uwaUpdated).Error
	assert.Equal(t.T(), true, uwaUpdated.IsFavorite.Bool)

	//remove from favorites
	data = `{"id":1,"action":"remove"}`
	body = []byte(data)
	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/v1/withdraw-address/favorite", bytes.NewReader(body))
	token = "Bearer " + t.userActor.Token

	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result = response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)
	//get from db and check
	uwaUpdated = &userwithdrawaddress.UserWithdrawAddress{}
	err = t.db.Where("id = ?", int64(1)).First(uwaUpdated).Error
	assert.Equal(t.T(), false, uwaUpdated.IsFavorite.Bool)

	//scenario when id not found
	data = `{"id":2,"action":"add"}`
	body = []byte(data)
	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/v1/withdraw-address/favorite", bytes.NewReader(body))
	token = "Bearer " + t.userActor.Token

	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result = response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), false, result.Status)
	assert.Equal(t.T(), "address not found", result.Message)

}

func (t *UserWithdrawAddressTests) Test_Delete() {
	uwa := &userwithdrawaddress.UserWithdrawAddress{
		ID:      1,
		UserID:  t.userActor.ID,
		CoinID:  1,
		Address: "usdtAddress1",
		Label:   sql.NullString{String: "eth1", Valid: true},
	}

	err := t.db.Create(uwa).Error
	if err != nil {
		t.Fail(err.Error())
	}

	data := `{"ids":[1,2,3]}`
	body := []byte(data)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/withdraw-address/delete", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token

	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)
	//get from db and check
	uwaUpdated := &userwithdrawaddress.UserWithdrawAddress{}
	err = t.db.Where("id = ?", int64(1)).First(uwaUpdated).Error
	assert.Equal(t.T(), true, uwaUpdated.IsDeleted.Bool)
}

func TestUserWithdrawAddress(t *testing.T) {
	suite.Run(t, &UserWithdrawAddressTests{
		Suite: new(suite.Suite),
	})
}
