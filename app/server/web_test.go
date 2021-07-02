package server_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jumpei00/gostocktrade/app/server"
	"github.com/sirupsen/logrus"

	"github.com/jumpei00/gostocktrade/app/models"
	"github.com/jumpei00/gostocktrade/app/models/indicator"
	"github.com/jumpei00/gostocktrade/stock"
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
	models.DB, _ = gorm.Open(sqlite.Open("web_test.sqlite3"), &gorm.Config{
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
	backTestParam.BackTest().CreateBacktestResult()
}

func (suite *ModelsTestSuite) TearDownTest() {
	models.AllDeleteCandles()
	models.DeleteBacktestResult("VOO")
}

func (suite *ModelsTestSuite) TearDownSuite() {
	os.Remove("web_test.sqlite3")
}

func (suite *ModelsTestSuite) TestCandleGetAPIHadler() {
	// normal access
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/candles?get=true&symbol=VOO&period=100", nil)
	server.CandleGetAPIHandler(recorder, req)
	resp := recorder.Result()

	dframe := models.DataFrame{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&dframe)

	suite.Equal(200, resp.StatusCode)
	suite.Equal("application/json", resp.Header.Get("Content-Type"))
	suite.Equal("VOO", dframe.CandleFrame.Symbol)
	suite.NotEmpty(dframe.CandleFrame.Candles)
	suite.NotEmpty(dframe.OptimizedParamFrame.Param)
	suite.Nil(dframe.SignalFrame)
	suite.NotEmpty(dframe.TradeFrame.Trade)

	// signal access
	recorder = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/candles?symbol=VOO&ema=true&bb=true&macd=true&rsi=true&willr=true", nil)
	server.CandleGetAPIHandler(recorder, req)
	resp = recorder.Result()

	dframe = models.DataFrame{}
	dec = json.NewDecoder(resp.Body)
	dec.Decode(&dframe)

	suite.Equal(200, resp.StatusCode)
	suite.Equal("application/json", resp.Header.Get("Content-Type"))
	suite.Nil(dframe.CandleFrame)
	suite.Nil(dframe.OptimizedParamFrame)
	suite.NotEmpty(dframe.SignalFrame.Signals.EmaSignals)
	suite.NotEmpty(dframe.SignalFrame.Signals.BBSignals)
	suite.NotEmpty(dframe.SignalFrame.Signals.MacdSignals)
	suite.NotEmpty(dframe.SignalFrame.Signals.RsiSignals)
	suite.NotEmpty(dframe.SignalFrame.Signals.WillrSignals)
	suite.Nil(dframe.TradeFrame)

	// when no backtest data, example GOOGL
	recorder = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/candles?get=true&symbol=GOOGL&period=100", nil)
	server.CandleGetAPIHandler(recorder, req)
	resp = recorder.Result()

	dframe = models.DataFrame{}
	dec = json.NewDecoder(resp.Body)
	dec.Decode(&dframe)

	suite.Equal(200, resp.StatusCode)
	suite.Equal("application/json", resp.Header.Get("Content-Type"))
	suite.Equal("GOOGL", dframe.CandleFrame.Symbol)
	suite.NotEmpty(dframe.CandleFrame.Candles)
	suite.Nil(dframe.OptimizedParamFrame)
	suite.Nil(dframe.SignalFrame)
	suite.Nil(dframe.TradeFrame)

	// wrong request, when no symbol
	recorder = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/candles?get=true&period=100", nil)
	server.CandleGetAPIHandler(recorder, req)
	resp = recorder.Result()
	body, _ := io.ReadAll(resp.Body)

	suite.Equal(400, resp.StatusCode)
	suite.Equal("{\"error\":\"bad parameter(symbol)\"}", string(body))

	// wrong request, when no period
	recorder = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/candles?get=true&symbol=VOO", nil)
	server.CandleGetAPIHandler(recorder, req)
	resp = recorder.Result()
	body, _ = io.ReadAll(resp.Body)

	suite.Equal(400, resp.StatusCode)
	suite.Equal("{\"error\":\"bad parameter(get, symbol)\"}", string(body))

	// wrong request, when wrong ticker symbol, example symbol=DAMYTEST
	recorder = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/candles?get=true&symbol=DAMYTEST&period=100", nil)
	server.CandleGetAPIHandler(recorder, req)
	resp = recorder.Result()
	body, _ = io.ReadAll(resp.Body)

	suite.Equal(400, resp.StatusCode)
	suite.Equal("{\"error\":\"stock get error, symbol: DAMYTEST\"}", string(body))
}

func (suite *ModelsTestSuite) TestBacktestAPIHandler() {
	// normal access
	recorder := httptest.NewRecorder()
	jsonData, _ := json.Marshal(backTestParam)
	req := httptest.NewRequest("POST", "/backtest", bytes.NewReader(jsonData))
	server.BacktestAPIHandler(recorder, req)
	resp := recorder.Result()

	dframe := models.DataFrame{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&dframe)

	suite.Equal(200, resp.StatusCode)
	suite.Equal("application/json", resp.Header.Get("Content-Type"))
	suite.Nil(dframe.CandleFrame)
	suite.NotEmpty(dframe.OptimizedParamFrame.Param)
	suite.Equal("VOO", dframe.OptimizedParamFrame.Param.Symbol)
	suite.NotEmpty(dframe.TradeFrame.Trade)
}

func TestModels(t *testing.T) {
	suite.Run(t, new(ModelsTestSuite))
}
