// Package order_test tests the InQueueOrderManager. Covers:
//   - HandleInQueueOrders: delegates in-queue order processing to the engine for a given pair and price
//
// Test data: mocked Engine with HandleInQueueOrders expectation for BTC-USDT at price 50000.
package order_test

import (
	"context"
	"exchange-go/internal/mocks"
	"exchange-go/internal/order"
	"testing"
)

func TestInQueueOrderManager_HandleInQueueOrders(t *testing.T) {
	e := new(mocks.Engine)
	e.On("HandleInQueueOrders", "BTC-USDT", "50000").Once().Return(nil)
	logger := new(mocks.Logger)
	inQueueOrderManager := order.NewInQueueOrderManager(e, logger)
	ctx := context.Background()
	pairName := "BTC-USDT"
	price := "50000"
	inQueueOrderManager.HandleInQueueOrders(ctx, pairName, price)
	e.AssertExpectations(t)
}
