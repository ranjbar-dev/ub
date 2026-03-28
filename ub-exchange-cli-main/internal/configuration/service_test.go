// Package configuration_test tests the configuration service. Covers:
//   - Retrieving recaptcha keys with platform-specific selection (Android vs web)
//   - Fetching app version updates with force-update flags, key features, and bug fixes
//   - Sending contact-us requests via the communication service
//
// Test data: mock configuration and app version repositories, communication service,
// and config provider with recaptcha keys and support email fixtures.
package configuration_test

import (
	"database/sql"
	"exchange-go/internal/communication"
	"exchange-go/internal/configuration"
	"exchange-go/internal/mocks"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetRecaptchaKey(t *testing.T) {
	configurationRepository := new(mocks.ConfigurationRepository)
	appVersionRepository := new(mocks.AppVersionRepository)
	communicationService := new(mocks.CommunicationService)
	configs := new(mocks.Configs)
	configs.On("GetAndroidRecaptchaSiteKey").Once().Return("androidRecaptchaKey")
	configs.On("GetRecaptchaSiteKey").Once().Return("siteRecaptchaKey")

	logger := new(mocks.Logger)
	configurationService := configuration.NewConfigurationService(configurationRepository, appVersionRepository, communicationService, configs, logger)

	params := configuration.GetRecaptchaKeyParams{
		UserAgent: "ubandroidv1.0.2",
	}
	res, statusCode := configurationService.GetRecaptchaKey(params)
	assert.Equal(t, http.StatusOK, statusCode)
	recaptcha, ok := res.Data.(configuration.GetRecaptchaKeyResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}
	assert.Equal(t, "androidRecaptchaKey", recaptcha.Key)

	params = configuration.GetRecaptchaKeyParams{
		UserAgent: "moziila",
	}
	res, statusCode = configurationService.GetRecaptchaKey(params)
	assert.Equal(t, http.StatusOK, statusCode)
	recaptcha, ok = res.Data.(configuration.GetRecaptchaKeyResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}
	assert.Equal(t, "siteRecaptchaKey", recaptcha.Key)

	configs.AssertExpectations(t)
}

func TestService_GetAppVersion(t *testing.T) {
	configurationRepository := new(mocks.ConfigurationRepository)
	appVersionRepository := new(mocks.AppVersionRepository)

	appVersions := []configuration.AppVersion{
		{
			ID:          1,
			Version:     "1.3.4",
			VersionCode: 1.34,
			Platform:    "",
			ForceUpdate: true,
			KeyFeatures: sql.NullString{String: `["f1","f2"]`, Valid: true},
			ReleaseDate: time.Time{},
			BugFixes:    sql.NullString{String: `["b1","b2"]`, Valid: true},
			StoreURL:    "storeurl.com",
		},
		{
			ID:          2,
			Version:     "1.3.5",
			VersionCode: 1.35,
			Platform:    "",
			ForceUpdate: false,
			KeyFeatures: sql.NullString{String: `["f1","f2"]`, Valid: true},
			ReleaseDate: time.Time{},
			BugFixes:    sql.NullString{String: `["b1","b2"]`, Valid: true},
			StoreURL:    "storeurl.com",
		},
	}
	appVersionRepository.On("FindNewAppVersions", "android", 1.23, mock.Anything).Once().Return(appVersions, nil)
	communicationService := new(mocks.CommunicationService)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	configurationService := configuration.NewConfigurationService(configurationRepository, appVersionRepository, communicationService, configs, logger)

	params := configuration.AppVersionParams{
		Platform: "android",
		Version:  "1.2.3",
	}
	res, statusCode := configurationService.GetAppVersion(params)
	assert.Equal(t, http.StatusOK, statusCode)
	result, ok := res.Data.([]configuration.GetAppVersionResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	assert.Equal(t, "1.3.4", result[0].Version)
	assert.Equal(t, "b1", result[0].BugFixes[0])
	assert.Equal(t, "b2", result[0].BugFixes[1])
	assert.Equal(t, "f1", result[0].KeyFeatures[0])
	assert.Equal(t, "f2", result[0].KeyFeatures[1])
	assert.Equal(t, true, result[0].ForceUpdate)
	assert.Equal(t, "storeurl.com", result[0].URL)

	assert.Equal(t, "1.3.5", result[1].Version)
	assert.Equal(t, "b1", result[1].BugFixes[0])
	assert.Equal(t, "b2", result[1].BugFixes[1])
	assert.Equal(t, "f1", result[1].KeyFeatures[0])
	assert.Equal(t, "f2", result[1].KeyFeatures[1])
	assert.Equal(t, false, result[1].ForceUpdate)
	assert.Equal(t, "storeurl.com", result[0].URL)

	appVersionRepository.AssertExpectations(t)
}

func TestService_ContactUs(t *testing.T) {
	configurationRepository := new(mocks.ConfigurationRepository)
	appVersionRepository := new(mocks.AppVersionRepository)
	communicationService := new(mocks.CommunicationService)
	p := communication.ContactUsToAdminParams{
		Name:    "name",
		Email:   "email",
		Subject: "subject",
		Body:    "body",
	}

	adminUser := communication.CommunicatingUser{
		Email: "support@unitedbit.com",
		Phone: "",
	}
	communicationService.On("SendContactUsToAdmin", adminUser, p).Once().Return()
	configs := new(mocks.Configs)
	configs.On("GetString", "exchange.supportemail").Once().Return("support@unitedbit.com")
	logger := new(mocks.Logger)

	configurationService := configuration.NewConfigurationService(configurationRepository, appVersionRepository, communicationService, configs, logger)

	params := configuration.ContactUsParams{
		Name:    "name",
		Email:   "email",
		Subject: "subject",
		Body:    "body",
	}

	_, statusCode := configurationService.ContactUs(params)
	assert.Equal(t, http.StatusOK, statusCode)

	time.Sleep(20 * time.Millisecond)
	communicationService.AssertExpectations(t)

}
