// Package externalexchange_test tests the order-from-external service for ingesting
// external exchange data. Covers:
//   - Retrieving the last ingested order from an external exchange by pair ID
//   - Retrieving the last ingested trade from an external exchange by pair ID
//   - Creating new order records from external exchange data
//   - Creating new trade records from external exchange data
//
// Test data: mock order-from-external and trade-from-external repositories
// with ID auto-assignment via mock Run callbacks.
package externalexchange_test

import (
	"exchange-go/internal/externalexchange"
	"exchange-go/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrderFromExternalService_GetLastOrderFromExternalByPairID(t *testing.T) {
	orderFromExternalRepo := new(mocks.OrderFromExternalRepository)
	orderFromExternalRepo.On("GetLastOrderFromExternalByPairID", int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ofe := args.Get(1).(*externalexchange.OrderFromExternal)
		ofe.ID = 1
	})

	tradeFromExternalRepo := new(mocks.TradeFromExternalRepository)

	service := externalexchange.NewOrderFromExternalService(orderFromExternalRepo, tradeFromExternalRepo)

	o := &externalexchange.OrderFromExternal{}
	err := service.GetLastOrderFromExternalByPairID(int64(1), o)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), o.ID)
	orderFromExternalRepo.AssertExpectations(t)
}

func TestOrderFromExternalService_GetLastTradeFromExternalByPairID(t *testing.T) {
	orderFromExternalRepo := new(mocks.OrderFromExternalRepository)

	tradeFromExternalRepo := new(mocks.TradeFromExternalRepository)
	tradeFromExternalRepo.On("GetLastTradeFromExternalByPairID", int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		tfe := args.Get(1).(*externalexchange.TradeFromExternal)
		tfe.ID = 1
	})

	service := externalexchange.NewOrderFromExternalService(orderFromExternalRepo, tradeFromExternalRepo)

	tf := &externalexchange.TradeFromExternal{}
	err := service.GetLastTradeFromExternalByPairID(int64(1), tf)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), tf.ID)
	tradeFromExternalRepo.AssertExpectations(t)
}

func TestOrderFromExternalService_CreateOrder(t *testing.T) {
	orderFromExternalRepo := new(mocks.OrderFromExternalRepository)
	orderFromExternalRepo.On("Create", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ofe := args.Get(0).(*externalexchange.OrderFromExternal)
		ofe.ID = 1
	})

	tradeFromExternalRepo := new(mocks.TradeFromExternalRepository)

	service := externalexchange.NewOrderFromExternalService(orderFromExternalRepo, tradeFromExternalRepo)

	o := &externalexchange.OrderFromExternal{}
	err := service.CreateOrder(o)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), o.ID)

	orderFromExternalRepo.AssertExpectations(t)

}

func TestOrderFromExternalService_CreateTrade(t *testing.T) {
	orderFromExternalRepo := new(mocks.OrderFromExternalRepository)

	tradeFromExternalRepo := new(mocks.TradeFromExternalRepository)
	tradeFromExternalRepo.On("Create", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		tfe := args.Get(0).(*externalexchange.TradeFromExternal)
		tfe.ID = 1
	})

	service := externalexchange.NewOrderFromExternalService(orderFromExternalRepo, tradeFromExternalRepo)

	tf := &externalexchange.TradeFromExternal{}
	err := service.CreateTrade(tf)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), tf.ID)
	tradeFromExternalRepo.AssertExpectations(t)

}
