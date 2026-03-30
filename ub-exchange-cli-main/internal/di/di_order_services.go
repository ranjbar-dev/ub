package di

import (
	"exchange-go/internal/communication"
	"exchange-go/internal/currency"
	"exchange-go/internal/engine"
	"exchange-go/internal/externalexchange"
	"exchange-go/internal/externalexchangews"
	"exchange-go/internal/livedata"
	"exchange-go/internal/order"
	"exchange-go/internal/orderbook"
	"exchange-go/internal/platform"
	"exchange-go/internal/processor"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"

	"github.com/sarulabs/di"
	"gorm.io/gorm"
)

// DI registrations for order and trading domain services.
// These form the core order-matching pipeline:
//
//	decisionManager → orderEventsHandler → engineCommunicator → engine
//	forceTrader     → postOrderMatchingService, engineCommunicator
//	postOrderMatchingService → engineResultHandler → engine
//	engine          → engineCommunicator → orderEventsHandler
//	stopOrderSubmissionManager → wsDataProcessor
//	inQueueOrderManager        → wsDataProcessor
//	orderCreateManager → orderService, autoExchangeManager
//	adminOrderManager  → orderService
//
// Also contains externalExchange services, botAggregation, wsDataProcessor,
// tradeService, and the unmatchedOrderHandler.
func addDecisionManager() {
	mustAdd(di.Def{
		Name:  decisionManager,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			configs := ctn.Get(ConfigService).(platform.Configs)
			return order.NewDecisionManager(configs), nil
		},
	})
}

func addForceTrader() {
	mustAdd(di.Def{
		Name:  forceTrader,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			priceGen := ctn.Get(priceGenerator).(currency.PriceGenerator)
			currencySvc := ctn.Get(currencyService).(currency.Service)
			return order.NewForceTrader(priceGen, currencySvc), nil
		},
	})
}

func addPostOrderMatchingService() {
	mustAdd(di.Def{
		Name:  postOrderMatchingService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			orderRepo := ctn.Get(orderRepository).(order.Repository)
			userBalanceSvc := ctn.Get(userBalanceService).(userbalance.Service)
			forceTradrInst := ctn.Get(forceTrader).(order.ForceTrader)
			priceGen := ctn.Get(priceGenerator).(currency.PriceGenerator)
			tradeEventsHdl := ctn.Get(tradeEventsHandler).(order.TradeEventsHandler)
			centrifugoMgr := ctn.Get(centrifugoManager).(communication.CentrifugoManager)
			redisClient := ctn.Get(RedisClient).(platform.RedisClient)
			currencySvc := ctn.Get(currencyService).(currency.Service)
			userSvc := ctn.Get(userService).(user.Service)
			userLevelSvc := ctn.Get(userLevelService).(user.LevelService)
			configService := ctn.Get(ConfigService).(platform.Configs)
			logger := ctn.Get(LoggerService).(platform.Logger)

			srv := order.NewPostOrderMatchingService(
				dbCli,
				orderRepo,
				userBalanceSvc,
				forceTradrInst,
				priceGen,
				tradeEventsHdl,
				centrifugoMgr,
				redisClient,
				currencySvc,
				userSvc,
				userLevelSvc,
				configService,
				logger,
			)
			return srv, nil
		},
	})
}

func addEngineResultHandler() {
	mustAdd(di.Def{
		Name:  EngineResultHandler,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			postOrderMatchingSvc := ctn.Get(postOrderMatchingService).(order.PostOrderMatchingService)
			return order.NewEngineResultHandler(postOrderMatchingSvc), nil
		},
	})
}

func addEngine() {
	mustAdd(di.Def{
		Name:  engineService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			redisClient := ctn.Get(RedisClient).(platform.RedisClient)
			logger := ctn.Get(LoggerService).(platform.Logger)
			orderbookProvider := engine.NewRedisOrderBookProvider(redisClient, logger) //todo should be this here  or be an independent service??
			engineResultHandler := ctn.Get(EngineResultHandler).(order.EngineResultHandler)
			configs := ctn.Get(ConfigService).(platform.Configs)
			env := configs.GetEnv()
			srv := engine.NewEngine(
				redisClient,
				orderbookProvider,
				engineResultHandler,
				logger,
				env,
			)
			return srv, nil
		},
	})
}

func addEngineCommunicator() {
	mustAdd(di.Def{
		Name:  engineCommunicator,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			forceTradrInst := ctn.Get(forceTrader).(order.ForceTrader)
			engineSvc := ctn.Get(engineService).(engine.Engine)
			return order.NewEngineCommunicator(forceTradrInst, engineSvc), nil
		},
	})
}

func addOrderEventsHandler() {
	mustAdd(di.Def{
		Name:  orderEventsHandler,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			orderRedisMgr := ctn.Get(orderRedisManager).(order.RedisManager)
			decisionMgr := ctn.Get(decisionManager).(order.DecisionManager)
			centrifugoMgr := ctn.Get(centrifugoManager).(communication.CentrifugoManager)
			externalExchangeOrderSvc := ctn.Get(externalExchangeOrderService).(externalexchange.OrderService)
			engineCommInst := ctn.Get(engineCommunicator).(order.EngineCommunicator)
			postOrderMatchingSvc := ctn.Get(postOrderMatchingService).(order.PostOrderMatchingService)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := order.NewOrderEventsHandler(
				orderRedisMgr,
				decisionMgr,
				centrifugoMgr,
				externalExchangeOrderSvc,
				engineCommInst,
				postOrderMatchingSvc,
				logger,
			)
			return srv, nil
		},
	})
}

func addStopOrderSubmissionManager() {
	mustAdd(di.Def{
		Name:  StopOrderSubmissionManager,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			orderRepo := ctn.Get(orderRepository).(order.Repository)
			liveDataSvc := ctn.Get(liveDataService).(livedata.Service)
			orderRedisMgr := ctn.Get(orderRedisManager).(order.RedisManager)
			orderEventsHdl := ctn.Get(orderEventsHandler).(order.EventsHandler)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := order.NewStopOrderSubmissionManager(
				dbCli,
				orderRepo,
				liveDataSvc,
				orderRedisMgr,
				orderEventsHdl,
				logger,
			)
			return srv, nil
		},
	})
}

func addInQueueOrderManager() {
	mustAdd(di.Def{
		Name:  inQueueOrderManager,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			engineSvc := ctn.Get(engineService).(engine.Engine)
			logger := ctn.Get(LoggerService).(platform.Logger)
			return order.NewInQueueOrderManager(engineSvc, logger), nil
		},
	})
}

func addAdminOrderManager() {
	mustAdd(di.Def{
		Name:  adminOrderManager,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			currencySvc := ctn.Get(currencyService).(currency.Service)
			klineSvc := ctn.Get(klineService).(currency.KlineService)
			priceGen := ctn.Get(priceGenerator).(currency.PriceGenerator)
			postOrderMatchingSvc := ctn.Get(postOrderMatchingService).(order.PostOrderMatchingService)
			stopOrderSubmissionManger := ctn.Get(StopOrderSubmissionManager).(order.StopOrderSubmissionManager)
			orderEventsHdl := ctn.Get(orderEventsHandler).(order.EventsHandler)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := order.NewAdminOrderManager(
				currencySvc,
				klineSvc,
				priceGen,
				postOrderMatchingSvc,
				stopOrderSubmissionManger,
				orderEventsHdl,
				logger,
			)
			return srv, nil
		},
	})
}

func addOrderCreateManager() {
	mustAdd(di.Def{
		Name:  orderCreateManager,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			userBalanceSvc := ctn.Get(userBalanceService).(userbalance.Service)
			userLevelSvc := ctn.Get(userLevelService).(user.LevelService)
			priceGen := ctn.Get(priceGenerator).(currency.PriceGenerator)
			srv := order.NewOrderCreateManager(
				dbCli,
				userBalanceSvc,
				userLevelSvc,
				priceGen,
			)
			return srv, nil
		},
	})
}

func addOrderService() {
	mustAdd(di.Def{
		Name:  orderService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			orderRepo := ctn.Get(orderRepository).(order.Repository)
			orderCreateMgr := ctn.Get(orderCreateManager).(order.CreateManager)
			orderEventsHdl := ctn.Get(orderEventsHandler).(order.EventsHandler)
			currencySvc := ctn.Get(currencyService).(currency.Service)
			priceGen := ctn.Get(priceGenerator).(currency.PriceGenerator)
			userBalanceSvc := ctn.Get(userBalanceService).(userbalance.Service)
			orderRedisMgr := ctn.Get(orderRedisManager).(order.RedisManager)
			userConfigSvc := ctn.Get(userConfigService).(user.ConfigService)
			userPermissionManager := ctn.Get(permissionManager).(user.PermissionManager)
			adminOrderMgr := ctn.Get(adminOrderManager).(order.AdminOrderManager)
			engineCommInst := ctn.Get(engineCommunicator).(order.EngineCommunicator)
			configsService := ctn.Get(ConfigService).(platform.Configs)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := order.NewOrderService(
				dbCli,
				orderRepo,
				orderCreateMgr,
				orderEventsHdl,
				currencySvc,
				priceGen,
				userBalanceSvc,
				orderRedisMgr,
				userConfigSvc,
				userPermissionManager,
				adminOrderMgr,
				engineCommInst,
				configsService,
				logger,
			)
			return srv, nil
		},
	})
}

func addBotAggregationService() {
	mustAdd(di.Def{
		Name:  botAggregationService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			redisClient := ctn.Get(RedisClient).(platform.RedisClient)
			return order.NewBotAggregationService(redisClient), nil
		},
	})
}

func addTradeEventsHandler() {
	mustAdd(di.Def{
		Name:  tradeEventsHandler,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			botAggregationSvc := ctn.Get(botAggregationService).(order.BotAggregationService)
			configService := ctn.Get(ConfigService).(platform.Configs)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := order.NewTradeEventsHandler(botAggregationSvc, configService, logger)
			return srv, nil
		},
	})
}

func addTradeService() {
	mustAdd(di.Def{
		Name:  tradeService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			tradeRepo := ctn.Get(tradeRepository).(order.TradeRepository)
			return order.NewTradeService(tradeRepo), nil
		},
	})
}

func addUnmatchedOrderHandler() {
	mustAdd(di.Def{
		Name:  UnmatchedOrderHandler,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			redisClient := ctn.Get(RedisClient).(platform.RedisClient)
			orderRepo := ctn.Get(orderRepository).(order.Repository)
			engineCommInst := ctn.Get(engineCommunicator).(order.EngineCommunicator)
			configService := ctn.Get(ConfigService).(platform.Configs)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := order.NewUnmatchedOrdersHandler(
				redisClient,
				orderRepo,
				engineCommInst,
				configService,
				logger,
			)
			return srv, nil
		},
	})
}

func addExternalExchangeService() {
	mustAdd(di.Def{
		Name:  externalExchangeService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			externalExchangeRepo := ctn.Get(externalExchangeRepository).(externalexchange.Repository)
			redisClient := ctn.Get(RedisClient).(platform.RedisClient)
			httpCli := ctn.Get(httpClient).(platform.HTTPClient)
			priceGen := ctn.Get(priceGenerator).(currency.PriceGenerator)
			configService := ctn.Get(ConfigService).(platform.Configs)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := externalexchange.NewExternalExchangeService(
				externalExchangeRepo,
				redisClient,
				httpCli,
				priceGen,
				configService,
				logger,
			)
			return srv, nil
		},
	})
}

func addExternalExchangeOrderService() {
	mustAdd(di.Def{
		Name:  externalExchangeOrderService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			externalExchangeOrderRepo := ctn.Get(externalExchangeOrderRepository).(externalexchange.OrderRepository)
			externalExchangeSvc := ctn.Get(externalExchangeService).(externalexchange.Service)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := externalexchange.NewOrderService(
				externalExchangeOrderRepo,
				externalExchangeSvc,
				logger,
			)
			return srv, nil
		},
	})
}

func addOrderFromExternalService() {
	mustAdd(di.Def{
		Name:  orderFromExternalService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			orderFromExternalRepo := ctn.Get(orderFromExternalRepository).(externalexchange.OrderFromExternalRepository)
			tradeFromExternalRepo := ctn.Get(tradeFromExternalRepository).(externalexchange.TradeFromExternalRepository)
			srv := externalexchange.NewOrderFromExternalService(
				orderFromExternalRepo,
				tradeFromExternalRepo,
			)
			return srv, nil
		},
	})
}

func addExternalExchangeWsService() {
	mustAdd(di.Def{
		Name:  ExternalExchangeWsService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			wsCli := ctn.Get(wsClient).(platform.WsClient)
			wsDataProc := ctn.Get(wsDataProcessor).(processor.Processor)
			configsService := ctn.Get(ConfigService).(platform.Configs)
			logger := ctn.Get(LoggerService).(platform.Logger)
			currencySvc := ctn.Get(currencyService).(currency.Service)
			srv := externalexchangews.NewService(
				wsCli,
				wsDataProc,
				logger,
				configsService,
				currencySvc,
			)
			return srv, nil
		},
	})
}

func addWSDataProcessor() {
	mustAdd(di.Def{
		Name:  wsDataProcessor,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			redisClient := ctn.Get(RedisClient).(platform.RedisClient)
			liveDataSvc := ctn.Get(liveDataService).(livedata.Service)
			priceGen := ctn.Get(priceGenerator).(currency.PriceGenerator)
			klineSvc := ctn.Get(klineService).(currency.KlineService)
			orderbookSvc := ctn.Get(orderbookService).(orderbook.Service)
			centrifugoMgr := ctn.Get(centrifugoManager).(communication.CentrifugoManager)
			stopOrderSubmissionManager := ctn.Get(StopOrderSubmissionManager).(order.StopOrderSubmissionManager)
			inQueueOrderMgr := ctn.Get(inQueueOrderManager).(order.InQueueOrderManager)
			queueMgr := ctn.Get(queueManager).(communication.QueueManager)
			logger := ctn.Get(LoggerService).(platform.Logger)
			currencySvc := ctn.Get(currencyService).(currency.Service)
			srv := processor.NewProcessor(
				redisClient,
				liveDataSvc,
				priceGen,
				klineSvc,
				orderbookSvc,
				centrifugoMgr,
				stopOrderSubmissionManager,
				inQueueOrderMgr,
				queueMgr,
				logger,
				currencySvc,
			)
			return srv, nil
		},
	})
}
