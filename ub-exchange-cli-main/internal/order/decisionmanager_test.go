// Package order_test tests the DecisionManager. Covers:
//   - Market BUY order within pair exchange limit routed to our exchange
//   - Market BUY order exceeding pair exchange limit routed to external exchange
//   - All LIMIT orders always routed to our exchange regardless of amount
//
// Test data: mocked Configs returning test environment, currency pair with
// MaxOurExchangeLimit, and order structs with varying demanded amounts.
package order_test

import (
	"database/sql"
	"exchange-go/internal/currency"
	"exchange-go/internal/mocks"
	"exchange-go/internal/order"
	"exchange-go/internal/platform"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecisionManager_DecideOrderPlacement(t *testing.T) {
	configs := new(mocks.Configs)
	configs.On("GetEnv").Times(3).Return(platform.EnvTest)
	dm := order.NewDecisionManager(configs)
	pair := currency.Pair{
		ID:                  1,
		MaxOurExchangeLimit: "3.0",
	}
	o1 := order.Order{
		ID:             1,
		UserID:         1,
		Pair:           pair,
		Type:           "BUY",
		ExchangeType:   "MARKET",
		DemandedAmount: sql.NullString{String: "2.5", Valid: true},
	}
	orderPlace, err := dm.DecideOrderPlacement(o1)
	assert.Nil(t, err)
	assert.Equal(t, order.PlaceOurExchange, orderPlace)

	o2 := order.Order{
		ID:             1,
		UserID:         1,
		Pair:           pair,
		Type:           "BUY",
		ExchangeType:   "MARKET",
		DemandedAmount: sql.NullString{String: "3.5", Valid: true},
	}

	orderPlace, err = dm.DecideOrderPlacement(o2)
	assert.Nil(t, err)
	assert.Equal(t, order.PlaceExternalExchange, orderPlace)

	//all limit orders goes to our exchange
	o3 := order.Order{
		ID:             1,
		UserID:         1,
		Pair:           pair,
		Type:           "BUY",
		ExchangeType:   "LIMIT",
		Price:          sql.NullString{String: "50000", Valid: true},
		DemandedAmount: sql.NullString{String: "3.5", Valid: true},
	}

	orderPlace, err = dm.DecideOrderPlacement(o3)
	assert.Nil(t, err)
	assert.Equal(t, order.PlaceOurExchange, orderPlace)

}
