package order

import (
	"context"
	"encoding/json"
	"exchange-go/internal/communication"
	"exchange-go/internal/externalexchange"
	"exchange-go/internal/platform"
	"strings"

	"go.uber.org/zap"
)

// EventsHandler publishes order lifecycle events such as creation, cancellation,
// and admin fulfillment to downstream services and clients.
type EventsHandler interface {
	// HandleOrderCreation processes a newly created order by routing it to the matching
	// engine, an external exchange, or the Redis stop-order queue. When isForStopOrderSubmission
	// is true, the order is treated as a triggered stop order being submitted to the engine.
	HandleOrderCreation(o Order, isForStopOrderSubmission bool)
	// HandleOrderCancellation pushes a cancellation notification to the client via MQTT.
	HandleOrderCancellation(o Order)
	// HandleOrderFulfillByAdmin removes the order from the matching engine's order book
	// after an admin manually fulfills it.
	HandleOrderFulfillByAdmin(o Order)
}

type eventsHandler struct {
	redisManager                 RedisManager
	decisionManager              DecisionManager
	mqttManager                  communication.MqttManager
	externalExchangeOrderService externalexchange.OrderService
	engineCommunicator           EngineCommunicator
	postOrderMatchingService     PostOrderMatchingService
	logger                       platform.Logger
}

func (s *eventsHandler) HandleOrderCreation(o Order, isForStopOrderSubmission bool) {
	if o.IsStopOrder() && !isForStopOrderSubmission {
		ctx := context.Background()
		//add to redis only
		err := s.redisManager.AddStopOrderToQueue(ctx, o)
		if err != nil {
			s.logger.Error2("can not add stopOrder to queue", err,
				zap.String("service", "orderEventsHandler"),
				zap.String("method", "HandleOrderCreation"),
				zap.Int64("orderID", o.ID),
			)
		}
		go s.pushOrderToClient(o)
		return
	} else {
		place, err := s.decisionManager.DecideOrderPlacement(o)
		if err != nil {
			s.logger.Error2("can not decide where to place the order", err,
				zap.String("service", "orderEventsHandler"),
				zap.String("method", "HandleOrderCreation"),
				zap.Int64("orderID", o.ID),
			)
			return
		}

		if place == PlaceExternalExchange {
			externalExchangeOrderParams := externalexchange.UserOrderParams{
				PairID:       o.Pair.ID,
				PairName:     o.Pair.Name,
				Type:         o.Type,
				ExchangeType: o.ExchangeType,
				Amount:       o.getAmount(),
				Price:        o.Price.String,
				OrderID:      o.ID,
			}

			result, err := s.externalExchangeOrderService.CreateExternalExchangeOrderForUser(externalExchangeOrderParams)
			if err != nil {
				s.logger.Error2("can not create externalExchangeOrderForUser", err,
					zap.String("service", "orderEventsHandler"),
					zap.String("method", "HandleOrderCreation"),
					zap.Int64("orderID", o.ID),
				)
			}

			if result.IsOrderPlaced {
				s.doPostOrderPlacementInExternalExchange(o, result.ExternalExchangeID, result.ExternalExchangeOrderID, result.Data)
				go s.pushOrderToClient(o)
				return
			}
		}

		err = s.engineCommunicator.SubmitOrder(o)
		if err != nil {
			s.logger.Error2("can not submit order to engine", err,
				zap.String("service", "orderEventsHandler"),
				zap.String("method", "HandleOrderCreation"),
				zap.Int64("orderID", o.ID),
			)
			return
		}
	}
	if s.shouldPush(o) {
		go s.pushOrderToClient(o)

	}

}

func (s *eventsHandler) doPostOrderPlacementInExternalExchange(o Order, externalExchangeID int64, externalExchangeOrderID string, data string) {
	params := ExternalTradedOrderData{
		OrderID:                 o.ID,
		ExtraInfoID:             o.ExtraInfoID.Int64,
		Data:                    data,
		ExternalExchangeID:      externalExchangeID,
		ExternalExchangeOrderID: externalExchangeOrderID,
		Pair:                    o.Pair,
	}
	err := s.postOrderMatchingService.HandleExternalTradedOrder(params)
	if err != nil {
		s.logger.Error2("can not handle externalTradedOrder", err,
			zap.String("service", "orderEventsHandler"),
			zap.String("method", "doPostOrderPlacementInExternalExchange"),
			zap.Int64("orderID", o.ID),
		)
		return
	}
}

func (s *eventsHandler) HandleOrderCancellation(o Order) {
	go s.pushOrderToClient(o)

}

func (s *eventsHandler) HandleOrderFulfillByAdmin(o Order) {
	//TODO in case of error occurrence we should have backup plan to be sure the order is removed from order book
	err := s.engineCommunicator.RemoveOrder(o)
	if err != nil {
		s.logger.Error2("can not remove order from engine", err,
			zap.String("service", "orderEventsHandler"),
			zap.String("method", "HandleOrderFulfillByAdmin"),
			zap.Int64("orderID", o.ID),
		)
	}
}

type orderPushPayload struct {
	ID                 int64  `json:"id"`
	Amount             string `json:"amount"`
	Price              string `json:"price"`
	Status             string `json:"status"`
	OrderType          string `json:"type"`
	Pair               string `json:"pairCurrency"`
	UserPrivateChannel string `json:"-"`
}

func (s *eventsHandler) pushOrderToClient(o Order) {

	ctx := context.Background()
	channelName := o.User.PrivateChannelName

	status := strings.ToLower(o.Status)
	if o.isMarket() {
		status = strings.ToLower(StatusFilled)
	}

	p := orderPushPayload{
		ID:        o.ID,
		Amount:    o.getAmount(),
		Price:     o.Price.String,
		Status:    status,
		OrderType: strings.ToLower(o.Type),
		Pair:      o.Pair.Name,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		s.logger.Error2("can not push order to Client", err,
			zap.String("service", "orderEventsHandler"),
			zap.String("method", "pushOrderToClient"),
			zap.Int64("orderID", o.ID),
		)
	}

	s.mqttManager.PublishOrderToOpenOrders(ctx, channelName, payload)
}

func (s *eventsHandler) shouldPush(o Order) bool {
	if o.isMarket() {
		return false
	}

	return true
}

func NewOrderEventsHandler(rm RedisManager, dm DecisionManager, mqttM communication.MqttManager,
	ees externalexchange.OrderService, ec EngineCommunicator, pom PostOrderMatchingService,
	logger platform.Logger) EventsHandler {
	return &eventsHandler{
		redisManager:                 rm,
		decisionManager:              dm,
		mqttManager:                  mqttM,
		externalExchangeOrderService: ees,
		engineCommunicator:           ec,
		postOrderMatchingService:     pom,
		logger:                       logger,
	}
}
