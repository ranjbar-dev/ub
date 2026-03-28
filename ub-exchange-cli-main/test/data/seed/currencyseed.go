package seed

import (
	"database/sql"
	"exchange-go/internal/currency"

	"gorm.io/gorm"
)

func CurrencySeed(db *gorm.DB) {
	otherNetworksConfigs := `[{"code":"TRX","supportsWithdraw":true,"supportsDeposit":true,"completedNetworkName":"Tron (TRX)","fee":"2.5"}]`
	usdt := currency.Coin{
		ID:   1,
		Name: "Tether",
		Code: "USDT",
		CompletedNetworkName: sql.NullString{
			String: "Ethereum(ETH) ERC20",
			Valid:  true,
		},
		SubUnit:                        6,
		ShowSubUnit:                    6,
		IsActive:                       true,
		IsMain:                         false,
		Priority:                       1,
		ConversionRatio:                0.1,
		Image:                          "/images",
		MinimumWithdraw:                "10",
		MaximumWithdraw:                "20000",
		BlockchainNetwork:              sql.NullString{String: "ETH", Valid: true},
		SupportsWithdraw:               sql.NullBool{Bool: true, Valid: true},
		SupportsDeposit:                sql.NullBool{Bool: true, Valid: true},
		WithdrawalFee:                  sql.NullFloat64{Float64: 2, Valid: true},
		OtherBlockchainNetworksConfigs: sql.NullString{String: otherNetworksConfigs, Valid: true},
	}

	btc := currency.Coin{
		ID:               2,
		Name:             "Bitcoin",
		Code:             "BTC",
		SubUnit:          8,
		ShowSubUnit:      8,
		IsActive:         true,
		IsMain:           true,
		Priority:         2,
		ConversionRatio:  1,
		Image:            "/images",
		MinimumWithdraw:  "0.001",
		MaximumWithdraw:  "5.0",
		SupportsWithdraw: sql.NullBool{Bool: true, Valid: true},
		SupportsDeposit:  sql.NullBool{Bool: true, Valid: true},
		WithdrawalFee:    sql.NullFloat64{Float64: 0.0001, Valid: true},
	}

	eth := currency.Coin{
		ID:               3,
		Name:             "Ethereum",
		Code:             "ETH",
		SubUnit:          8,
		ShowSubUnit:      8,
		IsActive:         true,
		IsMain:           false,
		Priority:         3,
		ConversionRatio:  1,
		Image:            "/images",
		MinimumWithdraw:  "0.001",
		MaximumWithdraw:  "50.0",
		SupportsWithdraw: sql.NullBool{Bool: true, Valid: true},
		SupportsDeposit:  sql.NullBool{Bool: true, Valid: true},
		WithdrawalFee:    sql.NullFloat64{Float64: 0.0001, Valid: true},
	}

	grs := currency.Coin{
		ID:               4,
		Name:             "Groestlcoin",
		Code:             "GRS",
		SubUnit:          8,
		ShowSubUnit:      8,
		IsActive:         true,
		IsMain:           false,
		Priority:         4,
		ConversionRatio:  1,
		Image:            "/images",
		MinimumWithdraw:  "10.0",
		MaximumWithdraw:  "5000.0",
		SupportsWithdraw: sql.NullBool{Bool: true, Valid: true},
		SupportsDeposit:  sql.NullBool{Bool: true, Valid: true},
		WithdrawalFee:    sql.NullFloat64{Float64: 0.01, Valid: true},
	}

	dai := currency.Coin{
		ID:               5,
		Name:             "Dai",
		Code:             "DAI",
		SubUnit:          6,
		ShowSubUnit:      6,
		IsActive:         true,
		IsMain:           false,
		Priority:         5,
		ConversionRatio:  1,
		Image:            "/images",
		MinimumWithdraw:  "10",
		MaximumWithdraw:  "20000",
		SupportsWithdraw: sql.NullBool{Bool: true, Valid: true},
		SupportsDeposit:  sql.NullBool{Bool: true, Valid: true},
		WithdrawalFee:    sql.NullFloat64{Float64: 2, Valid: true},
	}

	trx := currency.Coin{
		ID:               6,
		Name:             "Tron",
		Code:             "TRX",
		SubUnit:          8,
		ShowSubUnit:      8,
		IsActive:         true,
		IsMain:           false,
		Priority:         5,
		ConversionRatio:  1,
		Image:            "/images",
		MinimumWithdraw:  "100",
		MaximumWithdraw:  "50000",
		SupportsWithdraw: sql.NullBool{Bool: true, Valid: true},
		SupportsDeposit:  sql.NullBool{Bool: true, Valid: true},
		WithdrawalFee:    sql.NullFloat64{Float64: 2, Valid: true},
	}

	db.Create(&usdt)
	db.Create(&btc)
	db.Create(&eth)
	db.Create(&grs)
	db.Create(&dai)
	db.Create(&trx)

	btcUsdt := currency.Pair{
		ID:                  1,
		Name:                "BTC-USDT",
		IsActive:            true,
		IsMain:              true,
		Spread:              1,
		ShowDigits:          6,
		BasisCoinID:         usdt.ID,
		BasisCoin:           usdt,
		DependentCoinID:     btc.ID,
		DependentCoin:       btc,
		MakerFee:            0.2,
		TakerFee:            0.3,
		TradeStatus:         "FULL_TRADE",
		AggregationStatus:   "RUN",
		MinimumOrderAmount:  sql.NullString{String: "10", Valid: true},
		MaxOurExchangeLimit: "3.0",
		BotRules:            sql.NullString{String: `{"buyValue":0.04,"sellValue":0.04,"type":"PERCENTAGE"}`, Valid: true},
	}

	ethUsdt := currency.Pair{
		ID:                  2,
		Name:                "ETH-USDT",
		IsActive:            true,
		IsMain:              true,
		Spread:              1,
		ShowDigits:          6,
		BasisCoinID:         usdt.ID,
		BasisCoin:           usdt,
		DependentCoinID:     eth.ID,
		DependentCoin:       eth,
		MakerFee:            0.2,
		TakerFee:            0.3,
		TradeStatus:         "FULL_TRADE",
		AggregationStatus:   "RUN",
		MaxOurExchangeLimit: "3.0",
	}

	ethBtc := currency.Pair{
		ID:                  3,
		Name:                "ETH-BTC",
		IsActive:            true,
		IsMain:              true,
		Spread:              1,
		ShowDigits:          6,
		BasisCoinID:         btc.ID,
		BasisCoin:           btc,
		DependentCoinID:     eth.ID,
		DependentCoin:       eth,
		MakerFee:            0.2,
		TakerFee:            0.3,
		TradeStatus:         "FULL_TRADE",
		AggregationStatus:   "RUN",
		MaxOurExchangeLimit: "3.0",
	}

	grsBtc := currency.Pair{
		ID:                  4,
		Name:                "GRS-BTC",
		IsActive:            true,
		IsMain:              true,
		Spread:              1,
		ShowDigits:          6,
		BasisCoinID:         btc.ID,
		BasisCoin:           btc,
		DependentCoinID:     grs.ID,
		DependentCoin:       grs,
		MakerFee:            0.2,
		TakerFee:            0.3,
		TradeStatus:         "FULL_TRADE",
		AggregationStatus:   "RUN",
		MaxOurExchangeLimit: "3.0",
	}

	usdtDai := currency.Pair{
		ID:                  5,
		Name:                "USDT-DAI",
		IsActive:            true,
		IsMain:              true,
		Spread:              1,
		ShowDigits:          6,
		BasisCoinID:         dai.ID,
		BasisCoin:           dai,
		DependentCoinID:     usdt.ID,
		DependentCoin:       usdt,
		MakerFee:            0.2,
		TakerFee:            0.3,
		TradeStatus:         "FULL_TRADE",
		AggregationStatus:   "RUN",
		MaxOurExchangeLimit: "10000.0",
	}

	btcDai := currency.Pair{
		ID:                  6,
		Name:                "BTC-DAI",
		IsActive:            true,
		IsMain:              true,
		Spread:              1,
		ShowDigits:          6,
		BasisCoinID:         dai.ID,
		BasisCoin:           dai,
		DependentCoinID:     btc.ID,
		DependentCoin:       btc,
		MakerFee:            0.2,
		TakerFee:            0.3,
		TradeStatus:         "FULL_TRADE",
		AggregationStatus:   "RUN",
		MaxOurExchangeLimit: "3.0",
	}

	db.Create(&btcUsdt)
	db.Create(&ethUsdt)
	db.Create(&ethBtc)
	db.Create(&grsBtc)
	db.Create(&usdtDai)
	db.Create(&btcDai)

}
