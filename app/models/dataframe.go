package models

import (
	"github.com/jumpei00/gostocktrade/app/models/indicator"
	"github.com/markcheno/go-talib"
	"github.com/sirupsen/logrus"
)

// DataFrame is data frame including candles, optimized parameters, signals
type DataFrame struct {
	*CandleFrame
	*OptimizedParamFrame
	*SignalFrame
	*TradeFrame
}

// NewDataFrame is constructor of DataFrame
func NewDataFrame() *DataFrame {
	return &DataFrame{}
}

// AddCandleFrame adds CandleFrame in DataFrame
func (dframe *DataFrame) AddCandleFrame(symbol string, limit int) {
	dframe.CandleFrame = GetCandleFrame(symbol, limit)
}

// AddSignalFrame adds SignalFrame in DataFrame
func (dframe *DataFrame) AddSignalFrame(symbol string, ema, bb, macd, rsi, willr bool) {
	dframe.SignalFrame = GetSignalFrame(symbol, ema, bb, macd, rsi, willr)
}

// AddOptimizedParamFrame adds OptimizedParamFrame in DataFrame
func (dframe *DataFrame) AddOptimizedParamFrame(symbol string) {
	dframe.OptimizedParamFrame = GetOptimizedParamFrame(symbol)
}

// AddTradeFrame adds TradeFrame in DataFrame
func (dframe *DataFrame) AddTradeFrame(symbol string) {
	dframe.TradeFrame = GetTradeState(symbol)
}

// SignalFrame is dataframe of SignalEvents
type SignalFrame struct {
	Signals *SignalEvents `json:"signals,omitempty"`
}

// OptimizedParamFrame is optimized params data frame
type OptimizedParamFrame struct {
	Param *OptimizedParam `json:"optimized_params,omitempty"`
}

// TradeFrame is Trade frame
type TradeFrame struct {
	Trade *Trade `json:"trade,omitempty"`
}

// CandleFrame is candle data frame
type CandleFrame struct {
	Symbol  string   `json:"symbol,omitempty"`
	Candles []Candle `json:"candles,omitempty"`
}

// Opens is open prices of candles
func (cframe *CandleFrame) Opens() []float64 {
	open := make([]float64, len(cframe.Candles))
	for i, candle := range cframe.Candles {
		open[i] = candle.Open
	}
	return open
}

// Highs is high prices of candles
func (cframe *CandleFrame) Highs() []float64 {
	high := make([]float64, len(cframe.Candles))
	for i, candle := range cframe.Candles {
		high[i] = candle.High
	}
	return high
}

// Lows is low prices of candles
func (cframe *CandleFrame) Lows() []float64 {
	low := make([]float64, len(cframe.Candles))
	for i, candle := range cframe.Candles {
		low[i] = candle.Low
	}
	return low
}

// Closes is close prices of candles
func (cframe *CandleFrame) Closes() []float64 {
	close := make([]float64, len(cframe.Candles))
	for i, candle := range cframe.Candles {
		close[i] = candle.Close
	}
	return close
}

// Volumes is volume prices of candles
func (cframe *CandleFrame) Volumes() []float64 {
	volume := make([]float64, len(cframe.Candles))
	for i, candle := range cframe.Candles {
		volume[i] = candle.Volume
	}
	return volume
}

// following, using for backtest
func (cframe *CandleFrame) optimizeEma(
	lowShort, highShort, lowLong, highLong int) (bestPerformance float64, bestShort, bestLong int) {
	logrus.Infof("Ema backtest start: paramas -> %v, %v, %v %v", lowShort, highShort, lowLong, highLong)

	profit := 0.0
	bestShort = 7
	bestLong = 14

	for short := lowShort; short <= highShort; short++ {
		for long := lowLong; long <= highLong; long++ {
			signals := cframe.backtestEma(1, short, long, nil)
			if signals == nil {
				continue
			}

			profit = signals.Profit()
			if bestPerformance < profit {
				bestPerformance = profit
				bestShort = short
				bestLong = long
			}
		}
	}

	logrus.Infof("Ema backtest end: results -> %v, %v, %v", bestPerformance, bestShort, bestLong)
	return bestPerformance, bestShort, bestLong
}

func (cframe *CandleFrame) backtestEma(startDay, short int, long int, lastSignal *indicator.EmaSignal) *indicator.EmaSignals {
	candles := cframe.Candles
	lenCandles := len(candles)

	if short >= lenCandles || long >= lenCandles {
		return nil
	}

	signals := indicator.EmaSignals{}
	// using at SignalTest
	if lastSignal != nil {
		signals.EmaSignals = append(signals.EmaSignals, *lastSignal)
	}

	shortEma := talib.Ema(cframe.Closes(), short)
	longEma := talib.Ema(cframe.Closes(), long)

	for day := startDay; day < lenCandles; day++ {
		if day < short || day < long {
			continue
		}

		if shortEma[day-1] < longEma[day-1] && shortEma[day] >= longEma[day] {
			signals.Buy(cframe.Symbol, candles[day].Time, candles[day].Close)
		}

		if shortEma[day-1] > longEma[day-1] && shortEma[day] <= longEma[day] {
			signals.Sell(cframe.Symbol, candles[day].Time, candles[day].Close)
		}
	}

	return &signals
}

func (cframe *CandleFrame) optimizeBB(
	lowN, highN int, lowK, highK float64) (bestPerformance float64, bestN int, bestK float64) {
	logrus.Infof("BB backtest start: paramas -> %v, %v, %v %v", lowN, highN, lowK, highK)

	profit := 0.0
	bestN = 20
	bestK = 2.0

	for n := lowN; n <= highN; n++ {
		for k := lowK; k <= highK; k += 0.1 {
			signals := cframe.backtestBB(1, n, k, nil)
			if signals == nil {
				continue
			}
			profit = signals.Profit()
			if bestPerformance < profit {
				bestPerformance = profit
				bestN = n
				bestK = k
			}
		}
	}

	logrus.Infof("BB backtest end: results -> %v, %v, %v", bestPerformance, bestN, bestK)
	return bestPerformance, bestN, bestK
}

func (cframe *CandleFrame) backtestBB(startDay, N int, K float64, lastSignal *indicator.BBSignal) *indicator.BBSignals {
	candles := cframe.Candles
	lenCandles := len(candles)

	if N >= lenCandles {
		return nil
	}

	signals := indicator.BBSignals{}
	// using at SignalTest
	if lastSignal != nil {
		signals.BBSignals = append(signals.BBSignals, *lastSignal)
	}

	upBand, _, lowBand := talib.BBands(cframe.Closes(), N, K, K, 0)

	for day := startDay; day < lenCandles; day++ {
		if day < N {
			continue
		}

		if candles[day-1].Close < lowBand[day-1] && candles[day].Close >= lowBand[day] {
			signals.Buy(cframe.Symbol, candles[day].Time, candles[day].Close)
		}

		if candles[day-1].Close > upBand[day-1] && candles[day].Close <= upBand[day] {
			signals.Sell(cframe.Symbol, candles[day].Time, candles[day].Close)
		}
	}

	return &signals
}

func (cframe *CandleFrame) optimizeMacd(
	lowFast, highFast, lowSlow, highSlow, lowSignal, highSignal int) (bestPerformance float64, bestFast, bestSlow, bestSignal int) {
	logrus.Infof("Macd backtest start: paramas -> %v, %v, %v %v, %v, %v", lowFast, highFast, lowSlow, highSlow, lowSignal, highSignal)

	profit := 0.0
	bestFast = 12
	bestSlow = 26
	bestSignal = 9

	for fast := lowFast; fast <= highFast; fast++ {
		for slow := lowSlow; slow <= highSlow; slow++ {
			for signal := lowSignal; signal <= highSignal; signal++ {
				signals := cframe.backtestMacd(1, fast, slow, signal, nil)
				if signals == nil {
					continue
				}
				profit = signals.Profit()
				if bestPerformance < profit {
					bestPerformance = profit
					bestFast = fast
					bestSlow = slow
					bestSignal = signal
				}

			}
		}
	}

	logrus.Infof("Macd backtest end: results -> %v, %v, %v %v", bestPerformance, bestFast, bestSlow, bestSignal)
	return bestPerformance, bestFast, bestSlow, bestSignal
}

func (cframe *CandleFrame) backtestMacd(startDay, fast, slow, signal int, lastSignal *indicator.MacdSignal) *indicator.MacdSignals {
	candles := cframe.Candles
	lenCandles := len(candles)

	if fast >= lenCandles || slow >= lenCandles || signal >= lenCandles {
		return nil
	}

	signals := indicator.MacdSignals{}
	// using at SignalTest
	if lastSignal != nil {
		signals.MacdSignals = append(signals.MacdSignals, *lastSignal)
	}

	macd, macdSignal, _ := talib.Macd(cframe.Closes(), fast, slow, signal)

	for day := startDay; day < lenCandles; day++ {
		if macd[day] < 0 && macdSignal[day] < 0 &&
			macd[day-1] < macdSignal[day-1] &&
			macd[day] >= macdSignal[day] {
			signals.Buy(cframe.Symbol, candles[day].Time, candles[day].Close)
		}

		if macd[day] > 0 && macdSignal[day] > 0 &&
			macd[day-1] > macdSignal[day-1] &&
			macd[day] <= macdSignal[day] {
			signals.Sell(cframe.Symbol, candles[day].Time, candles[day].Close)
		}
	}

	return &signals
}

func (cframe *CandleFrame) optimizeRsi(
	lowPeriod, highPeriod int,
	lowBuyThread, highBuyThread, lowSellThread, highSellThread float64) (bestPerformance float64, bestPeriod int, bestBuyThread, bestSellThread float64) {
	logrus.Infof("Rsi backtest start: paramas -> %v, %v, %v %v, %v, %v", lowPeriod, highPeriod, lowBuyThread, highBuyThread, lowSellThread, highSellThread)

	profit := 0.0
	bestPeriod = 14
	bestBuyThread = 30.0
	bestSellThread = 70.0

	for peirod := lowPeriod; peirod <= highPeriod; peirod++ {
		for buyThread := lowBuyThread; buyThread <= highBuyThread; buyThread++ {
			for sellThread := lowSellThread; sellThread <= highSellThread; sellThread++ {
				signals := cframe.backtestRsi(1, peirod, buyThread, sellThread, nil)
				if signals == nil {
					continue
				}
				profit = signals.Profit()
				if bestPerformance < profit {
					bestPerformance = profit
					bestPeriod = peirod
					bestBuyThread = buyThread
					bestSellThread = sellThread
				}
			}
		}
	}

	logrus.Infof("Rsi backtest end: results -> %v, %v, %v %v", bestPerformance, bestPeriod, bestBuyThread, bestSellThread)
	return bestPerformance, bestPeriod, bestBuyThread, bestSellThread
}

func (cframe *CandleFrame) backtestRsi(startDay, period int, buyThread, sellThread float64, lastSignal *indicator.RsiSignal) *indicator.RsiSignals {
	candles := cframe.Candles
	lenCandles := len(candles)

	if period >= lenCandles {
		return nil
	}

	signals := indicator.RsiSignals{}
	// using at SignalTest
	if lastSignal != nil {
		signals.RsiSignals = append(signals.RsiSignals, *lastSignal)
	}

	rsi := talib.Rsi(cframe.Closes(), period)

	for day := startDay; day < lenCandles; day++ {
		if rsi[day-1] == 0 || rsi[day-1] == 100 {
			continue
		}

		if rsi[day-1] < buyThread && rsi[day] >= buyThread {
			signals.Buy(cframe.Symbol, candles[day].Time, candles[day].Close)
		}

		if rsi[day-1] > sellThread && rsi[day] <= sellThread {
			signals.Sell(cframe.Symbol, candles[day].Time, candles[day].Close)
		}
	}

	return &signals
}

func (cframe *CandleFrame) optimizeWillr(
	lowPeriod, highPeriod int,
	lowBuyThread, highBuyThread, lowSellThread, highSellThread float64) (bestPerformance float64, bestPeriod int, bestBuyThread, bestSellThread float64) {
	logrus.Infof("Willr backtest start: paramas -> %v, %v, %v %v, %v, %v", lowPeriod, highPeriod, lowBuyThread, highBuyThread, lowSellThread, highSellThread)

	profit := 0.0
	bestPeriod = 10
	bestBuyThread = -20.0
	bestSellThread = -80.0

	for period := lowPeriod; period <= highPeriod; period++ {
		for buyThread := lowBuyThread; buyThread <= highBuyThread; buyThread++ {
			for sellThread := lowSellThread; sellThread <= highSellThread; sellThread++ {
				signals := cframe.backtestWillr(1, period, buyThread, sellThread, nil)
				if signals == nil {
					continue
				}
				profit = signals.Profit()
				if bestPerformance < profit {
					bestPerformance = profit
					bestPeriod = period
					bestBuyThread = buyThread
					bestSellThread = sellThread
				}
			}
		}
	}

	logrus.Infof("Willr backtest end: results -> %v, %v, %v %v", bestPerformance, bestPeriod, bestBuyThread, bestSellThread)
	return bestPerformance, bestPeriod, bestBuyThread, bestSellThread
}

func (cframe *CandleFrame) backtestWillr(startDay, period int, buyThread, sellThread float64, lastSignal *indicator.WillrSignal) *indicator.WillrSignals {
	candles := cframe.Candles
	lenCandles := len(candles)

	if period >= lenCandles {
		return nil
	}

	signals := indicator.WillrSignals{}
	// using at SignalTest
	if lastSignal != nil {
		signals.WillrSignals = append(signals.WillrSignals, *lastSignal)
	}

	willr := talib.WillR(cframe.Highs(), cframe.Lows(), cframe.Closes(), period)

	for day := startDay; day < lenCandles; day++ {
		if willr[day-1] == 0 || willr[day-1] == -100 {
			continue
		}

		if willr[day-1] < buyThread && willr[day] >= buyThread {
			signals.Buy(cframe.Symbol, candles[day].Time, candles[day].Close)
		}

		if willr[day-1] > sellThread && willr[day] <= sellThread {
			signals.Sell(cframe.Symbol, candles[day].Time, candles[day].Close)
		}
	}

	return &signals
}
