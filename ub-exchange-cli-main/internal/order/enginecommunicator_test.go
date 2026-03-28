// Package order_test tests the EngineCommunicator. Covers:
//   - SubmitOrder: computes min/max threshold prices via ForceTrader and submits to engine
//   - RemoveOrder: computes threshold prices and removes order from engine
//   - RetrieveOrder: computes threshold prices and retrieves order from engine
//
// Test data: mocked ForceTrader returning threshold prices for BTC-USDT BUY orders,
// mocked Engine with SubmitOrder/RemoveOrder/RetrieveOrder expectations.
package order_test

import (
	"database/sql"
	"exchange-go/internal/currency"
	"exchange-go/internal/mocks"
	"exchange-go/internal/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestEngineCommunicator_SubmitOrder(t *testing.T) {
	ft := new(mocks.ForceTrader)
	ft.On("GetMinAndMaxPrice", "BTC-USDT", "BUY", "50000.00000000").Once().Return("49950", "50050", nil)
	e := new(mocks.Engine)
	e.On("SubmitOrder", mock.Anything).Once().Return(nil)

	ec := order.NewEngineCommunicator(ft, e)
	pair := currency.Pair{
		ID:   1,
		Name: "BTC-USDT",
	}

	o := order.Order{
		ID:                 1,
		Pair:               pair,
		Type:               "BUY",
		DemandedAmount:     sql.NullString{String: "0.1", Valid: true},
		CurrentMarketPrice: sql.NullString{String: "50000.00000000", Valid: true},
	}

	err := ec.SubmitOrder(o)
	assert.Nil(t, err)
	ft.AssertExpectations(t)
	e.AssertExpectations(t)
}

func TestEngineCommunicator_RemoveOrder(t *testing.T) {
	ft := new(mocks.ForceTrader)
	ft.On("GetMinAndMaxPrice", "BTC-USDT", "BUY", "50000.00000000").Once().Return("49950", "50050", nil)
	e := new(mocks.Engine)
	e.On("RemoveOrder", mock.Anything).Once().Return(nil)

	ec := order.NewEngineCommunicator(ft, e)
	pair := currency.Pair{
		ID:   1,
		Name: "BTC-USDT",
	}

	o := order.Order{
		ID:                 1,
		Pair:               pair,
		Type:               "BUY",
		DemandedAmount:     sql.NullString{String: "0.1", Valid: true},
		CurrentMarketPrice: sql.NullString{String: "50000.00000000", Valid: true},
	}

	err := ec.RemoveOrder(o)
	assert.Nil(t, err)
	e.AssertExpectations(t)
}

func TestEngineCommunicator_RetrieveOrder(t *testing.T) {
	ft := new(mocks.ForceTrader)
	ft.On("GetMinAndMaxPrice", "BTC-USDT", "BUY", "50000.00000000").Once().Return("49950", "50050", nil)
	e := new(mocks.Engine)
	e.On("RetrieveOrder", mock.Anything).Once().Return(nil)

	ec := order.NewEngineCommunicator(ft, e)
	pair := currency.Pair{
		ID:   1,
		Name: "BTC-USDT",
	}

	o := order.Order{
		ID:                 1,
		Pair:               pair,
		Type:               "BUY",
		DemandedAmount:     sql.NullString{String: "0.1", Valid: true},
		CurrentMarketPrice: sql.NullString{String: "50000.00000000", Valid: true},
	}

	err := ec.RetrieveOrder(o)
	assert.Nil(t, err)
	e.AssertExpectations(t)
}
