package currency

import (
	"database/sql"
	"encoding/json"
	"time"
)

const (
	MaxOfMarketPricePercentage  = 0.001
	TraderBotRuleTypePercentage = "PERCENTAGE"
	TraderBotRuleTypeConst      = "CONST"
)

type OtherBlockchainNetworksConfig struct {
	Code                 string `json:"code"`
	SupportsWithdraw     bool   `json:"supportsWithdraw"`
	SupportsDeposit      bool   `json:"supportsDeposit"`
	CompletedNetworkName string `json:"completedNetworkName"`
	Fee                  string `json:"fee"`
}

type Coin struct {
	ID                             int64
	Name                           string
	Code                           string
	SubUnit                        int
	ShowSubUnit                    int
	IsActive                       bool
	IsMain                         bool
	Priority                       int32
	ConversionRatio                float64
	Image                          string         `gorm:"column:image_path"`
	SecondImage                    sql.NullString `gorm:"column:second_image_path"`
	MinimumWithdraw                string
	MaximumWithdraw                string
	CompletedNetworkName           sql.NullString
	BlockchainNetwork              sql.NullString
	WithdrawalFee                  sql.NullFloat64
	SupportsWithdraw               sql.NullBool `gorm:"column:supports_withdraw"`
	SupportsDeposit                sql.NullBool `gorm:"column:supports_deposit"`
	OtherBlockchainNetworksConfigs sql.NullString
	DepositComments                sql.NullString
	WithdrawComments               sql.NullString
}

func (Coin) TableName() string {
	return "currencies"
}

func (c Coin) GetOtherBlockchainNetworksConfigs() ([]OtherBlockchainNetworksConfig, error) {
	var networkConfigs []OtherBlockchainNetworksConfig
	if !c.OtherBlockchainNetworksConfigs.Valid {
		return networkConfigs, nil
	}
	err := json.Unmarshal([]byte(c.OtherBlockchainNetworksConfigs.String), &networkConfigs)
	return networkConfigs, err

}

// Repository provides data access for active cryptocurrency coins.
type Repository interface {
	// GetActiveCoins returns all coins that are currently active on the exchange.
	GetActiveCoins() []Coin
	// GetCoinByCode looks up a single coin by its ticker symbol (e.g. "BTC") and
	// populates the provided Coin pointer.
	GetCoinByCode(code string, coin *Coin) error
	// GetCoinsAlphabetically returns all coins sorted alphabetically by name.
	GetCoinsAlphabetically() []Coin
}

//in db {"buyValue":0.0002,"sellValue":0.0002,"type":"PERCENTAGE"}
type BotRules struct {
	Type      string  `json:"type"`
	BuyValue  float64 `json:"buyValue"`
	SellValue float64 `json:"sellValue"`
}

type Pair struct {
	ID                       int64
	Name                     string
	IsActive                 bool
	IsMain                   bool
	Spread                   float32 `gorm:"column:ohlc_spread"`
	ShowDigits               int
	BasisCoinID              int64 `gorm:"column:basis_currency_id"`
	BasisCoin                Coin  `gorm:"foreignKey:BasisCoinID"`
	DependentCoinID          int64 `gorm:"column:dependent_currency_id"`
	DependentCoin            Coin  `gorm:"foreignKey:DependentCoinID"`
	MakerFee                 float64
	TakerFee                 float64
	TradeStatus              string `gorm:"default:'FULL_TRADE'"`
	AggregationStatus        string `gorm:"default:'RUN'"`
	BotRules                 sql.NullString
	MinimumOrderAmount       sql.NullString
	MaxOurExchangeLimit      string
	BotOrdersAggregationTime sql.NullInt64
}

func (Pair) TableName() string {
	return "pair_currencies"
}

func (p Pair) GetBotRules() BotRules {
	rules := p.BotRules

	if rules.Valid && rules.String != "" {
		botRules := BotRules{}
		err := json.Unmarshal([]byte(rules.String), &botRules)
		if err == nil {
			return botRules
		}
	}

	return BotRules{
		Type:      TraderBotRuleTypePercentage,
		BuyValue:  MaxOfMarketPricePercentage,
		SellValue: MaxOfMarketPricePercentage,
	}
}

// PairRepository provides data access for trading pairs (e.g. BTC-USDT).
type PairRepository interface {
	// GetActivePairCurrenciesList returns all trading pairs that are currently active.
	GetActivePairCurrenciesList() []Pair
	// GetPairByID looks up a trading pair by its database ID.
	GetPairByID(id int64, pair *Pair) error
	// GetPairByName looks up a trading pair by its name (e.g. "BTC-USDT").
	GetPairByName(name string, pair *Pair) error
	// GetAllPairs returns every trading pair regardless of active status.
	GetAllPairs() []Pair
	// GetPairsByName returns all trading pairs whose names match the provided list.
	GetPairsByName(names []string) []Pair
}

type KlineSync struct {
	ID         int64
	PairID     int64 `gorm:"column:pair_currency_id"`
	TimeFrame  sql.NullString
	StartTime  sql.NullTime
	EndTime    sql.NullTime
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Type       string
	WithUpdate bool
}

func (KlineSync) TableName() string {
	return "ohlc_sync"
}

// KlineSyncRepository manages OHLC candlestick synchronization task records.
type KlineSyncRepository interface {
	// Create inserts a new kline sync task record.
	Create(klineSync *KlineSync) error
	// Update persists changes to an existing kline sync task record.
	Update(klineSync *KlineSync) error
	// GetKlineSyncsByStatusAndLimit returns up to limit sync tasks that match the given status.
	GetKlineSyncsByStatusAndLimit(status string, limit int) []KlineSync
}

type FavoritePair struct {
	UserID int
	PairID int64 `gorm:"column:pair_currency_id"`
}

func (FavoritePair) TableName() string {
	return "user_favorite_pair_currency"
}

// FavoritePairRepository manages user-specific favorite trading pair bookmarks.
type FavoritePairRepository interface {
	// Create adds a trading pair to the user's favorites.
	Create(favoritePair *FavoritePair) error
	// Delete removes a trading pair from the user's favorites.
	Delete(favoritePair *FavoritePair) error
	// GetFavoritePair looks up a single favorite pair entry for a given user and pair ID.
	GetFavoritePair(userID int, pairID int64, favoritePair *FavoritePair) error
	// GetUserFavoritePairs returns all favorite pairs for the specified user.
	GetUserFavoritePairs(userID int) []FavoritePairQueryFields
}
