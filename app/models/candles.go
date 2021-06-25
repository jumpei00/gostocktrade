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

// NewCandlesFromQuote converts Quote to slice of Candle due to creating in database,
// ex) [Date[1, 2, 3...], Open[1, 2, 3...]...] → [[Date[1], Open[1]...], [Date[2], Open[2]...]...]
// and return pointer of Candles(used as constructor)
// Because of using for frondend, this method also converts time to Unixtime
func NewCandlesFromQuote(adjStock *quote.Quote, Stock *quote.Quote) *Candles {
	candles := Candles{}
	for i := 0; i < len(Stock.Date); i++ {
		candles = append(candles, Candle{
			Time:   Stock.Date[i].Unix() * 1000,
			Open:   (math.Round(adjStock.Open[i]*100) / 100),
			High:   (math.Round(adjStock.High[i]*100) / 100),
			Low:    (math.Round(adjStock.Low[i]*100) / 100),
			Close:  (math.Round(Stock.Close[i]*100) / 100),
			Volume: (math.Round(adjStock.Volume[i]*100) / 100),
		})
	}

	return &candles
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

// AllDeleteCandles deletes all data of "candles" table
func AllDeleteCandles() {
	DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Candle{})
}

// CreateCandles creates candle data
func (cs *Candles) CreateCandles() {
	DB.Create(cs)
}

// Candle is daily stock candledata, also used as json
type Candle struct {
	ID     int     `json:"-"`
	Time   int64   `json:"time"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

// LastCandleTime returns a time of last candle
func LastCandleTime() (int64, error) {
	var candle Candle
	if err := DB.Last(&candle).Error; err != nil {
		return 0, err
	}
	return candle.Time, nil
}

// MatchTime returns ID mathed to Time field
func MatchTime(time int64) (int, error) {
	var candle Candle
	if err := DB.Where("Time = ?", time).First(&candle).Error; err != nil {
		return 0, err
	}
	return candle.ID, nil
}
