package order

import (
	"context"
	"exchange-go/internal/currency"
	"strconv"

	"github.com/shopspring/decimal"
)

// ForceTrader provides price-threshold logic to determine whether an order's price
// falls within the acceptable trading range defined by bot rules.
type ForceTrader interface {
	// ShouldForceTrade returns true if the given price is within the acceptable threshold
	// for the pair and order type, indicating the order should be force-traded.
	ShouldForceTrade(pairName string, orderType string, price string) (bool, error)
	// GetMinAndMaxPrice calculates and returns the minimum and maximum acceptable prices
	// for a pair and order type based on the current market price and bot rules.
	GetMinAndMaxPrice(pairName string, orderType string, givenPrice string) (string, string, error)
}

type forceTrader struct {
	priceGenerator  currency.PriceGenerator
	currencyService currency.Service
}

func (f *forceTrader) ShouldForceTrade(pairName string, orderType string, priceString string) (bool, error) {
	//market orders always should be traded
	if priceString == "" {
		return true, nil
	}
	//return false
	minPriceString, maxPriceString, err := f.GetMinAndMaxPrice(pairName, orderType, "")
	if err != nil {
		return false, err
	}

	minPriceDecimal, err := decimal.NewFromString(minPriceString)
	if err != nil {
		return false, err
	}

	maxPriceDecimal, err := decimal.NewFromString(maxPriceString)
	if err != nil {
		return false, err
	}

	priceDecimal, err := decimal.NewFromString(priceString)
	if err != nil {
		return false, err
	}

	if orderType == TypeSell {
		return priceDecimal.LessThanOrEqual(maxPriceDecimal), nil
	}

	return priceDecimal.GreaterThanOrEqual(minPriceDecimal), nil
}

func (f *forceTrader) GetMinAndMaxPrice(pairName string, orderType string, givenPrice string) (minPrice string, maxPrice string, err error) {
	ctx := context.Background()
	pairPriceString := givenPrice
	if pairPriceString == "" {
		var err error
		pairPriceString, err = f.priceGenerator.GetPrice(ctx, pairName)
		if err != nil {
			return minPrice, maxPrice, err
		}

	}

	pairPrice, err := strconv.ParseFloat(pairPriceString, 64)
	if err != nil {
		return minPrice, maxPrice, err
	}
	pair, err := f.currencyService.GetPairByName(pairName)
	if err != nil {
		return minPrice, maxPrice, err
	}

	botRules := pair.GetBotRules()

	if botRules.Type == currency.TraderBotRuleTypePercentage {
		minPriceFloat := (1 - botRules.BuyValue) * pairPrice
		maxPriceFloat := (1 + botRules.BuyValue) * pairPrice
		if orderType == TypeSell {
			minPriceFloat = (1 - botRules.SellValue) * pairPrice
			maxPriceFloat = (1 + botRules.SellValue) * pairPrice
		}
		minPriceString := strconv.FormatFloat(minPriceFloat, 'f', 8, 64)
		maxPriceString := strconv.FormatFloat(maxPriceFloat, 'f', 8, 64)
		return minPriceString, maxPriceString, nil
	} else {
		minPriceFloat := pairPrice - botRules.BuyValue
		maxPriceFloat := pairPrice + botRules.BuyValue
		if orderType == TypeSell {
			minPriceFloat = pairPrice - botRules.SellValue
			maxPriceFloat = pairPrice + botRules.SellValue
		}
		minPriceString := strconv.FormatFloat(minPriceFloat, 'f', 8, 64)
		maxPriceString := strconv.FormatFloat(maxPriceFloat, 'f', 8, 64)
		return minPriceString, maxPriceString, nil
	}

}

func NewForceTrader(priceGenerator currency.PriceGenerator, currencyService currency.Service) ForceTrader {
	return &forceTrader{
		priceGenerator:  priceGenerator,
		currencyService: currencyService,
	}

}
