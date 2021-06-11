package models

import (
	"math"
	"sort"

	"github.com/markcheno/go-quote"
	"gorm.io/gorm"
)

// Candles is slice of Candle
// Using this, create candle data in database
type Candles []Candle

// Candle is daily stock candledata, also used as json
type Candle struct {
	Time   int64   `json:"time"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

// NewCandlesFromQuote converts Quote to slice of Candle due to creating in database,
// ex) [Date[1, 2, 3...], Open[1, 2, 3...]...] → [[Date[1], Open[1]...], [Date[2], Open[2]...]...]
// and return pointer of Candles(used as constructor)
// Because of using for frondend, this method also converts time to Unixtime
func NewCandlesFromQuote(stock *quote.Quote) *Candles {
	candles := Candles{}
	for i := 0; i < len(stock.Date); i++ {
		candles = append(candles, Candle{
			Time:   stock.Date[i].Unix() * 1000,
			Open:   (math.Round(stock.Open[i]*100) / 100),
			High:   (math.Round(stock.High[i]*100) / 100),
			Low:    (math.Round(stock.Low[i]*100) / 100),
			Close:  (math.Round(stock.Close[i]*100) / 100),
			Volume: (math.Round(stock.Volume[i]*100) / 100),
		})
	}

	return &candles
}

// CreateCandles creates candle data
func (cs *Candles) CreateCandles() {
	DB.Create(cs)
}

// AllDeleteCandles deletes all data of "candles" table
func AllDeleteCandles() {
	DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Candle{})
}

// GetCandleFrame gets candle data for limit by descending
// After get data, return DataFrame stored in data
func GetCandleFrame(symbol string, limit int) *CandleFrame {
	var candles Candles
	DB.Order("time desc").Limit(limit).Find(&candles)
	sort.Slice(candles, func(i, j int) bool { return candles[i].Time < candles[j].Time })

	cframe := CandleFrame{}
	cframe.Symbol = symbol
	cframe.Candles = candles

	return &cframe
}

// GetDataFrameInCandles is warapper of GetCandles to return DataFrame
func GetDataFrameInCandles(symbol string, limit int) *DataFrame {
	cframe := GetCandleFrame(symbol, limit)

	dframe := DataFrame{}
	dframe.CandleFrame = cframe

	return &dframe
}
