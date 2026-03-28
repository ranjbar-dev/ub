package order

import (
	"exchange-go/internal/currency"

	"github.com/shopspring/decimal"
)

func (ps *postOrderMatchingService) shouldCreateTrade(trades []tempTrade, orderID int64, tradedWithOrderID int64) bool {
	for _, trade := range trades {
		if trade.sellOrderID == orderID || trade.sellOrderID == tradedWithOrderID ||
			trade.buyOrderID == orderID || trade.buyOrderID == tradedWithOrderID {
			return false
		}
	}
	return true
}

func (ps *postOrderMatchingService) GetFeePercentage(data MatchingNeededQueryFields, pair currency.Pair, isMarket bool, isMaker bool) float64 {
	//market orders always considered as taker
	if isMarket || !isMaker {
		return float64(data.TakerFeePercentage) * pair.TakerFee
	}

	return float64(data.MakerFeePercentage) * pair.MakerFee

}

func (ps *postOrderMatchingService) getFinalAmountsDecimal(tempOrder tempOrder) (demandedDecimal decimal.Decimal, payedByDecimal decimal.Decimal) {
	tradeAmountDecimal, _ := decimal.NewFromString(tempOrder.tradeAmount)
	tradePriceDecimal, _ := decimal.NewFromString(tempOrder.tradePrice)
	multiplication := tradeAmountDecimal.Mul(tradePriceDecimal)
	if tempOrder.orderType == TypeBuy {
		demandedDecimal = tradeAmountDecimal
		payedByDecimal = multiplication
	} else {
		demandedDecimal = multiplication
		payedByDecimal = tradeAmountDecimal
	}

	return demandedDecimal, payedByDecimal

}
