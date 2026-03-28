package di

import (
	"fmt"

	"github.com/sarulabs/di"
)

const (
	cacheService                     = "cacheService"
	dbClient                         = "dbClient"
	wsClient                         = "wsClient"
	mqttClient                       = "mqttClient"
	httpClient                       = "httpClient"
	mqttManager                      = "mqttManager"
	recaptchaManager                 = "recaptchaManager"
	liveDataService                  = "liveDataService"
	priceGenerator                   = "priceGenerator"
	klineService                     = "klineService"
	mqttAuthService                  = "mqttAuthService"
	authService                      = "authService"
	candleGRPCClient                 = "candleGRPCClient"
	klineSyncRepository              = "klineSyncRepository"
	appVersionRepository             = "appVersionRepository"
	orderRepository                  = "orderRepository"
	internalTransferRepository       = "internalTransferRepository"
	paymentRepository                = "paymentRepository"
	userConfigRepository             = "userConfigRepository"
	userWithdrawAddressRepository    = "userWithdrawAddressRepository"
	tradeFromExternalRepository      = "tradeFromExternalRepository"
	orderFromExternalRepository      = "orderFromExternalRepository"
	orderFromExternalService         = "orderFromExternalService"
	internalTransferService          = "internalTransferService"
	paymentService                   = "paymentService"
	userConfigService                = "userConfigService"
	userWithdrawAddressService       = "userWithdrawAddressService"
	withdrawEmailConfirmationManager = "withdrawEmailConfirmationManager"
	userBalanceRepository            = "userBalanceRepository"
	externalExchangeOrderRepository  = "externalExchangeOrderRepository"
	externalExchangeRepository       = "externalExchangeRepository"
	currencyRepository               = "currencyRepository"
	pairRepository                   = "pairRepository"
	favoritePairRepository           = "favoritePairRepository"
	usersPermissionsRepository       = "usersPermissionsRepository"
	permissionRepository             = "permissionRepository"
	orderRedisManager                = "orderRedisManager"
	wsDataProcessor                  = "wsDataProcessor"
	orderbookService                 = "orderbookService"
	inQueueOrderManager              = "inQueueOrderManager"
	orderEventsHandler               = "orderEventsHandler"
	tradeEventsHandler               = "tradeEventsHandler"
	userService                      = "userService"
	countryService                   = "countryService"
	twoFaManager                     = "twoFaManager"
	passwordEncoder                  = "passwordEncoder"
	rabbitmqClient                   = "rabbitmqClient"
	queueManager                     = "queueManager"
	communicationService             = "communicationService"
	phoneConfirmationManager         = "phoneConfirmationManager"
	jwtHandler                       = "jwtHandler"
	jwtService                       = "jwtService"
	userLevelService                 = "userLevelService"
	userLevelRepository              = "userLevelRepository"
	userRepository                   = "userRepository"
	userProfileRepository            = "userProfileRepository"
	profileImageRepository           = "profileImageRepository"
	loginHistoryRepository           = "loginHistoryRepository"
	loginHistoryService              = "loginHistoryService"
	authEventsHandler                = "authEventsHandler"
	countryRepository                = "countryRepository"
	configurationRepository          = "configurationRepository"
	botAggregationService            = "botAggregationService"
	engineCommunicator               = "engineCommunicator"
	forceTrader                      = "forceTrader"
	engineService                    = "engineService"
	postOrderMatchingService         = "postOrderMatchingService"
	userBalanceService               = "userBalanceService"
	walletAuthorizationService       = "walletAuthorizationService"
	walletService                    = "walletService"
	tradeService                     = "tradeService"
	orderService                     = "orderService"
	configurationService             = "configurationService"
	adminOrderManager                = "adminOrderManager"
	tradeRepository                  = "tradeRepository"
	permissionManager                = "permissionManager"
	currencyService                  = "currencyService"
	externalExchangeOrderService     = "externalExchangeOrderService"
	externalExchangeService          = "externalExchangeService"
	decisionManager                  = "decisionManager"
	ubCaptchaManager                 = "ubCaptchaManager"
	forgotPasswordManager            = "forgotPasswordManager"
	orderCreateManager               = "orderCreateManager"
	userWalletBalanceRepository      = "userWalletBalanceRepository"
	autoExchangeManager              = "autoExchangeManager"

	//exported could be used in main functions
	ConfigService                             = "configService"
	HTTPServer                                = "httpServer"
	UnmatchedOrderHandler                     = "unmatchedOrderHandler"
	RedisClient                               = "redisClient"
	LoggerService                             = "loggerService"
	StopOrderSubmissionManager                = "stopOrderSubmissionManager"
	EngineResultHandler                       = "engineResultHandler"
	ExternalExchangeWsService                 = "externalExchangeWsService"
	CheckWithdrawalsInExternalExchangeCommand = "checkWithdrawalsInExternalExchangeCommand"
	UpdateOrdersInExternalExchangeCommand     = "updateOrdersInExternalExchangeCommand"
	KlineSyncCommand                          = "klineSyncCommand"
	GenerateKlineSyncCommand                  = "generateKlineSyncCommand"
	RetrieveExternalOrdersToRedisCommand      = "retrieveExternalOrdersToRedisCommand"
	SubmitBotAggregatedOrderCommand           = "submitBotAggregatedOrderCommand"
	RetrieveOpenOrdersToRedisCommand          = "retrieveOpenOrdersToRedisCommand"
	GenerateAddressCommand                    = "generateAddressCommand"
	InitializeBalanceCommand                  = "initializeBalanceCommand"
	WsHealthCheckCommand                      = "wsHealthCheckCommand"
	SetUserLevelCommand                       = "setUserLevelCommand"
	DeleteCacheCommand                        = "deleteCacheCommand"
	UbCaptchaDecryptionCommand                = "ubCaptchaDecryptionCommand"
	UbCaptchaEncryptionCommand                = "ubCaptchaEncryptionCommand"
	UbCaptchaKeyGeneratorCommand              = "ubCaptchaKeyGeneratorCommand"
	UpdateUserWalletBalancesCommand           = "ubUpdateUserWalletBalancesCommand"
)

// Service dependency chain (critical path — order matching):
//
//	configService → all services via ctn.Get(ConfigService)
//	loggerService → most services for structured logging
//	dbClient      → all repositories (GORM *gorm.DB)
//	redisClient   → orderRedisManager, liveDataService, engine, postOrderMatchingService
//	cacheService  → pairRepository, userRepository, countryRepository, appVersionRepository, userLevelRepository
//
// Order matching critical path:
//
//	currencyRepository  → currencyService → orderbookService, orderService, paymentService
//	userRepository      → userService     → authService, userBalanceService, paymentService
//	orderRepository     → postOrderMatchingService, orderService, stopOrderSubmissionManager
//	botAggregationService  → tradeEventsHandler → postOrderMatchingService
//	forceTrader            → postOrderMatchingService, engineCommunicator
//	postOrderMatchingService → engineResultHandler → engine
//	engine                   → engineCommunicator → orderEventsHandler → orderService
//	orderCreateManager       → orderService, autoExchangeManager
//	stopOrderSubmissionManager → wsDataProcessor
//	inQueueOrderManager        → wsDataProcessor
var builder *di.Builder

// mustAdd registers a service definition with the DI builder.
// Panics if registration fails, ensuring misconfiguration is caught at startup.
func mustAdd(def di.Def) {
	if err := builder.Add(def); err != nil {
		panic(fmt.Sprintf("di: failed to register service %q: %v", def.Name, err))
	}
}

func NewContainer() di.Container {
	if builder == nil {
		builder, _ = di.NewBuilder(di.App)
	}

	// === Infrastructure (core platform: config, DB, cache, messaging clients) ===
	addConfigService()
	addCacheService()        // depends on: configService, loggerService
	addDBClient()            // depends on: configService
	addLogger()              // depends on: configService
	addWSClient()            // no deps
	addMQTTClient()          // depends on: configService, loggerService
	addRedisClient()         // depends on: configService
	addMQTTManager()         // depends on: mqttClient

	// === Live Data + Kline Services (depend on: redisClient, dbClient) ===
	addLiveDataService()     // depends on: redisClient
	addPriceGenerator()      // depends on: liveDataService, klineService, pairRepository
	addKlineSyncRepository() // depends on: dbClient
	addCandleGRPCClient()    // depends on: configService, loggerService
	addKlineService()        // depends on: klineSyncRepository, liveDataService, candleGRPCClient
	addHTTPClient()          // no deps

	// === Order Book + Order Core (depend on: redisClient, dbClient) ===
	addOrderbookService()  // depends on: liveDataService, httpClient, currencyService, loggerService
	addOrderRepository()   // depends on: dbClient
	addOrderRedisManager() // depends on: redisClient
	addDecisionManager()   // depends on: configService

	// === External Exchange (depend on: dbClient, redisClient, httpClient) ===
	addExternalExchangeOrderRepository() // depends on: dbClient
	addExternalExchangeRepository()      // depends on: dbClient
	addExternalExchangeService()         // depends on: externalExchangeRepository, redisClient, httpClient, priceGenerator, configService, loggerService
	addExternalExchangeOrderService()    // depends on: externalExchangeOrderRepository, externalExchangeService, loggerService
	addForceTrader()                     // depends on: priceGenerator, currencyService

	// === Currency + Balance Domain (depend on: dbClient, cacheService, redisClient) ===
	addUserBalanceRepository()      // depends on: dbClient
	addCurrencyRepository()         // depends on: dbClient
	addFavoritePairRepository()     // depends on: dbClient
	addPairRepository()             // depends on: dbClient, cacheService
	addCurrencyService()            // depends on: currencyRepository, liveDataService, priceGenerator, pairRepository, klineService, favoritePairRepository, configService, loggerService
	addPermissionRepository()       // depends on: dbClient
	addUsersPermissionsRepository() // depends on: dbClient
	addPermissionManager()          // depends on: usersPermissionsRepository, permissionRepository
	addWalletAuthorizationService() // depends on: redisClient, httpClient, configService, loggerService
	addWalletService()              // depends on: walletAuthorizationService, httpClient, configService, loggerService
	addUserBalanceService()         // depends on: dbClient, userBalanceRepository, currencyService, priceGenerator, permissionManager, walletService, userService, userWalletBalanceRepository, configService, loggerService
	addBotAggregationService()      // depends on: redisClient
	addTradeEventsHandler()         // depends on: botAggregationService, configService, loggerService

	// === User Domain (depend on: dbClient, cacheService, communicationService) ===
	addUserRepository()           // depends on: dbClient, cacheService
	addUserProfileRepository()    // depends on: dbClient
	addProfileImageRepository()   // depends on: dbClient
	addCountryRepository()        // depends on: dbClient, cacheService
	addCountryService()           // depends on: countryRepository, configService
	addTwoFaManager()             // no deps
	addPasswordEncoder()          // no deps
	addRabbitmqClient()           // depends on: configService, loggerService
	addQueueManager()             // depends on: rabbitmqClient, loggerService
	addCommunicationService()     // depends on: queueManager, loggerService
	addPhoneConfirmationManager() // depends on: redisClient, communicationService
	addJWTHandler()               // no deps
	addJWTService()               // depends on: configService, jwtHandler
	addUserService()              // depends on: dbClient, userRepository, userProfileRepository, profileImageRepository, countryService, twoFaManager, passwordEncoder, communicationService, phoneConfirmationManager, jwtService, configService, loggerService
	addUserLevelRepository()      // depends on: dbClient, cacheService
	addUserLevelService()         // depends on: userLevelRepository

	// === Order Engine Pipeline (critical path for order matching) ===
	addPostOrderMatchingService()   // depends on: dbClient, orderRepository, userBalanceService, forceTrader, priceGenerator, tradeEventsHandler, mqttManager, redisClient, currencyService, userService, userLevelService, configService, loggerService
	addEngineResultHandler()        // depends on: postOrderMatchingService
	addEngine()                     // depends on: redisClient, engineResultHandler, configService, loggerService
	addEngineCommunicator()         // depends on: forceTrader, engineService
	addOrderEventsHandler()         // depends on: orderRedisManager, decisionManager, mqttManager, externalExchangeOrderService, engineCommunicator, postOrderMatchingService, loggerService
	addStopOrderSubmissionManager() // depends on: dbClient, orderRepository, liveDataService, orderRedisManager, orderEventsHandler, loggerService
	addInQueueOrderManager()        // depends on: engineService, loggerService
	addWSDataProcessor()            // depends on: redisClient, liveDataService, priceGenerator, klineService, orderbookService, mqttManager, stopOrderSubmissionManager, inQueueOrderManager, queueManager, loggerService, currencyService
	addUnmatchedOrderHandler()      // depends on: redisClient, orderRepository, engineCommunicator, configService, loggerService

	// === Auth Domain (depend on: userService, communicationService) ===
	addDeleteCacheCommand()           // depends on: cacheService, loggerService
	addLoginHistoryRepository()       // depends on: dbClient
	addLoginHistoryService()          // depends on: loginHistoryRepository
	addAuthEventsHandler()            // depends on: loginHistoryService, communicationService, userService, configService, loggerService
	addUbCaptchaManager()             // depends on: loggerService
	addUbCaptchaDecryptionCommand()   // depends on: ubCaptchaManager, loggerService
	addUbCaptchaEncryptionCommand()   // depends on: ubCaptchaManager, loggerService
	addUbCaptchaKeyGeneratorCommand() // depends on: ubCaptchaManager, loggerService
	addForgotPasswordManager()        // depends on: redisClient, communicationService, configService
	addAppVersionRepository()         // depends on: dbClient, cacheService
	addRecaptchaManager()             // depends on: httpClient, configService, loggerService, ubCaptchaManager
	addAuthService()                  // depends on: dbClient, userRepository, userLevelService, permissionManager, userBalanceService, jwtService, passwordEncoder, communicationService, authEventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configService, loggerService
	addMqttAuthService()              // depends on: authService, configService, loggerService

	// === Payment Domain (depend on: userService, currencyService, walletService) ===
	addInternalTransferRepository()       // depends on: dbClient
	addInternalTransferService()          // depends on: internalTransferRepository
	addPaymentRepository()                // depends on: dbClient
	addUserConfigRepository()             // depends on: dbClient
	addUserConfigService()                // depends on: userConfigRepository
	addUserWithdrawAddressRepository()    // depends on: dbClient
	addUserWithdrawAddressService()       // depends on: dbClient, userWithdrawAddressRepository, currencyService, walletService, loggerService
	addWithdrawEmailConfirmationManager() // depends on: redisClient, communicationService
	addPaymentService()                   // depends on: dbClient, paymentRepository, currencyService, walletService, userConfigService, twoFaManager, withdrawEmailConfirmationManager, permissionManager, userService, userBalanceService, userWithdrawAddressService, communicationService, priceGenerator, internalTransferService, externalExchangeService, autoExchangeManager, mqttManager, configService, loggerService

	// === CLI Commands (external exchange, kline, trade, balance) ===
	addCheckWithdrawalsInExternalExchangeCommand() // depends on: paymentService, externalExchangeService, internalTransferService, configService, loggerService
	addTradeFromExternalRepository()               // depends on: dbClient
	addOrderFromExternalRepository()               // depends on: dbClient
	addOrderFromExternalService()                  // depends on: orderFromExternalRepository, tradeFromExternalRepository
	addUpdateOrdersInExternalExchangeCommand()     // depends on: currencyService, orderFromExternalService, externalExchangeService, configService, loggerService
	addSyncKlineCommand()                          // depends on: currencyService, klineService, externalExchangeService, queueManager, loggerService
	addGenerateKlineSyncCommand()                  // depends on: currencyService, klineService, loggerService
	addTradeRepository()                           // depends on: dbClient
	addTradeService()                              // depends on: tradeRepository
	addRetrieveExternalOrdersToRedisCommand()      // depends on: externalExchangeOrderService, botAggregationService, tradeService, loggerService
	addSubmitBotAggregatedOrderCommand()           // depends on: currencyService, botAggregationService, liveDataService, externalExchangeOrderService, loggerService
	addRetrieveOpenOrdersToRedisCommand()          // depends on: orderRepository, orderRedisManager, engineCommunicator, loggerService
	addGenerateAddressCommand()                    // depends on: userRepository, userBalanceRepository, userBalanceService, loggerService
	addInitializeBalanceCommand()                  // depends on: userRepository, currencyService, userBalanceRepository, userBalanceService, loggerService
	addAdminOrderManager()                         // depends on: currencyService, klineService, priceGenerator, postOrderMatchingService, stopOrderSubmissionManager, orderEventsHandler, loggerService
	addWsHealthCheckCommand()                      // depends on: liveDataService, loggerService
	addOrderCreateManager()                        // depends on: dbClient, userBalanceService, userLevelService, priceGenerator
	addOrderService()                              // depends on: dbClient, orderRepository, orderCreateManager, orderEventsHandler, currencyService, priceGenerator, userBalanceService, orderRedisManager, userConfigService, permissionManager, adminOrderManager, engineCommunicator, configService, loggerService

	// === Configuration Service ===
	addConfigurationRepository() // depends on: dbClient
	addConfigurationService()    // depends on: configurationRepository, appVersionRepository, communicationService, configService, loggerService

	// === HTTP Server (depends on: all domain services) ===
	addHTTPServer()

	// === External WebSocket + Remaining CLI Commands ===
	addExternalExchangeWsService()         // depends on: wsClient, wsDataProcessor, configService, loggerService, currencyService
	addSetUserLevelCommand()               // depends on: dbClient, userService, tradeService, pairRepository, currencyRepository, klineService, userLevelService, loggerService
	addUbUpdateUserWalletBalancesCommand() // depends on: userService, userBalanceService, currencyService, loggerService
	addUserWalletBalanceRepository()       // depends on: dbClient
	addAutoExchangeManager()              // depends on: dbClient, paymentRepository, orderCreateManager, orderEventsHandler, userService, currencyService, priceGenerator, loggerService

	return builder.Build()
}
