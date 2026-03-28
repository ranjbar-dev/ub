package user

import (
	"database/sql"
	"errors"
	"exchange-go/internal/communication"
	"exchange-go/internal/country"
	"exchange-go/internal/jwt"
	"exchange-go/internal/platform"
	"exchange-go/internal/response"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ttacon/libphonenumber"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	KycLevelMinimum       = 0
	KycLevel1Confirmation = 1
	KycLevel2Confirmation = 2

	StatusRegistered = "REGISTERED"
	StatusVerified   = "VERIFIED"

	AccountStatusUnblocked = "UNBLOCKED"
	AccountStatusBlocked   = "BLOCKED"

	ProfileStatusIncomplete         = "INCOMPLETE"
	ProfileStatusProcessing         = "PROCESSING"
	ProfileStatusConfirmed          = "CONFIRMED"
	ProfileStatusPartiallyConfirmed = "PARTIALLY_CONFIRMED"
	ProfileStatusRejected           = "REJECTED"

	ProfileImageTypeAddress  = "ADDRESS"
	ProfileImageTypeIdentity = "IDENTITY"

	ProfileImageSubtypeIdentityPassport      = "PASSPORT"
	ProfileImageSubtypeIdentityDriverLicense = "DRIVER_LICENSE"
	ProfileImageSubtypeIdentityIdentityCard  = "IDENTITY_CARD"

	ProfileImageSubtypeAddressBankStatement = "BANK_STATEMENT"
	ProfileImageSubtypeAddressUtilityBill   = "UTILITY_BILL"
	ProfileImageSubtypeAddressOther         = "OTHER"

	ProfileImageStatusIncomplete         = "INCOMPLETE"
	ProfileImageStatusProcessing         = "PROCESSING"
	ProfileImageStatusConfirmed          = "CONFIRMED"
	ProfileImageStatusPartiallyConfirmed = "PARTIALLY_CONFIRMED"
	ProfileImageStatusRejected           = "REJECTED"

	SecurityLevelLow           = "LOW"
	SecurityLevelMedium        = "MEDIUM"
	SecurityLevelHigh          = "HIGH"
	SecurityLevelLowMessage    = "we highly recommend to enable 2 factor authentication"
	SecurityLevelMediumMessage = "we highly recommend to verify your identity"
	SecurityLevelHighMessage   = "your security level is high"

	ProfilePhoneConfirmationTypeCall = "CALL"
	ProfilePhoneConfirmationTypeSms  = "SMS"

	UserProfileImagePathPrefix = "/assets/images/user-profile-images/"
	MaxUploadSize              = 5 * 1024 * 1024 // 5MB
)

type UploadImagesParams struct {
	Type         string `form:"type" binding:"required"`
	SubType      string `form:"sub_type" binding:"required"`
	FrontImageID int64  `form:"front_image_id"`
	BackImageID  int64  `form:"back_image_id"`
	IDCardCode   string `form:"id_card_code"`
	FrontImage   multipart.File
	FrontHeader  *multipart.FileHeader
	BackImage    multipart.File
	BackHeader   *multipart.FileHeader
}

type GetUsersFilters struct {
	Status string
}

type DeleteImageParams struct {
	ID int64 `json:"id" binding:"required"`
}

type UploadImageSingleResponse struct {
	ID   int64  `json:"id"`
	Path string `json:"path"`
}

// Service is the main user service providing profile management, two-factor authentication,
// password operations, SMS verification, email verification, and KYC image handling.
type Service interface {
	// GetUserProfile retrieves the profile associated with the given user.
	GetUserProfile(u User) (Profile, error)
	// GetUserProfileUsingTx retrieves the user profile within the provided database transaction.
	GetUserProfileUsingTx(tx *gorm.DB, u User) (Profile, error)
	// GetUserByID retrieves a user by their unique identifier.
	GetUserByID(ID int) (User, error)
	// SetUserProfile updates the user's profile with the given parameters.
	SetUserProfile(u *User, params SetUserProfileParams) (apiResponse response.APIResponse, statusCode int)
	// GetProfile returns the user's profile data as an API response.
	GetProfile(u *User) (apiResponse response.APIResponse, statusCode int)
	// GetUserData returns comprehensive user account data as an API response.
	GetUserData(u *User) (apiResponse response.APIResponse, statusCode int)
	// Get2FaBarcode generates a TOTP secret and QR code URL for Google Authenticator setup.
	Get2FaBarcode(u *User) (apiResponse response.APIResponse, statusCode int)
	// Enable2Fa activates two-factor authentication for the user after code verification.
	Enable2Fa(u *User, params Enable2FaParams) (apiResponse response.APIResponse, statusCode int)
	// Disable2Fa deactivates two-factor authentication for the user after code verification.
	Disable2Fa(u *User, params Disable2FaParams) (apiResponse response.APIResponse, statusCode int)
	// ChangePassword updates the user's password after verifying the current password.
	ChangePassword(u *User, params ChangePasswordParams) (apiResponse response.APIResponse, statusCode int)
	// SendSms initiates phone verification by sending an SMS code to the provided number.
	SendSms(u *User, params SendSmsParams) (apiResponse response.APIResponse, statusCode int)
	// EnableSms activates SMS-based phone verification for the user.
	EnableSms(u *User, params EnableSmsParams) (apiResponse response.APIResponse, statusCode int)
	// DisableSms deactivates SMS-based phone verification for the user.
	DisableSms(u *User, params DisableSmsParams) (apiResponse response.APIResponse, statusCode int)
	// SendVerificationEmail sends an email verification link to the user's registered email.
	SendVerificationEmail(u *User) (apiResponse response.APIResponse, statusCode int)
	// GetUsersDataForOrderMatching returns lightweight user data needed for order matching.
	GetUsersDataForOrderMatching(userIds []int) []UsersDataForOrderMatching
	// UploadImages handles uploading KYC document images for the user's profile.
	UploadImages(u *User, params UploadImagesParams) (apiResponse response.APIResponse, statusCode int)
	// DeleteImage removes a KYC document image from the user's profile.
	DeleteImage(u *User, params DeleteImageParams) (apiResponse response.APIResponse, statusCode int)
	// GetUsersByPagination returns a paginated list of users matching the given filters.
	GetUsersByPagination(page int64, pageSize int, filters map[string]interface{}) []User
}

type service struct {
	db                       *gorm.DB
	userRepository           Repository
	userProfileRepository    ProfileRepository
	profileImageRepository   ProfileImageRepository
	countryService           country.Service
	twoFaManager             TwoFaManager
	passwordEncoder          platform.PasswordEncoder
	communicationService     communication.Service
	phoneConfirmationManager PhoneConfirmationManager
	jwtService               jwt.Service
	configs                  platform.Configs
	logger                   platform.Logger
}

type UsersDataForOrderMatching struct {
	UserID             int
	UserEmail          string
	UserLevelID        int64
	UserPrivateChannel string
}

type ImagesQueryFields struct {
	ID   int64
	Type string
}

type SetUserProfileParams struct {
	FirstName     string `json:"first_name" binding:"required"`
	LastName      string `json:"last_name" binding:"required"`
	Gender        string `json:"gender" binding:"required,oneof='male' 'female' 'unknown'"`
	DateOfBirth   string `json:"date_of_birth" binding:"required"`
	Address       string `json:"address" binding:"required"`
	RegionAndCity string `json:"region_and_city" binding:"required"`
	PostalCode    string `json:"postal_code" binding:"required"`
	CountryID     int64  `json:"country_id" binding:"required,gt=0"`
}

type SetUserProfileResponse struct {
	ID            int64  `json:"id"`
	UpdatedAt     string `json:"updatedAt"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Gender        string `json:"gender"`
	DateOfBirth   string `json:"dateOfBirth"`
	Address       string `json:"address"`
	RegionAndCity string `json:"regionAndCity"`
	PostalCode    string `json:"postalCode"`
	CountryID     int64  `json:"country"`
	CountryName   string `json:"countryName"`
	Status        string `json:"status"`
}

type GetProfileResponse struct {
	ID                        int64                             `json:"id"`
	UpdatedAt                 string                            `json:"updatedAt"`
	FirstName                 string                            `json:"firstName"`
	LastName                  string                            `json:"lastName"`
	Gender                    string                            `json:"gender"`
	DateOfBirth               string                            `json:"dateOfBirth"`
	Address                   string                            `json:"address"`
	RegionAndCity             string                            `json:"regionAndCity"`
	PostalCode                string                            `json:"postalCode"`
	CountryID                 *int64                            `json:"country"`
	CountryName               *string                           `json:"countryName"`
	Status                    string                            `json:"status"`
	AdminComment              string                            `json:"adminComment"`
	UserProfileImages         []ProfileImageResponse            `json:"userProfileImages"`
	UserProfileImagesMetaData map[string][]profileImageMetadata `json:"userProfileImagesMetaData"`
}

type ProfileImageResponse struct {
	Image           string `json:"image"`
	Type            string `json:"type"`
	ID              int64  `json:"id"`
	IDCardCode      string `json:"idCardCode"`
	Status          string `json:"status"`
	RejectionReason string `json:"rejectionReason"`
	ImageID         int64  `json:"imageId"`
	MainImageID     *int64 `json:"mainImageId"`
	IsBack          bool   `json:"isBack"`
	SubType         string `json:"subType"`
	CreatedAt       string `json:"createdAt"`
}

type profileImageSubType struct {
	Name    string `json:"name"`
	HasBack bool   `json:"hasBack"`
}

type profileImageMetadata struct {
	Name     string                `json:"name"`
	SubTypes []profileImageSubType `json:"subTypes"`
}

type GetUserDataResponse struct {
	Email                string `json:"email"`
	UbID                 string `json:"ubId"`
	Phone                string `json:"phone"`
	KycLevel             string `json:"kycLevel"`
	KycStatus            string `json:"kycStatus"`
	KycLevelMessage      string `json:"kycLevelMessage"`
	SecurityLevel        string `json:"securityLevel"`
	SecurityLevelMessage string `json:"securityLevelMessage"`
	Google2faEnabled     bool   `json:"google2faEnabled"`
	Has2fa               bool   `json:"has2fa"`
	IsAccountVerified    bool   `json:"isAccountVerified"`
	ChannelName          string `json:"channelName"`
	ThemeID              int    `json:"themeId"`
	Theme                string `json:"theme"`
	ProfileStatus        string `json:"profileStatus"`
}

type Get2FaBarcodeResponse struct {
	URL  string `json:"url"`
	Code string `json:"code"`
}

type Enable2FaParams struct {
	Password  string `json:"password" binding:"required"`
	Code      string `json:"code" binding:"required"`
	IP        string
	UserAgent string
}

type Disable2FaParams struct {
	Password string `json:"password" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

type RequestUserAgentInfo struct {
	Device  string `json:"device"`
	IP      string `json:"ip"`
	Browser string `json:"browser"`
}

type ChangePasswordParams struct {
	OldPassword   string `json:"old_password" binding:"required"`
	NewPassword   string `json:"new_password" binding:"required"`
	Confirmed     string `json:"confirmed" binding:"required"`
	TwoFaCode     string `json:"2fa_code"`
	IP            string
	UserAgent     string
	UserAgentInfo RequestUserAgentInfo
}

type SendSmsParams struct {
	Phone string `json:"phone" binding:"required"`
}

type EnableSmsParams struct {
	Phone     string `json:"phone" binding:"required"`
	Code      string `json:"code" binding:"required"`
	TwoFaCode string `json:"2fa_code"`
	Password  string `json:"password"`
}

type DisableSmsParams struct {
	Phone     string `json:"phone" binding:"required"`
	Code      string `json:"code" binding:"required"`
	TwoFaCode string `json:"2fa_code"`
	Password  string `json:"password"`
}

func (s *service) SetUserProfile(u *User, params SetUserProfileParams) (apiResponse response.APIResponse, statusCode int) {
	country, err := s.countryService.GetCountryByID(params.CountryID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("can not get country by id", err,
			zap.String("service", "userService"),
			zap.String("method", "SetUserProfile"),
			zap.Int64("countryID", params.CountryID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || country.ID == 0 {
		return response.Error("country not found", http.StatusUnprocessableEntity, nil)
	}

	userProfile, err := s.GetUserProfile(*u)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("can not get profile of user ", err,
			zap.String("service", "userService"),
			zap.String("method", "SetUserProfile"),
			zap.Int("userID", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if userProfile.Status.Valid && userProfile.Status.String == ProfileStatusConfirmed {
		return response.Error("user profile has been set already", http.StatusUnprocessableEntity, nil)
	}

	userProfile.FirstName = sql.NullString{String: params.FirstName, Valid: true}
	userProfile.LastName = sql.NullString{String: params.LastName, Valid: true}
	userProfile.Gender = sql.NullString{String: strings.ToUpper(params.Gender), Valid: true}
	userProfile.DateOfBirth = sql.NullString{String: params.DateOfBirth, Valid: true}
	userProfile.Address = sql.NullString{String: params.Address, Valid: true}
	userProfile.RegionAndCity = sql.NullString{String: params.RegionAndCity, Valid: true}
	userProfile.PostalCode = sql.NullString{String: params.PostalCode, Valid: true}
	userProfile.CountryID = sql.NullInt64{Int64: params.CountryID, Valid: true}
	userProfile.Status = sql.NullString{String: ProfileStatusProcessing, Valid: true}

	if userProfile.UserID == 0 {
		userProfile.UserID = u.ID
	}

	err = s.db.Omit(clause.Associations).Save(&userProfile).Error

	if err != nil {
		s.logger.Error2("can not save userProfile", err,
			zap.String("service", "userService"),
			zap.String("method", "SetUserProfile"),
			zap.Int("userID", u.ID),
			zap.Int64("userProfileID", userProfile.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	res := SetUserProfileResponse{
		ID:            userProfile.ID,
		UpdatedAt:     userProfile.UpdatedAt.Format("2006-01-02 15:04:05"),
		FirstName:     userProfile.FirstName.String,
		LastName:      userProfile.LastName.String,
		Gender:        strings.ToLower(userProfile.Gender.String),
		DateOfBirth:   userProfile.DateOfBirth.String,
		Address:       userProfile.Address.String,
		RegionAndCity: userProfile.RegionAndCity.String,
		PostalCode:    userProfile.PostalCode.String,
		CountryID:     country.ID,
		CountryName:   country.Name.String,
		Status:        strings.ToLower(userProfile.Status.String),
	}

	return response.Success(res, "")
}

func (s *service) GetProfile(u *User) (apiResponse response.APIResponse, statusCode int) {
	userProfile, err := s.GetUserProfile(*u)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("can not get userProfile", err,
			zap.String("service", "userService"),
			zap.String("method", "GetProfile"),
			zap.Int("userID", u.ID),
		)
	}

	var imagesIds []int64
	if userProfile.ID > 0 {
		imagesData := s.profileImageRepository.GetLatestImagesDataByProfileID(userProfile.ID)
		for _, imageData := range imagesData {
			imagesIds = append(imagesIds, imageData.ID)
		}
	}

	var images []ProfileImage
	if len(imagesIds) > 0 {
		images = s.profileImageRepository.GetImagesByIds(imagesIds)
	}

	var finalImages []ProfileImage

	//first loop  we append the images that are not the backside
	//then in second loop we try to find the backside for them
	for _, image := range images {
		if !image.IsBack.Valid || image.IsBack.Bool == false {
			finalImages = append(finalImages, image)
		}
	}
	for _, image := range images {
		for _, finalImage := range finalImages {
			if image.MainImageID.Int64 == finalImage.ID {
				finalImages = append(finalImages, image)
			}
		}
	}

	profileImageResponses := make([]ProfileImageResponse, 0)

	for _, image := range finalImages {
		imagePath := ""
		if image.ConfirmationStatus.String == ProfileImageStatusProcessing {
			imagePath = s.configs.GetImagePath() + image.ImagePath
		}
		var mainImageID *int64
		if image.MainImageID.Valid {
			mainImageID = &image.MainImageID.Int64
		}
		img := ProfileImageResponse{
			Image:           imagePath,
			Type:            strings.ToLower(image.Type),
			ID:              image.ID,
			IDCardCode:      image.IDCardCode.String,
			Status:          strings.ToLower(image.ConfirmationStatus.String),
			RejectionReason: image.RejectionReason.String,
			ImageID:         image.ID,
			MainImageID:     mainImageID,
			IsBack:          image.IsBack.Bool,
			SubType:         strings.ToLower(image.SubType.String),
			CreatedAt:       image.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		profileImageResponses = append(profileImageResponses, img)
	}

	//return response

	var countryID *int64
	var countryName *string
	if userProfile.CountryID.Valid {
		c, _ := s.countryService.GetCountryByID(userProfile.CountryID.Int64)
		countryID = &userProfile.CountryID.Int64
		countryName = &c.Name.String
	}
	res := GetProfileResponse{
		ID:                        userProfile.ID,
		UpdatedAt:                 userProfile.UpdatedAt.Format("2006-01-02 15:04:05"),
		FirstName:                 userProfile.FirstName.String,
		LastName:                  userProfile.LastName.String,
		Gender:                    strings.ToLower(userProfile.Gender.String),
		DateOfBirth:               userProfile.DateOfBirth.String,
		Address:                   userProfile.Address.String,
		RegionAndCity:             userProfile.RegionAndCity.String,
		PostalCode:                userProfile.PostalCode.String,
		CountryID:                 countryID,
		CountryName:               countryName,
		Status:                    strings.ToLower(userProfile.Status.String),
		AdminComment:              userProfile.AdminComment.String,
		UserProfileImages:         profileImageResponses,
		UserProfileImagesMetaData: s.getProfileImagesMetaData(),
	}

	return response.Success(res, "")
}

func (s *service) getProfileImagesMetaData() map[string][]profileImageMetadata {
	metadata := []profileImageMetadata{
		{
			Name: strings.ToLower(ProfileImageTypeAddress),
			SubTypes: []profileImageSubType{
				{
					Name:    strings.ToLower(ProfileImageSubtypeAddressBankStatement),
					HasBack: false,
				},
				{
					Name:    strings.ToLower(ProfileImageSubtypeAddressUtilityBill),
					HasBack: false,
				},
				{
					Name:    strings.ToLower(ProfileImageSubtypeAddressOther),
					HasBack: false,
				},
			},
		},
		{
			Name: strings.ToLower(ProfileImageTypeIdentity),
			SubTypes: []profileImageSubType{
				{
					Name:    strings.ToLower(ProfileImageSubtypeIdentityIdentityCard),
					HasBack: true,
				},
				{
					Name:    strings.ToLower(ProfileImageSubtypeIdentityDriverLicense),
					HasBack: true,
				},
				{
					Name:    strings.ToLower(ProfileImageSubtypeIdentityPassport),
					HasBack: false,
				},
			},
		},
	}
	data := make(map[string][]profileImageMetadata)
	data["types"] = metadata
	return data
}

func (s *service) GetUserProfile(u User) (Profile, error) {
	up := Profile{}
	err := s.userProfileRepository.GetProfileByUserID(u.ID, &up)
	return up, err
}

func (s *service) GetUserProfileUsingTx(tx *gorm.DB, u User) (Profile, error) {
	up := Profile{}
	err := s.userProfileRepository.GetProfileByUserIDUsingTx(tx, u.ID, &up)
	return up, err
}

func (s *service) GetUserByID(ID int) (User, error) {
	u := User{}
	err := s.userRepository.GetUserByID(ID, &u)
	return u, err
}

func (s *service) GetUserData(u *User) (apiResponse response.APIResponse, statusCode int) {
	userProfile, err := s.GetUserProfile(*u)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("can not get userProfile", err,
			zap.String("service", "userService"),
			zap.String("method", "GetUserData"),
			zap.Int("userID", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	KycStatus := strings.ToLower(userProfile.Status.String)
	KycLevel := ""
	KycLevelMessage := ""

	if u.Kyc == KycLevelMinimum {
		KycLevel = "none"
		KycLevelMessage = "you identity,address, phone are not confirmed yet"
	} else {
		KycLevel = "profile completed"
		KycLevelMessage = "you have completed your document verification"
	}

	securityLevel := ""
	securityMessage := ""

	if !u.IsTwoFaEnabled {
		securityLevel = SecurityLevelLow
		securityMessage = SecurityLevelLowMessage
	} else {
		securityLevel = SecurityLevelMedium
		securityMessage = SecurityLevelMediumMessage
		if userProfile.Status.String == ProfileStatusConfirmed || userProfile.Status.String == ProfileStatusPartiallyConfirmed {
			securityLevel = SecurityLevelHigh
			securityMessage = SecurityLevelHighMessage
		}
	}

	hasTwoFa := false
	if u.Google2faSecretCode.String != "" {
		hasTwoFa = true
	}

	isAccountVerified := false
	if u.Status == StatusVerified {
		isAccountVerified = true
	}

	res := GetUserDataResponse{
		Email:                u.Email,
		UbID:                 u.UbID,
		Phone:                u.Phone.String,
		KycLevel:             KycLevel,
		KycStatus:            KycStatus,
		KycLevelMessage:      KycLevelMessage,
		SecurityLevel:        strings.ToLower(securityLevel),
		SecurityLevelMessage: securityMessage,
		Google2faEnabled:     u.IsTwoFaEnabled,
		Has2fa:               hasTwoFa,
		IsAccountVerified:    isAccountVerified,
		ChannelName:          u.PrivateChannelName,
		ThemeID:              0,
		Theme:                "default",
		ProfileStatus:        strings.ToLower(userProfile.Status.String),
	}

	return response.Success(res, "")
}

func (s *service) Get2FaBarcode(u *User) (apiResponse response.APIResponse, statusCode int) {
	if u.IsTwoFaEnabled {
		return response.Error("two fa is already enabled", http.StatusUnprocessableEntity, nil)
	}

	secretCode := ""
	url := ""

	if !u.Google2faSecretCode.Valid || u.Google2faSecretCode.String == "" {
		var err error
		secretCode, url, err = s.twoFaManager.GenerateSecretCode(*u)
		if err != nil {
			s.logger.Error2("can not generate secret code", err,
				zap.String("service", "userService"),
				zap.String("method", "Get2FaBarcode"),
				zap.Int("userID", u.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}

		u.Google2faSecretCode = sql.NullString{String: secretCode, Valid: true}
		err = s.db.Omit(clause.Associations).Save(u).Error
		if err != nil {
			s.logger.Error2("can not save user", err,
				zap.String("service", "userService"),
				zap.String("method", "Get2FaBarcode"),
				zap.Int("userID", u.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
	} else {
		var err error
		secretCode, url, err = s.twoFaManager.GenerateSecretCode(*u)
		if err != nil {
			s.logger.Error2("can not genrate secret code", err,
				zap.String("service", "userService"),
				zap.String("method", "Get2FaBarcode"),
				zap.Int("userID", u.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
	}

	//get url and code
	res := Get2FaBarcodeResponse{
		URL:  url,
		Code: secretCode,
	}

	return response.Success(res, "")
}

func (s *service) Enable2Fa(u *User, params Enable2FaParams) (apiResponse response.APIResponse, statusCode int) {
	if u.IsTwoFaEnabled {
		return response.Error("two fa is already enabled", http.StatusUnprocessableEntity, nil)
	}

	err := s.passwordEncoder.CompareHashAndPassword(u.Password, params.Password)
	if err != nil {
		return response.Error("password is not correct", http.StatusUnprocessableEntity, nil)
	}

	if !s.twoFaManager.CheckCode(*u, params.Code) {
		return response.Error("2fa code is not correct", http.StatusUnprocessableEntity, nil)
	}

	u.IsTwoFaEnabled = true
	//even enabling 2fa code will update this field because after enabling 2fa users must not able withdraw for 24 hours
	u.Google2faDisabledAt = sql.NullTime{Time: time.Now(), Valid: true}
	u.TwoFaChangedAt = sql.NullTime{Time: time.Now(), Valid: true}
	err = s.db.Omit(clause.Associations).Save(u).Error
	if err != nil {
		s.logger.Error2("can not save user", err,
			zap.String("service", "userService"),
			zap.String("method", "Enable2Fa"),
			zap.Int("userID", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	token, err := s.jwtService.IssueToken(u.Email, params.UserAgent, params.IP)
	if err != nil {
		s.logger.Error2("can not issue token for user", err,
			zap.String("service", "userService"),
			zap.String("method", "ChangePassword"),
			zap.Int("userID", u.ID),
		)
	}

	res := make(map[string]string, 0)
	res["token"] = token
	return response.Success(res, "")
}

func (s *service) Disable2Fa(u *User, params Disable2FaParams) (apiResponse response.APIResponse, statusCode int) {
	if !u.IsTwoFaEnabled {
		return response.Error("two fa is already disabled", http.StatusUnprocessableEntity, nil)
	}

	err := s.passwordEncoder.CompareHashAndPassword(u.Password, params.Password)
	if err != nil {
		return response.Error("password is not correct", http.StatusUnprocessableEntity, nil)
	}

	if !s.twoFaManager.CheckCode(*u, params.Code) {
		return response.Error("2fa code is not correct", http.StatusUnprocessableEntity, nil)
	}

	u.IsTwoFaEnabled = false
	u.Google2faDisabledAt = sql.NullTime{Time: time.Now(), Valid: true}
	u.Google2faSecretCode = sql.NullString{String: "", Valid: false}
	err = s.db.Omit(clause.Associations).Save(u).Error

	if err != nil {
		s.logger.Error2("can not save user", err,
			zap.String("service", "userService"),
			zap.String("method", "Disable2Fa"),
			zap.Int("userID", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	res := make(map[string]string, 0)
	return response.Success(res, "")
}

func (s *service) ChangePassword(u *User, params ChangePasswordParams) (apiResponse response.APIResponse, statusCode int) {
	if u.IsTwoFaEnabled {
		if strings.Trim(params.TwoFaCode, "") == "" {
			data := map[string]bool{
				"need2fa": true,
			}
			return response.Success(data, "")
		}

		if !s.twoFaManager.CheckCode(*u, params.TwoFaCode) {
			return response.Error("2fa authentication failed", http.StatusUnprocessableEntity, nil)
		}
	}

	err := s.passwordEncoder.CompareHashAndPassword(u.Password, params.OldPassword)
	if err != nil {
		return response.Error("old password is not correct", http.StatusUnprocessableEntity, nil)
	}

	//check new password and its confirmed are equal
	if params.Confirmed != params.NewPassword {
		return response.Error("new password and confirm password does not match", http.StatusUnprocessableEntity, nil)
	}

	passwordBytes, err := s.passwordEncoder.GenerateFromPassword(params.NewPassword)
	if err != nil {
		s.logger.Error2("can not generate password hash", err,
			zap.String("service", "userService"),
			zap.String("method", "ChangePassword"),
			zap.Int("userID", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	passwordHash := string(passwordBytes)

	u.Password = passwordHash
	u.PasswordChangedAt = sql.NullTime{Time: time.Now(), Valid: true}
	err = s.db.Omit(clause.Associations).Save(u).Error

	if err != nil {
		s.logger.Error2("can not save user", err,
			zap.String("service", "userService"),
			zap.String("method", "ChangePassword"),
			zap.Int("userID", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	token, err := s.jwtService.IssueToken(u.Email, params.UserAgent, params.IP)
	if err != nil {
		s.logger.Error2("can not issue token for user", err,
			zap.String("service", "userService"),
			zap.String("method", "ChangePassword"),
			zap.Int("userID", u.ID),
		)
	}

	emailParams := communication.PasswordChangedEmailParams{
		Email:         u.Email,
		CurrentIP:     params.UserAgentInfo.IP,
		CurrentDevice: params.UserAgentInfo.Device,
		ChangedDate:   time.Now().Format("2006-01-02 15:04:05"),
	}
	cu := communication.CommunicatingUser{
		Email: u.Email,
		Phone: "",
	}
	go s.communicationService.SendPasswordChangedEmail(cu, emailParams)
	res := make(map[string]string, 0)
	res["token"] = token
	return response.Success(res, "")
}

func (s *service) SendSms(u *User, params SendSmsParams) (apiResponse response.APIResponse, statusCode int) {
	num, err := libphonenumber.Parse(params.Phone, "")
	if err != nil {
		return response.Error("the phone you entered is not valid", http.StatusUnprocessableEntity, nil)
	}

	if !s.phoneConfirmationManager.IsAllowedToSendSms(*u) {
		return response.Error("only one sms per minute can be send", http.StatusUnprocessableEntity, nil)
	}

	phone := libphonenumber.Format(num, libphonenumber.E164)

	err = s.phoneConfirmationManager.GeneratePhoneConfirmationCodeAndSendSms(*u, phone)
	if err != nil {
		s.logger.Error2("can not generate phone confirmation and send sms", err,
			zap.String("service", "userService"),
			zap.String("method", "SendSms"),
			zap.String("phone", phone),
			zap.Int("userID", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	return response.Success(nil, "")
}

func (s *service) EnableSms(u *User, params EnableSmsParams) (apiResponse response.APIResponse, statusCode int) {
	if u.IsTwoFaEnabled {
		if strings.Trim(params.TwoFaCode, "") == "" {
			data := map[string]bool{
				"need2fa": true,
			}
			return response.Success(data, "")
		}

		if !s.twoFaManager.CheckCode(*u, params.TwoFaCode) {
			return response.Error("2fa authentication failed", http.StatusUnprocessableEntity, nil)
		}
	} else {
		err := s.passwordEncoder.CompareHashAndPassword(u.Password, params.Password)
		if err != nil {
			return response.Error("password is not correct", http.StatusUnprocessableEntity, nil)
		}
	}

	num, err := libphonenumber.Parse(params.Phone, "")
	if err != nil {
		return response.Error("the phone you entered is not valid ", http.StatusUnprocessableEntity, nil)
	}
	phone := libphonenumber.Format(num, libphonenumber.E164)

	if !s.phoneConfirmationManager.IsCodeCorrect(*u, phone, params.Code) {
		return response.Error("code is not correct", http.StatusUnprocessableEntity, nil)
	}

	tx := s.db.Begin()
	err = tx.Error
	if err != nil {
		s.logger.Error2("can not start transaction", err,
			zap.String("service", "userService"),
			zap.String("method", "EnableSms"),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	user := &User{}
	err = s.userRepository.GetUserByIDUsingTx(tx, u.ID, user)
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not get user by id", err,
			zap.String("service", "userService"),
			zap.String("method", "EnableSms"),
			zap.Int("userID", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	userProfile, err := s.GetUserProfileUsingTx(tx, *user)
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not get userProfile", err,
			zap.String("service", "userService"),
			zap.String("method", "EnableSms"),
			zap.Int("userID", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	user.Phone = sql.NullString{String: phone, Valid: true}
	userProfile.PhoneConfirmationType = sql.NullString{String: ProfilePhoneConfirmationTypeSms, Valid: true}

	err = tx.Omit(clause.Associations).Save(user).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not save user", err,
			zap.String("service", "userService"),
			zap.String("method", "EnableSms"),
			zap.Int("userID", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	err = tx.Omit(clause.Associations).Save(&userProfile).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not save userProfile", err,
			zap.String("service", "userService"),
			zap.String("method", "EnableSms"),
			zap.Int("userID", u.ID),
			zap.Int64("userProfileID", userProfile.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not commit transaction", err,
			zap.String("service", "userService"),
			zap.String("method", "EnableSms"),
			zap.Int("userID", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	go func() {
		err := s.phoneConfirmationManager.DeleteKey(*u)
		if err != nil {
			s.logger.Warn("can not delete phone confirmation from redis",
				zap.Error(err),
				zap.String("service", "userService"),
				zap.String("method", "EnableSms"),
				zap.Int("userID", u.ID),
			)
		}
	}()

	go s.calculateKyc(u, &userProfile)

	return response.Success(make(map[string]bool, 0), "")
}

func (s *service) DisableSms(u *User, params DisableSmsParams) (apiResponse response.APIResponse, statusCode int) {
	if !u.Phone.Valid {
		return response.Error("phone already disabled", http.StatusUnprocessableEntity, nil)
	}

	if u.IsTwoFaEnabled {
		if strings.Trim(params.TwoFaCode, "") == "" {
			data := map[string]bool{
				"need2fa": true,
			}
			return response.Success(data, "")
		}

		if !s.twoFaManager.CheckCode(*u, params.TwoFaCode) {
			return response.Error("2fa authentication failed", http.StatusUnprocessableEntity, nil)
		}
	} else {
		err := s.passwordEncoder.CompareHashAndPassword(u.Password, params.Password)
		if err != nil {
			return response.Error("password is not correct", http.StatusUnprocessableEntity, nil)
		}
	}

	num, err := libphonenumber.Parse(params.Phone, "")
	if err != nil {
		return response.Error("the phone you entered is not valid ", http.StatusUnprocessableEntity, nil)
	}
	phone := libphonenumber.Format(num, libphonenumber.E164)

	if u.Phone.String != phone {
		return response.Error("this phone is not submitted yet", http.StatusUnprocessableEntity, nil)
	}

	if !s.phoneConfirmationManager.IsCodeCorrect(*u, phone, params.Code) {
		return response.Error("code is not correct", http.StatusUnprocessableEntity, nil)
	}

	u.Phone = sql.NullString{String: "", Valid: false}
	err = s.db.Omit(clause.Associations).Save(u).Error
	if err != nil {
		s.logger.Error2("can not save user", err,
			zap.String("service", "userService"),
			zap.String("method", "DisableSms"),
			zap.Int("userID", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	go func() {
		err := s.phoneConfirmationManager.DeleteKey(*u)
		if err != nil {
			s.logger.Warn("can not delete phone confirmation from redis",
				zap.Error(err),
				zap.String("service", "userService"),
				zap.String("method", "EnableSms"),
				zap.Int("userID", u.ID),
			)
		}
	}()

	go s.calculateKyc(u, nil)

	return response.Success(nil, "")
}

func (s *service) calculateKyc(u *User, up *Profile) {
	if up == nil {
		userProfile, _ := s.GetUserProfile(*u)
		up = &userProfile
	}

	kyc := KycLevelMinimum
	if u.Phone.Valid && u.Phone.String != "" && up.Status.Valid && up.Status.String == ProfileStatusConfirmed {
		kyc = KycLevel1Confirmation
	}

	if u.Kyc != kyc {
		u.Kyc = kyc
		err := s.db.Omit(clause.Associations).Save(u).Error
		if err != nil {
			s.logger.Error2("can not save user", err,
				zap.String("service", "userService"),
				zap.String("method", "calculateKyc"),
				zap.Int("userID", u.ID),
			)
		}
	}
}

func (s *service) SendVerificationEmail(u *User) (apiResponse response.APIResponse, statusCode int) {
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
	go s.communicationService.SendVerificationEmailToUser(cu, link)

	return response.Success(nil, "")
}

func (s *service) GetUsersDataForOrderMatching(userIds []int) []UsersDataForOrderMatching {
	return s.userRepository.GetUsersDataForOrderMatching(userIds)
}

func (s *service) UploadImages(u *User, params UploadImagesParams) (apiResponse response.APIResponse, statusCode int) {
	params.Type = strings.ToUpper(params.Type)
	params.SubType = strings.ToUpper(params.SubType)
	if params.Type == "" || (params.Type != ProfileImageTypeAddress && params.Type != ProfileImageTypeIdentity) {
		return response.Error("type is not valid", http.StatusUnprocessableEntity, nil)
	}
	if params.SubType == "" ||
		(params.SubType != ProfileImageSubtypeAddressBankStatement &&
			params.SubType != ProfileImageSubtypeAddressUtilityBill &&
			params.SubType != ProfileImageSubtypeAddressOther &&
			params.SubType != ProfileImageSubtypeIdentityDriverLicense &&
			params.SubType != ProfileImageSubtypeIdentityIdentityCard &&
			params.SubType != ProfileImageSubtypeIdentityPassport) {
		return response.Error("sub type is not valid", http.StatusUnprocessableEntity, nil)
	}
	if params.FrontHeader == nil && params.BackHeader == nil {
		return response.Error("no file is uploaded", http.StatusUnprocessableEntity, nil)
	}
	if params.FrontHeader == nil && params.FrontImageID == 0 {
		return response.Error("can not set back image without providing front image", http.StatusUnprocessableEntity, nil)
	}

	up, err := s.GetUserProfile(*u)
	if err != nil {
		s.logger.Error2("can not get user profile", err,
			zap.String("service", "userService"),
			zap.String("method", "UploadImages"),
			zap.Int("userID", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	//check if user can upload
	errorMessage := s.checkIfUserCanUpload(up, params.Type)
	if errorMessage != "" {
		return response.Error(errorMessage, http.StatusUnprocessableEntity, nil)
	}
	res := make(map[string]UploadImageSingleResponse)
	tx := s.db.Begin()
	err = tx.Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error beginning transaction", err,
			zap.String("service", "userService"),
			zap.String("method", "UploadImages"),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	var frontProfileImage *ProfileImage
	if params.FrontHeader != nil {
		if params.FrontHeader.Size > MaxUploadSize {
			return response.Error("max upload size is 5M", http.StatusUnprocessableEntity, nil)
		}
		frontImageBuff := make([]byte, 512)
		_, err := params.FrontImage.Read(frontImageBuff)
		if err != nil {
			tx.Rollback()
			s.logger.Error2("error reading fornt image buff", err,
				zap.String("service", "userService"),
				zap.String("method", "UploadImages"),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
		fileType := http.DetectContentType(frontImageBuff)
		if fileType != "image/jpeg" && fileType != "image/png" && fileType != "application/pdf" {
			return response.Error("file type must be jpeg, png or pdf", http.StatusUnprocessableEntity, nil)
		}
		//When we read the first 512 bytes of the uploaded file in order to determine the content type,
		//the underlying file stream pointer moves forward by 512 bytes. When io.Copy() is called later,
		//it continues reading from that position resulting in a corrupted image file. The file.Seek()
		//method is used to return the pointer back to the start of the file so that io.Copy() starts from the beginning.
		_, err = params.FrontImage.Seek(0, io.SeekStart)
		if err != nil {
			tx.Rollback()
			s.logger.Error2("error seeking the frontImage", err,
				zap.String("service", "userService"),
				zap.String("method", "UploadImages"),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}

		originalName := params.FrontHeader.Filename
		newFilePath := UserProfileImagePathPrefix + uuid.NewString() + filepath.Ext(params.FrontHeader.Filename)
		err = os.MkdirAll("."+UserProfileImagePathPrefix, 0700)
		if err != nil {
			tx.Rollback()
			s.logger.Error2("can not create directory for images", err,
				zap.String("service", "userService"),
				zap.String("method", "UploadImages"),
				zap.Int("userID", u.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
		dest, err := os.Create("." + newFilePath)
		if err != nil {
			tx.Rollback()
			s.logger.Error2("can not create front image new file path", err,
				zap.String("service", "userService"),
				zap.String("method", "UploadImages"),
				zap.Int("userID", u.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
		_, err = io.Copy(dest, params.FrontImage)
		if err != nil {
			tx.Rollback()
			s.logger.Error2("can not copy front image", err,
				zap.String("service", "userService"),
				zap.String("method", "UploadImages"),
				zap.Int("userID", u.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)

		}
		frontProfileImage = &ProfileImage{
			UserProfileID:      sql.NullInt64{Int64: up.ID, Valid: true},
			Type:               params.Type,
			ImagePath:          newFilePath,
			ConfirmationStatus: sql.NullString{String: ProfileImageStatusProcessing, Valid: true},
			OriginalFileName:   originalName,
			IsBack:             sql.NullBool{Bool: false, Valid: true},
		}
		if params.IDCardCode != "" {
			frontProfileImage.IDCardCode = sql.NullString{String: params.IDCardCode, Valid: true}
		}
		if params.SubType != "" {
			frontProfileImage.SubType = sql.NullString{String: params.SubType, Valid: true}
		}
		err = tx.Omit(clause.Associations).Save(frontProfileImage).Error
		if err != nil {
			tx.Rollback()
			s.logger.Error2("can not save front image in db", err,
				zap.String("service", "userService"),
				zap.String("method", "UploadImages"),
				zap.Int("userID", u.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
		if params.Type == ProfileImageTypeAddress {
			up.AddressConfirmationStatus = sql.NullString{String: ProfileImageStatusIncomplete, Valid: true}
		} else {
			up.IdentityConfirmationStatus = sql.NullString{String: ProfileImageStatusIncomplete, Valid: true}
		}
		up.LastUploadedImageDate = time.Now()
		err = tx.Omit(clause.Associations).Save(up).Error
		if err != nil {
			tx.Rollback()
			s.logger.Error2("can not update user profile", err,
				zap.String("service", "userService"),
				zap.String("method", "UploadImages"),
				zap.Int("userID", u.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
		imagePath := s.configs.GetImagePath() + frontProfileImage.ImagePath
		res["frontImage"] = UploadImageSingleResponse{
			ID:   frontProfileImage.ID,
			Path: imagePath,
		}
	}

	if params.BackHeader != nil {
		if params.BackHeader.Size > MaxUploadSize {
			return response.Error("max upload size is 5M", http.StatusUnprocessableEntity, nil)
		}
		backImageBuff := make([]byte, 512)
		_, err := params.BackImage.Read(backImageBuff)
		if err != nil {
			tx.Rollback()
			s.logger.Error2("error reading fornt image buff", err,
				zap.String("service", "userService"),
				zap.String("method", "UploadImages"),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
		fileType := http.DetectContentType(backImageBuff)
		if fileType != "image/jpeg" && fileType != "image/png" && fileType != "application/pdf" {
			return response.Error("file type must be jpeg, png or pdf", http.StatusUnprocessableEntity, nil)
		}
		//When we read the first 512 bytes of the uploaded file in order to determine the content type,
		//the underlying file stream pointer moves forward by 512 bytes. When io.Copy() is called later,
		//it continues reading from that position resulting in a corrupted image file. The file.Seek()
		//method is used to return the pointer back to the start of the file so that io.Copy() starts from the beginning.
		_, err = params.BackImage.Seek(0, io.SeekStart)
		if err != nil {
			tx.Rollback()
			s.logger.Error2("error seeking the backImage", err,
				zap.String("service", "userService"),
				zap.String("method", "UploadImages"),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
		if frontProfileImage == nil {
			if params.FrontImageID == 0 {
				return response.Error("front image is not sent", http.StatusUnprocessableEntity, nil)
			}
			frontProfileImage = &ProfileImage{}
			err := s.profileImageRepository.GetImageByID(params.FrontImageID, frontProfileImage)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				s.logger.Error2("error getting frontImage by id", err,
					zap.String("service", "userService"),
					zap.String("method", "UploadImages"),
					zap.Int("userID", u.ID),
					zap.Int64("frontImageId", params.FrontImageID),
				)
				return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
			}
			if errors.Is(err, gorm.ErrRecordNotFound) || frontProfileImage.UserProfileID.Int64 != up.ID || frontProfileImage.IsBack.Bool {
				tx.Rollback()
				return response.Error("front image id not found", http.StatusUnprocessableEntity, nil)
			}
		}
		originalName := params.BackHeader.Filename
		newFilePath := UserProfileImagePathPrefix + uuid.NewString() + filepath.Ext(params.BackHeader.Filename)
		err = os.MkdirAll("."+UserProfileImagePathPrefix, 0700)
		if err != nil {
			tx.Rollback()
			s.logger.Error2("can not create directory for images", err,
				zap.String("service", "userService"),
				zap.String("method", "UploadImages"),
				zap.Int("userID", u.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
		dest, err := os.Create("." + newFilePath)
		if err != nil {
			tx.Rollback()
			s.logger.Error2("can not create back image new file path", err,
				zap.String("service", "userService"),
				zap.String("method", "UploadImages"),
				zap.Int("userID", u.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
		_, err = io.Copy(dest, params.BackImage)
		if err != nil {
			tx.Rollback()
			s.logger.Error2("can not copy back image", err,
				zap.String("service", "userService"),
				zap.String("method", "UploadImages"),
				zap.Int("userID", u.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
		backProfileImage := &ProfileImage{
			UserProfileID:      sql.NullInt64{Int64: up.ID, Valid: true},
			Type:               params.Type,
			ImagePath:          newFilePath,
			ConfirmationStatus: sql.NullString{String: ProfileImageStatusProcessing, Valid: true},
			OriginalFileName:   originalName,
			IsBack:             sql.NullBool{Bool: true, Valid: true},
			MainImageID:        sql.NullInt64{Int64: frontProfileImage.ID, Valid: true},
		}
		if params.IDCardCode != "" {
			backProfileImage.IDCardCode = sql.NullString{String: params.IDCardCode, Valid: true}
		}
		if params.SubType != "" {
			backProfileImage.SubType = sql.NullString{String: params.SubType, Valid: true}
		}
		err = tx.Omit(clause.Associations).Save(backProfileImage).Error
		if err != nil {
			tx.Rollback()
			s.logger.Error2("can not save back image in db", err,
				zap.String("service", "userService"),
				zap.String("method", "UploadImages"),
				zap.Int("userID", u.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}

		if params.Type == ProfileImageTypeAddress {
			up.AddressConfirmationStatus = sql.NullString{String: ProfileImageStatusProcessing, Valid: true}
		} else {
			up.IdentityConfirmationStatus = sql.NullString{String: ProfileImageStatusProcessing, Valid: true}
		}

		up.LastUploadedImageDate = time.Now()
		err = tx.Omit(clause.Associations).Save(up).Error
		if err != nil {
			tx.Rollback()
			s.logger.Error2("can not update user profile", err,
				zap.String("service", "userService"),
				zap.String("method", "UploadImages"),
				zap.Int("userID", u.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
		imagePath := s.configs.GetImagePath() + frontProfileImage.ImagePath
		res["backImage"] = UploadImageSingleResponse{
			ID:   backProfileImage.ID,
			Path: imagePath,
		}
	}
	if params.BackHeader == nil && frontProfileImage != nil && params.BackImageID != 0 {
		backProfileImage := &ProfileImage{}
		err := s.profileImageRepository.GetImageByID(params.BackImageID, backProfileImage)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			s.logger.Error2("error getting backImage by id", err,
				zap.String("service", "userService"),
				zap.String("method", "UploadImages"),
				zap.Int("userID", u.ID),
				zap.Int64("backImageId", params.BackImageID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
		if errors.Is(err, gorm.ErrRecordNotFound) || backProfileImage.UserProfileID.Int64 != up.ID || !backProfileImage.IsBack.Bool {
			tx.Rollback()
			return response.Error("back image id not found", http.StatusUnprocessableEntity, nil)
		}
		backProfileImage.MainImageID = sql.NullInt64{Int64: frontProfileImage.ID, Valid: true}
		err = tx.Omit(clause.Associations).Save(backProfileImage).Error
		if err != nil {
			tx.Rollback()
			s.logger.Error2("can not save back image in db", err,
				zap.String("service", "userService"),
				zap.String("method", "UploadImages"),
				zap.Int("userID", u.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not commit transaction", err,
			zap.String("service", "userService"),
			zap.String("method", "UploadImages"),
			zap.Int("userID", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	return response.Success(res, "")
}

func (s *service) checkIfUserCanUpload(up Profile, imageType string) string {
	if up.Status.String == ProfileImageStatusIncomplete {
		return "user profile information is not set yet"
	}
	if up.Status.String == ProfileImageStatusConfirmed {
		return "user profile information is already confirmed"
	}
	if imageType == ProfileImageTypeAddress {
		if up.AddressConfirmationStatus.String == ProfileStatusProcessing || up.AddressConfirmationStatus.String == ProfileStatusConfirmed {
			return "address is already in processing or confirmed"
		}
	}
	if imageType == ProfileImageTypeIdentity {
		if up.IdentityConfirmationStatus.String == ProfileStatusProcessing || up.IdentityConfirmationStatus.String == ProfileStatusConfirmed {
			return "address is already in processing or confirmed"
		}
	}
	return ""
}

func (s *service) DeleteImage(u *User, params DeleteImageParams) (apiResponse response.APIResponse, statusCode int) {
	up, err := s.GetUserProfile(*u)
	if err != nil {
		s.logger.Error2("can not get user profile", err,
			zap.String("service", "userService"),
			zap.String("method", "DeleteImage"),
			zap.Int("userID", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	pi := &ProfileImage{}
	err = s.profileImageRepository.GetImageByID(params.ID, pi)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("can not get image by id", err,
			zap.String("service", "userService"),
			zap.String("method", "DeleteImage"),
			zap.Int64("imageID", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || pi.UserProfileID.Int64 != up.ID {
		return response.Error("image not found", http.StatusUnprocessableEntity, nil)
	}
	if pi.ConfirmationStatus.String != ProfileImageStatusProcessing {
		return response.Error("can not delete already processed", http.StatusUnprocessableEntity, nil)
	}
	pi.IsDeleted = sql.NullBool{Bool: true, Valid: true}
	err = s.db.Omit(clause.Associations).Save(pi).Error
	if err != nil {
		s.logger.Error2("can not save profile image", err,
			zap.String("service", "userService"),
			zap.String("method", "DeleteImage"),
			zap.Int64("imageId", pi.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	res := make(map[string]string)
	return response.Success(res, "")
}

func (s *service) GetUsersByPagination(page int64, pageSize int, filters map[string]interface{}) []User {
	return s.userRepository.GetUsersByPagination(page, pageSize, filters)
}

func NewUserService(db *gorm.DB, userRepo Repository, profileRepo ProfileRepository, profileImageRepository ProfileImageRepository,
	countryService country.Service, twoFaManager TwoFaManager, passwordEncoder platform.PasswordEncoder,
	communicationService communication.Service, phoneConfirmationManager PhoneConfirmationManager,
	jwtService jwt.Service, configs platform.Configs, logger platform.Logger) Service {

	return &service{
		db:                       db,
		userRepository:           userRepo,
		userProfileRepository:    profileRepo,
		profileImageRepository:   profileImageRepository,
		countryService:           countryService,
		twoFaManager:             twoFaManager,
		passwordEncoder:          passwordEncoder,
		communicationService:     communicationService,
		phoneConfirmationManager: phoneConfirmationManager,
		jwtService:               jwtService,
		configs:                  configs,
		logger:                   logger,
	}

}
