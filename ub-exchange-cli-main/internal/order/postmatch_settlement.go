package order

import (
	"context"
	"database/sql"
	"fmt"

	"exchange-go/internal/currency"
	"exchange-go/internal/platform"
	"exchange-go/internal/userbalance"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (ps *postOrderMatchingService) doPreMatchingActions(pairName string) {
	//we empty the tempTrades in case former one exists in it
	//also set the currentMarketPrice
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.tempTrades = make([]tempTrade, 0)
	ps.tradesData = make([]TradeData, 0)
	ps.pushData = make([]orderPushPayload, 0)
	ps.currentMarketPrice, _ = ps.priceGenerator.GetPrice(context.Background(), pairName)
}

func (ps *postOrderMatchingService) HandlePostOrderMatching(doneOrders []CallBackOrderData, partial *CallBackOrderData, isFromAdmin bool) MatchingResult {
	var pairName string
	orderIdSet := make(map[int64]bool, len(doneOrders)+1)
	for _, doneOrder := range doneOrders {
		pairName = doneOrder.PairName
		orderIdSet[doneOrder.ID] = true
	}

	if partial != nil {
		pairName = partial.PairName
		orderIdSet[partial.ID] = true
	}

	orderIds := make([]int64, 0, len(orderIdSet))
	for id := range orderIdSet {
		orderIds = append(orderIds, id)
	}

	ps.doPreMatchingActions(pairName)

	var removingDoneOrderIds []int64
	tx := ps.db.Begin()
	err := tx.Error
	if err != nil {
		ps.logger.Error2("can not start db transaction", err,
			zap.String("service", "postOrderMatchingService"),
			zap.String("method", "HandlePostOrderMatching"),
		)
		return ps.prepareResult(
			doneOrders,
			partial,
			true,
			removingDoneOrderIds,
			err,
			nil,
		)
	}

	//following lines are written for the sake of performance and locking mechanism in mysql
	// first getting orders with write lock  then get users and userLevels neededData
	//then integrate them
	orderItems := ps.orderRepository.GetOrdersDataByIdsWithJoinUsingTx(tx, orderIds)
	var userIds []int
	for _, o := range orderItems {
		userIds = append(userIds, o.UserID)
	}
	usersData := ps.userService.GetUsersDataForOrderMatching(userIds)
	userLevelIdSet := make(map[int64]bool, len(usersData))
	for _, ud := range usersData {
		userLevelIdSet[ud.UserLevelID] = true
	}
	userLevelIds := make([]int64, 0, len(userLevelIdSet))
	for id := range userLevelIdSet {
		userLevelIds = append(userLevelIds, id)
	}
	userLevels := ps.userLevelService.GetLevelsByIds(userLevelIds)

	// Build lookup maps for O(1) enrichment
	type userInfo struct {
		UserEmail          string
		UserPrivateChannel string
		UserLevelID        int64
	}
	userDataMap := make(map[int]userInfo, len(usersData))
	for _, ud := range usersData {
		userDataMap[ud.UserID] = userInfo{
			UserEmail:          ud.UserEmail,
			UserPrivateChannel: ud.UserPrivateChannel,
			UserLevelID:        ud.UserLevelID,
		}
	}

	type levelFees struct {
		MakerFeePercentage float64
		TakerFeePercentage float64
	}
	userLevelMap := make(map[int64]levelFees, len(userLevels))
	for _, ul := range userLevels {
		userLevelMap[ul.ID] = levelFees{
			MakerFeePercentage: ul.MakerFeePercentage,
			TakerFeePercentage: ul.TakerFeePercentage,
		}
	}

	//we complete data for orderItems
	for i, o := range orderItems {
		if ud, ok := userDataMap[o.UserID]; ok {
			orderItems[i].UserEmail = ud.UserEmail
			orderItems[i].UserPrivateChannel = ud.UserPrivateChannel
			orderItems[i].UserLevelID = ud.UserLevelID
			if lf, ok := userLevelMap[ud.UserLevelID]; ok {
				orderItems[i].MakerFeePercentage = lf.MakerFeePercentage
				orderItems[i].TakerFeePercentage = lf.TakerFeePercentage
			}
		}
	}

	//here we check if any of orders in not open then it means we should not go further
	//and delete them from orderbook
	isPartialOrderOpen := true
	for _, orderItem := range orderItems {
		for _, doneOrder := range doneOrders {
			if orderItem.OrderID == doneOrder.ID {
				if orderItem.Status != StatusOpen {
					exists := false
					for _, id := range removingDoneOrderIds {
						if id == orderItem.OrderID {
							exists = true
							break
						}
					}
					if !exists {
						removingDoneOrderIds = append(removingDoneOrderIds, orderItem.OrderID)

					}
				}
			}
		}

		if partial != nil && orderItem.OrderID == partial.ID && orderItem.Status != StatusOpen {
			isPartialOrderOpen = false
		}
	}
	if len(removingDoneOrderIds) > 0 || !isPartialOrderOpen {
		tx.Rollback()
		partialOrderID := int64(0)
		if partial != nil {
			partialOrderID = partial.ID //just for logging
		}
		err := fmt.Errorf("at least one order status is not open")
		ps.logger.Error2("order status is not open", err,
			zap.String("service", "postOrderMatchingService"),
			zap.String("method", "HandlePostOrderMatching"),
			zap.Int64s("doneOrderIDs", removingDoneOrderIds),
			zap.Int64("partialOrderId", partialOrderID),
		)
		return ps.prepareResult(
			doneOrders,
			partial,
			isPartialOrderOpen,
			removingDoneOrderIds,
			err,
			nil,
		)
	}

	pair, _ := ps.currencyService.GetPairByName(pairName)

	orderGroups := ps.createGroups(tx, pair, orderItems, doneOrders, partial, isFromAdmin)

	var remainingPartialOrder *CallBackOrderData
	for _, group := range orderGroups {
		partialOrder, err := ps.handleOrderGroup(tx, group, pair, isFromAdmin)
		if partialOrder != nil {
			remainingPartialOrder = partialOrder
		}
		if err != nil {
			tx.Rollback()
			ps.logger.Error2("error handling OrderGroup", err,
				zap.String("service", "postOrderMatchingService"),
				zap.String("method", "HandlePostOrderMatching"),
				zap.String("pairName", pair.Name),
				zap.Bool("isFromAdmin", isFromAdmin),
				zap.Int64("orderId", group.orderItem.OrderID),
			)
			return ps.prepareResult(
				doneOrders,
				partial,
				true,
				removingDoneOrderIds,
				err,
				nil,
			)
		}
	}

	//this line are just for test purposes
	if ps.configs.GetEnv() == platform.EnvTest {
		if ps.configs.GetBool("commitError") {
			err = fmt.Errorf("test env error")
		}
	}
	if err != nil {
		tx.Rollback()
		return ps.prepareResult(
			doneOrders,
			partial,
			true,
			removingDoneOrderIds,
			err,
			nil,
		)
	}
	//end of test scnenarios
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		ps.logger.Error2("error commiting transaction", err,
			zap.String("service", "postOrderMatchingService"),
			zap.String("method", "HandlePostOrderMatching"),
			zap.String("pairName", pair.Name),
			zap.Bool("pairName", isFromAdmin),
		)
		return ps.prepareResult(
			doneOrders,
			partial,
			true,
			removingDoneOrderIds,
			err,
			nil,
		)
	}

	localTrades := make([]TradeData, len(ps.tradesData))
	copy(localTrades, ps.tradesData)
	localPush := make([]orderPushPayload, len(ps.pushData))
	copy(localPush, ps.pushData)
	go ps.tradeEventsHandler.HandleTradesCreation(localTrades, pair)
	go ps.pushDataToUsers(localPush)
	for _, doneOrder := range doneOrders {
		removingDoneOrderIds = append(removingDoneOrderIds, doneOrder.ID)
	}
	return ps.prepareResult(
		doneOrders,
		partial,
		true,
		removingDoneOrderIds,
		nil,
		remainingPartialOrder,
	)
}

func (ps *postOrderMatchingService) handleOrderGroup(tx *gorm.DB, group orderGroup, pair currency.Pair, isFromAdmin bool) (remainingPartialOrder *CallBackOrderData, err error) {
	var orders []Order
	var childTempOrders []tempOrder
	var mainOrder *tempOrder
	var partial *tempOrder

	demandedDiffDecimal := decimal.NewFromFloat(0)
	payedByDiffDecimal := decimal.NewFromFloat(0)
	frozenAmountReductionDecimal := decimal.NewFromFloat(0)

	isMarket := false
	if group.orderItem.Price == "" {
		isMarket = true
	}

	for _, tempOrder := range group.tempOrders {
		if tempOrder.isPartial {
			partial = &tempOrder
			break
		}
	}

	if partial != nil {
		if len(group.tempOrders) != 1 {
			mainOrder = &group.tempOrders[0]
			childTempOrders = group.tempOrders[1 : len(group.tempOrders)-1]
		}
	} else {
		mainOrder = &group.tempOrders[0]
		if len(group.tempOrders) > 1 {
			childTempOrders = group.tempOrders[1:]
		}
	}

	if mainOrder != nil {
		finalDemandedAmountDecimal, finalPayedByAmountDecimal := ps.getFinalAmountsDecimal(*mainOrder)
		feePercentage := ps.GetFeePercentage(group.orderItem, pair, isMarket, mainOrder.isMaker)
		feePercentageDecimal := decimal.NewFromFloat(feePercentage)
		feeDecimal := feePercentageDecimal.Mul(finalDemandedAmountDecimal)

		demandedDiffDecimal = demandedDiffDecimal.Add(finalDemandedAmountDecimal)
		demandedDiffDecimal = demandedDiffDecimal.Sub(feeDecimal) //fee amount is reduced from demanded
		payedByDiffDecimal = payedByDiffDecimal.Add(finalPayedByAmountDecimal)
		//doing this so reducing the frozen amount from user balance correctly
		if group.orderItem.Price != "" {
			orderPriceDecimal, _ := decimal.NewFromString(group.orderItem.Price)
			tradeAmountDecimal, _ := decimal.NewFromString(mainOrder.tradeAmount)
			if group.orderItem.OrderType == TypeBuy {
				frozenAmountReductionDecimal = frozenAmountReductionDecimal.Add(orderPriceDecimal.Mul(tradeAmountDecimal))
			} else {
				frozenAmountReductionDecimal = frozenAmountReductionDecimal.Add(tradeAmountDecimal)
			}
		} else {
			frozenAmountReductionDecimal = frozenAmountReductionDecimal.Add(finalPayedByAmountDecimal)
		}

		o, err := ps.updateMainOrder(tx, group.orderItem, *mainOrder, pair, isMarket)
		if err != nil {
			return nil, fmt.Errorf("error saving mainOrder %d: %w", group.orderItem.OrderID, err)
		}

		//create trade for orders
		sellOrderID := group.orderItem.OrderID
		buyOrderID := mainOrder.TradedWithOrderID
		if group.orderItem.OrderType == TypeBuy {
			buyOrderID = group.orderItem.OrderID
			sellOrderID = mainOrder.TradedWithOrderID
		}
		tt := tempTrade{
			price:       mainOrder.tradePrice,
			amount:      mainOrder.tradeAmount,
			buyOrderID:  buyOrderID,
			sellOrderID: sellOrderID,
			pair:        pair,
		}
		err = ps.handleTrade(tx, tt, pair, "", 0)
		if err != nil {
			return nil, fmt.Errorf("handleOrderGroup: handle main trade: %w", err)
		}

		orders = append(orders, o)
		ps.addToPushData(group.orderItem, pair.Name)

		parenOrder := o
		for _, childTempOrder := range childTempOrders {
			finalDemandedAmountDecimal, finalPayedByAmountDecimal := ps.getFinalAmountsDecimal(childTempOrder)
			feePercentage := ps.GetFeePercentage(group.orderItem, pair, isMarket, childTempOrder.isMaker)
			feePercentageDecimal := decimal.NewFromFloat(feePercentage)
			feeDecimal := feePercentageDecimal.Mul(finalDemandedAmountDecimal)
			demandedDiffDecimal = demandedDiffDecimal.Add(finalDemandedAmountDecimal)
			demandedDiffDecimal = demandedDiffDecimal.Sub(feeDecimal) //fee amount is reduced from demanded
			payedByDiffDecimal = payedByDiffDecimal.Add(finalPayedByAmountDecimal)
			//doing this so reducing the frozen amount from user balance correctly
			if group.orderItem.Price != "" {
				orderPriceDecimal, _ := decimal.NewFromString(group.orderItem.Price)
				tradeAmountDecimal, _ := decimal.NewFromString(childTempOrder.tradeAmount)
				if group.orderItem.OrderType == TypeBuy {
					frozenAmountReductionDecimal = frozenAmountReductionDecimal.Add(orderPriceDecimal.Mul(tradeAmountDecimal))
				} else {
					frozenAmountReductionDecimal = frozenAmountReductionDecimal.Add(tradeAmountDecimal)
				}
			} else {
				frozenAmountReductionDecimal = frozenAmountReductionDecimal.Add(finalPayedByAmountDecimal)
			}

			childOrder, err := ps.createChildOrder(tx, group.orderItem, childTempOrder, parenOrder, pair, StatusFilled)
			if err != nil {
				return nil, fmt.Errorf("handleOrderGroup: create child order: %w", err)
			}

			//create trade for orders
			sellOrderID := childOrder.ID
			buyOrderID := childTempOrder.TradedWithOrderID
			if group.orderItem.OrderType == TypeBuy {
				buyOrderID = childOrder.ID
				sellOrderID = childTempOrder.TradedWithOrderID
			}

			tt := tempTrade{
				price:       childTempOrder.tradePrice,
				amount:      childTempOrder.tradeAmount,
				buyOrderID:  buyOrderID,
				sellOrderID: sellOrderID,
				pair:        pair,
			}
			err = ps.handleTrade(tx, tt, pair, "", 0)
			if err != nil {
				return nil, fmt.Errorf("handleOrderGroup: handle child trade: %w", err)
			}

			orders = append(orders, childOrder)
			parenOrder = childOrder
		}
	}

	if partial != nil {
		var parentOrder *Order
		if len(orders) > 0 {
			parentOrder = &orders[len(orders)-1] //the last one is parent order
		}

		partialOrderHandlingResult := ps.handlePartialOrder(tx, group.orderItem, *partial, parentOrder, pair, isFromAdmin)
		if partialOrderHandlingResult.err != nil {
			return nil, partialOrderHandlingResult.err
		}
		order := partialOrderHandlingResult.order
		if partialOrderHandlingResult.isTraded {
			finalDemandedAmountDecimal, _ := decimal.NewFromString(order.FinalDemandedAmount.String)
			finalPayedByAmountDecimal, _ := decimal.NewFromString(order.FinalPayedByAmount.String)
			feePercentage := ps.GetFeePercentage(group.orderItem, pair, isMarket, false)
			feePercentageDecimal := decimal.NewFromFloat(feePercentage)
			feeDecimal := feePercentageDecimal.Mul(finalDemandedAmountDecimal)
			demandedDiffDecimal = demandedDiffDecimal.Add(finalDemandedAmountDecimal)
			demandedDiffDecimal = demandedDiffDecimal.Sub(feeDecimal) //fee amount is reduced from demanded
			payedByDiffDecimal = payedByDiffDecimal.Add(finalPayedByAmountDecimal)
			frozenAmountReductionDecimal = frozenAmountReductionDecimal.Add(finalPayedByAmountDecimal)
			orders = append(orders, order)

			tempTrade := tempTrade{
				price: order.TradePrice.String,
				pair:  pair,
			}
			if group.orderItem.OrderType == TypeBuy {
				tempTrade.buyOrderID = order.ID
				tempTrade.amount = order.FinalDemandedAmount.String
			} else {
				tempTrade.sellOrderID = order.ID
				tempTrade.amount = order.FinalPayedByAmount.String
			}

			userID := group.orderItem.UserID
			userEmail := group.orderItem.UserEmail
			err := ps.handleTrade(tx, tempTrade, pair, userEmail, userID)
			if err != nil {
				return nil, fmt.Errorf("handleOrderGroup: handle partial trade: %w", err)
			}
			ps.addToPushData(group.orderItem, pair.Name)
		} else {
			min, max, err := ps.forceTrader.GetMinAndMaxPrice(pair.Name, order.Type, ps.currentMarketPrice)
			if err != nil {
				return nil, fmt.Errorf("handleOrderGroup: get min/max price: %w", err)
			}
			quantity := order.DemandedAmount.String
			if order.Type == TypeSell {
				quantity = order.PayedByAmount.String
			}

			remainingPartialOrder = &CallBackOrderData{
				ID:                order.ID,
				PairName:          pair.Name,
				OrderType:         order.Type,
				Quantity:          quantity,
				Price:             order.Price.String,
				Timestamp:         order.CreatedAt.Unix(),
				TradedWithOrderID: 0,
				QuantityTraded:    "",
				TradePrice:        "",
				MarketPrice:       group.orderItem.MarketPrice,
				MinThresholdPrice: min,
				MaxThresholdPrice: max,
			}
		}

	}
	err = ps.createTransactions(tx, orders, pair)
	if err != nil {
		return nil, fmt.Errorf("handleOrderGroup: create transactions: %w", err)
	}

	err = ps.updateUserBalances(tx, group.userBalances, demandedDiffDecimal, payedByDiffDecimal, frozenAmountReductionDecimal, group.orderItem.OrderType, pair)
	if err != nil {
		return remainingPartialOrder, fmt.Errorf("handleOrderGroup: update user balances: %w", err)
	}
	return remainingPartialOrder, nil
}

func (ps *postOrderMatchingService) HandleExternalTradedOrder(data ExternalTradedOrderData) error {
	pair := data.Pair
	ctx := context.Background()
	tx := ps.db.Begin()
	err := tx.Error
	if err != nil {
		return fmt.Errorf("HandleExternalTradedOrder: begin tx: %w", err)
	}
	orderIds := []int64{data.OrderID}
	orderItems := ps.orderRepository.GetOrdersDataByIdsWithJoinUsingTx(tx, orderIds)
	orderItem := orderItems[0]

	isMarket := false
	if orderItem.Price == "" {
		isMarket = true
	}

	tradePrice, err := ps.priceGenerator.GetPrice(ctx, pair.Name)
	if err != nil {
		return fmt.Errorf("HandleExternalTradedOrder: get price: %w", err)
	}

	tradeAmount := orderItem.PayedByAmount
	if orderItem.OrderType == TypeBuy {
		tradeAmount = orderItem.DemandedAmount

	}

	tempOrder := tempOrder{
		tradePrice:  tradePrice,
		tradeAmount: tradeAmount,
		orderType:   orderItem.OrderType,
		marketPrice: orderItem.MarketPrice,
		isMaker:     false,
		isPartial:   false,
	}

	order, err := ps.updateMainOrder(tx, orderItem, tempOrder, pair, isMarket)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("HandleExternalTradedOrder: update main order: %w", err)
	}

	extraInfo := &ExtraInfo{
		ID:                        data.ExtraInfoID,
		ExternalExchangeOtherInfo: sql.NullString{String: data.Data, Valid: true},
		ExternalExchangeID:        sql.NullInt64{Int64: data.ExternalExchangeID, Valid: true},
		ExternalExchangeOrderID:   sql.NullString{String: data.ExternalExchangeOrderID, Valid: true},
	}

	err = tx.Model(extraInfo).Updates(extraInfo).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("HandleExternalTradedOrder: update extra info: %w", err)
	}
	orders := []Order{order}
	err = ps.createTransactions(tx, orders, pair)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("HandleExternalTradedOrder: create transactions: %w", err)
	}

	userIds := []int{orderItem.UserID}
	coinIds := []int64{pair.BasisCoinID, pair.DependentCoinID}
	allUserBalances := ps.userBalanceService.GetBalancesOfUsersForCoinsUsingTx(tx, userIds, coinIds)
	userBalances := [2]*userbalance.UserBalance{&allUserBalances[0], &allUserBalances[1]}
	demandedDecimal, _ := decimal.NewFromString(order.FinalDemandedAmount.String)
	payedByDecimal, _ := decimal.NewFromString(order.FinalPayedByAmount.String)
	frozenReductionDecimal, _ := decimal.NewFromString(order.FinalPayedByAmount.String)

	err = ps.updateUserBalances(tx, userBalances, demandedDecimal, payedByDecimal, frozenReductionDecimal, orderItem.OrderType, pair)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("HandleExternalTradedOrder: update user balances: %w", err)
	}

	tempTrade := tempTrade{
		price: order.TradePrice.String,
		pair:  pair,
	}
	if orderItem.OrderType == TypeBuy {
		tempTrade.buyOrderID = order.ID
		tempTrade.amount = order.FinalDemandedAmount.String
	} else {
		tempTrade.sellOrderID = order.ID
		tempTrade.amount = order.FinalPayedByAmount.String
	}

	err = ps.handleTrade(tx, tempTrade, pair, orderItem.UserEmail, orderItem.UserID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("HandleExternalTradedOrder: handle trade: %w", err)
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		ps.logger.Error2("error handling OrderGroup", err,
			zap.String("service", "postOrderMatchingService"),
			zap.String("method", "HandleExternalTradedOrder"),
			zap.Int64("orderID", data.OrderID),
		)
	}
	if err != nil {
		return fmt.Errorf("HandleExternalTradedOrder: commit tx: %w", err)
	}
	return nil
}
