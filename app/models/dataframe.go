package models

import "github.com/markcheno/go-talib"

// DataFrame is json dataframe used for response from webserver
// After get data or calculate data, those are stored in this struct
type DataFrame struct {
	Candles []Candle `json:"candles"`
	Smas    []Sma    `json:"smas,omitempty"`
	Emas    []Ema    `json:"emas,omitempty"`
	BBands  *BBands  `json:"bbands,omitempty"`
	Macd    *Macd    `json:"macd,omitempty"`
	Rsi     *Rsi     `json:"rsi,omitempty"`
	WillR   *WillR   `json:"willr,omitempty"`
}

// Opens is open prices of candles
func (df *DataFrame) Opens() []float64 {
	open := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		open[i] = candle.Open
	}
	return open
}

// Highs is high prices of candles
func (df *DataFrame) Highs() []float64 {
	high := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		high[i] = candle.High
	}
	return high
}

// Lows is low prices of candles
func (df *DataFrame) Lows() []float64 {
	low := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		low[i] = candle.Low
	}
	return low
}

// Closes is close prices of candles
func (df *DataFrame) Closes() []float64 {
	close := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		close[i] = candle.Close
	}
	return close
}

// Volumes is volume prices of candles
func (df *DataFrame) Volumes() []float64 {
	volume := make([]float64, len(df.Candles))
	for i, candle := range df.Candles {
		volume[i] = candle.Volume
	}
	return volume
}

// AddSma adds Sma data in DataFrame.Smas
func (df *DataFrame) AddSma(period int) bool {
	if period > len(df.Candles) {
		return false
	}

	df.Smas = append(df.Smas, Sma{
		Period: period,
		Values: talib.Sma(df.Closes(), period),
	})
	return true
}

// AddEma adds Emas data in DataFrame.Emas
func (df *DataFrame) AddEma(period int) bool {
	if period > len(df.Candles) {
		return false
	}

	df.Emas = append(df.Emas, Ema{
		Period: period,
		Values: talib.Ema(df.Closes(), period),
	})
	return true
}

// AddBBands adds Boringer Bands data in DataFrame.BBands
func (df *DataFrame) AddBBands(N int, K float64) bool {
	if N > len(df.Candles) {
		return false
	}

	up, mid, low := talib.BBands(df.Closes(), N, K, K, 0)
	df.BBands = &BBands{
		N:   N,
		K:   K,
		Up:  up,
		Mid: mid,
		Low: low,
	}
	return true
}

// AddMacd adds Macd data in DataFrame.Macd
func (df *DataFrame) AddMacd(fast, slow, signal int) bool {
	if len(df.Candles) < 1 {
		return false
	}

	macd, macdSignal, macdHist := talib.Macd(df.Closes(), fast, slow, signal)
	df.Macd = &Macd{
		Fast:       fast,
		Slow:       slow,
		Signal:     signal,
		Macd:       macd,
		MacdSignal: macdSignal,
		MacdHist:   macdHist,
	}
	return true
}

// AddRsi adds Rsi data in DataFrame.Rsi
func (df *DataFrame) AddRsi(period int) bool {
	if period > len(df.Candles) {
		return false
	}

	df.Rsi = &Rsi{
		Period: period,
		Values: talib.Rsi(df.Closes(), period),
	}
	return true
}

// AddWillR adds WilliamR data in DataFrame.WillR
func (df *DataFrame) AddWillR(period int) bool {
	if period > len(df.Candles) {
		return false
	}

	df.WillR = &WillR{
		Period: period,
		Values: talib.WillR(df.Highs(), df.Lows(), df.Closes(), period),
	}
	return true
}
