package test

import (
	"context"
	"database/sql"
	"exchange-go/internal/command"
	"exchange-go/internal/currency"
	"exchange-go/internal/di"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type SyncKlineCmd struct {
	*suite.Suite
	syncKlineCmd command.ConsoleCommand
	db           *gorm.DB
	userActor    *userActor
}

func (t *SyncKlineCmd) SetupSuite() {
	container := getContainer()
	t.syncKlineCmd = container.Get(di.KlineSyncCommand).(command.ConsoleCommand)
	t.db = getDb()
	t.userActor = getUserActor()
}

func (t *SyncKlineCmd) SetupTest() {
}

func (t *SyncKlineCmd) TearDownTest() {
}

func (t *SyncKlineCmd) TearDownSuite() {
}

func (t *SyncKlineCmd) TestRunWithoutFlags() {
	//insert klineSync in database
	st, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 10:10:10")
	et, _ := time.Parse("2006-01-02 15:04:05", "2021-01-02 10:10:10")
	ks := &currency.KlineSync{
		ID:         1,
		PairID:     1,
		TimeFrame:  sql.NullString{String: "1minute", Valid: true},
		StartTime:  sql.NullTime{Time: st, Valid: true},
		EndTime:    sql.NullTime{Time: et, Valid: true},
		Status:     currency.SyncStatusCreated,
		Type:       "AUTO",
		WithUpdate: false,
	}
	err := t.db.Create(ks).Error

	if err != nil {
		t.Fail(err.Error())
	}

	var flags []string
	t.syncKlineCmd.Run(context.Background(), flags)

	updatedKs := &currency.KlineSync{}
	_ = t.db.Where("id =?", ks.ID).Find(&updatedKs).Error
	assert.Equal(t.T(), currency.SyncStatusDone, updatedKs.Status)

}

func (t *SyncKlineCmd) TestRunWithFlags() {
	var flags []string
	flags = append(flags, "-pair=BTC-USDT")
	flags = append(flags, "-frame=1minute")
	flags = append(flags, "-start=2021-01-01 10:10:10")
	flags = append(flags, "-end=2021-01-02 10:10:10")
	flags = append(flags, "-update=true")
	t.syncKlineCmd.Run(context.Background(), flags)
}

func TestSyncKlineCmd(t *testing.T) {
	suite.Run(t, &SyncKlineCmd{
		Suite: new(suite.Suite),
	})
}
