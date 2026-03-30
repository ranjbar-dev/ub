package order

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

func (ps *postOrderMatchingService) addToPushData(orderItem MatchingNeededQueryFields, pairName string) {
	for _, item := range ps.pushData {
		if item.ID == orderItem.OrderID {
			return
		}
	}

	amount := orderItem.PayedByAmount
	if orderItem.OrderType == TypeBuy {
		amount = orderItem.DemandedAmount
	}
	payload := orderPushPayload{
		ID:                 orderItem.OrderID,
		Amount:             amount,
		Price:              orderItem.Price,
		Status:             strings.ToLower(StatusFilled),
		OrderType:          strings.ToLower(orderItem.OrderType),
		UserPrivateChannel: orderItem.UserPrivateChannel,
		Pair:               pairName,
	}
	ps.pushData = append(ps.pushData, payload)

}

//todo this method should not be here pushing data to client should be independent from this service
func (ps *postOrderMatchingService) pushDataToUsers(payloads []orderPushPayload) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ctx := context.Background()
	for _, item := range payloads {
		payload, err := json.Marshal(item)
		if err != nil {
			ps.logger.Error2("can not push order to clients", err,
				zap.String("service", "postOrderMatchingService"),
				zap.String("method", "pushDataToUsers"),
				zap.Int64("orderID", item.ID),
			)
		}
		ps.mqttManager.PublishOrderToOpenOrders(ctx, item.UserPrivateChannel, payload)
	}
}

func (ps *postOrderMatchingService) storeUnmatchedOrderID(orderID int64) {
	ctx := context.Background()
	orderIDString := strconv.FormatInt(orderID, 10)
	_, err := ps.rc.LPush(ctx, UnmatchedOrdersList, orderIDString)
	if err != nil {
		ps.logger.Error2("can not push to redis", err,
			zap.String("service", "postOrderMatchingService"),
			zap.String("method", "storeUnmatchedOrderId"),
			zap.Int64("orderID", orderID),
		)
	}
}
