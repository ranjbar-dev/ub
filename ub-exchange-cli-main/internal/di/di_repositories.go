package di

import (
	"exchange-go/internal/order"
	"exchange-go/internal/platform"
	"exchange-go/internal/repository"

	"github.com/sarulabs/di"
	"gorm.io/gorm"
)

// DI registrations for data-access repositories.
// All repositories depend on dbClient (GORM *gorm.DB).
// Some also depend on cacheService (pairRepository, userRepository,
// countryRepository, appVersionRepository, userLevelRepository).
// orderRedisManager is included here as it is also a data-access layer — it depends on redisClient.
// Repositories are stateless and can be registered in any order relative to each other.
func addOrderRepository() {
	mustAdd(di.Def{
		Name:  orderRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			return repository.NewOrderRepository(dbCli), nil
		},
	})
}

func addUserBalanceRepository() {
	mustAdd(di.Def{
		Name:  userBalanceRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			return repository.NewUserBalanceRepository(dbCli), nil
		},
	})
}

func addCurrencyRepository() {
	mustAdd(di.Def{
		Name:  currencyRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			return repository.NewCurrencyRepository(dbCli), nil
		},
	})
}

func addFavoritePairRepository() {
	mustAdd(di.Def{
		Name:  favoritePairRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			return repository.NewFavoritePairRepository(dbCli), nil
		},
	})
}

func addPairRepository() {
	mustAdd(di.Def{
		Name:  pairRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			cache := ctn.Get(cacheService).(platform.Cache)
			return repository.NewPairRepository(dbCli, cache), nil
		},
	})
}

func addPermissionRepository() {
	mustAdd(di.Def{
		Name:  permissionRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			return repository.NewPermissionRepository(dbCli), nil
		},
	})
}

func addUsersPermissionsRepository() {
	mustAdd(di.Def{
		Name:  usersPermissionsRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			return repository.NewUsersPermissionsRepository(dbCli), nil
		},
	})
}

func addUserRepository() {
	mustAdd(di.Def{
		Name:  userRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			cacheSvc := ctn.Get(cacheService).(platform.Cache)
			return repository.NewUserRepository(dbCli, cacheSvc), nil
		},
	})
}

func addUserProfileRepository() {
	mustAdd(di.Def{
		Name:  userProfileRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			return repository.NewUserProfileRepository(dbCli), nil
		},
	})
}

func addProfileImageRepository() {
	mustAdd(di.Def{
		Name:  profileImageRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			return repository.NewProfileImageRepository(dbCli), nil
		},
	})
}

func addCountryRepository() {
	mustAdd(di.Def{
		Name:  countryRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			cache := ctn.Get(cacheService).(platform.Cache)
			return repository.NewCountryRepository(dbCli, cache), nil
		},
	})
}

func addLoginHistoryRepository() {
	mustAdd(di.Def{
		Name:  loginHistoryRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			return repository.NewUserLoginHistoryRepository(dbCli), nil
		},
	})
}

func addUserWalletBalanceRepository() {
	mustAdd(di.Def{
		Name:  userWalletBalanceRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			db := ctn.Get(dbClient).(*gorm.DB)
			return repository.NewUserWalletBalanceRepository(db), nil
		},
	})
}

func addKlineSyncRepository() {
	mustAdd(di.Def{
		Name:  klineSyncRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			db := ctn.Get(dbClient).(*gorm.DB)
			return repository.NewKlineSyncRepository(db), nil
		},
	})
}

func addAppVersionRepository() {
	mustAdd(di.Def{
		Name:  appVersionRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			db := ctn.Get(dbClient).(*gorm.DB)
			cache := ctn.Get(cacheService).(platform.Cache)
			return repository.NewAppVersionRepository(db, cache), nil
		},
	})
}

func addInternalTransferRepository() {
	mustAdd(di.Def{
		Name:  internalTransferRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			return repository.NewInternalTransferRepository(dbCli), nil
		},
	})
}

func addPaymentRepository() {
	mustAdd(di.Def{
		Name:  paymentRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			return repository.NewPaymentRepository(dbCli), nil
		},
	})
}

func addUserConfigRepository() {
	mustAdd(di.Def{
		Name:  userConfigRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			return repository.NewUserConfigRepository(dbCli), nil
		},
	})
}

func addUserWithdrawAddressRepository() {
	mustAdd(di.Def{
		Name:  userWithdrawAddressRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			return repository.NewUserWithdrawAddressRepository(dbCli), nil
		},
	})
}

func addTradeFromExternalRepository() {
	mustAdd(di.Def{
		Name:  tradeFromExternalRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			return repository.NewTradeFromExternalRepository(dbCli), nil
		},
	})
}

func addOrderFromExternalRepository() {
	mustAdd(di.Def{
		Name:  orderFromExternalRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			return repository.NewOrderFromExternalRepository(dbCli), nil
		},
	})
}

func addTradeRepository() {
	mustAdd(di.Def{
		Name:  tradeRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			return repository.NewTradeRepository(dbCli), nil
		},
	})
}

func addConfigurationRepository() {
	mustAdd(di.Def{
		Name:  configurationRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			return repository.NewConfigurationRepository(dbCli), nil
		},
	})
}

func addExternalExchangeOrderRepository() {
	mustAdd(di.Def{
		Name:  externalExchangeOrderRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			return repository.NewExternalExchangeOrderRepository(dbCli), nil
		},
	})
}

func addExternalExchangeRepository() {
	mustAdd(di.Def{
		Name:  externalExchangeRepository,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			dbCli := ctn.Get(dbClient).(*gorm.DB)
			return repository.NewExternalExchangeRepository(dbCli), nil
		},
	})
}

func addOrderRedisManager() {
	mustAdd(di.Def{
		Name:  orderRedisManager,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			redisClient := ctn.Get(RedisClient).(platform.RedisClient)
			return order.NewRedisManager(redisClient), nil
		},
	})

}
