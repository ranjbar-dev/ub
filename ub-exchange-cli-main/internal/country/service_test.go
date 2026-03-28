// Package country_test tests the country service. Covers:
//   - Listing all countries with ID, name, full name, code, and image path
//   - Retrieving a single country by ID with populated fields
//
// Test data: mock country repository and config provider with four country
// fixtures using sql.NullString fields.
package country_test

import (
	"database/sql"
	"exchange-go/internal/country"
	"exchange-go/internal/mocks"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetCountries(t *testing.T) {
	countryRepo := new(mocks.CountryRepository)
	countryModels := []country.Country{
		{
			ID:       1,
			Name:     sql.NullString{String: "co1", Valid: true},
			FullName: sql.NullString{String: "country1", Valid: true},
			Code:     sql.NullString{String: "01", Valid: true},
		},
		{
			ID:       2,
			Name:     sql.NullString{String: "co2", Valid: true},
			FullName: sql.NullString{String: "country2", Valid: true},
			Code:     sql.NullString{String: "02", Valid: true},
		},
		{
			ID:       3,
			Name:     sql.NullString{String: "co3", Valid: true},
			FullName: sql.NullString{String: "country3", Valid: true},
			Code:     sql.NullString{String: "03", Valid: true},
		},
		{
			ID:       4,
			Name:     sql.NullString{String: "co4", Valid: true},
			FullName: sql.NullString{String: "country4", Valid: true},
			Code:     sql.NullString{String: "04", Valid: true},
		},
	}
	countryRepo.On("All", mock.Anything).Once().Return(countryModels)

	configs := new(mocks.Configs)
	configs.On("GetImagePath").Once().Return("localhost")

	countryService := country.NewCountryService(countryRepo, configs)
	res, statusCode := countryService.GetCountries()
	assert.Equal(t, http.StatusOK, statusCode)

	countries, ok := res.Data.([]country.GetCountryResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	for i, co := range countries {
		assert.Equal(t, int64(i+1), co.ID)
		assert.Equal(t, "co"+strconv.Itoa(i+1), co.Name)
		assert.Equal(t, "country"+strconv.Itoa(i+1), co.FullName)
		assert.Equal(t, "0"+strconv.Itoa(i+1), co.Code)
	}
	countryRepo.AssertExpectations(t)
}

func TestService_GetCountryById(t *testing.T) {
	countryRepo := new(mocks.CountryRepository)

	countryModel := &country.Country{}

	countryRepo.On("GetCountryByID", int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		countryModel = args.Get(1).(*country.Country)
		countryModel.ID = 1
		countryModel.Name = sql.NullString{String: "co1", Valid: true}
		countryModel.FullName = sql.NullString{String: "country1", Valid: true}
		countryModel.Code = sql.NullString{String: "01", Valid: true}
	})

	configs := new(mocks.Configs)
	configs.On("GetImagePath").Once().Return("localhost")

	countryService := country.NewCountryService(countryRepo, configs)
	c, err := countryService.GetCountryByID(int64(1))
	assert.Nil(t, err)
	assert.Equal(t, int64(1), c.ID)
	assert.Equal(t, "co1", c.Name.String)
	assert.Equal(t, "country1", c.FullName.String)
	assert.Equal(t, "01", c.Code.String)

	countryRepo.AssertExpectations(t)
}
