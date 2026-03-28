package repository

import (
	"exchange-go/internal/user"

	"gorm.io/gorm"
)

type profileImageRepository struct {
	db *gorm.DB
}

func (r *profileImageRepository) GetImagesByIds(ids []int64) []user.ProfileImage {
	var images []user.ProfileImage
	r.db.Where("id in ?", ids).Find(&images)
	return images
}

func (r *profileImageRepository) GetImageByID(id int64, pi *user.ProfileImage) error {
	return r.db.Where(user.ProfileImage{ID: id}).First(pi).Error
}

func (r *profileImageRepository) GetLatestImagesDataByProfileID(profileID int64) []user.ImagesQueryFields {
	var data []user.ImagesQueryFields
	r.db.Raw("SELECT type as Type,MAX(id) as Id FROM user_profile_image WHERE user_profile_id = ? AND (is_deleted is null OR is_deleted = 0) group by Type,is_back", profileID).Scan(&data)
	return data
}

func NewProfileImageRepository(db *gorm.DB) user.ProfileImageRepository {
	return &profileImageRepository{db}
}
