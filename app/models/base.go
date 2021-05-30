package models

import (
	"time"

	"github.com/jumpei00/gostocktrade/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is DBconnection
var DB *gorm.DB


// Candle is daily stock candledata
type Candle struct {
	gorm.Model
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

func init() {
	var err error
	DB, err = gorm.Open(sqlite.Open(config.Config.DBname), &gorm.Config{})
	if err != nil {
		logrus.Warnf("database open error: %v", err)
	}

	DB.AutoMigrate(&Candle{})
}
