package test

import (
	"context"
	"exchange-go/internal/command"
	"exchange-go/internal/currency"
	"exchange-go/internal/di"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type GenerateKlineSyncCmd struct {
	*suite.Suite
	generateKlineSyncCmd command.ConsoleCommand
	db                   *gorm.DB
	userActor            *userActor
}

func (t *GenerateKlineSyncCmd) SetupSuite() {
	container := getContainer()
	t.generateKlineSyncCmd = container.Get(di.GenerateKlineSyncCommand).(command.ConsoleCommand)
	t.db = getDb()
	t.userActor = getUserActor()
}

func (t *GenerateKlineSyncCmd) SetupTest() {
	t.db.Where("id > ?", 0).Delete(currency.KlineSync{})
}

func (t *GenerateKlineSyncCmd) TearDownTest() {
	t.db.Where("id > ?", 0).Delete(currency.KlineSync{})
}

func (t *GenerateKlineSyncCmd) TearDownSuite() {
}

func (t *GenerateKlineSyncCmd) TestRun() {
	ctx := context.Background()
	var flags []string
	t.generateKlineSyncCmd.Run(ctx, flags)

	var klineSyncs []currency.KlineSync
	err := t.db.Find(&klineSyncs).Error
	if err != nil {
		t.Fail(err.Error())
	}

	now := time.Now()
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	lastDay := now.Add(-24 * time.Hour)
	startTime := time.Date(lastDay.Year(), lastDay.Month(), lastDay.Day(), 0, 0, 0, 0, lastDay.Location())

	assert.Equal(t.T(), 24, len(klineSyncs)) //because we have 6 pair in  test env (look at currencyseed.go) and 4 time frame
	for _, item := range klineSyncs {
		assert.Equal(t.T(), false, item.WithUpdate)
		assert.Equal(t.T(), currency.SyncTypeAuto, item.Type)
		assert.Equal(t.T(), startTime.Format("2006-01-02 15:04:05"), item.StartTime.Time.Format("2006-01-02 15:04:05"))
		assert.Equal(t.T(), endTime.Format("2006-01-02 15:04:05"), item.EndTime.Time.Format("2006-01-02 15:04:05"))
	}

}

func TestGenerateKlineSyncCmd(t *testing.T) {
	suite.Run(t, &GenerateKlineSyncCmd{
		Suite: new(suite.Suite),
	})
}
