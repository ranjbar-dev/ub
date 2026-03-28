package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"exchange-go/internal/api"
	"exchange-go/internal/configuration"
	"exchange-go/internal/di"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ConfigurationTests struct {
	*suite.Suite
	httpServer http.Handler
	db         *gorm.DB
}

func (t *ConfigurationTests) SetupSuite() {
	container := getContainer()
	t.httpServer = container.Get(di.HTTPServer).(api.HTTPServer).GetEngine()
	t.db = getDb()
}

func (t *ConfigurationTests) SetupTest() {}

func (t *ConfigurationTests) TearDownTest() {}

func (t *ConfigurationTests) TearDownSuite() {

}

func (t *ConfigurationTests) TestRecaptchaKey() {
	//for website
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/main-data/common", nil)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	response := struct {
		Status  bool
		Message string
		Data    configuration.GetRecaptchaKeyResponse
	}{}
	err := json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "site_key", response.Data.Key) // this key is in config.yaml

	//for android
	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/main-data/common", nil)
	req.Header.Set("User-Agent", "ubandroidv1.2.3")
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	response = struct {
		Status  bool
		Message string
		Data    configuration.GetRecaptchaKeyResponse
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "android_site_key", response.Data.Key) // this key is in config.yaml

}

func (t *ConfigurationTests) TestGetAppVersion() {
	//insert data in db
	av := &configuration.AppVersion{
		Version:     "1.2.7",
		VersionCode: 1.27,
		Platform:    "android",
		ForceUpdate: true,
		BugFixes:    sql.NullString{String: `["b1","b2"]`, Valid: true},
		KeyFeatures: sql.NullString{String: `["f1","f2"]`, Valid: true},
		ReleaseDate: time.Now(),
		StoreURL:    "testurl.com",
	}

	err := t.db.Create(av).Error
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()

	queryParams := url.Values{}
	queryParams.Set("platform", "android")
	queryParams.Set("current_version", "1.2.3")
	paramsString := queryParams.Encode()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/main-data/version?"+paramsString, nil)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	response := struct {
		Status  bool
		Message string
		Data    []configuration.GetAppVersionResponse
	}{}

	err = json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "testurl.com", response.Data[0].URL)
	assert.Equal(t.T(), "1.2.7", response.Data[0].Version)
	assert.Equal(t.T(), true, response.Data[0].ForceUpdate)
	assert.Equal(t.T(), "b1", response.Data[0].BugFixes[0])
	assert.Equal(t.T(), "b2", response.Data[0].BugFixes[1])
	assert.Equal(t.T(), "f1", response.Data[0].KeyFeatures[0])
	assert.Equal(t.T(), "f2", response.Data[0].KeyFeatures[1])
}

func (t *ConfigurationTests) TestContactUs() {
	res := httptest.NewRecorder()
	data := `{"name":"test","email":"test@test.com","subject":"test","body":"test"}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/main-data/contact-us", bytes.NewReader(body))
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)
}

func TestConfigurations(t *testing.T) {
	suite.Run(t, &ConfigurationTests{
		Suite: new(suite.Suite),
	})

}
