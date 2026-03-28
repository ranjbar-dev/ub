package test

import (
	"context"
	"database/sql"
	"exchange-go/internal/command"
	"exchange-go/internal/di"
	"exchange-go/internal/payment"
	"exchange-go/internal/transaction"
	"exchange-go/internal/userbalance"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type CheckWithdrawalsInExternalExchangeCmd struct {
	*suite.Suite
	checkWithdrawalsInExternalExchangeCmd command.ConsoleCommand
	db                                    *gorm.DB
	userActor                             *userActor
}

func (t *CheckWithdrawalsInExternalExchangeCmd) SetupSuite() {
	container := getContainer()
	t.checkWithdrawalsInExternalExchangeCmd = container.Get(di.CheckWithdrawalsInExternalExchangeCommand).(command.ConsoleCommand)
	t.db = getDb()
	t.userActor = getUserActor()
}

func (t *CheckWithdrawalsInExternalExchangeCmd) SetupTest() {
	t.db.Where("id > ?", int64(0)).Delete(userbalance.UserBalance{})
	t.db.Where("id > ?", int64(0)).Delete(transaction.Transaction{})
	t.db.Where("id > ?", int64(0)).Delete(payment.ExtraInfo{})
	t.db.Where("id > ?", int64(0)).Delete(payment.Payment{})
}

func (t *CheckWithdrawalsInExternalExchangeCmd) TearDownTest() {}

func (t *CheckWithdrawalsInExternalExchangeCmd) TearDownSuite() {
	t.db.Where("id > ?", int64(0)).Delete(userbalance.UserBalance{})
	t.db.Where("id > ?", int64(0)).Delete(transaction.Transaction{})
	t.db.Where("id > ?", int64(0)).Delete(payment.ExtraInfo{})
	t.db.Where("id > ?", int64(0)).Delete(payment.Payment{})
}

func (t *CheckWithdrawalsInExternalExchangeCmd) TestRun_SuccessfulWithdrawal() {
	//insert data in database
	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "2.00000000",
		FrozenAmount:  "1.00000000",
	}
	err := t.db.Create(btcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	p := &payment.Payment{
		UserID:    t.userActor.ID,
		CoinID:    2, //BTC look at currencyseed.go file
		Type:      payment.TypeWithdraw,
		Status:    payment.StatusInProgress,
		Code:      "BTC",
		Amount:    sql.NullString{String: "1.00000000", Valid: true},
		FeeAmount: sql.NullString{String: "0.00100000", Valid: true},
	}

	err = t.db.Create(p).Error
	if err != nil {
		t.Fail(err.Error())
	}

	externalExchangeWithdrawID := "1"
	pei := &payment.ExtraInfo{
		ID:                         0,
		PaymentID:                  p.ID,
		ExternalExchangeWithdrawID: sql.NullString{String: externalExchangeWithdrawID, Valid: true},
	}

	err = t.db.Create(pei).Error
	if err != nil {
		t.Fail(err.Error())
	}

	var flags []string
	flags = append(flags, payment.StatusCompleted)

	t.checkWithdrawalsInExternalExchangeCmd.Run(context.Background(), flags)

	updatedPayment := &payment.Payment{}
	err = t.db.Where(payment.Payment{ID: p.ID}).First(updatedPayment).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), payment.StatusCompleted, updatedPayment.Status)
	assert.Equal(t.T(), "txId", updatedPayment.TxID.String)

	updatedExtraInfo := &payment.ExtraInfo{}
	err = t.db.Where(payment.ExtraInfo{ID: pei.ID}).First(updatedExtraInfo).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "test", updatedExtraInfo.ExternalExchangeWithdrawInfo.String)

	updatedBtcUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: btcUb.ID}).First(updatedBtcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "1.00000000", updatedBtcUb.Amount)
	assert.Equal(t.T(), "0.00000000", updatedBtcUb.FrozenAmount)

	//check transaction table
	var transactions []transaction.Transaction
	err = t.db.Where("user_id = ?", t.userActor.ID).Find(&transactions).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), 2, len(transactions))

	for _, tx := range transactions {
		switch tx.Type {
		case transaction.TypeWithdraw:
			assert.Equal(t.T(), "1.00000000", tx.Amount.String)
			assert.Equal(t.T(), p.ID, tx.PaymentID.Int64)
		case transaction.TypeWithdrawFee:
			assert.Equal(t.T(), "0.00100000", tx.Amount.String)
			assert.Equal(t.T(), p.ID, tx.PaymentID.Int64)
		default:
			t.Fail("we should not be in default case")
		}
	}
}

func (t *CheckWithdrawalsInExternalExchangeCmd) TestRun_FailedWithdrawal() {
	//insert data in database
	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "2.00000000",
		FrozenAmount:  "1.00000000",
	}
	err := t.db.Create(btcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	p := &payment.Payment{
		UserID:    t.userActor.ID,
		CoinID:    2, //BTC look at currencyseed.go file
		Type:      payment.TypeWithdraw,
		Status:    payment.StatusInProgress,
		Code:      "BTC",
		Amount:    sql.NullString{String: "1.00000000", Valid: true},
		FeeAmount: sql.NullString{String: "0.00100000", Valid: true},
	}

	err = t.db.Create(p).Error
	if err != nil {
		t.Fail(err.Error())
	}

	externalExchangeWithdrawID := "1"
	pei := &payment.ExtraInfo{
		ID:                         0,
		PaymentID:                  p.ID,
		ExternalExchangeWithdrawID: sql.NullString{String: externalExchangeWithdrawID, Valid: true},
	}

	err = t.db.Create(pei).Error
	if err != nil {
		t.Fail(err.Error())
	}

	var flags []string
	flags = append(flags, payment.StatusFailed)

	t.checkWithdrawalsInExternalExchangeCmd.Run(context.Background(), flags)

	updatedPayment := &payment.Payment{}
	err = t.db.Where(payment.Payment{ID: p.ID}).First(updatedPayment).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), payment.StatusFailed, updatedPayment.Status)
	assert.Equal(t.T(), "", updatedPayment.TxID.String)

	updatedExtraInfo := &payment.ExtraInfo{}
	err = t.db.Where(payment.ExtraInfo{ID: pei.ID}).First(updatedExtraInfo).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "test", updatedExtraInfo.ExternalExchangeWithdrawInfo.String)

	updatedBtcUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: btcUb.ID}).First(updatedBtcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "2.00000000", updatedBtcUb.Amount)
	assert.Equal(t.T(), "0.00000000", updatedBtcUb.FrozenAmount)

	//check transaction table
	var transactions []transaction.Transaction
	err = t.db.Where("user_id = ?", t.userActor.ID).Find(&transactions).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), 0, len(transactions))
}

type CryptoBalance struct {
	ID                 int64
	CurrencyID         sql.NullInt64
	ExternalExchangeID sql.NullInt64
	Type               string
	Address            sql.NullString
	Tag                sql.NullString
	MetaData           sql.NullString `gorm:"column:metaData"`
	FreeAmount         sql.NullString
	LockedAmount       sql.NullString
	CreatedAt          time.Time
	UpdatedAt          time.Time
	BlockchainNetwork  sql.NullString
}

func (CryptoBalance) TableName() string {
	return "crypto_balance"
}

func (t *CheckWithdrawalsInExternalExchangeCmd) TestRun_CheckInternalTransfers() {
	//insert data in database

	cr := &CryptoBalance{
		Type: "EXTERNAL",
	}

	err := t.db.Create(cr).Error
	if err != nil {
		t.Fail(err.Error())
	}

	//err := t.db.Raw(`INSERT INTO crypto_balance(id,currency_id,type,created_at,updated_at) VALUES(1,2,"EXTERNAL","2021-01-21 20:20:20","2021-01-21 20:20:20")`).Error
	//fmt.Println("err====>", err)
	//if err != nil {
	//	t.Fail(err.Error())
	//}
	metadata := `{"info":{"success":true,"id":"1"},"id":"1"}`

	it := &payment.InternalTransfer{
		FromBalanceID: cr.ID,
		Amount:        "2.00000000",
		Status:        payment.InternalTransferStatusInProgress,
		Network:       "",
		Metadata:      sql.NullString{String: metadata, Valid: true},
	}
	err = t.db.Create(it).Error
	if err != nil {
		t.Fail(err.Error())
	}

	var flags []string
	flags = append(flags, payment.StatusCompleted)

	t.checkWithdrawalsInExternalExchangeCmd.Run(context.Background(), flags)

	updatedInternalTransfer := &payment.InternalTransfer{}
	err = t.db.Where(payment.InternalTransfer{ID: it.ID}).First(updatedInternalTransfer).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), payment.StatusCompleted, updatedInternalTransfer.Status)
	assert.Equal(t.T(), "txId", updatedInternalTransfer.TxID.String)
}

func TestCheckWithdrawalsInExternalExchangeCmd(t *testing.T) {
	suite.Run(t, &CheckWithdrawalsInExternalExchangeCmd{
		Suite: new(suite.Suite),
	})
}
