package api

import (
	"context"
	"exchange-go/internal/api/adminhandler"
	"exchange-go/internal/api/handler"
	"exchange-go/internal/api/middleware"
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
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

const (
	DefaultAddr  = "0.0.0.0:8000"
	AdminAddr    = "0.0.0.0:8001"
	ReadTimeout  = 30 * time.Second
	WriteTimeout = 30 * time.Second
)


type Services struct {
	CountryService             country.Service
	Configs                    platform.Configs
	ConfigurationService       configuration.Service
	CurrencyService            currency.Service
	AuthService                auth.Service
	UserWithdrawAddressService userwithdrawaddress.Service
	WalletService              wallet.Service
	OrderService               order.Service
	UserBalanceService         userbalance.Service
	PaymentService             payment.Service
	CentrifugoTokenService     auth.CentrifugoTokenService
	UserService                user.Service
	OrderBookService           orderbook.Service
	DB                         *gorm.DB
	RedisClient                platform.RedisClient
}

// HTTPServer manages the public and admin HTTP API servers for the exchange platform.
type HTTPServer interface {
	// ListenAndServe starts the public-facing API server on the given address.
	ListenAndServe(address string) error
	// ListenAndServeAdmin starts the admin API server on the given address.
	ListenAndServeAdmin(address string) error
	// Shutdown gracefully stops both servers, waiting for in-flight requests to complete.
	Shutdown(ctx context.Context) error
	// GetEngine returns the HTTP handler for the public API (used for testing).
	GetEngine() http.Handler
	// GetAdminEngine returns the HTTP handler for the admin API (used for testing).
	GetAdminEngine() http.Handler
}

type httpServer struct {
	server      *http.Server
	engine      *gin.Engine
	adminServer *http.Server
	adminEngine *gin.Engine
	services    Services
}

func (s *httpServer) ListenAndServe(address string) error {
	s.server.Addr = address
	return s.server.ListenAndServe()
}

func (s *httpServer) ListenAndServeAdmin(address string) error {
	s.adminServer.Addr = address
	return s.adminServer.ListenAndServe()
}

func (s *httpServer) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return s.adminServer.Shutdown(ctx)
}

func (s *httpServer) GetEngine() http.Handler {
	return s.engine
}

func (s *httpServer) GetAdminEngine() http.Handler {
	return s.adminEngine
}

func NewHTTPServer(services Services, logger platform.Logger) HTTPServer {
	binding.Validator = new(defaultValidator)
	apiRouter := gin.New()
	adminRouter := gin.New()

	apiRouter.Use(globalRecover(logger, services.Configs))
	apiRouter.Use(middleware.RateLimiter(rate.Limit(10), 20))
	apiRouter.Use(middleware.BodyLimit(1 << 20))
	apiRouter.Use(middleware.Metrics())

	adminRouter.Use(globalRecover(logger, services.Configs))
	env := services.Configs.GetEnv()
	if strings.ToUpper(env) == platform.EnvProd {
		gin.SetMode(gin.ReleaseMode)
	}

	server := &http.Server{
		Addr:           DefaultAddr,
		Handler:        apiRouter,
		ReadTimeout:    ReadTimeout,
		WriteTimeout:   WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	adminServer := &http.Server{
		Addr:           AdminAddr,
		Handler:        adminRouter,
		ReadTimeout:    ReadTimeout,
		WriteTimeout:   WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	s := httpServer{
		server:      server,
		engine:      apiRouter,
		adminServer: adminServer,
		adminEngine: adminRouter,
		services:    services,
	}

	s.registerRoutes()
	s.registerAdminRoutes()
	return &s

}

func (s *httpServer) registerRoutes() {
	r := s.engine
	r.GET("/health", handler.HealthCheck())
	r.GET("/ready", handler.ReadinessCheck(s.services.DB, s.services.RedisClient))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	if s.services.Configs.GetEnv() == platform.EnvProd {
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"https://unitedbit.com", "https://www.unitedbit.com", "https://admin.unitedbit.com", "https://app.unitedbit.com", "https://m.unitedbit.com", "https://dev-m.unitedbit.com"},
			AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
			AllowHeaders:     []string{"Content-Type,Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With, Access-Control-Allow-Origin"},
			AllowCredentials: true,
		}))
	} else {
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
			AllowHeaders:     []string{"Content-Type,Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With, Access-Control-Allow-Origin"},
			AllowCredentials: true,
		}))
	}
	r.Static("/assets", "./assets")

	v1 := r.Group("/api/v1")
	{
		s.registerAuthRoutes(v1)
		s.registerCentrifugoRoutes(v1)
		s.registerMainDataRoutes(v1)
		s.registerCurrencyRoutes(v1)
		s.registerWithdrawAddressRoutes(v1)
		s.registerOrderRoutes(v1)
		v1.GET("/order-book", handler.OrderBook(s.services.OrderBookService))
		v1.GET("/trade-book", handler.TradeBook(s.services.OrderBookService))
		s.registerTradeRoutes(v1)
		s.registerUserBalanceRoutes(v1)
		s.registerPaymentRoutes(v1)
		s.registerUserRoutes(v1)
		s.registerUserProfileImageRoutes(v1)
	}
}

func (s *httpServer) registerAdminRoutes() {
	r := s.adminEngine
	v1 := r.Group("/api/v1")
	v1.Use(middleware.AdminAuthMiddleware(s.services.AuthService))
	{
		orderRoutes := v1.Group("/order")
		{
			orderRoutes.POST("/fulfill", adminhandler.FulFillOrder(s.services.OrderService))
		}

		paymentRoutes := v1.Group("/payment")
		{
			paymentRoutes.POST("/callback", adminhandler.Callback(s.services.PaymentService))
			paymentRoutes.POST("/update-withdraw", adminhandler.UpdateWithdraw(s.services.PaymentService))
			paymentRoutes.POST("/update-deposit", adminhandler.UpdateDeposit(s.services.PaymentService))
		}

		userBalanceRoutes := v1.Group("/user-balance")
		{
			userBalanceRoutes.POST("/update", adminhandler.UpdateUserBalance(s.services.UserBalanceService))
		}
	}

}

func globalRecover(logger platform.Logger, configs platform.Configs) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func(c *gin.Context) {
			message := http.StatusText(http.StatusInternalServerError)
			if rec := recover(); rec != nil {
				if configs.GetEnv() != platform.EnvProd {
				}
				// that recovery also handle XHR's
				// you need handle it
				err := fmt.Errorf("error 500")
				logger.Error2(fmt.Sprintf("error  500 in global recover %v", rec), err,
					zap.String("service", "httpServer"),
					zap.String("method", "globalRecover"),
				)
				response := handler.APIResponse{
					Status:  false,
					Message: message,
				}
				c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			}
		}(c)
		c.Next()
	}

}
