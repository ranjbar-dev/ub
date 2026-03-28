package user

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                     int
	Email                  string
	Password               string
	Kyc                    int
	Status                 string
	AccountStatus          string
	ExchangeNumber         int64 `gorm:"column:number_of_exchange"`
	Google2faSecretCode    sql.NullString
	Phone                  sql.NullString
	IsTwoFaEnabled         bool
	UbID                   string
	VerificationCode       string
	Google2faDisabledAt    sql.NullTime
	PrivateChannelName     string
	IsLevelManuallySet     sql.NullBool
	CreatedAt              time.Time
	UpdatedAt              time.Time
	ExchangeVolumeCoinCode string `gorm:"column:exchange_volume_currency;default:BTC"`
	ExchangeVolumeAmount   string
	UserLevelID            int64
	UserLevel              Level
	RefreshToken           sql.NullString
	RefreshTokenExpiry     time.Time
	PasswordChangedAt      sql.NullTime
	TwoFaChangedAt         sql.NullTime
}

type Role struct {
	ID   int64
	Name string
	Role string
}

type UserRole struct {
	UserID int64
	RoleID int64
}

func (UserRole) TableName() string {
	return "user_role"
}

type Profile struct {
	ID                         int64
	UserID                     int
	FirstName                  sql.NullString
	LastName                   sql.NullString
	CountryID                  sql.NullInt64
	Gender                     sql.NullString
	DateOfBirth                sql.NullString
	Address                    sql.NullString
	RegionAndCity              sql.NullString
	PostalCode                 sql.NullString
	IDCardCode                 sql.NullString
	Status                     sql.NullString
	AdminComment               sql.NullString
	IdentityConfirmationStatus sql.NullString
	AddressConfirmationStatus  sql.NullString
	PhoneConfirmationType      sql.NullString
	RegistrationIP             sql.NullString
	ReferKey                   sql.NullString
	TrustLevel                 int64
	AvatarImagePath            sql.NullString
	LastUploadedImageDate      time.Time
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}

func (p Profile) GetFullName() string {
	return p.FirstName.String + " " + p.LastName.String
}

func (Profile) TableName() string {
	return "user_profiles"
}

type Level struct {
	ID                        int64
	Name                      string
	MakerFeePercentage        float64
	TakerFeePercentage        float64
	WithdrawFeePercentage     float64
	DepositFeePercentage      float64
	Code                      int64
	ExchangeNumberLimit       int64
	ExchangeVolumeLimitAmount string
	ExchangeVolumeLimitCoin   string `gorm:"column:exchange_volume_limit_currency"`
	MinExchangeVolume         float64
	MaxExchangeVolume         float64
	MinKycLevel               int
}

func (Level) TableName() string {
	return "user_levels"
}

type Permission struct {
	ID   int64
	Name string
}

func (Permission) TableName() string {
	return "user_permissions"
}

type UsersPermissions struct {
	UserID           int
	UserPermissionID int64
	Permission       Permission `gorm:"foreignKey:UserPermissionID"`
}

func (UsersPermissions) TableName() string {
	return "users_permissions"
}

type Config struct {
	ID                                    int64
	UserID                                int
	Theme                                 string
	Mode                                  string
	CustomThemeData                       sql.NullString
	IsTradeNotificationEnabled            bool
	IsEmailVerificationForWithdrawEnabled bool
	IsTwoFaVerificationForWithdrawEnabled bool
	IsEmailVerificationForLoginEnabled    bool
	IsTwoFaVerificationForLoginEnabled    bool
	IsWhiteListEnabled                    bool
	IsReadOnly                            bool
}

func (Config) TableName() string {
	return "user_configs"
}

// Repository provides data access methods for user entities.
type Repository interface {
	// GetUserByUsername retrieves an active user by their username (email).
	GetUserByUsername(username string, user *User) error
	// GetEvenBlockedUserByEmail retrieves a user by email regardless of account block status.
	GetEvenBlockedUserByEmail(email string, user *User) error
	// GetUserByIDUsingTx retrieves a user by ID within the provided database transaction.
	GetUserByIDUsingTx(tx *gorm.DB, ID int, user *User) error
	// GetUserByID retrieves a user by their unique identifier.
	GetUserByID(ID int, user *User) error
	// GetAdminUserByUsername retrieves an admin user by their username (email).
	GetAdminUserByUsername(username string, user *User) error
	// GetUsersByPagination returns a paginated list of users matching the given filters.
	GetUsersByPagination(page int64, pageSize int, filters map[string]interface{}) []User
	// GetUserByVerificationCode retrieves a user by their email verification code.
	GetUserByVerificationCode(code string, u *User) error
	// GetUserByRefreshToken retrieves a user by their JWT refresh token.
	GetUserByRefreshToken(refreshToken string, user *User) error
	// GetUsersDataForOrderMatching returns lightweight user data needed for order matching by user IDs.
	GetUsersDataForOrderMatching(userIds []int) []UsersDataForOrderMatching
}

// UsersPermissionsRepository provides data access for user-permission assignments.
type UsersPermissionsRepository interface {
	// GetUserPermissions retrieves all permission assignments for the given user.
	GetUserPermissions(userID int) []UsersPermissions
}

// PermissionRepository provides data access for permission definitions.
type PermissionRepository interface {
	// GetAllPermissions retrieves all defined permissions in the system.
	GetAllPermissions() []Permission
}

// LevelRepository provides data access for user VIP level definitions.
type LevelRepository interface {
	// GetAllLevels retrieves all VIP levels ordered by code.
	GetAllLevels() []Level
	// GetLevelByID retrieves a VIP level by its unique identifier.
	GetLevelByID(id int64, level *Level) error
	// GetLevelsByIds retrieves multiple VIP levels by their IDs.
	GetLevelsByIds(ids []int64) []Level
	// GetLevelByCode retrieves a VIP level by its numeric code (e.g., 0–8).
	GetLevelByCode(code int64, level *Level) error
}

// ConfigRepository provides data access for user configuration preferences.
type ConfigRepository interface {
	// GetUserConfigByUserID retrieves the configuration settings for the given user.
	GetUserConfigByUserID(userID int, config *Config) error
}

// ProfileRepository provides data access for user profile records.
type ProfileRepository interface {
	// GetProfileByUserID retrieves the profile for the given user.
	GetProfileByUserID(userID int, profile *Profile) error
	// GetProfileByUserIDUsingTx retrieves the profile for the given user within a database transaction.
	GetProfileByUserIDUsingTx(tx *gorm.DB, userID int, profile *Profile) error
}

type ProfileImage struct {
	ID                 int64
	UserProfileID      sql.NullInt64
	Type               string
	ImagePath          string
	ConfirmationStatus sql.NullString
	OriginalFileName   string
	IDCardCode         sql.NullString
	RejectionReason    sql.NullString
	CreatedAt          time.Time
	UpdatedAt          time.Time
	SubType            sql.NullString
	IsDeleted          sql.NullBool
	MainImageID        sql.NullInt64
	IsBack             sql.NullBool
}

func (ProfileImage) TableName() string {
	return "user_profile_image"
}

// ProfileImageRepository provides data access for KYC document images.
type ProfileImageRepository interface {
	// GetImagesByIds retrieves multiple profile images by their IDs.
	GetImagesByIds(ids []int64) []ProfileImage
	// GetImageByID retrieves a single profile image by its unique identifier.
	GetImageByID(id int64, upi *ProfileImage) error
	// GetLatestImagesDataByProfileID retrieves the most recent image metadata for a user profile.
	GetLatestImagesDataByProfileID(profileID int64) []ImagesQueryFields
}

type LoginHistory struct {
	ID        int64
	UserID    sql.NullInt64
	Device    sql.NullString
	IP        sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     sql.NullString
	Type      string
	Password  sql.NullString
}

func (LoginHistory) TableName() string {
	return "user_login_history"
}

// LoginHistoryRepository provides data access for user login history records.
type LoginHistoryRepository interface {
	// Create persists a new login history entry.
	Create(loginHistory *LoginHistory) error
	// GetLastLoginHistoryByUserID retrieves the most recent login record for the given user.
	GetLastLoginHistoryByUserID(userID int, loginHistory *LoginHistory) error
}
