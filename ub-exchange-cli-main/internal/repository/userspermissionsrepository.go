package repository

import (
	"exchange-go/internal/user"

	"gorm.io/gorm"
)

type usersPermissionsRepository struct {
	db *gorm.DB
}

func (r *usersPermissionsRepository) GetUserPermissions(userID int) []user.UsersPermissions {
	var ups []user.UsersPermissions
	r.db.Joins("Permission").Where("user_id = ?", userID).Find(&ups)
	return ups
}

func NewUsersPermissionsRepository(db *gorm.DB) user.UsersPermissionsRepository {
	return &usersPermissionsRepository{db}
}
