package order

import (
	"context"
	"errors"
	"exchange-go/internal/currency"
	"exchange-go/internal/response"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"net/http"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *service) CancelOrder(user *user.User, params CancelOrderParams) (apiResponse response.APIResponse, statusCode int) {
	ctx := context.Background()
	tx := s.db.Begin()
	err := tx.Error
	if err != nil {
		s.logger.Error2("error in starting transaction", err,
			zap.String("service", "orderService"),
			zap.String("method", "CancelOrder"),
			zap.Int64("orderID", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	o := &Order{}
	err = s.orderRepository.GetOrderByIDUsingTx(tx, params.ID, o)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		s.logger.Error2("error in finding order", err,
			zap.String("service", "orderService"),
			zap.String("method", "CancelOrder"),
			zap.Int64("orderID", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || o.ID == 0 || o.UserID != user.ID {
		tx.Rollback()
		return response.Error("order not found", http.StatusUnprocessableEntity, nil)
	}

	if o.Status != StatusOpen {
		tx.Rollback()
		return response.Error("order status is not open", http.StatusUnprocessableEntity, nil)
	}

	pair, err := s.currencyService.GetPairByID(o.PairID)
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not get pair from db", err,
			zap.String("service", "orderService"),
			zap.String("method", "CancelOrder"),
			zap.Int64("orderID", o.ID),
			zap.Int64("pairID", o.PairID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	isActionAllowed := IsActionAllowed(pair, currency.ActionCancel)
	if !o.IsStopOrder() && !isActionAllowed {
		tx.Rollback()
		return response.Error("this action is not allowed at the moment", http.StatusUnprocessableEntity, nil)
	}

	payedByCoinID := pair.DependentCoinID
	if o.Type == TypeBuy {
		payedByCoinID = pair.BasisCoinID
	}

	ub := &userbalance.UserBalance{}
	err = s.userBalanceService.GetBalanceOfUserByCoinUsingTx(tx, user.ID, payedByCoinID, ub)
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error in getting user balance", err,
			zap.String("service", "orderService"),
			zap.String("method", "CancelOrder"),
			zap.Int64("orderID", o.ID),
			zap.Int("userID", user.ID),
			zap.Int64("coinID", payedByCoinID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	o.Status = StatusCanceled
	err = tx.Omit(clause.Associations).Save(o).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error in saving order", err,
			zap.String("service", "orderService"),
			zap.String("method", "CancelOrder"),
			zap.Int64("orderID", o.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	formerFrozenAmountDecimal, err := decimal.NewFromString(ub.FrozenAmount)
	payedByAmountDecimal, err := decimal.NewFromString(o.PayedByAmount.String)

	newFrozenBalance := formerFrozenAmountDecimal.Sub(payedByAmountDecimal)
	if newFrozenBalance.IsNegative() {
		tx.Rollback()
		s.logger.Error2("fronzen balance is negative", err,
			zap.String("service", "orderService"),
			zap.String("method", "CancelOrder"),
			zap.Int64("orderID", o.ID),
			zap.Int("userID", user.ID),
			zap.Int64("coinID", payedByCoinID),
			zap.String("payedBy", o.PayedByAmount.String),
			zap.String("frozenAmount", ub.FrozenAmount),
			zap.Bool("isStopOrder", o.IsStopOrder()),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		//newFrozenBalance = decimal.NewFromFloat(0.0)
	}
	newFrozenBalanceString := newFrozenBalance.StringFixed(8)
	ub.FrozenAmount = newFrozenBalanceString
	err = tx.Omit(clause.Associations).Save(ub).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error in saving userBalance", err,
			zap.String("service", "orderService"),
			zap.String("method", "CancelOrder"),
			zap.Int64("orderID", o.ID),
			zap.Int("userID", user.ID),
			zap.Int64("coinID", payedByCoinID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	//we need this in other flows
	o.Pair = pair
	o.User = *user

	if o.IsStopOrder() && !o.IsSubmitted.Bool {
		err := s.orderRedisManager.RemoveStopOrderFromQueue(ctx, *o, pair.Name)
		//if errors happen in removing stop order from redis we continue the process since the stop order will
		// eventually removed in stopOrderSubmissionManager
		if err != nil {
			s.logger.Error2("error in removing stopOrder from redis", err,
				zap.String("service", "orderService"),
				zap.String("method", "CancelOrder"),
				zap.Int64("orderID", o.ID),
			)
		}
	} else {
		//this one should sent to engine to remove the order from orderbook
		err := s.engineCommunicator.RemoveOrder(*o)
		if err != nil {
			tx.Rollback()
			s.logger.Error2("error in removing order from orderbook", err,
				zap.String("service", "orderService"),
				zap.String("method", "CancelOrder"),
				zap.Int64("orderID", o.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error in commiting transaction", err,
			zap.String("service", "orderService"),
			zap.String("method", "CancelOrder"),
			zap.Int64("orderID", o.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	s.eventsHandler.HandleOrderCancellation(*o)

	return response.Success(nil, "")
}
