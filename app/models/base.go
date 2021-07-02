package models

import (
	"github.com/jumpei00/gostocktrade/app/models/indicator"
	"github.com/jumpei00/gostocktrade/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is DBconnection
var DB *gorm.DB

// InitDB initializes DB
func InitDB() {
	var err error

	DB, err = gorm.Open(sqlite.Open(config.Config.DBname), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		logrus.Warnf("database open error: %v", err)
	}

	DB.AutoMigrate(
		&Candle{},
		&OptimizedParam{},
		&indicator.EmaSignal{},
		&indicator.BBSignal{},
		&indicator.MacdSignal{},
		&indicator.RsiSignal{},
		&indicator.WillrSignal{},
	)
}
