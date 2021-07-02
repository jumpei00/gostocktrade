package models_test

import (
	"github.com/jumpei00/gostocktrade/app/models"
)

func (suite *ModelsTestSuite) TestDataFrame() {
	// initializing
	suite.Op.CreateBacktestResult()

	dframe := models.NewDataFrame()
	suite.Nil(dframe.CandleFrame)
	suite.Nil(dframe.OptimizedParamFrame)
	suite.Nil(dframe.SignalFrame)
	suite.Nil(dframe.TradeFrame)

	dframe.AddCandleFrame("VOO", 100)
	suite.Equal("VOO", dframe.CandleFrame.Symbol)
	suite.Len(dframe.CandleFrame.Candles, 100)

	dframe.AddSignalFrame("VOO", true, false, false, false, false)
	suite.NotEmpty(dframe.SignalFrame.Signals.EmaSignals)
	suite.Empty(dframe.SignalFrame.Signals.BBSignals)
	suite.Empty(dframe.SignalFrame.Signals.MacdSignals)
	suite.Empty(dframe.SignalFrame.Signals.RsiSignals)
	suite.Empty(dframe.SignalFrame.Signals.WillrSignals)

	dframe.AddOptimizedParamFrame("DAMY")
	suite.Nil(dframe.OptimizedParamFrame.Param)
	dframe.AddOptimizedParamFrame("VOO")
	suite.NotEmpty(dframe.OptimizedParamFrame.Param)

	dframe.AddTradeFrame("VOO")
	suite.NotEmpty(dframe.TradeFrame.Trade)
}