package platform

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func NewDbClient(c Configs) (*gorm.DB, error) {
	//dbName := c.GetString("db.name")
	dsn := c.GetString("db.dsn")

	gormConfig := &gorm.Config{
		SkipDefaultTransaction: true,
	}
	env := c.GetEnv()
	if env == EnvProd || env == EnvTest {
		gormConfig.Logger = gormLogger.Default.LogMode(gormLogger.Silent)
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)

	if err != nil {
		return nil, err
	}
	//defer db.Close()
	return db, nil
}
