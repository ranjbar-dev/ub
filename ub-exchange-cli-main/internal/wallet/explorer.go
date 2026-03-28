package wallet

import (
	"fmt"
	"strings"
)

type explorerAddress struct {
	addressURL string
	txURL      string
}

var explorers = map[string]explorerAddress{
	"BTC": {
		addressURL: "https://www.blockchain.com/btc/address/%s",
		txURL:      "https://www.blockchain.com/btc/tx/%s",
	},
	"ETH": {
		addressURL: "https://etherscan.io/address/%s",
		txURL:      "https://etherscan.io/tx/%s",
	},
	"BCH": {
		addressURL: "https://www.blockchain.com/bch/address/%s",
		txURL:      "https://www.blockchain.com/bch/tx/%s",
	},
	"DASH": {
		addressURL: "https://explorer.dash.org/insight/address/%s",
		txURL:      "https://explorer.dash.org/insight/tx/%s",
	},
	"DGB": {
		addressURL: "https://digiexplorer.info/address/%s",
		txURL:      "https://digiexplorer.info/tx/%s",
	},
	"DOGE": {
		addressURL: "https://blockchair.com/dogecoin/address/%s",
		txURL:      "https://blockchair.com/dogecoin/transaction/%s",
	},
	"GRS": {
		addressURL: "https://groestlsight.groestlcoin.org/address/%s",
		txURL:      "https://groestlsight.groestlcoin.org/tx/%s",
	},
	"LTC": {
		addressURL: "https://explorer.zcha.in/accounts/%s",
		txURL:      "https://explorer.zcha.in/transactions/%s",
	},
	"TRX": {
		addressURL: "https://tronscan.org/#/address/%s",
		txURL:      "https://tronscan.org/#/transaction/%s",
	},

	//"USDT": {
	//	addressUrl: "https://etherscan.io/address/%s",
	//	txUrl:      "https://etherscan.io/tx/%s",
	//},
}

func GetTxExplorer(coin string, network string, txID string) string {
	explorerURL := ""

	if txID == "" {
		return explorerURL
	}

	if network != "" {
		explorerAddress, ok := explorers[strings.ToUpper(network)]
		if ok {
			return fmt.Sprintf(explorerAddress.txURL, txID)
		}

		return explorerURL

	} else {
		explorerAddress, ok := explorers[strings.ToUpper(coin)]
		if ok {
			return fmt.Sprintf(explorerAddress.txURL, txID)
		}

		return explorerURL
	}

}

func GetAddressExplorer(coin string, network string, address string) string {
	explorerURL := ""
	if address == "" {
		return explorerURL
	}

	if network != "" {
		explorerAddress, ok := explorers[strings.ToUpper(network)]
		if ok {
			return fmt.Sprintf(explorerAddress.addressURL, address)
		}

		return explorerURL

	} else {
		explorerAddress, ok := explorers[strings.ToUpper(coin)]
		if ok {
			return fmt.Sprintf(explorerAddress.addressURL, address)
		}

		return explorerURL
	}

}
