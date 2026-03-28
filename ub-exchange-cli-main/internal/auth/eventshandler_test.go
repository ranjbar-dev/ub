// Package auth_test tests the authentication events handler. Covers:
//   - Logging user login attempts to the login history service
//   - Detecting IP address changes and sending notification emails
//
// Test data: testify mocks for LoginHistoryService, CommunicationService,
// UserService, Configs, and Logger with sample user profiles, IP addresses,
// and email parameters.
package auth_test

import (
	"database/sql"
	"exchange-go/internal/auth"
	"exchange-go/internal/communication"
	"exchange-go/internal/mocks"
	"exchange-go/internal/user"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

func TestEventsHandler_LogUserLogin(t *testing.T) {
	loginHistoryService := new(mocks.LoginHistoryService)
	loginHistoryService.On("CreateLoginHistory", mock.Anything).Once().Return(nil)
	communicationService := new(mocks.CommunicationService)
	userService := new(mocks.UserService)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	eventsHandler := auth.NewAuthEventsHandler(loginHistoryService, communicationService, userService, configs, logger)
	params := auth.LogUserLoginParams{
		UserID:   1,
		Device:   "WEB",
		IP:       "127.0.0.1",
		Email:    "test@test.com",
		Type:     "SUCCESSFULL",
		Password: "",
	}
	eventsHandler.LogUserLogin(params)
	loginHistoryService.AssertExpectations(t)
}

func TestEventsHandler_NotifyIfUserIpHasChanged(t *testing.T) {
	loginHistoryService := new(mocks.LoginHistoryService)
	loginHistoryService.On("GetLastLoginHistoryByUserID", 1, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		lh := args.Get(1).(*user.LoginHistory)
		lh.ID = 1
		lh.IP = sql.NullString{String: "127.0.0.1", Valid: true}
	})

	communicationService := new(mocks.CommunicationService)
	emailParams := communication.SendIPChangedEmailParams{
		Email:       "test@test.com",
		FullName:    "test test",
		LastIP:      "127.0.0.1",
		CurrentIP:   "127.0.0.2",
		Device:      "WEB",
		ChangedDate: time.Now().Format("2006-01-02 15:04:05"),
	}
	communicationService.On("SendIPChangedEmail", mock.Anything, emailParams).Once().Return()

	up := user.Profile{
		FirstName: sql.NullString{String: "test", Valid: true},
		LastName:  sql.NullString{String: "test", Valid: true},
	}
	userService := new(mocks.UserService)
	userService.On("GetUserProfile", mock.Anything).Once().Return(up, nil)

	configs := new(mocks.Configs)
	configs.On("GetString", "wallet.username").Once().Return("wallet@wallet.com")

	logger := new(mocks.Logger)

	eventsHandler := auth.NewAuthEventsHandler(loginHistoryService, communicationService, userService, configs, logger)
	u := user.User{
		ID:     1,
		Email:  "test@test.com",
		Status: user.StatusVerified,
	}
	params := auth.NotifyIPChangeParams{
		User:   u,
		Device: "WEB",
		IP:     "127.0.0.2",
	}
	eventsHandler.NotifyIfUserIPHasChanged(params)

	loginHistoryService.AssertExpectations(t)
	communicationService.AssertExpectations(t)
	userService.AssertExpectations(t)
	configs.AssertExpectations(t)
}
