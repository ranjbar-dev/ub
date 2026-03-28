// Package auth_test tests the authentication service. Covers:
//   - Login with invalid username, wrong password, failed recaptcha, and blocked account
//   - Login with 2FA enabled (missing code, valid code) and without 2FA
//   - Login with expired refresh token renewal
//   - User registration with recaptcha failure, duplicate email, and successful flow
//   - GetUser and GetAdminUser including password-changed-after-token-issuance checks
//   - Forgot-password generation, code validation, and password update
//   - Email verification (missing code, successful)
//   - Refresh-token exchange (success, expired, blocked account)
//
// Test data: testify mocks for UserRepository, JwtService, AuthEventsHandler,
// RecaptchaManager, TwoFaManager, CommunicationService, and go-sqlmock for
// GORM database interactions.
package auth_test

import (
	"database/sql"
	"exchange-go/internal/auth"
	"exchange-go/internal/mocks"
	"exchange-go/internal/platform"
	"exchange-go/internal/response"
	"exchange-go/internal/user"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// queryMatcher is a test helper that satisfies the sqlmock.QueryMatcher
// interface by unconditionally matching any expected SQL against any actual SQL.
type queryMatcher struct {
}

// Match always returns nil, allowing any SQL query to match during tests.
func (queryMatcher) Match(expectedSQL, actualSQL string) error {
	return nil
}

func TestService_Login_UserNameDoesNotExist(t *testing.T) {
	db := &gorm.DB{}
	userRepository := new(mocks.UserRepository)
	userRepository.On("GetUserByUsername", "test@test.com", mock.Anything).Once().Return(gorm.ErrRecordNotFound)

	pe := platform.NewPasswordEncoder()
	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)
	jts := new(mocks.JwtService)
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	logUserLoginParams := auth.LogUserLoginParams{
		UserID:   0,
		Device:   "WEB",
		Email:    "test@test.com",
		IP:       "127.0.0.1",
		Type:     "FAILED",
		Password: "",
	}
	eventsHandler.On("LogUserLogin", logUserLoginParams).Once().Return()

	forgotPasswordManager := new(mocks.ForgotPasswordManager)
	recaptchaManager := new(mocks.RecaptchaManager)
	recaptchaManager.On("CheckRecaptcha", "noCheckingCurrently", "WEB", "127.0.0.1").Once().Return(true, nil)

	twoFaManager := new(mocks.TwoFaManager)

	configs := new(mocks.Configs)
	configs.On("GetString", "wallet.username").Once().Return("wallet@wallet.com")
	logger := new(mocks.Logger)
	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	params := auth.LoginParams{
		Username:  "test@test.com",
		Password:  "123456789",
		Recaptcha: "noCheckingCurrently",
		TwoFaCode: "",
		EmailCode: "",
		IP:        "127.0.0.1",
		UserAgent: "someuseragent",
	}

	res, statusCode := authService.Login(params)
	assert.Equal(t, http.StatusUnauthorized, statusCode)
	loginResponse := res.(auth.LoginFailedResponseData)

	assert.Equal(t, "invalid credentials", loginResponse.Message)
	assert.Equal(t, 401, loginResponse.Code)
	time.Sleep(20 * time.Millisecond)
	eventsHandler.AssertExpectations(t)
	recaptchaManager.AssertExpectations(t)
	userRepository.AssertExpectations(t)
	configs.AssertExpectations(t)

}

func TestService_Login_PasswordIsNotCorrect(t *testing.T) {
	db := &gorm.DB{}
	userRepository := new(mocks.UserRepository)
	pe := platform.NewPasswordEncoder()
	encodedPassword, _ := pe.GenerateFromPassword("1234567891") //is not the same as the user one entered
	u := &user.User{}
	userRepository.On("GetUserByUsername", "test@test.com", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		u = args.Get(1).(*user.User)
		u.ID = 1
		u.Password = string(encodedPassword)
		u.AccountStatus = user.AccountStatusUnblocked
		u.IsTwoFaEnabled = false
		u.RefreshToken = sql.NullString{String: "refreshtoken", Valid: true}
		u.RefreshTokenExpiry = time.Now().Add(-3 * time.Hour)
	})

	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)
	jts := new(mocks.JwtService)
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	logUserLoginParams := auth.LogUserLoginParams{
		UserID:   1,
		Device:   "WEB",
		Email:    "test@test.com",
		IP:       "127.0.0.1",
		Type:     "FAILED",
		Password: "",
	}
	eventsHandler.On("LogUserLogin", logUserLoginParams).Once().Return()

	forgotPasswordManager := new(mocks.ForgotPasswordManager)

	recaptchaManager := new(mocks.RecaptchaManager)
	recaptchaManager.On("CheckRecaptcha", "noCheckingCurrently", "WEB", "127.0.0.1").Once().Return(true, nil)

	twoFaManager := new(mocks.TwoFaManager)

	configs := new(mocks.Configs)
	configs.On("GetString", "wallet.username").Once().Return("wallet@wallet.com")
	logger := new(mocks.Logger)
	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	params := auth.LoginParams{
		Username:  "test@test.com",
		Password:  "123456789",
		Recaptcha: "noCheckingCurrently",
		TwoFaCode: "",
		EmailCode: "",
		IP:        "127.0.0.1",
		UserAgent: "someuseragent",
	}

	res, statusCode := authService.Login(params)
	assert.Equal(t, http.StatusUnauthorized, statusCode)
	loginResponse := res.(auth.LoginFailedResponseData)

	assert.Equal(t, "invalid credentials", loginResponse.Message)
	assert.Equal(t, 401, loginResponse.Code)

	time.Sleep(20 * time.Millisecond)
	eventsHandler.AssertExpectations(t)
	recaptchaManager.AssertExpectations(t)
	userRepository.AssertExpectations(t)
	configs.AssertExpectations(t)
}

func TestService_Login_RecaptchaFailed(t *testing.T) {
	db := &gorm.DB{}
	userRepository := new(mocks.UserRepository)
	pe := platform.NewPasswordEncoder()
	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)
	jts := new(mocks.JwtService)
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	forgotPasswordManager := new(mocks.ForgotPasswordManager)

	recaptchaManager := new(mocks.RecaptchaManager)
	recaptchaManager.On("CheckRecaptcha", "noCheckingCurrently", "WEB", "127.0.0.1").Once().Return(false, nil)

	twoFaManager := new(mocks.TwoFaManager)

	configs := new(mocks.Configs)
	configs.On("GetString", "wallet.username").Once().Return("wallet@wallet.com")
	logger := new(mocks.Logger)
	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	params := auth.LoginParams{
		Username:  "test@test.com",
		Password:  "123456789",
		Recaptcha: "noCheckingCurrently",
		TwoFaCode: "",
		EmailCode: "",
		IP:        "127.0.0.1",
		UserAgent: "someuseragent",
	}

	res, statusCode := authService.Login(params)
	assert.Equal(t, http.StatusUnauthorized, statusCode)
	loginResponse := res.(auth.LoginFailedResponseData)

	assert.Equal(t, "recaptcha failed", loginResponse.Message)
	assert.Equal(t, 401, loginResponse.Code)

	recaptchaManager.AssertExpectations(t)
	configs.AssertExpectations(t)
}

func TestService_Login_AccountBlocked(t *testing.T) {
	db := &gorm.DB{}
	userRepository := new(mocks.UserRepository)
	pe := platform.NewPasswordEncoder()
	encodedPassword, _ := pe.GenerateFromPassword("123456789") //is not the same as the user one entered
	u := &user.User{}
	userRepository.On("GetUserByUsername", "test@test.com", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		u = args.Get(1).(*user.User)
		u.ID = 1
		u.Password = string(encodedPassword)
		u.AccountStatus = user.AccountStatusBlocked
		u.IsTwoFaEnabled = false
		u.RefreshToken = sql.NullString{String: "refreshtoken", Valid: true}
		u.RefreshTokenExpiry = time.Now().Add(-3 * time.Hour)
	})

	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)
	jts := new(mocks.JwtService)
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	logUserLoginParams := auth.LogUserLoginParams{
		UserID:   1,
		Device:   "WEB",
		Email:    "test@test.com",
		IP:       "127.0.0.1",
		Type:     "FAILED",
		Password: "",
	}
	eventsHandler.On("LogUserLogin", logUserLoginParams).Once().Return()

	forgotPasswordManager := new(mocks.ForgotPasswordManager)
	recaptchaManager := new(mocks.RecaptchaManager)
	recaptchaManager.On("CheckRecaptcha", "noCheckingCurrently", "WEB", "127.0.0.1").Once().Return(true, nil)

	twoFaManager := new(mocks.TwoFaManager)

	configs := new(mocks.Configs)
	configs.On("GetString", "wallet.username").Once().Return("wallet@wallet.com")
	logger := new(mocks.Logger)
	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	params := auth.LoginParams{
		Username:  "test@test.com",
		Password:  "123456789",
		Recaptcha: "noCheckingCurrently",
		TwoFaCode: "",
		EmailCode: "",
		IP:        "127.0.0.1",
		UserAgent: "someuseragent",
	}

	res, statusCode := authService.Login(params)
	assert.Equal(t, http.StatusUnauthorized, statusCode)
	loginResponse := res.(auth.LoginFailedResponseData)

	assert.Equal(t, "account is blocked", loginResponse.Message)
	assert.Equal(t, 401, loginResponse.Code)

	time.Sleep(20 * time.Millisecond)
	eventsHandler.AssertExpectations(t)
	userRepository.AssertExpectations(t)
	recaptchaManager.AssertExpectations(t)
	configs.AssertExpectations(t)
}

func TestService_Login_2FaEnabled_2FaCodeNotProvided(t *testing.T) {
	db := &gorm.DB{}
	userRepository := new(mocks.UserRepository)
	pe := platform.NewPasswordEncoder()
	encodedPassword, _ := pe.GenerateFromPassword("123456789")
	u := &user.User{}
	userRepository.On("GetUserByUsername", "test@test.com", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		u = args.Get(1).(*user.User)
		u.ID = 1
		u.Password = string(encodedPassword)
		u.AccountStatus = user.AccountStatusUnblocked
		u.IsTwoFaEnabled = true
		u.RefreshToken = sql.NullString{String: "refreshtoken", Valid: true}
		u.RefreshTokenExpiry = time.Now().Add(3 * time.Hour)
	})

	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)
	jts := new(mocks.JwtService)
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	forgotPasswordManager := new(mocks.ForgotPasswordManager)
	recaptchaManager := new(mocks.RecaptchaManager)
	recaptchaManager.On("CheckRecaptcha", "noCheckingCurrently", "WEB", "127.0.0.1").Once().Return(true, nil)

	twoFaManager := new(mocks.TwoFaManager)

	configs := new(mocks.Configs)
	configs.On("GetString", "wallet.username").Twice().Return("wallet@wallet.com")
	logger := new(mocks.Logger)
	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	params := auth.LoginParams{
		Username:  "test@test.com",
		Password:  "123456789",
		Recaptcha: "noCheckingCurrently",
		TwoFaCode: "",
		EmailCode: "",
		IP:        "127.0.0.1",
		UserAgent: "someuseragent",
	}

	res, statusCode := authService.Login(params)
	assert.Equal(t, http.StatusOK, statusCode)
	loginResponse, ok := res.(auth.LoginResponseData)
	if !ok {
		t.Error("can not cast interface to struct")
		t.Fail()
	}

	assert.Equal(t, "", loginResponse.Token)
	assert.Equal(t, "", loginResponse.RefreshToken)
	assert.Equal(t, true, loginResponse.Need2fa)
	assert.Equal(t, false, loginResponse.NeedEmailCode)

	recaptchaManager.AssertExpectations(t)
	userRepository.AssertExpectations(t)
	configs.AssertExpectations(t)

}

func TestService_Login_Successful_No2Fa(t *testing.T) {
	db := &gorm.DB{}
	userRepository := new(mocks.UserRepository)
	pe := platform.NewPasswordEncoder()
	encodedPassword, _ := pe.GenerateFromPassword("123456789")
	u := &user.User{}
	userRepository.On("GetUserByUsername", "test@test.com", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		u = args.Get(1).(*user.User)
		u.ID = 1
		u.Password = string(encodedPassword)
		u.AccountStatus = user.AccountStatusUnblocked
		u.IsTwoFaEnabled = false
		u.RefreshToken = sql.NullString{String: "refreshtoken", Valid: true}
		u.RefreshTokenExpiry = time.Now().Add(3 * time.Hour)
	})

	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)
	jts := new(mocks.JwtService)
	jts.On("IssueToken", "test@test.com", mock.Anything, "127.0.0.1").Once().Return("token", nil)
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	logUserLoginParams := auth.LogUserLoginParams{
		UserID:   1,
		Device:   "WEB",
		Email:    "test@test.com",
		IP:       "127.0.0.1",
		Type:     "SUCCESSFUL",
		Password: "",
	}

	eventsHandler.On("LogUserLogin", logUserLoginParams).Once().Return()
	eventsHandler.On("NotifyIfUserIPHasChanged", mock.Anything).Once().Return()

	forgotPasswordManager := new(mocks.ForgotPasswordManager)

	recaptchaManager := new(mocks.RecaptchaManager)
	recaptchaManager.On("CheckRecaptcha", "noCheckingCurrently", "WEB", "127.0.0.1").Once().Return(true, nil)

	twoFaManager := new(mocks.TwoFaManager)

	configs := new(mocks.Configs)
	configs.On("GetString", "wallet.username").Twice().Return("wallet@wallet.com")
	logger := new(mocks.Logger)
	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	params := auth.LoginParams{
		Username:  "test@test.com",
		Password:  "123456789",
		Recaptcha: "noCheckingCurrently",
		TwoFaCode: "",
		EmailCode: "",
		IP:        "127.0.0.1",
		UserAgent: "someuseragent",
	}

	res, statusCode := authService.Login(params)
	assert.Equal(t, http.StatusOK, statusCode)
	loginResponse, ok := res.(auth.LoginResponseData)
	if !ok {
		t.Error("can not cast interface to struct")
		t.Fail()
	}

	assert.Equal(t, "token", loginResponse.Token)
	assert.Equal(t, "refreshtoken", loginResponse.RefreshToken)
	assert.Equal(t, false, loginResponse.Need2fa)
	assert.Equal(t, false, loginResponse.NeedEmailCode)
	time.Sleep(20 * time.Millisecond)
	eventsHandler.AssertExpectations(t)
	jts.AssertExpectations(t)
	recaptchaManager.AssertExpectations(t)
	userRepository.AssertExpectations(t)
	configs.AssertExpectations(t)

}

func TestService_Login_Successful_No2Fa_ExpiredRefreshToken(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(12, 1))

	userRepository := new(mocks.UserRepository)
	pe := platform.NewPasswordEncoder()
	encodedPassword, _ := pe.GenerateFromPassword("123456789")
	u := &user.User{}
	userRepository.On("GetUserByUsername", "test@test.com", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		u = args.Get(1).(*user.User)
		u.ID = 1
		u.Password = string(encodedPassword)
		u.AccountStatus = user.AccountStatusUnblocked
		u.IsTwoFaEnabled = false
		u.RefreshToken = sql.NullString{String: "refreshtoken", Valid: true}
		u.RefreshTokenExpiry = time.Now().Add(-3 * time.Hour)
	})

	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)
	jts := new(mocks.JwtService)
	jts.On("IssueToken", "test@test.com", mock.Anything, "127.0.0.1").Once().Return("token", nil)
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	logUserLoginParams := auth.LogUserLoginParams{
		UserID:   1,
		Device:   "WEB",
		Email:    "test@test.com",
		IP:       "127.0.0.1",
		Type:     "SUCCESSFUL",
		Password: "",
	}
	eventsHandler.On("LogUserLogin", logUserLoginParams).Once().Return()
	eventsHandler.On("NotifyIfUserIPHasChanged", mock.Anything).Once().Return()

	forgotPasswordManager := new(mocks.ForgotPasswordManager)

	recaptchaManager := new(mocks.RecaptchaManager)
	recaptchaManager.On("CheckRecaptcha", "noCheckingCurrently", "WEB", "127.0.0.1").Once().Return(true, nil)

	twoFaManager := new(mocks.TwoFaManager)

	configs := new(mocks.Configs)
	configs.On("GetString", "wallet.username").Twice().Return("wallet@wallet.com")
	logger := new(mocks.Logger)
	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	params := auth.LoginParams{
		Username:  "test@test.com",
		Password:  "123456789",
		Recaptcha: "noCheckingCurrently",
		TwoFaCode: "",
		EmailCode: "",
		IP:        "127.0.0.1",
		UserAgent: "someuseragent",
	}

	res, statusCode := authService.Login(params)
	assert.Equal(t, http.StatusOK, statusCode)
	loginResponse, ok := res.(auth.LoginResponseData)
	if !ok {
		t.Error("can not cast interface to struct")
		t.Fail()
	}

	assert.Equal(t, "token", loginResponse.Token)
	assert.NotEqual(t, "refreshtoken", loginResponse.RefreshToken)
	assert.Equal(t, false, loginResponse.Need2fa)
	assert.Equal(t, false, loginResponse.NeedEmailCode)

	time.Sleep(20 * time.Millisecond)
	eventsHandler.AssertExpectations(t)
	jts.AssertExpectations(t)
	recaptchaManager.AssertExpectations(t)
	userRepository.AssertExpectations(t)
	configs.AssertExpectations(t)

}

func TestService_Login_Successful_2FaEnabled(t *testing.T) {
	db := &gorm.DB{}
	userRepository := new(mocks.UserRepository)
	pe := platform.NewPasswordEncoder()
	encodedPassword, _ := pe.GenerateFromPassword("123456789")
	u := &user.User{}
	userRepository.On("GetUserByUsername", "test@test.com", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		u = args.Get(1).(*user.User)
		u.ID = 1
		u.Password = string(encodedPassword)
		u.AccountStatus = user.AccountStatusUnblocked
		u.IsTwoFaEnabled = true
		u.RefreshToken = sql.NullString{String: "refreshtoken", Valid: true}
		u.RefreshTokenExpiry = time.Now().Add(3 * time.Hour)
	})

	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)
	jts := new(mocks.JwtService)
	jts.On("IssueToken", "test@test.com", mock.Anything, "127.0.0.1").Once().Return("token", nil)
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	logUserLoginParams := auth.LogUserLoginParams{
		UserID:   1,
		Device:   "WEB",
		Email:    "test@test.com",
		IP:       "127.0.0.1",
		Type:     "SUCCESSFUL",
		Password: "",
	}
	eventsHandler.On("LogUserLogin", logUserLoginParams).Once().Return()
	eventsHandler.On("NotifyIfUserIPHasChanged", mock.Anything).Once().Return()

	forgotPasswordManager := new(mocks.ForgotPasswordManager)
	recaptchaManager := new(mocks.RecaptchaManager)
	recaptchaManager.On("CheckRecaptcha", "noCheckingCurrently", "WEB", "127.0.0.1").Once().Return(true, nil)

	twoFaManager := new(mocks.TwoFaManager)
	twoFaManager.On("CheckCode", mock.Anything, "123456").Once().Return(true)

	configs := new(mocks.Configs)
	configs.On("GetString", "wallet.username").Twice().Return("wallet@wallet.com")
	logger := new(mocks.Logger)
	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	params := auth.LoginParams{
		Username:  "test@test.com",
		Password:  "123456789",
		Recaptcha: "noCheckingCurrently",
		TwoFaCode: "123456",
		EmailCode: "",
		IP:        "127.0.0.1",
		UserAgent: "someuseragent",
	}

	res, statusCode := authService.Login(params)
	assert.Equal(t, http.StatusOK, statusCode)
	loginResponse, ok := res.(auth.LoginResponseData)
	if !ok {
		t.Error("can not cast interface to struct")
		t.Fail()
	}

	assert.Equal(t, "token", loginResponse.Token)
	assert.Equal(t, "refreshtoken", loginResponse.RefreshToken)
	assert.Equal(t, true, loginResponse.Need2fa)
	assert.Equal(t, false, loginResponse.NeedEmailCode)

	time.Sleep(20 * time.Millisecond)
	eventsHandler.AssertExpectations(t)
	twoFaManager.AssertExpectations(t)
	jts.AssertExpectations(t)
	recaptchaManager.AssertExpectations(t)
	userRepository.AssertExpectations(t)
	configs.AssertExpectations(t)

}

func TestService_Register_RecaptchaFailed(t *testing.T) {
	db := &gorm.DB{}
	userRepository := new(mocks.UserRepository)
	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)
	jts := new(mocks.JwtService)
	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	forgotPasswordManager := new(mocks.ForgotPasswordManager)

	recaptchaManager := new(mocks.RecaptchaManager)
	recaptchaManager.On("CheckRecaptcha", "recaptcha", "WEB", "127.0.0.1").Once().Return(false, nil)

	twoFaManager := new(mocks.TwoFaManager)

	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	params := auth.RegisterParams{
		Email:     "test@test.com",
		Password:  "123456789",
		IP:        "127.0.0.1",
		UserAgent: "some userAgent",
		Recaptcha: "recaptcha",
	}

	res, statusCode := authService.Register(params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "recaptcha failed", res.Message)
	recaptchaManager.AssertExpectations(t)
}

func TestService_Register_EmailAlreadyTaken(t *testing.T) {
	db := &gorm.DB{}
	userRepository := new(mocks.UserRepository)
	userRepository.On("GetEvenBlockedUserByEmail", "test@test.com", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*user.User)
		u.ID = 1
	})
	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)
	jts := new(mocks.JwtService)
	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	forgotPasswordManager := new(mocks.ForgotPasswordManager)

	recaptchaManager := new(mocks.RecaptchaManager)
	recaptchaManager.On("CheckRecaptcha", "recaptcha", "WEB", "127.0.0.1").Once().Return(true, nil)

	twoFaManager := new(mocks.TwoFaManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	params := auth.RegisterParams{
		Email:     "test@test.com",
		Password:  "123456789",
		IP:        "127.0.0.1",
		UserAgent: "some userAgent",
		Recaptcha: "recaptcha",
	}

	res, statusCode := authService.Register(params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "email is already taken", res.Message)
	recaptchaManager.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestService_Register_Successful(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()

	dbMock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("INSERT INTO user_profiles").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("INSERT INTO user_balances").WillReturnResult(sqlmock.NewResult(10, 1))
	dbMock.ExpectExec("INSERT INTO users_permissions").WillReturnResult(sqlmock.NewResult(10, 1))

	dbMock.ExpectCommit()

	userRepository := new(mocks.UserRepository)
	userRepository.On("GetEvenBlockedUserByEmail", "test@test.com", mock.Anything).Once().Return(nil)
	userLevelService := new(mocks.UserLevelService)
	userLevelService.On("GetLevelByCode", int64(user.UserLevelVip0Code)).Once().Return(user.Level{ID: 1}, nil)
	userPermissionManager := new(mocks.UserPermissionManager)
	permissions := []user.Permission{
		{
			ID:   1,
			Name: user.PermissionWithdraw,
		},
		{
			ID:   2,
			Name: user.PermissionDeposit,
		},
		{
			ID:   3,
			Name: user.PermissionExchange,
		},
	}
	userPermissionManager.On("GetAllPermissions").Once().Return(permissions)

	userBalanceService := new(mocks.UserBalanceService)
	userBalanceService.On("GenerateBalancesAndAddressForUser", mock.Anything).Once().Return()

	jts := new(mocks.JwtService)
	pe := platform.NewPasswordEncoder()

	communicationService := new(mocks.CommunicationService)
	communicationService.On("SendVerificationEmailToUser", mock.Anything, mock.Anything).Once().Return()

	eventsHandler := new(mocks.AuthEventsHandler)
	forgotPasswordManager := new(mocks.ForgotPasswordManager)

	recaptchaManager := new(mocks.RecaptchaManager)
	recaptchaManager.On("CheckRecaptcha", "recaptcha", "WEB", "127.0.0.1").Once().Return(true, nil)

	twoFaManager := new(mocks.TwoFaManager)

	configs := new(mocks.Configs)
	configs.On("GetDomain").Once().Return("localhost")

	logger := new(mocks.Logger)

	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	params := auth.RegisterParams{
		Email:     "test@test.com",
		Password:  "123456789",
		IP:        "127.0.0.1",
		UserAgent: "some userAgent",
		Recaptcha: "recaptcha",
	}

	res, statusCode := authService.Register(params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)
	time.Sleep(20 * time.Millisecond)

	recaptchaManager.AssertExpectations(t)
	userRepository.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	communicationService.AssertExpectations(t)
	userLevelService.AssertExpectations(t)
	userPermissionManager.AssertExpectations(t)
	configs.AssertExpectations(t)
}

func TestService_GetUser(t *testing.T) {
	db := &gorm.DB{}
	userRepository := new(mocks.UserRepository)
	userRepository.On("GetUserByUsername", "username", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*user.User)
		u.ID = 2
	})

	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)

	jts := new(mocks.JwtService)
	jts.On("GetUsernameFromToken", "token").Once().Return("username", time.Now().Add(-10*time.Hour), nil)

	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	forgotPasswordManager := new(mocks.ForgotPasswordManager)
	recaptchaManager := new(mocks.RecaptchaManager)
	twoFaManager := new(mocks.TwoFaManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	token := "token"
	loggedInUser, err := authService.GetUser(token)
	assert.Nil(t, err)
	assert.Equal(t, 2, loggedInUser.ID)

	jts.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestService_GetUser_password_changed_after_token_issuance(t *testing.T) {
	db := &gorm.DB{}
	userRepository := new(mocks.UserRepository)
	userRepository.On("GetUserByUsername", "username", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*user.User)
		u.ID = 2
		u.PasswordChangedAt = sql.NullTime{Time: time.Now(), Valid: true}
	})

	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)

	jts := new(mocks.JwtService)
	jts.On("GetUsernameFromToken", "token").Once().Return("username", time.Now().Add(-10*time.Hour), nil)

	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	forgotPasswordManager := new(mocks.ForgotPasswordManager)
	recaptchaManager := new(mocks.RecaptchaManager)
	twoFaManager := new(mocks.TwoFaManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	token := "token"
	loggedInUser, err := authService.GetUser(token)
	assert.NotNil(t, err)
	assert.Nil(t, loggedInUser)
	assert.Equal(t, "invalid token", err.Error())

	jts.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestService_GetAdminUser(t *testing.T) {
	db := &gorm.DB{}
	userRepository := new(mocks.UserRepository)
	userRepository.On("GetAdminUserByUsername", "username", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*user.User)
		u.ID = 2
	})

	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)

	jts := new(mocks.JwtService)
	jts.On("GetUsernameFromToken", "token").Once().Return("username", time.Now().Add(-10*time.Hour), nil)

	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	forgotPasswordManager := new(mocks.ForgotPasswordManager)
	recaptchaManager := new(mocks.RecaptchaManager)
	twoFaManager := new(mocks.TwoFaManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	token := "token"
	loggedInUser, err := authService.GetAdminUser(token)
	assert.Nil(t, err)
	assert.Equal(t, 2, loggedInUser.ID)

	jts.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestService_GetAdminUser_password_changed_after_token_issuance(t *testing.T) {
	db := &gorm.DB{}
	userRepository := new(mocks.UserRepository)
	userRepository.On("GetAdminUserByUsername", "username", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*user.User)
		u.ID = 2
		u.PasswordChangedAt = sql.NullTime{Time: time.Now(), Valid: true}
	})

	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)

	jts := new(mocks.JwtService)
	jts.On("GetUsernameFromToken", "token").Once().Return("username", time.Now().Add(-10*time.Hour), nil)

	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	forgotPasswordManager := new(mocks.ForgotPasswordManager)
	recaptchaManager := new(mocks.RecaptchaManager)
	twoFaManager := new(mocks.TwoFaManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	token := "token"
	loggedInUser, err := authService.GetAdminUser(token)
	assert.NotNil(t, err)
	assert.Nil(t, loggedInUser)
	assert.Equal(t, "invalid token", err.Error())

	jts.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestService_ForgotPassword_Successful(t *testing.T) {
	db := &gorm.DB{}
	userRepository := new(mocks.UserRepository)
	userRepository.On("GetUserByUsername", "test@test.com", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*user.User)
		u.ID = 1
	})

	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)

	jts := new(mocks.JwtService)

	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	forgotPasswordManager := new(mocks.ForgotPasswordManager)
	forgotPasswordManager.On("GenerateForgotPasswordAndSendEmail", mock.Anything, "WEB", "127.0.0.1").Once().Return(nil)
	recaptchaManager := new(mocks.RecaptchaManager)
	recaptchaManager.On("CheckRecaptcha", "noCheckingCurrently", "WEB", "127.0.0.1").Once().Return(true, nil)
	twoFaManager := new(mocks.TwoFaManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	forgotPasswordParams := auth.ForgotPasswordParams{
		Email:     "test@test.com",
		Recaptcha: "noCheckingCurrently",
		IP:        "127.0.0.1",
		UserAgent: "userAgent",
	}
	_, statusCode := authService.ForgotPassword(forgotPasswordParams)
	assert.Equal(t, http.StatusOK, statusCode)

	forgotPasswordManager.AssertExpectations(t)
	recaptchaManager.AssertExpectations(t)
	userRepository.AssertExpectations(t)

}

func TestService_ForgotPasswordUpdate_Failed_PasswordAndConfirmedAreNotTheSame(t *testing.T) {
	db := &gorm.DB{}
	userRepository := new(mocks.UserRepository)
	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)
	jts := new(mocks.JwtService)
	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	forgotPasswordManager := new(mocks.ForgotPasswordManager)
	recaptchaManager := new(mocks.RecaptchaManager)
	twoFaManager := new(mocks.TwoFaManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	forgotPasswordUpdateParams := auth.ForgotPasswordUpdateParams{
		Email:     "test@test.com",
		Password:  "123456789",
		Confirmed: "1234567890",
		Code:      "123456",
		IP:        "127.0.0.1",
		UserAgent: "userAgent",
	}
	res, statusCode := authService.ForgotPasswordUpdate(forgotPasswordUpdateParams)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, "new password and confirm password does not match", res.Message)
}

func TestService_ForgotPasswordUpdate_Failed_CodeIsNotCorrect(t *testing.T) {
	db := &gorm.DB{}
	userRepository := new(mocks.UserRepository)
	userRepository.On("GetUserByUsername", "test@test.com", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*user.User)
		u.ID = 1
	})

	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)
	jts := new(mocks.JwtService)
	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	forgotPasswordManager := new(mocks.ForgotPasswordManager)
	forgotPasswordManager.On("IsCodeCorrect", mock.Anything, "123456").Once().Return(false)
	recaptchaManager := new(mocks.RecaptchaManager)
	twoFaManager := new(mocks.TwoFaManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	forgotPasswordUpdateParams := auth.ForgotPasswordUpdateParams{
		Email:     "test@test.com",
		Password:  "123456789",
		Confirmed: "123456789",
		Code:      "123456",
		IP:        "127.0.0.1",
		UserAgent: "userAgent",
	}
	res, statusCode := authService.ForgotPasswordUpdate(forgotPasswordUpdateParams)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, "code is not correct", res.Message)

	forgotPasswordManager.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestService_ForgotPasswordUpdate_Successful(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(12, 1))

	userRepository := new(mocks.UserRepository)
	userRepository.On("GetUserByUsername", "test@test.com", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*user.User)
		u.ID = 1
	})

	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)

	jts := new(mocks.JwtService)

	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	communicationService.On("SendPasswordChangedEmail", mock.Anything, mock.Anything).Once().Return()

	eventsHandler := new(mocks.AuthEventsHandler)
	forgotPasswordManager := new(mocks.ForgotPasswordManager)
	forgotPasswordManager.On("IsCodeCorrect", mock.Anything, "123456").Once().Return(true)
	forgotPasswordManager.On("DeleteKey", mock.Anything).Once().Return(nil)

	recaptchaManager := new(mocks.RecaptchaManager)
	twoFaManager := new(mocks.TwoFaManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	forgotPasswordUpdateParams := auth.ForgotPasswordUpdateParams{
		Email:     "test@test.com",
		Password:  "123456789",
		Confirmed: "123456789",
		Code:      "123456",
		IP:        "127.0.0.1",
		UserAgent: "userAgent",
	}
	_, statusCode := authService.ForgotPasswordUpdate(forgotPasswordUpdateParams)
	assert.Equal(t, http.StatusOK, statusCode)

	time.Sleep(50 * time.Millisecond)
	communicationService.AssertExpectations(t)
	forgotPasswordManager.AssertExpectations(t)
	recaptchaManager.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestService_VerifyUser_Fail_CodeNotFound(t *testing.T) {
	db := &gorm.DB{}
	userRepository := new(mocks.UserRepository)
	userRepository.On("GetUserByVerificationCode", "verificationCode", mock.Anything).Once().Return(gorm.ErrRecordNotFound)

	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)

	jts := new(mocks.JwtService)

	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	forgotPasswordManager := new(mocks.ForgotPasswordManager)
	recaptchaManager := new(mocks.RecaptchaManager)
	twoFaManager := new(mocks.TwoFaManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	verifyEmailParams := auth.VerifyEmailParams{
		Code: "verificationCode",
	}
	res, statusCode := authService.VerifyEmail(verifyEmailParams)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, "code not found", res.Message)

	userRepository.AssertExpectations(t)
}

func TestService_VerifyUser_Successful(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(12, 1))

	userRepository := new(mocks.UserRepository)
	userRepository.On("GetUserByVerificationCode", "verificationCode", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*user.User)
		u.ID = 1
	})

	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)

	jts := new(mocks.JwtService)

	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	forgotPasswordManager := new(mocks.ForgotPasswordManager)
	recaptchaManager := new(mocks.RecaptchaManager)
	twoFaManager := new(mocks.TwoFaManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	verifyEmailParams := auth.VerifyEmailParams{
		Code: "verificationCode",
	}
	_, statusCode := authService.VerifyEmail(verifyEmailParams)
	assert.Equal(t, http.StatusOK, statusCode)

	userRepository.AssertExpectations(t)
}

func TestService_GetTokenByRefreshToken_Successful(t *testing.T) {
	db := &gorm.DB{}
	userRepository := new(mocks.UserRepository)
	userRepository.On("GetUserByRefreshToken", "refreshtoken", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*user.User)
		u.ID = 1
		u.Email = "test@test.com"
		u.AccountStatus = user.AccountStatusUnblocked
		u.IsTwoFaEnabled = true
		u.RefreshToken = sql.NullString{String: "refreshtoken", Valid: true}
		u.RefreshTokenExpiry = time.Now().Add(3 * time.Hour)
	})

	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)

	jts := new(mocks.JwtService)
	jts.On("IssueToken", "test@test.com", mock.Anything, "127.0.0.1").Once().Return("token", nil)

	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	forgotPasswordManager := new(mocks.ForgotPasswordManager)
	recaptchaManager := new(mocks.RecaptchaManager)
	twoFaManager := new(mocks.TwoFaManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	params := auth.GetTokenByRefreshTokenParams{
		Refresh:   "refreshtoken",
		IP:        "127.0.0.1",
		UserAgent: "web",
	}
	res, statusCode := authService.GetTokenByRefreshToken(params)
	assert.Equal(t, http.StatusOK, statusCode)
	loginResponse, ok := res.(auth.LoginResponseData)
	if !ok {
		t.Error("can not cast interface to struct")
		t.Fail()
	}

	assert.Equal(t, "token", loginResponse.Token)
	assert.Equal(t, "refreshtoken", loginResponse.RefreshToken)
	assert.Equal(t, false, loginResponse.Need2fa)
	assert.Equal(t, false, loginResponse.NeedEmailCode)

	jts.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestService_GetTokenByRefreshToken_Failed_ExpiryTimePassed(t *testing.T) {
	db := &gorm.DB{}
	userRepository := new(mocks.UserRepository)
	userRepository.On("GetUserByRefreshToken", "refreshtoken", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*user.User)
		u.ID = 1
		u.Email = "test@test.com"
		u.AccountStatus = user.AccountStatusUnblocked
		u.IsTwoFaEnabled = true
		u.RefreshToken = sql.NullString{String: "refreshtoken", Valid: true}
		u.RefreshTokenExpiry = time.Now().Add(-3 * time.Hour)
	})

	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)

	jts := new(mocks.JwtService)

	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	forgotPasswordManager := new(mocks.ForgotPasswordManager)
	recaptchaManager := new(mocks.RecaptchaManager)
	twoFaManager := new(mocks.TwoFaManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	params := auth.GetTokenByRefreshTokenParams{
		Refresh:   "refreshtoken",
		IP:        "127.0.0.1",
		UserAgent: "web",
	}
	res, statusCode := authService.GetTokenByRefreshToken(params)
	result, ok := res.(response.APIResponse)
	if !ok {
		t.Error("can not cast interface to struct")
		t.Fail()
	}
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, "refresh token expired", result.Message)

	userRepository.AssertExpectations(t)
}

func TestService_GetTokenByRefreshToken_Failed_AccountBlocked(t *testing.T) {
	db := &gorm.DB{}
	userRepository := new(mocks.UserRepository)
	userRepository.On("GetUserByRefreshToken", "refreshtoken", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*user.User)
		u.ID = 1
		u.Email = "test@test.com"
		u.AccountStatus = user.AccountStatusBlocked
		u.IsTwoFaEnabled = true
		u.RefreshToken = sql.NullString{String: "refreshtoken", Valid: true}
		u.RefreshTokenExpiry = time.Now().Add(3 * time.Hour)
	})

	userLevelService := new(mocks.UserLevelService)
	userPermissionManager := new(mocks.UserPermissionManager)
	userBalanceService := new(mocks.UserBalanceService)

	jts := new(mocks.JwtService)

	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	eventsHandler := new(mocks.AuthEventsHandler)
	forgotPasswordManager := new(mocks.ForgotPasswordManager)
	recaptchaManager := new(mocks.RecaptchaManager)
	twoFaManager := new(mocks.TwoFaManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	authService := auth.NewAuthService(db, userRepository, userLevelService, userPermissionManager, userBalanceService, jts,
		pe, communicationService, eventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, configs, logger)

	params := auth.GetTokenByRefreshTokenParams{
		Refresh:   "refreshtoken",
		IP:        "127.0.0.1",
		UserAgent: "web",
	}
	res, statusCode := authService.GetTokenByRefreshToken(params)
	result, ok := res.(auth.LoginFailedResponseData)
	if !ok {
		t.Error("can not cast interface to struct")
		t.Fail()
	}
	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, "account is blocked", result.Message)

	userRepository.AssertExpectations(t)
}
