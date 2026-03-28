// Package order_test tests the PostOrderMatchingService. Covers:
//   - All-limit order matching with no partial remainder
//   - All-limit order matching with partial order without remaining quantity
//   - Single market order matching with only partial without remaining
//   - All-limit order matching with partial order and remaining quantity
//   - Partial order matching triggered from admin fulfillment
//   - Error when matched orders are not in OPEN status
//   - Error when partial order is a market order
//
// Test data: sqlmock MySQL DB, mocked order/trade repositories, user balance,
// currency, MQTT, live data, trade events handler, and platform configs.
package order_test

import (
	"exchange-go/internal/currency"
	"exchange-go/internal/mocks"
	"exchange-go/internal/order"
	"exchange-go/internal/platform"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// postOrderMatchingQueryMatcher is a sqlmock.QueryMatcher that accepts any SQL
// query, allowing tests to focus on business logic rather than exact SQL strings.
type postOrderMatchingQueryMatcher struct {
}

// Match always returns nil, effectively disabling SQL query matching so that
// any expected SQL statement is considered a match against any actual SQL.
func (postOrderMatchingQueryMatcher) Match(expectedSQL, actualSQL string) error {
	return nil
}

func TestPostOrderMatchingService_HandlePostOrderMatching_AllLimits_NoPartial(t *testing.T) {
	qm := postOrderMatchingQueryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	//these are for updating orders themselves
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(1, 1))

	//these are for creating and updating child order
	dbMock.ExpectExec("INSERT INTO orders").WillReturnResult(sqlmock.NewResult(4, 1))
	dbMock.ExpectExec("INSERT INTO orders_extra_info").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(1, 1))

	//2 userbalance update for every order demanded and payedBy coin
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))

	//3 transaction for every order demanded ,payedBy and fee
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))

	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))

	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))

	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))

	dbMock.ExpectExec("INSERT INTO trades").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO trades").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	orderFields := []order.MatchingNeededQueryFields{
		{
			OrderID:           1,
			Price:             "50000",
			DemandedAmount:    "2.0",
			OrderType:         "BUY",
			OrderExchangeType: "LIMIT",
			PayedByAmount:     "100000",
			Path:              "1,",
			Status:            order.StatusOpen,
			CreatedAt:         time.Now().Format("2006-01-02 15:04:05"),
			UserID:            1,
			UserAgentInfo:     "",
		},
		{
			OrderID:           2,
			Price:             "50000",
			DemandedAmount:    "50000",
			OrderType:         "SELL",
			OrderExchangeType: "LIMIT",
			PayedByAmount:     "1.0",
			Path:              "2,",
			Status:            order.StatusOpen,
			CreatedAt:         time.Now().Format("2006-01-02 15:04:05"),
			UserID:            2,
			UserAgentInfo:     "",
		},
		{
			OrderID:           3,
			Price:             "50000",
			DemandedAmount:    "50000",
			OrderType:         "SELL",
			OrderExchangeType: "LIMIT",
			PayedByAmount:     "1.0",
			Path:              "2,",
			Status:            order.StatusOpen,
			CreatedAt:         time.Now().Format("2006-01-02 15:04:05"),
			UserID:            3,
			UserAgentInfo:     "",
		},
	}
	orderRepo := new(mocks.OrderRepository)
	orderRepo.On("GetOrdersDataByIdsWithJoinUsingTx", mock.Anything, []int64{1, 2, 3}).Once().Return(orderFields)

	basisCoin := currency.Coin{
		ID:   1,
		Name: "USDT",
	}

	dependentCoin := currency.Coin{
		ID:   2,
		Name: "BTC",
	}

	userBalanceService := new(mocks.UserBalanceService)
	ubs := []userbalance.UserBalance{
		{
			ID:           1,
			UserID:       1,
			CoinID:       1,
			Coin:         basisCoin,
			Amount:       "150000",
			FrozenAmount: "100000",
		},
		{
			ID:           2,
			UserID:       1,
			CoinID:       2,
			Coin:         dependentCoin,
			Amount:       "1",
			FrozenAmount: "0",
		},

		{
			ID:           3,
			UserID:       2,
			CoinID:       1,
			Coin:         basisCoin,
			Amount:       "0",
			FrozenAmount: "0",
		},
		{
			ID:           4,
			UserID:       2,
			CoinID:       2,
			Coin:         dependentCoin,
			Amount:       "2.0",
			FrozenAmount: "1.0",
		},

		{
			ID:           5,
			UserID:       3,
			CoinID:       1,
			Coin:         basisCoin,
			Amount:       "0",
			FrozenAmount: "0",
		},
		{
			ID:           6,
			UserID:       3,
			CoinID:       2,
			Coin:         dependentCoin,
			Amount:       "2.0",
			FrozenAmount: "1.0",
		},
	}
	userBalanceService.On("GetBalancesOfUsersForCoinsUsingTx", mock.Anything, []int{1, 2, 3}, []int64{1, 2}).Once().Return(ubs)

	forceTrader := new(mocks.ForceTrader)
	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetPrice", mock.Anything, "BTC-USDT").Once().Return("50000", nil)

	tradeEventsHandler := new(mocks.TradeEventsHandler)
	tradeEventsHandler.On("HandleTradesCreation", mock.Anything, mock.Anything).Once().Return()

	mqttManager := new(mocks.MqttManager)
	mqttManager.On("PublishOrderToOpenOrders", mock.Anything, mock.Anything, mock.Anything).Times(3).Return()

	redisClient := new(mocks.RedisClient)
	configs := new(mocks.Configs)
	configs.On("GetEnv").Once().Return(platform.EnvTest)
	configs.On("GetBool", "commitError").Once().Return(false)
	logger := new(mocks.Logger)
	currencyService := new(mocks.CurrencyService)
	pair := currency.Pair{
		ID:              1,
		Name:            "BTC-USDT",
		BasisCoinID:     1,
		BasisCoin:       basisCoin,
		DependentCoinID: 2,
		DependentCoin:   dependentCoin,
		MakerFee:        0.3,
		TakerFee:        0.3,
	}
	currencyService.On("GetPairByName", "BTC-USDT").Once().Return(pair, nil)
	userService := new(mocks.UserService)
	usersData := []user.UsersDataForOrderMatching{
		{
			UserID:             1,
			UserEmail:          "user1@test.com",
			UserLevelID:        1,
			UserPrivateChannel: "",
		},
		{
			UserID:             2,
			UserEmail:          "user2@test.com",
			UserLevelID:        1,
			UserPrivateChannel: "",
		},
		{
			UserID:             3,
			UserEmail:          "user3@test.com",
			UserLevelID:        1,
			UserPrivateChannel: "",
		},
	}
	userService.On("GetUsersDataForOrderMatching", []int{1, 2, 3}).Once().Return(usersData)
	userLevelService := new(mocks.UserLevelService)
	userLevels := []user.Level{
		{
			ID:                 1,
			MakerFeePercentage: 0.5,
			TakerFeePercentage: 1,
		},
	}
	userLevelService.On("GetLevelsByIds", []int64{1}).Once().Return(userLevels)
	poms := order.NewPostOrderMatchingService(db, orderRepo, userBalanceService, forceTrader, priceGenerator, tradeEventsHandler, mqttManager, redisClient, currencyService, userService, userLevelService, configs, logger)
	doneOrdersData := []order.CallBackOrderData{
		{
			ID:                1,
			PairName:          "BTC-USDT",
			OrderType:         "BUY",
			Quantity:          "2.0",
			Price:             "50000",
			Timestamp:         time.Now().Unix(),
			TradedWithOrderID: 2,
			QuantityTraded:    "1.0",
			TradePrice:        "50000",
			MinThresholdPrice: "49500",
			MaxThresholdPrice: "50500",
		},
		{
			ID:                1,
			PairName:          "BTC-USDT",
			OrderType:         "BUY",
			Quantity:          "1.0",
			Price:             "50000",
			Timestamp:         time.Now().Unix(),
			TradedWithOrderID: 3,
			QuantityTraded:    "1.0",
			TradePrice:        "50000",
			MinThresholdPrice: "49500",
			MaxThresholdPrice: "50500",
		},
		{
			ID:                2,
			PairName:          "BTC-USDT",
			OrderType:         "SELL",
			Quantity:          "1.0",
			Price:             "50000",
			Timestamp:         time.Now().Unix(),
			TradedWithOrderID: 1,
			QuantityTraded:    "1.0",
			TradePrice:        "50000",
			MinThresholdPrice: "49500",
			MaxThresholdPrice: "50500",
		},
		{
			ID:                3,
			PairName:          "BTC-USDT",
			OrderType:         "SELL",
			Quantity:          "0.5",
			Price:             "50000",
			Timestamp:         time.Now().Unix(),
			TradedWithOrderID: 1,
			QuantityTraded:    "1.0",
			TradePrice:        "50000",
			MinThresholdPrice: "49500",
			MaxThresholdPrice: "50500",
		},
	}

	matchingResult := poms.HandlePostOrderMatching(doneOrdersData, nil, false)
	assert.Nil(t, matchingResult.Err)
	assert.Nil(t, matchingResult.RemainingPartialOrder)
	time.Sleep(100 * time.Millisecond)
	currencyService.AssertExpectations(t)
	orderRepo.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	tradeEventsHandler.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
	userService.AssertExpectations(t)
	userLevelService.AssertExpectations(t)
}

func TestPostOrderMatchingService_HandlePostOrderMatching_AllLimits_WithPartial_WithoutRemaining(t *testing.T) {
	qm := postOrderMatchingQueryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	//these are for updating orders themselves
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(2, 1))
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(3, 1))

	//these are for creating and updating child order
	dbMock.ExpectExec("INSERT INTO orders").WillReturnResult(sqlmock.NewResult(11, 1))
	dbMock.ExpectExec("INSERT INTO orders_extra_info").WillReturnResult(sqlmock.NewResult(11, 1))
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(11, 1))

	//these are for second child order because of partial order
	dbMock.ExpectExec("INSERT INTO orders").WillReturnResult(sqlmock.NewResult(15, 1))
	dbMock.ExpectExec("INSERT INTO orders_extra_info").WillReturnResult(sqlmock.NewResult(15, 1))
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(15, 1))

	//2 userbalance update for every order demanded and payedBy coin
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))

	//3 transaction for every order demanded ,payedBy and fee,this scenario the main order would have 2 child so we have 5 orders

	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))

	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))

	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))

	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))

	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))

	//one trade for partial which would be traded with bot
	dbMock.ExpectExec("INSERT INTO trades").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO trades").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO trades").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	orderFields := []order.MatchingNeededQueryFields{
		{
			OrderID:           1,
			Price:             "50000",
			DemandedAmount:    "2.5",
			OrderType:         "BUY",
			OrderExchangeType: "LIMIT",
			PayedByAmount:     "125000",
			Path:              "1,",
			Status:            order.StatusOpen,
			CreatedAt:         time.Now().Format("2006-01-02 15:04:05"),
			UserID:            1,
			UserAgentInfo:     "",
		},
		{
			OrderID:           2,
			Price:             "50000",
			DemandedAmount:    "50000",
			OrderType:         "SELL",
			OrderExchangeType: "LIMIT",
			PayedByAmount:     "1.0",
			Path:              "2,",
			Status:            order.StatusOpen,
			CreatedAt:         time.Now().Format("2006-01-02 15:04:05"),
			UserID:            2,
			UserAgentInfo:     "",
		},
		{
			OrderID:           3,
			Price:             "50000",
			DemandedAmount:    "50000",
			OrderType:         "SELL",
			OrderExchangeType: "LIMIT",
			PayedByAmount:     "1.0",
			Path:              "2,",
			Status:            order.StatusOpen,
			CreatedAt:         time.Now().Format("2006-01-02 15:04:05"),
			UserID:            3,
			UserAgentInfo:     "",
		},
	}
	orderRepo := new(mocks.OrderRepository)
	orderRepo.On("GetOrdersDataByIdsWithJoinUsingTx", mock.Anything, []int64{1, 2, 3}).Once().Return(orderFields)

	basisCoin := currency.Coin{
		ID:   1,
		Name: "USDT",
	}

	dependentCoin := currency.Coin{
		ID:   2,
		Name: "BTC",
	}

	userBalanceService := new(mocks.UserBalanceService)
	ubs := []userbalance.UserBalance{
		{
			ID:           1,
			UserID:       1,
			CoinID:       1,
			Coin:         basisCoin,
			Amount:       "150000",
			FrozenAmount: "100000",
		},
		{
			ID:           2,
			UserID:       1,
			CoinID:       2,
			Coin:         dependentCoin,
			Amount:       "1",
			FrozenAmount: "0",
		},

		{
			ID:           3,
			UserID:       2,
			CoinID:       1,
			Coin:         basisCoin,
			Amount:       "0",
			FrozenAmount: "0",
		},
		{
			ID:           4,
			UserID:       2,
			CoinID:       2,
			Coin:         dependentCoin,
			Amount:       "2.0",
			FrozenAmount: "1.0",
		},

		{
			ID:           5,
			UserID:       3,
			CoinID:       1,
			Coin:         basisCoin,
			Amount:       "0",
			FrozenAmount: "0",
		},
		{
			ID:           6,
			UserID:       3,
			CoinID:       2,
			Coin:         dependentCoin,
			Amount:       "2.0",
			FrozenAmount: "1.0",
		},
	}
	userBalanceService.On("GetBalancesOfUsersForCoinsUsingTx", mock.Anything, []int{1, 2, 3}, []int64{1, 2}).Once().Return(ubs)

	forceTrader := new(mocks.ForceTrader)

	forceTrader.On("ShouldForceTrade", "BTC-USDT", "BUY", "50000").Once().Return(true, nil)
	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetPrice", mock.Anything, "BTC-USDT").Once().Return("50000", nil)

	tradeEventsHandler := new(mocks.TradeEventsHandler)
	tradeEventsHandler.On("HandleTradesCreation", mock.Anything, mock.Anything).Once().Return()

	mqttManager := new(mocks.MqttManager)
	mqttManager.On("PublishOrderToOpenOrders", mock.Anything, mock.Anything, mock.Anything).Times(3).Return()

	redisClient := new(mocks.RedisClient)
	configs := new(mocks.Configs)
	configs.On("GetEnv").Once().Return(platform.EnvTest)
	configs.On("GetBool", "commitError").Once().Return(false)
	logger := new(mocks.Logger)
	currencyService := new(mocks.CurrencyService)
	pair := currency.Pair{
		ID:              1,
		Name:            "BTC-USDT",
		BasisCoinID:     1,
		BasisCoin:       basisCoin,
		DependentCoinID: 2,
		DependentCoin:   dependentCoin,
		MakerFee:        0.3,
		TakerFee:        0.3,
	}
	currencyService.On("GetPairByName", "BTC-USDT").Once().Return(pair, nil)
	userService := new(mocks.UserService)
	usersData := []user.UsersDataForOrderMatching{
		{
			UserID:             1,
			UserEmail:          "user1@test.com",
			UserLevelID:        1,
			UserPrivateChannel: "",
		},
		{
			UserID:             2,
			UserEmail:          "user2@test.com",
			UserLevelID:        1,
			UserPrivateChannel: "",
		},
		{
			UserID:             3,
			UserEmail:          "user3@test.com",
			UserLevelID:        1,
			UserPrivateChannel: "",
		},
	}
	userService.On("GetUsersDataForOrderMatching", []int{1, 2, 3}).Once().Return(usersData)
	userLevelService := new(mocks.UserLevelService)
	userLevels := []user.Level{
		{
			ID:                 1,
			MakerFeePercentage: 0.5,
			TakerFeePercentage: 1,
		},
	}
	userLevelService.On("GetLevelsByIds", []int64{1}).Once().Return(userLevels)
	poms := order.NewPostOrderMatchingService(db, orderRepo, userBalanceService, forceTrader, priceGenerator, tradeEventsHandler, mqttManager, redisClient, currencyService, userService, userLevelService, configs, logger)
	doneOrdersData := []order.CallBackOrderData{
		{
			ID:                1,
			PairName:          "BTC-USDT",
			OrderType:         "BUY",
			Quantity:          "2.5",
			Price:             "50000",
			Timestamp:         time.Now().Unix(),
			TradedWithOrderID: 2,
			QuantityTraded:    "1.0",
			TradePrice:        "50000",
			MinThresholdPrice: "49500",
			MaxThresholdPrice: "50500",
		},
		{
			ID:                1,
			PairName:          "BTC-USDT",
			OrderType:         "BUY",
			Quantity:          "1.0",
			Price:             "50000",
			Timestamp:         time.Now().Unix(),
			TradedWithOrderID: 3,
			QuantityTraded:    "1.0",
			TradePrice:        "50000",
			MinThresholdPrice: "49500",
			MaxThresholdPrice: "50500",
		},
		{
			ID:                2,
			PairName:          "BTC-USDT",
			OrderType:         "SELL",
			Quantity:          "1.0",
			Price:             "50000",
			Timestamp:         time.Now().Unix(),
			TradedWithOrderID: 1,
			QuantityTraded:    "1.0",
			TradePrice:        "50000",
			MinThresholdPrice: "49500",
			MaxThresholdPrice: "50500",
		},
		{
			ID:                3,
			PairName:          "BTC-USDT",
			OrderType:         "SELL",
			Quantity:          "1.0",
			Price:             "50000",
			Timestamp:         time.Now().Unix(),
			TradedWithOrderID: 1,
			QuantityTraded:    "1.0",
			TradePrice:        "50000",
			MinThresholdPrice: "49500",
			MaxThresholdPrice: "50500",
		},
	}

	partial := order.CallBackOrderData{
		ID:                1,
		PairName:          "BTC-USDT",
		OrderType:         "BUY",
		Quantity:          "0.5",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: 0,
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	matchingResult := poms.HandlePostOrderMatching(doneOrdersData, &partial, false)
	assert.Nil(t, matchingResult.Err)
	assert.Nil(t, matchingResult.RemainingPartialOrder)
	currencyService.AssertExpectations(t)
	time.Sleep(50 * time.Millisecond)
	orderRepo.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	tradeEventsHandler.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
	userService.AssertExpectations(t)
	userLevelService.AssertExpectations(t)
}

func TestPostOrderMatchingService_HandlePostOrderMatching_SingleMarketOrder_OnlyPartial_WithoutRemaining(t *testing.T) {
	qm := postOrderMatchingQueryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	//these are for updating orders themselves
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(1, 1))

	//2 userbalance update for every order demanded and payedBy coin
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))

	//3 transaction for every order demanded ,payedBy and fee,this scenario the main order would have 2 child so we have 5 orders
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))

	//one trade for partial which would be traded with bot
	dbMock.ExpectExec("INSERT INTO trades").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	orderFields := []order.MatchingNeededQueryFields{
		{
			OrderID:           1,
			Price:             "",
			DemandedAmount:    "2.5",
			OrderType:         "BUY",
			OrderExchangeType: "MARKET",
			PayedByAmount:     "125000",
			Path:              "1,",
			Status:            order.StatusOpen,
			CreatedAt:         time.Now().Format("2006-01-02 15:04:05"),
			UserID:            1,
			UserAgentInfo:     "",
		},
	}
	orderRepo := new(mocks.OrderRepository)
	orderRepo.On("GetOrdersDataByIdsWithJoinUsingTx", mock.Anything, []int64{1}).Once().Return(orderFields)

	basisCoin := currency.Coin{
		ID:   1,
		Name: "USDT",
	}

	dependentCoin := currency.Coin{
		ID:   2,
		Name: "BTC",
	}

	userBalanceService := new(mocks.UserBalanceService)
	ubs := []userbalance.UserBalance{
		{
			ID:           1,
			UserID:       1,
			CoinID:       1,
			Coin:         basisCoin,
			Amount:       "150000",
			FrozenAmount: "100000",
		},
		{
			ID:           2,
			UserID:       1,
			CoinID:       2,
			Coin:         dependentCoin,
			Amount:       "1",
			FrozenAmount: "0",
		},
	}
	userBalanceService.On("GetBalancesOfUsersForCoinsUsingTx", mock.Anything, []int{1}, []int64{1, 2}).Once().Return(ubs)

	forceTrader := new(mocks.ForceTrader)

	forceTrader.On("ShouldForceTrade", "BTC-USDT", "BUY", "").Once().Return(true, nil)
	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetPrice", mock.Anything, "BTC-USDT").Once().Return("50000", nil)

	tradeEventsHandler := new(mocks.TradeEventsHandler)
	tradeEventsHandler.On("HandleTradesCreation", mock.Anything, mock.Anything).Once().Return()

	mqttManager := new(mocks.MqttManager)
	mqttManager.On("PublishOrderToOpenOrders", mock.Anything, mock.Anything, mock.Anything).Once().Return()

	redisClient := new(mocks.RedisClient)
	configs := new(mocks.Configs)
	configs.On("GetEnv").Once().Return(platform.EnvTest)
	configs.On("GetBool", "commitError").Once().Return(false)
	logger := new(mocks.Logger)
	currencyService := new(mocks.CurrencyService)
	pair := currency.Pair{
		ID:              1,
		Name:            "BTC-USDT",
		BasisCoinID:     1,
		BasisCoin:       basisCoin,
		DependentCoinID: 2,
		DependentCoin:   dependentCoin,
		MakerFee:        0.3,
		TakerFee:        0.3,
	}
	currencyService.On("GetPairByName", "BTC-USDT").Once().Return(pair, nil)
	userService := new(mocks.UserService)
	usersData := []user.UsersDataForOrderMatching{
		{
			UserID:             1,
			UserEmail:          "user1@test.com",
			UserLevelID:        1,
			UserPrivateChannel: "",
		},
	}
	userService.On("GetUsersDataForOrderMatching", []int{1}).Once().Return(usersData)
	userLevelService := new(mocks.UserLevelService)
	userLevels := []user.Level{
		{
			ID:                 1,
			MakerFeePercentage: 0.5,
			TakerFeePercentage: 1,
		},
	}
	userLevelService.On("GetLevelsByIds", []int64{1}).Once().Return(userLevels)
	poms := order.NewPostOrderMatchingService(db, orderRepo, userBalanceService, forceTrader, priceGenerator, tradeEventsHandler, mqttManager, redisClient, currencyService, userService, userLevelService, configs, logger)
	var doneOrdersData []order.CallBackOrderData

	partial := order.CallBackOrderData{
		ID:                1,
		PairName:          "BTC-USDT",
		OrderType:         "BUY",
		Quantity:          "2.5",
		Price:             "",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: 0,
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	matchingResult := poms.HandlePostOrderMatching(doneOrdersData, &partial, false)
	assert.Nil(t, matchingResult.Err)
	assert.Nil(t, matchingResult.RemainingPartialOrder)
	currencyService.AssertExpectations(t)
	time.Sleep(100 * time.Millisecond)
	orderRepo.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	forceTrader.AssertExpectations(t)
	tradeEventsHandler.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
	userService.AssertExpectations(t)
	userLevelService.AssertExpectations(t)
}

func TestPostOrderMatchingService_HandlePostOrderMatching_AllLimits_WithPartial_WithRemaining(t *testing.T) {
	qm := postOrderMatchingQueryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	//these are for updating orders themselves
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(1, 1))

	//these are for creating and updating child order
	dbMock.ExpectExec("INSERT INTO orders").WillReturnResult(sqlmock.NewResult(11, 1))
	dbMock.ExpectExec("INSERT INTO orders_extra_info").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(1, 1))

	dbMock.ExpectExec("INSERT INTO orders").WillReturnResult(sqlmock.NewResult(13, 1))
	dbMock.ExpectExec("INSERT INTO orders_extra_info").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(1, 1))

	//2 userbalance update for every order demanded and payedBy coin
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))

	//3 transaction for every order demanded ,payedBy and fee,this scenario the main order would have 2 child so we have 5 orders

	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))

	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))

	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))

	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))

	//one trade for partial which would be traded with bot
	dbMock.ExpectExec("INSERT INTO trades").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO trades").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	orderFields := []order.MatchingNeededQueryFields{
		{
			OrderID:           1,
			Price:             "50000",
			DemandedAmount:    "2.5",
			OrderType:         "BUY",
			OrderExchangeType: "LIMIT",
			PayedByAmount:     "125000",
			Path:              "1,",
			Status:            order.StatusOpen,
			CreatedAt:         time.Now().Format("2006-01-02 15:04:05"),
			UserID:            1,
			UserAgentInfo:     "",
		},
		{
			OrderID:           2,
			Price:             "50000",
			DemandedAmount:    "50000",
			OrderType:         "SELL",
			OrderExchangeType: "LIMIT",
			PayedByAmount:     "1.0",
			Path:              "2,",
			Status:            order.StatusOpen,
			CreatedAt:         time.Now().Format("2006-01-02 15:04:05"),
			UserID:            2,
			UserAgentInfo:     "",
		},
		{
			OrderID:           3,
			Price:             "50000",
			DemandedAmount:    "50000",
			OrderType:         "SELL",
			OrderExchangeType: "LIMIT",
			PayedByAmount:     "1.0",
			Path:              "2,",
			Status:            order.StatusOpen,
			CreatedAt:         time.Now().Format("2006-01-02 15:04:05"),
			UserID:            3,
			UserAgentInfo:     "",
		},
	}
	orderRepo := new(mocks.OrderRepository)
	orderRepo.On("GetOrdersDataByIdsWithJoinUsingTx", mock.Anything, []int64{1, 2, 3}).Once().Return(orderFields)

	basisCoin := currency.Coin{
		ID:   1,
		Name: "USDT",
	}

	dependentCoin := currency.Coin{
		ID:   2,
		Name: "BTC",
	}

	userBalanceService := new(mocks.UserBalanceService)
	ubs := []userbalance.UserBalance{
		{
			ID:           1,
			UserID:       1,
			CoinID:       1,
			Coin:         basisCoin,
			Amount:       "150000",
			FrozenAmount: "100000",
		},
		{
			ID:           2,
			UserID:       1,
			CoinID:       2,
			Coin:         dependentCoin,
			Amount:       "1",
			FrozenAmount: "0",
		},

		{
			ID:           3,
			UserID:       2,
			CoinID:       1,
			Coin:         basisCoin,
			Amount:       "0",
			FrozenAmount: "0",
		},
		{
			ID:           4,
			UserID:       2,
			CoinID:       2,
			Coin:         dependentCoin,
			Amount:       "2.0",
			FrozenAmount: "1.0",
		},

		{
			ID:           5,
			UserID:       3,
			CoinID:       1,
			Coin:         basisCoin,
			Amount:       "0",
			FrozenAmount: "0",
		},
		{
			ID:           6,
			UserID:       3,
			CoinID:       2,
			Coin:         dependentCoin,
			Amount:       "2.0",
			FrozenAmount: "1.0",
		},
	}
	userBalanceService.On("GetBalancesOfUsersForCoinsUsingTx", mock.Anything, []int{1, 2, 3}, []int64{1, 2}).Once().Return(ubs)

	forceTrader := new(mocks.ForceTrader)

	forceTrader.On("ShouldForceTrade", "BTC-USDT", "BUY", "50000").Once().Return(false, nil)
	forceTrader.On("GetMinAndMaxPrice", "BTC-USDT", "BUY", "50000").Once().Return("49500", "50500", nil)
	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetPrice", mock.Anything, "BTC-USDT").Once().Return("50000", nil)

	tradeEventsHandler := new(mocks.TradeEventsHandler)
	tradeEventsHandler.On("HandleTradesCreation", mock.Anything, mock.Anything).Once().Return()

	mqttManager := new(mocks.MqttManager)
	mqttManager.On("PublishOrderToOpenOrders", mock.Anything, mock.Anything, mock.Anything).Times(3).Return()

	redisClient := new(mocks.RedisClient)
	configs := new(mocks.Configs)
	configs.On("GetEnv").Once().Return(platform.EnvTest)
	configs.On("GetBool", "commitError").Once().Return(false)
	logger := new(mocks.Logger)
	currencyService := new(mocks.CurrencyService)
	pair := currency.Pair{
		ID:              1,
		Name:            "BTC-USDT",
		BasisCoinID:     1,
		BasisCoin:       basisCoin,
		DependentCoinID: 2,
		DependentCoin:   dependentCoin,
		MakerFee:        0.3,
		TakerFee:        0.3,
	}
	currencyService.On("GetPairByName", "BTC-USDT").Once().Return(pair, nil)
	userService := new(mocks.UserService)
	usersData := []user.UsersDataForOrderMatching{
		{
			UserID:             1,
			UserEmail:          "user1@test.com",
			UserLevelID:        1,
			UserPrivateChannel: "",
		},
		{
			UserID:             2,
			UserEmail:          "user2@test.com",
			UserLevelID:        1,
			UserPrivateChannel: "",
		},
		{
			UserID:             3,
			UserEmail:          "user3@test.com",
			UserLevelID:        1,
			UserPrivateChannel: "",
		},
	}
	userService.On("GetUsersDataForOrderMatching", []int{1, 2, 3}).Once().Return(usersData)
	userLevelService := new(mocks.UserLevelService)
	userLevels := []user.Level{
		{
			ID:                 1,
			MakerFeePercentage: 0.5,
			TakerFeePercentage: 1,
		},
	}
	userLevelService.On("GetLevelsByIds", []int64{1}).Once().Return(userLevels)
	poms := order.NewPostOrderMatchingService(db, orderRepo, userBalanceService, forceTrader, priceGenerator, tradeEventsHandler, mqttManager, redisClient, currencyService, userService, userLevelService, configs, logger)
	doneOrdersData := []order.CallBackOrderData{
		{
			ID:                1,
			PairName:          "BTC-USDT",
			OrderType:         "BUY",
			Quantity:          "2.5",
			Price:             "50000",
			Timestamp:         time.Now().Unix(),
			TradedWithOrderID: 2,
			QuantityTraded:    "1.0",
			TradePrice:        "50000",
			MinThresholdPrice: "49500",
			MaxThresholdPrice: "50500",
		},
		{
			ID:                1,
			PairName:          "BTC-USDT",
			OrderType:         "BUY",
			Quantity:          "1.0",
			Price:             "50000",
			Timestamp:         time.Now().Unix(),
			TradedWithOrderID: 3,
			QuantityTraded:    "1.0",
			TradePrice:        "50000",
			MinThresholdPrice: "49500",
			MaxThresholdPrice: "50500",
		},
		{
			ID:                2,
			PairName:          "BTC-USDT",
			OrderType:         "SELL",
			Quantity:          "1.0",
			Price:             "50000",
			Timestamp:         time.Now().Unix(),
			TradedWithOrderID: 1,
			QuantityTraded:    "1.0",
			TradePrice:        "50000",
			MinThresholdPrice: "49500",
			MaxThresholdPrice: "50500",
		},
		{
			ID:                3,
			PairName:          "BTC-USDT",
			OrderType:         "SELL",
			Quantity:          "1.0",
			Price:             "50000",
			Timestamp:         time.Now().Unix(),
			TradedWithOrderID: 1,
			QuantityTraded:    "1.0",
			TradePrice:        "50000",
			MinThresholdPrice: "49500",
			MaxThresholdPrice: "50500",
		},
	}

	partial := order.CallBackOrderData{
		ID:                1,
		PairName:          "BTC-USDT",
		OrderType:         "BUY",
		Quantity:          "0.5",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: 0,
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	matchingResult := poms.HandlePostOrderMatching(doneOrdersData, &partial, false)

	assert.Nil(t, matchingResult.Err)
	assert.Equal(t, int64(1), matchingResult.RemainingPartialOrder.ID)
	assert.Equal(t, "BTC-USDT", matchingResult.RemainingPartialOrder.PairName)
	assert.Equal(t, "BUY", matchingResult.RemainingPartialOrder.OrderType)
	assert.Equal(t, "0.50000000", matchingResult.RemainingPartialOrder.Quantity)
	assert.Equal(t, "50000", matchingResult.RemainingPartialOrder.Price)
	assert.Equal(t, int64(0), matchingResult.RemainingPartialOrder.TradedWithOrderID)
	assert.Equal(t, "", matchingResult.RemainingPartialOrder.QuantityTraded)
	assert.Equal(t, "", matchingResult.RemainingPartialOrder.TradePrice)
	assert.Equal(t, "49500", matchingResult.RemainingPartialOrder.MinThresholdPrice)
	assert.Equal(t, "50500", matchingResult.RemainingPartialOrder.MaxThresholdPrice)

	time.Sleep(100 * time.Millisecond)
	currencyService.AssertExpectations(t)
	orderRepo.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	forceTrader.AssertExpectations(t)
	tradeEventsHandler.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
	userService.AssertExpectations(t)
	userLevelService.AssertExpectations(t)
}

func TestPostOrderMatchingService_HandlePostOrderMatching_Partial_FromAdmin(t *testing.T) {
	qm := postOrderMatchingQueryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	//these are for updating orders themselves
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(1, 1))

	//2 userbalance update for every order demanded and payedBy coin
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))

	//3 transaction for every order demanded ,payedBy and fee,this scenario the main order would have 2 child so we have 5 orders

	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))

	//one trade for partial which would be traded with bot
	dbMock.ExpectExec("INSERT INTO trades").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	orderFields := []order.MatchingNeededQueryFields{
		{
			OrderID:           1,
			Price:             "50000",
			DemandedAmount:    "2.5",
			OrderType:         "BUY",
			OrderExchangeType: "LIMIT",
			PayedByAmount:     "125000",
			Path:              "1,",
			Status:            order.StatusOpen,
			CreatedAt:         time.Now().Format("2006-01-02 15:04:05"),
			UserID:            1,
			UserAgentInfo:     "",
		},
	}
	orderRepo := new(mocks.OrderRepository)
	orderRepo.On("GetOrdersDataByIdsWithJoinUsingTx", mock.Anything, []int64{1}).Once().Return(orderFields)

	basisCoin := currency.Coin{
		ID:   1,
		Name: "USDT",
	}

	dependentCoin := currency.Coin{
		ID:   2,
		Name: "BTC",
	}

	userBalanceService := new(mocks.UserBalanceService)
	ubs := []userbalance.UserBalance{
		{
			ID:           1,
			UserID:       1,
			CoinID:       1,
			Coin:         basisCoin,
			Amount:       "150000",
			FrozenAmount: "100000",
		},
		{
			ID:           2,
			UserID:       1,
			CoinID:       2,
			Coin:         dependentCoin,
			Amount:       "1",
			FrozenAmount: "0",
		},
	}
	userBalanceService.On("GetBalancesOfUsersForCoinsUsingTx", mock.Anything, []int{1}, []int64{1, 2}).Once().Return(ubs)

	forceTrader := new(mocks.ForceTrader)
	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetPrice", mock.Anything, "BTC-USDT").Once().Return("50000", nil)

	tradeEventsHandler := new(mocks.TradeEventsHandler)
	tradeEventsHandler.On("HandleTradesCreation", mock.Anything, mock.Anything).Once().Return()

	mqttManager := new(mocks.MqttManager)
	mqttManager.On("PublishOrderToOpenOrders", mock.Anything, mock.Anything, mock.Anything).Once().Return()
	redisClient := new(mocks.RedisClient)
	configs := new(mocks.Configs)
	configs.On("GetEnv").Once().Return(platform.EnvTest)
	configs.On("GetBool", "commitError").Once().Return(false)
	logger := new(mocks.Logger)
	currencyService := new(mocks.CurrencyService)
	pair := currency.Pair{
		ID:              1,
		Name:            "BTC-USDT",
		BasisCoinID:     1,
		BasisCoin:       basisCoin,
		DependentCoinID: 2,
		DependentCoin:   dependentCoin,
		MakerFee:        0.3,
		TakerFee:        0.3,
	}
	currencyService.On("GetPairByName", "BTC-USDT").Once().Return(pair, nil)
	userService := new(mocks.UserService)
	usersData := []user.UsersDataForOrderMatching{
		{
			UserID:             1,
			UserEmail:          "user1@test.com",
			UserLevelID:        1,
			UserPrivateChannel: "",
		},
	}
	userService.On("GetUsersDataForOrderMatching", []int{1}).Once().Return(usersData)
	userLevelService := new(mocks.UserLevelService)
	userLevels := []user.Level{
		{
			ID:                 1,
			MakerFeePercentage: 0.5,
			TakerFeePercentage: 1,
		},
	}
	userLevelService.On("GetLevelsByIds", []int64{1}).Once().Return(userLevels)
	poms := order.NewPostOrderMatchingService(db, orderRepo, userBalanceService, forceTrader, priceGenerator, tradeEventsHandler, mqttManager, redisClient, currencyService, userService, userLevelService, configs, logger)
	var doneOrdersData []order.CallBackOrderData

	partial := order.CallBackOrderData{
		ID:                1,
		PairName:          "BTC-USDT",
		OrderType:         "BUY",
		Quantity:          "0.5",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: 0,
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}

	matchingResult := poms.HandlePostOrderMatching(doneOrdersData, &partial, true)
	assert.Nil(t, matchingResult.Err)
	assert.Nil(t, matchingResult.RemainingPartialOrder)
	currencyService.AssertExpectations(t)
	time.Sleep(100 * time.Millisecond)
	orderRepo.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	tradeEventsHandler.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
	userService.AssertExpectations(t)
	userLevelService.AssertExpectations(t)
}

func TestPostOrderMatchingService_Error_OrdersAreNotOpen(t *testing.T) {
	qm := postOrderMatchingQueryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()

	orderFields := []order.MatchingNeededQueryFields{
		{
			OrderID:           1,
			Price:             "50000",
			DemandedAmount:    "2.5",
			OrderType:         "BUY",
			OrderExchangeType: "LIMIT",
			PayedByAmount:     "125000",
			Path:              "1,",
			Status:            order.StatusFilled,
			CreatedAt:         time.Now().Format("2006-01-02 15:04:05"),
			UserID:            1,
			UserAgentInfo:     "",
		},
		{
			OrderID:           2,
			Price:             "50000",
			DemandedAmount:    "50000",
			OrderType:         "SELL",
			OrderExchangeType: "LIMIT",
			PayedByAmount:     "1.0",
			Path:              "2,",
			Status:            order.StatusFilled,
			CreatedAt:         time.Now().Format("2006-01-02 15:04:05"),
			UserID:            2,
			UserAgentInfo:     "",
		},
		{
			OrderID:           3,
			Price:             "50000",
			DemandedAmount:    "50000",
			OrderType:         "SELL",
			OrderExchangeType: "LIMIT",
			PayedByAmount:     "1.0",
			Path:              "2,",
			Status:            order.StatusFilled,
			CreatedAt:         time.Now().Format("2006-01-02 15:04:05"),
			UserID:            3,
			UserAgentInfo:     "",
		},
	}
	orderRepo := new(mocks.OrderRepository)
	orderRepo.On("GetOrdersDataByIdsWithJoinUsingTx", mock.Anything, []int64{1, 2, 3}).Once().Return(orderFields)

	userBalanceService := new(mocks.UserBalanceService)
	forceTrader := new(mocks.ForceTrader)
	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetPrice", mock.Anything, "BTC-USDT").Once().Return("50000", nil)
	tradeEventsHandler := new(mocks.TradeEventsHandler)
	mqttManager := new(mocks.MqttManager)
	redisClient := new(mocks.RedisClient)
	configs := new(mocks.Configs)
	configs.On("GetEnv").Once().Return(platform.EnvTest)
	configs.On("GetBool", "commitError").Once().Return(false)
	logger := new(mocks.Logger)
	logger.On("Error2", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	currencyService := new(mocks.CurrencyService)
	userService := new(mocks.UserService)
	usersData := []user.UsersDataForOrderMatching{
		{
			UserID:             1,
			UserEmail:          "user1@test.com",
			UserLevelID:        1,
			UserPrivateChannel: "",
		},
		{
			UserID:             2,
			UserEmail:          "user2@test.com",
			UserLevelID:        1,
			UserPrivateChannel: "",
		},
		{
			UserID:             3,
			UserEmail:          "user3@test.com",
			UserLevelID:        1,
			UserPrivateChannel: "",
		},
	}
	userService.On("GetUsersDataForOrderMatching", []int{1, 2, 3}).Once().Return(usersData)
	userLevelService := new(mocks.UserLevelService)
	userLevels := []user.Level{
		{
			ID:                 1,
			MakerFeePercentage: 0.5,
			TakerFeePercentage: 1,
		},
	}
	userLevelService.On("GetLevelsByIds", []int64{1}).Once().Return(userLevels)
	poms := order.NewPostOrderMatchingService(db, orderRepo, userBalanceService, forceTrader, priceGenerator, tradeEventsHandler, mqttManager, redisClient, currencyService, userService, userLevelService, configs, logger)
	doneOrdersData := []order.CallBackOrderData{
		{
			ID:                1,
			PairName:          "BTC-USDT",
			OrderType:         "BUY",
			Quantity:          "2.5",
			Price:             "50000",
			Timestamp:         time.Now().Unix(),
			TradedWithOrderID: 2,
			QuantityTraded:    "1.0",
			TradePrice:        "50000",
			MinThresholdPrice: "49500",
			MaxThresholdPrice: "50500",
		},
		{
			ID:                1,
			PairName:          "BTC-USDT",
			OrderType:         "BUY",
			Quantity:          "1.0",
			Price:             "50000",
			Timestamp:         time.Now().Unix(),
			TradedWithOrderID: 3,
			QuantityTraded:    "1.0",
			TradePrice:        "50000",
			MinThresholdPrice: "49500",
			MaxThresholdPrice: "50500",
		},
		{
			ID:                2,
			PairName:          "BTC-USDT",
			OrderType:         "SELL",
			Quantity:          "1.0",
			Price:             "50000",
			Timestamp:         time.Now().Unix(),
			TradedWithOrderID: 1,
			QuantityTraded:    "1.0",
			TradePrice:        "50000",
			MinThresholdPrice: "49500",
			MaxThresholdPrice: "50500",
		},
		{
			ID:                3,
			PairName:          "BTC-USDT",
			OrderType:         "SELL",
			Quantity:          "1.0",
			Price:             "50000",
			Timestamp:         time.Now().Unix(),
			TradedWithOrderID: 1,
			QuantityTraded:    "1.0",
			TradePrice:        "50000",
			MinThresholdPrice: "49500",
			MaxThresholdPrice: "50500",
		},
	}

	partial := order.CallBackOrderData{
		ID:                1,
		PairName:          "BTC-USDT",
		OrderType:         "BUY",
		Quantity:          "0.5",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: 0,
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	matchingResult := poms.HandlePostOrderMatching(doneOrdersData, &partial, false)
	assert.NotNil(t, matchingResult.Err)
	assert.Nil(t, matchingResult.RemainingPartialOrder)
	assert.Equal(t, int64(1), matchingResult.RemovingDoneOrderIds[0])
	assert.Equal(t, int64(2), matchingResult.RemovingDoneOrderIds[1])
	assert.Equal(t, int64(3), matchingResult.RemovingDoneOrderIds[2])
	orderRepo.AssertExpectations(t)
	priceGenerator.AssertExpectations(t)
	logger.AssertExpectations(t)
	userService.AssertExpectations(t)
	userLevelService.AssertExpectations(t)
}

func TestPostOrderMatchingService_Error_Error_PartialIsMarket(t *testing.T) {
	qm := postOrderMatchingQueryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	//these are for updating orders themselves
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(1, 1))

	//2 userbalance update for every order demanded and payedBy coin
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))

	//3 transaction for every order demanded ,payedBy and fee,this scenario the main order would have 2 child so we have 5 orders
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(1, 1))

	//one trade for partial which would be traded with bot
	dbMock.ExpectExec("INSERT INTO trades").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	orderFields := []order.MatchingNeededQueryFields{
		{
			OrderID:           1,
			Price:             "",
			DemandedAmount:    "2.5",
			OrderType:         "BUY",
			OrderExchangeType: "MARKET",
			PayedByAmount:     "125000",
			Path:              "1,",
			Status:            order.StatusOpen,
			CreatedAt:         time.Now().Format("2006-01-02 15:04:05"),
			UserID:            1,
			UserAgentInfo:     "",
		},
	}
	orderRepo := new(mocks.OrderRepository)
	orderRepo.On("GetOrdersDataByIdsWithJoinUsingTx", mock.Anything, []int64{1}).Once().Return(orderFields)

	basisCoin := currency.Coin{
		ID:   1,
		Name: "USDT",
	}

	dependentCoin := currency.Coin{
		ID:   2,
		Name: "BTC",
	}

	userBalanceService := new(mocks.UserBalanceService)
	ubs := []userbalance.UserBalance{
		{
			ID:           1,
			UserID:       1,
			CoinID:       1,
			Coin:         basisCoin,
			Amount:       "150000",
			FrozenAmount: "100000",
		},
		{
			ID:           2,
			UserID:       1,
			CoinID:       2,
			Coin:         dependentCoin,
			Amount:       "1",
			FrozenAmount: "0",
		},
	}
	userBalanceService.On("GetBalancesOfUsersForCoinsUsingTx", mock.Anything, []int{1}, []int64{1, 2}).Once().Return(ubs)

	forceTrader := new(mocks.ForceTrader)

	forceTrader.On("ShouldForceTrade", "BTC-USDT", "BUY", "").Once().Return(true, nil)
	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetPrice", mock.Anything, "BTC-USDT").Once().Return("50000", nil)

	tradeEventsHandler := new(mocks.TradeEventsHandler)

	mqttManager := new(mocks.MqttManager)

	redisClient := new(mocks.RedisClient)
	redisClient.On("LPush", mock.Anything, order.UnmatchedOrdersList, "1").Once().Return(int64(1), nil)
	configs := new(mocks.Configs)
	configs.On("GetEnv").Once().Return(platform.EnvTest)
	configs.On("GetBool", "commitError").Once().Return(true)
	logger := new(mocks.Logger)
	currencyService := new(mocks.CurrencyService)
	pair := currency.Pair{
		ID:              1,
		Name:            "BTC-USDT",
		BasisCoinID:     1,
		BasisCoin:       basisCoin,
		DependentCoinID: 2,
		DependentCoin:   dependentCoin,
		MakerFee:        0.3,
		TakerFee:        0.3,
	}
	currencyService.On("GetPairByName", "BTC-USDT").Once().Return(pair, nil)
	userService := new(mocks.UserService)
	usersData := []user.UsersDataForOrderMatching{
		{
			UserID:             1,
			UserEmail:          "user1@test.com",
			UserLevelID:        1,
			UserPrivateChannel: "",
		},
	}
	userService.On("GetUsersDataForOrderMatching", []int{1}).Once().Return(usersData)
	userLevelService := new(mocks.UserLevelService)
	userLevels := []user.Level{
		{
			ID:                 1,
			MakerFeePercentage: 0.5,
			TakerFeePercentage: 1,
		},
	}
	userLevelService.On("GetLevelsByIds", []int64{1}).Once().Return(userLevels)
	poms := order.NewPostOrderMatchingService(db, orderRepo, userBalanceService, forceTrader, priceGenerator, tradeEventsHandler, mqttManager, redisClient, currencyService, userService, userLevelService, configs, logger)
	var doneOrdersData []order.CallBackOrderData

	partial := order.CallBackOrderData{
		ID:                1,
		PairName:          "BTC-USDT",
		OrderType:         "BUY",
		Quantity:          "2.5",
		Price:             "",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: 0,
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	matchingResult := poms.HandlePostOrderMatching(doneOrdersData, &partial, false)
	time.Sleep(100 * time.Millisecond) //to be sure gorutine would be done
	assert.NotNil(t, matchingResult.Err)
	assert.Nil(t, matchingResult.RemainingPartialOrder)
	currencyService.AssertExpectations(t)
	orderRepo.AssertExpectations(t)
	redisClient.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	priceGenerator.AssertExpectations(t)
	forceTrader.AssertExpectations(t)
	userService.AssertExpectations(t)
	userLevelService.AssertExpectations(t)
}
