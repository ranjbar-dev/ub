package order

import (
	"exchange-go/internal/platform"
	"exchange-go/internal/response"
	"exchange-go/internal/user"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

func (s *service) GetTradesHistory(user *user.User, params GetTradesHistoryParams, isFullHistory bool) (apiResponse response.APIResponse, statusCode int) {
	tradesHistory := make([]TradeHistoryResponse, 0)
	filters := s.getFiltersForTradesHistory(params, isFullHistory)
	filters.UserID = user.ID
	orders := s.orderRepository.GetUserTradedOrders(filters)

	for _, o := range orders {
		orderType := o.Type
		coinNames := strings.Split(o.Pair.Name, "-")
		basisCoin := coinNames[1]
		dependentCoin := coinNames[0]

		executedCoin := dependentCoin
		demandedCoin := basisCoin

		finalPayedByDecimal, _ := decimal.NewFromString(o.FinalPayedByAmount.String)
		finalDemandedDecimal, _ := decimal.NewFromString(o.FinalDemandedAmount.String)
		feeDecimal := decimal.NewFromFloat(o.FeePercentage.Float64)
		fee := finalDemandedDecimal.Mul(feeDecimal).String()

		amount := ""
		executed := ""

		if orderType == TypeBuy {
			demandedCoin = dependentCoin
			executedCoin = demandedCoin
			amount = finalPayedByDecimal.String()
			executed = finalDemandedDecimal.String()
		} else {
			amount = finalDemandedDecimal.String()
			executed = finalPayedByDecimal.String()
		}

		th := TradeHistoryResponse{
			ID:        o.ID,
			CreatedAt: o.CreatedAt.Format("2006-01-02 15:04:05"),
			Pair:      o.Pair.Name,
			OrderType: strings.ToLower(o.Type),
			Price:     o.TradePrice.String,
			SubUnit:   8,
			Executed:  executed + " " + executedCoin,
			Fee:       fee + " " + demandedCoin,
			Amount:    amount + " " + basisCoin,
			Total:     amount + " " + basisCoin,
		}

		tradesHistory = append(tradesHistory, th)
	}

	return response.Success(tradesHistory, "")
}

func (s *service) getFiltersForTradesHistory(params GetTradesHistoryParams, isFullHistory bool) TradeHistoryFilters {
	var filters TradeHistoryFilters
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

	if s.configs.GetEnv() == platform.EnvTest {
		filters.PageSize = 3
	}

	return filters
}

func (s *service) getStartAndEndDateFromPeriod(period string) (string, string) {

	startDate := ""
	endDate := ""
	switch period {
	case HistoryPeriod1Day:
		t := time.Now()
		startDate = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, t.Nanosecond(), t.Location()).Format("2006-01-02 15:04:05")
		endDate = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, t.Nanosecond(), t.Location()).Format("2006-01-02 15:04:05")
		return startDate, endDate
	case HistoryPeriod1Week:
		t := time.Now()
		t2 := time.Now()
		t = t.Add(-1 * 7 * 24 * time.Hour)
		startDate = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, t.Nanosecond(), t.Location()).Format("2006-01-02 15:04:05")
		endDate = time.Date(t2.Year(), t2.Month(), t2.Day(), 23, 59, 59, t2.Nanosecond(), t2.Location()).Format("2006-01-02 15:04:05")
		return startDate, endDate
	case HistoryPeriod1Month:
		t := time.Now()
		t2 := time.Now()
		t = t.Add(-1 * 30 * 24 * time.Hour)
		startDate = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, t.Nanosecond(), t.Location()).Format("2006-01-02 15:04:05")
		endDate = time.Date(t2.Year(), t2.Month(), t2.Day(), 23, 59, 59, t.Nanosecond(), t.Location()).Format("2006-01-02 15:04:05")
		return startDate, endDate
	case HistoryPeriod3Month:
		t := time.Now()
		t2 := time.Now()
		t = t.Add(-1 * 90 * 24 * time.Hour)
		startDate = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, t.Nanosecond(), t.Location()).Format("2006-01-02 15:04:05")
		endDate = time.Date(t2.Year(), t2.Month(), t2.Day(), 23, 59, 59, t.Nanosecond(), t.Location()).Format("2006-01-02 15:04:05")
	default:

	}

	return startDate, endDate
}

func (s *service) isValidPeriod(period string) bool {
	switch period {
	case HistoryPeriod1Day, HistoryPeriod1Week, HistoryPeriod1Month, HistoryPeriod3Month:
		return true
	default:
		return false
	}
}
