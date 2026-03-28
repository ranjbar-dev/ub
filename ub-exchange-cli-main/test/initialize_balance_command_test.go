package test

import (
	"context"
	"exchange-go/internal/command"
	"exchange-go/internal/di"
	"exchange-go/internal/userbalance"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type InitializeBalanceCmd struct {
	*suite.Suite
	initializeBalanceCmd command.ConsoleCommand
	db                   *gorm.DB
	userActor            *userActor
}

func (t *InitializeBalanceCmd) SetupSuite() {
	container := getContainer()
	t.initializeBalanceCmd = container.Get(di.InitializeBalanceCommand).(command.ConsoleCommand)
	t.db = getDb()
	t.userActor = getUserActor()
}

func (t *InitializeBalanceCmd) SetupTest() {}

func (t *InitializeBalanceCmd) TearDownTest() {}

func (t *InitializeBalanceCmd) TearDownSuite() {
	t.db.Where("user_id = ?", t.userActor.ID).Delete(userbalance.UserBalance{})
}

func (t *InitializeBalanceCmd) TestRun() {
	ctx := context.Background()

	var flags []string

	flags = append(flags, "-coin=BTC")

	t.initializeBalanceCmd.Run(ctx, flags)

	var userBalances []userbalance.UserBalance

	err := t.db.Where("user_id = ?", t.userActor.ID).Find(&userBalances).Error
	if err != nil {
		t.Fail(err.Error())
	}

	btcUb := userBalances[0]
	assert.Equal(t.T(), "0.0", btcUb.Amount)
	assert.Equal(t.T(), "BTCAddress", btcUb.Address.String)

}

func TestInitializeBalanceCmd(t *testing.T) {
	suite.Run(t, &InitializeBalanceCmd{
		Suite: new(suite.Suite),
	})
}
