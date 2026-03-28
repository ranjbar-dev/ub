package user

import (
	"github.com/shopspring/decimal"
)

const (
	UserLevelVip0Code = 0
	UserLevelVip1Code = 1
	UserLevelVip2Code = 2
	UserLevelVip3Code = 3
	UserLevelVip4Code = 4
	UserLevelVip5Code = 5
	UserLevelVip6Code = 6
	UserLevelVip7Code = 7
	UserLevelVip8Code = 8
)

// LevelService provides VIP level management including level recalculation
// based on trading volume and level data access.
type LevelService interface {
	// RecalculateUserLevel computes the appropriate VIP level for a user based on
	// their KYC status and cumulative exchange volume.
	RecalculateUserLevel(u User, exchangeVolume string) Level
	// GetLevelByID retrieves a VIP level by its unique identifier.
	GetLevelByID(id int64) (Level, error)
	// GetLevelsByIds retrieves multiple VIP levels by their IDs.
	GetLevelsByIds(ids []int64) []Level
	// GetLevelByCode retrieves a VIP level by its numeric code (e.g., 0–8).
	GetLevelByCode(code int64) (Level, error)
}

type levelService struct {
	levelRepo LevelRepository
}

func (s *levelService) GetLevelByID(id int64) (Level, error) {
	level := Level{}
	err := s.levelRepo.GetLevelByID(id, &level)
	return level, err
}

func (s *levelService) GetLevelByCode(code int64) (Level, error) {
	level := Level{}
	err := s.levelRepo.GetLevelByCode(code, &level)
	return level, err
}

func (s *levelService) RecalculateUserLevel(u User, exchangeVolume string) Level {
	allLevels := s.levelRepo.GetAllLevels()
	kyc := u.Kyc
	var validLevels []Level

	for _, level := range allLevels {
		if kyc >= level.MinKycLevel {
			validLevels = append(validLevels, level)
		}
	}

	firstUserLevel := allLevels[0]
	lastUserLevel := allLevels[len(allLevels)-1]

	exchangeVolumeDecimal, _ := decimal.NewFromString(exchangeVolume)

	for _, validLevel := range validLevels {
		minExchangeDecimal := decimal.NewFromFloat(validLevel.MinExchangeVolume)
		maxExchangeDecimal := decimal.NewFromFloat(validLevel.MaxExchangeVolume)
		if exchangeVolumeDecimal.GreaterThanOrEqual(minExchangeDecimal) && exchangeVolumeDecimal.LessThan(maxExchangeDecimal) {
			return validLevel
		}

		if validLevel.ID == lastUserLevel.ID {
			if exchangeVolumeDecimal.GreaterThanOrEqual(maxExchangeDecimal) {
				return validLevel
			}
		}
	}

	return firstUserLevel
}

func (s *levelService) GetLevelsByIds(ids []int64) []Level {
	return s.levelRepo.GetLevelsByIds(ids)
}

func NewUserLevelService(levelRepo LevelRepository) LevelService {
	return &levelService{
		levelRepo: levelRepo,
	}
}
