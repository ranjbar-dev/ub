// Package userwithdrawaddress_test tests the user withdraw address service. Covers:
//   - Listing all withdraw addresses for a user with coin and network details
//   - Creating new addresses with validation (coin not found, network not found, invalid address)
//   - Successful address creation for single-network coins
//   - Adding addresses to favorites (not found and success scenarios)
//   - Deleting multiple withdraw addresses by ID
//   - Retrieving former addresses filtered by coin code (with coin-not-found error case)
//   - Looking up withdraw addresses by address string
//   - Saving new address records through the repository
//
// Test data: mock repositories, currency service, wallet service, go-sqlmock for GORM,
// and pre-built withdraw address fixtures with ETH/BTC/USDT coins and networks.
package userwithdrawaddress_test

import (
	"database/sql"
	"exchange-go/internal/currency"
	"exchange-go/internal/mocks"
	"exchange-go/internal/user"
	"exchange-go/internal/userwithdrawaddress"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestService_GetWithdrawAddresses(t *testing.T) {
	db := &gorm.DB{}
	repo := new(mocks.UserWithdrawAddressRepository)
	withdrawAddresses := []userwithdrawaddress.UserWithdrawAddress{
		{
			ID:         1,
			Coin:       currency.Coin{Code: "ETH", Name: "Ethereum"},
			Address:    "ethAddress1",
			Label:      sql.NullString{String: "eth1", Valid: true},
			IsFavorite: sql.NullBool{Bool: true, Valid: true},
		},
		{
			ID:      2,
			Coin:    currency.Coin{Code: "BTC", Name: "Bitcoin"},
			Address: "btcAddress1",
			Label:   sql.NullString{String: "btc1", Valid: true},
		},
		{
			ID:      1,
			Coin:    currency.Coin{Code: "USDT", Name: "Tether"},
			Address: "usdtAddress1",
			Label:   sql.NullString{String: "usdt1", Valid: true},
			Network: sql.NullString{String: "ETH", Valid: true},
		},
	}
	repo.On("GetUserWithdrawAddresses", mock.Anything).Once().Return(withdrawAddresses)
	cs := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	logger := new(mocks.Logger)

	s := userwithdrawaddress.NewUserWithdrawAddressService(db, repo, cs, walletService, logger)
	u := &user.User{}
	params := userwithdrawaddress.GetWithdrawAddressesParams{}
	res, statusCode := s.GetWithdrawAddresses(u, params)
	assert.Equal(t, http.StatusOK, statusCode)

	resp, ok := res.Data.([]userwithdrawaddress.GetWithdrawAddressesResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}
	address1 := resp[0]
	assert.Equal(t, "ethAddress1", address1.Address)
	assert.Equal(t, "eth1", address1.Label)
	assert.Equal(t, true, address1.IsFavorite)
	assert.Equal(t, "ETH", address1.Coin)
	assert.Equal(t, "Ethereum", address1.Name)
	assert.Equal(t, "", address1.Network)

	address2 := resp[1]
	assert.Equal(t, "btcAddress1", address2.Address)
	assert.Equal(t, "btc1", address2.Label)
	assert.Equal(t, false, address2.IsFavorite)
	assert.Equal(t, "BTC", address2.Coin)
	assert.Equal(t, "Bitcoin", address2.Name)
	assert.Equal(t, "", address2.Network)

	address3 := resp[2]
	assert.Equal(t, "usdtAddress1", address3.Address)
	assert.Equal(t, "usdt1", address3.Label)
	assert.Equal(t, false, address3.IsFavorite)
	assert.Equal(t, "USDT", address3.Coin)
	assert.Equal(t, "Tether", address3.Name)
	assert.Equal(t, "ETH", address3.Network)

	repo.AssertExpectations(t)

}

func TestService_CreateNewAddress_Fail_CoinNotFound(t *testing.T) {
	db := &gorm.DB{}
	repo := new(mocks.UserWithdrawAddressRepository)
	cs := new(mocks.CurrencyService)
	cs.On("GetCoinByCode", "BTCQ").Once().Return(currency.Coin{}, gorm.ErrRecordNotFound)
	walletService := new(mocks.WalletService)
	logger := new(mocks.Logger)
	s := userwithdrawaddress.NewUserWithdrawAddressService(db, repo, cs, walletService, logger)

	u := &user.User{}
	params := userwithdrawaddress.CreateAddressParams{
		Coin:    "BTCq",
		Label:   "label",
		Address: "address",
		Network: "network",
	}
	res, statusCode := s.CreateNewAddress(u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "coin not found", res.Message)

	cs.AssertExpectations(t)
}

func TestService_CreateNewAddress_Fail_NetworkCoinNotFound(t *testing.T) {
	db := &gorm.DB{}
	repo := new(mocks.UserWithdrawAddressRepository)
	cs := new(mocks.CurrencyService)
	cs.On("GetCoinByCode", "USDT").Once().Return(currency.Coin{ID: 1}, nil)
	cs.On("GetCoinByCode", "ETHQ").Once().Return(currency.Coin{}, gorm.ErrRecordNotFound)
	walletService := new(mocks.WalletService)
	logger := new(mocks.Logger)
	s := userwithdrawaddress.NewUserWithdrawAddressService(db, repo, cs, walletService, logger)

	u := &user.User{}
	params := userwithdrawaddress.CreateAddressParams{
		Coin:    "USDT",
		Label:   "label",
		Address: "address",
		Network: "ethq",
	}
	res, statusCode := s.CreateNewAddress(u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "network not found", res.Message)

	cs.AssertExpectations(t)
}

func TestService_CreateNewAddress_Fail_AddressIsNotValid(t *testing.T) {
	db := &gorm.DB{}
	repo := new(mocks.UserWithdrawAddressRepository)
	cs := new(mocks.CurrencyService)
	cs.On("GetCoinByCode", "USDT").Once().Return(currency.Coin{ID: 1, Code: "USDT"}, nil)
	cs.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 2}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(false, nil)
	logger := new(mocks.Logger)
	s := userwithdrawaddress.NewUserWithdrawAddressService(db, repo, cs, walletService, logger)

	u := &user.User{}
	params := userwithdrawaddress.CreateAddressParams{
		Coin:    "USDT",
		Label:   "label",
		Address: "ethAddress1",
		Network: "eth",
	}
	res, statusCode := s.CreateNewAddress(u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "address is not valid", res.Message)

	cs.AssertExpectations(t)
	walletService.AssertExpectations(t)
}

func TestService_CreateNewAddress_Successful(t *testing.T) {
	db := &gorm.DB{}
	repo := new(mocks.UserWithdrawAddressRepository)
	repo.On("Create", mock.Anything).Once().Return(nil)
	cs := new(mocks.CurrencyService)
	cs.On("GetCoinByCode", "BTC").Once().Return(currency.Coin{ID: 1, Code: "BTC", Name: "Bitcoin"}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "BTC", "btcAddress1", "").Once().Return(true, nil)
	logger := new(mocks.Logger)

	s := userwithdrawaddress.NewUserWithdrawAddressService(db, repo, cs, walletService, logger)
	u := &user.User{}
	params := userwithdrawaddress.CreateAddressParams{
		Coin:    "BTC",
		Label:   "label",
		Address: "btcAddress1",
		Network: "",
	}

	res, statusCode := s.CreateNewAddress(u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	resp, ok := res.Data.([]userwithdrawaddress.CreateAddressResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}
	assert.Equal(t, "btcAddress1", resp[0].Address)
	assert.Equal(t, "label", resp[0].Label)
	assert.Equal(t, false, resp[0].IsFavorite)
	assert.Equal(t, "BTC", resp[0].Coin)
	assert.Equal(t, "Bitcoin", resp[0].Name)
	assert.Equal(t, "", resp[0].Network)
	repo.AssertExpectations(t)
}

func TestService_AddToFavorites_NotFound(t *testing.T) {
	db := &gorm.DB{}
	repo := new(mocks.UserWithdrawAddressRepository)
	repo.On("GetUserWithdrawAddressByID", int64(1), mock.Anything).Once().Return(gorm.ErrRecordNotFound)
	cs := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	logger := new(mocks.Logger)
	s := userwithdrawaddress.NewUserWithdrawAddressService(db, repo, cs, walletService, logger)

	u := &user.User{}
	params := userwithdrawaddress.AddToFavoritesParams{
		ID:     1,
		Action: "add",
	}
	res, statusCode := s.AddToFavorites(u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "address not found", res.Message)
	repo.AssertExpectations(t)
}

// queryMatcher is a sqlmock.QueryMatcher that accepts any SQL query,
// allowing tests to focus on service logic rather than exact SQL statements.
type queryMatcher struct {
}

// Match satisfies the sqlmock.QueryMatcher interface by unconditionally matching
// any expected SQL against any actual SQL.
func (queryMatcher) Match(expectedSQL, actualSQL string) error {
	return nil
}

func TestService_AddToFavorites_Successful(t *testing.T) {
	qm := queryMatcher{}
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
	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectExec("UPDATE user_withdraw_address").WillReturnResult(sqlmock.NewResult(1, 1))

	repo := new(mocks.UserWithdrawAddressRepository)
	repo.On("GetUserWithdrawAddressByID", int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		uwa := args.Get(1).(*userwithdrawaddress.UserWithdrawAddress)
		uwa.ID = 1
		uwa.UserID = 1
	})
	cs := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	logger := new(mocks.Logger)
	s := userwithdrawaddress.NewUserWithdrawAddressService(db, repo, cs, walletService, logger)

	u := &user.User{
		ID: 1,
	}
	params := userwithdrawaddress.AddToFavoritesParams{
		ID:     1,
		Action: "add",
	}
	res, statusCode := s.AddToFavorites(u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)

	repo.AssertExpectations(t)

}

func TestService_Delete(t *testing.T) {
	qm := queryMatcher{}
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
	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectExec("UPDATE user_withdraw_address").WillReturnResult(sqlmock.NewResult(1, 1))

	withdrawAddresses := []userwithdrawaddress.UserWithdrawAddress{
		{
			ID:         1,
			Coin:       currency.Coin{Code: "ETH", Name: "Ethereum"},
			Address:    "ethAddress1",
			Label:      sql.NullString{String: "eth1", Valid: true},
			IsFavorite: sql.NullBool{Bool: true, Valid: true},
		},
		{
			ID:      2,
			Coin:    currency.Coin{Code: "BTC", Name: "Bitcoin"},
			Address: "btcAddress1",
			Label:   sql.NullString{String: "btc1", Valid: true},
		},
		{
			ID:      1,
			Coin:    currency.Coin{Code: "USDT", Name: "Tether"},
			Address: "usdtAddress1",
			Label:   sql.NullString{String: "usdt1", Valid: true},
			Network: sql.NullString{String: "ETH", Valid: true},
		},
	}
	repo := new(mocks.UserWithdrawAddressRepository)
	repo.On("GetUserWithdrawAddressesByIds", 1, []int64{1, 2, 3}).Once().Return(withdrawAddresses)
	cs := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	logger := new(mocks.Logger)
	s := userwithdrawaddress.NewUserWithdrawAddressService(db, repo, cs, walletService, logger)

	u := &user.User{
		ID: 1,
	}
	params := userwithdrawaddress.DeleteParams{
		Ids: []int64{1, 2, 3},
	}
	res, statusCode := s.Delete(u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	repo.AssertExpectations(t)

}

func TestService_GetFormerAddresses_CoinNotFound(t *testing.T) {
	db := &gorm.DB{}
	repo := new(mocks.UserWithdrawAddressRepository)
	cs := new(mocks.CurrencyService)
	cs.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{}, gorm.ErrRecordNotFound)
	walletService := new(mocks.WalletService)
	logger := new(mocks.Logger)

	s := userwithdrawaddress.NewUserWithdrawAddressService(db, repo, cs, walletService, logger)
	u := &user.User{
		ID: 1,
	}
	params := userwithdrawaddress.GetFormerAddressesParams{
		Coin: "ETH",
	}
	res, statusCode := s.GetFormerAddresses(u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "coin not found", res.Message)

	cs.AssertExpectations(t)
}

func TestService_GetFormerAddresses(t *testing.T) {
	db := &gorm.DB{}
	repo := new(mocks.UserWithdrawAddressRepository)
	withdrawAddresses := []userwithdrawaddress.UserWithdrawAddress{
		{
			ID:         1,
			Coin:       currency.Coin{Code: "ETH", Name: "Ethereum"},
			Address:    "ethAddress1",
			Label:      sql.NullString{String: "eth1", Valid: true},
			IsFavorite: sql.NullBool{Bool: true, Valid: true},
		},
		{
			ID:      1,
			Coin:    currency.Coin{Code: "ETH", Name: "Ethereum"},
			Address: "ethAddress2",
			Label:   sql.NullString{String: "eth2", Valid: true},
		},
		{
			ID:      1,
			Coin:    currency.Coin{Code: "ETH", Name: "Ethereum"},
			Address: "ethAddress3",
			Label:   sql.NullString{String: "eth3", Valid: true},
		},
	}
	repo.On("GetUserWithdrawAddressesByCoinID", 1, int64(1)).Once().Return(withdrawAddresses)
	cs := new(mocks.CurrencyService)
	cs.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 1, Code: "ETH", Name: "Ethereum"}, nil)
	walletService := new(mocks.WalletService)
	logger := new(mocks.Logger)

	s := userwithdrawaddress.NewUserWithdrawAddressService(db, repo, cs, walletService, logger)
	u := &user.User{
		ID: 1,
	}
	params := userwithdrawaddress.GetFormerAddressesParams{
		Coin: "ETH",
	}
	res, statusCode := s.GetFormerAddresses(u, params)
	assert.Equal(t, http.StatusOK, statusCode)

	resp, ok := res.Data.([]userwithdrawaddress.GetWithdrawAddressesResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}
	address1 := resp[0]
	assert.Equal(t, "ethAddress1", address1.Address)
	assert.Equal(t, "eth1", address1.Label)
	assert.Equal(t, true, address1.IsFavorite)
	assert.Equal(t, "ETH", address1.Coin)
	assert.Equal(t, "Ethereum", address1.Name)
	assert.Equal(t, "", address1.Network)

	address2 := resp[1]
	assert.Equal(t, "ethAddress2", address2.Address)
	assert.Equal(t, "eth2", address2.Label)
	assert.Equal(t, false, address2.IsFavorite)
	assert.Equal(t, "ETH", address2.Coin)
	assert.Equal(t, "Ethereum", address2.Name)
	assert.Equal(t, "", address2.Network)

	address3 := resp[2]

	assert.Equal(t, "ethAddress3", address3.Address)
	assert.Equal(t, "eth3", address3.Label)
	assert.Equal(t, false, address3.IsFavorite)
	assert.Equal(t, "ETH", address3.Coin)
	assert.Equal(t, "Ethereum", address2.Name)
	assert.Equal(t, "", address2.Network)
	repo.AssertExpectations(t)
}

func TestService_GetUserWithdrawAddressesByAddress(t *testing.T) {
	db := &gorm.DB{}
	repo := new(mocks.UserWithdrawAddressRepository)
	withdrawAddresses := []userwithdrawaddress.UserWithdrawAddress{
		{
			ID:         1,
			Coin:       currency.Coin{Code: "ETH", Name: "Ethereum"},
			Address:    "ethAddress1",
			Label:      sql.NullString{String: "eth1", Valid: true},
			IsFavorite: sql.NullBool{Bool: true, Valid: true},
		},
	}
	repo.On("GetUserWithdrawAddressesByAddress", 1, int64(1), "btcAddress1").Once().Return(withdrawAddresses)
	cs := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	logger := new(mocks.Logger)

	s := userwithdrawaddress.NewUserWithdrawAddressService(db, repo, cs, walletService, logger)
	u := &user.User{
		ID: 1,
	}
	coin := currency.Coin{
		ID: 1,
	}
	address := "btcAddress1"
	withdrawAddressesResult := s.GetUserWithdrawAddressesByAddress(u, coin, address)

	address1 := withdrawAddressesResult[0]
	assert.Equal(t, "ethAddress1", address1.Address)
	assert.Equal(t, sql.NullString(sql.NullString{String: "eth1", Valid: true}), address1.Label)
	assert.Equal(t, sql.NullBool{Bool: true, Valid: true}, address1.IsFavorite)
	assert.Equal(t, sql.NullString{String: "", Valid: false}, address1.Network)
	repo.AssertExpectations(t)
}

func TestService_SaveNewAddress(t *testing.T) {
	db := &gorm.DB{}
	repo := new(mocks.UserWithdrawAddressRepository)
	repo.On("Create", mock.Anything).Once().Return(nil)
	cs := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	logger := new(mocks.Logger)
	s := userwithdrawaddress.NewUserWithdrawAddressService(db, repo, cs, walletService, logger)
	u := &user.User{}
	coin := currency.Coin{}
	params := userwithdrawaddress.CreateAddressParams{
		Coin:    "BTC",
		Label:   "label",
		Address: "btcAddress1",
		Network: "",
	}

	uwa, err := s.SaveNewAddress(u, coin, params)
	assert.Nil(t, err)

	assert.Equal(t, sql.NullString{String: "label", Valid: true}, uwa.Label)
	assert.Equal(t, sql.NullBool{Bool: false, Valid: true}, uwa.IsFavorite)
	assert.Equal(t, "btcAddress1", uwa.Address)
	repo.AssertExpectations(t)
}
