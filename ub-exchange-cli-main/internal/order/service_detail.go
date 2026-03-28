package order

import (
	"errors"
	"exchange-go/internal/response"
	"exchange-go/internal/user"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (s *service) FulfillOrder(adminUser *user.User, params FulfillOrderParams) (apiResponse response.APIResponse, statusCode int) {
	order := &Order{}
	err := s.orderRepository.GetOrderByID(params.ID, order)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("error in finding order", err,
			zap.String("service", "orderService"),
			zap.String("method", "TryToFulfillOrder"),
			zap.Int64("orderID", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || order.ID == 0 {
		return response.Error("order not found", http.StatusUnprocessableEntity, nil)
	}

	if order.Status != StatusOpen {
		return response.Error("order status is not open", http.StatusUnprocessableEntity, nil)
	}

	err = s.adminOrderManager.TryToFulfillOrder(*order)
	if err != nil {
		return response.Error(err.Error(), http.StatusUnprocessableEntity, nil)
	}

	return response.Success(nil, "")
}

func (s *service) GetOrderDetail(u *user.User, params GetOrderDetailParams) (apiResponse response.APIResponse, statusCode int) {
	details := make([]DetailResponse, 0)
	orders := s.orderRepository.GetUserOrderDetailsByID(params.ID, u.ID)
	//create response and
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

		d := DetailResponse{
			CreatedAt: o.CreatedAt.Format("2006-01-02 15:04:05"),
			Pair:      o.Pair.Name,
			Type:      strings.ToLower(o.Type),
			SubUnit:   8,
			Price:     o.TradePrice.String,
			Executed:  executed + " " + executedCoin,
			Fee:       fee + " " + demandedCoin,
			Amount:    amount + " " + basisCoin,
		}
		details = append(details, d)
	}
	return response.Success(details, "")
}
