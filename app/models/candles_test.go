package models_test

import (
	"github.com/jumpei00/gostocktrade/app/models"
	"github.com/jumpei00/gostocktrade/stock"
)

func (suite *ModelsTestSuite) TestCreateCandles() {
	adjStock, _ := stock.GetStockData("VOO", 10, true)
	Stock, _ := stock.GetStockData("VOO", 10, false)
	candles := models.NewCandlesFromQuote(adjStock, Stock)

	suite.NotEmpty(candles)

	models.AllDeleteCandles()
	candles.CreateCandles()
}

func (suite *ModelsTestSuite) TestGetCandleFrame() {
	cframe := models.GetCandleFrame("VOO", 500)
	time := []int64{}
	for _, t := range cframe.Candles {
		time = append(time, t.Time)
	}

	suite.Equal("VOO", cframe.Symbol)
	suite.IsIncreasing(time)
}

func (suite *ModelsTestSuite) TestLastCandleTime() {
	cframe := models.GetCandleFrame("VOO", 500)
	lastTime := cframe.Candles[len(cframe.Candles)-1].Time
	lastCandleTime, err := models.LastCandleTime()

	suite.Equal(lastTime, lastCandleTime)
	suite.Nil(err)
}

func (suite *ModelsTestSuite) TestMatchTime() {
	cframe := models.GetCandleFrame("VOO", 500)
	firstCandle := cframe.Candles[0]
	lastCandle := cframe.Candles[len(cframe.Candles)-1]

	firstMatch, err1 := models.MatchTime(firstCandle.Time)
	lastMatch, err2 := models.MatchTime(lastCandle.Time)

	suite.Equal(firstCandle.ID, firstMatch)
	suite.Nil(err1)
	suite.Equal(lastCandle.ID, lastMatch)
	suite.Nil(err2)

	wrongMatch, err := models.MatchTime(firstCandle.Time+1)

	suite.Equal(0, wrongMatch)
	suite.NotNil(err)
}

func (suite *ModelsTestSuite) TestAllDeleteCandles() {
	models.AllDeleteCandles()
	cframe := models.GetCandleFrame("VOO", 10)

	suite.Empty(cframe.Candles)
}
