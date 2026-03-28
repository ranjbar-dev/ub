package test

import (
	"context"
	"database/sql"
	"exchange-go/internal/command"
	"exchange-go/internal/di"
	"exchange-go/internal/externalexchange"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UpdateOrdersInExternalExchangeCmd struct {
	*suite.Suite
	updateOrdersInExternalExchangeCmd command.ConsoleCommand
	db                                *gorm.DB
	userActor                         *userActor
}

func (t *UpdateOrdersInExternalExchangeCmd) SetupSuite() {
	container := getContainer()
	t.updateOrdersInExternalExchangeCmd = container.Get(di.UpdateOrdersInExternalExchangeCommand).(command.ConsoleCommand)
	t.db = getDb()
	t.userActor = getUserActor()
}

func (t *UpdateOrdersInExternalExchangeCmd) SetupTest() {}

func (t *UpdateOrdersInExternalExchangeCmd) TearDownTest() {}

func (t *UpdateOrdersInExternalExchangeCmd) TearDownSuite() {}

func (t *UpdateOrdersInExternalExchangeCmd) TestRun() {
	lastOrderFromExternal := &externalexchange.OrderFromExternal{
		PairID:          sql.NullInt64{Int64: 1, Valid: true},
		ExternalOrderID: 0,
		ClientOrderID:   "1",
		Type:            "BUY",
		ExchangeType:    "MARKET",
		Status:          sql.NullString{String: "completed", Valid: true},
		Time:            sql.NullTime{Time: time.Now(), Valid: true},
		Timestamp:       sql.NullInt64{Int64: time.Now().Unix() * 1000, Valid: true},
	}
	err := t.db.Create(lastOrderFromExternal).Error
	if err != nil {
		t.Fail(err.Error())
	}

	lastTradeFromExternal := &externalexchange.TradeFromExternal{
		OrderID:         sql.NullInt64{},
		ExternalTradeID: 0,
		Timestamp:       sql.NullInt64{Int64: time.Now().Unix() * 1000, Valid: true},
	}
	err = t.db.Create(lastTradeFromExternal).Error
	if err != nil {
		t.Fail(err.Error())
	}

	var flags []string
	t.updateOrdersInExternalExchangeCmd.Run(context.Background(), flags)

	newOrderFromExternal := &externalexchange.OrderFromExternal{}
	err = t.db.Where("id <> ?", lastOrderFromExternal.ID).First(newOrderFromExternal).Error
	if err != nil {
		t.Fail(err.Error())
	}

	// the value are returned from externalExchangeService:fetchOrders for test env
	assert.Equal(t.T(), "1.00000000", newOrderFromExternal.Amount.String)
	assert.Equal(t.T(), "30000.00000000", newOrderFromExternal.Price.String)
	assert.Equal(t.T(), "COMPLETED", newOrderFromExternal.Status.String)
	assert.Equal(t.T(), "test", newOrderFromExternal.MetaData.String)
	assert.Equal(t.T(), "BUY", newOrderFromExternal.Type)
	assert.Equal(t.T(), "MARKET", newOrderFromExternal.ExchangeType)

	//newTradeFromExternal := &externalexchange.TradeFromExternal{}
	//err = t.db.Where("id <> ?", lastTradeFromExternal.ID).First(newTradeFromExternal).Error
	//if err != nil {
	//	t.Fail(err.Error())
	//}
	//// the value are returned from externalExchangeService:fetchTrades for test env
	//assert.Equal(t.T(), "1.00000000", newTradeFromExternal.Amount.String)
	//assert.Equal(t.T(), "30000.00000000", newTradeFromExternal.Price.String)
	//assert.Equal(t.T(), "test", newTradeFromExternal.MetaData.String)
	//assert.Equal(t.T(), "BTC", newTradeFromExternal.Coin.String)
	//assert.Equal(t.T(), "0.00010000", newTradeFromExternal.Commission.String)

}

func TestUpdateOrdersInExternalExchangeCmd(t *testing.T) {
	suite.Run(t, &UpdateOrdersInExternalExchangeCmd{
		Suite: new(suite.Suite),
	})
}
