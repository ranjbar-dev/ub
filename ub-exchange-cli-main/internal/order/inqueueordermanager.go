package order

import (
	"context"
	"exchange-go/internal/engine"
	"exchange-go/internal/platform"
	"fmt"

	"go.uber.org/zap"
)

// InQueueOrderManager handles stop orders that are triggered when the market
// price reaches their stop-point level.
type InQueueOrderManager interface {
	// HandleInQueueOrders processes queued stop orders for a given pair when the
	// current price reaches the trigger level, submitting them to the matching engine.
	HandleInQueueOrders(ctx context.Context, pairName string, price string)
}

type inQueueOrderManager struct {
	engine engine.Engine
	logger platform.Logger
}

func (s *inQueueOrderManager) HandleInQueueOrders(ctx context.Context, pairName string, price string) {
	if price == "" {
		err := platform.NonSentryError{Err: fmt.Errorf("price is empty")}
		s.logger.Error2("empty price", err,
			zap.String("service", "inQueueOrderManager"),
			zap.String("method", "HandleInQueueOrders"),
			zap.String("pairName", pairName),
		)
		return
	}
	err := s.engine.HandleInQueueOrders(pairName, price)
	if err != nil {
		s.logger.Error2("can not handle in queue orders", err,
			zap.String("service", "inQueueOrderManager"),
			zap.String("method", "HandleInQueueOrders"),
			zap.String("pairName", pairName),
			zap.String("price", price),
		)
	}
}

func NewInQueueOrderManager(e engine.Engine, logger platform.Logger) InQueueOrderManager {
	return &inQueueOrderManager{
		engine: e,
		logger: logger,
	}
}
