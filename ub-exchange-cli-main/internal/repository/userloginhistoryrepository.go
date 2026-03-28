package repository

import (
	"exchange-go/internal/user"

	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

type userLoginHistoryRepository struct {
	db *gorm.DB
}

func (r *userLoginHistoryRepository) Create(loginHistory *user.LoginHistory) error {
	return r.db.Omit(clause.Associations).Create(loginHistory).Error
}

func (r *userLoginHistoryRepository) GetLastLoginHistoryByUserID(userID int, loginHistory *user.LoginHistory) error {
	return r.db.Where("user_id = ?", userID).Order("id desc").First(loginHistory).Error
}

func NewUserLoginHistoryRepository(db *gorm.DB) user.LoginHistoryRepository {
	return &userLoginHistoryRepository{db}
}
