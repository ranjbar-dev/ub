// Package wallet_test tests the wallet service for blockchain wallet operations. Covers:
//   - Validating cryptocurrency addresses via external wallet API
//   - Generating deposit addresses for users by coin type
//   - Sending blockchain transactions with coin, amount, destination, and fee parameters
//   - Querying address balances from the wallet backend
//
// Test data: mock wallet authorization service, HTTP client, and config provider
// with JSON response fixtures for address validation, generation, transaction, and balance endpoints.
package wallet_test

import (
	"exchange-go/internal/mocks"
	"exchange-go/internal/wallet"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_IsAddressValid(t *testing.T) {
	as := new(mocks.WalletAuthorizationService)
	as.On("GetToken", mock.Anything).Once().Return("token", nil)
	httpClient := new(mocks.HttpClient)
	body := []byte("{" +
		"\"status\": true," +
		"\"message\": \"\"," +
		"\"data\": {\"isValid\":true}" +
		"}")
	httpClient.On("HTTPPost", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(body, http.Header{}, http.StatusOK, nil)
	configs := new(mocks.Configs)
	configs.On("GetString", mock.Anything).Times(3).Return("")
	configs.On("GetEnv").Once().Return("test")
	logger := new(mocks.Logger)
	s := wallet.NewWalletService(as, httpClient, configs, logger)
	coin := "BTC"
	address := "address"
	network := "network"
	isValid, err := s.IsAddressValid(coin, address, network)
	assert.Nil(t, err)
	assert.True(t, isValid)
	//as.AssertExpectations(t)
	//httpClient.AssertExpectations(t)
}

func TestService_GetAddressForUser(t *testing.T) {
	as := new(mocks.WalletAuthorizationService)
	as.On("GetToken", mock.Anything).Once().Return("token", nil)
	httpClient := new(mocks.HttpClient)
	body := []byte("{" +
		"\"status\": true," +
		"\"message\": \"\"," +
		"\"data\": {\"address\":\"address\"}" +
		"}")
	httpClient.On("HttpPost", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(body, http.Header{}, http.StatusOK, nil)
	configs := new(mocks.Configs)
	configs.On("GetString", mock.Anything).Times(3).Return("")
	configs.On("GetEnv").Once().Return("test")
	logger := new(mocks.Logger)
	s := wallet.NewWalletService(as, httpClient, configs, logger)
	coin := "BTC"

	userCode := "someUniqueId"
	address, err := s.GetAddressForUser(coin, userCode)
	assert.Nil(t, err)
	assert.Equal(t, "BTCAddress", address)
	//as.AssertExpectations(t)
	//httpClient.AssertExpectations(t)
}

func TestService_SendTransaction(t *testing.T) {
	as := new(mocks.WalletAuthorizationService)
	as.On("GetToken", mock.Anything).Once().Return("token", nil)
	httpClient := new(mocks.HttpClient)
	body := []byte("{" +
		"\"status\": true," +
		"\"message\": \"\"," +
		"\"data\": {\"txId\":\"txId\"}" +
		"}")
	httpClient.On("HttpPost", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(body, http.Header{}, http.StatusOK, nil)
	configs := new(mocks.Configs)
	configs.On("GetString", mock.Anything).Times(3).Return("")
	configs.On("GetEnv").Once().Return("test")
	logger := new(mocks.Logger)
	s := wallet.NewWalletService(as, httpClient, configs, logger)

	tx, err := s.SendTransaction("BTC", "1.0", "toAddress", "ETH", "0.001")
	assert.Nil(t, err)
	assert.Equal(t, "BTCTxId", tx)
	//as.AssertExpectations(t)
	//httpClient.AssertExpectations(t)
}

func TestService_GetAddressBalance(t *testing.T) {
	as := new(mocks.WalletAuthorizationService)
	as.On("GetToken", mock.Anything).Once().Return("token", nil)
	httpClient := new(mocks.HttpClient)
	body := []byte("{" +
		"\"status\": true," +
		"\"message\": \"\"," +
		"\"data\": {\"balance\":\"0.1\"}" +
		"}")
	httpClient.On("HTTPPost", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(body, http.Header{}, http.StatusOK, nil)
	configs := new(mocks.Configs)
	configs.On("GetString", mock.Anything).Times(3).Return("")
	configs.On("GetEnv").Once().Return("prod")
	logger := new(mocks.Logger)
	s := wallet.NewWalletService(as, httpClient, configs, logger)

	balance, err := s.GetAddressBalance("BTC", "BTC", "btcAddress", true)
	assert.Nil(t, err)
	assert.Equal(t, "0.1", balance)

}
