package configuration

import (
	"encoding/json"
	"errors"
	"exchange-go/internal/communication"
	"exchange-go/internal/platform"
	"exchange-go/internal/response"
	"exchange-go/internal/userdevice"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppVersionParams struct {
	Platform string `form:"platform" binding:"required,oneof='android' 'ios'"`
	Version  string `form:"current_version" binding:"required"`
}

type ContactUsParams struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

type GetRecaptchaKeyParams struct {
	UserAgent string
}

type GetRecaptchaKeyResponse struct {
	Key string `json:"recaptchaSiteKey"`
}

type GetAppVersionResponse struct {
	Version     string   `json:"version"`
	ForceUpdate bool     `json:"forceUpdate"`
	KeyFeatures []string `json:"keyFeatures"`
	BugFixes    []string `json:"bugFixes"`
	ReleaseDate string   `json:"releaseDate"`
	URL         string   `json:"url"`
}

// Service provides the public API for application configuration, including
// reCAPTCHA keys, app version checks, and the contact-us form.
type Service interface {
	// GetRecaptchaKey returns the appropriate reCAPTCHA site key based on the client's user agent.
	GetRecaptchaKey(params GetRecaptchaKeyParams) (apiResponse response.APIResponse, statusCode int)
	// GetAppVersion returns available app updates for the specified platform and version.
	GetAppVersion(params AppVersionParams) (apiResponse response.APIResponse, statusCode int)
	// ContactUs sends a contact-us message from a user to the admin support email.
	ContactUs(params ContactUsParams) (apiResponse response.APIResponse, statusCode int)
}

type service struct {
	repository           Repository
	appVersionRepository AppVersionRepository
	communicationService communication.Service
	configs              platform.Configs
	logger               platform.Logger
}

func (s *service) GetRecaptchaKey(params GetRecaptchaKeyParams) (apiResponse response.APIResponse, statusCode int) {
	device := userdevice.GetDeviceUsingUserAgent(params.UserAgent)
	key := ""
	if device == userdevice.DeviceAndroid {
		key = s.configs.GetAndroidRecaptchaSiteKey()
	} else {
		key = s.configs.GetRecaptchaSiteKey()
	}
	recaptcha := GetRecaptchaKeyResponse{Key: key}

	return response.Success(recaptcha, "")
}

func (s *service) GetAppVersion(params AppVersionParams) (apiResponse response.APIResponse, statusCode int) {
	result := make([]GetAppVersionResponse, 0)
	newAppVersions := make([]AppVersion, 0)
	versionString := params.Version
	lastDotIndex := strings.LastIndex(versionString, ".")
	version := versionString[:lastDotIndex] + versionString[lastDotIndex+1:]
	versionCode, err := strconv.ParseFloat(version, 64)
	if err != nil {
		s.logger.Error2("error converting version string to float", err,
			zap.String("service", "configurationService"),
			zap.String("method", "GetAppVersion"),
			zap.String("version", version),
			zap.String("platform", params.Platform),
		)
	} else {
		newAppVersions, err = s.appVersionRepository.FindNewAppVersions(params.Platform, versionCode)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error2("error finding new app versions", err,
				zap.String("service", "configurationService"),
				zap.String("method", "GetAppVersion"),
				zap.String("version", version),
				zap.String("platform", params.Platform),
			)
		}
	}

	for _, newAppVersion := range newAppVersions {
		keyFeatures := make([]string, 0)
		bugFixes := make([]string, 0)
		_ = json.Unmarshal([]byte(newAppVersion.KeyFeatures.String), &keyFeatures)
		_ = json.Unmarshal([]byte(newAppVersion.BugFixes.String), &bugFixes)

		newestVersion := GetAppVersionResponse{
			Version:     newAppVersion.Version,
			ForceUpdate: newAppVersion.ForceUpdate,
			KeyFeatures: keyFeatures,
			BugFixes:    bugFixes,
			ReleaseDate: newAppVersion.ReleaseDate.Format("2006-01-02"),
			URL:         newAppVersion.StoreURL,
		}
		result = append(result, newestVersion)
	}

	return response.Success(result, "")
}

func (s *service) ContactUs(params ContactUsParams) (apiResponse response.APIResponse, statusCode int) {
	mailParams := communication.ContactUsToAdminParams{
		Name:    params.Name,
		Email:   params.Email,
		Subject: params.Subject,
		Body:    params.Body,
	}

	cu := communication.CommunicatingUser{
		Email: s.configs.GetString("exchange.supportemail"),
		Phone: "",
	}
	platform.SafeGo(s.logger, "configuration.SendContactUsToAdmin", func() {
		s.communicationService.SendContactUsToAdmin(cu, mailParams)
	})
	return response.Success(nil, "")
}

func NewConfigurationService(repository Repository, appVersionRepository AppVersionRepository, communicationService communication.Service,
	configs platform.Configs, logger platform.Logger) Service {
	return &service{
		repository:           repository,
		appVersionRepository: appVersionRepository,
		communicationService: communicationService,
		configs:              configs,
		logger:               logger,
	}
}
