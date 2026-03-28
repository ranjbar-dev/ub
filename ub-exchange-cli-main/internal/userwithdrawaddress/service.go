package userwithdrawaddress

import (
	"database/sql"
	"errors"
	"exchange-go/internal/currency"
	"exchange-go/internal/platform"
	"exchange-go/internal/response"
	"exchange-go/internal/user"
	"exchange-go/internal/wallet"
	"net/http"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Service provides the public API for managing user withdrawal addresses, including
// CRUD operations, favorites, and former address lookups.
type Service interface {
	// GetWithdrawAddresses returns a paginated list of the user's withdrawal addresses.
	GetWithdrawAddresses(user *user.User, params GetWithdrawAddressesParams) (apiResponse response.APIResponse, statusCode int)
	// CreateNewAddress creates and validates a new withdrawal address for the user.
	CreateNewAddress(user *user.User, params CreateAddressParams) (apiResponse response.APIResponse, statusCode int)
	// AddToFavorites marks a withdrawal address as a favorite for quick access.
	AddToFavorites(user *user.User, params AddToFavoritesParams) (apiResponse response.APIResponse, statusCode int)
	// Delete soft-deletes a withdrawal address from the user's address book.
	Delete(user *user.User, params DeleteParams) (apiResponse response.APIResponse, statusCode int)
	// GetFormerAddresses returns addresses that were previously used for withdrawals.
	GetFormerAddresses(user *user.User, params GetFormerAddressesParams) (apiResponse response.APIResponse, statusCode int)
	// GetUserWithdrawAddressesByAddress looks up addresses matching the given user, coin, and address.
	GetUserWithdrawAddressesByAddress(user *user.User, coin currency.Coin, address string) []UserWithdrawAddress
	// SaveNewAddress creates and persists a new withdrawal address, returning the saved record.
	SaveNewAddress(user *user.User, coin currency.Coin, params CreateAddressParams) (UserWithdrawAddress, error)
}

type service struct {
	db              *gorm.DB
	repository      Repository
	currencyService currency.Service
	walletService   wallet.Service
	logger          platform.Logger
}

type GetWithdrawAddressesParams struct {
	Label    string `form:"label"`
	CoinID   int64  `form:"currency_id"`
	Coin     string `form:"code"`
	Address  string `form:"address"`
	Page     int64  `form:"page"`
	PageSize int    `form:"page_size"`
}

type GetWithdrawAddressesFilters struct {
	UserID   int
	Label    string
	Address  string
	CoinID   int64
	Coin     string
	Page     int64
	PageSize int
}

type GetWithdrawAddressesResponse struct {
	ID         int64  `json:"id"`
	Address    string `json:"address"`
	Label      string `json:"label"`
	IsFavorite bool   `json:"isFavorite"`
	Coin       string `json:"code"`
	Name       string `json:"name"`
	Network    string `json:"network"`
}

type CreateAddressParams struct {
	Coin    string `json:"code" binding:"required"`
	Label   string `json:"label"  binding:"required"`
	Address string `json:"address" binding:"required"`
	Network string `json:"network"`
}

type CreateAddressResponse struct {
	ID         int64  `json:"id"`
	Address    string `json:"address"`
	Label      string `json:"label"`
	IsFavorite bool   `json:"isFavorite"`
	Coin       string `json:"code"`
	Name       string `json:"name"`
	Network    string `json:"network"`
}

type AddToFavoritesParams struct {
	ID     int64  `json:"id" binding:"required"`
	Action string `json:"action" binding:"required,oneof='add' 'remove'"`
}

type DeleteParams struct {
	Ids []int64 `json:"ids" binding:"required"`
}

type GetFormerAddressesParams struct {
	Coin string `form:"code"`
}

func (s *service) GetWithdrawAddresses(u *user.User, params GetWithdrawAddressesParams) (apiResponse response.APIResponse, statusCode int) {
	result := make([]GetWithdrawAddressesResponse, 0)
	filters := s.getFiltersForWithdrawAddresses(params)
	filters.UserID = u.ID
	//get from repo
	withdrawAddresses := s.repository.GetUserWithdrawAddresses(filters)
	for _, a := range withdrawAddresses {
		r := GetWithdrawAddressesResponse{
			ID:         a.ID,
			Address:    a.Address,
			Label:      a.Label.String,
			IsFavorite: a.IsFavorite.Bool,
			Coin:       a.Coin.Code,
			Name:       a.Coin.Name,
			Network:    a.Network.String,
		}
		result = append(result, r)
	}
	return response.Success(result, "")

}

func (s *service) getFiltersForWithdrawAddresses(params GetWithdrawAddressesParams) GetWithdrawAddressesFilters {
	var filters GetWithdrawAddressesFilters
	filters.Coin = strings.ToUpper(params.Coin)
	filters.CoinID = params.CoinID
	filters.Label = strings.Trim(params.Label, "")
	filters.Address = strings.Trim(params.Address, "")

	if params.Page >= 0 {
		filters.Page = params.Page
	}

	filters.PageSize = params.PageSize

	if filters.PageSize == 0 {
		filters.PageSize = 20
	}

	if filters.PageSize > 50 {
		filters.PageSize = 50
	}

	return filters

}

func (s *service) CreateNewAddress(user *user.User, params CreateAddressParams) (apiResponse response.APIResponse, statusCode int) {
	coinCode := strings.ToUpper(strings.Trim(params.Coin, ""))
	network := strings.ToUpper(strings.Trim(params.Network, ""))

	coin, err := s.currencyService.GetCoinByCode(coinCode)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("can not get coin by code", err,
			zap.String("service", "userWithdrawAddressService"),
			zap.String("method", "CreateNewAddress"),
			zap.String("coinCode", coinCode),
			zap.Int("userID", user.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || coin.ID == 0 {
		return response.Error("coin not found", http.StatusUnprocessableEntity, nil)
	}

	if network != "" {
		// check if network exists
		parentCoin, err := s.currencyService.GetCoinByCode(network)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error2("can not get coin by code for parent", err,
				zap.String("service", "userWithdrawAddressService"),
				zap.String("method", "CreateNewAddress"),
				zap.String("coinCode", network),
				zap.Int("userID", user.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}

		if errors.Is(err, gorm.ErrRecordNotFound) || parentCoin.ID == 0 {
			return response.Error("network not found", http.StatusUnprocessableEntity, nil)
		}
	}

	if network == "" && coin.BlockchainNetwork.Valid {
		network = strings.ToUpper(coin.BlockchainNetwork.String)
	}

	isValid, err := s.walletService.IsAddressValid(coin.Code, params.Address, network)
	if err != nil {
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if !isValid {
		return response.Error("address is not valid", http.StatusUnprocessableEntity, nil)
	}

	uwa, err := s.SaveNewAddress(user, coin, params)
	if err != nil {
		return response.Error("can not save now", http.StatusUnprocessableEntity, nil)
	}
	result := make([]CreateAddressResponse, 0)
	res := CreateAddressResponse{
		ID:         uwa.ID,
		Address:    uwa.Address,
		Label:      uwa.Label.String,
		IsFavorite: false,
		Coin:       coin.Code,
		Name:       coin.Name,
		Network:    params.Network,
	}
	result = append(result, res)
	return response.Success(result, "")
}

func (s *service) AddToFavorites(user *user.User, params AddToFavoritesParams) (apiResponse response.APIResponse, statusCode int) {
	id := params.ID
	action := params.Action
	uwa := &UserWithdrawAddress{}
	err := s.repository.GetUserWithdrawAddressByID(id, uwa)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("can not get userWithdrawAddress by id", err,
			zap.String("service", "userWithdrawAddressService"),
			zap.String("method", "AddToFavorites"),
			zap.Int64("userWithdrawAddressID", id),
			zap.Int("userID", user.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || uwa.ID == 0 || uwa.IsDeleted.Bool == true || uwa.UserID != user.ID {
		return response.Error("address not found", http.StatusUnprocessableEntity, nil)
	}

	if action == "add" {
		if !uwa.IsFavorite.Valid || uwa.IsFavorite.Bool != true {
			uwa.IsFavorite = sql.NullBool{Bool: true, Valid: true}
			err = s.db.Omit(clause.Associations).Save(uwa).Error
		}
	} else {
		if uwa.IsFavorite.Valid && uwa.IsFavorite.Bool == true {
			uwa.IsFavorite = sql.NullBool{Bool: false, Valid: true}
			err = s.db.Omit(clause.Associations).Save(uwa).Error
		}

	}

	if err != nil {
		s.logger.Error2("can not save user withdraw address", err,
			zap.String("service", "userWithdrawAddressService"),
			zap.String("method", "AddToFavorites"),
			zap.Int64("userWithdrawAddressID", id),
			zap.Int("userId", user.ID),
			zap.String("action", action),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	res := make(map[string]string, 0)
	return response.Success(res, "")
}

func (s *service) Delete(u *user.User, params DeleteParams) (apiResponse response.APIResponse, statusCode int) {
	res := make(map[string]string)
	ids := params.Ids
	if len(ids) < 1 {
		return response.Success(res, "")
	}

	withdrawAddresses := s.repository.GetUserWithdrawAddressesByIds(u.ID, ids)

	var shouldDeletedIds []int64
	for _, a := range withdrawAddresses {
		if !a.IsDeleted.Valid || a.IsDeleted.Bool == false {
			shouldDeletedIds = append(shouldDeletedIds, a.ID)
		}

	}
	err := s.db.Model(&UserWithdrawAddress{}).Where("id IN ?", shouldDeletedIds).Update("is_deleted", true).Error
	if err != nil {
		s.logger.Error2("can not delete user withdraw address", err,
			zap.String("service", "userWithdrawAddressService"),
			zap.String("method", "Delete"),
			zap.Int64s("userWithdrawAddressIDs", shouldDeletedIds),
			zap.Int("userID", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	return response.Success(res, "")
}

func (s *service) GetFormerAddresses(u *user.User, params GetFormerAddressesParams) (apiResponse response.APIResponse, statusCode int) {
	coinCode := strings.ToUpper(strings.Trim(params.Coin, ""))
	coin, err := s.currencyService.GetCoinByCode(coinCode)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("can not get coin by code", err,
			zap.String("service", "userWithdrawAddressService"),
			zap.String("method", "GetFormerAddresses"),
			zap.String("coinCode", coinCode),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || coin.ID == 0 {
		return response.Error("coin not found", http.StatusUnprocessableEntity, nil)
	}

	withdrawAddresses := s.repository.GetUserWithdrawAddressesByCoinID(u.ID, coin.ID)
	result := make([]GetWithdrawAddressesResponse, 0)
	for _, a := range withdrawAddresses {
		r := GetWithdrawAddressesResponse{
			ID:         a.ID,
			Address:    a.Address,
			Label:      a.Label.String,
			IsFavorite: a.IsFavorite.Bool,
			Coin:       a.Coin.Code,
			Name:       a.Coin.Name,
			Network:    a.Network.String,
		}
		result = append(result, r)
	}
	return response.Success(result, "")
}

func (s *service) SaveNewAddress(user *user.User, coin currency.Coin, params CreateAddressParams) (UserWithdrawAddress, error) {
	network := sql.NullString{String: "", Valid: false}
	if params.Network != "" {
		network = sql.NullString{String: params.Network, Valid: true}
	}
	uwa := UserWithdrawAddress{
		UserID:     user.ID,
		CoinID:     coin.ID,
		Address:    params.Address,
		Label:      sql.NullString{String: params.Label, Valid: true},
		IsDeleted:  sql.NullBool{Bool: false, Valid: true},
		IsFavorite: sql.NullBool{Bool: false, Valid: true},
		Network:    network,
	}
	err := s.repository.Create(&uwa)
	return uwa, err
}

func (s *service) GetUserWithdrawAddressesByAddress(user *user.User, coin currency.Coin, address string) []UserWithdrawAddress {
	return s.repository.GetUserWithdrawAddressesByAddress(user.ID, coin.ID, address)
}

func NewUserWithdrawAddressService(db *gorm.DB, repository Repository, currencyService currency.Service, walletService wallet.Service,
	logger platform.Logger) Service {
	return &service{
		db:              db,
		repository:      repository,
		currencyService: currencyService,
		walletService:   walletService,
		logger:          logger,
	}
}
