package payment

import (
	"context"
	"database/sql"
	"errors"
	"exchange-go/internal/currency"
	"exchange-go/internal/order"
	"exchange-go/internal/platform"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	FailureTypeLogical   = "LOGICAL"
	FailureTypeException = "EXCEPTION"
)

// AutoExchangeManager automatically converts deposited cryptocurrency into the
// user's preferred auto-exchange coin after a deposit is confirmed.
type AutoExchangeManager interface {
	// AutoExchange converts the deposited payment amount into the user's preferred
	// currency using the current market price.
	AutoExchange(p *Payment, ub *userbalance.UserBalance)
}

type autoExchangeManager struct {
	db                 *gorm.DB
	paymentRepository  Repository
	orderCreateManager order.CreateManager
	orderEventsHandler order.EventsHandler
	userService        user.Service
	currencyService    currency.Service
	priceGenerator     currency.PriceGenerator
	logger             platform.Logger
}

func (s *autoExchangeManager) AutoExchange(p *Payment, ub *userbalance.UserBalance) {
	if !s.shouldAutoExchange(p, ub) {
		return
	}
	var err error
	var errType string
	var reason string
	ei := &ExtraInfo{}
	err = s.paymentRepository.GetExtraInfoByPaymentID(p.ID, ei)
	if err != nil {
		s.logger.Error2("can not get extra info for payment", err,
			zap.String("service", "autoExchangeService"),
			zap.String("method", "AutoExchange"),
			zap.Int64("paymentId", p.ID),
		)
		errType = FailureTypeException
		return
	}

	defer func() {
		s.updateExtraInfoInCaseOfError(ei, err, errType, reason)
	}()
	amount := p.Amount.String
	userAgentInfo := order.UserAgentInfo{
		Device:  "",
		IP:      ei.IP.String,
		Browser: "",
	}
	u, err := s.userService.GetUserByID(p.UserID)
	if err != nil {
		s.logger.Error2("can not get user", err,
			zap.String("service", "autoExchangeService"),
			zap.String("method", "AutoExchange"),
			zap.String("amount", amount),
			zap.Int("userId", u.ID),
			zap.Int64("paymentId", p.ID),
		)
		errType = FailureTypeException
		return
	}

	pair, orderType, err := s.getPairAndType(ub)
	if err != nil {
		s.logger.Error2("error in getting pair and orderType", err,
			zap.String("service", "autoExchangeService"),
			zap.String("method", "AutoExchange"),
			zap.String("amount", amount),
			zap.Int("userId", u.ID),
			zap.Int64("paymentId", p.ID),
		)
		errType = FailureTypeException
		return
	}
	currentPrice, err := s.priceGenerator.GetPrice(context.Background(), pair.Name)
	if err != nil {
		s.logger.Error2("can not get price for pair", err,
			zap.String("service", "autoExchangeService"),
			zap.String("method", "AutoExchange"),
			zap.String("amount", amount),
			zap.Int64("pairId", pair.ID),
			zap.String("type", orderType),
			zap.Int("userId", u.ID),
			zap.Int64("paymentId", p.ID),
		)
		errType = FailureTypeException
		return
	}

	data := order.CreateRequiredData{
		User:           &u,
		Pair:           &pair,
		Amount:         amount,
		OrderType:      orderType,
		ExchangeType:   strings.ToUpper(order.ExchangeTypeMarket),
		Price:          "",
		StopPointPrice: "",
		UserAgentInfo:  userAgentInfo,
		CurrentPrice:   currentPrice,
		IsInstant:      true,
		IsFastExchange: false,
		IsAutoExchange: true,
	}
	o, err := s.orderCreateManager.CreateOrder(data)
	if err != nil {
		if !errors.Is(err, platform.OrderCreateValidationError{}) {
			s.logger.Error2("error in create order", err,
				zap.String("service", "autoExchangeService"),
				zap.String("method", "AutoExchange"),
				zap.String("amount", amount),
				zap.Int64("pairId", pair.ID),
				zap.String("type", orderType),
				zap.Int("userId", u.ID),
				zap.Int64("paymentId", p.ID),
			)
			errType = FailureTypeException
			return
		}
		errType = FailureTypeLogical
		reason = err.Error()
		return
	}

	o.Pair = pair
	o.User = u
	platform.SafeGo(s.logger, "payment.HandleOrderCreation", func() {
		s.orderEventsHandler.HandleOrderCreation(*o, false)
	})
	//todo update the payment itself and set the orderid
	ei.AutoExchangeOrderID = sql.NullInt64{Int64: o.ID, Valid: true}
	err = s.db.Omit(clause.Associations).Save(ei).Error
	if err != nil {
		s.logger.Error2("error saving extraInfo", err,
			zap.String("service", "autoExchangeService"),
			zap.String("method", "AutoExchange"),
			zap.String("amount", amount),
			zap.Int64("pairId", pair.ID),
			zap.String("type", orderType),
			zap.Int("userId", u.ID),
			zap.Int64("paymentId", p.ID),
		)
		//this error is not important for defered function because
		//the order is created successfully so we set err to nil
		err = nil
	}
}

func (s *autoExchangeManager) shouldAutoExchange(p *Payment, ub *userbalance.UserBalance) bool {
	if p.Type == TypeDeposit && p.Status == StatusCompleted &&
		ub.AutoExchangeCoin.Valid && ub.AutoExchangeCoin.String != "" {
		return true
	}
	return false
}

func (s *autoExchangeManager) getPairAndType(ub *userbalance.UserBalance) (currency.Pair, string, error) {
	allPairs := s.currencyService.GetActivePairCurrenciesList()
	for _, p := range allPairs {
		if p.BasisCoin.Code == ub.BalanceCoin && p.DependentCoin.Code == ub.AutoExchangeCoin.String {
			orderType := order.TypeBuy
			return p, orderType, nil
		}
		if p.BasisCoin.Code == ub.AutoExchangeCoin.String && p.DependentCoin.Code == ub.BalanceCoin {
			orderType := order.TypeSell
			return p, orderType, nil
		}
	}

	return currency.Pair{}, "", fmt.Errorf("could not found pair for coin %s and exchangeCoin %s", ub.BalanceCoin, ub.AutoExchangeCoin.String)
}

func (s *autoExchangeManager) updateExtraInfoInCaseOfError(ei *ExtraInfo, err error, errorType string, reason string) {
	if err == nil {
		return
	}
	ei.AutoExchangeFailureType = sql.NullString{String: errorType, Valid: true}
	//due to security reason we do not save exception type reason
	if errorType == FailureTypeLogical {
		ei.AutoExchangeFailureReason = sql.NullString{String: reason, Valid: true}
	}
	saveErr := s.db.Omit(clause.Associations).Save(ei).Error
	if saveErr != nil {
		s.logger.Error2("error in create order", saveErr,
			zap.String("service", "autoExchangeService"),
			zap.String("method", "updateExtraInfoInCaseOfError"),
			zap.Int64("paymentId", ei.PaymentID),
		)
	}

}

func NewAutoExchangeManger(db *gorm.DB, paymentRepository Repository, orderCreateManager order.CreateManager, orderEventsHandler order.EventsHandler, userService user.Service, currencyService currency.Service, priceGenerator currency.PriceGenerator, logger platform.Logger) AutoExchangeManager {
	return &autoExchangeManager{
		db:                 db,
		paymentRepository:  paymentRepository,
		orderCreateManager: orderCreateManager,
		orderEventsHandler: orderEventsHandler,
		userService:        userService,
		currencyService:    currencyService,
		priceGenerator:     priceGenerator,
		logger:             logger,
	}
}
