package di

import (
	"exchange-go/internal/auth"
	"exchange-go/internal/communication"
	"exchange-go/internal/configuration"
	"exchange-go/internal/country"
	"exchange-go/internal/currency"
	"exchange-go/internal/currency/candle"
	"exchange-go/internal/externalexchange"
	"exchange-go/internal/jwt"
	"exchange-go/internal/livedata"
	"exchange-go/internal/order"
	"exchange-go/internal/orderbook"
	"exchange-go/internal/payment"
	"exchange-go/internal/platform"
	"exchange-go/internal/repository"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"exchange-go/internal/userwithdrawaddress"
	"exchange-go/internal/wallet"

	"github.com/sarulabs/di"
	"gorm.io/gorm"
)

// DI registrations for domain services.
// Services depend on repositories and infrastructure. Key dependency chains:
//   - currencyService  ← currencyRepository, pairRepository, liveDataService, priceGenerator, klineService
//   - userService      ← userRepository, userProfileRepository, communicationService, jwtService
//   - userBalanceService ← userBalanceRepository, currencyService, permissionManager, walletService, userService
//   - authService      ← userService, userBalanceService, permissionManager, authEventsHandler
//   - paymentService   ← currencyService, walletService, userService, userBalanceService, autoExchangeManager
//   - autoExchangeManager ← orderCreateManager, orderEventsHandler (cross-domain, registered in di_order_services.go)
func addLiveDataService() {
	mustAdd(di.Def{
		Name:  liveDataService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			redisClient := ctn.Get(RedisClient).(platform.RedisClient)
			return livedata.NewLiveDataService(redisClient), nil
		},
	})
}

func addPriceGenerator() {
	mustAdd(di.Def{
		Name:  priceGenerator,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			liveDataSvc := ctn.Get(liveDataService).(livedata.Service)
			klineSvc := ctn.Get(klineService).(currency.KlineService)
			pairRepo := ctn.Get(pairRepository).(currency.PairRepository)
			return currency.NewPriceGenerator(liveDataSvc, klineSvc, pairRepo), nil
		},
	})
}

func addCandleGRPCClient() {
	mustAdd(di.Def{
		Name:  candleGRPCClient,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			configsService := ctn.Get(ConfigService).(platform.Configs)
			loggerService := ctn.Get(LoggerService).(platform.Logger)
			return candle.NewCandleGRPCClient(configsService, loggerService), nil
		},
	})
}

func addKlineService() {
	mustAdd(di.Def{
		Name:  klineService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			klineSyncRepo := ctn.Get(klineSyncRepository).(currency.KlineSyncRepository)
			liveDataSvc := ctn.Get(liveDataService).(livedata.Service)
			candleGRPCCli := ctn.Get(candleGRPCClient).(candle.CandleGRPCClient)
			return currency.NewKlineService(klineSyncRepo, liveDataSvc, candleGRPCCli), nil
		},
	})
}

func addCurrencyService() {
	mustAdd(di.Def{
		Name:  currencyService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			currencyRepo := ctn.Get(currencyRepository).(currency.Repository)
			liveDataSvc := ctn.Get(liveDataService).(livedata.Service)
			priceGen := ctn.Get(priceGenerator).(currency.PriceGenerator)
			pairRepo := ctn.Get(pairRepository).(currency.PairRepository)
			klineSvc := ctn.Get(klineService).(currency.KlineService)
			favoritePairRepo := ctn.Get(favoritePairRepository).(currency.FavoritePairRepository)
			configService := ctn.Get(ConfigService).(platform.Configs)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := currency.NewCurrencyService(
				currencyRepo,
				liveDataSvc,
				priceGen,
				pairRepo,
				klineSvc,
				favoritePairRepo,
				configService,
				logger,
			)
			return srv, nil
		},
	})
}

func addOrderbookService() {
	mustAdd(di.Def{
		Name:  orderbookService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			liveDataSvc := ctn.Get(liveDataService).(livedata.Service)
			httpCli := ctn.Get(httpClient).(platform.HTTPClient)
			currencySvc := ctn.Get(currencyService).(currency.Service)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := orderbook.NewOrderBookService(
				httpCli,
				liveDataSvc,
				currencySvc,
				logger,
			)
			return srv, nil
		},
	})
}

func addCountryService() {
	mustAdd(di.Def{
		Name:  countryService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			countryRepo := ctn.Get(countryRepository).(country.Repository)
			configService := ctn.Get(ConfigService).(platform.Configs)
			return country.NewCountryService(countryRepo, configService), nil
		},
	})
}

func addCommunicationService() {
	mustAdd(di.Def{
		Name:  communicationService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			queueMgr := ctn.Get(queueManager).(communication.QueueManager)
			logger := ctn.Get(LoggerService).(platform.Logger)
			return communication.NewCommunicationService(queueMgr, logger), nil
		},
	})
}

func addPhoneConfirmationManager() {
	mustAdd(di.Def{
		Name:  phoneConfirmationManager,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			redisClient := ctn.Get(RedisClient).(platform.RedisClient)
			communicationSvc := ctn.Get(communicationService).(communication.Service)
			return user.NewPhoneConfirmationManager(redisClient, communicationSvc), nil
		},
	})
}

func addJWTService() {
	mustAdd(di.Def{
		Name:  jwtService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(ConfigService).(platform.Configs)
			jwtHdl := ctn.Get(jwtHandler).(platform.JwtHandler)
			return jwt.NewJwtService(configService, jwtHdl), nil
		},
	})
}

func addUserService() {
	mustAdd(di.Def{
		Name:  userService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			userRepo := ctn.Get(userRepository).(user.Repository)
			userProfileRepo := ctn.Get(userProfileRepository).(user.ProfileRepository)
			ProfileImageRepository := ctn.Get(profileImageRepository).(user.ProfileImageRepository)
			countrySvc := ctn.Get(countryService).(country.Service)
			twoFaMgr := ctn.Get(twoFaManager).(user.TwoFaManager)
			passwordEnc := ctn.Get(passwordEncoder).(platform.PasswordEncoder)
			communicationSvc := ctn.Get(communicationService).(communication.Service)
			phoneConfirmationManger := ctn.Get(phoneConfirmationManager).(user.PhoneConfirmationManager)
			jwtSvc := ctn.Get(jwtService).(jwt.Service)
			configService := ctn.Get(ConfigService).(platform.Configs)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := user.NewUserService(
				dbCli,
				userRepo,
				userProfileRepo,
				ProfileImageRepository,
				countrySvc,
				twoFaMgr,
				passwordEnc,
				communicationSvc,
				phoneConfirmationManger,
				jwtSvc,
				configService,
				logger,
			)
			return srv, nil
		},
	})
}

func addUserLevelRepository() {
	mustAdd(di.Def{
		Name:  userLevelRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			cache := ctn.Get(cacheService).(platform.Cache)
			return repository.NewUserLevelRepository(dbCli, cache), nil
		},
	})
}

func addUserLevelService() {
	mustAdd(di.Def{
		Name:  userLevelService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			userLevelRepo := ctn.Get(userLevelRepository).(user.LevelRepository)
			return user.NewUserLevelService(userLevelRepo), nil
		},
	})
}

func addForgotPasswordManager() {
	mustAdd(di.Def{
		Name:  forgotPasswordManager,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			redisClient := ctn.Get(RedisClient).(platform.RedisClient)
			communicationSvc := ctn.Get(communicationService).(communication.Service)
			configService := ctn.Get(ConfigService).(platform.Configs)
			srv := user.NewForgotPasswordManager(
				redisClient,
				communicationSvc,
				configService,
			)
			return srv, nil
		},
	})
}

func addRecaptchaManager() {
	mustAdd(di.Def{
		Name:  recaptchaManager,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			httpCli := ctn.Get(httpClient).(platform.HTTPClient)
			configsService := ctn.Get(ConfigService).(platform.Configs)
			logger := ctn.Get(LoggerService).(platform.Logger)
			ubCaptchaMgr := ctn.Get(ubCaptchaManager).(user.UbCaptchaManager)
			srv := user.NewRecaptchaManager(
				httpCli,
				configsService,
				logger,
				ubCaptchaMgr,
			)
			return srv, nil
		},
	})
}

func addAuthService() {
	mustAdd(di.Def{
		Name:  authService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			db := ctn.Get(dbClient).(*gorm.DB)
			userRepo := ctn.Get(userRepository).(user.Repository)
			userLevelSvc := ctn.Get(userLevelService).(user.LevelService)
			userPermissionManger := ctn.Get(permissionManager).(user.PermissionManager)
			userBalanceSvc := ctn.Get(userBalanceService).(userbalance.Service)
			jwtSvc := ctn.Get(jwtService).(jwt.Service)
			passwordEnc := ctn.Get(passwordEncoder).(platform.PasswordEncoder)
			communicationSvc := ctn.Get(communicationService).(communication.Service)
			authEventsHdl := ctn.Get(authEventsHandler).(auth.EventsHandler)
			forgotPasswordMgr := ctn.Get(forgotPasswordManager).(user.ForgotPasswordManager)
			recaptchaMgr := ctn.Get(recaptchaManager).(user.RecaptchaManager)
			twoFaMgr := ctn.Get(twoFaManager).(user.TwoFaManager)
			configsService := ctn.Get(ConfigService).(platform.Configs)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := auth.NewAuthService(
				db,
				userRepo,
				userLevelSvc,
				userPermissionManger,
				userBalanceSvc,
				jwtSvc,
				passwordEnc,
				communicationSvc,
				authEventsHdl,
				forgotPasswordMgr,
				recaptchaMgr,
				twoFaMgr,
				configsService,
				logger,
			)

			return srv, nil
		},
	})
}

func addCentrifugoTokenService() {
	mustAdd(di.Def{
		Name:  centrifugoTokenService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			configsService := ctn.Get(ConfigService).(platform.Configs)
			return auth.NewCentrifugoTokenService(configsService), nil
		},
	})
}

func addLoginHistoryService() {
	mustAdd(di.Def{
		Name:  loginHistoryService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			loginHistoryRepo := ctn.Get(loginHistoryRepository).(user.LoginHistoryRepository)
			return user.NewLoginHistoryService(loginHistoryRepo), nil
		},
	})
}

func addAuthEventsHandler() {
	mustAdd(di.Def{
		Name:  authEventsHandler,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			loginHistorySvc := ctn.Get(loginHistoryService).(user.LoginHistoryService)
			communicationSvc := ctn.Get(communicationService).(communication.Service)
			userSvc := ctn.Get(userService).(user.Service)
			configService := ctn.Get(ConfigService).(platform.Configs)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := auth.NewAuthEventsHandler(
				loginHistorySvc,
				communicationSvc,
				userSvc,
				configService,
				logger,
			)
			return srv, nil
		},
	})
}

func addUbCaptchaManager() {
	mustAdd(di.Def{
		Name:  ubCaptchaManager,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			logger := ctn.Get(LoggerService).(platform.Logger)
			return user.NewUbCaptchaManager(logger), nil
		},
	})
}

func addWalletAuthorizationService() {
	mustAdd(di.Def{
		Name:  walletAuthorizationService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			redisClient := ctn.Get(RedisClient).(platform.RedisClient)
			httpCli := ctn.Get(httpClient).(platform.HTTPClient)
			configService := ctn.Get(ConfigService).(platform.Configs)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := wallet.NewAuthorizationService(
				redisClient,
				logger,
				httpCli,
				configService,
			)
			return srv, nil
		},
	})
}

func addWalletService() {
	mustAdd(di.Def{
		Name:  walletService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			walletAuthSvc := ctn.Get(walletAuthorizationService).(wallet.AuthorizationService)
			httpCli := ctn.Get(httpClient).(platform.HTTPClient)
			configService := ctn.Get(ConfigService).(platform.Configs)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := wallet.NewWalletService(
				walletAuthSvc,
				httpCli,
				configService,
				logger,
			)
			return srv, nil
		},
	})
}

func addUserBalanceService() {
	mustAdd(di.Def{
		Name:  userBalanceService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			userBalanceRepo := ctn.Get(userBalanceRepository).(userbalance.Repository)
			currencySvc := ctn.Get(currencyService).(currency.Service)
			priceGen := ctn.Get(priceGenerator).(currency.PriceGenerator)
			permissionMgr := ctn.Get(permissionManager).(user.PermissionManager)
			walletSvc := ctn.Get(walletService).(wallet.Service)
			configService := ctn.Get(ConfigService).(platform.Configs)
			logger := ctn.Get(LoggerService).(platform.Logger)
			userSvc := ctn.Get(userService).(user.Service)
			userWalletBalanceRepo := ctn.Get(userWalletBalanceRepository).(userbalance.UserWalletBalanceRepository)
			srv := userbalance.NewBalanceService(
				dbCli,
				userBalanceRepo,
				currencySvc,
				priceGen,
				permissionMgr,
				walletSvc,
				userSvc,
				userWalletBalanceRepo,
				configService,
				logger,
			)
			return srv, nil
		},
	})
}

func addUserConfigService() {
	mustAdd(di.Def{
		Name:  userConfigService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			userConfigRepo := ctn.Get(userConfigRepository).(user.ConfigRepository)
			return user.NewUserConfigService(userConfigRepo), nil
		},
	})
}

func addUserWithdrawAddressService() {
	mustAdd(di.Def{
		Name:  userWithdrawAddressService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			userWithdrawAddressRepo := ctn.Get(userWithdrawAddressRepository).(userwithdrawaddress.Repository)
			currencySvc := ctn.Get(currencyService).(currency.Service)
			walletSvc := ctn.Get(walletService).(wallet.Service)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := userwithdrawaddress.NewUserWithdrawAddressService(
				dbCli,
				userWithdrawAddressRepo,
				currencySvc,
				walletSvc,
				logger,
			)
			return srv, nil
		},
	})
}

func addWithdrawEmailConfirmationManager() {
	mustAdd(di.Def{
		Name:  withdrawEmailConfirmationManager,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			redisClient := ctn.Get(RedisClient).(platform.RedisClient)
			communicationSvc := ctn.Get(communicationService).(communication.Service)
			return payment.NewWithdrawEmailConfirmationManager(redisClient, communicationSvc), nil
		},
	})
}

func addPaymentService() {
	mustAdd(di.Def{
		Name:  paymentService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			paymentRepo := ctn.Get(paymentRepository).(payment.Repository)
			currencySvc := ctn.Get(currencyService).(currency.Service)
			walletSvc := ctn.Get(walletService).(wallet.Service)
			userConfigSvc := ctn.Get(userConfigService).(user.ConfigService)
			twoFaMgr := ctn.Get(twoFaManager).(user.TwoFaManager)
			withdrawEmailConfirmMgr := ctn.Get(withdrawEmailConfirmationManager).(payment.WithdrawEmailConfirmationManager)
			permissionMgr := ctn.Get(permissionManager).(user.PermissionManager)
			userSvc := ctn.Get(userService).(user.Service)
			userBalanceSvc := ctn.Get(userBalanceService).(userbalance.Service)
			userWithdrawAddrSvc := ctn.Get(userWithdrawAddressService).(userwithdrawaddress.Service)
			communicationSvc := ctn.Get(communicationService).(communication.Service)
			priceGen := ctn.Get(priceGenerator).(currency.PriceGenerator)
			internalTransferSvc := ctn.Get(internalTransferService).(payment.InternalTransferService)
			externalExchangeSvc := ctn.Get(externalExchangeService).(externalexchange.Service)
			autoExchangeMgr := ctn.Get(autoExchangeManager).(payment.AutoExchangeManager)
			centrifugoMgr := ctn.Get(centrifugoManager).(communication.CentrifugoManager)
			configService := ctn.Get(ConfigService).(platform.Configs)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := payment.NewPaymentService(
				dbCli,
				paymentRepo,
				currencySvc,
				walletSvc,
				userConfigSvc,
				twoFaMgr,
				withdrawEmailConfirmMgr,
				permissionMgr,
				userWithdrawAddrSvc,
				userSvc,
				userBalanceSvc,
				communicationSvc,
				priceGen,
				internalTransferSvc,
				externalExchangeSvc,
				autoExchangeMgr,
				centrifugoMgr,
				configService,
				logger,
			)
			return srv, nil
		},
	})
}

func addInternalTransferService() {
	mustAdd(di.Def{
		Name:  internalTransferService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			internalTransferRepo := ctn.Get(internalTransferRepository).(payment.InternalTransferRepository)
			return payment.NewInternalTransferService(internalTransferRepo), nil
		},
	})
}

func addAutoExchangeManager() {
	mustAdd(di.Def{
		Name:  autoExchangeManager,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			paymentRepo := ctn.Get(paymentRepository).(payment.Repository)
			orderCreateMgr := ctn.Get(orderCreateManager).(order.CreateManager)
			orderEventsHdl := ctn.Get(orderEventsHandler).(order.EventsHandler)
			userSvc := ctn.Get(userService).(user.Service)
			currencySvc := ctn.Get(currencyService).(currency.Service)
			priceGen := ctn.Get(priceGenerator).(currency.PriceGenerator)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := payment.NewAutoExchangeManger(
				dbCli,
				paymentRepo,
				orderCreateMgr,
				orderEventsHdl,
				userSvc,
				currencySvc,
				priceGen,
				logger,
			)
			return srv, nil
		},
	})
}

func addConfigurationService() {
	mustAdd(di.Def{
		Name:  configurationService,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			configurationRepo := ctn.Get(configurationRepository).(configuration.Repository)
			appVersionRepo := ctn.Get(appVersionRepository).(configuration.AppVersionRepository)
			communicationSvc := ctn.Get(communicationService).(communication.Service)
			configsService := ctn.Get(ConfigService).(platform.Configs)
			logger := ctn.Get(LoggerService).(platform.Logger)
			srv := configuration.NewConfigurationService(
				configurationRepo,
				appVersionRepo,
				communicationSvc,
				configsService,
				logger,
			)
			return srv, nil
		},
	})
}

func addTwoFaManager() {
	mustAdd(di.Def{
		Name:  twoFaManager,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return user.NewTwoFaManager(), nil
		},
	})
}

func addPermissionManager() {
	mustAdd(di.Def{
		Name:  permissionManager,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			usersPermissionsRepo := ctn.Get(usersPermissionsRepository).(user.UsersPermissionsRepository)
			permissionRepo := ctn.Get(permissionRepository).(user.PermissionRepository)
			return user.NewUserPermissionManager(usersPermissionsRepo, permissionRepo), nil
		},
	})
}
