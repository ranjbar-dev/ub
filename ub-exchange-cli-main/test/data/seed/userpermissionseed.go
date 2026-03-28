package seed

import (
	"exchange-go/internal/user"
	"gorm.io/gorm"
)

func UserPermissionSeed(db *gorm.DB) {
	up1 := user.Permission{
		ID:   1,
		Name: user.PermissionDeposit,
	}

	up2 := user.Permission{
		ID:   2,
		Name: user.PermissionWithdraw,
	}

	up3 := user.Permission{
		ID:   3,
		Name: user.PermissionExchange,
	}

	up4 := user.Permission{
		ID:   4,
		Name: user.PermissionFiatDeposit,
	}

	up5 := user.Permission{
		ID:   5,
		Name: user.PermissionFiatWithdraw,
	}

	db.Create(&up1)
	db.Create(&up2)
	db.Create(&up3)
	db.Create(&up4)
	db.Create(&up5)

}
