package user

import "strings"

const (
	PermissionDeposit      = "DEPOSIT"
	PermissionWithdraw     = "WITHDRAW"
	PermissionExchange     = "EXCHANGE"
	PermissionFiatDeposit  = "FIAT_DEPOSIT"
	PermissionFiatWithdraw = "FIAT_WITHDRAW"
)

// PermissionManager provides methods for checking and listing user permissions
// such as DEPOSIT, WITHDRAW, EXCHANGE, FIAT_DEPOSIT, and FIAT_WITHDRAW.
type PermissionManager interface {
	// IsPermissionGrantedToUserFor checks whether the given user has been granted
	// the specified permission (case-insensitive comparison).
	IsPermissionGrantedToUserFor(user User, permissionName string) bool
	// GetAllPermissions retrieves all defined permissions in the system.
	GetAllPermissions() []Permission
}

type permissionManager struct {
	usersPermissionsRepo UsersPermissionsRepository
	permissionRepository PermissionRepository
}

func (pm *permissionManager) GetAllPermissions() []Permission {
	return pm.permissionRepository.GetAllPermissions()
}

func (pm *permissionManager) IsPermissionGrantedToUserFor(user User, permissionName string) bool {
	usersPermissions := pm.usersPermissionsRepo.GetUserPermissions(user.ID)
	for _, up := range usersPermissions {
		if strings.ToUpper(up.Permission.Name) == strings.ToUpper(permissionName) {
			return true
		}
	}
	return false
}

func NewUserPermissionManager(usersPermissionsRepo UsersPermissionsRepository, permissionRepository PermissionRepository) PermissionManager {
	return &permissionManager{
		usersPermissionsRepo: usersPermissionsRepo,
		permissionRepository: permissionRepository,
	}

}
