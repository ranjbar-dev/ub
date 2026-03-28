package currency

import (
	"context"
	"exchange-go/internal/livedata"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
)

// PriceGenerator converts cryptocurrency prices across different base currencies
// (USDT, BTC) using live market data with fallback to historical kline data.
type PriceGenerator interface {
	// GetBTCUSDTPrice returns the current BTC-USDT exchange rate.
	GetBTCUSDTPrice(ctx context.Context) (string, error)
	// GetAmountBasedOnUSDT converts the given coin amount to its USDT equivalent value.
	GetAmountBasedOnUSDT(ctx context.Context, coin string, amount string) (string, error)
	// GetAmountBasedOnBTC converts the given coin amount to its BTC equivalent value.
	GetAmountBasedOnBTC(ctx context.Context, coin string, amount string) (string, error)
	// GetPairPriceBasedOnUSDT returns a pair's current price normalized to USDT.
	GetPairPriceBasedOnUSDT(ctx context.Context, pairName string) (string, error)
	// GetPrice returns the current price for the specified trading pair.
	GetPrice(ctx context.Context, pairName string) (string, error)
}

type priceGenerator struct {
	liveDataService livedata.Service
	klineService    KlineService
	pairRepository  PairRepository
	activePairs     []Pair
}

func (p *priceGenerator) GetBTCUSDTPrice(ctx context.Context) (string, error) {
	return p.getPrice(ctx, PairBTCUSDT)
}

func (p *priceGenerator) GetAmountBasedOnUSDT(ctx context.Context, coin string, amount string) (string, error) {
	if coin == USDT {
		return amount, nil
	}

	for _, pair := range p.activePairs {
		if pair.BasisCoin.Code == USDT && pair.DependentCoin.Code == coin {
			price, err := p.getPrice(ctx, pair.Name)
			priceFloat, err := strconv.ParseFloat(price, 64)
			if err != nil {
				return "", err
			}
			amountFloat, err := strconv.ParseFloat(amount, 64)
			if err != nil {
				return "", err
			}
			finalAmount := strconv.FormatFloat(priceFloat*amountFloat, 'f', 8, 64)
			return finalAmount, nil
		}
		//for coins like DAI which we do not have a pair based on USDT but reverse exists USDT-DAI
		if pair.DependentCoin.Code == USDT && pair.BasisCoin.Code == coin {
			price, err := p.getPrice(ctx, pair.Name)
			if err != nil {
				return "", nil
			}
			priceFloat, err := strconv.ParseFloat(price, 64)
			if err != nil {
				return "", err
			}
			amountFloat64, _ := strconv.ParseFloat(amount, 64)
			finalAmount := strconv.FormatFloat(amountFloat64/priceFloat, 'f', 8, 64)
			return finalAmount, nil
		}
	}

	//for some coins like GRS we do not have a pair  like GRS-USDT or USDT-GRS
	// for this coins we try to find a mediator pair which is based on BTC
	priceBasedOnBtc, err := p.GetAmountBasedOnBTC(ctx, coin, "1")
	if err != nil {
		return "", err
	}

	priceBasedOnBtcFloat, err := strconv.ParseFloat(priceBasedOnBtc, 64)
	if err != nil {
		return "", err
	}

	BTCUSDTPrice, err := p.GetBTCUSDTPrice(ctx)
	if err != nil {
		return "", err
	}

	BTCUSDTPriceFloat, err := strconv.ParseFloat(BTCUSDTPrice, 64)
	if err != nil {
		return "", err
	}

	amountFloat64, _ := strconv.ParseFloat(amount, 64)

	finalAmount := strconv.FormatFloat(priceBasedOnBtcFloat*BTCUSDTPriceFloat*amountFloat64, 'f', 8, 64)
	return finalAmount, nil

}

func (p *priceGenerator) GetAmountBasedOnBTC(ctx context.Context, coin string, amount string) (string, error) {
	if coin == BTC {
		return amount, nil
	}

	for _, pair := range p.activePairs {
		if pair.BasisCoin.Code == BTC && pair.DependentCoin.Code == coin {
			price, err := p.getPrice(ctx, pair.Name)
			priceFloat, err := strconv.ParseFloat(price, 64)
			if err != nil {
				return "", err
			}
			amountFloat, err := strconv.ParseFloat(amount, 64)
			if err != nil {
				return "", err
			}
			finalAmount := strconv.FormatFloat(priceFloat*amountFloat, 'f', 8, 64)
			return finalAmount, nil
		}
		// for coin like USDT or DAI we try to find reverse pair BTC-USDT or BTC-DAI
		if pair.DependentCoin.Code == BTC && pair.BasisCoin.Code == coin {
			price, err := p.getPrice(ctx, pair.Name)
			if err != nil {
				return "", nil
			}
			priceFloat, err := strconv.ParseFloat(price, 64)
			if err != nil {
				return "", err
			}

			amountFloat64, _ := strconv.ParseFloat(amount, 64)
			finalAmount := strconv.FormatFloat(amountFloat64/priceFloat, 'f', 8, 64)
			return finalAmount, nil
		}
	}

	//we should never reach here
	return "", fmt.Errorf("no pair found based on btc or reverse")

}

func (p *priceGenerator) getPrice(ctx context.Context, pairName string) (string, error) {
	price, err := p.liveDataService.GetPrice(ctx, pairName)
	if err != nil && err != redis.Nil {
		return "", err
	}

	if err == redis.Nil {
		//no price in redis we get from kline service
		price, err = p.getPriceFromKline(ctx, pairName)
		if err != nil {
			return "", err
		}
	}

	priceDecimal, _ := decimal.NewFromString(price)
	return priceDecimal.StringFixed(8), nil
}

func (p *priceGenerator) getPriceFromKline(ctx context.Context, pairName string) (string, error) {
	price, err := p.klineService.GetLastPriceForPair(ctx, pairName, time.Now())
	return price, err
}

func (p *priceGenerator) GetPairPriceBasedOnUSDT(ctx context.Context, pairName string) (string, error) {
	for _, pair := range p.activePairs {
		if pair.Name == pairName {
			pairPrice, err := p.getPrice(ctx, pairName)
			if err != nil {
				return "", err
			}

			if err != nil {
				return "", err
			}

			if pair.BasisCoin.Code == USDT {
				return pairPrice, nil
			}

			pairPriceFloat, err := strconv.ParseFloat(pairPrice, 64)

			if pair.BasisCoin.Code == BTC {
				BTCUSDTPrice, err := p.GetBTCUSDTPrice(ctx)
				if err != nil {
					return "", err
				}
				BTCUSDTPriceFloat, err := strconv.ParseFloat(BTCUSDTPrice, 64)
				if err != nil {
					return "", err
				}

				finalPrice := pairPriceFloat * BTCUSDTPriceFloat
				return strconv.FormatFloat(finalPrice, 'f', 8, 64), nil
			}

			priceBasedOnUSDT, err := p.GetAmountBasedOnUSDT(ctx, pair.BasisCoin.Code, "1")
			if err != nil {
				return "", err
			}

			priceBasedOnUSDTFloat, err := strconv.ParseFloat(priceBasedOnUSDT, 64)
			if err != nil {
				return "", err
			}

			finalPrice := pairPriceFloat * priceBasedOnUSDTFloat
			return strconv.FormatFloat(finalPrice, 'f', 8, 64), nil

		}
	}

	return "", fmt.Errorf("no pair found based on btc or reverse")
}

func (p *priceGenerator) GetPrice(ctx context.Context, pairName string) (string, error) {
	return p.getPrice(ctx, pairName)
}

func NewPriceGenerator(liveDataService livedata.Service, klineService KlineService, pairRepository PairRepository) PriceGenerator {
	activePairs := pairRepository.GetActivePairCurrenciesList()
	return &priceGenerator{
		liveDataService: liveDataService,
		klineService:    klineService,
		pairRepository:  pairRepository,
		activePairs:     activePairs,
	}
}
