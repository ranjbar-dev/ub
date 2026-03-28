package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"exchange-go/internal/api"
	"exchange-go/internal/country"
	"exchange-go/internal/di"
	"exchange-go/internal/platform"
	"exchange-go/internal/user"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserTests struct {
	*suite.Suite
	httpServer  http.Handler
	db          *gorm.DB
	redisClient *redis.Client
	userActor   *userActor
}

func (t *UserTests) SetupSuite() {
	container := getContainer()
	t.httpServer = container.Get(di.HTTPServer).(api.HTTPServer).GetEngine()
	t.db = getDb()
	t.redisClient = getRedis()
	t.userActor = getUserActor()

}

func (t *UserTests) SetupTest() {

}

func (t *UserTests) TearDownTest() {
}

func (t *UserTests) TearDownSuite() {

}

func (t *UserTests) Test_AA_SetUserProfile() {
	c := country.Country{
		ID:        1,
		Name:      sql.NullString{String: "co1", Valid: true},
		FullName:  sql.NullString{String: "country1", Valid: true},
		Code:      sql.NullString{String: "01", Valid: true},
		ImagePath: sql.NullString{String: "/images", Valid: true},
	}

	err := t.db.Create(&c).Error
	if err != nil {
		t.Fail(err.Error())
	}

	data := `{` +
		`"first_name":"firstNameTest",` +
		`"last_name":"lastNameTest",` +
		`"gender":"male",` +
		`"date_of_birth":"2002-01-02",` +
		`"address":"testAddress",` +
		`"region_and_city":"test and test",` +
		`"postal_code":"testPostalCode",` +
		`"country_id":1` +
		`}`

	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/set-user-profile", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    user.SetUserProfileResponse
	}{}

	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	assert.Equal(t.T(), "firstNameTest", result.Data.FirstName)
	assert.Equal(t.T(), "lastNameTest", result.Data.LastName)
	assert.Equal(t.T(), "male", result.Data.Gender)
	assert.Equal(t.T(), "2002-01-02", result.Data.DateOfBirth)
	assert.Equal(t.T(), "testAddress", result.Data.Address)
	assert.Equal(t.T(), "test and test", result.Data.RegionAndCity)
	assert.Equal(t.T(), "testPostalCode", result.Data.PostalCode)
	assert.Equal(t.T(), int64(1), result.Data.CountryID)
	assert.Equal(t.T(), "co1", result.Data.CountryName)
	assert.Equal(t.T(), "processing", result.Data.Status)

	//check if the user profile is in database
	up := &user.Profile{}
	err = t.db.Where(user.Profile{UserID: t.userActor.ID}).Find(up).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "firstNameTest", up.FirstName.String)
	assert.Equal(t.T(), "lastNameTest", up.LastName.String)
	assert.Equal(t.T(), "MALE", up.Gender.String)
	assert.Equal(t.T(), "2002-01-02", up.DateOfBirth.String)
	assert.Equal(t.T(), "testAddress", up.Address.String)
	assert.Equal(t.T(), "test and test", up.RegionAndCity.String)
	assert.Equal(t.T(), "testPostalCode", up.PostalCode.String)
	assert.Equal(t.T(), int64(1), up.CountryID.Int64)
	assert.Equal(t.T(), "PROCESSING", up.Status.String)

}

func (t *UserTests) Test_AB_GetProfile() {
	up := &user.Profile{}
	err := t.db.Where(user.Profile{UserID: t.userActor.ID}).Find(up).Error
	if err != nil {
		t.Fail(err.Error())
	}

	//insert userProfileImages
	images := []user.ProfileImage{
		{
			ID:                 1,
			UserProfileID:      sql.NullInt64{Int64: up.ID, Valid: true},
			Type:               user.ProfileImageTypeIdentity,
			ImagePath:          "/image1.jpg",
			ConfirmationStatus: sql.NullString{String: user.ProfileImageStatusProcessing, Valid: true},
			OriginalFileName:   "image1.jpg",
			IDCardCode:         sql.NullString{String: "123456", Valid: true},
			SubType:            sql.NullString{String: user.ProfileImageSubtypeIdentityIdentityCard, Valid: true},
			IsBack:             sql.NullBool{Bool: false, Valid: true},
		},
		{
			ID:                 2,
			UserProfileID:      sql.NullInt64{Int64: up.ID, Valid: true},
			Type:               user.ProfileImageTypeIdentity,
			ImagePath:          "/image2.jpg",
			ConfirmationStatus: sql.NullString{String: user.ProfileImageStatusProcessing, Valid: true},
			OriginalFileName:   "image2.jpg",
			SubType:            sql.NullString{String: user.ProfileImageSubtypeIdentityIdentityCard, Valid: true},
			MainImageID:        sql.NullInt64{Int64: 1, Valid: true},
			IsBack:             sql.NullBool{Bool: true, Valid: true},
		},
		{
			ID:                 3,
			UserProfileID:      sql.NullInt64{Int64: up.ID, Valid: true},
			Type:               user.ProfileImageTypeAddress,
			ImagePath:          "/image3.jpg",
			ConfirmationStatus: sql.NullString{String: user.ProfileImageStatusProcessing, Valid: true},
			OriginalFileName:   "image3.jpg",
			SubType:            sql.NullString{String: user.ProfileImageSubtypeAddressUtilityBill, Valid: true},
			IsBack:             sql.NullBool{Bool: false, Valid: true},
		},
	}

	err = t.db.Create(images).Error
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/get-user-profile", nil)
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    user.GetProfileResponse
	}{}

	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	assert.Equal(t.T(), "firstNameTest", result.Data.FirstName)
	assert.Equal(t.T(), "lastNameTest", result.Data.LastName)
	assert.Equal(t.T(), "male", result.Data.Gender)
	assert.Equal(t.T(), "2002-01-02", result.Data.DateOfBirth)
	assert.Equal(t.T(), "testAddress", result.Data.Address)
	assert.Equal(t.T(), "test and test", result.Data.RegionAndCity)
	assert.Equal(t.T(), "testPostalCode", result.Data.PostalCode)
	assert.Equal(t.T(), int64(1), *result.Data.CountryID)
	assert.Equal(t.T(), "co1", *result.Data.CountryName)
	assert.Equal(t.T(), "processing", result.Data.Status)
	assert.Equal(t.T(), "", result.Data.AdminComment)
	//assert.Equal(t.T(), "", result.Data.UserProfileImages)
	for _, image := range result.Data.UserProfileImages {
		switch image.ID {
		case int64(1):
			assert.Equal(t.T(), "identity", image.Type)
			assert.Equal(t.T(), int64(1), image.ID)
			assert.Equal(t.T(), "123456", image.IDCardCode)
			assert.Equal(t.T(), "processing", image.Status)
			assert.Equal(t.T(), "", image.RejectionReason)
			assert.Equal(t.T(), int64(1), image.ImageID)
			assert.Nil(t.T(), image.MainImageID)
			assert.Equal(t.T(), false, image.IsBack)
			assert.Equal(t.T(), "identity_card", image.SubType)
		case int64(2):
			assert.Equal(t.T(), "identity", image.Type)
			assert.Equal(t.T(), int64(2), image.ID)
			assert.Equal(t.T(), "", image.IDCardCode)
			assert.Equal(t.T(), "processing", image.Status)
			assert.Equal(t.T(), "", image.RejectionReason)
			assert.Equal(t.T(), int64(2), image.ImageID)
			assert.Equal(t.T(), int64(1), *image.MainImageID)
			assert.Equal(t.T(), true, image.IsBack)
			assert.Equal(t.T(), "identity_card", image.SubType)

		case int64(3):
			assert.Equal(t.T(), "address", image.Type)
			assert.Equal(t.T(), int64(3), image.ID)
			assert.Equal(t.T(), "", image.IDCardCode)
			assert.Equal(t.T(), "processing", image.Status)
			assert.Equal(t.T(), "", image.RejectionReason)
			assert.Equal(t.T(), int64(3), image.ImageID)
			assert.Nil(t.T(), image.MainImageID)
			assert.Equal(t.T(), false, image.IsBack)
			assert.Equal(t.T(), "utility_bill", image.SubType)

		default:
			t.Fail("we should not be in default case")
		}
	}

}

func (t *UserTests) Test_AC_EnableSms() {
	data := `{` +
		`"phone":"+989121234567"` +
		`}`

	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/sms-send", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)
	//enable sms

	//get the code from redis
	//update user to need twoFa
	err := t.db.Model(&user.User{}).Where("id =?", t.userActor.ID).Updates(user.User{IsTwoFaEnabled: true}).Error
	if err != nil {
		t.Fail(err.Error())
	}

	userIDString := strconv.Itoa(t.userActor.ID)
	redisData, err := t.redisClient.HGetAll(context.Background(), "phone-confirmation:"+userIDString).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	smsCode := redisData["code"]

	twoFaCode, err := totp.GenerateCode("HWOAQZBGXCKJZQVH", time.Now()) //secret from main_test file
	if err != nil {
		t.Fail(err.Error())
	}

	data = `{` +
		`"phone":"+989121234567",` +
		`"code":"` + smsCode + `",` +
		`"2fa_code":"` + twoFaCode + `",` +
		`"password":""` +
		`}`

	res = httptest.NewRecorder()
	body = []byte(data)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/user/sms-enable", bytes.NewReader(body))
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	//because we have goroutine we sleep here to all the things be done
	time.Sleep(300 * time.Millisecond)
	//check the database
	u := &user.User{}
	err = t.db.Where(user.User{ID: t.userActor.ID}).Find(u).Error
	if err != nil {
		t.Fail(err.Error())
	}

	up := &user.Profile{}
	err = t.db.Where(user.Profile{UserID: t.userActor.ID}).Find(up).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "+989121234567", u.Phone.String)
	assert.Equal(t.T(), user.KycLevelMinimum, u.Kyc)
	assert.Equal(t.T(), "SMS", up.PhoneConfirmationType.String)

}

func (t *UserTests) Test_AD_GetUserData() {
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/user-data", nil)
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    user.GetUserDataResponse
	}{}

	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	assert.Equal(t.T(), "test@test.com", result.Data.Email) //from the main_test file
	assert.Equal(t.T(), "", result.Data.UbID)
	assert.Equal(t.T(), "+989121234567", result.Data.Phone)
	assert.Equal(t.T(), "none", result.Data.KycLevel)
	assert.Equal(t.T(), "processing", result.Data.KycStatus)
	assert.Equal(t.T(), "you identity,address, phone are not confirmed yet", result.Data.KycLevelMessage)
	assert.Equal(t.T(), "medium", result.Data.SecurityLevel)
	assert.Equal(t.T(), "we highly recommend to verify your identity", result.Data.SecurityLevelMessage)
	assert.Equal(t.T(), true, result.Data.Google2faEnabled)
	assert.Equal(t.T(), true, result.Data.Has2fa)
	assert.Equal(t.T(), true, result.Data.IsAccountVerified)
	assert.Equal(t.T(), "userActorPrivateChannel", result.Data.ChannelName)
	assert.Equal(t.T(), 0, result.Data.ThemeID)
	assert.Equal(t.T(), "default", result.Data.Theme)
	assert.Equal(t.T(), "processing", result.Data.ProfileStatus)
}

func (t *UserTests) Test_AE_DisableSms() {
	data := `{` +
		`"phone":"+989121234567"` +
		`}`

	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/sms-send", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	//enable sms

	//get the code from redis
	userIDString := strconv.Itoa(t.userActor.ID)
	redisData, err := t.redisClient.HGetAll(context.Background(), "phone-confirmation:"+userIDString).Result()
	if err != nil {
		t.Fail(err.Error())
	}

	smsCode := redisData["code"]

	twoFaCode, err := totp.GenerateCode("HWOAQZBGXCKJZQVH", time.Now()) //secret from main_test file
	if err != nil {
		t.Fail(err.Error())
	}

	data = `{` +
		`"phone":"+989121234567",` +
		`"code":"` + smsCode + `",` +
		`"2fa_code":"` + twoFaCode + `",` +
		`"password":""` +
		`}`

	res = httptest.NewRecorder()
	body = []byte(data)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/user/sms-disable", bytes.NewReader(body))
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	//because we have goroutine we sleep here to all the things be done
	time.Sleep(300 * time.Millisecond)
	//check the database
	u := &user.User{}
	err = t.db.Where(user.User{ID: t.userActor.ID}).Find(u).Error
	if err != nil {
		t.Fail(err.Error())
	}

	up := &user.Profile{}
	err = t.db.Where(user.Profile{UserID: t.userActor.ID}).Find(up).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), false, u.Phone.Valid)
	assert.Equal(t.T(), "", u.Phone.String)
	assert.Equal(t.T(), user.KycLevelMinimum, u.Kyc)

}

func (t *UserTests) TestGet2FaBarcodeThenEnable2FaThenDisable() {
	//first we get the secret code
	newUserActor := getNewUserActor()

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/google-2fa-barcode", nil)
	token := "Bearer " + newUserActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    user.Get2FaBarcodeResponse
	}{}

	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	assert.Equal(t.T(), 16, len(result.Data.Code))

	//check the database
	u := &user.User{}
	err = t.db.Where(user.User{ID: newUserActor.ID}).Find(u).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), result.Data.Code, u.Google2faSecretCode.String)

	//here we enable twoFa

	code, err := totp.GenerateCode(u.Google2faSecretCode.String, time.Now())
	if err != nil {
		t.Fail(err.Error())
	}

	data := `{` +
		`"password":"123456789",` +
		`"code":"` + code + `"` +
		`}`

	res = httptest.NewRecorder()
	body := []byte(data)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/user/google-2fa-enable", bytes.NewReader(body))
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	assert.Equal(t.T(), http.StatusOK, res.Code)
	enableResult := struct {
		Status  bool
		Message string
		Data    map[string]string
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &enableResult)
	if err != nil {
		t.Fail(err.Error())
	}
	newToken := enableResult.Data["token"]
	assert.NotEmpty(t.T(), newToken)

	updatedUser := &user.User{}
	err = t.db.Where(user.User{ID: newUserActor.ID}).Find(updatedUser).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), true, updatedUser.IsTwoFaEnabled)
	assert.Equal(t.T(), true, updatedUser.TwoFaChangedAt.Valid)

	//u.PasswordChangedAt = sql.NullTime{Time: time.Now(), Valid: false}
	//t.db.Save(u)

	//here we disable twoFa

	data = `{` +
		`"password":"123456789",` +
		`"code":"` + code + `"` +
		`}`

	res = httptest.NewRecorder()
	body = []byte(data)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/user/google-2fa-disable", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+newToken)
	t.httpServer.ServeHTTP(res, req)

	assert.Equal(t.T(), http.StatusOK, res.Code)

	updatedUser = &user.User{}
	err = t.db.Where(user.User{ID: newUserActor.ID}).Find(updatedUser).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), false, updatedUser.IsTwoFaEnabled)
	assert.Equal(t.T(), true, updatedUser.Google2faDisabledAt.Valid)
	assert.Equal(t.T(), false, updatedUser.Google2faSecretCode.Valid)

	//doing this so other tests could be passed
	u.TwoFaChangedAt = sql.NullTime{Time: time.Now(), Valid: false}
	t.db.Save(u)
}

func (t *UserTests) TestChangePassword_2FaIsEnabled() {
	code, err := totp.GenerateCode("HWOAQZBGXCKJZQVH", time.Now()) //secret from main_test file
	if err != nil {
		t.Fail(err.Error())
	}

	data := `{` +
		`"old_password":"123456789",` +
		`"new_password":"987654321",` +
		`"confirmed":"987654321",` +
		`"2fa_code":"` + code + `"` +
		`}`

	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/change-password", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)
	result := struct {
		Status  bool
		Message string
		Data    map[string]string
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.NotEmpty(t.T(), result.Data["token"])
	u := &user.User{}
	err = t.db.Where(user.User{ID: t.userActor.ID}).Find(u).Error
	if err != nil {
		t.Fail(err.Error())
	}

	passwordEncoder := platform.NewPasswordEncoder()
	err = passwordEncoder.CompareHashAndPassword(u.Password, "987654321")

	assert.Nil(t.T(), err)

	//here we set the passwordChangedAt to Null so other tests won't fail
	u.PasswordChangedAt = sql.NullTime{Time: time.Now(), Valid: false}
	t.db.Save(u)
}

func (t *UserTests) TestChangePassword_2FaIsDisabled() {
	newUserActor := getNewUserActor()

	data := `{` +
		`"old_password":"123456789",` +
		`"new_password":"987654321",` +
		`"confirmed":"987654321",` +
		`"2fa_code":""` +
		`}`

	res := httptest.NewRecorder()
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/change-password", bytes.NewReader(body))
	token := "Bearer " + newUserActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)
	result := struct {
		Status  bool
		Message string
		Data    map[string]string
	}{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.NotEmpty(t.T(), result.Data["token"])

	u := &user.User{}
	err = t.db.Where(user.User{ID: newUserActor.ID}).Find(u).Error
	if err != nil {
		t.Fail(err.Error())
	}

	passwordEncoder := platform.NewPasswordEncoder()
	err = passwordEncoder.CompareHashAndPassword(u.Password, "987654321")

	assert.Nil(t.T(), err)

	//here we set the passwordChangedAt to Null so other tests won't fail
	u.PasswordChangedAt = sql.NullTime{Time: time.Now(), Valid: false}
	t.db.Save(u)
}

func (t *UserTests) TestUploadImages_BackAndFrontProvided() {
	newUserActor := getNewUserActor()
	up := &user.Profile{
		UserID:     newUserActor.ID,
		TrustLevel: 0,
	}

	err := t.db.Create(&up).Error
	if err != nil {
		t.Fail(err.Error())
	}

	path := "./data/image.png"
	frontImage, err := os.Open(path)
	if err != nil {
		t.Fail(err.Error())
	}
	defer frontImage.Close()
	backImage, err := os.Open(path)
	if err != nil {
		t.Fail(err.Error())
	}
	defer backImage.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	frontPart, err := writer.CreateFormFile("front_image", filepath.Base(path))
	if err != nil {
		t.Fail(err.Error())
	}
	_, err = io.Copy(frontPart, frontImage)
	if err != nil {
		t.Fail(err.Error())
	}
	backPart, err := writer.CreateFormFile("back_image", filepath.Base(path))
	if err != nil {
		t.Fail(err.Error())
	}
	_, err = io.Copy(backPart, backImage)
	if err != nil {
		t.Fail(err.Error())
	}
	writer.WriteField("type", "identity")
	writer.WriteField("sub_type", "passport")
	writer.WriteField("id_card_code", "123456789")
	err = writer.Close()
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user-profile-image/multiple-upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	token := "Bearer " + newUserActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)
	result := struct {
		Status  bool
		Message string
		Data    map[string]user.UploadImageSingleResponse
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	frontImageID := result.Data["frontImage"].ID
	backImageID := result.Data["backImage"].ID
	assert.NotEmpty(t.T(), frontImageID)
	assert.NotEmpty(t.T(), result.Data["frontImage"].Path)

	assert.NotEmpty(t.T(), backImageID)
	assert.NotEmpty(t.T(), result.Data["backImage"].Path)

	updatedProfile := &user.Profile{}
	err = t.db.Where("id = ? ", up.ID).First(updatedProfile).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "PROCESSING", updatedProfile.IdentityConfirmationStatus.String)

	var images []user.ProfileImage
	err = t.db.Where("user_profile_id = ? ", up.ID).Find(&images).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), 2, len(images))
	for _, img := range images {
		if img.ID == frontImageID {
			assert.Equal(t.T(), "PROCESSING", img.ConfirmationStatus.String)
			assert.Equal(t.T(), "123456789", img.IDCardCode.String)
			assert.Equal(t.T(), "IDENTITY", img.Type)
			assert.Equal(t.T(), "PASSPORT", img.SubType.String)
			assert.Equal(t.T(), true, img.IsBack.Valid)
			assert.Equal(t.T(), false, img.IsBack.Bool)
		}

		if img.ID == backImageID {
			assert.Equal(t.T(), "PROCESSING", img.ConfirmationStatus.String)
			assert.Equal(t.T(), "123456789", img.IDCardCode.String)
			assert.Equal(t.T(), "IDENTITY", img.Type)
			assert.Equal(t.T(), "PASSPORT", img.SubType.String)
			assert.Equal(t.T(), true, img.IsBack.Valid)
			assert.Equal(t.T(), true, img.IsBack.Bool)
			assert.Equal(t.T(), frontImageID, img.MainImageID.Int64)
		}
	}

	files, _ := os.ReadDir("." + user.UserProfileImagePathPrefix)
	assert.Equal(t.T(), 2, len(files))

	// delete image folder in tests
	err = os.RemoveAll("./assets")
	if err != nil {
		t.Fail(err.Error())
	}

}

func (t *UserTests) TestUploadImages_OnlyBackProvided() {
	newUserActor := getNewUserActor()
	up := &user.Profile{
		UserID:     newUserActor.ID,
		TrustLevel: 0,
	}

	err := t.db.Create(&up).Error
	if err != nil {
		t.Fail(err.Error())
	}

	upi := &user.ProfileImage{
		UserProfileID:      sql.NullInt64{Int64: up.ID, Valid: true},
		Type:               "ADDRESS",
		ImagePath:          "/somepath",
		ConfirmationStatus: sql.NullString{String: "PROCESSING", Valid: true},
		OriginalFileName:   "somename.png",
		SubType:            sql.NullString{String: "PASSPORT", Valid: true},
		IsBack:             sql.NullBool{Bool: false, Valid: true},
	}

	err = t.db.Create(&upi).Error
	if err != nil {
		t.Fail(err.Error())
	}

	path := "./data/image.png"
	frontImage, err := os.Open(path)
	if err != nil {
		t.Fail(err.Error())
	}
	defer frontImage.Close()
	backImage, err := os.Open(path)
	if err != nil {
		t.Fail(err.Error())
	}
	defer backImage.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	backPart, err := writer.CreateFormFile("back_image", filepath.Base(path))
	if err != nil {
		t.Fail(err.Error())
	}
	_, err = io.Copy(backPart, backImage)
	if err != nil {
		t.Fail(err.Error())
	}
	frontImageIDString := strconv.FormatInt(upi.ID, 10)
	writer.WriteField("type", "identity")
	writer.WriteField("sub_type", "passport")
	writer.WriteField("id_card_code", "123456789")
	writer.WriteField("front_image_id", frontImageIDString)
	err = writer.Close()
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user-profile-image/multiple-upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	token := "Bearer " + newUserActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)
	result := struct {
		Status  bool
		Message string
		Data    map[string]user.UploadImageSingleResponse
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	backImageID := result.Data["backImage"].ID

	assert.NotEmpty(t.T(), backImageID)
	assert.NotEmpty(t.T(), result.Data["backImage"].Path)

	updatedProfile := &user.Profile{}
	err = t.db.Where("id = ? ", up.ID).First(updatedProfile).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "PROCESSING", updatedProfile.IdentityConfirmationStatus.String)

	var images []user.ProfileImage
	err = t.db.Where("user_profile_id = ? ", up.ID).Find(&images).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), 2, len(images))
	for _, img := range images {
		if img.ID == backImageID {
			assert.Equal(t.T(), "PROCESSING", img.ConfirmationStatus.String)
			assert.Equal(t.T(), "123456789", img.IDCardCode.String)
			assert.Equal(t.T(), "IDENTITY", img.Type)
			assert.Equal(t.T(), "PASSPORT", img.SubType.String)
			assert.Equal(t.T(), true, img.IsBack.Valid)
			assert.Equal(t.T(), true, img.IsBack.Bool)
			assert.Equal(t.T(), upi.ID, img.MainImageID.Int64)
		}
	}

	files, _ := os.ReadDir("." + user.UserProfileImagePathPrefix)
	assert.Equal(t.T(), 1, len(files))

	// delete image folder in tests
	err = os.RemoveAll("./assets")
	if err != nil {
		t.Fail(err.Error())
	}

}

func (t *UserTests) TestDeleteImage() {
	newUserActor := getNewUserActor()
	up := &user.Profile{
		UserID:     newUserActor.ID,
		TrustLevel: 0,
	}

	err := t.db.Create(&up).Error
	if err != nil {
		t.Fail(err.Error())
	}

	upi := &user.ProfileImage{
		UserProfileID:      sql.NullInt64{Int64: up.ID, Valid: true},
		Type:               "ADDRESS",
		ImagePath:          "/deletingpath.png",
		ConfirmationStatus: sql.NullString{String: "PROCESSING", Valid: true},
		OriginalFileName:   "somename.png",
		SubType:            sql.NullString{String: "PASSPORT", Valid: true},
		IsBack:             sql.NullBool{Bool: false, Valid: true},
	}

	err = t.db.Create(&upi).Error
	if err != nil {
		t.Fail(err.Error())
	}

	data := fmt.Sprintf(`{"id":%d}`, upi.ID)
	body := []byte(data)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user-profile-image/delete", bytes.NewReader(body))
	token := "Bearer " + newUserActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	updatedProfileImage := &user.ProfileImage{}
	err = t.db.Where("id = ? ", upi.ID).First(updatedProfileImage).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), true, updatedProfileImage.IsDeleted.Valid)
	assert.Equal(t.T(), true, updatedProfileImage.IsDeleted.Bool)

}

func TestUser(t *testing.T) {
	suite.Run(t, &UserTests{
		Suite: new(suite.Suite),
	})
}
