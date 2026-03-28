package repository

import (
	"context"
	"exchange-go/internal/currency"
	"exchange-go/internal/platform"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type pairRepository struct {
	db    *gorm.DB
	cache platform.Cache
}

func (pr *pairRepository) GetPairByID(id int64, pair *currency.Pair) error {
	cacheKey := fmt.Sprintf("pair:%d", id)
	ctx := context.Background()
	if err := pr.cache.Get(ctx, cacheKey, pair); err == nil {
		return nil
	}
	err := pr.db.Joins("BasisCoin").Joins("DependentCoin").Where(currency.Pair{ID: id}).First(pair).Error
	if err == nil {
		_ = pr.cache.Set(ctx, cacheKey, pair, time.Duration(1*time.Minute), nil)
	}
	return err
}

func (pr *pairRepository) GetPairByName(name string, pair *currency.Pair) error {
	cacheKey := fmt.Sprintf("pair:%s", name)
	ctx := context.Background()
	if err := pr.cache.Get(ctx, cacheKey, pair); err == nil {
		return nil
	}
	err := pr.db.Joins("BasisCoin").Joins("DependentCoin").Where(currency.Pair{Name: name}).First(pair).Error
	if err == nil {
		_ = pr.cache.Set(ctx, cacheKey, pair, time.Duration(1*time.Minute), nil)
	}
	return err
}

func (pr *pairRepository) GetActivePairCurrenciesList() []currency.Pair {
	var pairs []currency.Pair
	pr.db.Joins("BasisCoin").Joins("DependentCoin").Where(currency.Pair{IsActive: true}).Find(&pairs)
	return pairs
}

func (pr *pairRepository) GetAllPairs() []currency.Pair {
	var pairs []currency.Pair
	pr.db.Joins("BasisCoin").Joins("DependentCoin").Find(&pairs)
	return pairs
}

func (pr *pairRepository) GetPairsByName(names []string) []currency.Pair {
	cacheKey := "pairsByName:" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(names)), ","), "[]")
	ctx := context.Background()
	var pairs []currency.Pair
	if err := pr.cache.Get(ctx, cacheKey, &pairs); err == nil {
		return pairs
	}
	pr.db.Joins("BasisCoin").Joins("DependentCoin").Where("pair_currencies.name in ?", names).Find(&pairs)
	_ = pr.cache.Set(ctx, cacheKey, &pairs, time.Duration(20*time.Minute), nil)
	return pairs
}

func NewPairRepository(db *gorm.DB, cache platform.Cache) currency.PairRepository {
	return &pairRepository{
		db:    db,
		cache: cache,
	}
}
