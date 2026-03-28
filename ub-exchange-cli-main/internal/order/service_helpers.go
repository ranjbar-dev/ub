package order

import "strings"

func mapSideToType(side string) string {
	if side == SideAsk {
		return TypeSell
	}
	return TypeBuy
}

func mapTypeToSide(orderType string) string {
	if orderType == TypeBuy {
		return SideBid
	}
	return SideAsk
}

func removeUnderline(text string) string {
	return strings.Replace(text, "_", " ", -1)
}
