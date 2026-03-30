// Package order_test tests the OrderEventsHandler. Covers:
//   - HandleOrderCreation routing to our exchange via engine communicator
//   - HandleOrderCreation routing to external exchange with post-trade handling
//   - HandleOrderCreation for stop orders queued in Redis
//   - HandleOrderCreation for stop order submission to our exchange
//   - HandleOrderCreation for stop order submission to external exchange
//   - HandleOrderCreation fallback when external exchange returns placement failure
//   - HandleOrderCancellation publishing cancellation via Centrifugo
//   - HandleOrderFulfillByAdmin removing order from engine
//
// Test data: mocked Redis manager, decision manager, Centrifugo manager, external exchange
// service, engine communicator, post-order matching service, and logger.
package order_test

import (
	"database/sql"
	"exchange-go/internal/externalexchange"
	"exchange-go/internal/mocks"
	"exchange-go/internal/order"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

func TestEventsHandler_HandleOrderCreation_OurExchange(t *testing.T) {

	redisManager := new(mocks.OrderRedisManager)
	decisionManager := new(mocks.DecisionManager)
	decisionManager.On("DecideOrderPlacement", mock.Anything).Once().Return(order.PlaceOurExchange, nil)

	mqttManager := new(mocks.CentrifugoManager)
	mqttManager.On("PublishOrderToOpenOrders", mock.Anything, mock.Anything, mock.Anything).Once().Return()

	externalExchangeOrderService := new(mocks.ExternalExchangeOrderService)
	engineCommunicator := new(mocks.EngineCommunicator)
	engineCommunicator.On("SubmitOrder", mock.Anything).Once().Return(nil)
	postOrderMatchingService := new(mocks.PostOrderMatchingService)
	logger := new(mocks.Logger)

	eh := order.NewOrderEventsHandler(redisManager, decisionManager, mqttManager, externalExchangeOrderService, engineCommunicator,
		postOrderMatchingService, logger)

	o := order.Order{
		Price: sql.NullString{String: "50000.0", Valid: true},
	}
	eh.HandleOrderCreation(o, false)
	time.Sleep(50 * time.Millisecond)
	decisionManager.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
	engineCommunicator.AssertExpectations(t)

}

func TestEventsHandler_HandleOrderCreation_ExternalExchange(t *testing.T) {
	redisManager := new(mocks.OrderRedisManager)
	decisionManager := new(mocks.DecisionManager)
	decisionManager.On("DecideOrderPlacement", mock.Anything).Once().Return(order.PlaceExternalExchange, nil)
	mqttManager := new(mocks.CentrifugoManager)
	mqttManager.On("PublishOrderToOpenOrders", mock.Anything, mock.Anything, mock.Anything).Once().Return()

	externalExchangeOrderService := new(mocks.ExternalExchangeOrderService)
	orderResult := externalexchange.UserOrderResult{IsOrderPlaced: true}
	externalExchangeOrderService.On("CreateExternalExchangeOrderForUser", mock.Anything).Once().Return(orderResult, nil)

	engineCommunicator := new(mocks.EngineCommunicator)
	postOrderMatchingService := new(mocks.PostOrderMatchingService)
	postOrderMatchingService.On("HandleExternalTradedOrder", mock.Anything).Once().Return(nil)
	logger := new(mocks.Logger)

	eh := order.NewOrderEventsHandler(redisManager, decisionManager, mqttManager, externalExchangeOrderService, engineCommunicator,
		postOrderMatchingService, logger)

	o := order.Order{}
	eh.HandleOrderCreation(o, false)
	time.Sleep(50 * time.Millisecond)
	decisionManager.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
	externalExchangeOrderService.AssertExpectations(t)
	postOrderMatchingService.AssertExpectations(t)
}

func TestEventsHandler_HandleOrderCreation_ForStopOrder(t *testing.T) {
	redisManager := new(mocks.OrderRedisManager)
	redisManager.On("AddStopOrderToQueue", mock.Anything, mock.Anything).Once().Return(nil)
	decisionManager := new(mocks.DecisionManager)
	mqttManager := new(mocks.CentrifugoManager)
	mqttManager.On("PublishOrderToOpenOrders", mock.Anything, mock.Anything, mock.Anything).Once().Return()
	externalExchangeOrderService := new(mocks.ExternalExchangeOrderService)

	engineCommunicator := new(mocks.EngineCommunicator)
	postOrderMatchingService := new(mocks.PostOrderMatchingService)
	logger := new(mocks.Logger)

	eh := order.NewOrderEventsHandler(redisManager, decisionManager, mqttManager, externalExchangeOrderService, engineCommunicator,
		postOrderMatchingService, logger)

	o := order.Order{
		StopPointPrice: sql.NullString{String: "50000", Valid: true},
	}
	eh.HandleOrderCreation(o, false)
	time.Sleep(50 * time.Millisecond)
	redisManager.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
}

func TestEventsHandler_HandleOrderCreation_ForStopOrderSubmission_OurExchange(t *testing.T) {
	redisManager := new(mocks.OrderRedisManager)
	decisionManager := new(mocks.DecisionManager)
	decisionManager.On("DecideOrderPlacement", mock.Anything).Once().Return(order.PlaceOurExchange, nil)
	mqttManager := new(mocks.CentrifugoManager)
	mqttManager.On("PublishOrderToOpenOrders", mock.Anything, mock.Anything, mock.Anything).Once().Return()
	externalExchangeOrderService := new(mocks.ExternalExchangeOrderService)
	engineCommunicator := new(mocks.EngineCommunicator)
	engineCommunicator.On("SubmitOrder", mock.Anything).Once().Return(nil)
	postOrderMatchingService := new(mocks.PostOrderMatchingService)
	logger := new(mocks.Logger)

	eh := order.NewOrderEventsHandler(redisManager, decisionManager, mqttManager, externalExchangeOrderService, engineCommunicator,
		postOrderMatchingService, logger)

	o := order.Order{
		Price:          sql.NullString{String: "51000", Valid: true},
		StopPointPrice: sql.NullString{String: "50000", Valid: true},
	}
	eh.HandleOrderCreation(o, true)
	time.Sleep(50 * time.Millisecond)
	decisionManager.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
	engineCommunicator.AssertExpectations(t)
}

func TestEventsHandler_HandleOrderCreation_ForStopOrderSubmission_ExternalExchange(t *testing.T) {
	redisManager := new(mocks.OrderRedisManager)
	decisionManager := new(mocks.DecisionManager)
	decisionManager.On("DecideOrderPlacement", mock.Anything).Once().Return(order.PlaceExternalExchange, nil)
	mqttManager := new(mocks.CentrifugoManager)
	mqttManager.On("PublishOrderToOpenOrders", mock.Anything, mock.Anything, mock.Anything).Once().Return()

	externalExchangeOrderService := new(mocks.ExternalExchangeOrderService)
	orderResult := externalexchange.UserOrderResult{IsOrderPlaced: true}
	externalExchangeOrderService.On("CreateExternalExchangeOrderForUser", mock.Anything).Once().Return(orderResult, nil)

	engineCommunicator := new(mocks.EngineCommunicator)
	postOrderMatchingService := new(mocks.PostOrderMatchingService)
	postOrderMatchingService.On("HandleExternalTradedOrder", mock.Anything).Once().Return(nil)
	logger := new(mocks.Logger)

	eh := order.NewOrderEventsHandler(redisManager, decisionManager, mqttManager, externalExchangeOrderService, engineCommunicator,
		postOrderMatchingService, logger)

	o := order.Order{}
	eh.HandleOrderCreation(o, true)
	time.Sleep(50 * time.Millisecond)
	decisionManager.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
	externalExchangeOrderService.AssertExpectations(t)
	postOrderMatchingService.AssertExpectations(t)
}

func TestEventsHandler_HandleOrderCreation_ExternalExchangeReturnsError(t *testing.T) {
	redisManager := new(mocks.OrderRedisManager)
	decisionManager := new(mocks.DecisionManager)
	decisionManager.On("DecideOrderPlacement", mock.Anything).Once().Return(order.PlaceExternalExchange, nil)
	mqttManager := new(mocks.CentrifugoManager)
	externalExchangeOrderService := new(mocks.ExternalExchangeOrderService)
	orderResult := externalexchange.UserOrderResult{IsOrderPlaced: false}
	externalExchangeOrderService.On("CreateExternalExchangeOrderForUser", mock.Anything).Once().Return(orderResult, nil)

	engineCommunicator := new(mocks.EngineCommunicator)
	engineCommunicator.On("SubmitOrder", mock.Anything).Once().Return(nil)
	postOrderMatchingService := new(mocks.PostOrderMatchingService)
	logger := new(mocks.Logger)

	eh := order.NewOrderEventsHandler(redisManager, decisionManager, mqttManager, externalExchangeOrderService, engineCommunicator,
		postOrderMatchingService, logger)

	o := order.Order{}

	eh.HandleOrderCreation(o, true)
	decisionManager.AssertExpectations(t)
	externalExchangeOrderService.AssertExpectations(t)
	engineCommunicator.AssertExpectations(t)
}

func TestEventsHandler_HandleOrderCancellation(t *testing.T) {
	redisManager := new(mocks.OrderRedisManager)
	decisionManager := new(mocks.DecisionManager)
	mqttManager := new(mocks.CentrifugoManager)
	mqttManager.On("PublishOrderToOpenOrders", mock.Anything, mock.Anything, mock.Anything).Once().Return()
	externalExchangeOrderService := new(mocks.ExternalExchangeOrderService)
	engineCommunicator := new(mocks.EngineCommunicator)
	postOrderMatchingService := new(mocks.PostOrderMatchingService)
	logger := new(mocks.Logger)
	eh := order.NewOrderEventsHandler(redisManager, decisionManager, mqttManager, externalExchangeOrderService, engineCommunicator,
		postOrderMatchingService, logger)
	o := order.Order{}
	eh.HandleOrderCancellation(o)
	time.Sleep(50 * time.Millisecond)
	engineCommunicator.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
}

func TestEventsHandler_HandleOrderFulfillByAdmin(t *testing.T) {
	redisManager := new(mocks.OrderRedisManager)
	decisionManager := new(mocks.DecisionManager)
	mqttManager := new(mocks.CentrifugoManager)
	externalExchangeOrderService := new(mocks.ExternalExchangeOrderService)
	engineCommunicator := new(mocks.EngineCommunicator)
	engineCommunicator.On("RemoveOrder", mock.Anything).Once().Return(nil)
	postOrderMatchingService := new(mocks.PostOrderMatchingService)
	logger := new(mocks.Logger)
	eh := order.NewOrderEventsHandler(redisManager, decisionManager, mqttManager, externalExchangeOrderService, engineCommunicator,
		postOrderMatchingService, logger)
	o := order.Order{}
	eh.HandleOrderFulfillByAdmin(o)
	engineCommunicator.AssertExpectations(t)
}
