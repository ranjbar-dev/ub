package configuration

import (
	"database/sql"
	"time"
)

type Configuration struct {
	ID        int64
	Key       string `gorm:"column:k"`
	Value     string `gorm:"column:v"`
	GroupCode string
}

// Repository is a placeholder interface for configuration persistence operations.
type Repository interface {
}

type AppVersion struct {
	ID          int64
	Version     string
	VersionCode float64
	Platform    string
	ForceUpdate bool
	KeyFeatures sql.NullString
	ReleaseDate time.Time
	BugFixes    sql.NullString
	StoreURL    string
}

func (AppVersion) TableName() string {
	return "app_version"
}

// AppVersionRepository provides data access for mobile application version records.
type AppVersionRepository interface {
	// FindNewAppVersion finds a single newer app version for the given platform and version code.
	FindNewAppVersion(platform string, versionCode float64, appVersion *AppVersion) error
	// FindNewAppVersions returns all app versions newer than versionCode for the given platform.
	FindNewAppVersions(platform string, versionCode float64) ([]AppVersion, error)
}
