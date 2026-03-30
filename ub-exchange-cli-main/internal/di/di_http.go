package di

import (
	"exchange-go/internal/api"
	"exchange-go/internal/auth"
	"exchange-go/internal/configuration"
	"exchange-go/internal/country"
	"exchange-go/internal/currency"
	"exchange-go/internal/order"
	"exchange-go/internal/orderbook"
	"exchange-go/internal/payment"
	"exchange-go/internal/platform"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"exchange-go/internal/userwithdrawaddress"
	"exchange-go/internal/wallet"

	"github.com/sarulabs/di"
)

// DI registrations for the HTTP API server.
// The HTTP server is the outermost layer — it depends on all domain services
// (authService, orderService, paymentService, userService, currencyService,
// userBalanceService, walletService, configurationService, orderbookService,
// userWithdrawAddressService, countryService, centrifugoTokenService).
// It must be registered after all domain services are registered.
func addHTTPServer() {
	mustAdd(di.Def{
		Name:  HTTPServer,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			logger := ctn.Get(LoggerService).(platform.Logger)
			services := api.Services{
				CountryService:             ctn.Get(countryService).(country.Service),
				Configs:                    ctn.Get(ConfigService).(platform.Configs),
				ConfigurationService:       ctn.Get(configurationService).(configuration.Service),
				CurrencyService:            ctn.Get(currencyService).(currency.Service),
				AuthService:                ctn.Get(authService).(auth.Service),
				UserWithdrawAddressService: ctn.Get(userWithdrawAddressService).(userwithdrawaddress.Service),
				WalletService:              ctn.Get(walletService).(wallet.Service),
				OrderService:               ctn.Get(orderService).(order.Service),
				UserBalanceService:         ctn.Get(userBalanceService).(userbalance.Service),
				PaymentService:             ctn.Get(paymentService).(payment.Service),
				CentrifugoTokenService:     ctn.Get(centrifugoTokenService).(auth.CentrifugoTokenService),
				UserService:                ctn.Get(userService).(user.Service),
				OrderBookService:           ctn.Get(orderbookService).(orderbook.Service),
			}

			return api.NewHTTPServer(services, logger), nil
		},
	})
}
