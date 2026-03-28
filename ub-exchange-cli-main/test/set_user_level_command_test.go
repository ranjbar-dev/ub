package test

import (
	"context"
	"database/sql"
	"exchange-go/internal/command"
	"exchange-go/internal/di"
	"exchange-go/internal/order"
	"exchange-go/internal/user"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type SetUserLevelCmd struct {
	*suite.Suite
	userLevelCmd command.ConsoleCommand
	db           *gorm.DB
	userActor    *userActor
}

func (t *SetUserLevelCmd) SetupSuite() {
	container := getContainer()
	t.userLevelCmd = container.Get(di.SetUserLevelCommand).(command.ConsoleCommand)
	t.db = getDb()
	t.userActor = getUserActor()
}

func (t *SetUserLevelCmd) SetupTest() {}

func (t *SetUserLevelCmd) TearDownTest() {
}

func (t *SetUserLevelCmd) TearDownSuite() {
}

func (t *SetUserLevelCmd) TestRun() {
	//insert order and trade to db for user actor
	o1 := &order.Order{
		ID:             1,
		UserID:         t.userActor.ID,
		Type:           order.TypeBuy,
		ExchangeType:   order.ExchangeTypeLimit,
		Price:          sql.NullString{String: "50000.00000000", Valid: true},
		Status:         order.StatusFilled,
		DemandedAmount: sql.NullString{String: "0.10000000", Valid: true},
		DemandedCoin:   "BTC",
		PayedByAmount:  sql.NullString{String: "5000.00000000", Valid: true},
		PayedByCoin:    "USDT",
		PairID:         1,
	}

	o2 := &order.Order{
		ID:             2,
		UserID:         t.userActor.ID,
		Type:           order.TypeSell,
		ExchangeType:   order.ExchangeTypeLimit,
		Price:          sql.NullString{String: "2000.00000000", Valid: true},
		Status:         order.StatusFilled,
		DemandedAmount: sql.NullString{String: "5000.00000000", Valid: true},
		DemandedCoin:   "USDT",
		PayedByAmount:  sql.NullString{String: "2.50000000", Valid: true},
		PayedByCoin:    "BTC",
		PairID:         2,
	}

	orders := []*order.Order{o1, o2}
	err := t.db.Create(orders).Error
	if err != nil {
		t.Fail(err.Error())
	}

	now := time.Now()
	lastDay := now.Add(-1 * 24 * time.Hour)
	t1 := &order.Trade{
		Price:      sql.NullString{String: "50000", Valid: true},
		Amount:     sql.NullString{String: "0.1", Valid: true},
		PairID:     1, //BTC-USDT
		BuyOrderID: sql.NullInt64{Int64: o1.ID, Valid: true},
		CreatedAt:  lastDay,
	}

	t2 := &order.Trade{
		Price:       sql.NullString{String: "2000", Valid: true},
		Amount:      sql.NullString{String: "2.5", Valid: true},
		PairID:      2, //ETH-USDT
		SellOrderID: sql.NullInt64{Int64: o2.ID, Valid: true},
		CreatedAt:   lastDay,
	}

	trades := []*order.Trade{t1, t2}

	err = t.db.Create(trades).Error
	if err != nil {
		t.Fail(err.Error())
	}

	//run Command
	ctx := context.Background()
	var flags []string
	t.userLevelCmd.Run(ctx, flags)

	//check the exchange number and value in database

	updatedUser := &user.User{}
	err = t.db.Where("id = ?", t.userActor.ID).First(updatedUser).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), int64(2), updatedUser.ExchangeNumber)
	assert.Equal(t.T(), "0.20000000", updatedUser.ExchangeVolumeAmount)

}

func TestSetUserLevelCmd(t *testing.T) {
	suite.Run(t, &SetUserLevelCmd{
		Suite: new(suite.Suite),
	})
}
