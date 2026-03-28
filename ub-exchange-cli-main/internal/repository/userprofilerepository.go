package repository

import (
	"exchange-go/internal/user"

	"gorm.io/gorm"
)

type userProfileRepository struct {
	db *gorm.DB
}

func (ur *userProfileRepository) GetProfileByUserID(userID int, profile *user.Profile) error {
	return ur.db.Where(&user.Profile{UserID: userID}).First(profile).Error
}

func (ur *userProfileRepository) GetProfileByUserIDUsingTx(tx *gorm.DB, userID int, profile *user.Profile) error {
	return tx.Where(&user.Profile{UserID: userID}).First(profile).Error
}

func NewUserProfileRepository(db *gorm.DB) user.ProfileRepository {
	return &userProfileRepository{db}
}
