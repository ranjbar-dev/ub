package order

import (
	"exchange-go/internal/response"
	"exchange-go/internal/user"
	"sort"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

func (s *service) GetOpenOrders(user *user.User, params GetOpenOrdersParams) (apiResponse response.APIResponse, statusCode int) {
	pairID := params.PairID
	leafOrders := s.orderRepository.GetUserOpenOrders(user.ID, pairID)

	var allAncestorsIds []int64
	for _, leafOrder := range leafOrders {
		path := leafOrder.Path.String
		parentOrderIds := strings.Split(path, ",")
		for _, parentOrderID := range parentOrderIds {
			parentOrderIDInt64, _ := strconv.ParseInt(parentOrderID, 10, 64)
			if parentOrderIDInt64 != leafOrder.ID {
				allAncestorsIds = append(allAncestorsIds, parentOrderIDInt64)
			}
		}

	}

	var allAncestorOrders []Order
	if len(allAncestorsIds) > 0 {
		allAncestorOrders = s.orderRepository.GetOrdersAncestors(allAncestorsIds)
	}
	sort.Slice(leafOrders, func(i, j int) bool {
		return leafOrders[i].CreatedAt.Sub(leafOrders[j].CreatedAt) > 0
	})

	openOrders := make([]OpenOrdersResponse, 0)

	for _, o := range leafOrders {
		ancestorIds := strings.Split(o.Path.String, ",")
		ancestorsLen := len(ancestorIds)
		if ancestorsLen > 2 {
			//this if means the order has parents and we try to aggregate all of them
			var result OpenOrdersResponse
			executedAmountDecimal := decimal.NewFromFloat(0)
			firstAncestorID, _ := strconv.ParseInt(ancestorIds[0], 10, 64)
			i := 0
			for _, ancestorOrder := range allAncestorOrders {
				if s.isInAncestors(ancestorOrder, ancestorIds) {
					i++
					if ancestorOrder.ID == firstAncestorID {
						result = s.getSingleOpenOrderData(ancestorOrder)
					}

					executedAmount := "0"
					if ancestorOrder.Type == TypeBuy {
						executedAmount = ancestorOrder.FinalPayedByAmount.String

					} else {
						executedAmount = ancestorOrder.FinalDemandedAmount.String
					}

					finalDecimal, _ := decimal.NewFromString(executedAmount)
					executedAmountDecimal = executedAmountDecimal.Add(finalDecimal)
					if i == ancestorsLen-1 {
						break
					}
				}
			}

			executedPercentageDecimal := decimal.NewFromFloat(0)
			if !result.total.IsZero() {
				executedPercentageDecimal = executedAmountDecimal.Div(result.total).Mul(decimal.NewFromInt(100))
			}

			result.ID = o.ID //we substitute the ancestor order id with open child order so the client can cancel the last open child by sending the id
			result.Executed = executedPercentageDecimal.StringFixed(2) + " %"
			openOrders = append(openOrders, result)

		} else {
			//here it means the order is single
			result := s.getSingleOpenOrderData(o)
			openOrders = append(openOrders, result)
		}
	}

	sort.Slice(openOrders, func(i, j int) bool {
		return openOrders[i].createdAt.Sub(openOrders[j].createdAt) > 0
	})

	return response.Success(openOrders, "")
}

func (s *service) isInAncestors(o Order, ancestorIds []string) bool {
	for _, id := range ancestorIds {
		idInt64, _ := strconv.ParseInt(id, 10, 64)
		if idInt64 == o.ID {
			return true
		}
	}
	return false
}

func (s *service) getSingleOpenOrderData(o Order) OpenOrdersResponse {
	amount := o.getAmount()
	totalAmount := decimal.NewFromFloat(0)
	price := o.Price.String
	if !o.isMarket() {
		priceDecimal, _ := decimal.NewFromString(o.Price.String)
		amountDecimal, _ := decimal.NewFromString(amount)
		if o.Type == TypeBuy {
			totalAmount, _ = decimal.NewFromString(o.PayedByAmount.String)
		} else {
			totalAmount = amountDecimal.Mul(priceDecimal)
		}
	} else {
		totalAmount, _ = decimal.NewFromString(o.PayedByAmount.String)
		amount = ""
	}

	//here it means the order is single
	mainType := MainTypeOrder
	if o.IsStopOrder() {
		mainType = MainTypeStopOrder
	}
	oor := OpenOrdersResponse{
		MainType:         mainType,
		OrderType:        strings.ToLower(o.ExchangeType),
		Pair:             o.Pair.Name,
		ID:               o.ID,
		Side:             strings.ToLower(o.Type),
		Price:            price,
		SubUnit:          8,
		Amount:           amount,
		Total:            totalAmount.StringFixed(8),
		total:            totalAmount,
		Executed:         "0.00 %", //since there is no parent and order is still open
		CreatedAt:        o.CreatedAt.Format("2006-01-02 15:04:05"),
		createdAt:        o.CreatedAt,
		TriggerCondition: s.getTriggerConditionForOrder(o),
	}

	return oor
}

func (s *service) getTriggerConditionForOrder(o Order) string {
	if !o.IsStopOrder() {
		return ""
	}
	if o.Type == TypeBuy {
		return ">= " + o.StopPointPrice.String
	}
	return "<= " + o.StopPointPrice.String
}
