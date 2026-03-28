package test

import (
	"context"
	"database/sql"
	"exchange-go/internal/command"
	"exchange-go/internal/di"
	"exchange-go/internal/userbalance"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UpdateUserWalletBalancesCmd struct {
	*suite.Suite
	updateUserWalletBalancesCmd command.ConsoleCommand
	db                          *gorm.DB
	redisClient                 *redis.Client
	userActor                   *userActor
}

func (t *UpdateUserWalletBalancesCmd) SetupSuite() {
	container := getContainer()
	t.updateUserWalletBalancesCmd = container.Get(di.UpdateUserWalletBalancesCommand).(command.ConsoleCommand)
	t.db = getDb()
	t.userActor = getUserActor()
}

func (t *UpdateUserWalletBalancesCmd) SetupTest() {}

func (t *UpdateUserWalletBalancesCmd) TearDownTest() {}

func (t *UpdateUserWalletBalancesCmd) TearDownSuite() {
	t.db.Where("user_id = ?", t.userActor.ID).Delete(userbalance.UserBalance{})
}

func (t *UpdateUserWalletBalancesCmd) TestRun() {
	usdtUb := &userbalance.UserBalance{
		UserID:         t.userActor.ID,
		CoinID:         1, //for usdt from currency seed
		Address:        sql.NullString{String: "ETHAddress", Valid: true},
		FrozenBalance:  "",
		BalanceCoin:    "USDT",
		Status:         userbalance.StatusEnabled,
		Amount:         "10000.00",
		FrozenAmount:   "1000.00",
		OtherAddresses: sql.NullString{String: "[{\"code\":\"TRX\",\"address\":\"TRXAddress\"}]", Valid: true},
	}

	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		Address:       sql.NullString{String: "BTCAddress", Valid: true},
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
	t.updateUserWalletBalancesCmd.Run(ctx, flags)

	usdtErc20UserWalletBalance := &userbalance.UserWalletBalance{}
	//3 is for ETH
	err = t.db.Where("user_id = ? and currency_id = ? and network_currency_id = ?", t.userActor.ID, usdtUb.CoinID, 3).First(usdtErc20UserWalletBalance).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "0.1", usdtErc20UserWalletBalance.Balance)

	usdtTrc20UserWalletBalance := &userbalance.UserWalletBalance{}
	//6 is for TRX
	err = t.db.Where("user_id = ? and currency_id = ? and network_currency_id = ?", t.userActor.ID, usdtUb.CoinID, 6).First(usdtTrc20UserWalletBalance).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "0.1", usdtTrc20UserWalletBalance.Balance)

	btcUserWalletBalance := &userbalance.UserWalletBalance{}
	//2 is for BTC
	err = t.db.Where("user_id = ? and currency_id = ? and network_currency_id = ?", t.userActor.ID, btcUb.CoinID, 2).First(btcUserWalletBalance).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "0.1", btcUserWalletBalance.Balance)

}

func TestUpdateUserWalletBalancesCmd(t *testing.T) {
	suite.Run(t, &UpdateUserWalletBalancesCmd{
		Suite: new(suite.Suite),
	})
}
