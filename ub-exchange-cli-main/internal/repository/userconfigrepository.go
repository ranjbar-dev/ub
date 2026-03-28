package repository

import (
	"exchange-go/internal/user"

	"gorm.io/gorm"
)

type userConfigRepository struct {
	db *gorm.DB
}

func (ucr *userConfigRepository) GetUserConfigByUserID(userID int, config *user.Config) error {
	return ucr.db.Where("user_id = ?", userID).First(&config).Error
}

func NewUserConfigRepository(db *gorm.DB) user.ConfigRepository {
	return &userConfigRepository{db}
}
