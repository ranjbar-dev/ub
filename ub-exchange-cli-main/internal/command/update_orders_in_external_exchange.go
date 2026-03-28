package command

import (
	"context"
	"database/sql"
	"errors"
	"exchange-go/internal/currency"
	"exchange-go/internal/externalexchange"
	"exchange-go/internal/platform"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type updateOrdersInExternalExchangeCmd struct {
	currencyService          currency.Service
	orderFromExternalService externalexchange.OrderFromExternalService
	externalExchangeService  externalexchange.Service
	configs                  platform.Configs
	logger                   platform.Logger
}

func (cmd *updateOrdersInExternalExchangeCmd) Run(ctx context.Context, flags []string) {
	fmt.Println("start of update orders in external exchange command")
	pairs := cmd.currencyService.GetActivePairCurrenciesList()
	for _, pair := range pairs {
		//fetching orders
		lastOrderFromExternal := &externalexchange.OrderFromExternal{}
		err := cmd.orderFromExternalService.GetLastOrderFromExternalByPairID(pair.ID, lastOrderFromExternal)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			cmd.logger.Error2("can not get last order from external", err,
				zap.String("service", "updateOrdersInExternalExchangeCmd"),
				zap.String("method", "Run"),
				zap.String("pairName", pair.Name),
			)
			continue
		}
		since := int64(0)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			since = lastOrderFromExternal.Timestamp.Int64 + 1
		}

		if cmd.configs.GetEnv() == platform.EnvProd {
			//todo check rate-limit instead of sleep
			time.Sleep(5 * time.Second) //just not to be banned by binance
		}

		fetchOrdersParams := externalexchange.FetchOrdersParams{
			Pair: pair.Name,
			From: since,
		}
		externalOrders, err := cmd.externalExchangeService.FetchOrders(fetchOrdersParams)
		if err != nil {
			cmd.logger.Error2("can not fetch orders", err,
				zap.String("service", "updateOrdersInExternalExchangeCmd"),
				zap.String("method", "Run"),
				zap.String("pairName", pair.Name),
			)
			continue
		}
		for _, o := range externalOrders {
			orderFromExternal := &externalexchange.OrderFromExternal{
				PairID:          sql.NullInt64{Int64: pair.ID, Valid: true},
				ExternalOrderID: o.OrderID,
				ClientOrderID:   o.ClientOrderID,
				Type:            o.Type,
				ExchangeType:    o.ExchangeType,
				Price:           sql.NullString{String: o.Price, Valid: true},
				Amount:          sql.NullString{String: o.Amount, Valid: true},
				Status:          sql.NullString{String: o.Status, Valid: true},
				MetaData:        sql.NullString{String: o.Data, Valid: true},
				Time:            sql.NullTime{Time: o.DateTime, Valid: true},
				Timestamp:       sql.NullInt64{Int64: o.Timestamp, Valid: true},
			}

			err := cmd.orderFromExternalService.CreateOrder(orderFromExternal)
			if err != nil {
				cmd.logger.Error2("can not create order", err,
					zap.String("service", "updateOrdersInExternalExchangeCmd"),
					zap.String("method", "Run"),
					zap.String("pairName", pair.Name),
				)
				continue
			}
		}

		////fetching trades
		////Updating trades disabled
		//lastTradeFromExternal := &externalexchange.TradeFromExternal{}
		//err = cmd.orderFromExternalService.GetLastTradeFromExternalByPairId(pair.ID, lastTradeFromExternal)
		//if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		//	cmd.logger.Error("error in updateOrdersInExternalExchangeCmd", err)
		//	continue
		//}
		//since = int64(0)
		//if !errors.Is(err, gorm.ErrRecordNotFound) {
		//	since = lastTradeFromExternal.Timestamp.Int64 + 1
		//}
		//
		//if cmd.configs.GetEnv() == platform.EnvProd {
		//	//todo check rate-limit instead of sleep
		//	time.Sleep(3 * time.Second) //just not to be banned by binance
		//}
		//
		//fetchTradesParams := externalexchange.FetchTradesParams{
		//	Pair: pair.Name,
		//	From: since,
		//}
		//externalTrades, err := cmd.externalExchangeService.FetchTrades(fetchTradesParams)
		//if err != nil {
		//	cmd.logger.Error("error in updateOrdersInExternalExchangeCmd", err)
		//	continue
		//}
		//
		//for _, t := range externalTrades {
		//
		//	tradeFromExternal := &externalexchange.TradeFromExternal{
		//		OrderId:         sql.NullInt64{}, //todo handle this later
		//		ExternalTradeId: t.Id,
		//		Price:           sql.NullString{String: t.Price, Valid: true},
		//		Amount:          sql.NullString{String: t.Amount, Valid: true},
		//		Commission:      sql.NullString{String: t.Commission, Valid: true},
		//		Coin:            sql.NullString{String: t.Coin, Valid: true},
		//		MetaData:        sql.NullString{String: t.Data, Valid: true},
		//		Timestamp:       sql.NullInt64{Int64: t.Timestamp, Valid: true},
		//	}
		//
		//	err := cmd.orderFromExternalService.CreateTrade(tradeFromExternal)
		//	if err != nil {
		//		cmd.logger.Error("error in updateOrdersInExternalExchangeCmd", err)
		//		continue
		//	}
		//
		//}

	}

	fmt.Println("end of update orders in external exchange command")
}

func NewUpdateOrdersInExternalExchangeCmd(currencyService currency.Service,
	orderFromExternalService externalexchange.OrderFromExternalService, externalExchangeService externalexchange.Service,
	configs platform.Configs, logger platform.Logger) ConsoleCommand {
	return &updateOrdersInExternalExchangeCmd{
		currencyService:          currencyService,
		orderFromExternalService: orderFromExternalService,
		externalExchangeService:  externalExchangeService,
		configs:                  configs,
		logger:                   logger,
	}
}
