package order

import (
	"context"
	"database/sql"
	"encoding/json"
	"exchange-go/internal/currency"
	"exchange-go/internal/platform"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"fmt"
	"strconv"
	"sync"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateManager validates and creates new orders, including balance checks,
// fee calculations, and database persistence.
type CreateManager interface {
	// CreateOrder validates the provided order data, verifies the user's balance and level,
	// then persists the order and freezes the required funds.
	CreateOrder(data CreateRequiredData) (order *Order, err error)
}

type createManager struct {
	db                 *gorm.DB
	userBalanceService userbalance.Service
	userLevelService   user.LevelService
	priceGenerator     currency.PriceGenerator
	mutex              *sync.Mutex
	data               CreateRequiredData
}

type CreateRequiredData struct {
	User           *user.User
	Pair           *currency.Pair
	Amount         string
	OrderType      string
	ExchangeType   string
	Price          string
	UserAgentInfo  UserAgentInfo
	StopPointPrice string
	CurrentPrice   string
	IsInstant      bool
	IsFastExchange bool
	IsAutoExchange bool
	//parentOrder    *Order

	//these are would be set later
	isOrderMarket bool
	isStopOrder   bool
	payedBy       decimal.Decimal
	demanded      decimal.Decimal
}

func (cm *createManager) CreateOrder(data CreateRequiredData) (order *Order, err error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.data = data
	err = cm.validateRequiredData()
	if err != nil {
		return nil, err
	}

	err = cm.setExtraRequiredData()
	if err != nil {
		return nil, err
	}

	isGreater, err := cm.isGreaterThanOrEqualMinimumOrderAmount()
	if err != nil {
		return nil, err
	}
	if !isGreater {
		return nil, newOrderCreateValidationError("the minimum order amount must be more than " + cm.data.Pair.MinimumOrderAmount.String + " " + cm.data.Pair.BasisCoin.Code)
	}

	allows, err := cm.doesUserLevelAllows()
	if err != nil {
		return nil, err
	}
	if !allows {
		return nil, newOrderCreateValidationError("your user level is low to place this order. please verify your identity to boost up your level")
	}

	o, err := cm.saveInDb()
	if err != nil {
		return nil, err
	}

	return o, nil

}

func (cm *createManager) validateRequiredData() error {
	if cm.data.User == nil {
		return newOrderCreateValidationError("user not found")
	}

	if cm.data.Pair == nil {
		return newOrderCreateValidationError("pair currency not found")
	}

	if cm.data.Amount == "" {
		return newOrderCreateValidationError("amount is not valid")
	}
	amountDecimal, err := decimal.NewFromString(cm.data.Amount)
	if err != nil || !amountDecimal.IsPositive() {
		return newOrderCreateValidationError("amount is not valid")
	}

	orderType := cm.data.OrderType
	if orderType != TypeBuy && orderType != TypeSell {
		return newOrderCreateValidationError("type is not valid")
	}

	exchangeType := cm.data.ExchangeType
	if exchangeType != ExchangeTypeLimit && exchangeType != ExchangeTypeMarket {
		return newOrderCreateValidationError("exchange type is not valid")
	}

	if exchangeType == ExchangeTypeLimit {
		if cm.data.Price == "" {
			return newOrderCreateValidationError("price is not valid")
		}
		priceDecimal, err := decimal.NewFromString(cm.data.Price)
		if err != nil || !priceDecimal.IsPositive() {
			return newOrderCreateValidationError("price is not valid")
		}
	}

	if cm.data.StopPointPrice != "" {
		if exchangeType != ExchangeTypeLimit {
			return newOrderCreateValidationError("stop order must be limit")
		}
		stopPointPriceDecimal, err := decimal.NewFromString(cm.data.StopPointPrice)
		if err != nil || !stopPointPriceDecimal.IsPositive() {
			return newOrderCreateValidationError("stop point price is not valid")
		}
	}

	return nil

}

func (cm *createManager) setExtraRequiredData() error {
	//for market orders we do not have tradePrice so we use currentPrice
	amount := cm.data.Amount

	price := cm.data.Price
	orderType := cm.data.OrderType

	if cm.data.Price == "" {
		cm.data.isOrderMarket = true
	}

	if cm.data.StopPointPrice != "" {
		cm.data.isStopOrder = true
	}

	isOrderMarket := cm.data.isOrderMarket
	isInstant := cm.data.IsInstant
	if isOrderMarket {
		price = cm.data.CurrentPrice
	}

	amountDecimal, err := decimal.NewFromString(amount)
	if err != nil {
		return err
	}

	priceDecimal, err := decimal.NewFromString(price)

	if err != nil {
		return err
	}

	if isOrderMarket {
		if cm.data.OrderType == TypeBuy {
			if isInstant {
				cm.data.demanded = amountDecimal.Div(priceDecimal)
				cm.data.payedBy = amountDecimal
			} else {
				cm.data.demanded = amountDecimal
				cm.data.payedBy = amountDecimal.Mul(priceDecimal)
			}
		} else {
			cm.data.demanded = amountDecimal.Mul(priceDecimal)
			cm.data.payedBy = amountDecimal
		}
	} else {
		if orderType == TypeBuy {
			cm.data.demanded = amountDecimal
			cm.data.payedBy = amountDecimal.Mul(priceDecimal)
		} else {
			cm.data.payedBy = amountDecimal
			cm.data.demanded = amountDecimal.Mul(priceDecimal)
		}
	}
	return nil

}

func (cm *createManager) isGreaterThanOrEqualMinimumOrderAmount() (bool, error) {
	minimumOrderAmount := cm.data.Pair.MinimumOrderAmount.String
	minimumOrderDecimal, err := decimal.NewFromString(minimumOrderAmount)
	if err != nil {
		return false, err
	}

	orderType := cm.data.OrderType
	if orderType == TypeBuy {
		return minimumOrderDecimal.LessThanOrEqual(cm.data.payedBy), nil
	}
	return minimumOrderDecimal.LessThanOrEqual(cm.data.demanded), nil
}

func (cm *createManager) doesUserLevelAllows() (bool, error) {
	ctx := context.Background()
	u := cm.data.User
	userLevel, _ := cm.userLevelService.GetLevelByID(u.UserLevelID)
	demandedDecimal := cm.data.demanded
	demandedCoin := cm.data.Pair.DependentCoin
	if cm.data.OrderType == TypeSell {
		demandedCoin = cm.data.Pair.BasisCoin
	}
	if u.Kyc == user.KycLevelMinimum {
		exchangeVolumeDecimal, _ := decimal.NewFromString(u.ExchangeVolumeAmount)                   //in btc
		exchangeVolumeLimitDecimal, _ := decimal.NewFromString(userLevel.ExchangeVolumeLimitAmount) //in btc
		orderAmountBasedOnBtc, _ := cm.priceGenerator.GetAmountBasedOnBTC(ctx, demandedCoin.Code, demandedDecimal.StringFixed(8))
		orderAmountBasedOnBtcDecimal, _ := decimal.NewFromString(orderAmountBasedOnBtc)
		newExchangeVolumeDecimal := exchangeVolumeDecimal.Add(orderAmountBasedOnBtcDecimal)
		exchangeNumberLimit := userLevel.ExchangeNumberLimit
		if newExchangeVolumeDecimal.GreaterThan(exchangeVolumeLimitDecimal) || u.ExchangeNumber+1 > exchangeNumberLimit {
			return false, nil
		}
	}
	return true, nil
}

func (cm *createManager) saveInDb() (*Order, error) {
	tx := cm.db.Begin()
	err := tx.Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	coinID := cm.data.Pair.BasisCoin.ID
	if cm.data.OrderType == TypeSell {
		coinID = cm.data.Pair.DependentCoin.ID
	}
	ba := &userbalance.UserBalance{}
	err = cm.userBalanceService.GetBalanceOfUserByCoinUsingTx(tx, cm.data.User.ID, coinID, ba)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	baAmountDecimal, err := decimal.NewFromString(ba.Amount)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	baFrozenAmountDecimal, err := decimal.NewFromString(ba.FrozenAmount)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	realBalance := baAmountDecimal.Sub(baFrozenAmountDecimal)
	if cm.data.payedBy.GreaterThan(realBalance) {
		tx.Rollback()
		return nil, newOrderCreateValidationError("user balance is not enough")
	}

	uai, err := json.Marshal(cm.data.UserAgentInfo)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	extraInfo := &ExtraInfo{
		UserAgentInfo: sql.NullString{String: string(uai), Valid: true},
	}
	if cm.data.IsAutoExchange {
		extraInfo.AutoExchange = sql.NullBool{Bool: true, Valid: true}
	}
	err = tx.Omit(clause.Associations).Create(&extraInfo).Error //create orderExtraInfo
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	price := sql.NullString{String: "", Valid: false}
	if cm.data.Price != "" {
		price.Valid = true
		priceDecimal, _ := decimal.NewFromString(cm.data.Price)
		price.String = priceDecimal.StringFixed(8)
	}

	currentMarketPrice := sql.NullString{String: cm.data.CurrentPrice, Valid: true}
	stopPointPrice := sql.NullString{String: "", Valid: false}
	if cm.data.StopPointPrice != "" {
		stopPointPrice.Valid = true
		stopPointPriceDecimal, _ := decimal.NewFromString(cm.data.StopPointPrice)
		stopPointPrice.String = stopPointPriceDecimal.StringFixed(8)

		//we do not set current market price for stop orders
		currentMarketPrice.Valid = false
		currentMarketPrice.String = ""
	}

	demanded := cm.data.demanded.StringFixed(8)
	payedBy := cm.data.payedBy.StringFixed(8)
	demandedMoneyCoin := cm.data.Pair.DependentCoin.Code
	PayedByMoneyCoin := cm.data.Pair.BasisCoin.Code

	if cm.data.OrderType == TypeSell {
		demandedMoneyCoin = cm.data.Pair.BasisCoin.Code
		PayedByMoneyCoin = cm.data.Pair.DependentCoin.Code
	}

	o := &Order{
		UserID:             cm.data.User.ID,
		Type:               cm.data.OrderType,
		ExchangeType:       cm.data.ExchangeType,
		Price:              price,
		Status:             StatusOpen,
		DemandedAmount:     sql.NullString{String: demanded, Valid: true},
		DemandedCoin:       demandedMoneyCoin,
		PayedByAmount:      sql.NullString{String: payedBy, Valid: true},
		PayedByCoin:        PayedByMoneyCoin,
		PairID:             cm.data.Pair.ID,
		ExtraInfoID:        sql.NullInt64{Int64: extraInfo.ID, Valid: true},
		Level:              sql.NullInt64{Int64: 1, Valid: true},
		StopPointPrice:     stopPointPrice,
		CurrentMarketPrice: currentMarketPrice,
		IsFastExchange:     cm.data.IsFastExchange,
	}

	err = tx.Omit(clause.Associations).Create(o).Error //create order
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//set data for tree model
	path := strconv.FormatInt(o.ID, 10) + ","
	o.Path = sql.NullString{String: path, Valid: true}
	err = tx.Omit(clause.Associations).Save(o).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//freeze balance
	frozenDecimal, err := decimal.NewFromString(ba.FrozenAmount)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	newFrozenAmountDecimal := frozenDecimal.Add(cm.data.payedBy)
	// Guard: frozen amount must never exceed total balance
	if newFrozenAmountDecimal.GreaterThan(baAmountDecimal) {
		tx.Rollback()
		return nil, newOrderCreateValidationError("insufficient balance for freeze")
	}
	newFrozenAmount := newFrozenAmountDecimal.StringFixed(8)
	ba.FrozenAmount = newFrozenAmount
	err = tx.Omit(clause.Associations).Save(ba).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return o, nil
}

func NewOrderCreateManager(db *gorm.DB, userBalanceService userbalance.Service, userLevelService user.LevelService, pg currency.PriceGenerator) CreateManager {
	return &createManager{
		db:                 db,
		userBalanceService: userBalanceService,
		userLevelService:   userLevelService,
		priceGenerator:     pg,
		mutex:              &sync.Mutex{},
	}
}

func newOrderCreateValidationError(message string) platform.OrderCreateValidationError {
	return platform.OrderCreateValidationError{Err: fmt.Errorf(message)}
}
