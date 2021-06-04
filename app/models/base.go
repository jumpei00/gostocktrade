package models

import (
	"github.com/jumpei00/gostocktrade/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is DBconnection
var DB *gorm.DB

// Candle is daily stock candledata, also used as json
type Candle struct {
	Time   int64   `json:"time"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

func init() {
	var err error
	DB, err = gorm.Open(sqlite.Open(config.Config.DBname), &gorm.Config{})
	if err != nil {
		logrus.Warnf("database open error: %v", err)
	}

	DB.AutoMigrate(&Candle{})
}
