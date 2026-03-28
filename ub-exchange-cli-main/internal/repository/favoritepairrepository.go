package repository

import (
	"exchange-go/internal/currency"

	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

type favoritePairRepository struct {
	db *gorm.DB
}

func (r *favoritePairRepository) Create(favoritePair *currency.FavoritePair) error {
	return r.db.Omit(clause.Associations).Create(favoritePair).Error
}

func (r *favoritePairRepository) Delete(favoritePair *currency.FavoritePair) error {
	return r.db.Where("user_id = ? and pair_currency_id =?", favoritePair.UserID, favoritePair.PairID).Delete(favoritePair).Error
}

func (r *favoritePairRepository) GetFavoritePair(userID int, pairID int64, favoritePair *currency.FavoritePair) error {
	return r.db.Where("user_id = ? and pair_currency_id = ?", userID, pairID).First(favoritePair).Error
}

func (r *favoritePairRepository) GetUserFavoritePairs(userID int) []currency.FavoritePairQueryFields {
	var result []currency.FavoritePairQueryFields
	r.db.Table("user_favorite_pair_currency as ufpc").
		Joins("join pair_currencies as p on p.id = ufpc.pair_currency_id").
		Select("ufpc.pair_currency_id as PairID, p.name as PairName").
		Where("ufpc.user_id = ?", userID).Scan(&result)
	return result
}

func NewFavoritePairRepository(db *gorm.DB) currency.FavoritePairRepository {
	return &favoritePairRepository{db}
}
