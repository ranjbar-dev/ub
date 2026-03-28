package currency

import "strings"

const (
	TradeStatusFullTrade = "FULL_TRADE"
	TradeStatusSellOnly  = "SELL_ONLY"
	TradeStatusBuyOnly   = "BUY_ONLY"
	TradeStatusCloseOnly = "CLOSE_ONLY"
	ActionCancel         = "CANCEL"
	ActionSell           = "SELL"
	ActionBuy            = "BUY"
)

func IsActionAllowed(pair Pair, action string) bool {
	tradeStatus := mapActionToTradeStatus(action)
	pairTradeStatus := pair.TradeStatus

	if pairTradeStatus == TradeStatusFullTrade {
		return true
	}

	if tradeStatus == TradeStatusCloseOnly { //in all situation user can close his order
		return true
	}

	if pairTradeStatus == TradeStatusBuyOnly && (tradeStatus == TradeStatusBuyOnly || tradeStatus == TradeStatusCloseOnly) {
		return true
	}

	if pairTradeStatus == TradeStatusSellOnly && (tradeStatus == TradeStatusSellOnly || tradeStatus == TradeStatusCloseOnly) {
		return true
	}

	if pairTradeStatus == TradeStatusCloseOnly && tradeStatus == TradeStatusCloseOnly {
		return true
	}

	return false

}

func mapActionToTradeStatus(action string) string {

	switch strings.ToUpper(action) {
	case ActionBuy:
		return TradeStatusBuyOnly

	case ActionSell:
		return TradeStatusSellOnly

	default: //for cancel action
		return TradeStatusCloseOnly
	}

}
