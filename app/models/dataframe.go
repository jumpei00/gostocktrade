package models

import (
	"github.com/jumpei00/gostocktrade/app/models/indicator"
	"github.com/markcheno/go-talib"
)

// DataFrame is data frame including candles, optimized parameters, signals
type DataFrame struct {
	*CandleFrame
	*SignalFrame
	*OptimizedParam
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

// SignalFrame is dataframe of SignalEvents
type SignalFrame struct {
	Signals *SignalEvents `json:"signals,omitempty"`
}

// OptimizedParamFrame is optimized params data frame
type OptimizedParamFrame struct {
	Param *OptimizedParam `json:"optimized_params,omitempty"`
}

// IndicatorFrame is json frame of indicator data
// After calculate data, those are stored in this struct
type IndicatorFrame struct {
	*CandleFrame
	Smas   []indicator.Sma   `json:"smas,omitempty"`
	Emas   []indicator.Ema   `json:"emas,omitempty"`
	BBands *indicator.BBands `json:"bbands,omitempty"`
	Macd   *indicator.Macd   `json:"macd,omitempty"`
	Rsi    *indicator.Rsi    `json:"rsi,omitempty"`
	WillR  *indicator.Willr  `json:"willr,omitempty"`
}

// NewIndicator is constractor of IndicatorFrame,
// and embeded CandleFrame, but not json
func NewIndicator(symbol string, limit int) *IndicatorFrame {
	iframe := IndicatorFrame{
		CandleFrame: GetCandleFrame(symbol, limit),
	}
	return &iframe
}

// AddSma adds Sma data in IndicatorFrame.Smas
func (iframe *IndicatorFrame) AddSma(period int) bool {
	if period > len(iframe.Candles) {
		return false
	}

	iframe.Smas = append(iframe.Smas, indicator.Sma{
		Period: period,
		Values: talib.Sma(iframe.Closes(), period),
	})
	return true
}

// AddEma adds Emas data in IndicatorFrame.Emas
func (iframe *IndicatorFrame) AddEma(period int) bool {
	if period > len(iframe.Candles) {
		return false
	}

	iframe.Emas = append(iframe.Emas, indicator.Ema{
		Period: period,
		Values: talib.Ema(iframe.Closes(), period),
	})
	return true
}

// AddBBands adds Boringer Bands data in IndicatorFrame.BBands
func (iframe *IndicatorFrame) AddBBands(N int, K float64) bool {
	if N > len(iframe.Candles) {
		return false
	}

	up, mid, low := talib.BBands(iframe.Closes(), N, K, K, 0)
	iframe.BBands = &indicator.BBands{
		N:   N,
		K:   K,
		Up:  up,
		Mid: mid,
		Low: low,
	}
	return true
}

// AddMacd adds Macd data in IndicatorFrame.Macd
func (iframe *IndicatorFrame) AddMacd(fast, slow, signal int) bool {
	if len(iframe.Candles) < 1 {
		return false
	}

	macd, macdSignal, macdHist := talib.Macd(iframe.Closes(), fast, slow, signal)
	iframe.Macd = &indicator.Macd{
		Fast:       fast,
		Slow:       slow,
		Signal:     signal,
		Macd:       macd,
		MacdSignal: macdSignal,
		MacdHist:   macdHist,
	}
	return true
}

// AddRsi adds Rsi data in IndicatorFrame.Rsi
func (iframe *IndicatorFrame) AddRsi(period int) bool {
	if period > len(iframe.Candles) {
		return false
	}

	iframe.Rsi = &indicator.Rsi{
		Period: period,
		Values: talib.Rsi(iframe.Closes(), period),
	}
	return true
}

// AddWillr adds WilliamR data in IndicatorFrame.WillR
func (iframe *IndicatorFrame) AddWillr(period int) bool {
	if period > len(iframe.Candles) {
		return false
	}

	iframe.WillR = &indicator.Willr{
		Period: period,
		Values: talib.WillR(iframe.Highs(), iframe.Lows(), iframe.Closes(), period),
	}
	return true
}
