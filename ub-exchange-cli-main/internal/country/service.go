package country

import (
	"context"
	"exchange-go/internal/platform"
	"exchange-go/internal/response"
)

// Service provides the public API for country listing and lookup operations.
type Service interface {
	// GetCountries returns all countries formatted for the API response.
	GetCountries() (apiResponse response.APIResponse, statusCode int)
	// GetCountryByID retrieves a single country by its database ID.
	GetCountryByID(id int64) (Country, error)
}

type service struct {
	repo    Repository
	configs platform.Configs
}

type GetCountryResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"fullName"`
	Code     string `json:"code"`
	Image    string `json:"image"`
}

func (s *service) GetCountries() (apiResponse response.APIResponse, statusCode int) {
	ctx := context.Background()
	countries := s.repo.All(ctx)
	var allCountries = make([]GetCountryResponse, 0)
	domain := s.configs.GetImagePath()
	for _, country := range countries {
		pc := GetCountryResponse{
			ID:       country.ID,
			Name:     country.Name.String,
			FullName: country.FullName.String,
			Code:     country.Code.String,
			Image:    domain + country.ImagePath.String,
		}
		allCountries = append(allCountries, pc)
	}

	return response.Success(allCountries, "")
}

func (s *service) GetCountryByID(id int64) (Country, error) {
	c := Country{}
	err := s.repo.GetCountryByID(id, &c)
	return c, err
}

func NewCountryService(repo Repository, configs platform.Configs) Service {
	return &service{
		repo:    repo,
		configs: configs,
	}
}
