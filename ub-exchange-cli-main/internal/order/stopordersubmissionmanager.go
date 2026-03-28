package order

import (
	"context"
	"database/sql"
	"encoding/json"
	"exchange-go/internal/livedata"
	"exchange-go/internal/platform"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// StopOrderSubmissionManager submits triggered stop orders to the matching engine
// when the market price crosses their stop-point threshold.
type StopOrderSubmissionManager interface {
	// Submit retrieves all triggered stop orders for a pair based on the price movement
	// between formerPrice and price, marks them as submitted, and routes them to the engine.
	Submit(ctx context.Context, pairName string, price string, formerPrice string)
	// SubmitOrderInDb marks a stop order as submitted in the database and records
	// the current market price at the time of submission.
	SubmitOrderInDb(ctx context.Context, order *Order, price string) error
}

type stopOrderSubmissionManager struct {
	db                *gorm.DB
	orderRepository   Repository
	liveDataService   livedata.Service
	orderRedisManager RedisManager
	eventsHandler     EventsHandler
	logger            platform.Logger
}

func (s *stopOrderSubmissionManager) Submit(ctx context.Context, pairName string, price string, formerPrice string) {
	isRising, shouldCalculate, err := s.isPriceRising(price, formerPrice)
	if err != nil {
		s.logger.Error2("can not check if error is Rising or not", err,
			zap.String("service", "stopOrderSubmissionManager"),
			zap.String("method", "Submit"),
			zap.String("pairName", pairName),
			zap.String("price", price),
			zap.String("formerPrice", formerPrice),
		)
		return
	}

	if !shouldCalculate {
		return
	}

	orders, err := s.getStopOrders(ctx, pairName, formerPrice, price, isRising)
	if err != nil {
		s.logger.Error2("can not get stop orders", err,
			zap.String("service", "stopOrderSubmissionManager"),
			zap.String("method", "Submit"),
			zap.String("pairName", pairName),
			zap.String("price", price),
			zap.String("formerPrice", formerPrice),
		)
		return
	}

	for _, order := range orders {
		if order.Status == StatusOpen {
			err = s.SubmitOrderInDb(ctx, &order, price)
			if err != nil {
				s.logger.Error2("can not submit order in db", err,
					zap.String("service", "stopOrderSubmissionManager"),
					zap.String("method", "Submit"),
					zap.Int64("orderID", order.ID),
				)
				continue
			}
			err = s.orderRedisManager.RemoveStopOrderFromQueue(ctx, order, pairName)
			if err != nil {
				s.logger.Error2("can not remove stopOrder from redis", err,
					zap.String("service", "stopOrderSubmissionManager"),
					zap.String("method", "Submit"),
					zap.Int64("orderID", order.ID),
				)
				continue
			} else {
				s.eventsHandler.HandleOrderCreation(order, true)
			}
		} else {
			//removing non open stop orders from redis
			err = s.orderRedisManager.RemoveStopOrderFromQueue(ctx, order, pairName)
			if err != nil {
				s.logger.Error2("can not remove non open stopOrder from redis", err,
					zap.String("service", "stopOrderSubmissionManager"),
					zap.String("method", "Submit"),
					zap.Int64("orderID", order.ID),
				)
				continue
			}
		}
	}

}

func (s *stopOrderSubmissionManager) SubmitOrderInDb(ctx context.Context, order *Order, price string) error {
	order.IsSubmitted = sql.NullBool{Bool: true, Valid: true}
	order.CurrentMarketPrice = sql.NullString{String: price, Valid: true}
	err := s.db.Omit(clause.Associations).Save(order).Error
	return err
}

func (s *stopOrderSubmissionManager) getStopOrders(ctx context.Context, pairName string, formerPrice string, price string, isRising bool) ([]Order, error) {
	var orders []Order
	var orderIds []int64
	res, err := s.orderRedisManager.GetStopOrdersFromQueue(ctx, pairName, formerPrice, price, isRising)

	if err != nil {
		return orders, err
	}

	for _, z := range res {
		var orderID int64
		err := json.Unmarshal([]byte(z.Member.(string)), &orderID)
		if err == nil {
			orderIds = append(orderIds, orderID)
		} else {
			s.logger.Error2("can not unmarshal json", err,
				zap.String("service", "stopOrderSubmissionManager"),
				zap.String("method", "getStopOrders"),
				zap.String("pairName", pairName),
				zap.String("price", price),
				zap.String("formerPrice", formerPrice),
			)
		}
	}

	if len(orderIds) > 0 {
		orders = s.orderRepository.GetOrdersByIds(orderIds)
	}
	return orders, nil
}

/**
 * shouldCalculated would be false if the prices are equal or there is error
 */
func (s *stopOrderSubmissionManager) isPriceRising(price string, formerPrice string) (isRising bool, shouldCalculated bool, err error) {
	if price != "" && formerPrice != "" {
		priceDecimal, err := decimal.NewFromString(price)
		if err != nil {
			return false, false, err
		}

		formerPriceDecimal, err := decimal.NewFromString(formerPrice)
		if err != nil {
			return false, false, err
		}

		if priceDecimal.Equal(formerPriceDecimal) {
			return false, false, nil
		}

		if priceDecimal.LessThan(formerPriceDecimal) {
			return false, true, nil
		}

		return true, true, nil
	}

	return false, false, nil

}

func NewStopOrderSubmissionManager(db *gorm.DB, orderRepository Repository, liveDataService livedata.Service, orderRedisManager RedisManager, eventsHandler EventsHandler, logger platform.Logger) StopOrderSubmissionManager {
	return &stopOrderSubmissionManager{
		db:                db,
		orderRepository:   orderRepository,
		liveDataService:   liveDataService,
		orderRedisManager: orderRedisManager,
		eventsHandler:     eventsHandler,
		logger:            logger,
	}
}
