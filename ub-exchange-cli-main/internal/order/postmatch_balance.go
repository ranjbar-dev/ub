package order

import (
	"database/sql"

	"exchange-go/internal/currency"
	"exchange-go/internal/transaction"
	"exchange-go/internal/userbalance"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (ps *postOrderMatchingService) updateUserBalances(tx *gorm.DB, userBalances [2]*userbalance.UserBalance,
	demandedDecimal decimal.Decimal, payedByDecimal decimal.Decimal, frozenReductionDecimal decimal.Decimal,
	orderType string, pair currency.Pair) error {

	var payedByUb *userbalance.UserBalance
	var demandedUb *userbalance.UserBalance

	if orderType == TypeBuy {
		if userBalances[0].CoinID == pair.DependentCoinID {
			payedByUb = userBalances[1]
			demandedUb = userBalances[0]
		} else {
			payedByUb = userBalances[0]
			demandedUb = userBalances[1]
		}
	}

	if orderType == TypeSell {
		if userBalances[0].CoinID == pair.DependentCoinID {
			payedByUb = userBalances[0]
			demandedUb = userBalances[1]
		} else {
			payedByUb = userBalances[1]
			demandedUb = userBalances[0]

		}
	}

	formerDemandedDecimal, _ := decimal.NewFromString(demandedUb.Amount)

	finalDemanded := formerDemandedDecimal.Add(demandedDecimal).StringFixed(8)
	newDemandedUb := &userbalance.UserBalance{
		ID:     demandedUb.ID,
		Amount: finalDemanded,
	}
	err := tx.Model(newDemandedUb).Updates(newDemandedUb).Error
	if err != nil {
		return err
	}
	//in case the user of the buy and sell order be the same person
	//we handle the balance using the pointer, here we update this
	//pointer so the other balance has the updated data
	demandedUb.Amount = finalDemanded

	formerPayedByDecimal, _ := decimal.NewFromString(payedByUb.Amount)
	formerFrozenPayedByDecimal, _ := decimal.NewFromString(payedByUb.FrozenAmount)
	finalPayedBy := formerPayedByDecimal.Sub(payedByDecimal).StringFixed(8)
	finalFrozenPayedBy := formerFrozenPayedByDecimal.Sub(frozenReductionDecimal).StringFixed(8)
	newPayedByUb := &userbalance.UserBalance{
		ID:           payedByUb.ID,
		Amount:       finalPayedBy,
		FrozenAmount: finalFrozenPayedBy,
	}
	err = tx.Model(newPayedByUb).Updates(newPayedByUb).Error
	//in case the user of the buy and sell order be the same person
	//we handle the balance using the pointer, here we update this
	//pointer so the other balance has the updated data
	payedByUb.Amount = finalPayedBy
	payedByUb.FrozenAmount = finalFrozenPayedBy
	return err

}

func (ps *postOrderMatchingService) createTransactions(tx *gorm.DB, orders []Order, pair currency.Pair) error {
	for _, order := range orders {
		orderType := order.Type
		userID := order.UserID
		//create transaction for order demanded
		demandedCoinID := pair.BasisCoin.ID
		demandedCoinName := pair.BasisCoin.Code
		demandedAmount := order.FinalDemandedAmount.String
		if orderType == TypeBuy {
			demandedCoinID = pair.DependentCoin.ID
			demandedCoinName = pair.DependentCoin.Code
		}
		demandedTransaction := &transaction.Transaction{
			UserID:    userID,
			CoinID:    demandedCoinID,
			OrderID:   sql.NullInt64{Int64: order.ID, Valid: true},
			Type:      transaction.TypeDemanded,
			Amount:    sql.NullString{String: demandedAmount, Valid: true},
			CoinName:  demandedCoinName,
			PaymentID: sql.NullInt64{Int64: 0, Valid: false},
		}
		err := tx.Omit(clause.Associations).Create(demandedTransaction).Error
		if err != nil {
			return err
		}

		//create transaction for order payedBy
		payedByCoinID := pair.DependentCoin.ID
		payedByCoinName := pair.DependentCoin.Code
		payedByAmount := order.FinalPayedByAmount.String
		if orderType == TypeBuy {
			payedByCoinID = pair.BasisCoin.ID
			payedByCoinName = pair.BasisCoin.Code
		}
		payedByTransaction := &transaction.Transaction{
			UserID:    userID,
			CoinID:    payedByCoinID,
			OrderID:   sql.NullInt64{Int64: order.ID, Valid: true},
			Type:      transaction.TypePayedBy,
			Amount:    sql.NullString{String: payedByAmount, Valid: true},
			CoinName:  payedByCoinName,
			PaymentID: sql.NullInt64{Int64: 0, Valid: false},
		}
		err = tx.Omit(clause.Associations).Create(payedByTransaction).Error
		if err != nil {
			return err
		}

		//create transaction for order fee
		feeCoinID := pair.BasisCoin.ID
		feeCoinName := pair.BasisCoin.Code
		finalDemandedDecimal, err := decimal.NewFromString(order.FinalDemandedAmount.String)
		if err != nil {
			return err
		}
		feePercentageDecimal := decimal.NewFromFloat(order.FeePercentage.Float64)
		fee := finalDemandedDecimal.Mul(feePercentageDecimal).String()
		if orderType == TypeBuy {
			feeCoinID = pair.DependentCoin.ID
			feeCoinName = pair.DependentCoin.Code
		}

		feeType := transaction.TypeTakerFee
		if order.IsMaker.Valid && order.IsMaker.Bool {
			feeType = transaction.TypeMakerFee
		}
		feeTransaction := &transaction.Transaction{
			UserID:    userID,
			CoinID:    feeCoinID,
			OrderID:   sql.NullInt64{Int64: order.ID, Valid: true},
			Type:      feeType,
			Amount:    sql.NullString{String: fee, Valid: true},
			CoinName:  feeCoinName,
			PaymentID: sql.NullInt64{Int64: 0, Valid: false},
		}
		err = tx.Omit(clause.Associations).Create(feeTransaction).Error
		if err != nil {
			return err
		}

	}

	//for taker
	return nil

}
