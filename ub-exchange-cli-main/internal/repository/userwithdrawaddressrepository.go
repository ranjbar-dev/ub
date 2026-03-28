package repository

import (
	"exchange-go/internal/userwithdrawaddress"

	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

type userWithdrawAddressRepository struct {
	db *gorm.DB
}

func (r *userWithdrawAddressRepository) GetUserWithdrawAddressesByCoinID(userID int, coinID int64) []userwithdrawaddress.UserWithdrawAddress {
	var withdrawAddresses []userwithdrawaddress.UserWithdrawAddress
	q := r.db.Joins("Coin").
		Where("user_withdraw_address.user_id =?", userID).
		Where("user_withdraw_address.is_deleted =? or user_withdraw_address.is_deleted is null", false).
		Where("user_withdraw_address.currency_id =?", coinID)
	q.Order("user_withdraw_address.id desc").Find(&withdrawAddresses)

	return withdrawAddresses
}

func (r *userWithdrawAddressRepository) GetUserWithdrawAddressByID(id int64, uwa *userwithdrawaddress.UserWithdrawAddress) error {
	return r.db.Where(userwithdrawaddress.UserWithdrawAddress{ID: id}).First(uwa).Error
}

func (r *userWithdrawAddressRepository) GetUserWithdrawAddressesByIds(userID int, ids []int64) []userwithdrawaddress.UserWithdrawAddress {
	var withdrawAddresses []userwithdrawaddress.UserWithdrawAddress
	r.db.Where("id IN ? AND  user_id = ? and (is_deleted IS NULL OR is_deleted=?)", ids, userID, false).Find(&withdrawAddresses)
	return withdrawAddresses
}

func (r *userWithdrawAddressRepository) GetUserWithdrawAddresses(filters userwithdrawaddress.GetWithdrawAddressesFilters) []userwithdrawaddress.UserWithdrawAddress {
	var withdrawAddresses []userwithdrawaddress.UserWithdrawAddress
	q := r.db.Joins("Coin").
		Where("user_withdraw_address.user_id =?", filters.UserID).
		Where("user_withdraw_address.is_deleted =? or user_withdraw_address.is_deleted is null", false)

	if filters.CoinID != 0 {
		q.Where("user_withdraw_address.currency_id =?", filters.CoinID)
	}

	if filters.Coin != "" {
		q.Where("Coin.code = ?", filters.Coin)
	}

	if filters.Label != "" {
		q.Where("user_withdraw_address.label like ?", "%"+filters.Label+"%")
	}

	if filters.Address != "" {
		q.Where("user_withdraw_address.address like ?", "%"+filters.Address+"%")
	}

	q.Order("user_withdraw_address.id desc").Offset(int(filters.Page) * filters.PageSize).Limit(filters.PageSize)

	q.Find(&withdrawAddresses)
	return withdrawAddresses
}

func (r *userWithdrawAddressRepository) Create(uwa *userwithdrawaddress.UserWithdrawAddress) error {
	return r.db.Omit(clause.Associations).Create(uwa).Error
}

func (r *userWithdrawAddressRepository) GetUserWithdrawAddressesByAddress(userID int, coinID int64, address string) []userwithdrawaddress.UserWithdrawAddress {
	var userWithdrawAddresses []userwithdrawaddress.UserWithdrawAddress
	r.db.Where(userwithdrawaddress.UserWithdrawAddress{UserID: userID, CoinID: coinID, Address: address}).Find(&userWithdrawAddresses)
	return userWithdrawAddresses
}

func NewUserWithdrawAddressRepository(db *gorm.DB) userwithdrawaddress.Repository {
	return &userWithdrawAddressRepository{db}
}
