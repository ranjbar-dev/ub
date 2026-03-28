package repository

import (
	"exchange-go/internal/platform"
	"exchange-go/internal/user"

	"gorm.io/gorm"
)

type userRepository struct {
	db    *gorm.DB
	cache platform.Cache
}

func (ur *userRepository) GetUsersByPagination(page int64, pageSize int, filters map[string]interface{}) []user.User {
	var users []user.User
	offset := int(page) * pageSize
	ur.db.Where(filters).Offset(offset).Limit(pageSize).Order("id asc").Find(&users)
	return users
}

func (ur *userRepository) GetUserByUsername(username string, u *user.User) error {
	//cacheKey := fmt.Sprintf("user_by_username:%s", username)
	//ctx := context.Background()
	//if err := ur.cache.Get(ctx, cacheKey, u); err == nil {
	//return nil
	//}
	err := ur.db.Where(&user.User{Email: username, AccountStatus: user.AccountStatusUnblocked}).First(u).Error
	//_ = ur.cache.Set(ctx, cacheKey, u, time.Duration(3*time.Minute), nil)
	return err
}

func (ur *userRepository) GetEvenBlockedUserByEmail(email string, u *user.User) error {
	return ur.db.Where(&user.User{Email: email}).First(u).Error
}

func (ur *userRepository) GetUserByIDUsingTx(tx *gorm.DB, ID int, u *user.User) error {
	return tx.Where(&user.User{ID: ID}).First(u).Error
}

func (ur *userRepository) GetUserByID(ID int, u *user.User) error {
	return ur.db.Where(&user.User{ID: ID}).First(u).Error
}

func (ur *userRepository) GetAdminUserByUsername(username string, u *user.User) error {
	return ur.db.Joins("join user_role on user_role.user_id = users.id").
		Where(&user.User{Email: username}).First(u).Error
}

func (ur *userRepository) GetUserByVerificationCode(code string, u *user.User) error {
	return ur.db.Where(&user.User{VerificationCode: code}).First(u).Error
}

func (ur *userRepository) GetUserByRefreshToken(refreshToken string, u *user.User) error {
	return ur.db.Where("refresh_token = ?", refreshToken).First(u).Error
}

func (ur *userRepository) GetUsersDataForOrderMatching(userIds []int) []user.UsersDataForOrderMatching {
	var result []user.UsersDataForOrderMatching
	ur.db.Table("users").Where("users.id IN ?", userIds).Select("" +
		"users.id as UserID," +
		"users.email as UserEmail," +
		"users.user_level_id as UserLevelID," +
		"users.private_channel_name as UserPrivateChannel").
		Scan(&result)
	return result
}

func NewUserRepository(db *gorm.DB, cache platform.Cache) user.Repository {
	return &userRepository{
		db:    db,
		cache: cache,
	}
}
