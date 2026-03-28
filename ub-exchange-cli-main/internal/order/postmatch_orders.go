package order

import (
	"database/sql"
	"strconv"
	"time"

	"exchange-go/internal/currency"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (ps *postOrderMatchingService) updateMainOrder(tx *gorm.DB, orderItem MatchingNeededQueryFields, tempOrder tempOrder, pair currency.Pair, isMarket bool) (Order, error) {
	finalDemandedAmountDecimal, finalPayedByAmountDecimal := ps.getFinalAmountsDecimal(tempOrder)
	//demandedAmountDecimal, _ := decimal.NewFromString(orderItem.DemandedAmount)
	//payedByAmountDecimal, _ := decimal.NewFromString(orderItem.PayedByAmount)

	//these two if is for handling very small diferrences in amount*price with payedBy and demanded
	//if finalDemandedAmountDecimal.GreaterThan(demandedAmountDecimal) {
	//	finalDemandedAmountDecimal = demandedAmountDecimal
	//}
	//
	//if finalPayedByAmountDecimal.GreaterThan(payedByAmountDecimal) {
	//	finalPayedByAmountDecimal = payedByAmountDecimal
	//}

	finalDemandedAmount := finalDemandedAmountDecimal.StringFixed(8)
	finalPayedByAmount := finalPayedByAmountDecimal.StringFixed(8)
	feePercentage := ps.GetFeePercentage(orderItem, pair, isMarket, tempOrder.isMaker)

	updatingOrder := &Order{
		ID:                  orderItem.OrderID,
		IsMaker:             sql.NullBool{Bool: tempOrder.isMaker, Valid: true},
		FeePercentage:       sql.NullFloat64{Float64: feePercentage, Valid: true},
		TradePrice:          sql.NullString{String: tempOrder.tradePrice, Valid: true},
		Status:              StatusFilled,
		FinalDemandedAmount: sql.NullString{String: finalDemandedAmount, Valid: true},
		FinalPayedByAmount:  sql.NullString{String: finalPayedByAmount, Valid: true},
		IsTradedWithBot:     sql.NullBool{Bool: false, Valid: true},
	}
	err := tx.Model(updatingOrder).Updates(updatingOrder).Error

	//we need this field for child orders
	updatingOrder.Path = sql.NullString{String: orderItem.Path, Valid: true}
	updatingOrder.UserID = orderItem.UserID
	updatingOrder.Type = orderItem.OrderType
	updatingOrder.DemandedAmount = sql.NullString{String: orderItem.DemandedAmount, Valid: true}
	updatingOrder.PayedByAmount = sql.NullString{String: orderItem.PayedByAmount, Valid: true}

	if orderItem.OrderType == TypeBuy {
		updatingOrder.DemandedCoin = pair.DependentCoin.Code
		updatingOrder.PayedByCoin = pair.BasisCoin.Code
	} else {
		updatingOrder.DemandedCoin = pair.BasisCoin.Code
		updatingOrder.PayedByCoin = pair.DependentCoin.Code
	}

	if orderItem.Price != "" {
		updatingOrder.Price = sql.NullString{String: orderItem.Price, Valid: true}
	}

	return *updatingOrder, err
}

func (ps *postOrderMatchingService) createChildOrder(tx *gorm.DB, orderItem MatchingNeededQueryFields, tempOrder tempOrder, parentOrder Order, pair currency.Pair, status string) (Order, error) {
	finalDemandedAmountDecimal, finalPayedByAmountDecimal := ps.getFinalAmountsDecimal(tempOrder)
	finalDemandedAmount := finalDemandedAmountDecimal.StringFixed(8)
	finalPayedByAmount := finalPayedByAmountDecimal.StringFixed(8)
	feePercentage := ps.GetFeePercentage(orderItem, pair, parentOrder.isMarket(), tempOrder.isMaker)

	extraInfo := &ExtraInfo{
		UserAgentInfo: sql.NullString{String: orderItem.UserAgentInfo, Valid: true},
	}

	err := tx.Omit(clause.Associations).Create(&extraInfo).Error //create orderExtraInfo
	if err != nil {
		return Order{}, err
	}

	parentOrderDemandedAmountDecimal, err := decimal.NewFromString(parentOrder.DemandedAmount.String)
	if err != nil {
		return Order{}, err
	}

	parentOrderPayedByAmountDecimal, err := decimal.NewFromString(parentOrder.PayedByAmount.String)
	if err != nil {
		return Order{}, err
	}

	parentOrderFinalDemandedAmountDecimal, err := decimal.NewFromString(parentOrder.FinalDemandedAmount.String)
	if err != nil {
		return Order{}, err
	}

	parentOrderFinalPayedByAmountDecimal, err := decimal.NewFromString(parentOrder.FinalPayedByAmount.String)
	if err != nil {
		return Order{}, err
	}

	demandedAmount := ""
	payedByAmount := ""
	//in case that parent order price is different from child price
	// this calculation handles the right demanded and payedby for child
	if parentOrder.Price.String != "" {
		if tempOrder.orderType == TypeBuy {
			demandedAmount = parentOrderDemandedAmountDecimal.Sub(parentOrderFinalDemandedAmountDecimal).StringFixed(8)
			parentTradePriceDecimal, _ := decimal.NewFromString(parentOrder.TradePrice.String)
			parentTradedAmountDecimal := parentOrderFinalPayedByAmountDecimal.Div(parentTradePriceDecimal)
			parentPriceDecimal, _ := decimal.NewFromString(parentOrder.Price.String)
			payedByAmount = parentOrderPayedByAmountDecimal.Sub(parentTradedAmountDecimal.Mul(parentPriceDecimal)).StringFixed(8)
		} else {
			parentTradePriceDecimal, _ := decimal.NewFromString(parentOrder.TradePrice.String)
			parentTradedAmountDecimal := parentOrderFinalDemandedAmountDecimal.Div(parentTradePriceDecimal)
			parentPriceDecimal, _ := decimal.NewFromString(parentOrder.Price.String)
			demandedAmount = parentOrderDemandedAmountDecimal.Sub(parentTradedAmountDecimal.Mul(parentPriceDecimal)).StringFixed(8)
			payedByAmount = parentOrderPayedByAmountDecimal.Sub(parentOrderFinalPayedByAmountDecimal).StringFixed(8)
		}
	} else {
		payedByAmountDecimal := parentOrderPayedByAmountDecimal.Sub(parentOrderFinalPayedByAmountDecimal)
		payedByAmount = parentOrderPayedByAmountDecimal.Sub(parentOrderFinalPayedByAmountDecimal).StringFixed(8)
		marketPriceDecimal, _ := decimal.NewFromString(tempOrder.marketPrice)
		if tempOrder.orderType == TypeBuy {
			demandedAmount = payedByAmountDecimal.Div(marketPriceDecimal).StringFixed(8)
		} else {
			demandedAmount = payedByAmountDecimal.Mul(marketPriceDecimal).StringFixed(8)
		}
	}

	finalDemandedAmountNil := sql.NullString{String: "", Valid: false}
	finalPayedByAmountNil := sql.NullString{String: "", Valid: false}
	tradePriceNil := sql.NullString{String: "", Valid: false}
	isMakerNil := sql.NullBool{Bool: false, Valid: false}
	feePercentageNil := sql.NullFloat64{Float64: 0.0, Valid: false}
	isTradedWithBotNil := sql.NullBool{Bool: false, Valid: false}
	if status == StatusFilled {
		if tempOrder.tradeAmount != "" {
			finalDemandedAmountNil = sql.NullString{String: finalDemandedAmount, Valid: true}
			finalPayedByAmountNil = sql.NullString{String: finalPayedByAmount, Valid: true}
			isTradedWithBotNil = sql.NullBool{Bool: false, Valid: true}
		} else {
			//this is for partial market orders
			finalDemandedAmountNil = sql.NullString{String: demandedAmount, Valid: true}
			finalPayedByAmountNil = sql.NullString{String: payedByAmount, Valid: true}
			isTradedWithBotNil = sql.NullBool{Bool: true, Valid: true}
		}
		tradePriceNil = sql.NullString{String: tempOrder.tradePrice, Valid: true}
		isMakerNil = sql.NullBool{Bool: tempOrder.isMaker, Valid: true}
		feePercentageNil = sql.NullFloat64{Float64: feePercentage, Valid: true}

	}

	childLevel := parentOrder.Level.Int64 + 1

	currentMarketPriceNil := sql.NullString{String: currentMarketPrice, Valid: true}

	childOrder := &Order{
		UserID:              parentOrder.UserID,
		ParentID:            sql.NullInt64{Int64: parentOrder.ID, Valid: true},
		Type:                tempOrder.orderType,
		ExchangeType:        orderItem.OrderExchangeType,
		Price:               parentOrder.Price,
		Status:              status,
		DemandedAmount:      sql.NullString{String: demandedAmount, Valid: true},
		DemandedCoin:        parentOrder.DemandedCoin,
		PayedByAmount:       sql.NullString{String: payedByAmount, Valid: true},
		PayedByCoin:         parentOrder.PayedByCoin,
		PairID:              pair.ID,
		ExtraInfoID:         sql.NullInt64{Int64: extraInfo.ID, Valid: true},
		TradePrice:          tradePriceNil,
		IsMaker:             isMakerNil,
		FeePercentage:       feePercentageNil,
		Level:               sql.NullInt64{Int64: childLevel, Valid: true},
		FinalDemandedAmount: finalDemandedAmountNil,
		FinalPayedByAmount:  finalPayedByAmountNil,
		IsTradedWithBot:     isTradedWithBotNil,
		CurrentMarketPrice:  currentMarketPriceNil,
	}

	err = tx.Omit(clause.Associations).Create(childOrder).Error //create order
	if err != nil {
		return *childOrder, err
	}
	childPath := parentOrder.Path.String + strconv.FormatInt(childOrder.ID, 10) + ","
	childOrder.Path = sql.NullString{String: childPath, Valid: true}
	err = tx.Omit(clause.Associations).Save(childOrder).Error
	return *childOrder, err

}

func (ps *postOrderMatchingService) handleTrade(tx *gorm.DB, tempTrade tempTrade, pair currency.Pair, userEmail string, userID int) error {
	if !ps.shouldCreateTrade(tempTrades, tempTrade.sellOrderID, tempTrade.buyOrderID) {
		return nil
	}

	buyOrderIDNil := sql.NullInt64{Int64: 0, Valid: false}
	sellOrderIDNil := sql.NullInt64{Int64: 0, Valid: false}
	botOrderTypeNil := sql.NullString{String: "", Valid: false}
	if tempTrade.sellOrderID == 0 {
		botOrderTypeNil = sql.NullString{String: TypeSell, Valid: true}
	} else {
		sellOrderIDNil = sql.NullInt64{Int64: tempTrade.sellOrderID, Valid: true}
	}

	if tempTrade.buyOrderID == 0 {
		botOrderTypeNil = sql.NullString{String: TypeBuy, Valid: true}
	} else {
		buyOrderIDNil = sql.NullInt64{Int64: tempTrade.buyOrderID, Valid: true}
	}

	priceDecimal, _ := decimal.NewFromString(tempTrade.price)
	amountDecimal, _ := decimal.NewFromString(tempTrade.amount)
	price := priceDecimal.StringFixed(8)
	amount := amountDecimal.StringFixed(8)
	trade := &Trade{
		Price:        sql.NullString{String: price, Valid: true},
		Amount:       sql.NullString{String: amount, Valid: true},
		PairID:       tempTrade.pair.ID,
		BuyOrderID:   buyOrderIDNil,
		SellOrderID:  sellOrderIDNil,
		BotOrderType: botOrderTypeNil,
	}

	err := tx.Omit(clause.Associations).Create(trade).Error
	if err != nil {
		return err
	}
	tempTrades = append(tempTrades, tempTrade)
	tradesData = append(tradesData, TradeData{
		Trade:     *trade,
		UserEmail: userEmail,
		UserID:    userID},
	)
	return nil
}

func (ps *postOrderMatchingService) prepareResult(
	doneOrders []CallBackOrderData,
	partial *CallBackOrderData,
	isPartialOpen bool,
	removingDoneOrderIds []int64,
	err error,
	remainingPartial *CallBackOrderData,
) MatchingResult {
	if err == nil {
		return MatchingResult{
			Err:                   nil,
			RemainingPartialOrder: remainingPartial,
			RemovingDoneOrderIds:  removingDoneOrderIds,
		}
	}
	result := MatchingResult{}
	result.Err = err
	result.RemovingDoneOrderIds = removingDoneOrderIds

	//becase the matching had error we check if partial order is independent order or part of one of done orders
	//if it is independent and limit order we push that to orderbook and if it is market order we save it in
	//redis to handle it later
	if partial != nil {
		if partial.Price != "" {
			if isPartialOpen && !partial.IsAlreadyInOrderBook {
				//the order is limit we push it to orderbook
				result.RemainingPartialOrder = partial
			}
		} else {
			//the order is market we save it in redis and handle it later
			go ps.storeUnmatchedOrderID(partial.ID)
		}

	}

	return result
}

func (ps *postOrderMatchingService) handlePartialOrder(tx *gorm.DB, orderItem MatchingNeededQueryFields, partial tempOrder, parentOrder *Order, pair currency.Pair, isFromAdmin bool) partialOrderHandlingResult {
	shouldForceTrade := true
	var err error
	if !isFromAdmin {
		shouldForceTrade, err = ps.forceTrader.ShouldForceTrade(pair.Name, orderItem.OrderType, orderItem.Price)
		if err != nil {
			return partialOrderHandlingResult{
				err:      err,
				isTraded: false,
			}
		}
	}

	if !shouldForceTrade {
		if parentOrder != nil {
			childOrder, err := ps.createChildOrder(tx, orderItem, partial, *parentOrder, pair, StatusOpen)
			if err != nil {
				return partialOrderHandlingResult{
					err:      err,
					isTraded: false,
					order:    childOrder,
				}
			}

			return partialOrderHandlingResult{
				err:      nil,
				isTraded: false,
				order:    childOrder,
			}
		} else {
			createdAt, _ := time.Parse("2006-01-02T15:04:05Z", orderItem.CreatedAt)
			order := Order{
				ID:             orderItem.OrderID,
				Type:           orderItem.OrderType,
				Price:          sql.NullString{String: orderItem.Price, Valid: true},
				CreatedAt:      createdAt,
				DemandedAmount: sql.NullString{String: orderItem.DemandedAmount, Valid: true},
				PayedByAmount:  sql.NullString{String: orderItem.PayedByAmount, Valid: true},
			}
			return partialOrderHandlingResult{
				err:      nil,
				isTraded: false,
				order:    order,
			}
		}
	}

	isMarket := false
	tradePrice := orderItem.Price
	if isFromAdmin {
		tradePrice = partial.tradePrice
	}

	if tradePrice == "" {
		isMarket = true
		tradePrice = currentMarketPrice
		if err != nil {
			return partialOrderHandlingResult{
				err:      err,
				isTraded: false,
			}
		}
	}

	feePercentage := ps.GetFeePercentage(orderItem, pair, isMarket, false)
	if parentOrder == nil {
		order := &Order{
			ID:                  orderItem.OrderID,
			IsMaker:             sql.NullBool{Bool: false, Valid: true},
			FeePercentage:       sql.NullFloat64{Float64: feePercentage, Valid: true},
			TradePrice:          sql.NullString{String: tradePrice, Valid: true},
			Status:              StatusFilled,
			FinalDemandedAmount: sql.NullString{String: orderItem.DemandedAmount, Valid: true},
			FinalPayedByAmount:  sql.NullString{String: orderItem.PayedByAmount, Valid: true},
			IsTradedWithBot:     sql.NullBool{Bool: true, Valid: true},
		}
		err = tx.Model(order).Updates(order).Error

		//set extra data needed in our flow

		order.UserID = orderItem.UserID
		order.Type = orderItem.OrderType
		createdAt, _ := time.Parse("2006-01-02T15:04:05Z", orderItem.CreatedAt)
		order.CreatedAt = createdAt

		if err != nil {
			return partialOrderHandlingResult{
				err:      err,
				isTraded: false,
			}
		}

		return partialOrderHandlingResult{
			err:      nil,
			isTraded: true,
			order:    *order,
		}
	} else {
		partial.tradePrice = tradePrice
		childOrder, err := ps.createChildOrder(tx, orderItem, partial, *parentOrder, pair, StatusFilled)
		if err != nil {
			return partialOrderHandlingResult{
				err:      err,
				isTraded: false,
			}
		}
		return partialOrderHandlingResult{
			err:      nil,
			isTraded: true,
			order:    childOrder,
		}

	}

}
