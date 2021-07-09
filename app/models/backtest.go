package models

import (
	"math"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/jumpei00/gostocktrade/app/models/indicator"
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

// BackTest excecutes backtest
// Caution, the Symbol in BackTestParam is the same to ticker symbol of the candle data,
// if those are different, deal with frontend process
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
		EmaSignals:       cframe.backtestEma(1, bpEmaShort, bpEmaLong, nil).EmaSignals,
		BBSignals:        cframe.backtestBB(1, bpBBn, bpBBk, nil).BBSignals,
		MacdSignals:      cframe.backtestMacd(1, bpMacdFast, bpMacdSlow, bpMacdSignal, nil).MacdSignals,
		RsiSignals:       cframe.backtestRsi(1, bpRsiPeriod, bpRsiBuy, bpRsiSell, nil).RsiSignals,
		WillrSignals:     cframe.backtestWillr(1, bpWillrPeriod, bpWillrBuy, bpWillrSell, nil).WillrSignals,
	}

	return &op
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

// CreateBacktestResult creates new backtest results, but before create, you delete existing data, beforehand
func (op *OptimizedParam) CreateBacktestResult() error {
	if err := DB.Create(op).Error; err != nil {
		return err
	}
	return nil
}
