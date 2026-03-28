// Package user_test tests the user permission manager. Covers:
//   - IsPermissionGrantedToUserFor: verifying granted (WITHDRAW, EXCHANGE) and denied (DEPOSIT) permissions
//   - GetAllPermissions: retrieving the full list of available permissions
//
// Test data: testify mocks for UsersPermissionsRepository and
// UserPermissionRepository with sample EXCHANGE, WITHDRAW, and DEPOSIT
// permission entries.
package user_test

import (
	"exchange-go/internal/mocks"
	"exchange-go/internal/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermissionManager_IsPermissionGrantedToUserFor(t *testing.T) {
	usersPermissionsRepo := new(mocks.UsersPermissionsRepository)
	usersPermissions := []user.UsersPermissions{
		{
			UserID:           1,
			UserPermissionID: 1,
			Permission:       user.Permission{Name: user.PermissionWithdraw},
		},
		{
			UserID:           1,
			UserPermissionID: 2,
			Permission:       user.Permission{Name: user.PermissionExchange},
		},
	}
	usersPermissionsRepo.On("GetUserPermissions", 1).Times(3).Return(usersPermissions)

	permissionsRepo := new(mocks.UserPermissionRepository)

	permissionManager := user.NewUserPermissionManager(usersPermissionsRepo, permissionsRepo)

	u := user.User{
		ID: 1,
	}
	isGranted := permissionManager.IsPermissionGrantedToUserFor(u, user.PermissionWithdraw)
	assert.Equal(t, true, isGranted)

	isGranted = permissionManager.IsPermissionGrantedToUserFor(u, user.PermissionDeposit)
	assert.Equal(t, false, isGranted)

	isGranted = permissionManager.IsPermissionGrantedToUserFor(u, user.PermissionExchange)
	assert.Equal(t, true, isGranted)
	usersPermissionsRepo.AssertExpectations(t)
}

func TestPermissionManager_GetAllPermissions(t *testing.T) {
	usersPermissionsRepo := new(mocks.UsersPermissionsRepository)
	permissions := []user.Permission{
		{
			ID:   1,
			Name: user.PermissionExchange,
		},
		{
			ID:   2,
			Name: user.PermissionWithdraw,
		},
		{
			ID:   3,
			Name: user.PermissionDeposit,
		},
	}
	permissionsRepo := new(mocks.UserPermissionRepository)
	permissionsRepo.On("GetAllPermissions").Once().Return(permissions)
	permissionManager := user.NewUserPermissionManager(usersPermissionsRepo, permissionsRepo)

	result := permissionManager.GetAllPermissions()
	assert.Equal(t, 3, len(result))

	for _, p := range result {
		switch p.ID {
		case int64(1):
			assert.Equal(t, "EXCHANGE", p.Name)
		case int64(2):
			assert.Equal(t, "WITHDRAW", p.Name)
		case int64(3):
			assert.Equal(t, "DEPOSIT", p.Name)
		default:
			t.Fatal("we should not be in default case")
		}

	}

	usersPermissionsRepo.AssertExpectations(t)
}
