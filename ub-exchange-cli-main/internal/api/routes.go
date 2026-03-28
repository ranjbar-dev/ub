package api

import (
	"exchange-go/internal/api/handler"
	"exchange-go/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

// registerAuthRoutes sets up /auth endpoints (login, register, password reset).
func (s *httpServer) registerAuthRoutes(v1 *gin.RouterGroup) {
	auth := v1.Group("/auth")
	{
		auth.POST("/login", handler.Login(s.services.AuthService))
		auth.POST("/register", handler.Register(s.services.AuthService))
		auth.POST("/forgot-password", handler.ForgotPassword(s.services.AuthService))
		auth.POST("/forgot-password/update", handler.ForgotPasswordUpdate(s.services.AuthService))
		auth.POST("/verify", handler.VerifyEmail(s.services.AuthService))
	}
}

// registerMqttAuthRoutes sets up /emqtt endpoints for MQTT broker authentication.
func (s *httpServer) registerMqttAuthRoutes(v1 *gin.RouterGroup) {
	mqttAuth := v1.Group("/emqtt")
	{
		mqttAuth.POST("/login", handler.MqttLogin(s.services.MqttAuthService))
		mqttAuth.POST("/acl", handler.MqttACL(s.services.MqttAuthService))
		mqttAuth.POST("/superuser", handler.MqttSuperUser(s.services.MqttAuthService))
	}
}

// registerMainDataRoutes sets up /main-data public endpoints (health check, country list, app version).
func (s *httpServer) registerMainDataRoutes(v1 *gin.RouterGroup) {
	mainData := v1.Group("/main-data")
	{
		mainData.GET("/check", handler.Check())
		mainData.GET("/country-list", handler.Countries(s.services.CountryService))
		mainData.GET("/common", handler.GetRecaptchaKey(s.services.ConfigurationService))
		mainData.GET("/version", handler.GetAppVersion(s.services.ConfigurationService))
		mainData.POST("/contact-us", handler.ContactUs(s.services.ConfigurationService))
	}
}

// registerCurrencyRoutes sets up /currencies endpoints (pairs, fees, favorites).
func (s *httpServer) registerCurrencyRoutes(v1 *gin.RouterGroup) {
	currencies := v1.Group("/currencies")
	{
		currencies.GET("", handler.GetCurrencies(s.services.CurrencyService))
		currencies.GET("/pairs", handler.GetPairs(s.services.CurrencyService))
		currencies.GET("/pairs-list", handler.GetPairsList(s.services.CurrencyService))
		currencies.GET("/pairs-statistic", handler.GetPairsStatistic(s.services.CurrencyService))
		currencies.GET("/pairs-ratio", handler.GetPairRatio(s.services.CurrencyService))
		currencies.GET("/fees", handler.GetFees(s.services.CurrencyService))

		favorite := currencies.Group("/favorite")
		favorite.Use(middleware.AuthMiddleware(s.services.AuthService))
		{
			favorite.POST("", handler.AddOrRemoveFavoritePair(s.services.CurrencyService))
		}

		getFavorites := currencies.Group("/favorite-pairs")
		getFavorites.Use(middleware.AuthMiddleware(s.services.AuthService))
		{
			getFavorites.GET("", handler.GetFavoritePairs(s.services.CurrencyService))
		}
	}
}

// registerWithdrawAddressRoutes sets up /withdraw-address endpoints (auth required).
func (s *httpServer) registerWithdrawAddressRoutes(v1 *gin.RouterGroup) {
	userWithdrawAddress := v1.Group("/withdraw-address")
	userWithdrawAddress.Use(middleware.AuthMiddleware(s.services.AuthService))
	{
		userWithdrawAddress.GET("", handler.GetWithdrawAddresses(s.services.UserWithdrawAddressService))
		userWithdrawAddress.GET("former-addresses", handler.GetFormerAddresses(s.services.UserWithdrawAddressService))
		userWithdrawAddress.POST("/new", handler.NewWithdrawAddress(s.services.UserWithdrawAddressService))
		userWithdrawAddress.POST("/favorite", handler.AddToFavorites(s.services.UserWithdrawAddressService))
		userWithdrawAddress.POST("/delete", handler.Delete(s.services.UserWithdrawAddressService))
	}
}

// registerOrderRoutes sets up /order endpoints (create, cancel, history — auth required).
func (s *httpServer) registerOrderRoutes(v1 *gin.RouterGroup) {
	orderRoutes := v1.Group("/order")
	orderRoutes.Use(middleware.AuthMiddleware(s.services.AuthService))
	{
		orderRoutes.POST("/create", handler.CreateOrder(s.services.OrderService))
		orderRoutes.POST("/cancel", handler.CancelOrder(s.services.OrderService))
		orderRoutes.GET("/open-orders", handler.OpenOrders(s.services.OrderService))
		orderRoutes.GET("/history", handler.OrdersHistory(s.services.OrderService))
		orderRoutes.GET("/full-history", handler.FullOrdersHistory(s.services.OrderService))
		orderRoutes.GET("/detail", handler.GetOrderDetail(s.services.OrderService))
	}
}

// registerTradeRoutes sets up /trade endpoints (trade history — auth required).
func (s *httpServer) registerTradeRoutes(v1 *gin.RouterGroup) {
	trade := v1.Group("/trade")
	trade.Use(middleware.AuthMiddleware(s.services.AuthService))
	{
		trade.GET("/history", handler.TradesHistory(s.services.OrderService))
		trade.GET("/full-history", handler.FullTradesHistory(s.services.OrderService))
	}
}

// registerUserBalanceRoutes sets up /user-balance endpoints (auth required).
func (s *httpServer) registerUserBalanceRoutes(v1 *gin.RouterGroup) {
	userBalance := v1.Group("/user-balance")
	userBalance.Use(middleware.AuthMiddleware(s.services.AuthService))
	{
		userBalance.GET("/pair-balance", handler.PairBalances(s.services.UserBalanceService))
		userBalance.GET("/balance", handler.AllBalances(s.services.UserBalanceService))
		userBalance.GET("/withdraw-deposit", handler.WithdrawAndDeposit(s.services.UserBalanceService))
		userBalance.POST("/auto-exchange", handler.SetAutoExchange(s.services.UserBalanceService))
	}
}

// registerPaymentRoutes sets up /crypto-payment endpoints (withdraw, deposit — auth required).
func (s *httpServer) registerPaymentRoutes(v1 *gin.RouterGroup) {
	cryptoPayment := v1.Group("/crypto-payment")
	cryptoPayment.Use(middleware.AuthMiddleware(s.services.AuthService))
	{
		cryptoPayment.GET("", handler.GetPayments(s.services.PaymentService))
		cryptoPayment.GET("/detail", handler.GetPaymentDetail(s.services.PaymentService))
		cryptoPayment.POST("/pre-withdraw", handler.PreWithdraw(s.services.PaymentService))
		cryptoPayment.POST("/withdraw", handler.Withdraw(s.services.PaymentService))
		cryptoPayment.POST("/cancel", handler.Cancel(s.services.PaymentService))
	}
}

// registerUserRoutes sets up /user endpoints (profile, 2FA, SMS, password — auth required).
func (s *httpServer) registerUserRoutes(v1 *gin.RouterGroup) {
	userRoutes := v1.Group("/user")
	userRoutes.Use(middleware.AuthMiddleware(s.services.AuthService))
	{
		userRoutes.POST("/set-user-profile", handler.SetUserProfile(s.services.UserService))
		userRoutes.GET("/get-user-profile", handler.GetUserProfile(s.services.UserService))
		userRoutes.GET("/user-data", handler.GetUserData(s.services.UserService))
		userRoutes.GET("/google-2fa-barcode", handler.Get2FaBarcode(s.services.UserService))
		userRoutes.POST("/google-2fa-enable", handler.Enable2Fa(s.services.UserService))
		userRoutes.POST("/google-2fa-disable", handler.Disable2Fa(s.services.UserService))
		userRoutes.POST("/change-password", handler.ChangePassword(s.services.UserService))
		userRoutes.POST("/sms-send", handler.SendSms(s.services.UserService))
		userRoutes.POST("/sms-enable", handler.EnableSms(s.services.UserService))
		userRoutes.POST("/sms-disable", handler.DisableSms(s.services.UserService))
		userRoutes.POST("/send-verification-email", handler.SendVerificationEmail(s.services.UserService))
	}
}

// registerUserProfileImageRoutes sets up /user-profile-image upload/delete endpoints (auth required).
func (s *httpServer) registerUserProfileImageRoutes(v1 *gin.RouterGroup) {
	userProfileImagesRoutes := v1.Group("/user-profile-image")
	userProfileImagesRoutes.Use(middleware.AuthMiddleware(s.services.AuthService))
	{
		userProfileImagesRoutes.POST("/multiple-upload", handler.MultipleUpload(s.services.UserService))
		userProfileImagesRoutes.POST("/delete", handler.DeleteProfileImage(s.services.UserService))
	}
}
