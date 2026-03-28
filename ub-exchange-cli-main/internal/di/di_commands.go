package di

import (
	"exchange-go/internal/command"
	"exchange-go/internal/communication"
	"exchange-go/internal/currency"
	"exchange-go/internal/externalexchange"
	"exchange-go/internal/livedata"
	"exchange-go/internal/order"
	"exchange-go/internal/payment"
	"exchange-go/internal/platform"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"

	"github.com/sarulabs/di"
	"gorm.io/gorm"
)

// DI registrations for CLI console commands.
// Commands depend on domain services and repositories — they are the outermost layer.
// Each command implements the ConsoleCommand interface and is registered in main.go.
// Commands have no dependants; they are terminal nodes in the dependency graph.
func addDeleteCacheCommand() {
	mustAdd(di.Def{
		Name:  DeleteCacheCommand,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			cache := ctn.Get(cacheService).(platform.Cache)
			logger := ctn.Get(LoggerService).(platform.Logger)
			return command.NewDeleteCacheCmd(cache, logger), nil
		},
	})
}

func addUbCaptchaDecryptionCommand() {
	mustAdd(di.Def{
		Name:  UbCaptchaDecryptionCommand,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			ubCaptchaMgr := ctn.Get(ubCaptchaManager).(user.UbCaptchaManager)
			logger := ctn.Get(LoggerService).(platform.Logger)
			return command.NewUbCaptchaDecryptionCmd(ubCaptchaMgr, logger), nil
		},
	})
}

func addUbCaptchaEncryptionCommand() {
	mustAdd(di.Def{
		Name:  UbCaptchaEncryptionCommand,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			ubCaptchaMgr := ctn.Get(ubCaptchaManager).(user.UbCaptchaManager)
			logger := ctn.Get(LoggerService).(platform.Logger)
			return command.NewUbCaptchaEncryptionCmd(ubCaptchaMgr, logger), nil
		},
	})
}

func addUbCaptchaKeyGeneratorCommand() {
	mustAdd(di.Def{
		Name:  UbCaptchaKeyGeneratorCommand,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			ubCaptchaMgr := ctn.Get(ubCaptchaManager).(user.UbCaptchaManager)
			logger := ctn.Get(LoggerService).(platform.Logger)
			return command.NewUbCaptchaKeyGeneratorCmd(ubCaptchaMgr, logger), nil
		},
	})
}

func addUbUpdateUserWalletBalancesCommand() {
	mustAdd(di.Def{
		Name:  UpdateUserWalletBalancesCommand,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			userSvc := ctn.Get(userService).(user.Service)
			userBalanceSvc := ctn.Get(userBalanceService).(userbalance.Service)
			currencySvc := ctn.Get(currencyService).(currency.Service)
			logger := ctn.Get(LoggerService).(platform.Logger)
			return command.NewUbUpdateUserWalletBalances(userSvc, userBalanceSvc, currencySvc, logger), nil
		},
	})
}

func addCheckWithdrawalsInExternalExchangeCommand() {
	mustAdd(di.Def{
		Name:  CheckWithdrawalsInExternalExchangeCommand,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			paymentSvc := ctn.Get(paymentService).(payment.Service)
			externalExchangeSvc := ctn.Get(externalExchangeService).(externalexchange.Service)
			internalTransferSvc := ctn.Get(internalTransferService).(payment.InternalTransferService)
			configService := ctn.Get(ConfigService).(platform.Configs)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := command.NewCheckWithdrawalsCmd(
				paymentSvc,
				externalExchangeSvc,
				internalTransferSvc,
				configService,
				logger,
			)
			return srv, nil
		},
	})

}

func addUpdateOrdersInExternalExchangeCommand() {
	mustAdd(di.Def{
		Name:  UpdateOrdersInExternalExchangeCommand,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			currencySvc := ctn.Get(currencyService).(currency.Service)
			orderFromExternalSvc := ctn.Get(orderFromExternalService).(externalexchange.OrderFromExternalService)
			externalExchangeSvc := ctn.Get(externalExchangeService).(externalexchange.Service)
			configService := ctn.Get(ConfigService).(platform.Configs)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := command.NewUpdateOrdersInExternalExchangeCmd(
				currencySvc,
				orderFromExternalSvc,
				externalExchangeSvc,
				configService,
				logger,
			)
			return srv, nil
		},
	})
}

func addSyncKlineCommand() {
	mustAdd(di.Def{
		Name:  KlineSyncCommand,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			currencySvc := ctn.Get(currencyService).(currency.Service)
			klineSvc := ctn.Get(klineService).(currency.KlineService)
			externalExchangeSvc := ctn.Get(externalExchangeService).(externalexchange.Service)
			queueMgr := ctn.Get(queueManager).(communication.QueueManager)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := command.NewSyncKlineCmd(
				currencySvc,
				klineSvc,
				externalExchangeSvc,
				queueMgr,
				logger,
			)
			return srv, nil
		},
	})
}

func addGenerateKlineSyncCommand() {
	mustAdd(di.Def{
		Name:  GenerateKlineSyncCommand,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			currencySvc := ctn.Get(currencyService).(currency.Service)
			klineSvc := ctn.Get(klineService).(currency.KlineService)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := command.NewGenerateKlineSyncCmd(
				currencySvc,
				klineSvc,
				logger,
			)
			return srv, nil
		},
	})
}

func addRetrieveExternalOrdersToRedisCommand() {
	mustAdd(di.Def{
		Name:  RetrieveExternalOrdersToRedisCommand,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			externalExchangeOrderSvc := ctn.Get(externalExchangeOrderService).(externalexchange.OrderService)
			botAggregationSvc := ctn.Get(botAggregationService).(order.BotAggregationService)
			tradeSvc := ctn.Get(tradeService).(order.TradeService)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := command.NewRetrieveExternalOrdersToRedisCmd(
				externalExchangeOrderSvc,
				botAggregationSvc,
				tradeSvc,
				logger,
			)
			return srv, nil
		},
	})
}

func addSubmitBotAggregatedOrderCommand() {
	mustAdd(di.Def{
		Name:  SubmitBotAggregatedOrderCommand,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			currencySvc := ctn.Get(currencyService).(currency.Service)
			botAggregationSvc := ctn.Get(botAggregationService).(order.BotAggregationService)
			liveDataSvc := ctn.Get(liveDataService).(livedata.Service)
			externalExchangeOrderSvc := ctn.Get(externalExchangeOrderService).(externalexchange.OrderService)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := command.NewSubmitBotAggregatedCmd(
				currencySvc,
				botAggregationSvc,
				liveDataSvc,
				externalExchangeOrderSvc,
				logger,
			)
			return srv, nil
		},
	})
}

func addRetrieveOpenOrdersToRedisCommand() {
	mustAdd(di.Def{
		Name:  RetrieveOpenOrdersToRedisCommand,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			orderRepo := ctn.Get(orderRepository).(order.Repository)
			orderRedisMgr := ctn.Get(orderRedisManager).(order.RedisManager)
			engineCommInst := ctn.Get(engineCommunicator).(order.EngineCommunicator)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := command.NewRetrieveOrderToRedisCmd(
				orderRepo,
				orderRedisMgr,
				engineCommInst,
				logger,
			)
			return srv, nil
		},
	})
}

func addGenerateAddressCommand() {
	mustAdd(di.Def{
		Name:  GenerateAddressCommand,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			userRepo := ctn.Get(userRepository).(user.Repository)
			userBalanceRepo := ctn.Get(userBalanceRepository).(userbalance.Repository)
			userBalanceSvc := ctn.Get(userBalanceService).(userbalance.Service)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := command.NewGenerateAddressCmd(
				userRepo,
				userBalanceRepo,
				userBalanceSvc,
				logger,
			)
			return srv, nil
		},
	})
}

func addInitializeBalanceCommand() {
	mustAdd(di.Def{
		Name:  InitializeBalanceCommand,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			userRepo := ctn.Get(userRepository).(user.Repository)
			currencySvc := ctn.Get(currencyService).(currency.Service)
			userBalanceRepo := ctn.Get(userBalanceRepository).(userbalance.Repository)
			userBalanceSvc := ctn.Get(userBalanceService).(userbalance.Service)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := command.NewInitializedBalanceCmd(
				userRepo,
				currencySvc,
				userBalanceRepo,
				userBalanceSvc,
				logger,
			)
			return srv, nil
		},
	})
}

func addWsHealthCheckCommand() {
	mustAdd(di.Def{
		Name:  WsHealthCheckCommand,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			liveDataSvc := ctn.Get(liveDataService).(livedata.Service)
			logger := ctn.Get(LoggerService).(platform.Logger)
			return command.NewWsHealthCheckCmd(liveDataSvc, logger), nil
		},
	})
}

func addSetUserLevelCommand() {
	mustAdd(di.Def{
		Name:  SetUserLevelCommand,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			userSvc := ctn.Get(userService).(user.Service)
			tradeSvc := ctn.Get(tradeService).(order.TradeService)
			pairRepo := ctn.Get(pairRepository).(currency.PairRepository)
			currencyRepo := ctn.Get(currencyRepository).(currency.Repository)
			klineSvc := ctn.Get(klineService).(currency.KlineService)
			userLevelSvc := ctn.Get(userLevelService).(user.LevelService)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := command.NewSetUserLevelCmd(
				dbCli,
				userSvc,
				tradeSvc,
				pairRepo,
				currencyRepo,
				klineSvc,
				userLevelSvc,
				logger,
			)
			return srv, nil
		},
	})
}
