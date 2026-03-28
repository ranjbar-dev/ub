package auth

import (
	"crypto/md5"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"errors"
	"exchange-go/internal/communication"
	"exchange-go/internal/jwt"
	"exchange-go/internal/platform"
	"exchange-go/internal/response"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"exchange-go/internal/userdevice"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LoginResponseData struct {
	Token         string `json:"token"`
	RefreshToken  string `json:"refreshToken"`
	Need2fa       bool   `json:"need2fa"`
	NeedEmailCode bool   `json:"needEmailCode"`
	IsNewDevice   bool   `json:"isNewDevice"`
}

type LoginFailedResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type LoginParams struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Recaptcha string `json:"recaptcha" binding:"required"`
	TwoFaCode string `json:"2fa_code"`
	EmailCode string `json:"email_code"`
	IP        string
	UserAgent string
}

type RegisterParams struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	Recaptcha string `json:"recaptcha" binding:"required"`
	ReferKey  string `json:"refer_key"`
	IP        string
	UserAgent string
}

type ForgotPasswordParams struct {
	Email     string `json:"email" binding:"required,email"`
	Recaptcha string `json:"recaptcha" binding:"required"`
	IP        string
	UserAgent string
}

type ForgotPasswordUpdateParams struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	Confirmed string `json:"confirmed" binding:"required"`
	Code      string `json:"code" binding:"required"`
	IP        string
	UserAgent string
}

type VerifyEmailParams struct {
	Code string `json:"code" binding:"required"`
}

type GetTokenByRefreshTokenParams struct {
	Refresh   string `json:"refresh" binding:"required"`
	IP        string
	UserAgent string
}

// Service is the core authentication service that handles user login, registration,
// token management, and password recovery flows.
type Service interface {
	// Login authenticates a user with credentials and returns JWT tokens.
	Login(params LoginParams) (response interface{}, statusCode int)
	// Register creates a new user account and returns the API response.
	Register(params RegisterParams) (apiResponse response.APIResponse, statusCode int)
	// GetUser validates a JWT token and returns the associated user. Used by auth middleware.
	GetUser(token string) (*user.User, error)
	// GetAdminUser validates a JWT token and returns the user only if they have admin privileges.
	GetAdminUser(token string) (*user.User, error)
	// ForgotPassword initiates the password reset flow by sending a reset email.
	ForgotPassword(params ForgotPasswordParams) (apiResponse response.APIResponse, statusCode int)
	// ForgotPasswordUpdate completes the password reset flow with a new password and verification code.
	ForgotPasswordUpdate(params ForgotPasswordUpdateParams) (apiResponse response.APIResponse, statusCode int)
	// VerifyEmail confirms a user's email address using the provided verification code.
	VerifyEmail(params VerifyEmailParams) (apiResponse response.APIResponse, statusCode int)
	// GetTokenByRefreshToken issues a new JWT access token from a valid refresh token.
	GetTokenByRefreshToken(params GetTokenByRefreshTokenParams) (response interface{}, statusCode int)
}

type service struct {
	db                    *gorm.DB
	userRepository        user.Repository
	userLevelService      user.LevelService
	userPermissionManager user.PermissionManager
	userBalanceService    userbalance.Service
	jwtService            jwt.Service
	passwordEncoder       platform.PasswordEncoder
	communicationService  communication.Service
	eventsHandler         EventsHandler
	forgotPasswordManager user.ForgotPasswordManager
	recaptchaManager      user.RecaptchaManager
	twoFaManager          user.TwoFaManager
	configs               platform.Configs
	logger                platform.Logger
}

type loginExtraData struct {
	need2fa       bool
	needEmailCode bool
	isNewDevice   bool
}

//the response of the login api is different from other apis for backward compatibility
func (s *service) Login(params LoginParams) (resp interface{}, statusCode int) {
	device := userdevice.GetDeviceUsingUserAgent(params.UserAgent)
	shouldLogActivity := true
	logUserLoginParams := LogUserLoginParams{
		UserID: 0,
		Device: device,
		IP:     params.IP,
		Email:  params.Username,
		Type:   user.UserLoginHistoryTypeFailed,
	}
	defer func() {
		if shouldLogActivity {
			go s.eventsHandler.LogUserLogin(logUserLoginParams)
		}
	}()
	params.Username = strings.Trim(params.Username, "")
	params.Password = strings.Trim(params.Password, "")
	params.TwoFaCode = strings.Trim(params.TwoFaCode, "")
	params.EmailCode = strings.Trim(params.EmailCode, "")

	statusCode = http.StatusUnauthorized
	response := LoginResponseData{}

	//checking recaptcha for non wallet users
	if !s.isWalletUser(params.Username) {
		isRecaptchaSuccessful, err := s.recaptchaManager.CheckRecaptcha(params.Recaptcha, device, params.IP)
		if err != nil {
			shouldLogActivity = false
			s.logger.Error2("error in checking recaptcha", err,
				zap.String("service", "authService"),
				zap.String("method", "Login"),
				zap.String("username", params.Username),
			)
			statusCode = http.StatusInternalServerError
			return response, statusCode
		}

		if !isRecaptchaSuccessful {

			shouldLogActivity = false
			return LoginFailedResponseData{Code: statusCode, Message: "recaptcha failed"}, statusCode
		}
	}

	u := &user.User{}
	err := s.userRepository.GetUserByUsername(params.Username, u)
	logUserLoginParams.UserID = u.ID

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("error in getting user ", err,
			zap.String("service", "authService"),
			zap.String("method", "Login"),
			zap.String("username", params.Username),
		)
		statusCode = http.StatusInternalServerError
		return response, statusCode
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return LoginFailedResponseData{Code: statusCode, Message: "invalid credentials"}, statusCode
	}
	err = s.passwordEncoder.CompareHashAndPassword(u.Password, params.Password)
	if err != nil {
		return LoginFailedResponseData{Code: statusCode, Message: "invalid credentials"}, statusCode
	}

	if u.AccountStatus == user.AccountStatusBlocked {
		return LoginFailedResponseData{Code: statusCode, Message: "account is blocked"}, statusCode
	}

	//the following lines check if user enabled 2fa and check it if it is correct
	extraData := s.getExtraDataForLogin(*u)

	errorMessage, shouldLogin := s.shouldLoginUser(*u, extraData, params)

	if errorMessage != "" {
		return LoginFailedResponseData{Code: statusCode, Message: errorMessage}, statusCode
	}

	statusCode = http.StatusOK
	if shouldLogin {
		token, err := s.jwtService.IssueToken(params.Username, params.UserAgent, params.IP)
		if err != nil {
			statusCode = http.StatusInternalServerError
			s.logger.Error2("error in issuing token", err,
				zap.String("service", "authService"),
				zap.String("method", "Login"),
				zap.String("username", params.Username),
			)
			return response, statusCode

		}
		refreshToken := ""
		if token != "" {
			refreshToken, err = s.getRefreshTokenForUser(u)
		}
		if err != nil {
			s.logger.Error2("error in getting user ", err,
				zap.String("service", "authService"),
				zap.String("method", "Login"),
				zap.Int("userID", u.ID),
				zap.String("username", params.Username),
			)
		}
		logUserLoginParams.Type = user.UserLoginHistoryTypeSuccessful
		response = LoginResponseData{
			Token:         token,
			RefreshToken:  refreshToken,
			Need2fa:       extraData.need2fa,
			NeedEmailCode: extraData.needEmailCode,
			IsNewDevice:   false,
		}
		notifyParams := NotifyIPChangeParams{
			User:   *u,
			Device: device,
			IP:     params.IP,
		}
		go s.eventsHandler.NotifyIfUserIPHasChanged(notifyParams)
		return response, statusCode
	}

	shouldLogActivity = false
	response = LoginResponseData{
		Token:         "",
		RefreshToken:  "",
		Need2fa:       extraData.need2fa,
		NeedEmailCode: extraData.needEmailCode,
		IsNewDevice:   false,
	}

	return response, statusCode

}

func (s *service) shouldLoginUser(u user.User, data loginExtraData, params LoginParams) (errorMessage string, shouldLogin bool) {
	if !data.need2fa && !data.needEmailCode {
		return "", true
	}

	if data.need2fa {
		if params.TwoFaCode != "" {
			if s.twoFaManager.CheckCode(u, params.TwoFaCode) {
				shouldLogin = true
			} else {
				errorMessage = "2fa code is not correct"
				shouldLogin = false
			}
		} else {
			shouldLogin = false
		}
	}
	return errorMessage, shouldLogin
	//todo we should check the email code later when the feature implemented in front
}

func (s *service) isWalletUser(username string) bool {
	walletUsername := s.configs.GetString("wallet.username")
	return walletUsername == username
}

func (s *service) getExtraDataForLogin(u user.User) (data loginExtraData) {
	//currently we do not care about isEmailNeeded and isNewDevice these are the feature should be implemented later
	if s.isWalletUser(u.Email) {
		return data
	}

	if u.IsTwoFaEnabled {
		data.need2fa = true
	}

	return data
}

func (s *service) getRefreshTokenForUser(u *user.User) (refreshToken string, err error) {
	if u.RefreshTokenExpiry.Before(time.Now()) {
		bs := make([]byte, 32)
		rand.Read(bs) //this method never returns error read the doc
		refreshToken := hex.EncodeToString(bs)
		u.RefreshToken = sql.NullString{String: refreshToken, Valid: true}
		u.RefreshTokenExpiry = time.Now().Add(6 * 30 * 24 * time.Hour) //6 month
		err = s.db.Omit(clause.Associations).Save(u).Error
		if err != nil {
			return "", err
		}
	}
	return u.RefreshToken.String, nil
}

func (s *service) Register(params RegisterParams) (apiResponse response.APIResponse, statusCode int) {
	email := strings.Trim(params.Email, "")
	password := strings.Trim(params.Password, "")
	if email == "" {
		return response.Error("please insert a valid email", http.StatusUnprocessableEntity, nil)
	}

	if password == "" {
		return response.Error("please insert a valid password", http.StatusUnprocessableEntity, nil)
	}

	device := userdevice.GetDeviceUsingUserAgent(params.UserAgent)
	//check captcha
	isRecaptchaSuccessful, err := s.recaptchaManager.CheckRecaptcha(params.Recaptcha, device, params.IP)
	if err != nil {
		s.logger.Error2("error in checking recaptcha", err,
			zap.String("service", "authService"),
			zap.String("method", "Register"),
			zap.String("email", email),
		)
		return response.Error("something went wrong with recaptcha checking", http.StatusUnprocessableEntity, nil)
	}

	if !isRecaptchaSuccessful {
		return response.Error("recaptcha failed", http.StatusUnprocessableEntity, nil)
	}

	formerUser := &user.User{}
	err = s.userRepository.GetEvenBlockedUserByEmail(email, formerUser)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("error getting user ", err,
			zap.String("service", "authService"),
			zap.String("method", "Register"),
			zap.String("email", email),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if formerUser.ID > 0 {
		return response.Error("email is already taken", http.StatusUnprocessableEntity, nil)
	}

	tx := s.db.Begin()
	err = tx.Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error beggining transaction", err,
			zap.String("service", "authService"),
			zap.String("method", "Register"),
			zap.String("email", email),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	passwordBytes, err := s.passwordEncoder.GenerateFromPassword(password)
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error generating password hash", err,
			zap.String("service", "authService"),
			zap.String("method", "Register"),
			zap.String("email", email),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	passwordHash := string(passwordBytes)

	verificationCode := s.generateVerificationCodeForUser()

	privateChannelName := s.generatePrivateChannelNameForUser()

	defaultUserLevel, err := s.userLevelService.GetLevelByCode(user.UserLevelVip0Code)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		s.logger.Error2("error getting level ", err,
			zap.String("service", "authService"),
			zap.String("method", "Register"),
			zap.Int("code", user.UserLevelVip0Code),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	userLevelID := defaultUserLevel.ID

	newUser := &user.User{
		Email:                  email,
		Password:               passwordHash,
		Kyc:                    0,
		Status:                 user.StatusRegistered,
		AccountStatus:          user.AccountStatusUnblocked,
		ExchangeNumber:         0,
		IsTwoFaEnabled:         false,
		UbID:                   s.randSeq(10), //just to be sure its unique we would update this later
		VerificationCode:       verificationCode,
		PrivateChannelName:     privateChannelName,
		ExchangeVolumeCoinCode: "BTC",
		ExchangeVolumeAmount:   "0.0",
		UserLevelID:            userLevelID,
	}

	err = tx.Omit(clause.Associations).Save(newUser).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error saving user", err,
			zap.String("service", "authService"),
			zap.String("method", "Register"),
			zap.String("email", email),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	//update newUser to set ubID
	ubID := s.getUbID(*newUser)
	newUser.UbID = ubID
	err = tx.Omit(clause.Associations).Save(newUser).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error updating user", err,
			zap.String("service", "authService"),
			zap.String("method", "Register"),
			zap.String("email", email),
			zap.Int("userID", newUser.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	//create user profile
	userProfile := &user.Profile{
		UserID:         newUser.ID,
		RegistrationIP: sql.NullString{String: params.IP, Valid: true},
		ReferKey:       sql.NullString{String: params.ReferKey, Valid: true},
		TrustLevel:     0,
	}
	err = tx.Omit(clause.Associations).Save(userProfile).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error saving userProfile", err,
			zap.String("service", "authService"),
			zap.String("method", "Register"),
			zap.String("email", email),
			zap.Int("userID", newUser.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	permissions := s.userPermissionManager.GetAllPermissions()
	var ups []user.UsersPermissions
	for _, p := range permissions {
		up := user.UsersPermissions{
			UserID:           newUser.ID,
			UserPermissionID: p.ID,
		}
		ups = append(ups, up)
	}

	err = tx.Omit(clause.Associations).Create(ups).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error creating userPermissions", err,
			zap.String("service", "authService"),
			zap.String("method", "Register"),
			zap.String("email", email),
			zap.Int("userID", newUser.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error commiting transaction", err,
			zap.String("service", "authService"),
			zap.String("method", "Register"),
			zap.String("email", email),
			zap.Int("userID", newUser.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	//send verification code to user
	go s.sendVerificationEmail(*newUser)

	//create balances and addresses
	go s.userBalanceService.GenerateBalancesAndAddressForUser(*newUser)

	res := make(map[string]string)
	return response.Success(res, "")
}

func (s *service) ForgotPassword(params ForgotPasswordParams) (apiResponse response.APIResponse, statusCode int) {
	email := strings.Trim(params.Email, "")

	device := userdevice.GetDeviceUsingUserAgent(params.UserAgent)
	// this email is used for checking the health of our email sending status
	// so we do not check recaptcha for it
	if email != "behkamegit@gmail.com" {
		isRecaptchaSuccessful, err := s.recaptchaManager.CheckRecaptcha(params.Recaptcha, device, params.IP)
		if err != nil {
			s.logger.Error2("error in checking recaptcha", err,
				zap.String("service", "authService"),
				zap.String("method", "ForgotPassword"),
				zap.String("email", email),
				zap.String("device", device),
			)
			return response.Error("something went wrong with recaptcha checking", http.StatusUnprocessableEntity, nil)
		}

		if !isRecaptchaSuccessful {
			return response.Error("recaptcha failed", http.StatusUnprocessableEntity, nil)
		}

	}

	u := &user.User{}
	err := s.userRepository.GetUserByUsername(email, u)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("error in getting user", err,
			zap.String("service", "authService"),
			zap.String("method", "ForgotPassword"),
			zap.String("email", email),
			zap.String("device", device),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if u.ID > 0 { //it means we have user with this email
		err := s.forgotPasswordManager.GenerateForgotPasswordAndSendEmail(*u, device, params.IP)
		if err != nil {
			s.logger.Error2("error generating fogotPassword and sending email", err,
				zap.String("service", "authService"),
				zap.String("method", "ForgotPassword"),
				zap.String("email", email),
				zap.String("device", device),
			)
		}
	}

	return response.Success(nil, "")
}

func (s *service) ForgotPasswordUpdate(params ForgotPasswordUpdateParams) (apiResponse response.APIResponse, statusCode int) {
	email := strings.Trim(params.Email, "")
	password := strings.Trim(params.Password, "")
	confirmed := strings.Trim(params.Confirmed, "")
	code := strings.Trim(params.Code, "")

	if password != confirmed {
		return response.Error("new password and confirm password does not match", http.StatusUnprocessableEntity, nil)
	}
	u := &user.User{}
	err := s.userRepository.GetUserByUsername(email, u)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("error getting user ", err,
			zap.String("service", "authService"),
			zap.String("method", "ForgotPasswordUpdate"),
			zap.String("email", email),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || u.ID == 0 {
		return response.Error("email not found", http.StatusUnprocessableEntity, nil)
	}

	isCodeCorrect := s.forgotPasswordManager.IsCodeCorrect(*u, code)
	if !isCodeCorrect {
		return response.Error("code is not correct", http.StatusUnprocessableEntity, nil)
	}

	passwordBytes, err := s.passwordEncoder.GenerateFromPassword(password)
	if err != nil {
		s.logger.Error2("error generating password hash", err,
			zap.String("service", "authService"),
			zap.String("method", "ForgotPasswordUpdate"),
			zap.String("email", email),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	passwordHash := string(passwordBytes)
	u.Password = passwordHash
	u.PasswordChangedAt = sql.NullTime{Time: time.Now(), Valid: true}
	err = s.db.Omit(clause.Associations).Save(u).Error
	if err != nil {
		s.logger.Error2("error saving user", err,
			zap.String("service", "authService"),
			zap.String("method", "ForgotPasswordUpdate"),
			zap.String("email", email),
			zap.Int("userID", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	go func() {
		err := s.forgotPasswordManager.DeleteKey(*u)
		if err != nil {
			s.logger.Warn("error deleting redis key",
				zap.Error(err),
				zap.String("service", "authService"),
				zap.String("method", "ForgotPasswordUpdate"),
				zap.Int("userID", u.ID),
			)
		}
	}()

	device := userdevice.GetDeviceUsingUserAgent(params.UserAgent)
	emailParams := communication.PasswordChangedEmailParams{
		Email:         u.Email,
		CurrentIP:     params.IP,
		CurrentDevice: device,
		ChangedDate:   time.Now().Format("2006-01-02 15:04:05"),
	}
	cu := communication.CommunicatingUser{
		Email: u.Email,
		Phone: "",
	}
	go s.communicationService.SendPasswordChangedEmail(cu, emailParams)

	return response.Success(nil, "")
}

func (s *service) VerifyEmail(params VerifyEmailParams) (apiResponse response.APIResponse, statusCode int) {
	code := strings.Trim(params.Code, "")

	u := &user.User{}
	err := s.userRepository.GetUserByVerificationCode(code, u)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("error getting user", err,
			zap.String("service", "authService"),
			zap.String("method", "VerifyEmail"),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || u.ID == 0 {
		return response.Error("code not found", http.StatusUnprocessableEntity, nil)
	}

	if u.Status != user.StatusVerified {
		u.Status = user.StatusVerified
		err := s.db.Omit(clause.Associations).Save(u).Error
		if err != nil {
			s.logger.Error2("error saving user", err,
				zap.String("service", "authService"),
				zap.String("method", "VerifyEmail"),
				zap.Int("userID", u.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
	}

	return response.Success(nil, "")
}

func (s *service) GetTokenByRefreshToken(params GetTokenByRefreshTokenParams) (resp interface{}, statusCode int) {
	refreshToken := strings.Trim(params.Refresh, "")
	statusCode = http.StatusUnauthorized
	u := &user.User{}
	err := s.userRepository.GetUserByRefreshToken(refreshToken, u)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("error getting user", err,
			zap.String("service", "authService"),
			zap.String("method", "GetTokenByRefreshToken"),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || u.ID == 0 {
		return response.Error("user not found", http.StatusUnprocessableEntity, nil)
	}

	if u.AccountStatus == user.AccountStatusBlocked {
		return LoginFailedResponseData{Code: statusCode, Message: "account is blocked"}, statusCode
	}
	if u.RefreshTokenExpiry.Before(time.Now()) {
		return response.Error("refresh token expired", http.StatusUnprocessableEntity, nil)
	}
	token, err := s.jwtService.IssueToken(u.Email, params.UserAgent, params.IP)
	if err != nil {
		s.logger.Error2("error issuing token", err,
			zap.String("service", "authService"),
			zap.String("method", "GetTokenByRefreshToken"),
			zap.String("email", u.Email),
			zap.Int("userID", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	resp = LoginResponseData{
		Token:         token,
		RefreshToken:  refreshToken,
		Need2fa:       false,
		NeedEmailCode: false,
		IsNewDevice:   false,
	}

	return resp, http.StatusOK
}

func (s *service) GetUser(token string) (*user.User, error) {
	//todo maybe we should sent better errors to inform the user what was the problem of authentication
	username, issuedAt, err := s.jwtService.GetUsernameFromToken(token)
	if err != nil {
		return nil, err
	}
	loggedInUser := &user.User{}
	err = s.userRepository.GetUserByUsername(username, loggedInUser)

	if err != nil || loggedInUser.ID == 0 {
		return nil, err
	}
	//here it means the user has changed the password after the issuance of this token so we cosider it as invalid
	if loggedInUser.PasswordChangedAt.Valid && loggedInUser.PasswordChangedAt.Time.Sub(issuedAt) > 0 {
		return nil, fmt.Errorf("invalid token")
	}
	//here it means the user has enabled twoFa after the issuance of this token so we cosider it as invalid
	if loggedInUser.TwoFaChangedAt.Valid && loggedInUser.TwoFaChangedAt.Time.Sub(issuedAt) > 0 {
		return nil, fmt.Errorf("invalid token")
	}
	return loggedInUser, nil
}

func (s *service) GetAdminUser(token string) (*user.User, error) {
	username, issuedAt, err := s.jwtService.GetUsernameFromToken(token)
	if err != nil {
		return nil, err
	}
	loggedInUser := &user.User{}
	err = s.userRepository.GetAdminUserByUsername(username, loggedInUser)

	if err != nil || loggedInUser.ID == 0 {
		return nil, err
	}
	//here it means the user has changed the password after the issuance of this token so we cosider it as invalid
	if loggedInUser.PasswordChangedAt.Valid && loggedInUser.PasswordChangedAt.Time.After(issuedAt) {
		return nil, fmt.Errorf("invalid token")
	}
	return loggedInUser, nil
}

func (s *service) generateVerificationCodeForUser() string {
	uuidCode := uuid.NewString()
	sha := sha1.New()
	_, _ = io.WriteString(sha, uuidCode)
	shaResult := hex.EncodeToString(sha.Sum(nil))

	h := md5.New()
	_, _ = io.WriteString(h, shaResult)
	result := hex.EncodeToString(h.Sum(nil))
	return result
}

func (s *service) generatePrivateChannelNameForUser() string {
	uuidCode := uuid.NewString()
	sha := sha1.New()
	_, _ = io.WriteString(sha, uuidCode)
	shaResult := hex.EncodeToString(sha.Sum(nil))

	h := md5.New()
	_, _ = io.WriteString(h, shaResult)
	result := hex.EncodeToString(h.Sum(nil))
	return result
}

func (s *service) getUbID(u user.User) string {
	//todo this method is not a proper solution to generate the ubId try to change this later
	add1 := 11121
	add2 := 146
	coefficient := 10
	ubID := ((u.ID + add1) * coefficient) + add2
	return "ub" + strconv.Itoa(ubID)
}

func (s *service) sendVerificationEmail(u user.User) {
	domain := s.configs.GetDomain()
	link := ""
	if strings.HasSuffix(domain, "/") {
		link = domain + "auth/verify?code=" + u.VerificationCode
	} else {
		link = domain + "/auth/verify?code=" + u.VerificationCode
	}
	cu := communication.CommunicatingUser{
		Email: u.Email,
		Phone: "",
	}
	s.communicationService.SendVerificationEmailToUser(cu, link)
}

func (s *service) randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func NewAuthService(db *gorm.DB, userRepository user.Repository, userLevelService user.LevelService,
	userPermissionManager user.PermissionManager, userBalanceService userbalance.Service, jts jwt.Service,
	pe platform.PasswordEncoder, communicationService communication.Service, eventsHandler EventsHandler,
	forgotPasswordManager user.ForgotPasswordManager, recaptchaManager user.RecaptchaManager,
	twoFaManager user.TwoFaManager, configs platform.Configs, logger platform.Logger) Service {
	return &service{
		db:                    db,
		userRepository:        userRepository,
		userLevelService:      userLevelService,
		userPermissionManager: userPermissionManager,
		userBalanceService:    userBalanceService,
		jwtService:            jts,
		passwordEncoder:       pe,
		communicationService:  communicationService,
		eventsHandler:         eventsHandler,
		forgotPasswordManager: forgotPasswordManager,
		recaptchaManager:      recaptchaManager,
		twoFaManager:          twoFaManager,
		configs:               configs,
		logger:                logger,
	}
}
