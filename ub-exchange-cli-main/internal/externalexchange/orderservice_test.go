// Package externalexchange_test tests the external exchange order service. Covers:
//   - Creating external exchange orders for bot trading (market orders with buy/sell amounts)
//   - Creating external exchange orders for user trading with order placement result verification
//
// Test data: mock external exchange order repository, external exchange service,
// and logger with BTC-USDT market order fixtures and order placement results.
package externalexchange_test

import (
	"exchange-go/internal/externalexchange"
	"exchange-go/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrderService_CreateExternalExchangeOrderForBot(t *testing.T) {
	externalExchangeOrderRepo := new(mocks.ExternalExchangeOrderRepository)
	externalExchangeOrderRepo.On("Create", mock.Anything).Once().Return(nil)
	externalExchangeOrderRepo.On("Update", mock.Anything).Once().Return(nil)
	externalExchangeService := new(mocks.ExternalExchangeService)
	data := externalexchange.ExternalOrderData{
		Pair:         "BTC-USDT",
		Type:         "BUY",
		Amount:       "0.10000000",
		Price:        "",
		ExchangeType: "MARKET",
	}
	orderPlacementResult := externalexchange.OrderPlacementResult{
		IsOrderPlaced:           true,
		ExternalExchangeOrderID: "1",
		ExternalExchangeID:      1,
		Data:                    "{}",
	}
	externalExchangeService.On("OrderPlacement", data).Once().Return(orderPlacementResult, nil)
	logger := new(mocks.Logger)
	service := externalexchange.NewOrderService(externalExchangeOrderRepo, externalExchangeService, logger)
	params := externalexchange.BotOrderParams{
		PairID:       1,
		PairName:     "BTC-USDT",
		Type:         "BUY",
		ExchangeType: "MARKET",
		Amount:       "0.10000000",
		Price:        "",
		BuyAmount:    "0.15000000",
		BuyPrice:     "",
		SellAmount:   "0.05000000",
		SellPrice:    "",
		LastTradeID:  1,
		OrderIds:     []string{"1"},
	}
	isOrderPlaced, err := service.CreateExternalExchangeOrderForBot(params)
	assert.Nil(t, err)
	assert.True(t, isOrderPlaced)

	externalExchangeOrderRepo.AssertExpectations(t)
	externalExchangeService.AssertExpectations(t)
}

func TestOrderService_CreateExternalExchangeOrderForUser(t *testing.T) {
	externalExchangeOrderRepo := new(mocks.ExternalExchangeOrderRepository)
	externalExchangeOrderRepo.On("Create", mock.Anything).Once().Return(nil)
	externalExchangeOrderRepo.On("Update", mock.Anything).Once().Return(nil)
	externalExchangeService := new(mocks.ExternalExchangeService)
	data := externalexchange.ExternalOrderData{
		Pair:         "BTC-USDT",
		Type:         "BUY",
		Amount:       "0.10000000",
		Price:        "",
		ExchangeType: "MARKET",
	}
	orderPlacementResult := externalexchange.OrderPlacementResult{
		IsOrderPlaced:           true,
		ExternalExchangeOrderID: "1",
		ExternalExchangeID:      1,
		Data:                    "{}",
	}
	externalExchangeService.On("OrderPlacement", data).Once().Return(orderPlacementResult, nil)
	logger := new(mocks.Logger)
	service := externalexchange.NewOrderService(externalExchangeOrderRepo, externalExchangeService, logger)
	params := externalexchange.UserOrderParams{
		PairID:       1,
		PairName:     "BTC-USDT",
		Type:         "BUY",
		ExchangeType: "MARKET",
		Amount:       "0.10000000",
		Price:        "",
	}
	result, err := service.CreateExternalExchangeOrderForUser(params)
	assert.Nil(t, err)
	assert.True(t, result.IsOrderPlaced)
	assert.Equal(t, int64(1), result.ExternalExchangeID)
	assert.Equal(t, "1", result.ExternalExchangeOrderID)

	externalExchangeOrderRepo.AssertExpectations(t)
	externalExchangeService.AssertExpectations(t)

}
