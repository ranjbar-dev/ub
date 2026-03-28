// Package order_test tests the ForceTrader. Covers:
//   - GetMinAndMaxPrice with PERCENTAGE bot rules: applies percentage spread to market price
//   - GetMinAndMaxPrice with CONST bot rules: applies fixed-value spread to market price
//   - ShouldForceTrade: returns true when price is within threshold, false when outside,
//     and correctly handles boundary values for both BUY and SELL order types
//
// Test data: mocked PriceGenerator returning fixed market price, mocked CurrencyService
// with BTC-USDT pair carrying JSON-encoded bot rules (PERCENTAGE and CONST types).
package order_test

import (
	"database/sql"
	"exchange-go/internal/currency"
	"exchange-go/internal/mocks"
	"exchange-go/internal/order"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestForceTrader_GetMinAndMaxPrice_Percentage(t *testing.T) {
	pg := new(mocks.PriceGenerator)
	pg.On("GetPrice", mock.Anything, mock.Anything).Once().Return("50000", nil)
	botRules := `{"buyValue":0.01,"sellValue":0.01,"type":"PERCENTAGE"}`
	BTCUSDT := currency.Pair{
		ID:       1,
		Name:     "BTC-USDT",
		BotRules: sql.NullString{String: botRules, Valid: true},
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetPairByName", "BTC-USDT").Once().Return(BTCUSDT, nil)
	ft := order.NewForceTrader(pg, currencyService)
	pairName := "BTC-USDT"
	orderType := "BUY"
	min, max, err := ft.GetMinAndMaxPrice(pairName, orderType, "")
	assert.Nil(t, err)
	assert.Equal(t, "49500.00000000", min)
	assert.Equal(t, "50500.00000000", max)
	pg.AssertExpectations(t)
	currencyService.AssertExpectations(t)
}

func TestForceTrader_GetMinAndMaxPrice_Const(t *testing.T) {
	pg := new(mocks.PriceGenerator)
	pg.On("GetPrice", mock.Anything, mock.Anything).Once().Return("50000", nil)
	botRules := `{"buyValue":10,"sellValue":10,"type":"CONST"}`
	BTCUSDT := currency.Pair{
		ID:       1,
		Name:     "BTC-USDT",
		BotRules: sql.NullString{String: botRules, Valid: true},
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetPairByName", "BTC-USDT").Once().Return(BTCUSDT, nil)
	ft := order.NewForceTrader(pg, currencyService)
	pairName := "BTC-USDT"
	orderType := "BUY"
	min, max, err := ft.GetMinAndMaxPrice(pairName, orderType, "")
	assert.Nil(t, err)
	assert.Equal(t, "49990.00000000", min)
	assert.Equal(t, "50010.00000000", max)
	pg.AssertExpectations(t)
	currencyService.AssertExpectations(t)
}

func TestForceTrader_ShouldForceTrade(t *testing.T) {
	pg := new(mocks.PriceGenerator)
	pg.On("GetPrice", mock.Anything, mock.Anything).Times(5).Return("50000", nil)
	botRules := `{"buyValue":0.01,"sellValue":0.01,"type":"PERCENTAGE"}`
	BTCUSDT := currency.Pair{
		ID:       1,
		Name:     "BTC-USDT",
		BotRules: sql.NullString{String: botRules, Valid: true},
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetPairByName", "BTC-USDT").Times(5).Return(BTCUSDT, nil)
	ft := order.NewForceTrader(pg, currencyService)
	pairName := "BTC-USDT"
	orderType := "BUY"
	price := "50000"
	shouldForceTrade, err := ft.ShouldForceTrade(pairName, orderType, price)
	assert.Nil(t, err)
	assert.Equal(t, true, shouldForceTrade)

	price = "40000"
	shouldForceTrade, err = ft.ShouldForceTrade(pairName, orderType, price)
	assert.Nil(t, err)
	assert.Equal(t, false, shouldForceTrade)

	orderType = "SELL"
	price = "60000"
	shouldForceTrade, err = ft.ShouldForceTrade(pairName, orderType, price)
	assert.Nil(t, err)
	assert.Equal(t, false, shouldForceTrade)

	//testing border numbers
	orderType = "SELL"
	price = "49500"
	shouldForceTrade, err = ft.ShouldForceTrade(pairName, orderType, price)
	assert.Nil(t, err)
	assert.Equal(t, true, shouldForceTrade)

	orderType = "BUY"
	price = "50500"
	shouldForceTrade, err = ft.ShouldForceTrade(pairName, orderType, price)
	assert.Nil(t, err)
	assert.Equal(t, true, shouldForceTrade)

	pg.AssertExpectations(t)
	currencyService.AssertExpectations(t)
}
