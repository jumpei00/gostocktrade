package models_test

import (
	"os"
	"testing"

	"github.com/jumpei00/gostocktrade/stock"
	"github.com/sirupsen/logrus"

	"github.com/jumpei00/gostocktrade/app/models"
	"github.com/jumpei00/gostocktrade/app/models/indicator"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var backTestParam = models.BackTestParam{
	Symbol: "VOO",
	Period: 500,
	Ema: &indicator.EmaBacktestParam{
		EmaShortLow:  5,
		EmaShortHigh: 15,
		EmaLongLow:   15,
		EmaLongHigh:  30,
	},
	BB: &indicator.BBBacktestParam{
		BBnLow:  10,
		BBnHigh: 30,
		BBkLow:  1.5,
		BBkHigh: 2.5,
	},
	Macd: &indicator.MacdBacktestParam{
		MacdFastLow:    5,
		MacdFastHigh:   20,
		MacdSlowLow:    20,
		MacdSlowHigh:   35,
		MacdSignalLow:  5,
		MacdSignalHigh: 20,
	},
	Rsi: &indicator.RsiBacktestParam{
		RsiPeriodLow:      5,
		RsiPeriodHigh:     50,
		RsiBuyThreadLow:   20,
		RsiBuyThreadHigh:  35,
		RsiSellThreadLow:  65,
		RsiSellThreadHigh: 80,
	},
	Willr: &indicator.WillrBacktestParam{
		WillrPeriodLow:      5,
		WillrPeriodHigh:     50,
		WillrBuyThreadLow:   -90,
		WillrBuyThreadHigh:  -75,
		WillrSellThreadLow:  -25,
		WillrSellThreadHigh: -10,
	},
}

type ModelsTestSuite struct {
	suite.Suite
	Candles *models.Candles
	Op      *models.OptimizedParam
}

func (suite *ModelsTestSuite) SetupSuite() {
	logrus.SetLevel(logrus.ErrorLevel)
	models.DB, _ = gorm.Open(sqlite.Open("models_test.sqlite3"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	models.DB.AutoMigrate(
		&models.Candle{},
		&models.OptimizedParam{},
		&indicator.EmaSignal{},
		&indicator.BBSignal{},
		&indicator.MacdSignal{},
		&indicator.RsiSignal{},
		&indicator.WillrSignal{},
	)

	adjStock, _ := stock.GetStockData("VOO", 500, true)
	Stock, _ := stock.GetStockData("VOO", 500, false)
	suite.Candles = models.NewCandlesFromQuote(adjStock, Stock)
}

func (suite *ModelsTestSuite) SetupTest() {
	suite.Candles.CreateCandles()
	suite.Op = backTestParam.BackTest()
}

func (suite *ModelsTestSuite) TearDownTest() {
	models.AllDeleteCandles()
}

func (suite *ModelsTestSuite) TearDownSuite() {
	os.Remove("models_test.sqlite3")
}

func TestModels(t *testing.T) {
	suite.Run(t, new(ModelsTestSuite))
}
