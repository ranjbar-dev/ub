package order

import (
	"context"
	"exchange-go/internal/currency"
	"exchange-go/internal/platform"
	"fmt"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// AdminOrderManager provides admin-only operations for manually fulfilling orders
// that have not been matched by the engine.
type AdminOrderManager interface {
	// TryToFulfillOrder verifies that the order's price has been reached in market history,
	// then manually fulfills it through the post-order matching service.
	TryToFulfillOrder(o Order) error
}

type adminOrderManager struct {
	currencyService            currency.Service
	klineService               currency.KlineService
	priceGenerator             currency.PriceGenerator
	postOrderMatchingService   PostOrderMatchingService
	stopOrderSubmissionManager StopOrderSubmissionManager
	eventsHandler              EventsHandler
	logger                     platform.Logger
}

func (s *adminOrderManager) TryToFulfillOrder(order Order) error {
	if !order.IsStopOrder() {
		isTouched, price, err := s.isOrderPriceTouched(order)
		if err != nil {
			s.logger.Error2("can not check if the order price is touched", err,
				zap.String("service", "adminOrderManager"),
				zap.String("method", "TryToFulfillOrder"),
				zap.Int64("orderID", order.ID),
			)

			return fmt.Errorf("something went wrong")
		}
		if !isTouched {
			return fmt.Errorf("price is not touched since the time the order is created")
		}

		err = s.fulfillOrder(order, price)
		s.eventsHandler.HandleOrderFulfillByAdmin(order)
		return err
	}

	isStopPointPriceTouched, isPriceTouched, err := s.isStopOrderPriceTouched(order)
	if err != nil {
		s.logger.Error2("can not check if the stopOrder price is touched", err,
			zap.String("service", "adminOrderManager"),
			zap.String("method", "TryToFulfillOrder"),
			zap.Int64("orderID", order.ID),
		)
		return fmt.Errorf("something went wrong")
	}

	if !isStopPointPriceTouched {
		return fmt.Errorf("stop point price is not touched since the time the order is created")
	}

	// first we submit the stop order then try to fulfill it
	ctx := context.Background()
	if !order.IsSubmitted.Valid || !order.IsSubmitted.Bool {
		currentPrice, err := s.priceGenerator.GetPrice(ctx, order.Pair.Name)
		if err != nil {
			s.logger.Error2("can not the currentPrice", err,
				zap.String("service", "adminOrderManager"),
				zap.String("method", "TryToFulfillOrder"),
				zap.Int64("orderID", order.ID),
			)
			return fmt.Errorf("something went wrong")
		}
		err = s.stopOrderSubmissionManager.SubmitOrderInDb(ctx, &order, currentPrice)
		if err != nil {
			s.logger.Error2("can not submit order in db", err,
				zap.String("service", "adminOrderManager"),
				zap.String("method", "TryToFulfillOrder"),
				zap.Int64("orderID", order.ID),
			)
			return fmt.Errorf("something went wrong")
		}
	}

	if isPriceTouched {
		err = s.fulfillOrder(order, "")
		return err
	}

	return nil
}

func (s *adminOrderManager) isOrderPriceTouched(o Order) (isTouched bool, price string, err error) {
	from := o.CreatedAt
	pairID := o.PairID

	pair, err := s.currencyService.GetPairByID(pairID)
	if err != nil {
		return false, "", err
	}
	highAndLow, err := s.klineService.GetHighAndLowPriceFromDateForPairByPairName(pair.Name, from)
	if err != nil {
		return false, "", err
	}

	minPriceDecimal, err := decimal.NewFromString(highAndLow.Low)
	if err != nil {
		return false, "", err
	}

	maxPriceDecimal, err := decimal.NewFromString(highAndLow.High)
	if err != nil {
		return false, "", err
	}

	if o.isMarket() {
		if o.Type == TypeBuy {
			price = highAndLow.Low
		} else {
			price = highAndLow.High
		}
		return true, price, nil
	}

	orderPrice := o.Price.String
	orderPriceDecimal, _ := decimal.NewFromString(orderPrice)
	if o.Type == TypeBuy {
		return orderPriceDecimal.GreaterThanOrEqual(minPriceDecimal), orderPrice, nil
	}

	return orderPriceDecimal.LessThanOrEqual(maxPriceDecimal), orderPrice, nil

}

func (s *adminOrderManager) isStopOrderPriceTouched(o Order) (isStopPointPriceTouched bool, isPriceTouched bool, err error) {
	from := o.CreatedAt
	pairID := o.PairID
	pair, err := s.currencyService.GetPairByID(pairID)
	if err != nil {
		return false, false, err
	}

	stopPointPriceDecimal, _ := decimal.NewFromString(o.StopPointPrice.String)
	orderPriceDecimal, _ := decimal.NewFromString(o.Price.String)

	highAndLow, err := s.klineService.GetHighAndLowPriceFromDateForPairByPairName(pair.Name, from)
	if err != nil {
		return false, false, err
	}

	minPriceDecimal, err := decimal.NewFromString(highAndLow.Low)
	if err != nil {
		return false, false, err
	}

	maxPriceDecimal, err := decimal.NewFromString(highAndLow.High)
	if err != nil {
		return false, false, err
	}

	isStopPointPriceTouched = minPriceDecimal.LessThanOrEqual(stopPointPriceDecimal) && maxPriceDecimal.GreaterThanOrEqual(stopPointPriceDecimal)

	if o.Type == TypeBuy {
		isPriceTouched = orderPriceDecimal.GreaterThanOrEqual(minPriceDecimal)

	} else {
		isPriceTouched = orderPriceDecimal.LessThanOrEqual(maxPriceDecimal)

	}

	return isStopPointPriceTouched, isPriceTouched, nil

}

func (s *adminOrderManager) fulfillOrder(order Order, price string) error {
	if !order.isMarket() {
		price = order.Price.String
	}
	ordersData := make([]CallBackOrderData, 0)
	partial := &CallBackOrderData{
		ID:                order.ID,
		PairName:          order.Pair.Name,
		OrderType:         order.Type,
		Quantity:          order.getAmount(),
		Price:             price,
		Timestamp:         order.CreatedAt.Unix(),
		TradedWithOrderID: 0,
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "",
		MaxThresholdPrice: "",
	}
	result := s.postOrderMatchingService.HandlePostOrderMatching(ordersData, partial, true)
	if result.Err != nil {
		return fmt.Errorf("can not fulfill order")

	}
	return nil
}

func NewAdminOrderManager(currencyService currency.Service, klineService currency.KlineService, priceGenerator currency.PriceGenerator, postOrderMatchingService PostOrderMatchingService, stopOrderSubmissionManager StopOrderSubmissionManager, eventsHandler EventsHandler, logger platform.Logger) AdminOrderManager {
	return &adminOrderManager{
		currencyService:            currencyService,
		klineService:               klineService,
		priceGenerator:             priceGenerator,
		postOrderMatchingService:   postOrderMatchingService,
		stopOrderSubmissionManager: stopOrderSubmissionManager,
		eventsHandler:              eventsHandler,
		logger:                     logger,
	}
}
