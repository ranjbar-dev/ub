// Package order_test tests the UnmatchedOrdersHandler. Covers:
//   - Match: pops an unmatched order ID from Redis, loads the order from the repository,
//     and resubmits it to the engine communicator if the order is still OPEN
//
// Test data: mocked RedisClient with RPop from unmatched orders list, mocked order
// repository returning an OPEN order, mocked engine communicator, and test environment configs.
package order_test

import (
	"exchange-go/internal/mocks"
	"exchange-go/internal/order"
	"exchange-go/internal/platform"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestUnmatchedOrdershandler_Match(t *testing.T) {
	rc := new(mocks.RedisClient)
	rc.On("RPop", mock.Anything, order.UnmatchedOrdersList).Once().Return("1", nil)

	var o *order.Order
	orderRepository := new(mocks.OrderRepository)

	orderRepository.On("GetOrderByID", int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		o = args.Get(1).(*order.Order)
		o.ID = 1
		o.Status = order.StatusOpen
	})

	engineCommunicator := new(mocks.EngineCommunicator)
	engineCommunicator.On("SubmitOrder", mock.Anything).Once().Return(nil)
	configs := new(mocks.Configs)
	configs.On("GetEnv").Twice().Return(platform.EnvTest)
	logger := new(mocks.Logger)
	unmatchedOrdersHandler := order.NewUnmatchedOrdersHandler(rc, orderRepository, engineCommunicator, configs, logger)
	unmatchedOrdersHandler.Match()

	rc.AssertExpectations(t)
	orderRepository.AssertExpectations(t)
	engineCommunicator.AssertExpectations(t)
}
