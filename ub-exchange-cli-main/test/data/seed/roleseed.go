package seed

import (
	"exchange-go/internal/user"

	"gorm.io/gorm"
)

func RoleSeed(db *gorm.DB) {
	superAdminRole := user.Role{
		ID:   1,
		Name: "super admin",
		Role: "ROLE_SUPER_ADMIN",
	}

	adminRole := user.Role{
		ID:   2,
		Name: "super admin",
		Role: "ROLE_ADMIN",
	}

	db.Create(&superAdminRole)
	db.Create(&adminRole)
}
