package repository

import (
	"context"
	"exchange-go/internal/configuration"
	"exchange-go/internal/platform"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type appVersionRepository struct {
	db    *gorm.DB
	cache platform.Cache
}

func (r *appVersionRepository) FindNewAppVersion(platform string, versionCode float64, appVersion *configuration.AppVersion) error {
	cacheKey := fmt.Sprintf("app_version:%f", versionCode)
	ctx := context.Background()
	if err := r.cache.Get(ctx, cacheKey, appVersion); err == nil {
		return nil
	}
	err := r.db.Where("version_code > ?", versionCode).Order("id desc").Limit(1).First(appVersion).Error
	_ = r.cache.Set(ctx, cacheKey, appVersion, time.Duration(3*time.Hour), nil)
	return err
}

func (r *appVersionRepository) FindNewAppVersions(platform string, versionCode float64) ([]configuration.AppVersion, error) {
	cacheKey := fmt.Sprintf("app_versions:%s-%f", platform, versionCode)
	ctx := context.Background()
	var appVersions []configuration.AppVersion
	if err := r.cache.Get(ctx, cacheKey, &appVersions); err == nil {
		return appVersions, nil
	}
	err := r.db.Where("version_code > ? AND platform = ?", versionCode, platform).Order("version_code asc").Find(&appVersions).Error
	_ = r.cache.Set(ctx, cacheKey, &appVersions, time.Duration(3*time.Hour), nil)
	return appVersions, err
}

func NewAppVersionRepository(db *gorm.DB, cache platform.Cache) configuration.AppVersionRepository {
	return &appVersionRepository{
		db:    db,
		cache: cache,
	}
}
