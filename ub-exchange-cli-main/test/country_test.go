package test

import (
	"database/sql"
	"encoding/json"
	"exchange-go/internal/api"
	"exchange-go/internal/country"
	"exchange-go/internal/di"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type CountryTests struct {
	*suite.Suite
	httpServer http.Handler
	db         *gorm.DB
}

func (t *CountryTests) SetupSuite() {
	container := getContainer()
	t.httpServer = container.Get(di.HTTPServer).(api.HTTPServer).GetEngine()
	t.db = getDb()

}

func (t *CountryTests) SetupTest() {}

func (t *CountryTests) TearDownTest() {}

func (t *CountryTests) TearDownSuite() {
	err := t.db.Where("id > ?", int64(0)).Delete(country.Country{}).Error
	if err != nil {
		t.Fail(err.Error())
	}
}

func (t *CountryTests) TestCountries() {
	c := []country.Country{
		{
			ID:        1,
			Name:      sql.NullString{String: "co1", Valid: true},
			FullName:  sql.NullString{String: "country1", Valid: true},
			Code:      sql.NullString{String: "01", Valid: true},
			ImagePath: sql.NullString{String: "/images", Valid: true},
		},
		{
			ID:        2,
			Name:      sql.NullString{String: "co2", Valid: true},
			FullName:  sql.NullString{String: "country2", Valid: true},
			Code:      sql.NullString{String: "02", Valid: true},
			ImagePath: sql.NullString{String: "/images", Valid: true},
		},
		{
			ID:        3,
			Name:      sql.NullString{String: "co3", Valid: true},
			FullName:  sql.NullString{String: "country3", Valid: true},
			Code:      sql.NullString{String: "03", Valid: true},
			ImagePath: sql.NullString{String: "/images", Valid: true},
		},
		{
			ID:        4,
			Name:      sql.NullString{String: "co4", Valid: true},
			FullName:  sql.NullString{String: "country4", Valid: true},
			Code:      sql.NullString{String: "04", Valid: true},
			ImagePath: sql.NullString{String: "/images", Valid: true},
		},
	}
	err := t.db.Create(&c).Error
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/main-data/country-list", nil)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	response := struct {
		Status  bool
		Message string
		Data    []country.GetCountryResponse
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		t.Fail(err.Error())
	}

	for i, co := range response.Data {
		assert.Equal(t.T(), int64(i+1), co.ID)
		assert.Equal(t.T(), "co"+strconv.Itoa(i+1), co.Name)
		assert.Equal(t.T(), "country"+strconv.Itoa(i+1), co.FullName)
		assert.Equal(t.T(), "0"+strconv.Itoa(i+1), co.Code)
	}
}

func TestCountry(t *testing.T) {
	suite.Run(t, &CountryTests{
		Suite: new(suite.Suite),
	})

}
