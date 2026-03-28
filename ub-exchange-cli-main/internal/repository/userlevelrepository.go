package repository

import (
	"context"
	"exchange-go/internal/platform"
	"exchange-go/internal/user"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type userLevelRepository struct {
	db    *gorm.DB
	cache platform.Cache
}

func (ulr *userLevelRepository) GetAllLevels() []user.Level {
	cacheKey := "allLevels"
	ctx := context.Background()
	var levels []user.Level
	if err := ulr.cache.Get(ctx, cacheKey, &levels); err == nil {
		return levels
	}
	err := ulr.db.Order("id desc").Find(&levels).Error
	if err == nil {
		_ = ulr.cache.Set(ctx, cacheKey, &levels, time.Duration(1*time.Hour), nil)
	}
	return levels
}

func (ulr *userLevelRepository) GetLevelByID(id int64, level *user.Level) error {
	cacheKey := fmt.Sprintf("level:%d", id)
	ctx := context.Background()
	if err := ulr.cache.Get(ctx, cacheKey, level); err == nil {
		return nil
	}
	err := ulr.db.Where("id = ?", id).Find(&level).Error
	if err == nil {
		_ = ulr.cache.Set(ctx, cacheKey, level, time.Duration(1*time.Hour), nil)
	}
	return err
}

func (ulr *userLevelRepository) GetLevelByCode(code int64, level *user.Level) error {
	cacheKey := fmt.Sprintf("levelByCode:%d", code)
	ctx := context.Background()
	if err := ulr.cache.Get(ctx, cacheKey, level); err == nil {
		return nil
	}
	err := ulr.db.Where("code = ?", code).Find(&level).Error
	if err == nil {
		_ = ulr.cache.Set(ctx, cacheKey, level, time.Duration(1*time.Hour), nil)
	}
	return err
}

func (ulr *userLevelRepository) GetLevelsByIds(ids []int64) []user.Level {
	cacheKey := "levelByIds:" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")
	ctx := context.Background()
	var levels []user.Level
	if err := ulr.cache.Get(ctx, cacheKey, &levels); err == nil {
		return levels
	}
	err := ulr.db.Where("id IN ?", ids).Find(&levels).Error
	if err == nil {
		_ = ulr.cache.Set(ctx, cacheKey, &levels, time.Duration(30*time.Minute), nil)
	}
	return levels
}

func NewUserLevelRepository(db *gorm.DB, cache platform.Cache) user.LevelRepository {
	return &userLevelRepository{
		db:    db,
		cache: cache,
	}
}
