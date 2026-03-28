package order

import (
	"exchange-go/internal/currency"
	"exchange-go/internal/platform"

	"github.com/shopspring/decimal"
)

// DecisionManager determines whether an order should be routed to the internal
// matching engine or to an external exchange (e.g., Binance).
type DecisionManager interface {
	// DecideOrderPlacement evaluates the order size against the pair's exchange limit
	// and returns the placement destination (internal or external).
	DecideOrderPlacement(o Order) (string, error)
}

type decisionManager struct {
	configs platform.Configs
}

func (d *decisionManager) DecideOrderPlacement(o Order) (string, error) {
	place := PlaceOurExchange
	if d.configs.GetEnv() == platform.EnvDev {
		return place, nil
	}
	limit := o.Pair.MaxOurExchangeLimit
	aggregationStatus := o.Pair.AggregationStatus
	if aggregationStatus == currency.AggregationStatusPause || aggregationStatus == currency.AggregationStatusStop {
		return place, nil
	}
	limitDecimal, err := decimal.NewFromString(limit)
	if err != nil {
		return "", err
	}
	AmountDecimal, err := decimal.NewFromString(o.getAmount())
	if err != nil {
		return "", err
	}
	if AmountDecimal.GreaterThan(limitDecimal) && o.isMarket() {
		place = PlaceExternalExchange
	}
	return place, nil
}

func NewDecisionManager(configs platform.Configs) DecisionManager {
	return &decisionManager{
		configs: configs,
	}
}
