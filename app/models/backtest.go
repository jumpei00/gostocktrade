package models

import (
	"math"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/jumpei00/gostocktrade/app/models/indicator"
	"github.com/markcheno/go-talib"
)

// BackTestParam recieves some parameters used for backtest at json
type BackTestParam struct {
	Symbol string                        `json:"symbol"`
	Period int                           `json:"period"`
	Ema    *indicator.EmaBacktestParam   `json:"ema"`
	BB     *indicator.BBBacktestParam    `json:"bb"`
	Macd   *indicator.MacdBacktestParam  `json:"macd"`
	Rsi    *indicator.RsiBacktestParam   `json:"rsi"`
	Willr  *indicator.WillrBacktestParam `json:"willr"`
}

// OptimizedParam is stored to optimized parameter for backtest,
// also has relationships a part of signal results of backtest.
type OptimizedParam struct {
	ID               int                     `gorm:"primary_key" json:"-"`
	Timestamp        int64                   `json:"timestamp"`
	Symbol           string                  `json:"symbol"`
	EmaPerformance   float64                 `json:"ema_performance"`
	EmaShort         int                     `json:"ema_short"`
	EmaLong          int                     `json:"ema_long"`
	BBPerformance    float64                 `json:"bb_performance"`
	BBn              int                     `json:"bb_n"`
	BBk              float64                 `json:"bb_k"`
	MacdPerformance  float64                 `json:"macd_performance"`
	MacdFast         int                     `json:"macd_fast"`
	MacdSlow         int                     `json:"macd_slow"`
	MacdSignal       int                     `json:"macd_signal"`
	RsiPerformance   float64                 `json:"rsi_performance"`
	RsiPeriod        int                     `json:"rsi_period"`
	RsiBuyThread     float64                 `json:"rsi_buythread"`
	RsiSellThread    float64                 `json:"rsi_sellthread"`
	WillrPerformance float64                 `json:"willr_performance"`
	WillrPeriod      int                     `json:"willr_period"`
	WillrBuyThread   float64                 `json:"willr_buythread"`
	WillrSellThread  float64                 `json:"willr_sellthread"`
	EmaSignals       []indicator.EmaSignal   `gorm:"foreignKey:Symbol;references:Symbol" json:"-"`
	BBSignals        []indicator.BBSignal    `gorm:"foreignKey:Symbol;references:Symbol" json:"-"`
	MacdSignals      []indicator.MacdSignal  `gorm:"foreignKey:Symbol;references:Symbol" json:"-"`
	RsiSignals       []indicator.RsiSignal   `gorm:"foreignKey:Symbol;references:Symbol" json:"-"`
	WillrSignals     []indicator.WillrSignal `gorm:"foreignKey:Symbol;references:Symbol" json:"-"`
}

// CreateBacktestResult creates new backtest results, but before create, you delete existing data, beforehand
func (op *OptimizedParam) CreateBacktestResult() error {
	if err := DB.Create(op).Error; err != nil {
		return err
	}
	return nil
}

// DeleteBacktestResult deletes all exiting data for symbol
func DeleteBacktestResult(symbol string) {
	DB.Delete(OptimizedParam{}, "Symbol LIKE ?", "%"+symbol+"%")
	DB.Delete(indicator.EmaSignal{}, "Symbol LIKE ?", "%"+symbol+"%")
	DB.Delete(indicator.BBSignal{}, "Symbol LIKE ?", "%"+symbol+"%")
	DB.Delete(indicator.MacdSignal{}, "Symbol LIKE ?", "%"+symbol+"%")
	DB.Delete(indicator.RsiSignal{}, "Symbol LIKE ?", "%"+symbol+"%")
	DB.Delete(indicator.WillrSignal{}, "Symbol LIKE ?", "%"+symbol+"%")
}

// GetOptimizedParamFrame returns OptimizedParamFrame including OptimizedParam for symbol
func GetOptimizedParamFrame(symbol string) *OptimizedParamFrame {
	var op OptimizedParam
	var opframe OptimizedParamFrame

	err := DB.First(&op, OptimizedParam{Symbol: symbol})
	if err.Error != nil {
		// Not Found
		opframe.Param = nil
		return &opframe
	}

	opframe.Param = &op
	return &opframe
}

// BackTest excecutes backtest
func (bt *BackTestParam) BackTest() *OptimizedParam {
	DeleteBacktestResult(bt.Symbol)

	cframe := GetCandleFrame(bt.Symbol, bt.Period)
	logrus.Infof("backtest start: %v, %v", bt.Symbol, bt.Period)

	bpEma, bpEmaShort, bpEmaLong := cframe.optimizeEma(
		bt.Ema.EmaShortLow, bt.Ema.EmaShortHigh, bt.Ema.EmaLongLow, bt.Ema.EmaLongHigh)
	bpBB, bpBBn, bpBBk := cframe.optimizeBB(
		bt.BB.BBnLow, bt.BB.BBnHigh, bt.BB.BBkLow, bt.BB.BBkHigh)
	bpMacd, bpMacdFast, bpMacdSlow, bpMacdSignal := cframe.optimizeMacd(
		bt.Macd.MacdFastLow, bt.Macd.MacdFastHigh, bt.Macd.MacdSlowLow, bt.Macd.MacdSlowHigh,
		bt.Macd.MacdSignalLow, bt.Macd.MacdSignalHigh)
	bpRsi, bpRsiPeriod, bpRsiBuy, bpRsiSell := cframe.optimizeRsi(
		bt.Rsi.RsiPeriodLow, bt.Rsi.RsiPeriodHigh, bt.Rsi.RsiBuyThreadLow, bt.Rsi.RsiBuyThreadHigh,
		bt.Rsi.RsiSellThreadLow, bt.Rsi.RsiSellThreadHigh)
	bpWillr, bpWillrPeriod, bpWillrBuy, bpWillrSell := cframe.optimizeWillr(
		bt.Willr.WillrPeriodLow, bt.Willr.WillrPeriodHigh, bt.Willr.WillrBuyThreadLow, bt.Willr.WillrBuyThreadHigh,
		bt.Willr.WillrSellThreadLow, bt.Willr.WillrSellThreadHigh)

	op := OptimizedParam{
		Timestamp:        time.Now().Unix() * 1000,
		Symbol:           bt.Symbol,
		EmaPerformance:   math.Round(bpEma*100) / 100,
		EmaShort:         bpEmaShort,
		EmaLong:          bpEmaLong,
		BBPerformance:    math.Round(bpBB*100) / 100,
		BBn:              bpBBn,
		BBk:              math.Round(bpBBk*10) / 10,
		MacdPerformance:  math.Round(bpMacd*100) / 100,
		MacdFast:         bpMacdFast,
		MacdSlow:         bpMacdSlow,
		MacdSignal:       bpMacdSignal,
		RsiPerformance:   math.Round(bpRsi*100) / 100,
		RsiPeriod:        bpRsiPeriod,
		RsiBuyThread:     bpRsiBuy,
		RsiSellThread:    bpRsiSell,
		WillrPerformance: math.Round(bpWillr*100) / 100,
		WillrPeriod:      bpWillrPeriod,
		WillrBuyThread:   bpWillrBuy,
		WillrSellThread:  bpWillrSell,
		EmaSignals:       cframe.backtestEma(1, bpEmaShort, bpEmaLong).EmaSignals,
		BBSignals:        cframe.backtestBB(1, bpBBn, bpBBk).BBSignals,
		MacdSignals:      cframe.backtestMacd(1, bpMacdFast, bpMacdSlow, bpMacdSignal).MacdSignals,
		RsiSignals:       cframe.backtestRsi(1, bpRsiPeriod, bpRsiBuy, bpRsiSell).RsiSignals,
		WillrSignals:     cframe.backtestWillr(1, bpWillrPeriod, bpWillrBuy, bpWillrSell).WillrSignals,
	}

	return &op
}

func (cframe *CandleFrame) optimizeEma(
	lowShort, highShort, lowLong, highLong int) (bestPerformance float64, bestShort, bestLong int) {
	logrus.Infof("Ema backtest start: paramas -> %v, %v, %v %v", lowShort, highShort, lowLong, highLong)

	profit := 0.0
	bestShort = 7
	bestLong = 14

	for short := lowShort; short <= highShort; short++ {
		for long := lowLong; long <= highLong; long++ {
			signals := cframe.backtestEma(1, short, long)
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

func (cframe *CandleFrame) backtestEma(startDay, short int, long int) *indicator.EmaSignals {
	candles := cframe.Candles
	lenCandles := len(candles)

	if short >= lenCandles || long >= lenCandles {
		return nil
	}

	signals := indicator.EmaSignals{}
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
			signals := cframe.backtestBB(1, n, k)
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

func (cframe *CandleFrame) backtestBB(startDay, N int, K float64) *indicator.BBSignals {
	candles := cframe.Candles
	lenCandles := len(candles)

	if N >= lenCandles {
		return nil
	}

	signals := indicator.BBSignals{}
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
				signals := cframe.backtestMacd(1, fast, slow, signal)
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

func (cframe *CandleFrame) backtestMacd(startDay, fast, slow, signal int) *indicator.MacdSignals {
	candles := cframe.Candles
	lenCandles := len(candles)

	if fast >= lenCandles || slow >= lenCandles || signal >= lenCandles {
		return nil
	}

	signals := indicator.MacdSignals{}
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
				signals := cframe.backtestRsi(1, peirod, buyThread, sellThread)
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

func (cframe *CandleFrame) backtestRsi(startDay, period int, buyThread, sellThread float64) *indicator.RsiSignals {
	candles := cframe.Candles
	lenCandles := len(candles)

	if period >= lenCandles {
		return nil
	}

	signals := indicator.RsiSignals{}
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
				signals := cframe.backtestWillr(1, period, buyThread, sellThread)
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

func (cframe *CandleFrame) backtestWillr(startDay, period int, buyThread, sellThread float64) *indicator.WillrSignals {
	candles := cframe.Candles
	lenCandles := len(candles)

	if period >= lenCandles {
		return nil
	}

	signals := indicator.WillrSignals{}
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
