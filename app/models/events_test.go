package models_test

import (
	"github.com/jumpei00/gostocktrade/app/models"
	"github.com/jumpei00/gostocktrade/app/models/indicator"
)

func (suite *ModelsTestSuite) TestGetTradeState() {
	// initializing
	suite.Op.CreateBacktestResult()

	// As test, create Ema signal due to doing same time to last candle time
	lastTime, _ := models.LastCandleTime()
	emaSignal := indicator.EmaSignal{
		Symbol: "VOO",
		Time:   lastTime,
		Price:  100,
		Action: indicator.BUY,
	}
	models.DB.Create(&emaSignal)

	// following, start test
	tradeFrame := models.GetTradeState("VOO")
	suite.Equal("BUY", tradeFrame.Trade.LastEmaTrade)
	suite.True(tradeFrame.Trade.IsEmaToday)

	// Delete Ema Signals
	models.DB.Delete(indicator.EmaSignal{}, "Symbol LIKE ?", "%VOO%")
	tradeFrame = models.GetTradeState("VOO")
	suite.Equal(indicator.NOTRADE, tradeFrame.Trade.LastEmaTrade)
	suite.False(tradeFrame.Trade.IsEmaToday)

	models.DeleteBacktestResult("VOO")
}

func (suite *ModelsTestSuite) TestGetSignalFrame() {
	// initializing
	suite.Op.CreateBacktestResult()

	signalFrame := models.GetSignalFrame("VOO", false, false, false, false, false)
	suite.Nil(signalFrame.Signals)

	signalFrame = models.GetSignalFrame("VOO", true, false, false, false, false)
	suite.NotNil(signalFrame.Signals)
	suite.NotEmpty(signalFrame.Signals.EmaSignals)

	signalFrame = models.GetSignalFrame("VOO", true, true, true, true, true)
	suite.NotNil(signalFrame.Signals)
	suite.NotEmpty(signalFrame.Signals.EmaSignals)
	suite.NotEmpty(signalFrame.Signals.BBSignals)
	suite.NotEmpty(signalFrame.Signals.MacdSignals)
	suite.NotEmpty(signalFrame.Signals.RsiSignals)
	suite.NotEmpty(signalFrame.Signals.WillrSignals)

	models.DeleteBacktestResult("VOO")
}

func (suite *ModelsTestSuite) TestSignalTest() {
	// no optimized params and signal data
	suite.False(models.SignalTest("VOO", 500))

	// initializing
	suite.Op.CreateBacktestResult()

	suite.True(models.SignalTest("VOO", 500))

	trades := models.GetTradeState("VOO").Trade
	signals := models.GetSignalFrame("VOO", true, true, true, true, true).Signals
	signalsLastTime := signals.LastSignalTimes()
	candleLastTime, _ := models.LastCandleTime()
	if len(signals.EmaSignals) != 0 {
		suite.Equal(signals.EmaSignals[len(signals.EmaSignals)-1].Action, trades.LastEmaTrade)
		suite.Equal(signalsLastTime["emaTime"] == candleLastTime, trades.IsEmaToday)
	} else {
		suite.Equal(indicator.NOTRADE, trades.LastEmaTrade)
		suite.False(trades.IsEmaToday)
	}

	if len(signals.BBSignals) != 0 {
		suite.Equal(signals.BBSignals[len(signals.BBSignals)-1].Action, trades.LastBBTrade)
		suite.Equal(signalsLastTime["bbTime"] == candleLastTime, trades.IsBBToday)
	} else {
		suite.Equal(indicator.NOTRADE, trades.LastBBTrade)
		suite.False(trades.IsBBToday)
	}

	if len(signals.MacdSignals) != 0 {
		suite.Equal(signals.MacdSignals[len(signals.MacdSignals)-1].Action, trades.LastMacdTrade)
		suite.Equal(signalsLastTime["macdTime"] == candleLastTime, trades.IsMacdToday)
	} else {
		suite.Equal(indicator.NOTRADE, trades.LastMacdTrade)
		suite.False(trades.IsMacdToday)
	}

	if len(signals.RsiSignals) != 0 {
		suite.Equal(signals.RsiSignals[len(signals.RsiSignals)-1].Action, trades.LastRsiTrade)
		suite.Equal(signalsLastTime["rsiTime"] == candleLastTime, trades.IsRsiToday)
	} else {
		suite.Equal(indicator.NOTRADE, trades.LastRsiTrade)
		suite.False(trades.IsRsiToday)
	}

	if len(signals.WillrSignals) != 0 {
		suite.Equal(signals.WillrSignals[len(signals.WillrSignals)-1].Action, trades.LastWillrTrade)
		suite.Equal(signalsLastTime["willrTime"] == candleLastTime, trades.IsWillrToday)
	} else {
		suite.Equal(indicator.NOTRADE, trades.LastWillrTrade)
		suite.False(trades.IsWillrToday)
	}

	models.DeleteBacktestResult("VOO")
}

func (suite *ModelsTestSuite) TestLastSignalTimes() {
	// initializing
	suite.Op.CreateBacktestResult()

	signalEvents := models.GetSignalFrame("VOO", true, true, true, true, true).Signals
	lastTimeMap := signalEvents.LastSignalTimes()

	suite.Equal(signalEvents.EmaSignals[len(signalEvents.EmaSignals)-1].Time, lastTimeMap["emaTime"])
	suite.Equal(signalEvents.BBSignals[len(signalEvents.BBSignals)-1].Time, lastTimeMap["bbTime"])
	suite.Equal(signalEvents.MacdSignals[len(signalEvents.MacdSignals)-1].Time, lastTimeMap["macdTime"])
	suite.Equal(signalEvents.RsiSignals[len(signalEvents.RsiSignals)-1].Time, lastTimeMap["rsiTime"])
	suite.Equal(signalEvents.WillrSignals[len(signalEvents.WillrSignals)-1].Time, lastTimeMap["willrTime"])

	models.DeleteBacktestResult("VOO")
}
