package repository

import (
	"exchange-go/internal/user"
	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func (r *permissionRepository) GetAllPermissions() []user.Permission {
	var permissions []user.Permission
	r.db.Find(&permissions)
	return permissions
}

func NewPermissionRepository(db *gorm.DB) user.PermissionRepository {
	return &permissionRepository{db}
}
