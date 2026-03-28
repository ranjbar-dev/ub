package test

import (
	"bytes"
	"context"
	"encoding/json"
	"exchange-go/internal/api"
	"exchange-go/internal/currency"
	"exchange-go/internal/di"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type CurrencyTests struct {
	*suite.Suite
	httpServer  http.Handler
	db          *gorm.DB
	redisClient *redis.Client
	userActor   *userActor
}

func (t *CurrencyTests) SetupSuite() {
	container := getContainer()
	t.httpServer = container.Get(di.HTTPServer).(api.HTTPServer).GetEngine()
	t.db = getDb()
	t.redisClient = getRedis()
	t.userActor = getUserActor()

	ctx := context.Background()

	t.redisClient.HMSet(ctx, "live_data:pair_currency:BTC-USDT", "price", "50000")
	t.redisClient.HMSet(ctx, "live_data:pair_currency:BTC-USDT", "change_price_percentage", "1.2")

	t.redisClient.HMSet(ctx, "live_data:pair_currency:ETH-USDT", "price", "2000")
	t.redisClient.HMSet(ctx, "live_data:pair_currency:ETH-USDT", "change_price_percentage", "2.2")

	t.redisClient.HMSet(ctx, "live_data:pair_currency:ETH-BTC", "price", "0.1")
	t.redisClient.HMSet(ctx, "live_data:pair_currency:ETH-BTC", "change_price_percentage", "2.3")

	t.redisClient.HMSet(ctx, "live_data:pair_currency:GRS-BTC", "price", "0.000013")
	t.redisClient.HMSet(ctx, "live_data:pair_currency:GRS-BTC", "change_price_percentage", "0.1")

	t.redisClient.HMSet(ctx, "live_data:pair_currency:USDT-DAI", "price", "1.0")
	t.redisClient.HMSet(ctx, "live_data:pair_currency:USDT-DAI", "change_price_percentage", "0.01")

	t.redisClient.HMSet(ctx, "live_data:pair_currency:BTC-DAI", "price", "50000")
	t.redisClient.HMSet(ctx, "live_data:pair_currency:BTC-DAI", "change_price_percentage", "1.2")

}

func (t *CurrencyTests) SetupTest() {

}

func (t *CurrencyTests) TearDownTest() {

}

func (t *CurrencyTests) TearDownSuite() {

}

func (t *CurrencyTests) TestGetPairs() {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/currencies/pairs", nil)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	response := struct {
		Status  bool
		Message string
		Data    []currency.PairGroup
	}{}

	err := json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		t.Fail(err.Error())
	}

	for _, pg := range response.Data {
		if pg.Coin == "USDT" {
			assert.Equal(t.T(), "USDT", pg.Coin)
			assert.Equal(t.T(), "Tether", pg.Name)
			for _, p := range pg.Pairs {
				if p.Name == "BTC-USDT" {
					assert.Equal(t.T(), "BTC", p.DependentCoin)
					assert.Equal(t.T(), "Bitcoin", p.DependentName)
					assert.Equal(t.T(), "50000", p.Price)
					assert.Equal(t.T(), "1.2", p.Percentage)
					assert.Equal(t.T(), "50000.00000000", p.EquivalentPrice)
				}

				if p.Name == "ETH-USDT" {
					assert.Equal(t.T(), "ETH", p.DependentCoin)
					assert.Equal(t.T(), "Ethereum", p.DependentName)
					assert.Equal(t.T(), "2000", p.Price)
					assert.Equal(t.T(), "2.2", p.Percentage)
					assert.Equal(t.T(), "2000.00000000", p.EquivalentPrice)
				}
			}
		}

		if pg.Coin == "BTC" {
			assert.Equal(t.T(), "BTC", pg.Coin)
			assert.Equal(t.T(), "Bitcoin", pg.Name)
			for _, p := range pg.Pairs {
				if p.Name == "GRS-BTC" {
					assert.Equal(t.T(), "GRS", p.DependentCoin)
					assert.Equal(t.T(), "Groestlcoin", p.DependentName)
					assert.Equal(t.T(), "0.000013", p.Price)
					assert.Equal(t.T(), "0.1", p.Percentage)
					assert.Equal(t.T(), "0.65000000", p.EquivalentPrice)
				}

				if p.Name == "ETH-BTC" {
					assert.Equal(t.T(), "ETH", p.DependentCoin)
					assert.Equal(t.T(), "Ethereum", p.DependentName)
					assert.Equal(t.T(), "0.1", p.Price)
					assert.Equal(t.T(), "2.3", p.Percentage)
					assert.Equal(t.T(), "5000.00000000", p.EquivalentPrice)
				}
			}

		}

		if pg.Coin == "DAI" {
			assert.Equal(t.T(), "DAI", pg.Coin)
			assert.Equal(t.T(), "Dai", pg.Name)
			for _, p := range pg.Pairs {
				if p.Name == "USDT-DAI" {
					assert.Equal(t.T(), "USDT", p.DependentCoin)
					assert.Equal(t.T(), "Tether", p.DependentName)
					assert.Equal(t.T(), "1.0", p.Price)
					assert.Equal(t.T(), "0.01", p.Percentage)
					assert.Equal(t.T(), "1.00000000", p.EquivalentPrice)
				}

				if p.Name == "BTC-DAI" {
					assert.Equal(t.T(), "BTC", p.DependentCoin)
					assert.Equal(t.T(), "Bitcoin", p.DependentName)
					assert.Equal(t.T(), "50000", p.Price)
					assert.Equal(t.T(), "1.2", p.Percentage)
					assert.Equal(t.T(), "50000.00000000", p.EquivalentPrice)
				}
			}

		}
	}
}

func (t *CurrencyTests) TestGetCurrenciesList() {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/currencies", nil)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	response := struct {
		Status  bool
		Message string
		Data    currency.GetCurrenciesResponse
	}{}

	err := json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), 1, len(response.Data.MainCoins))
	assert.Equal(t.T(), 6, len(response.Data.Coins))

	//data are inserted in cuurencyseed.go
	for _, c := range response.Data.Coins {
		switch c.Code {
		case "USDT":
			assert.Equal(t.T(), "Tether", c.Name)
			assert.Equal(t.T(), "ETH", c.MainNetwork)
			assert.Equal(t.T(), "Ethereum(ETH) ERC20", c.CompletedNetworkName)
			assert.Equal(t.T(), "TRX", c.OtherBlockChainNetworks[0].Code)
			assert.Equal(t.T(), "Tron (TRX)", c.OtherBlockChainNetworks[0].Name)
			assert.Equal(t.T(), "Tron (TRX)", c.OtherBlockChainNetworks[0].CompletedNetworkName)
		case "BTC":
			assert.Equal(t.T(), "Bitcoin", c.Name)
			assert.Equal(t.T(), "", c.MainNetwork)
			assert.Equal(t.T(), "Bitcoin(BTC)", c.CompletedNetworkName)
		case "ETH":
			assert.Equal(t.T(), "Ethereum", c.Name)
			assert.Equal(t.T(), "", c.MainNetwork)
			assert.Equal(t.T(), "Ethereum(ETH)", c.CompletedNetworkName)
		case "GRS":
			assert.Equal(t.T(), "Groestlcoin", c.Name)
			assert.Equal(t.T(), "", c.MainNetwork)
			assert.Equal(t.T(), "Groestlcoin(GRS)", c.CompletedNetworkName)
		case "DAI":
			assert.Equal(t.T(), "Dai", c.Name)
			assert.Equal(t.T(), "", c.MainNetwork)
			assert.Equal(t.T(), "Dai(DAI)", c.CompletedNetworkName)
		case "TRX":
			assert.Equal(t.T(), "Tron", c.Name)
			assert.Equal(t.T(), "", c.MainNetwork)
			assert.Equal(t.T(), "Tron(TRX)", c.CompletedNetworkName)
		default:
			t.Fail("we should not be in default case")
		}
	}

	for _, c := range response.Data.MainCoins {
		switch c.Code {
		case "BTC":
			assert.Equal(t.T(), "Bitcoin", c.Name)
			assert.Equal(t.T(), "", c.MainNetwork)
		default:
			t.Fail("we should not be in default case")
		}

	}
}

func (t *CurrencyTests) TestGetPairsStatistic() {
	queryParams := url.Values{}
	queryParams.Set("pair_currencies", "BTC-USDT|ETH-USDT")
	paramsString := queryParams.Encode()

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/currencies/pairs-statistic?"+paramsString, nil)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	response := struct {
		Status  bool
		Message string
		Data    map[string][]currency.PairStatisticGroup
	}{}

	err := json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		t.Fail(err.Error())
	}

	for _, pg := range response.Data["pairs"] {
		if pg.Coin == "USDT" {
			assert.Equal(t.T(), "USDT", pg.Coin)
			assert.Equal(t.T(), "Tether", pg.Name)
			for _, p := range pg.Pairs {
				if p.Name == "BTC-USDT" {
					assert.Equal(t.T(), "BTC", p.DependentCoin)
					assert.Equal(t.T(), "Bitcoin", p.DependentName)
					assert.Equal(t.T(), "50000", p.Price)
					assert.Equal(t.T(), "1.2", p.Percentage)
					assert.Equal(t.T(), "50000.00000000", p.EquivalentPrice)

					assert.Equal(t.T(), "30000.00000000", p.TrendData[0].Price)
				}

				if p.Name == "ETH-USDT" {
					assert.Equal(t.T(), "ETH", p.DependentCoin)
					assert.Equal(t.T(), "Ethereum", p.DependentName)
					assert.Equal(t.T(), "2000", p.Price)
					assert.Equal(t.T(), "2.2", p.Percentage)
					assert.Equal(t.T(), "2000.00000000", p.EquivalentPrice)

					assert.Equal(t.T(), "2000.00000000", p.TrendData[0].Price)
				}
			}
		}
	}

}

func (t *CurrencyTests) TestAddOrRemoveFavoritePair() {
	//first we delete all data in user_favorite_pair_currency table
	err := t.db.Where("user_id = ?", t.userActor.ID).Delete(currency.FavoritePair{}).Error
	if err != nil {
		t.Fail(err.Error())
	}

	//add
	res := httptest.NewRecorder()
	data := `{"pair_currency_id":1,"action":"add"}` // is BTC-USDT from currencyseed.go
	body := []byte(data)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/currencies/favorite", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	fp := &currency.FavoritePair{}
	err = t.db.Where("user_id = ?", t.userActor.ID).First(fp).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), int64(1), fp.PairID)

	//remove
	res = httptest.NewRecorder()
	data = `{"pair_currency_id":1,"action":"remove"}` // is BTC-USDT from currencyseed.go
	body = []byte(data)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/currencies/favorite", bytes.NewReader(body))
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	fp = &currency.FavoritePair{}
	err = t.db.Where("user_id = ? and pair_currency_id = ?", t.userActor.ID, int64(1)).First(fp).Error
	assert.Equal(t.T(), gorm.ErrRecordNotFound, err)

}

func (t *CurrencyTests) TestGetFavoritePairs() {
	favoritePairs := []currency.FavoritePair{
		{
			UserID: t.userActor.ID,
			PairID: 1,
		},
		{
			UserID: t.userActor.ID,
			PairID: 2,
		},
	}

	err := t.db.Create(favoritePairs).Error
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/currencies/favorite-pairs", nil)
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	response := struct {
		Status  bool
		Message string
		Data    []currency.GetFavoritePairsResponse
	}{}

	err = json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), 2, len(response.Data))

	assert.Equal(t.T(), int64(1), response.Data[0].ID)
	assert.Equal(t.T(), "BTC-USDT", response.Data[0].Name)

	assert.Equal(t.T(), int64(2), response.Data[1].ID)
	assert.Equal(t.T(), "ETH-USDT", response.Data[1].Name)

}

func (t *CurrencyTests) TestGetPairRatio() {
	queryParams := url.Values{}
	queryParams.Set("pair", "BTC-USDT")
	paramsString := queryParams.Encode()

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/currencies/pairs-ratio?"+paramsString, nil)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	response := struct {
		Status  bool
		Message string
		Data    map[string]float64
	}{}

	err := json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), float64(50000), response.Data["ratio"])

	queryParams = url.Values{}
	queryParams.Set("pair", "USDT-BTC")
	paramsString = queryParams.Encode()
	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/currencies/pairs-ratio?"+paramsString, nil)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	response = struct {
		Status  bool
		Message string
		Data    map[string]float64
	}{}

	err = json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), float64(0.00002), response.Data["ratio"])

	queryParams = url.Values{}
	queryParams.Set("pair", "ETH-BTC")
	paramsString = queryParams.Encode()

	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/currencies/pairs-ratio?"+paramsString, nil)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	response = struct {
		Status  bool
		Message string
		Data    map[string]float64
	}{}

	err = json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), float64(0.04), response.Data["ratio"])
}

func (t *CurrencyTests) TestGetPairsList() {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/currencies/pairs-list", nil)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	response := struct {
		Status  bool
		Message string
		Data    []currency.PairListGroup
	}{}

	err := json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		t.Fail(err.Error())
	}

	for _, pg := range response.Data {
		if pg.Coin == "USDT" {
			assert.Equal(t.T(), "USDT", pg.Coin)
			assert.Equal(t.T(), "Tether", pg.Name)
			for _, p := range pg.Pairs {
				if p.PairName == "BTC-USDT" {
					assert.Equal(t.T(), "BTC", p.DependentCoin)
					assert.Equal(t.T(), "USDT", p.BasisCoin)
				}

				if p.PairName == "ETH-USDT" {
					assert.Equal(t.T(), "ETH", p.DependentCoin)
					assert.Equal(t.T(), "USDT", p.BasisCoin)
				}
			}
		}

		if pg.Coin == "BTC" {
			assert.Equal(t.T(), "BTC", pg.Coin)
			assert.Equal(t.T(), "Bitcoin", pg.Name)
			for _, p := range pg.Pairs {
				if p.PairName == "GRS-BTC" {
					assert.Equal(t.T(), "GRS", p.DependentCoin)
					assert.Equal(t.T(), "BTC", p.BasisCoin)
				}

				if p.PairName == "ETH-BTC" {
					assert.Equal(t.T(), "ETH", p.DependentCoin)
					assert.Equal(t.T(), "BTC", p.BasisCoin)
				}
			}

		}

		if pg.Coin == "DAI" {
			assert.Equal(t.T(), "DAI", pg.Coin)
			assert.Equal(t.T(), "Dai", pg.Name)
			for _, p := range pg.Pairs {
				if p.PairName == "USDT-DAI" {
					assert.Equal(t.T(), "USDT", p.DependentCoin)
					assert.Equal(t.T(), "DAI", p.BasisCoin)
				}

				if p.PairName == "BTC-DAI" {
					assert.Equal(t.T(), "BTC", p.DependentCoin)
					assert.Equal(t.T(), "DAI", p.BasisCoin)
				}
			}

		}
	}

}

func (t *CurrencyTests) TestGetPFees() {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/currencies/fees", nil)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	response := struct {
		Status  bool
		Message string
		Data    currency.GetFeesResponse
	}{}

	err := json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		t.T().Fail()
	}

	assert.IsType(t.T(), []currency.CoinFee{}, response.Data.Coins)
	assert.IsType(t.T(), []currency.PairFee{}, response.Data.Pairs)
}

func TestCurrency(t *testing.T) {
	suite.Run(t, &CurrencyTests{
		Suite: new(suite.Suite),
	})

}
