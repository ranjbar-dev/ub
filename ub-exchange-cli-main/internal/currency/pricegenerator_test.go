// Package currency_test tests the price generator component. Covers:
//   - Converting coin amounts to USDT equivalents (direct pairs and cross-pair routing)
//   - Converting coin amounts to BTC equivalents (direct and indirect via cross pairs)
//   - Fetching BTC-USDT price with Redis cache-miss fallback to kline service
//   - Calculating pair prices in USDT terms (including BTC-based cross pairs)
//   - Retrieving raw pair prices from live data
//
// Test data: table-driven test cases for USDT/BTC/GRS/DAI/ETH conversions,
// mock live data, kline service, and pair repository with multi-pair active lists.
package currency_test

import (
	"context"
	"exchange-go/internal/currency"
	"exchange-go/internal/mocks"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var amountBasedOnUSDTTable = []struct {
	coin   string
	amount string
	output string
	err    error
}{
	{"USDT", "1", "1", nil},
	{"BTC", "1", "40000.00000000", nil},
	{"ETH", "1", "1500.00000000", nil},
	{"DAI", "1", "1.01010101", nil},
	{"GRS", "1", "0.46400000", nil},
}

func TestPriceGenerator_GetAmountBasedOnUSDT(t *testing.T) {
	liveData := new(mocks.LiveData)
	liveData.On("GetPrice", mock.Anything, "BTC-USDT").Twice().Return("40000", nil)
	liveData.On("GetPrice", mock.Anything, "ETH-USDT").Once().Return("1500", nil)
	liveData.On("GetPrice", mock.Anything, "USDT-DAI").Once().Return("0.99", nil)
	liveData.On("GetPrice", mock.Anything, "GRS-BTC").Once().Return("0.0000116", nil)
	klineService := new(mocks.KlineService)
	pairRepository := new(mocks.PairRepository)
	activePairs := []currency.Pair{
		{
			ID:   1,
			Name: "BTC-USDT",
			BasisCoin: currency.Coin{
				Code: "USDT",
			},
			DependentCoin: currency.Coin{
				Code: "BTC",
			},
		},
		{
			ID:   2,
			Name: "ETH-USDT",
			BasisCoin: currency.Coin{
				Code: "USDT",
			},
			DependentCoin: currency.Coin{
				Code: "ETH",
			},
		},
		{
			ID:         3,
			Name:       "USDT-DAI",
			ShowDigits: 6,
			BasisCoin: currency.Coin{
				Code: "DAI",
			},
			DependentCoin: currency.Coin{
				Code: "USDT",
			},
		},
		{
			ID:   4,
			Name: "GRS-BTC",
			BasisCoin: currency.Coin{
				Code: "BTC",
			},
			DependentCoin: currency.Coin{
				Code: "GRS",
			},
		},
		{
			ID:   6,
			Name: "BTC-DAI",
			BasisCoin: currency.Coin{
				Code: "DAI",
			},
			DependentCoin: currency.Coin{
				Code: "BTC",
			},
		},
	}

	pairRepository.On("GetActivePairCurrenciesList").Once().Return(activePairs)

	priceGenerator := currency.NewPriceGenerator(liveData, klineService, pairRepository)
	ctx := context.Background()

	for _, item := range amountBasedOnUSDTTable {
		amount, err := priceGenerator.GetAmountBasedOnUSDT(ctx, item.coin, item.amount)
		assert.Nil(t, err)
		assert.Equal(t, item.output, amount)
	}
	liveData.AssertExpectations(t)
	klineService.AssertExpectations(t)
	pairRepository.AssertExpectations(t)
}

var amountBasedOnBTCTable = []struct {
	coin   string
	amount string
	output string
	err    error
}{
	{"USDT", "1", "0.00002500", nil},
	{"BTC", "1", "1", nil},
	{"ETH", "1", "0.03000000", nil},
	{"DAI", "1", "0.00002500", nil},
	{"GRS", "1", "0.00001160", nil},
}

func TestPriceGenerator_GetAmountBasedOnBTC(t *testing.T) {
	liveData := new(mocks.LiveData)
	liveData.On("GetPrice", mock.Anything, "BTC-USDT").Once().Return("40000", nil)
	liveData.On("GetPrice", mock.Anything, "BTC-DAI").Once().Return("40000", nil)
	liveData.On("GetPrice", mock.Anything, "GRS-BTC").Once().Return("0.0000116", nil)
	liveData.On("GetPrice", mock.Anything, "ETH-BTC").Once().Return("0.03", nil)
	klineService := new(mocks.KlineService)
	pairRepository := new(mocks.PairRepository)
	activePairs := []currency.Pair{
		{
			ID:   1,
			Name: "BTC-USDT",
			BasisCoin: currency.Coin{
				Code: "USDT",
			},
			DependentCoin: currency.Coin{
				Code: "BTC",
			},
		},
		{
			ID:   2,
			Name: "ETH-USDT",
			BasisCoin: currency.Coin{
				Code: "USDT",
			},
			DependentCoin: currency.Coin{
				Code: "ETH",
			},
		},
		{
			ID:         3,
			Name:       "USDT-DAI",
			ShowDigits: 6,
			BasisCoin: currency.Coin{
				Code: "DAI",
			},
			DependentCoin: currency.Coin{
				Code: "USDT",
			},
		},
		{
			ID:   4,
			Name: "GRS-BTC",
			BasisCoin: currency.Coin{
				Code: "BTC",
			},
			DependentCoin: currency.Coin{
				Code: "GRS",
			},
		},
		{
			ID:   5,
			Name: "ETH-BTC",
			BasisCoin: currency.Coin{
				Code: "BTC",
			},
			DependentCoin: currency.Coin{
				Code: "ETH",
			},
		},
		{
			ID:   6,
			Name: "BTC-DAI",
			BasisCoin: currency.Coin{
				Code: "DAI",
			},
			DependentCoin: currency.Coin{
				Code: "BTC",
			},
		},
	}

	pairRepository.On("GetActivePairCurrenciesList").Once().Return(activePairs)

	priceGenerator := currency.NewPriceGenerator(liveData, klineService, pairRepository)
	ctx := context.Background()
	for _, item := range amountBasedOnBTCTable {
		amount, err := priceGenerator.GetAmountBasedOnBTC(ctx, item.coin, item.amount)
		assert.Nil(t, err)
		assert.Equal(t, item.output, amount)
	}
	liveData.AssertExpectations(t)
	klineService.AssertExpectations(t)
	pairRepository.AssertExpectations(t)
}

func TestPriceGenerator_GetBTCUSDTPrice(t *testing.T) {
	liveData := new(mocks.LiveData)
	liveData.On("GetPrice", mock.Anything, "BTC-USDT").Once().Return("", redis.Nil)
	klineService := new(mocks.KlineService)
	klineService.On("GetLastPriceForPair", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return("40000", nil)
	pairRepository := new(mocks.PairRepository)
	activePairs := []currency.Pair{
		{
			ID:   1,
			Name: "BTC-USDT",
			BasisCoin: currency.Coin{
				Code: "USDT",
			},
			DependentCoin: currency.Coin{
				Code: "BTC",
			},
		},
		{
			ID:   2,
			Name: "ETH-USDT",
			BasisCoin: currency.Coin{
				Code: "USDT",
			},
			DependentCoin: currency.Coin{
				Code: "ETH",
			},
		},
		{
			ID:         3,
			Name:       "USDT-DAI",
			ShowDigits: 6,
			BasisCoin: currency.Coin{
				Code: "DAI",
			},
			DependentCoin: currency.Coin{
				Code: "USDT",
			},
		},
		{
			ID:   4,
			Name: "GRS-BTC",
			BasisCoin: currency.Coin{
				Code: "BTC",
			},
			DependentCoin: currency.Coin{
				Code: "GRS",
			},
		},
		{
			ID:   5,
			Name: "ETH-BTC",
			BasisCoin: currency.Coin{
				Code: "BTC",
			},
			DependentCoin: currency.Coin{
				Code: "ETH",
			},
		},
		{
			ID:   6,
			Name: "BTC-DAI",
			BasisCoin: currency.Coin{
				Code: "DAI",
			},
			DependentCoin: currency.Coin{
				Code: "BTC",
			},
		},
	}

	pairRepository.On("GetActivePairCurrenciesList").Once().Return(activePairs)

	priceGenerator := currency.NewPriceGenerator(liveData, klineService, pairRepository)
	ctx := context.Background()

	price, err := priceGenerator.GetBTCUSDTPrice(ctx)
	assert.Nil(t, err)
	assert.Equal(t, "40000.00000000", price)

	liveData.AssertExpectations(t)
	klineService.AssertExpectations(t)

}

var pairPriceBasedOnUSDTTable = []struct {
	pairName string
	output   string
	err      error
}{
	{"BTC-USDT", "40000.00000000", nil},
	{"ETH-USDT", "1500.00000000", nil},
	{"GRS-BTC", "0.46400000", nil},
	{"USDT-DAI", "1.00000000", nil},
	{"ETH-BTC", "1200.00000000", nil},
}

func TestPriceGenerator_GetPairPriceBasedOnUSDT(t *testing.T) {
	liveData := new(mocks.LiveData)
	liveData.On("GetPrice", mock.Anything, "BTC-USDT").Times(3).Return("40000", nil)
	liveData.On("GetPrice", mock.Anything, "ETH-USDT").Once().Return("1500", nil)
	liveData.On("GetPrice", mock.Anything, "USDT-DAI").Twice().Return("0.99999999", nil)
	liveData.On("GetPrice", mock.Anything, "GRS-BTC").Once().Return("0.0000116", nil)
	liveData.On("GetPrice", mock.Anything, "ETH-BTC").Once().Return("0.03", nil)
	klineService := new(mocks.KlineService)
	pairRepository := new(mocks.PairRepository)

	activePairs := []currency.Pair{
		{
			ID:   1,
			Name: "BTC-USDT",
			BasisCoin: currency.Coin{
				Code: "USDT",
			},
			DependentCoin: currency.Coin{
				Code: "BTC",
			},
		},
		{
			ID:   2,
			Name: "ETH-USDT",
			BasisCoin: currency.Coin{
				Code: "USDT",
			},
			DependentCoin: currency.Coin{
				Code: "ETH",
			},
		},
		{
			ID:         3,
			Name:       "USDT-DAI",
			ShowDigits: 6,
			BasisCoin: currency.Coin{
				Code: "DAI",
			},
			DependentCoin: currency.Coin{
				Code: "USDT",
			},
		},
		{
			ID:   4,
			Name: "GRS-BTC",
			BasisCoin: currency.Coin{
				Code: "BTC",
			},
			DependentCoin: currency.Coin{
				Code: "GRS",
			},
		},
		{
			ID:   5,
			Name: "ETH-BTC",
			BasisCoin: currency.Coin{
				Code: "BTC",
			},
			DependentCoin: currency.Coin{
				Code: "ETH",
			},
		},
		{
			ID:   6,
			Name: "BTC-DAI",
			BasisCoin: currency.Coin{
				Code: "DAI",
			},
			DependentCoin: currency.Coin{
				Code: "BTC",
			},
		},
	}

	pairRepository.On("GetActivePairCurrenciesList").Once().Return(activePairs)
	priceGenerator := currency.NewPriceGenerator(liveData, klineService, pairRepository)
	ctx := context.Background()
	for _, item := range pairPriceBasedOnUSDTTable {
		amount, err := priceGenerator.GetPairPriceBasedOnUSDT(ctx, item.pairName)
		assert.Nil(t, err)
		assert.Equal(t, item.output, amount)
	}
	liveData.AssertExpectations(t)
	klineService.AssertExpectations(t)
	pairRepository.AssertExpectations(t)
}

var getPriceTable = []struct {
	pairName string
	output   string
	err      error
}{
	{"BTC-USDT", "50000.00000000", nil},
	{"ETH-USDT", "2000.00000000", nil},
	{"GRS-BTC", "0.00001160", nil},
	{"USDT-DAI", "0.99999999", nil},
	{"ETH-BTC", "0.03000000", nil},
}

func TestPriceGenerator_GetPrice(t *testing.T) {
	liveData := new(mocks.LiveData)
	liveData.On("GetPrice", mock.Anything, "BTC-USDT").Once().Return("50000", nil)
	liveData.On("GetPrice", mock.Anything, "ETH-USDT").Once().Return("2000", nil)
	liveData.On("GetPrice", mock.Anything, "GRS-BTC").Once().Return("0.0000116", nil)
	liveData.On("GetPrice", mock.Anything, "USDT-DAI").Once().Return("0.99999999", nil)
	liveData.On("GetPrice", mock.Anything, "ETH-BTC").Once().Return("0.03", nil)
	klineService := new(mocks.KlineService)
	pairRepository := new(mocks.PairRepository)
	activePairs := []currency.Pair{
		{
			ID:   1,
			Name: "BTC-USDT",
			BasisCoin: currency.Coin{
				Code: "USDT",
			},
			DependentCoin: currency.Coin{
				Code: "BTC",
			},
		},
		{
			ID:   2,
			Name: "ETH-USDT",
			BasisCoin: currency.Coin{
				Code: "USDT",
			},
			DependentCoin: currency.Coin{
				Code: "ETH",
			},
		},
		{
			ID:         3,
			Name:       "USDT-DAI",
			ShowDigits: 6,
			BasisCoin: currency.Coin{
				Code: "DAI",
			},
			DependentCoin: currency.Coin{
				Code: "USDT",
			},
		},
		{
			ID:   4,
			Name: "GRS-BTC",
			BasisCoin: currency.Coin{
				Code: "BTC",
			},
			DependentCoin: currency.Coin{
				Code: "GRS",
			},
		},
		{
			ID:   5,
			Name: "ETH-BTC",
			BasisCoin: currency.Coin{
				Code: "BTC",
			},
			DependentCoin: currency.Coin{
				Code: "ETH",
			},
		},
		{
			ID:   6,
			Name: "BTC-DAI",
			BasisCoin: currency.Coin{
				Code: "DAI",
			},
			DependentCoin: currency.Coin{
				Code: "BTC",
			},
		},
	}

	pairRepository.On("GetActivePairCurrenciesList").Once().Return(activePairs)

	priceGenerator := currency.NewPriceGenerator(liveData, klineService, pairRepository)
	ctx := context.Background()

	for _, item := range getPriceTable {
		amount, err := priceGenerator.GetPrice(ctx, item.pairName)
		assert.Nil(t, err)
		assert.Equal(t, item.output, amount)
	}
	liveData.AssertExpectations(t)
	klineService.AssertExpectations(t)
	pairRepository.AssertExpectations(t)
}
