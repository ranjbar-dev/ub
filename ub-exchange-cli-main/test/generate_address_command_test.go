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

type GenerateAddressCmd struct {
	*suite.Suite
	generateAddressCmd command.ConsoleCommand
	db                 *gorm.DB
	userActor          *userActor
}

func (t *GenerateAddressCmd) SetupSuite() {
	container := getContainer()
	t.generateAddressCmd = container.Get(di.GenerateAddressCommand).(command.ConsoleCommand)
	t.db = getDb()
	t.userActor = getUserActor()
}

func (t *GenerateAddressCmd) SetupTest() {}

func (t *GenerateAddressCmd) TearDownTest() {}

func (t *GenerateAddressCmd) TearDownSuite() {
	t.db.Where("user_id = ?", t.userActor.ID).Delete(userbalance.UserBalance{})
}

func (t *GenerateAddressCmd) TestRun() {
	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "10000.00",
		FrozenAmount:  "1000.00",
	}

	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.1",
		FrozenAmount:  "0",
	}
	err := t.db.Create(usdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	err = t.db.Create(btcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	//run Command
	ctx := context.Background()
	var flags []string
	t.generateAddressCmd.Run(ctx, flags)

	updatedUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where("id = ?", usdtUb.ID).First(updatedUsdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "ETHAddress", updatedUsdtUb.Address.String)

	updatedBtcUb := &userbalance.UserBalance{}
	err = t.db.Where("id = ?", btcUb.ID).First(updatedBtcUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	assert.Equal(t.T(), "BTCAddress", updatedBtcUb.Address.String)

}

func TestGenerateAddressCmd(t *testing.T) {
	suite.Run(t, &GenerateAddressCmd{
		Suite: new(suite.Suite),
	})
}
