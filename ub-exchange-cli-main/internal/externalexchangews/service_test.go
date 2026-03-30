// Package externalexchangews_test tests the external exchange WebSocket service factory. Covers:
//   - Resolving the active external exchange WebSocket implementation by config name
//   - Constructing a Binance WebSocket instance with active trading pairs
//
// Test data: mock WebSocket client, processor, config provider returning "binance",
// and currency service providing a single active BTC-USDT pair.
package externalexchangews_test

import (
	"exchange-go/internal/currency"
	"exchange-go/internal/externalexchangews"
	"exchange-go/internal/mocks"
	"testing"
)

func TestService_GetActiveExternalExchangeWs(t *testing.T) {
	wsClient := new(mocks.WsClient)
	processor := new(mocks.Processor)
	logger := new(mocks.Logger)
	configs := new(mocks.Configs)
	configs.On("GetActiveExternalExchange").Once().Return("binance")
	activePairs := []currency.Pair{
		{
			ID:              1,
			Name:            "BTC-USDT",
			IsActive:        true,
			Spread:          1.5,
			ShowDigits:      6,
			BasisCoinID:     1,
			DependentCoinID: 2,
		},
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetActivePairCurrenciesList").Once().Return(activePairs)
	externalExchangeService := externalexchangews.NewService(wsClient, processor, logger, configs, currencyService)
	_, err := externalExchangeService.GetActiveExternalExchangeWs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	configs.AssertExpectations(t)
}
