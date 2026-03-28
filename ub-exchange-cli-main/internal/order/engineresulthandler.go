package order

import (
	"exchange-go/internal/engine"
	"fmt"
	"strconv"
)

// EngineResultHandler processes trade results returned by the matching engine
// after orders have been matched.
type EngineResultHandler interface {
	// CallBack receives the list of fully matched orders and an optional partially filled
	// order from the engine, delegates post-trade settlement, and returns the matching result.
	CallBack(doneOrders []engine.Order, partialOrder *engine.Order) engine.MatchingResult
}

type engineResultHandler struct {
	postOrderMatchingService PostOrderMatchingService
}

func (rh *engineResultHandler) CallBack(doneOrders []engine.Order, partialOrder *engine.Order) engine.MatchingResult {
	var result []CallBackOrderData
	for _, doneOrder := range doneOrders {
		id, _ := strconv.ParseInt(doneOrder.ID, 10, 64)
		tradedWithOrderID, _ := strconv.ParseInt(doneOrder.TradedWithOrderID, 10, 64)
		orderData := CallBackOrderData{
			ID:                id,
			PairName:          doneOrder.Pair,
			OrderType:         mapSideToType(doneOrder.Side),
			Quantity:          doneOrder.Quantity,
			Price:             doneOrder.Price,
			Timestamp:         doneOrder.Timestamp,
			TradedWithOrderID: tradedWithOrderID,
			QuantityTraded:    doneOrder.QuantityTraded,
			TradePrice:        doneOrder.TradePrice,
			MarketPrice:       doneOrder.MarketPrice,
			MinThresholdPrice: doneOrder.MinThresholdPrice,
			MaxThresholdPrice: doneOrder.MaxThresholdPrice,
		}
		result = append(result, orderData)
	}

	var partialCallBackOrder *CallBackOrderData
	if partialOrder != nil {
		partialOrderID, _ := strconv.ParseInt(partialOrder.ID, 10, 64)
		partialCallBackOrder = &CallBackOrderData{
			ID:                   partialOrderID,
			PairName:             partialOrder.Pair,
			OrderType:            mapSideToType(partialOrder.Side),
			Quantity:             partialOrder.Quantity,
			Price:                partialOrder.Price,
			Timestamp:            partialOrder.Timestamp,
			TradedWithOrderID:    0,
			QuantityTraded:       partialOrder.QuantityTraded,
			TradePrice:           partialOrder.TradePrice,
			MarketPrice:          partialOrder.MarketPrice,
			IsAlreadyInOrderBook: partialOrder.IsAlreadyInOrderBook,
			MinThresholdPrice:    partialOrder.MinThresholdPrice,
			MaxThresholdPrice:    partialOrder.MaxThresholdPrice,
		}

	}

	matchingResult := rh.postOrderMatchingService.HandlePostOrderMatching(result, partialCallBackOrder, false)
	var remainingPartialOrder *engine.Order
	if matchingResult.RemainingPartialOrder != nil {
		idString := fmt.Sprintf("%011s", strconv.FormatInt(matchingResult.RemainingPartialOrder.ID, 10))
		remainingPartialOrder = &engine.Order{
			Pair:              matchingResult.RemainingPartialOrder.PairName,
			ID:                idString,
			Side:              mapTypeToSide(matchingResult.RemainingPartialOrder.OrderType),
			Quantity:          matchingResult.RemainingPartialOrder.Quantity,
			Price:             matchingResult.RemainingPartialOrder.Price,
			Timestamp:         matchingResult.RemainingPartialOrder.Timestamp,
			TradedWithOrderID: "",
			QuantityTraded:    "",
			TradePrice:        matchingResult.RemainingPartialOrder.TradePrice,
			MarketPrice:       matchingResult.RemainingPartialOrder.MarketPrice,
			MinThresholdPrice: matchingResult.RemainingPartialOrder.MinThresholdPrice,
			MaxThresholdPrice: matchingResult.RemainingPartialOrder.MaxThresholdPrice,
		}
	}

	engineMatchingResult := engine.MatchingResult{
		Err:                   matchingResult.Err,
		RemainingPartialOrder: remainingPartialOrder,
		RemovingDoneOrderIds:  matchingResult.RemovingDoneOrderIds,
	}
	return engineMatchingResult
}

func NewEngineResultHandler(poms PostOrderMatchingService) EngineResultHandler {
	return &engineResultHandler{
		postOrderMatchingService: poms,
	}
}
