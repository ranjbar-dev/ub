// Package currency_test tests the currency service. Covers:
//   - Listing active trading pairs with live price data and USDT equivalents
//   - Retrieving pairs by ID and by name
//   - Looking up coins by code
//   - Fetching active currencies and full currency lists
//   - Pair statistics with 24-hour high/low prices and volume
//   - Adding and removing favorite pairs (including idempotent duplicates)
//   - Pair ratio calculations (USDT-based, inverse USDT, and non-USDT cross pairs)
//   - Retrieving pairs list and trading fee structures
//
// Test data: mock repositories, live data service, price generator, kline service,
// candle gRPC client, and user-favorite pair repository with pre-built pair/coin fixtures.
package currency_test

import (
	"context"
	"database/sql"
	"exchange-go/internal/currency"
	"exchange-go/internal/currency/candle"
	"exchange-go/internal/livedata"
	"exchange-go/internal/mocks"
	"exchange-go/internal/user"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestService_GetPairs(t *testing.T) {
	repo := new(mocks.CurrencyRepository)
	pairPricesData := []livedata.RedisPairPriceData{
		{
			PairName:   "BTC-USDT",
			Price:      "50000",
			Percentage: "2.1",
		},
		{
			PairName:   "ETH-USDT",
			Price:      "2000",
			Percentage: "1.1",
		},
		{
			PairName:   "GRS-BTC",
			Price:      "0.000013",
			Percentage: "0.1",
		},
		{
			PairName:   "USDT-DAI",
			Price:      "0.99",
			Percentage: "0.01",
		},
		{
			PairName:   "BTC-DAI",
			Price:      "50000",
			Percentage: "2.1",
		},
	}
	liveData := new(mocks.LiveData)

	liveData.On("GetPairsPriceData", mock.Anything, mock.Anything).Once().Return(pairPricesData, nil)

	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetPairPriceBasedOnUSDT", mock.Anything, "BTC-USDT").Once().Return("50000", nil)
	priceGenerator.On("GetPairPriceBasedOnUSDT", mock.Anything, "ETH-USDT").Once().Return("2000", nil)
	priceGenerator.On("GetPairPriceBasedOnUSDT", mock.Anything, "USDT-DAI").Once().Return("0.99", nil)
	priceGenerator.On("GetPairPriceBasedOnUSDT", mock.Anything, "GRS-BTC").Once().Return("0.747", nil)
	priceGenerator.On("GetPairPriceBasedOnUSDT", mock.Anything, "BTC-DAI").Once().Return("50000", nil)

	activePairs := []currency.Pair{
		{
			ID:   1,
			Name: "BTC-USDT",
			BasisCoin: currency.Coin{
				ID:   1,
				Code: "USDT",
				Name: "Tether",
			},
			DependentCoin: currency.Coin{
				ID:   2,
				Code: "BTC",
				Name: "Bitcoin",
			},
		},
		{
			ID:   2,
			Name: "ETH-USDT",
			BasisCoin: currency.Coin{
				ID:   1,
				Code: "USDT",
				Name: "Tether",
			},
			DependentCoin: currency.Coin{
				ID:   3,
				Code: "ETH",
				Name: "Ethereum",
			},
		},
		{
			ID:         3,
			Name:       "USDT-DAI",
			ShowDigits: 6,
			BasisCoin: currency.Coin{
				ID:   4,
				Code: "DAI",
				Name: "Dai",
			},
			DependentCoin: currency.Coin{
				ID:   1,
				Code: "USDT",
				Name: "Tether",
			},
		},
		{
			ID:   4,
			Name: "GRS-BTC",
			BasisCoin: currency.Coin{
				ID:   2,
				Code: "BTC",
				Name: "Bitcoin",
			},
			DependentCoin: currency.Coin{
				ID:   5,
				Code: "GRS",
				Name: "Groestlcoin",
			},
		},
		{
			ID:   6,
			Name: "BTC-DAI",
			BasisCoin: currency.Coin{
				ID:   4,
				Code: "DAI",
				Name: "Dai",
			},
			DependentCoin: currency.Coin{
				ID:   2,
				Code: "BTC",
				Name: "Bitcoin",
			},
		},
	}

	pairRepo := new(mocks.PairRepository)
	pairRepo.On("GetActivePairCurrenciesList").Once().Return(activePairs)

	klineService := new(mocks.KlineService)
	favoritePairRepository := new(mocks.FavoritePairRepository)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	currencyService := currency.NewCurrencyService(repo, liveData, priceGenerator, pairRepo, klineService, favoritePairRepository, configs, logger)
	res, statusCode := currencyService.GetPairs()
	assert.Equal(t, http.StatusOK, statusCode)
	pairGroups, ok := res.Data.([]currency.PairGroup)
	if !ok {
		t.Error("can not cast response to struct")
	}
	for _, pg := range pairGroups {
		if pg.Coin == "USDT" {
			assert.Equal(t, int64(1), pg.ID)
			assert.Equal(t, "USDT", pg.Coin)
			assert.Equal(t, "Tether", pg.Name)
			for _, p := range pg.Pairs {
				if p.Name == "BTC-USDT" {
					assert.Equal(t, int64(1), p.ID)
					assert.Equal(t, "BTC", p.DependentCoin)
					assert.Equal(t, "Bitcoin", p.DependentName)
					assert.Equal(t, int64(2), p.DependentID)
					assert.Equal(t, "50000", p.Price)
					assert.Equal(t, "2.1", p.Percentage)
					assert.Equal(t, "50000", p.EquivalentPrice)
				}

				if p.Name == "ETH-USDT" {
					assert.Equal(t, int64(2), p.ID)
					assert.Equal(t, "ETH", p.DependentCoin)
					assert.Equal(t, "Ethereum", p.DependentName)
					assert.Equal(t, int64(3), p.DependentID)
					assert.Equal(t, "2000", p.Price)
					assert.Equal(t, "1.1", p.Percentage)
					assert.Equal(t, "2000", p.EquivalentPrice)
				}
			}
		}

		if pg.Coin == "BTC" {
			assert.Equal(t, int64(2), pg.ID)
			assert.Equal(t, "BTC", pg.Coin)
			assert.Equal(t, "Bitcoin", pg.Name)
			for _, p := range pg.Pairs {
				if p.Name == "GRS-BTC" {
					assert.Equal(t, int64(4), p.ID)
					assert.Equal(t, "GRS", p.DependentCoin)
					assert.Equal(t, "Groestlcoin", p.DependentName)
					assert.Equal(t, int64(5), p.DependentID)
					assert.Equal(t, "0.000013", p.Price)
					assert.Equal(t, "0.1", p.Percentage)
					assert.Equal(t, "0.747", p.EquivalentPrice)
				}
			}

		}

		if pg.Coin == "DAI" {
			assert.Equal(t, int64(4), pg.ID)
			assert.Equal(t, "DAI", pg.Coin)
			assert.Equal(t, "Dai", pg.Name)
			for _, p := range pg.Pairs {
				if p.Name == "USDT-DAI" {
					assert.Equal(t, int64(3), p.ID)
					assert.Equal(t, "USDT", p.DependentCoin)
					assert.Equal(t, "Tether", p.DependentName)
					assert.Equal(t, int64(1), p.DependentID)
					assert.Equal(t, "0.99", p.Price)
					assert.Equal(t, "0.01", p.Percentage)
					assert.Equal(t, "0.99", p.EquivalentPrice)
				}

				if p.Name == "BTC-DAI" {
					assert.Equal(t, int64(6), p.ID)
					assert.Equal(t, "BTC", p.DependentCoin)
					assert.Equal(t, "Bitcoin", p.DependentName)
					assert.Equal(t, int64(2), p.DependentID)
					assert.Equal(t, "50000", p.Price)
					assert.Equal(t, "2.1", p.Percentage)
					assert.Equal(t, "50000", p.EquivalentPrice)
				}
			}

		}
	}

	pairRepo.AssertExpectations(t)
	liveData.AssertExpectations(t)
	priceGenerator.AssertExpectations(t)
}

func TestService_GetPairByID(t *testing.T) {
	repo := new(mocks.CurrencyRepository)
	liveData := new(mocks.LiveData)
	priceGenerator := new(mocks.PriceGenerator)
	pairRepo := new(mocks.PairRepository)
	pairRepo.On("GetPairByID", int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		pair := args.Get(1).(*currency.Pair)
		pair.Name = "BTC-USDT"

	})

	klineService := new(mocks.KlineService)
	favoritePairRepository := new(mocks.FavoritePairRepository)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	currencyService := currency.NewCurrencyService(repo, liveData, priceGenerator, pairRepo, klineService, favoritePairRepository, configs, logger)
	pair, err := currencyService.GetPairByID(int64(1))

	assert.Nil(t, err)
	assert.Equal(t, "BTC-USDT", pair.Name)

	pairRepo.AssertExpectations(t)
}

func TestService_GetPairByName(t *testing.T) {
	repo := new(mocks.CurrencyRepository)
	liveData := new(mocks.LiveData)
	priceGenerator := new(mocks.PriceGenerator)
	pairRepo := new(mocks.PairRepository)
	pairRepo.On("GetPairByName", "BTC-USDT", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		pair := args.Get(1).(*currency.Pair)
		pair.ID = 1

	})

	klineService := new(mocks.KlineService)
	favoritePairRepository := new(mocks.FavoritePairRepository)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	currencyService := currency.NewCurrencyService(repo, liveData, priceGenerator, pairRepo, klineService, favoritePairRepository, configs, logger)
	pair, err := currencyService.GetPairByName("BTC-USDT")

	assert.Nil(t, err)
	assert.Equal(t, int64(1), pair.ID)

	pairRepo.AssertExpectations(t)
}

func TestService_GetCoinByCode(t *testing.T) {
	repo := new(mocks.CurrencyRepository)
	repo.On("GetCoinByCode", "BTC", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		coin := args.Get(1).(*currency.Coin)
		coin.ID = 1

	})
	liveData := new(mocks.LiveData)
	priceGenerator := new(mocks.PriceGenerator)
	pairRepo := new(mocks.PairRepository)

	klineService := new(mocks.KlineService)
	favoritePairRepository := new(mocks.FavoritePairRepository)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	currencyService := currency.NewCurrencyService(repo, liveData, priceGenerator, pairRepo, klineService, favoritePairRepository, configs, logger)
	coin, err := currencyService.GetCoinByCode("BTC")

	assert.Nil(t, err)
	assert.Equal(t, int64(1), coin.ID)

	repo.AssertExpectations(t)
}

func TestService_GetActiveCurrencies(t *testing.T) {
	repo := new(mocks.CurrencyRepository)
	coins := []currency.Coin{
		{
			ID:   1,
			Name: "Bitcoin",
			Code: "BTC",
		},
		{
			ID:   2,
			Name: "Ethereum",
			Code: "ETH",
		},
		{
			ID:   3,
			Name: "Tether",
			Code: "USDT",
		},
	}
	repo.On("GetActiveCoins").Once().Return(coins)
	liveData := new(mocks.LiveData)
	priceGenerator := new(mocks.PriceGenerator)
	pairRepo := new(mocks.PairRepository)

	klineService := new(mocks.KlineService)
	favoritePairRepository := new(mocks.FavoritePairRepository)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	currencyService := currency.NewCurrencyService(repo, liveData, priceGenerator, pairRepo, klineService, favoritePairRepository, configs, logger)
	result := currencyService.GetActiveCoins()

	assert.Equal(t, 3, len(result))

	for _, c := range result {
		switch c.ID {
		case int64(1):
			assert.Equal(t, "Bitcoin", c.Name)
			assert.Equal(t, "BTC", c.Code)
		case int64(2):
			assert.Equal(t, "Ethereum", c.Name)
			assert.Equal(t, "ETH", c.Code)
		case int64(3):
			assert.Equal(t, "Tether", c.Name)
			assert.Equal(t, "USDT", c.Code)
		default:
			t.Fatal("we should not be in default case")
		}

	}

	repo.AssertExpectations(t)
}

func TestService_GetCurrenciesList(t *testing.T) {
	repo := new(mocks.CurrencyRepository)
	coins := []currency.Coin{
		{
			ID:     1,
			Name:   "Bitcoin",
			Code:   "BTC",
			IsMain: true,
		},
		{
			ID:   2,
			Name: "Ethereum",
			Code: "ETH",
		},
		{
			ID:                             3,
			Name:                           "Tether",
			Code:                           "USDT",
			BlockchainNetwork:              sql.NullString{String: "ETH", Valid: true},
			CompletedNetworkName:           sql.NullString{String: "Ethereum(ETH) ERC20", Valid: true},
			OtherBlockchainNetworksConfigs: sql.NullString{String: `[{"code":"TRX","supportsWithdraw":true,"supportsDeposit":true,"completedNetworkName":"Tron(TRX) trc20"}]`, Valid: true},
		},
	}
	repo.On("GetCoinsAlphabetically").Once().Return(coins)
	liveData := new(mocks.LiveData)
	priceGenerator := new(mocks.PriceGenerator)
	pairRepo := new(mocks.PairRepository)

	klineService := new(mocks.KlineService)
	favoritePairRepository := new(mocks.FavoritePairRepository)
	configs := new(mocks.Configs)
	configs.On("GetImagePath").Once().Return("")
	logger := new(mocks.Logger)
	currencyService := currency.NewCurrencyService(repo, liveData, priceGenerator, pairRepo, klineService, favoritePairRepository, configs, logger)
	res, statusCode := currencyService.GetCurrenciesList()

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	result, ok := res.Data.(currency.GetCurrenciesResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}
	assert.Equal(t, 3, len(result.Coins))
	assert.Equal(t, 1, len(result.MainCoins))

	for _, c := range result.Coins {
		switch c.ID {
		case int64(1):
			assert.Equal(t, "Bitcoin", c.Name)
			assert.Equal(t, "BTC", c.Code)
			assert.Equal(t, "Bitcoin(BTC)", c.CompletedNetworkName)
		case int64(2):
			assert.Equal(t, "Ethereum", c.Name)
			assert.Equal(t, "ETH", c.Code)
			assert.Equal(t, "Ethereum(ETH)", c.CompletedNetworkName)
		case int64(3):
			assert.Equal(t, "Tether", c.Name)
			assert.Equal(t, "USDT", c.Code)
			assert.Equal(t, "ETH", c.MainNetwork)
			assert.Equal(t, "Ethereum(ETH) ERC20", c.CompletedNetworkName)
			assert.Equal(t, "TRX", c.OtherBlockChainNetworks[0].Code)
			assert.Equal(t, "Tron(TRX) trc20", c.OtherBlockChainNetworks[0].Name)
			assert.Equal(t, "Tron(TRX) trc20", c.OtherBlockChainNetworks[0].CompletedNetworkName)

		default:
			t.Fatal("we should not be in default case")
		}

	}

	assert.Equal(t, int64(1), result.MainCoins[0].ID)
	assert.Equal(t, "BTC", result.MainCoins[0].Code)
	assert.Equal(t, "Bitcoin", result.MainCoins[0].Name)

	repo.AssertExpectations(t)
	configs.AssertExpectations(t)

}

func TestService_GetPairsStatistic(t *testing.T) {
	repo := new(mocks.CurrencyRepository)

	pairPricesData := []livedata.RedisPairPriceData{
		{
			PairName:   "BTC-USDT",
			Price:      "50000",
			Percentage: "2.1",
		},
		{
			PairName:   "ETH-USDT",
			Price:      "2000",
			Percentage: "1.1",
		},
		{
			PairName:   "GRS-BTC",
			Price:      "0.000013",
			Percentage: "0.1",
		},
		{
			PairName:   "USDT-DAI",
			Price:      "0.99",
			Percentage: "0.01",
		},
		{
			PairName:   "BTC-DAI",
			Price:      "50000",
			Percentage: "2.1",
		},
	}
	liveData := new(mocks.LiveData)

	liveData.On("GetPairsPriceData", mock.Anything, mock.Anything).Once().Return(pairPricesData, nil)

	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetPairPriceBasedOnUSDT", mock.Anything, "BTC-USDT").Once().Return("50000", nil)
	priceGenerator.On("GetPairPriceBasedOnUSDT", mock.Anything, "ETH-USDT").Once().Return("2000", nil)
	priceGenerator.On("GetPairPriceBasedOnUSDT", mock.Anything, "USDT-DAI").Once().Return("0.99", nil)
	priceGenerator.On("GetPairPriceBasedOnUSDT", mock.Anything, "GRS-BTC").Once().Return("0.747", nil)
	priceGenerator.On("GetPairPriceBasedOnUSDT", mock.Anything, "BTC-DAI").Once().Return("50000", nil)

	pairs := []currency.Pair{
		{
			ID:   1,
			Name: "BTC-USDT",
			BasisCoin: currency.Coin{
				ID:   1,
				Code: "USDT",
				Name: "Tether",
			},
			DependentCoin: currency.Coin{
				ID:   2,
				Code: "BTC",
				Name: "Bitcoin",
			},
		},
		{
			ID:   2,
			Name: "ETH-USDT",
			BasisCoin: currency.Coin{
				ID:   1,
				Code: "USDT",
				Name: "Tether",
			},
			DependentCoin: currency.Coin{
				ID:   3,
				Code: "ETH",
				Name: "Ethereum",
			},
		},
		{
			ID:         3,
			Name:       "USDT-DAI",
			ShowDigits: 6,
			BasisCoin: currency.Coin{
				ID:   4,
				Code: "DAI",
				Name: "Dai",
			},
			DependentCoin: currency.Coin{
				ID:   1,
				Code: "USDT",
				Name: "Tether",
			},
		},
		{
			ID:   4,
			Name: "GRS-BTC",
			BasisCoin: currency.Coin{
				ID:   2,
				Code: "BTC",
				Name: "Bitcoin",
			},
			DependentCoin: currency.Coin{
				ID:   5,
				Code: "GRS",
				Name: "Groestlcoin",
			},
		},
		{
			ID:   6,
			Name: "BTC-DAI",
			BasisCoin: currency.Coin{
				ID:   4,
				Code: "DAI",
				Name: "Dai",
			},
			DependentCoin: currency.Coin{
				ID:   2,
				Code: "BTC",
				Name: "Bitcoin",
			},
		},
	}

	pairRepo := new(mocks.PairRepository)
	pairRepo.On("GetPairsByName", []string{"BTC-USDT", "ETH-USDT", "GRS-BTC", "USDT-DAI", "BTC-DAI"}).Once().Return(pairs)

	endTime, err := time.Parse("2006-01-02", "2021-06-06")

	if err != nil {
		t.Error("can not cast response to struct")
	}

	trends := []candle.CandleTrend{
		{
			Pair:      "BTC-USDT",
			Price:     "35000",
			StartTime: "",
			EndTime:   endTime.Format("2006-01-02 15:04:05"),
		},
		{
			Pair:      "ETH-USDT",
			Price:     "2000",
			StartTime: "",
			EndTime:   endTime.Format("2006-01-02 15:04:05"),
		},
	}
	klineService := new(mocks.KlineService)
	klineService.On("GetKlineTrends", []string{"BTC-USDT", "ETH-USDT", "GRS-BTC", "USDT-DAI", "BTC-DAI"}, currency.Timeframe1hour, mock.Anything, mock.Anything).Once().Return(trends, nil)

	favoritePairRepository := new(mocks.FavoritePairRepository)
	configs := new(mocks.Configs)
	configs.On("GetImagePath").Once().Return("")
	logger := new(mocks.Logger)
	currencyService := currency.NewCurrencyService(repo, liveData, priceGenerator, pairRepo, klineService, favoritePairRepository, configs, logger)
	params := currency.GetPairsStatisticParams{
		PairNames: "BTC-USDT|ETH-USDT|GRS-BTC|USDT-DAI|BTC-DAI",
	}

	res, statusCode := currencyService.GetPairsStatistic(params)
	assert.Equal(t, http.StatusOK, statusCode)
	resultData, ok := res.Data.(map[string][]currency.PairStatisticGroup)
	if !ok {
		t.Error("can not cast response to struct")
	}

	pairStatisticGroups := resultData["pairs"]

	for _, pg := range pairStatisticGroups {
		if pg.Coin == "USDT" {
			assert.Equal(t, int64(1), pg.ID)
			assert.Equal(t, "USDT", pg.Coin)
			assert.Equal(t, "Tether", pg.Name)
			for _, p := range pg.Pairs {
				if p.Name == "BTC-USDT" {
					assert.Equal(t, int64(1), p.ID)
					assert.Equal(t, "BTC", p.DependentCoin)
					assert.Equal(t, "Bitcoin", p.DependentName)
					assert.Equal(t, int64(2), p.DependentID)
					assert.Equal(t, "50000", p.Price)
					assert.Equal(t, "2.1", p.Percentage)
					assert.Equal(t, "50000", p.EquivalentPrice)
					assert.Equal(t, "35000", p.TrendData[0].Price)
					assert.Equal(t, "2021-06-06 00:00:00", p.TrendData[0].Time)
				}

				if p.Name == "ETH-USDT" {
					assert.Equal(t, int64(2), p.ID)
					assert.Equal(t, "ETH", p.DependentCoin)
					assert.Equal(t, "Ethereum", p.DependentName)
					assert.Equal(t, int64(3), p.DependentID)
					assert.Equal(t, "2000", p.Price)
					assert.Equal(t, "1.1", p.Percentage)
					assert.Equal(t, "2000", p.EquivalentPrice)
					assert.Equal(t, "2000", p.TrendData[0].Price)
					assert.Equal(t, "2021-06-06 00:00:00", p.TrendData[0].Time)
				}
			}
		}

		if pg.Coin == "BTC" {
			assert.Equal(t, int64(2), pg.ID)
			assert.Equal(t, "BTC", pg.Coin)
			assert.Equal(t, "Bitcoin", pg.Name)
			for _, p := range pg.Pairs {
				if p.Name == "GRS-BTC" {
					assert.Equal(t, int64(4), p.ID)
					assert.Equal(t, "GRS", p.DependentCoin)
					assert.Equal(t, "Groestlcoin", p.DependentName)
					assert.Equal(t, int64(5), p.DependentID)
					assert.Equal(t, "0.000013", p.Price)
					assert.Equal(t, "0.1", p.Percentage)
					assert.Equal(t, "0.747", p.EquivalentPrice)
				}
			}
		}

		if pg.Coin == "DAI" {
			assert.Equal(t, int64(4), pg.ID)
			assert.Equal(t, "DAI", pg.Coin)
			assert.Equal(t, "Dai", pg.Name)
			for _, p := range pg.Pairs {
				if p.Name == "USDT-DAI" {
					assert.Equal(t, int64(3), p.ID)
					assert.Equal(t, "USDT", p.DependentCoin)
					assert.Equal(t, "Tether", p.DependentName)
					assert.Equal(t, int64(1), p.DependentID)
					assert.Equal(t, "0.99", p.Price)
					assert.Equal(t, "0.01", p.Percentage)
					assert.Equal(t, "0.99", p.EquivalentPrice)
				}

				if p.Name == "BTC-DAI" {
					assert.Equal(t, int64(6), p.ID)
					assert.Equal(t, "BTC", p.DependentCoin)
					assert.Equal(t, "Bitcoin", p.DependentName)
					assert.Equal(t, int64(2), p.DependentID)
					assert.Equal(t, "50000", p.Price)
					assert.Equal(t, "2.1", p.Percentage)
					assert.Equal(t, "50000", p.EquivalentPrice)
				}
			}

		}
	}

	pairRepo.AssertExpectations(t)
	liveData.AssertExpectations(t)
	priceGenerator.AssertExpectations(t)
	klineService.AssertExpectations(t)

}

func TestService_AddOrRemoveFavoritePair_Add(t *testing.T) {
	repo := new(mocks.CurrencyRepository)
	liveData := new(mocks.LiveData)
	priceGenerator := new(mocks.PriceGenerator)

	pairRepo := new(mocks.PairRepository)
	pair := &currency.Pair{}
	pairRepo.On("GetPairByID", int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		pair = args.Get(1).(*currency.Pair)
		pair.ID = 1
	})

	klineService := new(mocks.KlineService)

	favoritePairRepository := new(mocks.FavoritePairRepository)
	favoritePairRepository.On("GetFavoritePair", 1, int64(1), mock.Anything).Once().Return(gorm.ErrRecordNotFound)
	favoritePairRepository.On("Create", mock.Anything).Once().Return(nil)

	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	currencyService := currency.NewCurrencyService(repo, liveData, priceGenerator, pairRepo, klineService, favoritePairRepository, configs, logger)
	user := &user.User{
		ID: 1,
	}
	params := currency.FavoriteParams{
		PairID: 1,
		Action: "add",
	}
	currencyService.AddOrRemoveFavoritePair(user, params)

	pairRepo.AssertExpectations(t)
	favoritePairRepository.AssertExpectations(t)
}

func TestService_AddOrRemoveFavoritePair_Add_AlreadyAdded(t *testing.T) {
	repo := new(mocks.CurrencyRepository)
	liveData := new(mocks.LiveData)
	priceGenerator := new(mocks.PriceGenerator)

	pairRepo := new(mocks.PairRepository)
	pair := &currency.Pair{}
	pairRepo.On("GetPairByID", int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		pair = args.Get(1).(*currency.Pair)
		pair.ID = 1
	})

	klineService := new(mocks.KlineService)

	favoritePairRepository := new(mocks.FavoritePairRepository)
	favoritePairRepository.On("GetFavoritePair", 1, int64(1), mock.Anything).Once().Return(nil)

	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	currencyService := currency.NewCurrencyService(repo, liveData, priceGenerator, pairRepo, klineService, favoritePairRepository, configs, logger)
	user := &user.User{
		ID: 1,
	}
	params := currency.FavoriteParams{
		PairID: 1,
		Action: "add",
	}
	currencyService.AddOrRemoveFavoritePair(user, params)

	pairRepo.AssertExpectations(t)
	favoritePairRepository.AssertExpectations(t)
}

func TestService_AddOrRemoveFavoritePair_Remove(t *testing.T) {
	repo := new(mocks.CurrencyRepository)
	liveData := new(mocks.LiveData)
	priceGenerator := new(mocks.PriceGenerator)

	pairRepo := new(mocks.PairRepository)
	pair := &currency.Pair{}
	pairRepo.On("GetPairByID", int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		pair = args.Get(1).(*currency.Pair)
		pair.ID = 1
	})

	klineService := new(mocks.KlineService)

	favoritePairRepository := new(mocks.FavoritePairRepository)
	favoritePairRepository.On("GetFavoritePair", 1, int64(1), mock.Anything).Once().Return(nil)
	favoritePairRepository.On("Delete", mock.Anything).Once().Return(nil)

	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	currencyService := currency.NewCurrencyService(repo, liveData, priceGenerator, pairRepo, klineService, favoritePairRepository, configs, logger)
	user := &user.User{
		ID: 1,
	}
	params := currency.FavoriteParams{
		PairID: 1,
		Action: "remove",
	}
	currencyService.AddOrRemoveFavoritePair(user, params)

	pairRepo.AssertExpectations(t)
	favoritePairRepository.AssertExpectations(t)
}

func TestService_AddOrRemoveFavoritePair_Remove_AlreadyRemoved(t *testing.T) {
	repo := new(mocks.CurrencyRepository)
	liveData := new(mocks.LiveData)
	priceGenerator := new(mocks.PriceGenerator)

	pairRepo := new(mocks.PairRepository)
	pair := &currency.Pair{}
	pairRepo.On("GetPairByID", int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		pair = args.Get(1).(*currency.Pair)
		pair.ID = 1
	})

	klineService := new(mocks.KlineService)

	favoritePairRepository := new(mocks.FavoritePairRepository)
	favoritePairRepository.On("GetFavoritePair", 1, int64(1), mock.Anything).Once().Return(gorm.ErrRecordNotFound)

	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	currencyService := currency.NewCurrencyService(repo, liveData, priceGenerator, pairRepo, klineService, favoritePairRepository, configs, logger)
	user := &user.User{
		ID: 1,
	}
	params := currency.FavoriteParams{
		PairID: 1,
		Action: "remove",
	}
	currencyService.AddOrRemoveFavoritePair(user, params)

	pairRepo.AssertExpectations(t)
	favoritePairRepository.AssertExpectations(t)
}

func TestService_GetFavoritePairs(t *testing.T) {
	repo := new(mocks.CurrencyRepository)
	liveData := new(mocks.LiveData)
	priceGenerator := new(mocks.PriceGenerator)
	pairRepo := new(mocks.PairRepository)
	klineService := new(mocks.KlineService)

	favoritePairs := []currency.FavoritePairQueryFields{
		{
			PairID:   1,
			PairName: "BTC-USDT",
		},
		{
			PairID:   2,
			PairName: "ETH-USDT",
		},
	}
	favoritePairRepository := new(mocks.FavoritePairRepository)
	favoritePairRepository.On("GetUserFavoritePairs", 1).Once().Return(favoritePairs)

	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	currencyService := currency.NewCurrencyService(repo, liveData, priceGenerator, pairRepo, klineService, favoritePairRepository, configs, logger)
	user := &user.User{
		ID: 1,
	}
	res, statusCode := currencyService.GetFavoritePairs(user)
	assert.Equal(t, http.StatusOK, statusCode)
	favoritePairsResult, ok := res.Data.([]currency.GetFavoritePairsResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	assert.Equal(t, 2, len(favoritePairsResult))

	assert.Equal(t, int64(1), favoritePairsResult[0].ID)
	assert.Equal(t, "BTC-USDT", favoritePairsResult[0].Name)
	assert.Equal(t, int64(2), favoritePairsResult[1].ID)
	assert.Equal(t, "ETH-USDT", favoritePairsResult[1].Name)

	favoritePairRepository.AssertExpectations(t)
}

func TestService_GetPairRatio_SecondOneIsUSDT(t *testing.T) {
	repo := new(mocks.CurrencyRepository)
	firstCoin := &currency.Coin{}
	repo.On("GetCoinByCode", "BTC", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		firstCoin = args.Get(1).(*currency.Coin)
		firstCoin.Code = "BTC"
		firstCoin.ID = 1

	})

	secondCoin := &currency.Coin{}
	repo.On("GetCoinByCode", "USDT", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		secondCoin = args.Get(1).(*currency.Coin)
		secondCoin.ID = 2
	})

	liveData := new(mocks.LiveData)
	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetAmountBasedOnUSDT", context.Background(), "BTC", "1.0").Once().Return("40000", nil)

	pairRepo := new(mocks.PairRepository)
	klineService := new(mocks.KlineService)
	favoritePairRepository := new(mocks.FavoritePairRepository)

	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	currencyService := currency.NewCurrencyService(repo, liveData, priceGenerator, pairRepo, klineService, favoritePairRepository, configs, logger)
	params := currency.GetPairRatioParams{
		PairName: "BTC-USDT",
	}
	res, statusCode := currencyService.GetPairRatio(params)
	assert.Equal(t, http.StatusOK, statusCode)
	pairRatioResult, ok := res.Data.(map[string]float64)
	if !ok {
		t.Error("can not cast response to struct")
	}

	assert.Equal(t, float64(40000), pairRatioResult["ratio"])

	repo.AssertExpectations(t)
	priceGenerator.AssertExpectations(t)
}

func TestService_GetPairRatio_FirstOneIsUSDT(t *testing.T) {
	repo := new(mocks.CurrencyRepository)
	firstCoin := &currency.Coin{}
	repo.On("GetCoinByCode", "USDT", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		firstCoin = args.Get(1).(*currency.Coin)
		firstCoin.Code = "USDT"
		firstCoin.ID = 2
	})

	secondCoin := &currency.Coin{}
	repo.On("GetCoinByCode", "BTC", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		secondCoin = args.Get(1).(*currency.Coin)
		secondCoin.ID = 1
		secondCoin.Code = "BTC"
	})

	liveData := new(mocks.LiveData)
	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetAmountBasedOnUSDT", context.Background(), "BTC", "1.0").Once().Return("40000", nil)

	pairRepo := new(mocks.PairRepository)
	klineService := new(mocks.KlineService)
	favoritePairRepository := new(mocks.FavoritePairRepository)

	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	currencyService := currency.NewCurrencyService(repo, liveData, priceGenerator, pairRepo, klineService, favoritePairRepository, configs, logger)
	params := currency.GetPairRatioParams{
		PairName: "USDT-BTC",
	}
	res, statusCode := currencyService.GetPairRatio(params)
	assert.Equal(t, http.StatusOK, statusCode)
	pairRatioResult, ok := res.Data.(map[string]float64)
	if !ok {
		t.Error("can not cast response to struct")
	}

	assert.Equal(t, float64(0.000025), pairRatioResult["ratio"])

	repo.AssertExpectations(t)
	priceGenerator.AssertExpectations(t)
}

func TestService_GetPairRatio_NonUSDT(t *testing.T) {
	repo := new(mocks.CurrencyRepository)
	firstCoin := &currency.Coin{}
	repo.On("GetCoinByCode", "ETH", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		firstCoin = args.Get(1).(*currency.Coin)
		firstCoin.Code = "ETH"
		firstCoin.ID = 2
	})

	secondCoin := &currency.Coin{}
	repo.On("GetCoinByCode", "BTC", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		secondCoin = args.Get(1).(*currency.Coin)
		secondCoin.ID = 1
		secondCoin.Code = "BTC"
	})

	liveData := new(mocks.LiveData)
	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetAmountBasedOnUSDT", context.Background(), "BTC", "1.0").Once().Return("40000", nil)
	priceGenerator.On("GetAmountBasedOnUSDT", context.Background(), "ETH", "1.0").Once().Return("2000", nil)

	pairRepo := new(mocks.PairRepository)
	klineService := new(mocks.KlineService)
	favoritePairRepository := new(mocks.FavoritePairRepository)

	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	currencyService := currency.NewCurrencyService(repo, liveData, priceGenerator, pairRepo, klineService, favoritePairRepository, configs, logger)
	params := currency.GetPairRatioParams{
		PairName: "ETH-BTC",
	}
	res, statusCode := currencyService.GetPairRatio(params)
	assert.Equal(t, http.StatusOK, statusCode)
	pairRatioResult, ok := res.Data.(map[string]float64)
	if !ok {
		t.Error("can not cast response to struct")
	}

	assert.Equal(t, float64(0.05), pairRatioResult["ratio"])

	repo.AssertExpectations(t)
	priceGenerator.AssertExpectations(t)
}

func TestService_GetPairsList(t *testing.T) {
	repo := new(mocks.CurrencyRepository)
	liveData := new(mocks.LiveData)

	priceGenerator := new(mocks.PriceGenerator)
	activePairs := []currency.Pair{
		{
			ID:   1,
			Name: "BTC-USDT",
			BasisCoin: currency.Coin{
				ID:   1,
				Code: "USDT",
				Name: "Tether",
			},
			DependentCoin: currency.Coin{
				ID:   2,
				Code: "BTC",
				Name: "Bitcoin",
			},
		},
		{
			ID:   2,
			Name: "ETH-USDT",
			BasisCoin: currency.Coin{
				ID:   1,
				Code: "USDT",
				Name: "Tether",
			},
			DependentCoin: currency.Coin{
				ID:   3,
				Code: "ETH",
				Name: "Ethereum",
			},
		},
		{
			ID:         3,
			Name:       "USDT-DAI",
			ShowDigits: 6,
			BasisCoin: currency.Coin{
				ID:   4,
				Code: "DAI",
				Name: "Dai",
			},
			DependentCoin: currency.Coin{
				ID:   1,
				Code: "USDT",
				Name: "Tether",
			},
		},
		{
			ID:   4,
			Name: "GRS-BTC",
			BasisCoin: currency.Coin{
				ID:   2,
				Code: "BTC",
				Name: "Bitcoin",
			},
			DependentCoin: currency.Coin{
				ID:   5,
				Code: "GRS",
				Name: "Groestlcoin",
			},
		},
		{
			ID:   6,
			Name: "BTC-DAI",
			BasisCoin: currency.Coin{
				ID:   4,
				Code: "DAI",
				Name: "Dai",
			},
			DependentCoin: currency.Coin{
				ID:   2,
				Code: "BTC",
				Name: "Bitcoin",
			},
		},
	}

	pairRepo := new(mocks.PairRepository)
	pairRepo.On("GetActivePairCurrenciesList").Once().Return(activePairs)

	klineService := new(mocks.KlineService)
	favoritePairRepository := new(mocks.FavoritePairRepository)

	configs := new(mocks.Configs)
	configs.On("GetImagePath").Once().Return("")

	logger := new(mocks.Logger)
	currencyService := currency.NewCurrencyService(repo, liveData, priceGenerator, pairRepo, klineService, favoritePairRepository, configs, logger)
	res, statusCode := currencyService.GetPairsList()
	assert.Equal(t, http.StatusOK, statusCode)
	pairListGroups, ok := res.Data.([]currency.PairListGroup)
	if !ok {
		t.Error("can not cast response to struct")
	}
	for _, pg := range pairListGroups {
		if pg.Coin == "USDT" {
			assert.Equal(t, int64(1), pg.ID)
			assert.Equal(t, "USDT", pg.Coin)
			assert.Equal(t, "Tether", pg.Name)
			for _, p := range pg.Pairs {
				if p.PairName == "BTC-USDT" {
					assert.Equal(t, int64(1), p.PairID)
					assert.Equal(t, "BTC", p.DependentCoin)
					assert.Equal(t, int64(2), p.DependentID)
					assert.Equal(t, "USDT", p.BasisCoin)

				}

				if p.PairName == "ETH-USDT" {
					assert.Equal(t, int64(2), p.PairID)
					assert.Equal(t, "ETH", p.DependentCoin)
					assert.Equal(t, int64(3), p.DependentID)
					assert.Equal(t, "USDT", p.BasisCoin)
				}
			}
		}

		if pg.Coin == "BTC" {
			assert.Equal(t, int64(2), pg.ID)
			assert.Equal(t, "BTC", pg.Coin)
			assert.Equal(t, "Bitcoin", pg.Name)
			for _, p := range pg.Pairs {
				if p.PairName == "GRS-BTC" {
					assert.Equal(t, int64(4), p.PairID)
					assert.Equal(t, "GRS", p.DependentCoin)
					assert.Equal(t, int64(5), p.DependentID)
					assert.Equal(t, "BTC", p.BasisCoin)
				}
			}

		}

		if pg.Coin == "DAI" {
			assert.Equal(t, int64(4), pg.ID)
			assert.Equal(t, "DAI", pg.Coin)
			assert.Equal(t, "Dai", pg.Name)
			for _, p := range pg.Pairs {
				if p.PairName == "USDT-DAI" {
					assert.Equal(t, int64(3), p.PairID)
					assert.Equal(t, "USDT", p.DependentCoin)
					assert.Equal(t, int64(1), p.DependentID)
					assert.Equal(t, "DAI", p.BasisCoin)
				}

				if p.PairName == "BTC-DAI" {
					assert.Equal(t, int64(6), p.PairID)
					assert.Equal(t, "BTC", p.DependentCoin)
					assert.Equal(t, int64(2), p.DependentID)
					assert.Equal(t, "DAI", p.BasisCoin)
				}
			}

		}
	}

	pairRepo.AssertExpectations(t)
}

func TestService_GetFees(t *testing.T) {
	repo := new(mocks.CurrencyRepository)
	liveData := new(mocks.LiveData)
	priceGenerator := new(mocks.PriceGenerator)
	pairRepo := new(mocks.PairRepository)
	klineService := new(mocks.KlineService)
	favoritePairRepository := new(mocks.FavoritePairRepository)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	configs.On("GetImagePath").Once().Return("http://exchange.test")

	activeCoins := []currency.Coin{
		{
			ID:              1,
			Name:            "Bitcoin",
			Code:            "BTC",
			Image:           "/images/btc.png",
			MinimumWithdraw: "0.01",
			WithdrawalFee: sql.NullFloat64{
				Float64: 0.001,
				Valid:   true,
			},
		},
		{
			ID:              1,
			Name:            "Tether",
			Code:            "USDT",
			Image:           "/images/usdt.png",
			MinimumWithdraw: "10",
			CompletedNetworkName: sql.NullString{
				String: "Ethereum(ETH) ERC20",
				Valid:  true,
			},
			BlockchainNetwork: sql.NullString{
				String: "ETH",
				Valid:  true,
			},
			WithdrawalFee: sql.NullFloat64{
				Float64: 5,
				Valid:   true,
			},
			OtherBlockchainNetworksConfigs: sql.NullString{
				String: `[{"code":"TRX","supportsWithdraw":true,"supportsDeposit":true,"completedNetworkName":"Tron(TRX) trc20","fee":"2.0"}]`,
				Valid:  true,
			},
		},
	}

	repo.On("GetActiveCoins").Once().Return(activeCoins)

	activePairs := []currency.Pair{
		{
			ID:       1,
			Name:     "BTC-USDT",
			MakerFee: 0.001,
			TakerFee: 0.002,
		},
		{
			ID:       2,
			Name:     "ETH-USDT",
			MakerFee: 0.01,
			TakerFee: 0.02,
		},
	}

	pairRepo.On("GetActivePairCurrenciesList").Once().Return(activePairs)

	currencyService := currency.NewCurrencyService(repo, liveData, priceGenerator, pairRepo, klineService, favoritePairRepository, configs, logger)
	res, statusCode := currencyService.GetFees()
	assert.Equal(t, http.StatusOK, statusCode)

	currencyFeeResponse, ok := res.Data.(currency.GetFeesResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	coinsFee := currencyFeeResponse.Coins
	assert.Equal(t, 2, len(coinsFee))

	//
	for _, coinFee := range coinsFee {
		if coinFee.Name == "Tether" {

			assert.Equal(t, "USDT", coinFee.Code)
			assert.Equal(t, "http://exchange.test/images/usdt.png", coinFee.Image)
			assert.Equal(t, "", coinFee.SecondImage)
			assert.Equal(t, float64(10), coinFee.MinWithdraw)

			feesData := coinFee.FeeData

			assert.Equal(t, 2, len(feesData))

			for _, feeData := range feesData {

				if feeData.Network == "ETH" {

					assert.Equal(t, "Ethereum(ETH) ERC20", feeData.CompletedNetworkName)
					assert.Equal(t, float64(5), feeData.WithdrawFee)

				} else if feeData.Network == "TRX" {

					assert.Equal(t, "Tron(TRX) trc20", feeData.CompletedNetworkName)
					assert.Equal(t, float64(2), feeData.WithdrawFee)

				} else {
					t.Error("unexpected network in feeData")
				}

			}

		} else if coinFee.Name == "Bitcoin" {

			assert.Equal(t, "BTC", coinFee.Code)
			assert.Equal(t, "http://exchange.test/images/btc.png", coinFee.Image)
			assert.Equal(t, "", coinFee.SecondImage)
			assert.Equal(t, 0.01, coinFee.MinWithdraw)

			feesData := coinFee.FeeData

			assert.Equal(t, 1, len(feesData))

			for _, feeData := range feesData {

				if feeData.Network == "" {

					assert.Equal(t, "", feeData.CompletedNetworkName)
					assert.Equal(t, 0.001, feeData.WithdrawFee)

				} else {
					t.Error("unexpected network in feeData")
				}

			}

		} else {
			t.Error("unexpected coin-fee response")
		}

	}

	pairsFee := currencyFeeResponse.Pairs
	assert.Equal(t, len(pairsFee), 2)

	for _, pairFee := range pairsFee {

		if pairFee.Name == "BTC-USDT" {

			assert.Equal(t, 0.001, pairFee.MakerFee)
			assert.Equal(t, 0.002, pairFee.TakerFee)

		} else if pairFee.Name == "ETH-USDT" {

			assert.Equal(t, 0.01, pairFee.MakerFee)
			assert.Equal(t, 0.02, pairFee.TakerFee)

		} else {
			t.Error("unexpected pair-fee")
		}

	}

	repo.AssertExpectations(t)
	pairRepo.AssertExpectations(t)
	configs.AssertExpectations(t)
	logger.AssertExpectations(t)
}
