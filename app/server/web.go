package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/jumpei00/gostocktrade/app/models"
	"github.com/jumpei00/gostocktrade/config"
	"github.com/jumpei00/gostocktrade/stock"
	"github.com/sirupsen/logrus"
)

// JSONError is json error massage
type JSONError struct {
	Error string `json:"error"`
}

func errorAPI(w http.ResponseWriter, message string, code int) {
	jsonMessage, err := json.Marshal(JSONError{Error: message})
	if err != nil {
		logrus.Warnf("error message create error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(jsonMessage)
}

// IndexAPIHandler returns index.html contents,
// when path is "/"
func IndexAPIHandler(w http.ResponseWriter, req *http.Request) {
	temp := template.Must(template.ParseFiles("templates/index.html"))
	temp.ExecuteTemplate(w, "index.html", nil)
}

// CandleGetAPIHandler gets stock data, optimized paramerters, signal data, and trade data,
// when path is "/candles"
func CandleGetAPIHandler(w http.ResponseWriter, req *http.Request) {
	logrus.Infof("candle get request: url -> %s", req.URL)

	get, _ := strconv.ParseBool(req.URL.Query().Get("get"))
	symbol := req.URL.Query().Get("symbol")
	period, err := strconv.Atoi(req.URL.Query().Get("period"))

	if symbol == "" {
		errorAPI(w, "bad parameter(symbol)", http.StatusBadRequest)
		return
	}

	if get && err != nil {
		errorAPI(w, "bad parameter(get, symbol)", http.StatusBadRequest)
		return
	}

	dframe := models.NewDataFrame()

	// Downloads stock data
	if get {
		adjStock, _ := stock.GetStockData(symbol, period, true)
		Stock, _ := stock.GetStockData(symbol, period, false)
		if len(adjStock.Date) == 0 || len(Stock.Date) == 0 {
			logrus.Warnf("stock get error, symbol: %v", symbol)
			errorAPI(w, fmt.Sprintf("stock get error, symbol: %v", symbol), http.StatusBadRequest)
			return
		}
		// After delete existing data, store stock data in DB
		models.AllDeleteCandles()
		models.NewCandlesFromQuote(adjStock, Stock).CreateCandles()
		dframe.AddCandleFrame(symbol, period)
		dframe.AddOptimizedParamFrame(symbol)
		if models.SignalTest(symbol, period) {
			dframe.AddTradeFrame(symbol)
		}
	}

	ema, _ := strconv.ParseBool(req.URL.Query().Get("ema"))
	bb, _ := strconv.ParseBool(req.URL.Query().Get("bb"))
	macd, _ := strconv.ParseBool(req.URL.Query().Get("macd"))
	rsi, _ := strconv.ParseBool(req.URL.Query().Get("rsi"))
	willr, _ := strconv.ParseBool(req.URL.Query().Get("willr"))

	dframe.AddSignalFrame(symbol, ema, bb, macd, rsi, willr)

	js, err := json.Marshal(dframe)
	if err != nil {
		logrus.Warnf("candle json error: %v", err)
		errorAPI(w, "candle json error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// BacktestAPIHandler executes backtest, returns optimized parameters, trade data,
// when path is "/backtest"
func BacktestAPIHandler(w http.ResponseWriter, req *http.Request) {
	logrus.Info("backtest request")
	dec := json.NewDecoder(req.Body)

	var bt models.BackTestParam
	if err := dec.Decode(&bt); err != nil {
		logrus.Warnf("backtest params error: %v", err)
		errorAPI(w, fmt.Sprintf("backtest params error: %v", err), http.StatusInternalServerError)
		return
	}

	if err := bt.BackTest().CreateBacktestResult(); err != nil {
		logrus.Warnf("backtest error: %v", err)
		errorAPI(w, fmt.Sprintf("backtest error: %v", err), http.StatusInternalServerError)
		return
	}

	dframe := models.NewDataFrame()
	dframe.AddOptimizedParamFrame(bt.Symbol)
	dframe.AddTradeFrame(bt.Symbol)

	js, err := json.Marshal(dframe)
	if err != nil {
		logrus.Warnf("optimized params json error: %v", err)
		errorAPI(w, "optimized params json error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// Run starts webserver
func Run() {
	logrus.Info("server start")
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", IndexAPIHandler)
	http.HandleFunc("/candles", CandleGetAPIHandler)
	http.HandleFunc("/backtest", BacktestAPIHandler)
	logrus.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil))
}
