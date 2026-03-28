package order

import (
	"exchange-go/internal/currency"
	"exchange-go/internal/userbalance"

	"gorm.io/gorm"
)

func (ps *postOrderMatchingService) createGroups(tx *gorm.DB, pair currency.Pair, orderItems []MatchingNeededQueryFields, doneOrders []CallBackOrderData, partial *CallBackOrderData, isFromAdmin bool) []orderGroup {
	var orderGroups []orderGroup
	var userIds []int
	for _, orderItem := range orderItems {
		userIds = append(userIds, orderItem.UserID)
	}

	coinIds := []int64{pair.BasisCoinID, pair.DependentCoinID}
	allUserBalances := ps.userBalanceService.GetBalancesOfUsersForCoinsUsingTx(tx, userIds, coinIds)

	for _, doneOrder := range doneOrders {
		orderID := doneOrder.ID
		tradedWithOrderID := doneOrder.TradedWithOrderID
		isMaker := false
		makerOrderID := ps.getTheMakerOrderID(doneOrders, orderID, tradedWithOrderID)
		if orderID == makerOrderID {
			isMaker = true
		}
		to := tempOrder{
			tradePrice:        doneOrder.TradePrice,
			tradeAmount:       doneOrder.QuantityTraded,
			isMaker:           isMaker,
			orderType:         doneOrder.OrderType,
			marketPrice:       doneOrder.MarketPrice,
			isPartial:         false,
			TradedWithOrderID: tradedWithOrderID,
		}

		orderItem := ps.getOrderByID(orderItems, orderID)
		exists := false
		for i, groupOrder := range orderGroups {
			if groupOrder.orderItem.OrderID == orderID {
				exists = true
				orderGroups[i].tempOrders = append(orderGroups[i].tempOrders, to)
				break
			}
		}
		if !exists {
			userBalances := [2]*userbalance.UserBalance{{}, {}}
			for i, ub := range allUserBalances {
				if ub.UserID == orderItem.UserID {
					if userBalances[0].ID == 0 { //checking if it is already does not exist
						userBalances[0] = &allUserBalances[i]
					} else {
						userBalances[1] = &allUserBalances[i]
					}
				}
			}

			newGroupOrder := orderGroup{
				tempOrders:   []tempOrder{to},
				orderItem:    orderItem,
				userBalances: userBalances,
			}
			orderGroups = append(orderGroups, newGroupOrder)
		}
	}

	if partial != nil {
		to := tempOrder{
			tradePrice:        "",
			tradeAmount:       "",
			orderType:         partial.OrderType,
			marketPrice:       partial.MarketPrice,
			isMaker:           false,
			isPartial:         true,
			TradedWithOrderID: 0,
		}
		if isFromAdmin {
			to.tradePrice = partial.Price
		}

		exists := false
		for i, groupOrder := range orderGroups {
			if groupOrder.orderItem.OrderID == partial.ID {
				orderGroups[i].tempOrders = append(orderGroups[i].tempOrders, to)
				exists = true
				break
			}
		}

		if !exists {
			orderItem := ps.getOrderByID(orderItems, partial.ID)
			userBalances := [2]*userbalance.UserBalance{{}, {}}
			for i, ub := range allUserBalances {
				if ub.UserID == orderItem.UserID {
					if userBalances[0].ID == 0 { //checking if it is already does not exist
						userBalances[0] = &allUserBalances[i]
					} else {
						userBalances[1] = &allUserBalances[i]
					}
				}
			}
			newGroupOrder := orderGroup{
				tempOrders:   []tempOrder{to},
				orderItem:    orderItem,
				userBalances: userBalances,
			}
			orderGroups = append(orderGroups, newGroupOrder)

		}
	}
	return orderGroups
}

func (ps *postOrderMatchingService) getTheMakerOrderID(ordersData []CallBackOrderData, orderID int64, tradedWithOrderID int64) (makerOrderID int64) {
	//var orderTimestamp int64
	var orderPrice string
	//var tradedWithOrderTimestamp int64
	var tradedWithOrderPrice string
	for _, orderData := range ordersData {
		if orderData.ID == orderID {
			//orderTimestamp = orderData.Timestamp
			orderPrice = orderData.Price
		}

		if orderData.ID == tradedWithOrderID {
			//tradedWithOrderTimestamp = orderData.Timestamp
			tradedWithOrderPrice = orderData.Price
		}
	}

	if orderPrice != "" && tradedWithOrderPrice != "" {
		//it means both orders are limit so the older one is maker
		//if orderTimestamp == tradedWithOrderTimestamp {
		//check the id here
		if orderID < tradedWithOrderID {
			return orderID
		}
		return tradedWithOrderID

		//}
		//if orderTimestamp < tradedWithOrderTimestamp {
		//	return orderId
		//}
		//return tradedWithOrderId
	}

	//the limit one is maker
	if orderPrice != "" {
		return orderID
	}

	return tradedWithOrderID

}

func (ps *postOrderMatchingService) getOrderByID(orderItems []MatchingNeededQueryFields, orderID int64) MatchingNeededQueryFields {
	for _, orderItem := range orderItems {
		if orderItem.OrderID == orderID {
			return orderItem
		}
	}
	//we should never reach here
	return MatchingNeededQueryFields{}
}
