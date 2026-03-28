package order

import (
	"exchange-go/internal/platform"
	"exchange-go/internal/response"
	"exchange-go/internal/user"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

func (s *service) GetOrdersHistory(user *user.User, params GetOrdersHistoryParams, isFullHistory bool) (apiResponse response.APIResponse, statusCode int) {
	//we should create filter for repository query
	ordersHistory := make([]HistoryResponse, 0)

	filters := s.getFiltersForOrdersHistory(params, isFullHistory)
	filters.UserID = user.ID
	fields := s.orderRepository.GetLeafOrders(filters)

	var allOrderIds []int64
	var childPaths [][]int64
	for _, field := range fields {
		var childPath []int64
		orderIds := strings.Split(field.Path, ",")
		for _, orderID := range orderIds {
			if orderID != "" {
				orderIDInt64, _ := strconv.ParseInt(orderID, 10, 64)
				childPath = append(childPath, orderIDInt64)
				allOrderIds = append(allOrderIds, orderIDInt64)
			}
		}
		childPaths = append(childPaths, childPath)
	}

	if len(allOrderIds) < 1 {
		return response.Success(ordersHistory, "")
	} else {
		orders := s.orderRepository.GetOrdersByIds(allOrderIds)
		ordersHistory = s.getHistoryResponse(childPaths, orders)
		return response.Success(ordersHistory, "")
	}
}

func (s *service) getHistoryResponse(childPaths [][]int64, orders []Order) []HistoryResponse {
	isInChildPath := func(parentIds []int64, orderID int64) bool {
		for _, parentID := range parentIds {
			if orderID == parentID {
				return true
			}
		}
		return false
	}

	var resp []HistoryResponse
	var orderGroups [][]Order
	for _, childPath := range childPaths {
		var group []Order
		for _, order := range orders {
			if isInChildPath(childPath, order.ID) {
				group = append(group, order)
			}
		}
		orderGroups = append(orderGroups, group)
	}

	for _, group := range orderGroups {
		sort.Slice(group, func(i, j int) bool {
			return group[i].ID < group[j].ID
		})
		mainOrder := group[0]
		childOrder := group[len(group)-1]
		mainType := MainTypeOrder
		if mainOrder.IsStopOrder() {
			mainType = MainTypeStopOrder
		}
		orderType := mainOrder.Type
		exchangeType := strings.ToLower(mainOrder.ExchangeType)
		if mainOrder.IsFastExchange {
			exchangeType = removeUnderline(strings.ToLower(ExchangeTypeFast))
		}

		coinNames := strings.Split(mainOrder.Pair.Name, "-")
		payedByCoin := coinNames[0]
		demandedCoin := coinNames[1]
		if orderType == TypeBuy {
			demandedCoin = coinNames[0]
			payedByCoin = coinNames[1]
		}

		averagePrice, executedPercent, total := s.getAveragePriceAndExecutedAndTotal(group)

		amountDecimal, _ := decimal.NewFromString(mainOrder.PayedByAmount.String)
		amount := amountDecimal.String()

		hr := HistoryResponse{
			MainType:         mainType,
			OrderType:        exchangeType,
			Pair:             mainOrder.Pair.Name,
			ID:               mainOrder.ID,
			Side:             strings.ToLower(orderType),
			Price:            mainOrder.Price.String,
			AveragePrice:     averagePrice,
			SubUnit:          8,
			Amount:           amount + " " + payedByCoin,
			Total:            total + " " + demandedCoin,
			Executed:         executedPercent + " %",
			CreatedAt:        mainOrder.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:        childOrder.UpdatedAt.Format("2006-01-02 15:04:05"),
			Status:           strings.ToLower(childOrder.Status),
			TriggerCondition: s.getTriggerConditionForOrder(mainOrder),
		}

		resp = append(resp, hr)
	}

	return resp
}

func (s *service) getAveragePriceAndExecutedAndTotal(group []Order) (averagePrice string, executedPercent string, total string) {
	//average price is sum of tradePrice * amount divided by sum of amount (amount is finalPayed for sell and finalDemanded for buy orders)
	mainOrder := group[0]
	orderType := mainOrder.Type
	executedDecimal := decimal.NewFromFloat(0)
	sumOfFinalAmountMultiplyPriceDecimal := decimal.NewFromFloat(0)
	sumOfFinalAmountDecimal := decimal.NewFromFloat(0)
	for _, o := range group {
		finalPayedByAmountDecimal := decimal.NewFromFloat(0)
		if o.FinalPayedByAmount.Valid {
			finalPayedByAmountDecimal, _ = decimal.NewFromString(o.FinalPayedByAmount.String)
		}

		finalDemandedAmountDecimal := decimal.NewFromFloat(0)
		if o.FinalDemandedAmount.Valid {
			finalDemandedAmountDecimal, _ = decimal.NewFromString(o.FinalDemandedAmount.String)
		}

		tradePriceDecimal, _ := decimal.NewFromString(o.TradePrice.String)
		//here it means order is sell
		if orderType == TypeBuy {
			sumOfFinalAmountMultiplyPriceDecimal = sumOfFinalAmountMultiplyPriceDecimal.Add(finalDemandedAmountDecimal.Mul(tradePriceDecimal))
			sumOfFinalAmountDecimal = sumOfFinalAmountDecimal.Add(finalDemandedAmountDecimal)
			executedDecimal = executedDecimal.Add(finalDemandedAmountDecimal)
		} else {
			sumOfFinalAmountMultiplyPriceDecimal = sumOfFinalAmountMultiplyPriceDecimal.Add(finalPayedByAmountDecimal.Mul(tradePriceDecimal))
			sumOfFinalAmountDecimal = sumOfFinalAmountDecimal.Add(finalPayedByAmountDecimal)
			executedDecimal = executedDecimal.Add(finalPayedByAmountDecimal)
		}

	}

	averagePriceDecimal := decimal.NewFromFloat(0)
	if sumOfFinalAmountDecimal.IsPositive() {
		averagePriceDecimal = sumOfFinalAmountMultiplyPriceDecimal.Div(sumOfFinalAmountDecimal)
		averagePrice = averagePriceDecimal.StringFixed(8)
	}

	executedPercent = "0"
	if orderType == TypeBuy {
		demandedDecimal := decimal.NewFromFloat(0)
		if mainOrder.DemandedAmount.Valid {
			demandedDecimal, _ = decimal.NewFromString(mainOrder.DemandedAmount.String)
		}
		if demandedDecimal.IsPositive() {
			executedPercent = executedDecimal.Div(demandedDecimal).Mul(decimal.NewFromInt(100)).StringFixed(2)
		}
		total = executedDecimal.String()
	} else {
		payedByDecimal := decimal.NewFromFloat(0)
		if mainOrder.PayedByAmount.Valid {
			payedByDecimal, _ = decimal.NewFromString(mainOrder.PayedByAmount.String)
		}
		if payedByDecimal.IsPositive() {
			executedPercent = executedDecimal.Div(payedByDecimal).Mul(decimal.NewFromInt(100)).StringFixed(2)
		}
		total = executedDecimal.Mul(averagePriceDecimal).String()
	}

	return averagePrice, executedPercent, total
}

func (s *service) getFiltersForOrdersHistory(params GetOrdersHistoryParams, isFullHistory bool) HistoryFilters {
	var filters HistoryFilters
	if params.PairID > 0 {
		filters.PairID = params.PairID
	} else if params.PairName != "" {
		filters.PairName = params.PairName
	}

	startDate := ""
	endDate := ""

	orderType := strings.ToUpper(params.Type)
	if orderType == TypeBuy || orderType == TypeSell {
		filters.Type = orderType
	}

	if isFullHistory {

		filters.IsFullHistory = true

		filters.LastID = params.LastID
		filters.Hide = params.Hide

		if params.Period != "" {
			period := strings.ToUpper(strings.Trim(params.Period, ""))
			if s.isValidPeriod(period) {
				startDate, endDate = s.getStartAndEndDateFromPeriod(period)
			}
		}

		_, err := time.Parse("2006-01-02 15:04:05", params.StartDate)
		if err == nil {
			startDate = params.StartDate
		}

		_, err = time.Parse("2006-01-02 15:04:05", params.EndDate)
		if err == nil {
			endDate = params.EndDate
		}

	} else {
		if params.Period != "" {
			period := strings.ToUpper(strings.Trim(params.Period, ""))
			if s.isValidPeriod(period) {
				startDate, endDate = s.getStartAndEndDateFromPeriod(period)
			}
		}

	}

	filters.StartDate = startDate
	filters.EndDate = endDate

	if params.IsFastExchange != nil {
		filters.IsFastExchange = params.IsFastExchange
	}

	if s.configs.GetEnv() == platform.EnvTest {
		filters.PageSize = 3
	}

	return filters
}
