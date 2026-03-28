// Package user_test tests the user service. Covers:
//   - GetUserProfile and SetUserProfile (country not found, already confirmed, success)
//   - GetProfile and GetUserData aggregation
//   - 2FA barcode generation, enabling, and disabling
//   - Password change with 2FA enabled, 2FA disabled, and wrong current password
//   - SMS phone confirmation: rate limiting, invalid format, successful send
//   - EnableSms and DisableSms with wrong code, 2FA enabled, and 2FA disabled
//   - Sending verification emails and deleting profile images
//
// Test data: testify mocks for UserRepository, UserProfileRepository,
// ProfileImageRepository, CountryService, TwoFaManager,
// PhoneConfirmationManager, JwtService, CommunicationService, and
// go-sqlmock for GORM database interactions.
package user_test

import (
	"database/sql"
	"exchange-go/internal/country"
	"exchange-go/internal/mocks"
	"exchange-go/internal/platform"
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

func TestService_GetUserProfile(t *testing.T) {
	db := &gorm.DB{}
	userRepo := new(mocks.UserRepository)
	userProfileRepo := new(mocks.UserProfileRepository)
	userProfileRepo.On("GetProfileByUserID", 1, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		profile := args.Get(1).(*user.Profile)
		profile.ID = 1
	})

	profileImageRepository := new(mocks.ProfileImageRepository)
	countryService := new(mocks.CountryService)
	twoFaManager := new(mocks.TwoFaManager)
	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	jwtSerivce := new(mocks.JwtService)
	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	u := user.User{
		ID: 1,
	}
	profile, err := userService.GetUserProfile(u)

	assert.Nil(t, err)
	assert.Equal(t, int64(1), profile.ID)

	userProfileRepo.AssertExpectations(t)
}

func TestService_SetUserProfile_CountryNotFound(t *testing.T) {
	db := &gorm.DB{}
	userRepo := new(mocks.UserRepository)
	userProfileRepo := new(mocks.UserProfileRepository)

	profileImageRepository := new(mocks.ProfileImageRepository)
	countryService := new(mocks.CountryService)
	countryService.On("GetCountryByID", int64(1)).Once().Return(country.Country{}, gorm.ErrRecordNotFound)
	twoFaManager := new(mocks.TwoFaManager)
	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	jwtSerivce := new(mocks.JwtService)
	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	u := user.User{
		ID: 1,
	}
	params := user.SetUserProfileParams{
		FirstName:     "test",
		LastName:      "test",
		Gender:        "male",
		DateOfBirth:   "2007-06-06",
		Address:       "test",
		RegionAndCity: "test and test",
		PostalCode:    "test",
		CountryID:     1,
	}

	res, statusCode := userService.SetUserProfile(&u, params)

	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "country not found", res.Message)

	countryService.AssertExpectations(t)
}

func TestService_SetUserProfile_AlreadyConfirmed(t *testing.T) {
	db := &gorm.DB{}
	userRepo := new(mocks.UserRepository)
	userProfileRepo := new(mocks.UserProfileRepository)
	userProfileRepo.On("GetProfileByUserID", 1, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		profile := args.Get(1).(*user.Profile)
		profile.ID = 1
		profile.Status = sql.NullString{String: user.ProfileStatusConfirmed, Valid: true}
	})

	profileImageRepository := new(mocks.ProfileImageRepository)
	countryService := new(mocks.CountryService)
	countryService.On("GetCountryByID", int64(1)).Once().Return(country.Country{ID: 1}, nil)
	twoFaManager := new(mocks.TwoFaManager)
	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	jwtSerivce := new(mocks.JwtService)
	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	u := user.User{
		ID: 1,
	}
	params := user.SetUserProfileParams{
		FirstName:     "test",
		LastName:      "test",
		Gender:        "male",
		DateOfBirth:   "2007-06-06",
		Address:       "test",
		RegionAndCity: "test and test",
		PostalCode:    "test",
		CountryID:     1,
	}

	res, statusCode := userService.SetUserProfile(&u, params)

	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "user profile has been set already", res.Message)

	countryService.AssertExpectations(t)
	userProfileRepo.AssertExpectations(t)
}

type queryMatcher struct {
}

func (queryMatcher) Match(expectedSQL, actualSQL string) error {
	return nil
}

func TestService_SetUserProfile_Successful(t *testing.T) {
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
	dbMock.ExpectExec("UPDATE user_profiles").WillReturnResult(sqlmock.NewResult(1, 1))

	userRepo := new(mocks.UserRepository)
	userProfileRepo := new(mocks.UserProfileRepository)
	userProfileRepo.On("GetProfileByUserID", 1, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		profile := args.Get(1).(*user.Profile)
		profile.ID = 1
		profile.Status = sql.NullString{String: user.ProfileStatusIncomplete, Valid: true}
	})

	profileImageRepository := new(mocks.ProfileImageRepository)
	countryService := new(mocks.CountryService)
	countryService.On("GetCountryByID", int64(1)).Once().Return(country.Country{ID: 1, Name: sql.NullString{String: "test", Valid: true}}, nil)
	twoFaManager := new(mocks.TwoFaManager)
	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	jwtSerivce := new(mocks.JwtService)
	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	u := user.User{
		ID: 1,
	}
	params := user.SetUserProfileParams{
		FirstName:     "test",
		LastName:      "test",
		Gender:        "male",
		DateOfBirth:   "2007-06-06",
		Address:       "test",
		RegionAndCity: "test and test",
		PostalCode:    "test",
		CountryID:     1,
	}

	res, statusCode := userService.SetUserProfile(&u, params)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	result, ok := res.Data.(user.SetUserProfileResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	assert.Equal(t, "test", result.FirstName)
	assert.Equal(t, "test", result.LastName)
	assert.Equal(t, "male", result.Gender)
	assert.Equal(t, "2007-06-06", result.DateOfBirth)
	assert.Equal(t, "test", result.Address)
	assert.Equal(t, "test and test", result.RegionAndCity)
	assert.Equal(t, "test", result.PostalCode)
	assert.Equal(t, int64(1), result.CountryID)
	assert.Equal(t, "test", result.CountryName)
	assert.Equal(t, "processing", result.Status)

	countryService.AssertExpectations(t)
	userProfileRepo.AssertExpectations(t)
}

func TestService_GetProfile(t *testing.T) {
	db := &gorm.DB{}
	userRepo := new(mocks.UserRepository)
	userProfileRepo := new(mocks.UserProfileRepository)
	userProfileRepo.On("GetProfileByUserID", 1, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		profile := args.Get(1).(*user.Profile)
		profile.ID = 1
		profile.Status = sql.NullString{String: user.ProfileStatusProcessing, Valid: true}
		profile.CountryID = sql.NullInt64{Int64: 1, Valid: true}

		profile.FirstName = sql.NullString{String: "test", Valid: true}
		profile.LastName = sql.NullString{String: "test", Valid: true}
		profile.Gender = sql.NullString{String: "male", Valid: true}
		profile.DateOfBirth = sql.NullString{String: "2007-06-06", Valid: true}
		profile.Address = sql.NullString{String: "test", Valid: true}
		profile.RegionAndCity = sql.NullString{String: "test and test", Valid: true}
		profile.PostalCode = sql.NullString{String: "test", Valid: true}
	})

	profileImageRepository := new(mocks.ProfileImageRepository)
	imagesData := []user.ImagesQueryFields{
		{
			ID:   1,
			Type: user.ProfileImageTypeIdentity,
		},
		{
			ID:   2,
			Type: user.ProfileImageTypeIdentity,
		},
		{
			ID:   3,
			Type: user.ProfileImageTypeAddress,
		},
	}
	profileImageRepository.On("GetLatestImagesDataByProfileID", int64(1)).Once().Return(imagesData)
	images := []user.ProfileImage{
		{
			ID:               1,
			UserProfileID:    sql.NullInt64{Int64: 1, Valid: true},
			Type:             user.ProfileImageTypeIdentity,
			ImagePath:        "/image1",
			OriginalFileName: "test",
			IDCardCode:       sql.NullString{String: "1234", Valid: true},
			SubType:          sql.NullString{String: user.ProfileImageSubtypeIdentityIdentityCard, Valid: true},
			IsDeleted:        sql.NullBool{Bool: false, Valid: true},
		},
		{
			ID:               2,
			UserProfileID:    sql.NullInt64{Int64: 1, Valid: true},
			Type:             user.ProfileImageTypeIdentity,
			ImagePath:        "/image2",
			OriginalFileName: "test2",
			IDCardCode:       sql.NullString{String: "1234", Valid: true},
			SubType:          sql.NullString{String: user.ProfileImageSubtypeIdentityIdentityCard, Valid: true},
			MainImageID:      sql.NullInt64{Int64: 1, Valid: true},
			IsBack:           sql.NullBool{Bool: true, Valid: true},
		},
		{
			ID:               3,
			UserProfileID:    sql.NullInt64{Int64: 1, Valid: true},
			Type:             user.ProfileImageTypeAddress,
			ImagePath:        "/image3",
			OriginalFileName: "test",
			SubType:          sql.NullString{String: user.ProfileImageSubtypeAddressUtilityBill, Valid: true},
			IsDeleted:        sql.NullBool{Bool: false, Valid: true},
		},
	}
	profileImageRepository.On("GetImagesByIds", []int64{1, 2, 3}).Once().Return(images)

	countryService := new(mocks.CountryService)
	countryService.On("GetCountryByID", int64(1)).Once().Return(country.Country{ID: 1, Name: sql.NullString{String: "test", Valid: true}}, nil)
	twoFaManager := new(mocks.TwoFaManager)
	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	jwtSerivce := new(mocks.JwtService)
	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	u := user.User{
		ID: 1,
	}

	res, statusCode := userService.GetProfile(&u)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	result, ok := res.Data.(user.GetProfileResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "test", result.FirstName)
	assert.Equal(t, "test", result.LastName)
	assert.Equal(t, "male", result.Gender)
	assert.Equal(t, "2007-06-06", result.DateOfBirth)
	assert.Equal(t, "test", result.Address)
	assert.Equal(t, "test and test", result.RegionAndCity)
	assert.Equal(t, "test", result.PostalCode)
	assert.Equal(t, int64(1), *result.CountryID)
	assert.Equal(t, "test", *result.CountryName)
	assert.Equal(t, "processing", result.Status)
	assert.Equal(t, "", result.AdminComment)

	img1 := result.UserProfileImages[0]
	assert.Equal(t, "identity", img1.Type)
	assert.Equal(t, int64(1), img1.ID)
	assert.Equal(t, "1234", img1.IDCardCode)
	assert.Equal(t, int64(1), img1.ImageID)
	assert.Equal(t, false, img1.IsBack)
	assert.Equal(t, "identity_card", img1.SubType)

	img2 := result.UserProfileImages[1]
	assert.Equal(t, "address", img2.Type)
	assert.Equal(t, int64(3), img2.ID)
	assert.Equal(t, "", img2.IDCardCode)
	assert.Equal(t, int64(3), img2.ImageID)
	assert.Equal(t, false, img2.IsBack)
	assert.Equal(t, "utility_bill", img2.SubType)

	img3 := result.UserProfileImages[2]
	assert.Equal(t, "identity", img3.Type)
	assert.Equal(t, int64(2), img3.ID)
	assert.Equal(t, "1234", img3.IDCardCode)
	assert.Equal(t, int64(2), img3.ImageID)
	assert.Equal(t, int64(1), *img3.MainImageID)
	assert.Equal(t, true, img3.IsBack)
	assert.Equal(t, "identity_card", img3.SubType)

	countryService.AssertExpectations(t)
	userProfileRepo.AssertExpectations(t)
	profileImageRepository.AssertExpectations(t)
}

func TestService_GetUserData(t *testing.T) {
	db := &gorm.DB{}
	userRepo := new(mocks.UserRepository)
	userProfileRepo := new(mocks.UserProfileRepository)
	userProfileRepo.On("GetProfileByUserID", 1, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		profile := args.Get(1).(*user.Profile)
		profile.ID = 1
		profile.Status = sql.NullString{String: user.ProfileStatusProcessing, Valid: true}
		profile.CountryID = sql.NullInt64{Int64: 1, Valid: true}
	})

	profileImageRepository := new(mocks.ProfileImageRepository)
	countryService := new(mocks.CountryService)
	twoFaManager := new(mocks.TwoFaManager)
	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	jwtSerivce := new(mocks.JwtService)
	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	u := user.User{
		ID:                  1,
		IsTwoFaEnabled:      true,
		Phone:               sql.NullString{String: "+989121234567", Valid: true},
		Email:               "test@test.com",
		UbID:                "ub-test",
		Google2faSecretCode: sql.NullString{String: "sd", Valid: true},
		Status:              user.StatusVerified,
		PrivateChannelName:  "testChannel",
	}

	res, statusCode := userService.GetUserData(&u)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	result, ok := res.Data.(user.GetUserDataResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	assert.Equal(t, "test@test.com", result.Email)
	assert.Equal(t, "ub-test", result.UbID)
	assert.Equal(t, "+989121234567", result.Phone)
	assert.Equal(t, "none", result.KycLevel)
	assert.Equal(t, "processing", result.KycStatus)
	assert.Equal(t, "you identity,address, phone are not confirmed yet", result.KycLevelMessage)
	assert.Equal(t, "medium", result.SecurityLevel)
	assert.Equal(t, "we highly recommend to verify your identity", result.SecurityLevelMessage)
	assert.Equal(t, true, result.Google2faEnabled)
	assert.Equal(t, true, result.Has2fa)
	assert.Equal(t, true, result.IsAccountVerified)
	assert.Equal(t, "testChannel", result.ChannelName)
	assert.Equal(t, 0, result.ThemeID)
	assert.Equal(t, "default", result.Theme)
	assert.Equal(t, "processing", result.ProfileStatus)

	userProfileRepo.AssertExpectations(t)

}

func TestService_Get2FaBarcode(t *testing.T) {
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
	dbMock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(1, 1))

	userRepo := new(mocks.UserRepository)
	userProfileRepo := new(mocks.UserProfileRepository)
	profileImageRepository := new(mocks.ProfileImageRepository)
	countryService := new(mocks.CountryService)
	twoFaManager := user.NewTwoFaManager()
	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	jwtSerivce := new(mocks.JwtService)
	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	u := user.User{
		ID:                  1,
		Email:               "test@tes.com",
		IsTwoFaEnabled:      false,
		Google2faSecretCode: sql.NullString{String: "", Valid: false},
	}

	res, statusCode := userService.Get2FaBarcode(&u)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	result, ok := res.Data.(user.Get2FaBarcodeResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	assert.Equal(t, 16, len(result.Code))
}

func TestService_Enable2Fa(t *testing.T) {
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
	dbMock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(1, 1))

	userRepo := new(mocks.UserRepository)
	userProfileRepo := new(mocks.UserProfileRepository)
	profileImageRepository := new(mocks.ProfileImageRepository)
	countryService := new(mocks.CountryService)
	twoFaManager := new(mocks.TwoFaManager)
	twoFaManager.On("CheckCode", mock.Anything, "123456").Once().Return(true)
	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	jwtSerivce := new(mocks.JwtService)
	jwtSerivce.On("IssueToken", "test@test.com", mock.Anything, mock.Anything).Once().Return("token", nil)
	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	passwordHash, _ := pe.GenerateFromPassword("123456789")

	u := user.User{
		Email:    "test@test.com",
		Password: string(passwordHash),
	}

	params := user.Enable2FaParams{
		Password: "123456789",
		Code:     "123456",
	}

	res, statusCode := userService.Enable2Fa(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	response, ok := res.Data.(map[string]string)
	if !ok {
		t.Error("can not cast data to map")
		t.Fail()
	}
	assert.Equal(t, "token", response["token"])

	twoFaManager.AssertExpectations(t)
	jwtSerivce.AssertExpectations(t)
}

func TestService_Disable2Fa(t *testing.T) {
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
	dbMock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(1, 1))

	userRepo := new(mocks.UserRepository)
	userProfileRepo := new(mocks.UserProfileRepository)
	profileImageRepository := new(mocks.ProfileImageRepository)
	countryService := new(mocks.CountryService)
	twoFaManager := new(mocks.TwoFaManager)
	twoFaManager.On("CheckCode", mock.Anything, "123456").Once().Return(true)
	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	jwtSerivce := new(mocks.JwtService)
	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	passwordHash, _ := pe.GenerateFromPassword("123456789")

	u := user.User{
		Password:       string(passwordHash),
		IsTwoFaEnabled: true,
	}

	params := user.Disable2FaParams{
		Password: "123456789",
		Code:     "123456",
	}

	res, statusCode := userService.Disable2Fa(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	twoFaManager.AssertExpectations(t)
}

func TestService_ChangePassword_TwoFaEnabled(t *testing.T) {
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
	dbMock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(1, 1))

	userRepo := new(mocks.UserRepository)
	userProfileRepo := new(mocks.UserProfileRepository)
	profileImageRepository := new(mocks.ProfileImageRepository)
	countryService := new(mocks.CountryService)
	twoFaManager := new(mocks.TwoFaManager)
	twoFaManager.On("CheckCode", mock.Anything, "123456").Once().Return(true)

	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	jwtSerivce := new(mocks.JwtService)
	jwtSerivce.On("IssueToken", "test@test.com", mock.Anything, mock.Anything).Once().Return("token", nil)
	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	communicationService.On("SendPasswordChangedEmail", mock.Anything, mock.Anything).Once().Return()

	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	passwordHash, _ := pe.GenerateFromPassword("123456789")

	u := user.User{
		Email:          "test@test.com",
		Password:       string(passwordHash),
		IsTwoFaEnabled: true,
	}

	params := user.ChangePasswordParams{
		OldPassword:   "123456789",
		NewPassword:   "123456",
		Confirmed:     "123456",
		TwoFaCode:     "123456",
		UserAgentInfo: user.RequestUserAgentInfo{},
	}

	res, statusCode := userService.ChangePassword(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	response, ok := res.Data.(map[string]string)
	if !ok {
		t.Error("can not cast data to map")
		t.Fail()
	}
	assert.Equal(t, "token", response["token"])

	twoFaManager.AssertExpectations(t)
	time.Sleep(20 * time.Millisecond)
	communicationService.AssertExpectations(t)
	jwtSerivce.AssertExpectations(t)

}

func TestService_ChangePassword_TwoFaDisabled(t *testing.T) {
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
	dbMock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(1, 1))

	userRepo := new(mocks.UserRepository)
	userProfileRepo := new(mocks.UserProfileRepository)
	profileImageRepository := new(mocks.ProfileImageRepository)
	countryService := new(mocks.CountryService)
	twoFaManager := new(mocks.TwoFaManager)

	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	jwtSerivce := new(mocks.JwtService)
	jwtSerivce.On("IssueToken", "test@test.com", mock.Anything, mock.Anything).Once().Return("token", nil)
	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	communicationService.On("SendPasswordChangedEmail", mock.Anything, mock.Anything).Once().Return()

	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	passwordHash, _ := pe.GenerateFromPassword("123456789")

	u := user.User{
		Email:          "test@test.com",
		Password:       string(passwordHash),
		IsTwoFaEnabled: false,
	}

	params := user.ChangePasswordParams{
		OldPassword:   "123456789",
		NewPassword:   "123456",
		Confirmed:     "123456",
		TwoFaCode:     "",
		UserAgentInfo: user.RequestUserAgentInfo{},
	}

	res, statusCode := userService.ChangePassword(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	response, ok := res.Data.(map[string]string)
	if !ok {
		t.Error("can not cast data to map")
		t.Fail()
	}
	assert.Equal(t, "token", response["token"])
	time.Sleep(20 * time.Millisecond)
	communicationService.AssertExpectations(t)
	jwtSerivce.AssertExpectations(t)
}

func TestService_ChangePassword_WrongPassword(t *testing.T) {
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
	dbMock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(1, 1))

	userRepo := new(mocks.UserRepository)
	userProfileRepo := new(mocks.UserProfileRepository)
	profileImageRepository := new(mocks.ProfileImageRepository)
	countryService := new(mocks.CountryService)
	twoFaManager := new(mocks.TwoFaManager)

	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	jwtSerivce := new(mocks.JwtService)
	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)

	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	passwordHash, _ := pe.GenerateFromPassword("123456782") //not the same as the one user entered

	u := user.User{
		Password:       string(passwordHash),
		IsTwoFaEnabled: false,
	}

	params := user.ChangePasswordParams{
		OldPassword:   "123456789",
		NewPassword:   "123456",
		Confirmed:     "123456",
		TwoFaCode:     "",
		UserAgentInfo: user.RequestUserAgentInfo{},
	}

	res, statusCode := userService.ChangePassword(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "old password is not correct", res.Message)

	communicationService.AssertExpectations(t)
}

func TestService_SendSms_LessThan60Seconds(t *testing.T) {
	db := &gorm.DB{}
	userRepo := new(mocks.UserRepository)
	userProfileRepo := new(mocks.UserProfileRepository)
	profileImageRepository := new(mocks.ProfileImageRepository)
	countryService := new(mocks.CountryService)
	twoFaManager := new(mocks.TwoFaManager)

	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	phoneConfirmationManager.On("IsAllowedToSendSms", mock.Anything).Once().Return(false)
	jwtSerivce := new(mocks.JwtService)

	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)

	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	u := user.User{
		ID: 1,
	}

	params := user.SendSmsParams{
		Phone: "+989121234567",
	}

	res, statusCode := userService.SendSms(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "only one sms per minute can be send", res.Message)

	phoneConfirmationManager.AssertExpectations(t)

}

func TestService_SendSms_PhoneFormatIsNotCorrect(t *testing.T) {
	db := &gorm.DB{}
	userRepo := new(mocks.UserRepository)
	userProfileRepo := new(mocks.UserProfileRepository)
	profileImageRepository := new(mocks.ProfileImageRepository)
	countryService := new(mocks.CountryService)
	twoFaManager := new(mocks.TwoFaManager)

	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	jwtSerivce := new(mocks.JwtService)

	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)

	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	u := user.User{
		ID: 1,
	}

	params := user.SendSmsParams{
		Phone: "123456",
	}

	res, statusCode := userService.SendSms(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "the phone you entered is not valid", res.Message)
}

func TestService_SendSms_Successful(t *testing.T) {
	db := &gorm.DB{}
	userRepo := new(mocks.UserRepository)
	userProfileRepo := new(mocks.UserProfileRepository)
	profileImageRepository := new(mocks.ProfileImageRepository)
	countryService := new(mocks.CountryService)
	twoFaManager := new(mocks.TwoFaManager)

	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	phoneConfirmationManager.On("IsAllowedToSendSms", mock.Anything).Once().Return(true)
	phoneConfirmationManager.On("GeneratePhoneConfirmationCodeAndSendSms", mock.Anything, "+989121234567").Once().Return(nil)
	jwtSerivce := new(mocks.JwtService)

	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)

	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	u := user.User{
		ID: 1,
	}

	params := user.SendSmsParams{
		Phone: "+989121234567",
	}

	res, statusCode := userService.SendSms(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	phoneConfirmationManager.AssertExpectations(t)
}

func TestService_EnableSms_WrongCode(t *testing.T) {
	db := &gorm.DB{}
	userRepo := new(mocks.UserRepository)
	userProfileRepo := new(mocks.UserProfileRepository)
	profileImageRepository := new(mocks.ProfileImageRepository)
	countryService := new(mocks.CountryService)
	twoFaManager := new(mocks.TwoFaManager)

	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	phoneConfirmationManager.On("IsCodeCorrect", mock.Anything, "+989121234567", "123456").Once().Return(false)
	jwtSerivce := new(mocks.JwtService)

	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)

	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	passwordHash, _ := pe.GenerateFromPassword("123456789")

	u := user.User{
		ID:             1,
		Password:       string(passwordHash),
		IsTwoFaEnabled: false,
	}

	params := user.EnableSmsParams{
		Phone:     "+989121234567",
		Code:      "123456",
		TwoFaCode: "",
		Password:  "123456789",
	}

	res, statusCode := userService.EnableSms(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "code is not correct", res.Message)

	phoneConfirmationManager.AssertExpectations(t)

}

func TestService_EnableSms_TwoFaEnabled(t *testing.T) {
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
	dbMock.ExpectBegin()
	dbMock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_profiles").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	userRepo := new(mocks.UserRepository)
	userRepo.On("GetUserByIDUsingTx", mock.Anything, 1, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(2).(*user.User)
		u.ID = 1
	})
	userProfileRepo := new(mocks.UserProfileRepository)
	userProfileRepo.On("GetProfileByUserIDUsingTx", mock.Anything, 1, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		profile := args.Get(2).(*user.Profile)
		profile.ID = 1
	})

	profileImageRepository := new(mocks.ProfileImageRepository)
	countryService := new(mocks.CountryService)
	twoFaManager := new(mocks.TwoFaManager)
	twoFaManager.On("CheckCode", mock.Anything, "123456").Once().Return(true)

	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	phoneConfirmationManager.On("IsCodeCorrect", mock.Anything, "+989121234567", "123456").Once().Return(true)
	phoneConfirmationManager.On("DeleteKey", mock.Anything).Once().Return(nil)
	jwtSerivce := new(mocks.JwtService)

	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)

	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	u := user.User{
		ID:             1,
		IsTwoFaEnabled: true,
	}

	params := user.EnableSmsParams{
		Phone:     "+989121234567",
		Code:      "123456",
		TwoFaCode: "123456",
		Password:  "123456789",
	}

	res, statusCode := userService.EnableSms(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	time.Sleep(100 * time.Millisecond)
	phoneConfirmationManager.AssertExpectations(t)
	twoFaManager.AssertExpectations(t)
	userProfileRepo.AssertExpectations(t)
}

func TestService_EnableSms_TwoFaDisabled(t *testing.T) {
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
	dbMock.ExpectBegin()
	dbMock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_profiles").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	userRepo := new(mocks.UserRepository)
	userRepo.On("GetUserByIDUsingTx", mock.Anything, 1, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(2).(*user.User)
		u.ID = 1
	})
	userProfileRepo := new(mocks.UserProfileRepository)
	userProfileRepo.On("GetProfileByUserIDUsingTx", mock.Anything, 1, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		profile := args.Get(2).(*user.Profile)
		profile.ID = 1
	})

	profileImageRepository := new(mocks.ProfileImageRepository)
	countryService := new(mocks.CountryService)
	twoFaManager := new(mocks.TwoFaManager)

	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	phoneConfirmationManager.On("IsCodeCorrect", mock.Anything, "+989121234567", "123456").Once().Return(true)
	phoneConfirmationManager.On("DeleteKey", mock.Anything).Once().Return(nil)
	jwtSerivce := new(mocks.JwtService)

	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)

	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	passwordHash, _ := pe.GenerateFromPassword("123456789")

	u := user.User{
		ID:             1,
		Password:       string(passwordHash),
		IsTwoFaEnabled: false,
	}

	params := user.EnableSmsParams{
		Phone:     "+989121234567",
		Code:      "123456",
		TwoFaCode: "",
		Password:  "123456789",
	}

	res, statusCode := userService.EnableSms(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	time.Sleep(50 * time.Millisecond)
	phoneConfirmationManager.AssertExpectations(t)
	twoFaManager.AssertExpectations(t)
	userProfileRepo.AssertExpectations(t)
}

func TestService_DisableSms_WrongCode(t *testing.T) {
	db := &gorm.DB{}
	userRepo := new(mocks.UserRepository)
	userProfileRepo := new(mocks.UserProfileRepository)
	profileImageRepository := new(mocks.ProfileImageRepository)
	countryService := new(mocks.CountryService)
	twoFaManager := new(mocks.TwoFaManager)

	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	phoneConfirmationManager.On("IsCodeCorrect", mock.Anything, "+989121234567", "123456").Once().Return(false)
	jwtSerivce := new(mocks.JwtService)

	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)

	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	passwordHash, _ := pe.GenerateFromPassword("123456789")

	u := user.User{
		ID:             1,
		Password:       string(passwordHash),
		IsTwoFaEnabled: false,
		Phone:          sql.NullString{String: "+989121234567", Valid: true},
	}

	params := user.DisableSmsParams{
		Phone:     "+989121234567",
		Code:      "123456",
		TwoFaCode: "",
		Password:  "123456789",
	}

	res, statusCode := userService.DisableSms(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "code is not correct", res.Message)

	phoneConfirmationManager.AssertExpectations(t)

}

func TestService_DisableSms_TwoFaEnabled(t *testing.T) {
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
	dbMock.ExpectBegin()
	dbMock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_profiles").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	userRepo := new(mocks.UserRepository)
	userProfileRepo := new(mocks.UserProfileRepository)
	userProfileRepo.On("GetProfileByUserID", 1, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		profile := args.Get(1).(*user.Profile)
		profile.ID = 1
	})

	profileImageRepository := new(mocks.ProfileImageRepository)
	countryService := new(mocks.CountryService)
	twoFaManager := new(mocks.TwoFaManager)
	twoFaManager.On("CheckCode", mock.Anything, "123456").Once().Return(true)

	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	phoneConfirmationManager.On("IsCodeCorrect", mock.Anything, "+989121234567", "123456").Once().Return(true)
	phoneConfirmationManager.On("DeleteKey", mock.Anything).Once().Return(nil)
	jwtSerivce := new(mocks.JwtService)

	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)

	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	u := user.User{
		Phone:          sql.NullString{String: "+989121234567", Valid: true},
		ID:             1,
		IsTwoFaEnabled: true,
	}

	params := user.DisableSmsParams{
		Phone:     "+989121234567",
		Code:      "123456",
		TwoFaCode: "123456",
		Password:  "123456789",
	}

	res, statusCode := userService.DisableSms(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	time.Sleep(50 * time.Millisecond)
	phoneConfirmationManager.AssertExpectations(t)
	twoFaManager.AssertExpectations(t)
	userProfileRepo.AssertExpectations(t)
}

func TestService_DisableSms_TwoFaDisabled(t *testing.T) {
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
	dbMock.ExpectBegin()
	dbMock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_profiles").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	userRepo := new(mocks.UserRepository)
	userProfileRepo := new(mocks.UserProfileRepository)
	userProfileRepo.On("GetProfileByUserID", 1, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		profile := args.Get(1).(*user.Profile)
		profile.ID = 1
	})

	profileImageRepository := new(mocks.ProfileImageRepository)
	countryService := new(mocks.CountryService)
	twoFaManager := new(mocks.TwoFaManager)

	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	phoneConfirmationManager.On("IsCodeCorrect", mock.Anything, "+989121234567", "123456").Once().Return(true)
	phoneConfirmationManager.On("DeleteKey", mock.Anything).Once().Return(nil)

	jwtSerivce := new(mocks.JwtService)

	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)

	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	passwordHash, _ := pe.GenerateFromPassword("123456789")

	u := user.User{
		ID:             1,
		Password:       string(passwordHash),
		IsTwoFaEnabled: false,
		Phone:          sql.NullString{String: "+989121234567", Valid: true},
	}

	params := user.DisableSmsParams{
		Phone:     "+989121234567",
		Code:      "123456",
		TwoFaCode: "",
		Password:  "123456789",
	}

	res, statusCode := userService.DisableSms(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	time.Sleep(50 * time.Millisecond)
	phoneConfirmationManager.AssertExpectations(t)
	twoFaManager.AssertExpectations(t)
	userProfileRepo.AssertExpectations(t)
}

func TestService_SendVerificationEmail(t *testing.T) {
	db := &gorm.DB{}
	userRepo := new(mocks.UserRepository)
	userProfileRepo := new(mocks.UserProfileRepository)
	profileImageRepository := new(mocks.ProfileImageRepository)
	countryService := new(mocks.CountryService)
	twoFaManager := new(mocks.TwoFaManager)
	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	jwtSerivce := new(mocks.JwtService)
	pe := platform.NewPasswordEncoder()

	communicationService := new(mocks.CommunicationService)
	communicationService.On("SendVerificationEmailToUser", mock.Anything, mock.Anything).Once().Return()

	configs := new(mocks.Configs)
	configs.On("GetDomain").Once().Return("localhost")

	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	passwordHash, _ := pe.GenerateFromPassword("123456789")

	u := user.User{
		ID:             1,
		Password:       string(passwordHash),
		IsTwoFaEnabled: false,
		Phone:          sql.NullString{String: "+989121234567", Valid: true},
	}

	res, statusCode := userService.SendVerificationEmail(&u)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	configs.AssertExpectations(t)
	time.Sleep(20 * time.Millisecond)
	communicationService.AssertExpectations(t)
}

func TestService_DeleteImage(t *testing.T) {
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
	dbMock.ExpectExec("UPDATE user_profile_image").WillReturnResult(sqlmock.NewResult(1, 1))
	userRepo := new(mocks.UserRepository)
	userProfileRepo := new(mocks.UserProfileRepository)
	userProfileRepo.On("GetProfileByUserID", 1, mock.Anything).Once().Return(nil)
	profileImageRepository := new(mocks.ProfileImageRepository)
	image := &user.ProfileImage{}
	profileImageRepository.On("GetImageByID", int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		image = args.Get(1).(*user.ProfileImage)
		image.ID = 1
		image.ConfirmationStatus = sql.NullString{String: "PROCESSING", Valid: true}
	})

	countryService := new(mocks.CountryService)
	twoFaManager := new(mocks.TwoFaManager)
	phoneConfirmationManager := new(mocks.PhoneConfirmationManager)
	jwtSerivce := new(mocks.JwtService)
	pe := platform.NewPasswordEncoder()
	communicationService := new(mocks.CommunicationService)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	userService := user.NewUserService(db, userRepo, userProfileRepo, profileImageRepository, countryService, twoFaManager, pe, communicationService, phoneConfirmationManager, jwtSerivce, configs, logger)

	u := user.User{
		ID: 1,
	}
	params := user.DeleteImageParams{
		ID: 1,
	}

	res, statusCode := userService.DeleteImage(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	assert.Equal(t, true, image.IsDeleted.Valid)
	assert.Equal(t, true, image.IsDeleted.Bool)

	userProfileRepo.AssertExpectations(t)
	profileImageRepository.AssertExpectations(t)
}
