package currency

import (
	"context"
	"errors"
	"exchange-go/internal/livedata"
	"exchange-go/internal/platform"
	"exchange-go/internal/response"
	"exchange-go/internal/user"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	USDT        = "USDT"
	BTC         = "BTC"
	PairBTCUSDT = "BTC-USDT"

	AggregationStatusStop  = "STOP"
	AggregationStatusRun   = "RUN"
	AggregationStatusPause = "PAUSE"
)

// Service provides the public API for currency and trading pair operations,
// including listings, lookups, statistics, favorites, and fee information.
type Service interface {
	// GetPairs returns all active trading pairs with their current live prices.
	GetPairs() (apiResponse response.APIResponse, statusCode int)
	// GetPairByID retrieves a single trading pair by its database ID.
	GetPairByID(id int64) (Pair, error)
	// GetPairByName retrieves a single trading pair by its name (e.g. "BTC-USDT").
	GetPairByName(name string) (Pair, error)
	// GetCoinByCode retrieves a single coin by its ticker symbol (e.g. "BTC").
	GetCoinByCode(code string) (Coin, error)
	// GetActiveCoins returns all coins that are currently active on the exchange.
	GetActiveCoins() []Coin
	// GetActivePairCurrenciesList returns all active trading pairs with their coin details.
	GetActivePairCurrenciesList() []Pair
	// GetCurrenciesList returns all coins with blockchain network information for the API.
	GetCurrenciesList() (apiResponse response.APIResponse, statusCode int)
	// GetPairsStatistic returns 24-hour statistics including price trends, volume, and
	// percentage changes for trading pairs.
	GetPairsStatistic(params GetPairsStatisticParams) (apiResponse response.APIResponse, statusCode int)
	// AddOrRemoveFavoritePair toggles a trading pair in the user's favorites list.
	AddOrRemoveFavoritePair(u *user.User, params FavoriteParams) (apiResponse response.APIResponse, statusCode int)
	// GetFavoritePairs returns all trading pairs the user has marked as favorites.
	GetFavoritePairs(u *user.User) (apiResponse response.APIResponse, statusCode int)
	// GetPairRatio calculates and returns the exchange ratio between two trading pairs.
	GetPairRatio(params GetPairRatioParams) (apiResponse response.APIResponse, statusCode int)
	// GetPairsList returns a simplified list of all active trading pairs.
	GetPairsList() (apiResponse response.APIResponse, statusCode int)
	// GetFees returns trading (maker/taker) and withdrawal fee schedules.
	GetFees() (apiResponse response.APIResponse, statusCode int)
}

type service struct {
	repository             Repository
	liveData               livedata.Service
	priceGenerator         PriceGenerator
	pairRepository         PairRepository
	klineService           KlineService
	favoritePairRepository FavoritePairRepository
	configs                platform.Configs
	logger                 platform.Logger
}

type GetPairsResponse struct {
	ID              int64  `json:"pairId"`
	Name            string `json:"pairName"`
	DependentCoin   string `json:"dependentCode"`
	DependentName   string `json:"dependentName"`
	DependentID     int64  `json:"dependentId"`
	IsMain          bool   `json:"isMain"`
	Price           string `json:"price"`
	EquivalentPrice string `json:"equivalentPrice"`
	Percentage      string `json:"percent"`
}

type PairGroup struct {
	ID    int64              `json:"id"`
	Coin  string             `json:"code"`
	Name  string             `json:"name"`
	Pairs []GetPairsResponse `json:"pairs"`
}

type GetPairsListResponse struct {
	PairID        int64  `json:"pairId"`
	PairName      string `json:"pairName"`
	DependentCoin string `json:"dependentCode"`
	BasisCoin     string `json:"basisCode"`
	DependentID   int64  `json:"dependentId"`
	ShowDigits    int    `json:"showDigits"`
	Image         string `json:"image"`
}

type PairListGroup struct {
	ID    int64                  `json:"id"`
	Coin  string                 `json:"code"`
	Name  string                 `json:"name"`
	Pairs []GetPairsListResponse `json:"pairs"`
}

type GetPairsStatisticResponse struct {
	ID              int64            `json:"pairId"`
	Name            string           `json:"pairName"`
	BasisCoin       string           `json:"basisCode"`
	BasisSubUnit    int              `json:"basisSubUnit"`
	DependentCoin   string           `json:"dependentCode"`
	DependentName   string           `json:"dependentName"`
	DependentID     int64            `json:"dependentId"`
	IsMain          bool             `json:"isMain"`
	IsFavorite      bool             `json:"isFavorite"`
	Price           string           `json:"price"`
	EquivalentPrice string           `json:"equivalentPrice"`
	Percentage      string           `json:"percent"`
	TrendData       []KlineTrendData `json:"trendData"`
	Volume          string           `json:"volume"`
	Image           string           `json:"image"`
	MakerFee        float64          `json:"makerFee"`
	TakerFee        float64          `json:"takerFee"`
	ShowDigits      int              `json:"showDigits"`
	SubUnit         int              `json:"subUnit"`
	LastUpdate      string           `json:"lastUpdate"`
}

type PairStatisticGroup struct {
	ID      int64                       `json:"id"`
	Coin    string                      `json:"code"`
	Name    string                      `json:"name"`
	SubUnit int                         `json:"subUnit"`
	Pairs   []GetPairsStatisticResponse `json:"pairs"`
}

type AverageDayPrice struct {
	PairID       int64
	Day          string
	AveragePrice string
}

type OtherBlockChainNetworksResponse struct {
	Code                 string `json:"code"`
	Name                 string `json:"name"`
	CompletedNetworkName string `json:"completedNetworkName"`
}

type SingleCoinResponse struct {
	ID                      int64                             `json:"id"`
	Code                    string                            `json:"code"`
	Name                    string                            `json:"name"`
	Image                   string                            `json:"image"`
	SecondImage             string                            `json:"secondImage"`
	BackgroundImage         string                            `json:"backgroundImage"`
	MainNetwork             string                            `json:"mainNetwork"`
	OtherBlockChainNetworks []OtherBlockChainNetworksResponse `json:"otherBlockChainNetworks"`
	ShowDigits              int                               `json:"showDigits"`
	CompletedNetworkName    string                            `json:"completedNetworkName"`
}

type GetCurrenciesResponse struct {
	Coins     []SingleCoinResponse `json:"currencies"`
	MainCoins []SingleCoinResponse `json:"mainCurrencies"`
}

type GetPairsStatisticParams struct {
	PairNames string `form:"pair_currencies"`
}

type FavoriteParams struct {
	PairID int64  `json:"pair_currency_id" binding:"required"`
	Action string `json:"action" binding:"required,oneof='add' 'remove'"`
}

type GetFavoritePairsResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type FavoritePairQueryFields struct {
	PairID   int64
	PairName string
}

type GetPairRatioParams struct {
	PairName string `form:"pair"`
}

type GetFeesResponse struct {
	Coins []CoinFee `json:"coins"`
	Pairs []PairFee `json:"pairs"`
}

type CoinFee struct {
	Name        string        `json:"name"`
	Code        string        `json:"code"`
	FeeData     []CoinFeeData `json:"feeData"`
	MinWithdraw float64       `json:"minWithdraw"`
	Image       string        `json:"image"`
	SecondImage string        `json:"secondImage"`
}

type CoinFeeData struct {
	Network              string  `json:"network"`
	CompletedNetworkName string  `json:"completedNetworkName"`
	WithdrawFee          float64 `json:"withdrawFee"`
}

type PairFee struct {
	Name     string  `json:"name"`
	MakerFee float64 `json:"makerFee"`
	TakerFee float64 `json:"takerFee"`
}

func (s *service) GetCurrenciesList() (apiResponse response.APIResponse, statusCode int) {
	allCoins := make([]SingleCoinResponse, 0)
	mainCoins := make([]SingleCoinResponse, 0)
	coins := s.repository.GetCoinsAlphabetically()
	imagePath := s.configs.GetImagePath()
	for _, coin := range coins {
		coinImage := imagePath + coin.Image
		secondCoinImage := ""
		if coin.SecondImage.Valid {
			secondCoinImage = imagePath + coin.SecondImage.String
		}
		backgroundImage := coinImage //todo this should be image from extra info handle this later
		otherBlockChainNetworks := make([]OtherBlockChainNetworksResponse, 0)
		otherBlockchainNetworksConfigs, err := coin.GetOtherBlockchainNetworksConfigs()
		if err != nil {
			s.logger.Error2("error getting otherBlockchainNetworksconfigs", err,
				zap.String("service", "currencyService"),
				zap.String("method", "GetCurrenciesList"),
				zap.String("coinCode", coin.Code),
			)
		}

		for _, config := range otherBlockchainNetworksConfigs {

			item := OtherBlockChainNetworksResponse{
				Code:                 config.Code,
				Name:                 config.CompletedNetworkName,
				CompletedNetworkName: config.CompletedNetworkName,
			}
			otherBlockChainNetworks = append(otherBlockChainNetworks, item)
		}

		completedNetworkName := ""
		if coin.CompletedNetworkName.Valid {
			completedNetworkName = coin.CompletedNetworkName.String
		} else {
			completedNetworkName = coin.Name + "(" + coin.Code + ")"
		}

		res := SingleCoinResponse{
			ID:                      coin.ID,
			Code:                    coin.Code,
			Name:                    coin.Name,
			Image:                   coinImage,
			SecondImage:             secondCoinImage,
			BackgroundImage:         backgroundImage,
			MainNetwork:             coin.BlockchainNetwork.String,
			OtherBlockChainNetworks: otherBlockChainNetworks,
			ShowDigits:              coin.ShowSubUnit,
			CompletedNetworkName:    completedNetworkName,
		}

		allCoins = append(allCoins, res)
		if coin.IsMain {
			mainCoins = append(mainCoins, res)
		}
	}

	result := GetCurrenciesResponse{
		Coins:     allCoins,
		MainCoins: mainCoins,
	}

	return response.Success(result, "")
}

func (s *service) GetPairs() (apiResponse response.APIResponse, statusCode int) {
	groups := make([]PairGroup, 0)
	var basisIds []int64
	var pairNames []string
	activePairs := s.pairRepository.GetActivePairCurrenciesList()

	for _, pair := range activePairs {
		pairNames = append(pairNames, pair.Name)
	}

	pairsPriceData, _ := s.liveData.GetPairsPriceData(context.Background(), pairNames)
	for _, pair := range activePairs {
		price := ""
		percentage := ""
		for _, pairPriceData := range pairsPriceData {
			if pairPriceData.PairName == pair.Name {
				price = pairPriceData.Price
				percentage = pairPriceData.Percentage
				break
			}
		}
		equivalentPrice, _ := s.priceGenerator.GetPairPriceBasedOnUSDT(context.Background(), pair.Name)
		basisID := pair.BasisCoin.ID
		partialPair := GetPairsResponse{
			ID:              pair.ID,
			Name:            pair.Name,
			DependentCoin:   pair.DependentCoin.Code,
			DependentName:   pair.DependentCoin.Name,
			DependentID:     pair.DependentCoin.ID,
			IsMain:          pair.IsMain,
			Price:           price,
			EquivalentPrice: equivalentPrice,
			Percentage:      percentage,
		}
		if !isInSlice(basisIds, basisID) {
			pairGroup := PairGroup{
				ID:   basisID,
				Coin: pair.BasisCoin.Code,
				Name: pair.BasisCoin.Name,
				Pairs: []GetPairsResponse{
					partialPair,
				},
			}
			groups = append(groups, pairGroup)
			basisIds = append(basisIds, basisID)
		} else {
			for i, g := range groups {
				if g.ID == basisID {
					groups[i].Pairs = append(g.Pairs, partialPair)
					break
				}
			}

		}
	}

	return response.Success(groups, "")
}

func (s *service) GetPairsList() (apiResponse response.APIResponse, statusCode int) {
	groups := make([]PairListGroup, 0)
	var basisIds []int64
	imagePath := s.configs.GetImagePath()
	activePairs := s.pairRepository.GetActivePairCurrenciesList()
	for _, p := range activePairs {
		coinImage := imagePath + p.DependentCoin.Image
		basisID := p.BasisCoin.ID
		partialPair := GetPairsListResponse{
			PairID:        p.ID,
			PairName:      p.Name,
			DependentCoin: p.DependentCoin.Code,
			BasisCoin:     p.BasisCoin.Code,
			DependentID:   p.DependentCoin.ID,
			ShowDigits:    p.ShowDigits,
			Image:         coinImage,
		}
		if !isInSlice(basisIds, basisID) {
			pairGroup := PairListGroup{
				ID:   basisID,
				Coin: p.BasisCoin.Code,
				Name: p.BasisCoin.Name,
				Pairs: []GetPairsListResponse{
					partialPair,
				},
			}
			groups = append(groups, pairGroup)
			basisIds = append(basisIds, basisID)
		} else {
			for i, g := range groups {
				if g.ID == basisID {
					groups[i].Pairs = append(g.Pairs, partialPair)
					break
				}
			}

		}
	}

	return response.Success(groups, "")
}

func (s *service) GetFees() (apiResponse response.APIResponse, statusCode int) {

	getFeesResponse := GetFeesResponse{
		Coins: s.getCoinsFee(),
		Pairs: s.getPairsFee(),
	}

	return response.Success(getFeesResponse, "")
}

func (s *service) getCoinsFee() []CoinFee {
	coinsFee := make([]CoinFee, 0)
	//coins image path
	imagePath := s.configs.GetImagePath()

	coins := s.GetActiveCoins()

	for _, coin := range coins {

		var minWithdraw float64
		minWithdraw, err := strconv.ParseFloat(coin.MinimumWithdraw, 64)
		if err != nil {
			minWithdraw = float64(0)
		}

		coinFeeData := make([]CoinFeeData, 0)

		if !coin.BlockchainNetwork.Valid {
			//Coin Is not Token
			feeData := CoinFeeData{
				Network:              "",
				CompletedNetworkName: "",
				WithdrawFee:          coin.WithdrawalFee.Float64,
			}

			coinFeeData = append(coinFeeData, feeData)

		} else {
			//Coin Is Token

			//Add main token network first
			feeData := CoinFeeData{
				Network:              coin.BlockchainNetwork.String,
				CompletedNetworkName: coin.CompletedNetworkName.String,
				WithdrawFee:          coin.WithdrawalFee.Float64,
			}

			coinFeeData = append(coinFeeData, feeData)

			//Check if it has other networks then add them
			otherBlockChainNetworksConfigs, err := coin.GetOtherBlockchainNetworksConfigs()
			if err != nil {
				//log this error
				s.logger.Error2("error getting otherBlockChainNetworksConfigs", err,
					zap.String("service", "currencyService"),
					zap.String("method", "getCoinFee"),
					zap.String("coinCode", coin.Code),
				)
			}

			for _, otherNetwork := range otherBlockChainNetworksConfigs {
				var otherNetworkWithdrawFee float64
				otherNetworkWithdrawFee, err := strconv.ParseFloat(otherNetwork.Fee, 64)
				if err != nil {
					s.logger.Error2("error getting otherBlockChainNetworks withdraw fee", err,
						zap.String("service", "currencyService"),
						zap.String("method", "getCoinFee"),
						zap.String("coinCode", coin.Code),
					)
					otherNetworkWithdrawFee = float64(0)
				}

				otherNetworkFeeData := CoinFeeData{
					Network:              otherNetwork.Code,
					CompletedNetworkName: otherNetwork.CompletedNetworkName,
					WithdrawFee:          otherNetworkWithdrawFee,
				}

				coinFeeData = append(coinFeeData, otherNetworkFeeData)
			}
		}

		coinImage := imagePath + coin.Image
		secondCoinImage := ""
		if coin.SecondImage.Valid {
			secondCoinImage = imagePath + coin.SecondImage.String
		}

		coinFee := CoinFee{
			Name:        coin.Name,
			Code:        coin.Code,
			FeeData:     coinFeeData,
			MinWithdraw: minWithdraw,
			Image:       coinImage,
			SecondImage: secondCoinImage,
		}

		coinsFee = append(coinsFee, coinFee)
	}

	return coinsFee
}

func (s *service) getPairsFee() []PairFee {
	pairsFee := make([]PairFee, 0)

	pairs := s.GetActivePairCurrenciesList()

	for _, pair := range pairs {
		pairFee := PairFee{
			Name:     pair.Name,
			MakerFee: pair.MakerFee,
			TakerFee: pair.TakerFee,
		}

		pairsFee = append(pairsFee, pairFee)
	}

	return pairsFee
}

func (s *service) GetPairsStatistic(params GetPairsStatisticParams) (apiResponse response.APIResponse, statusCode int) {
	pairNames := strings.Split(params.PairNames, "|")

	pairs := s.pairRepository.GetPairsByName(pairNames)
	from := time.Now().Add(-1 * 100 * time.Hour)
	to := time.Now()
	klineTrends, err := s.klineService.GetKlineTrends(pairNames, Timeframe1hour, from, to)
	if err != nil {
		s.logger.Warn("can not get kline trends",
			zap.Error(err),
			zap.String("service", "currencyService"),
			zap.String("method", "GetPairsStatistic"),
		)
	}
	var basisIds []int64
	ctx := context.Background()
	pairsPriceData, _ := s.liveData.GetPairsPriceData(ctx, pairNames)

	imagePath := s.configs.GetImagePath()
	groups := make([]PairStatisticGroup, 0)
	for _, p := range pairs {
		//klineTrends, err := s.klineService.GetKlineTrends(p.Name, TimeFrame_1HOUR, from, to )
		trendsData := make([]KlineTrendData, 0)
		for _, k := range klineTrends {
			if k.Pair == p.Name {
				td := KlineTrendData{
					PairID: p.ID,
					Price:  k.Price,
					Time:   k.EndTime,
				}
				trendsData = append(trendsData, td)
			}
		}
		price := ""
		percentage := ""
		volume := ""
		for _, pairPriceData := range pairsPriceData {
			if pairPriceData.PairName == p.Name {
				price = pairPriceData.Price
				percentage = pairPriceData.Percentage
				volume = pairPriceData.Volume
				break
			}
		}
		equivalentPrice, _ := s.priceGenerator.GetPairPriceBasedOnUSDT(ctx, p.Name)
		basisID := p.BasisCoin.ID
		image := imagePath + p.DependentCoin.Image
		partialPair := GetPairsStatisticResponse{
			ID:              p.ID,
			Name:            p.Name,
			BasisCoin:       p.BasisCoin.Code,
			BasisSubUnit:    p.BasisCoin.SubUnit,
			DependentCoin:   p.DependentCoin.Code,
			DependentName:   p.DependentCoin.Name,
			DependentID:     p.DependentCoin.ID,
			IsMain:          p.IsMain,
			IsFavorite:      false,
			Price:           price,
			EquivalentPrice: equivalentPrice,
			Percentage:      percentage,
			TrendData:       trendsData,
			Volume:          volume,
			Image:           image,
			MakerFee:        p.MakerFee,
			TakerFee:        p.TakerFee,
			ShowDigits:      p.BasisCoin.ShowSubUnit,
			SubUnit:         p.DependentCoin.SubUnit,
			LastUpdate:      time.Now().Format("2006-01-02 15:04:05"),
		}
		if !isInSlice(basisIds, basisID) {
			pairGroup := PairStatisticGroup{
				ID:      basisID,
				Coin:    p.BasisCoin.Code,
				Name:    p.BasisCoin.Name,
				SubUnit: p.BasisCoin.SubUnit,
				Pairs: []GetPairsStatisticResponse{
					partialPair,
				},
			}
			groups = append(groups, pairGroup)
			basisIds = append(basisIds, basisID)
		} else {
			for i, g := range groups {
				if g.ID == basisID {
					groups[i].Pairs = append(g.Pairs, partialPair)
					break
				}
			}
		}
	}

	data := map[string][]PairStatisticGroup{
		"pairs": groups,
	}

	return response.Success(data, "")
}

//func (s *service) getTrendDataForPairs(pairs []Pair, timeFrame string, from time.Time, to time.Time) map[int64][]KlineTrendData {
//	result := make(map[int64][]KlineTrendData, 0)
//	var pairIds []int64
//	for _, p := range pairs {
//		pairIds = append(pairIds, p.ID)
//		result[p.ID] = []KlineTrendData{}
//	}
//
//	trends := s.klineService.GetPairKlinesBetweenTimes(pairIds, timeFrame, from, to)
//	for _, t := range trends {
//		result[t.PairId] = append(result[t.PairId], t)
//	}
//
//	return result
//}

func isInSlice(data []int64, value int64) bool {
	for _, v := range data {
		if v == value {
			return true
		}
	}
	return false

}

func (s *service) AddOrRemoveFavoritePair(u *user.User, params FavoriteParams) (apiResponse response.APIResponse, statusCode int) {
	pair, err := s.GetPairByID(params.PairID)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("error getting pair by id", err,
			zap.String("service", "currencyService"),
			zap.String("method", "AddOrRemoveFavoritePair"),
			zap.Int64("pairID", params.PairID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || pair.ID == 0 {
		return response.Error("pair not found", http.StatusUnprocessableEntity, nil)
	}

	favoritePair := &FavoritePair{}
	err = s.favoritePairRepository.GetFavoritePair(u.ID, pair.ID, favoritePair)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("error getting favorite pair id", err,
			zap.String("service", "currencyService"),
			zap.String("method", "AddOrRemoveFavoritePair"),
			zap.Int64("pairID", pair.ID),
			zap.Int("userID", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if params.Action == "remove" {
		if errors.Is(err, gorm.ErrRecordNotFound) { //already does not exist
			return response.Success(nil, "")
		}

		err = s.favoritePairRepository.Delete(favoritePair)
		if err != nil {
			s.logger.Error2("error deleting favorite pair", err,
				zap.String("service", "currencyService"),
				zap.String("method", "AddOrRemoveFavoritePair"),
				zap.Int64("pairID", pair.ID),
				zap.Int("userID", u.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}

		return response.Success(nil, "")

	} else {
		if err == nil { //already added to user favorites
			return response.Success(nil, "")
		}

		newFavoritePair := &FavoritePair{
			UserID: u.ID,
			PairID: pair.ID,
		}

		err = s.favoritePairRepository.Create(newFavoritePair)
		if err != nil {
			s.logger.Error2("error creating favorite pair", err,
				zap.String("service", "currencyService"),
				zap.String("method", "AddOrRemoveFavoritePair"),
				zap.Int64("pairID", pair.ID),
				zap.Int("userID", u.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}

		return response.Success(nil, "")
	}
}

func (s *service) GetFavoritePairs(u *user.User) (apiResponse response.APIResponse, statusCode int) {
	result := make([]GetFavoritePairsResponse, 0)

	favoritePairs := s.favoritePairRepository.GetUserFavoritePairs(u.ID)

	for _, f := range favoritePairs {
		res := GetFavoritePairsResponse{
			ID:   f.PairID,
			Name: f.PairName,
		}
		result = append(result, res)
	}

	return response.Success(result, "")
}

func (s *service) GetPairRatio(params GetPairRatioParams) (apiResponse response.APIResponse, statusCode int) {

	pairName := strings.ToUpper(params.PairName)
	coins := strings.Split(pairName, "-")
	if len(coins) != 2 {
		return response.Error("pair is not valid", http.StatusUnprocessableEntity, nil)
	}

	firstCoinName := coins[0]
	secondCoinName := coins[1]

	result := make(map[string]float64)

	firstCoin := &Coin{}
	err := s.repository.GetCoinByCode(firstCoinName, firstCoin)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("error getting coin by code", err,
			zap.String("service", "currencyService"),
			zap.String("method", "GetPairRatio"),
			zap.String("firstCoinName", firstCoinName),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) || firstCoin.ID == 0 {
		return response.Error("coin not found", http.StatusUnprocessableEntity, nil)
	}

	secondCoin := &Coin{}
	err = s.repository.GetCoinByCode(secondCoinName, secondCoin)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("error getting coin by code", err,
			zap.String("service", "currencyService"),
			zap.String("method", "GetPairRatio"),
			zap.String("secondCoinName", secondCoinName),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) || secondCoin.ID == 0 {
		return response.Error("coin not found", http.StatusUnprocessableEntity, nil)
	}

	if secondCoinName == USDT {
		ratio, err := s.priceGenerator.GetAmountBasedOnUSDT(context.Background(), firstCoin.Code, "1.0")
		if err != nil {
			s.logger.Error2("can not get amount based on usdt", err,
				zap.String("service", "currencyService"),
				zap.String("method", "GetPairRatio"),
				zap.String("firstCoinCode", firstCoin.Code),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}

		ratioFloat, err := strconv.ParseFloat(ratio, 64)
		if err != nil {
			s.logger.Error2("can not parse float", err,
				zap.String("service", "currencyService"),
				zap.String("method", "GetPairRatio"),
				zap.String("ratio", ratio),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}

		result["ratio"] = ratioFloat

	} else if firstCoinName == USDT {
		ratio, err := s.priceGenerator.GetAmountBasedOnUSDT(context.Background(), secondCoin.Code, "1.0")
		if err != nil {
			s.logger.Error2("can not get amount based on usdt", err,
				zap.String("service", "currencyService"),
				zap.String("method", "GetPairRatio"),
				zap.String("ratio", ratio),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}

		ratioFloat, err := strconv.ParseFloat(ratio, 64)
		if err != nil {
			s.logger.Error2("can not parse float", err,
				zap.String("service", "currencyService"),
				zap.String("method", "GetPairRatio"),
				zap.String("ratio", ratio),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}

		if ratioFloat == float64(0) {
			return response.Error("ratio is not available", http.StatusUnprocessableEntity, nil)
		}

		result["ratio"] = 1 / ratioFloat

	} else {
		firstCoinRatio, err := s.priceGenerator.GetAmountBasedOnUSDT(context.Background(), firstCoin.Code, "1.0")
		if err != nil {
			s.logger.Error2("can not get anount based on usdt", err,
				zap.String("service", "currencyService"),
				zap.String("method", "GetPairRatio"),
				zap.String("firstCoinCode", firstCoin.Code),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}

		firstCoinRatioFloat, err := strconv.ParseFloat(firstCoinRatio, 64)
		if err != nil {
			s.logger.Error2("can not parse float", err,
				zap.String("service", "currencyService"),
				zap.String("method", "GetPairRatio"),
				zap.String("ratio", firstCoinRatio),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}

		secondCoinRatio, err := s.priceGenerator.GetAmountBasedOnUSDT(context.Background(), secondCoin.Code, "1.0")
		if err != nil {
			s.logger.Error2("can not get anount based on usdt", err,
				zap.String("service", "currencyService"),
				zap.String("method", "GetPairRatio"),
				zap.String("secondCoinCode", secondCoin.Code),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}

		secondCoinRatioFloat, err := strconv.ParseFloat(secondCoinRatio, 64)
		if err != nil {
			s.logger.Error2("can not parse float", err,
				zap.String("service", "currencyService"),
				zap.String("method", "GetPairRatio"),
				zap.String("ratio", secondCoinRatio),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}

		if secondCoinRatioFloat == float64(0) {
			return response.Error("ratio is not available", http.StatusUnprocessableEntity, nil)
		}

		ratioFloat := firstCoinRatioFloat / secondCoinRatioFloat
		result["ratio"] = ratioFloat
	}
	return response.Success(result, "")
}

func (s *service) GetActivePairCurrenciesList() []Pair {
	return s.pairRepository.GetActivePairCurrenciesList()
}

func GetPairByName(pairs []Pair, pairName string) Pair {
	for _, pair := range pairs {
		if pair.Name == pairName {
			return pair
		}
	}
	panic("we should never reach here")
}

func (s *service) GetPairByID(id int64) (Pair, error) {
	p := Pair{}
	err := s.pairRepository.GetPairByID(id, &p)
	return p, err
}

func (s *service) GetPairByName(name string) (Pair, error) {
	p := Pair{}
	err := s.pairRepository.GetPairByName(name, &p)
	return p, err
}

func (s *service) GetCoinByCode(code string) (Coin, error) {
	c := Coin{}
	err := s.repository.GetCoinByCode(code, &c)
	return c, err
}

func (s *service) GetActiveCoins() []Coin {
	return s.repository.GetActiveCoins()
}

func NewCurrencyService(repository Repository, liveData livedata.Service, priceGenerator PriceGenerator,
	pairRepository PairRepository, klineService KlineService, favoritePairRepository FavoritePairRepository,
	configs platform.Configs, logger platform.Logger) Service {
	return &service{
		repository:             repository,
		liveData:               liveData,
		priceGenerator:         priceGenerator,
		pairRepository:         pairRepository,
		klineService:           klineService,
		favoritePairRepository: favoritePairRepository,
		configs:                configs,
		logger:                 logger,
	}
}
