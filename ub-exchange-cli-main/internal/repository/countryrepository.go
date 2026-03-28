package repository

import (
	"context"
	"exchange-go/internal/country"
	"exchange-go/internal/platform"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type countryRepository struct {
	db    *gorm.DB
	cache platform.Cache
}

func (cr *countryRepository) All(ctx context.Context) []country.Country {
	cacheKey := "allCountries"
	var countries []country.Country
	if err := cr.cache.Get(ctx, cacheKey, &countries); err == nil {
		return countries
	}
	cr.db.Find(&countries)
	_ = cr.cache.Set(ctx, cacheKey, &countries, time.Duration(24*time.Hour), nil)
	return countries
}

func (cr *countryRepository) GetCountryByID(id int64, c *country.Country) error {
	cacheKey := fmt.Sprintf("country_by_id:%d", id)
	ctx := context.Background()
	if err := cr.cache.Get(ctx, cacheKey, c); err == nil {
		return nil
	}
	err := cr.db.Where(country.Country{ID: id}).First(c).Error
	_ = cr.cache.Set(ctx, cacheKey, c, time.Duration(1*time.Hour), nil)
	return err
}

func NewCountryRepository(db *gorm.DB, cache platform.Cache) country.Repository {
	return &countryRepository{db, cache}
}
