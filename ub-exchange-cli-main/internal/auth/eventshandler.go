package auth

import (
	"database/sql"
	"errors"
	"exchange-go/internal/communication"
	"exchange-go/internal/platform"
	"exchange-go/internal/user"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LogUserLoginParams struct {
	UserID   int
	Device   string
	IP       string
	Email    string
	Type     string
	Password string
}

type NotifyIPChangeParams struct {
	User   user.User
	Device string
	IP     string
}

// EventsHandler tracks security-related authentication events such as login attempts
// and IP address changes.
type EventsHandler interface {
	// LogUserLogin records a login attempt with the associated device and IP information.
	LogUserLogin(params LogUserLoginParams)
	// NotifyIfUserIPHasChanged sends an email alert when a user logs in from a new IP address.
	NotifyIfUserIPHasChanged(params NotifyIPChangeParams)
}

type eventsHandler struct {
	loginHistoryService  user.LoginHistoryService
	communicationService communication.Service
	userService          user.Service
	configs              platform.Configs
	logger               platform.Logger
}

func (s *eventsHandler) LogUserLogin(params LogUserLoginParams) {
	userIDNil := sql.NullInt64{Int64: 0, Valid: false}
	if params.UserID != 0 {
		userIDNil = sql.NullInt64{Int64: int64(params.UserID), Valid: true}
	}

	lh := &user.LoginHistory{
		UserID:   userIDNil,
		Device:   sql.NullString{String: params.Device, Valid: true},
		IP:       sql.NullString{String: params.IP, Valid: true},
		Email:    sql.NullString{String: params.Email, Valid: true},
		Type:     params.Type,
		Password: sql.NullString{String: "", Valid: false}, //always should be empty because we do not want to store user password
	}
	err := s.loginHistoryService.CreateLoginHistory(lh)
	if err != nil {
		s.logger.Error2("error in creating loginHistory", err,
			zap.String("service", "authEventsHandler"),
			zap.String("method", "LogUserLogin"),
			zap.Int("userID", params.UserID),
		)
	}
}

func (s *eventsHandler) NotifyIfUserIPHasChanged(params NotifyIPChangeParams) {
	if params.User.Status != user.StatusVerified {
		return

	}
	if s.isWalletUser(params.User.Email) {
		return
	}

	lh := &user.LoginHistory{}
	err := s.loginHistoryService.GetLastLoginHistoryByUserID(params.User.ID, lh)
	if err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		if !errors.Is(gorm.ErrRecordNotFound, err) {
			s.logger.Error2("error getting last login history", err,
				zap.String("service", "authEventsHandler"),
				zap.String("method", "NotifyIfUserIpHasChanged"),
				zap.Int("userID", params.User.ID),
			)
		}
		return
	}

	if lh.IP.String != params.IP {
		cu := communication.CommunicatingUser{
			Email: params.User.Email,
			Phone: "",
		}
		up, err := s.userService.GetUserProfile(params.User)
		if err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
			s.logger.Error2("error getting user profile", err,
				zap.String("service", "authEventsHandler"),
				zap.String("method", "NotifyIfUserIpHasChanged"),
				zap.Int("userID", params.User.ID),
			)
		}
		emailParams := communication.SendIPChangedEmailParams{
			Email:       params.User.Email,
			FullName:    up.GetFullName(),
			LastIP:      lh.IP.String,
			CurrentIP:   params.IP,
			Device:      params.Device,
			ChangedDate: time.Now().Format("2006-01-02 15:04:05"),
		}
		s.communicationService.SendIPChangedEmail(cu, emailParams)

	}
}

func (s *eventsHandler) isWalletUser(username string) bool {
	walletUsername := s.configs.GetString("wallet.username")
	return walletUsername == username
}

func NewAuthEventsHandler(loginHistoryService user.LoginHistoryService, communicationService communication.Service, userService user.Service, configs platform.Configs, logger platform.Logger) EventsHandler {
	return &eventsHandler{
		loginHistoryService:  loginHistoryService,
		communicationService: communicationService,
		userService:          userService,
		configs:              configs,
		logger:               logger,
	}
}
