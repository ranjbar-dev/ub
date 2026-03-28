package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"exchange-go/internal/api"
	"exchange-go/internal/auth"
	"exchange-go/internal/currency"
	"exchange-go/internal/di"
	"exchange-go/internal/platform"
	"exchange-go/internal/response"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AuthTests struct {
	*suite.Suite
	httpServer  http.Handler
	db          *gorm.DB
	redisClient *redis.Client
	userActor   *userActor
}

func (t *AuthTests) SetupSuite() {
	container := getContainer()
	t.httpServer = container.Get(di.HTTPServer).(api.HTTPServer).GetEngine()
	t.db = getDb()
	t.redisClient = getRedis()
	t.userActor = getUserActor()
}

func (t *AuthTests) SetupTest() {
}

func (t *AuthTests) TearDownTest() {
}

func (t *AuthTests) TearDownSuite() {

}

func (t *AuthTests) TestLogin_2faDisabled_Failed_WrongPassword() {
	newUserActor := getNewUserActor()
	data := `{"username":"` + newUserActor.Email + `","password":"1234567890","recaptcha":"recaptcha"}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))

	t.httpServer.ServeHTTP(res, req)
	result := auth.LoginFailedResponseData{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnauthorized, res.Code)
	assert.Equal(t.T(), "invalid credentials", result.Message)
	assert.Equal(t.T(), http.StatusUnauthorized, result.Code)

	time.Sleep(300 * time.Millisecond)

	lh := &user.LoginHistory{}
	err = t.db.Where("email = ?", newUserActor.Email).Order("id desc").First(lh).Error
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), int64(newUserActor.ID), lh.UserID.Int64)
	assert.Equal(t.T(), "FAILED", lh.Type)
	assert.Equal(t.T(), "WEB", lh.Device.String)
}

func (t *AuthTests) TestLogin_2faDisabled_Successful() {
	newUserActor := getNewUserActor()
	data := `{"username":"` + newUserActor.Email + `","password":"123456789","recaptcha":"recaptcha"}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))

	t.httpServer.ServeHTTP(res, req)
	result := auth.LoginResponseData{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), false, result.Need2fa)
	assert.Equal(t.T(), false, result.NeedEmailCode)
	assert.NotEqual(t.T(), "", result.Token)
	assert.NotEqual(t.T(), "", result.RefreshToken)

	//check if refresh token is inserted in db
	u := &user.User{}
	err = t.db.Where(user.User{Email: newUserActor.Email}).First(u).Error
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), true, u.RefreshToken.Valid)
	assert.NotEqual(t.T(), "", u.RefreshToken.String)
	if !u.RefreshTokenExpiry.After(time.Now()) {
		t.Fail(err.Error())
	}

	time.Sleep(300 * time.Millisecond)

	lh := &user.LoginHistory{}
	err = t.db.Where("email = ?", newUserActor.Email).Order("id desc").First(lh).Error
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), int64(newUserActor.ID), lh.UserID.Int64)
	assert.Equal(t.T(), "SUCCESSFUL", lh.Type)
	assert.Equal(t.T(), "WEB", lh.Device.String)
}

func (t *AuthTests) TestLogin_2faEnabled_Successful() {
	//to be sure that twoFa is enabled for user
	currentUser := &user.User{}
	err := t.db.Where(user.User{Email: t.userActor.Email}).First(currentUser).Error
	if err != nil {
		t.Fail(err.Error())
	}
	currentUser.IsTwoFaEnabled = true
	currentUser.Google2faDisabledAt = sql.NullTime{Time: time.Now().Add(-2 * time.Hour), Valid: true}
	currentUser.Google2faSecretCode = sql.NullString{String: "HWOAQZBGXCKJZQVH", Valid: true}
	err = t.db.Save(currentUser).Error
	if err != nil {
		t.Fail(err.Error())
	}

	//first we login without 2fa code
	data := `{"username":"` + t.userActor.Email + `","password":"123456789","recaptcha":"recaptcha"}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))

	t.httpServer.ServeHTTP(res, req)
	result := auth.LoginResponseData{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), true, result.Need2fa)
	assert.Equal(t.T(), false, result.NeedEmailCode)
	assert.Equal(t.T(), "", result.Token)
	assert.Equal(t.T(), "", result.RefreshToken)

	//here we generate 2fa code and login again
	twoFaCode, err := totp.GenerateCode("HWOAQZBGXCKJZQVH", time.Now())
	if err != nil {
		t.Fail(err.Error())
	}

	data = `{"username":"` + t.userActor.Email + `","password":"123456789","recaptcha":"recaptcha","2fa_code":"` + twoFaCode + `"}`
	res = httptest.NewRecorder()
	body = []byte(data)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))

	t.httpServer.ServeHTTP(res, req)
	result = auth.LoginResponseData{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), true, result.Need2fa)
	assert.Equal(t.T(), false, result.NeedEmailCode)
	assert.NotEqual(t.T(), "", result.Token)
	assert.NotEqual(t.T(), "", result.RefreshToken)

	//check if refresh token is inserted in db
	u := &user.User{}
	err = t.db.Where(user.User{Email: t.userActor.Email}).First(u).Error
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), true, u.RefreshToken.Valid)
	assert.NotEqual(t.T(), "", u.RefreshToken.String)
	if !u.RefreshTokenExpiry.After(time.Now()) {
		t.Fail(err.Error())
	}

	time.Sleep(300 * time.Millisecond)

	lh := &user.LoginHistory{}
	err = t.db.Where("email = ?", t.userActor.Email).Order("id desc").First(lh).Error
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), int64(t.userActor.ID), lh.UserID.Int64)
	assert.Equal(t.T(), "SUCCESSFUL", lh.Type)
	assert.Equal(t.T(), "WEB", lh.Device.String)
}

type registerValidationFailedScenarios struct {
	data         string
	reason       string
	errorMessage string
}

func (t *AuthTests) TestRegister_ValidationFail() {
	failedScenarios := []registerValidationFailedScenarios{
		{
			data:         `{"email":"","password":"123456789","recaptcha":"recaptcha"}`,
			reason:       "email not provided",
			errorMessage: "email is required",
		},
		{
			data:         `{"email":"test@test.com","password":"","recaptcha":"recaptcha"}`,
			reason:       "password not provided",
			errorMessage: "password is required",
		},
		{
			data:         `{"email":"test@test.com","password":"123456789","recaptcha":""}`,
			reason:       "recaptcha not provided",
			errorMessage: "recaptcha is required",
		},
	}

	for _, item := range failedScenarios {
		res := httptest.NewRecorder()
		body := []byte(item.data)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
		t.httpServer.ServeHTTP(res, req)
		result := response.APIResponse{}
		err := json.Unmarshal(res.Body.Bytes(), &result)
		if err != nil {
			t.Fail(err.Error())
		}

		assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
		assert.Equal(t.T(), item.errorMessage, result.Message)
	}

}

func (t *AuthTests) TestRegister_EmailAlreadyTaken() {
	formerUser := &user.User{
		Email:                  "alreadyTaken@test.com",
		Password:               "somepaswword",
		Kyc:                    0,
		Status:                 "REGISTERED",
		AccountStatus:          "UNBLOCKED",
		ExchangeNumber:         0,
		IsTwoFaEnabled:         false,
		UbID:                   "alreadyTakenUbId",
		VerificationCode:       "alreadyTakenVerificationCode",
		PrivateChannelName:     "privateChannelNameCode",
		ExchangeVolumeCoinCode: "",
		ExchangeVolumeAmount:   "",
		UserLevelID:            1,
	}

	err := t.db.Create(formerUser).Error
	if err != nil {
		t.Fail(err.Error())
	}

	data := `{"email":"alreadyTaken@test.com","password":"123456789","recaptcha":"recaptcha"}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)
	assert.Equal(t.T(), "email is already taken", result.Message)

}

func (t *AuthTests) TestRegister_Successful() {
	data := `{"email":"successfulRegister@test.com","password":"123456789","recaptcha":"recaptcha"}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), true, result.Status)
	assert.Equal(t.T(), "", result.Message)

	//sleeping to be sure all goroutines are done and we could check if userbalances inserted too
	time.Sleep(3 * time.Second)

	createdUser := &user.User{}
	err = t.db.Where(user.User{Email: "successfulRegister@test.com"}).First(createdUser).Error
	assert.Nil(t.T(), err)
	assert.NotEqual(t.T(), int64(0), createdUser.ID)
	assert.NotEqual(t.T(), "", createdUser.UbID)
	assert.NotEqual(t.T(), "", createdUser.PrivateChannelName)
	assert.NotEqual(t.T(), "", createdUser.VerificationCode)
	assert.Equal(t.T(), user.StatusRegistered, createdUser.Status)
	assert.Equal(t.T(), user.AccountStatusUnblocked, createdUser.AccountStatus)

	pe := platform.NewPasswordEncoder()
	err = pe.CompareHashAndPassword(createdUser.Password, "123456789")
	assert.Nil(t.T(), err)

	createdUserProfile := &user.Profile{}
	err = t.db.Where(user.Profile{UserID: createdUser.ID}).First(createdUserProfile).Error
	assert.Nil(t.T(), err)
	assert.NotEqual(t.T(), "", createdUserProfile.RegistrationIP.String)

	var ups []user.UsersPermissions
	t.db.Joins("Permission").Where("user_id = ?", createdUser.ID).Find(&ups)
	assert.Equal(t.T(), 5, len(ups))

	var coins []currency.Coin
	t.db.Where("is_active = ?", true).Find(&coins)
	var userBalances []userbalance.UserBalance
	t.db.Where("user_id = ?", createdUser.ID).Find(&userBalances)
	assert.Equal(t.T(), 6, len(userBalances))
	for _, ub := range userBalances {
		for _, c := range coins {
			if ub.CoinID == c.ID {
				if c.BlockchainNetwork.String != "" {
					assert.Equal(t.T(), c.BlockchainNetwork.String+"Address", ub.Address.String)
				} else {
					assert.Equal(t.T(), c.Code+"Address", ub.Address.String)
				}
			}
			assert.Equal(t.T(), "0.0", ub.Amount)
			assert.Equal(t.T(), "0.0", ub.FrozenAmount)

		}
	}

}

func (t *AuthTests) TestForgotPasswordAndUpdate() {
	newUserActor := getNewUserActor()
	data := `{"email":"` + newUserActor.Email + `","recaptcha":"recaptcha"}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/forgot-password", bytes.NewReader(body))

	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), true, result.Status)

	//get code from redis to restart the password
	userIDString := strconv.Itoa(newUserActor.ID)
	redisData, err := t.redisClient.HGetAll(context.Background(), "forgot-password:"+userIDString).Result()
	if err != nil {
		t.Fail(err.Error())
	}
	forgotPasswordCode, ok := redisData["code"]
	if !ok {
		t.Fail("forgot password key does not exist in redis")
	}

	//calling the forgot-password-update
	data = `{"email":"` + newUserActor.Email + `","password":"987654321","confirmed":"987654321","code":"` + forgotPasswordCode + `"}`
	res = httptest.NewRecorder()
	body = []byte(data)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/forgot-password/update", bytes.NewReader(body))

	t.httpServer.ServeHTTP(res, req)
	result = response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), true, result.Status)

	updatedUser := &user.User{}
	err = t.db.Where(user.User{ID: newUserActor.ID}).First(updatedUser).Error
	assert.Nil(t.T(), err)
	assert.True(t.T(), updatedUser.PasswordChangedAt.Valid)

	time.Sleep(300 * time.Millisecond)
	count, err := t.redisClient.Exists(context.Background(), "forgot-password:"+userIDString).Result()
	if err != nil && err != redis.Nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), int64(0), count)

	//finaly we login again to check if the password is changed
	data = `{"username":"` + newUserActor.Email + `","password":"987654321","recaptcha":"recaptcha"}`
	res = httptest.NewRecorder()
	body = []byte(data)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))

	t.httpServer.ServeHTTP(res, req)
	loginResult := auth.LoginResponseData{}
	err = json.Unmarshal(res.Body.Bytes(), &loginResult)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), false, loginResult.Need2fa)
	assert.Equal(t.T(), false, loginResult.NeedEmailCode)
	assert.NotEqual(t.T(), "", loginResult.Token)
	assert.NotEqual(t.T(), "", loginResult.RefreshToken)
}

func (t *AuthTests) TestVerifyEmail() {
	passwordEncoder := platform.NewPasswordEncoder()
	encodedPassword, _ := passwordEncoder.GenerateFromPassword("123456789")
	email := uuid.NewString() + "@test.com"
	ubID := randSeq(10)
	verificationCode := randSeq(12)
	PrivateChannelName := randSeq(14)
	u := user.User{
		Email:              email,
		Password:           string(encodedPassword),
		Kyc:                user.KycLevelMinimum,
		Status:             user.StatusRegistered,
		AccountStatus:      "UNBLOCKED",
		ExchangeNumber:     12,
		IsTwoFaEnabled:     false,
		UbID:               ubID,
		VerificationCode:   verificationCode,
		PrivateChannelName: PrivateChannelName,
		UserLevelID:        1,
	}
	err := t.db.Create(&u).Error
	if err != nil {
		t.Fail(err.Error())
	}

	data := `{"code":"` + verificationCode + `"}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/verify", bytes.NewReader(body))

	t.httpServer.ServeHTTP(res, req)
	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), true, result.Status)

	//check if user is verified in database
	updatedUser := &user.User{}
	err = t.db.Where(user.User{Email: email}).First(updatedUser).Error
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), user.StatusVerified, updatedUser.Status)
}

/*
func (t *AuthTests) TestGetTokenByRefreshToken() {
	passwordEncoder := platform.NewPasswordEncoder()
	encodedPassword, _ := passwordEncoder.GenerateFromPassword("123456789")
	email := uuid.NewString() + "@test.com"
	ubId := randSeq(10)
	verificationCode := randSeq(12)
	PrivateChannelName := randSeq(14)
	u := user.User{
		Email:              email,
		Password:           string(encodedPassword),
		Kyc:                user.KycLevelMinimum,
		Status:             user.StatusRegistered,
		AccountStatus:      user.StatusVerified,
		ExchangeNumber:     12,
		IsTwoFaEnabled:     false,
		UbId:               ubId,
		VerificationCode:   verificationCode,
		PrivateChannelName: PrivateChannelName,
		UserLevelId:        1,
		RefreshToken:       sql.NullString{String: "someUniqueRefreshToken", Valid: true},
		RefreshTokenExpiry: time.Now().Add(3 * time.Hour),
	}
	err := t.db.Create(&u).Error
	if err != nil {
		t.Fail(err.Error())
	}

	data := `{"refresh":"someUniqueRefreshToken"}`
	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(body))

	t.httpServer.ServeHTTP(res, req)
	result := auth.LoginResponseData{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), http.StatusOK, res.Code)
	assert.Equal(t.T(), false, result.Need2fa)
	assert.Equal(t.T(), false, result.NeedEmailCode)
	assert.NotEqual(t.T(), "", result.Token)
	assert.NotEqual(t.T(), "", result.RefreshToken)
}
*/

func TestAuth(t *testing.T) {
	suite.Run(t, &AuthTests{
		Suite: new(suite.Suite),
	})
}
