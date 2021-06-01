package models

import (
	"github.com/markcheno/go-quote"
	"gorm.io/gorm"
)

// Candles is slice of Candle
// Using this, create candle data in database
type Candles []Candle

// NewCandlesFromQuote converts Quote to slice of Candle due to creating in database,
// ex) [Date[1, 2, 3...], Open[1, 2, 3...]...] → [[Date[1], Open[1]...], [Date[2], Open[2]...]...]
// and return pointer of Candles(used as constructor)
func NewCandlesFromQuote(stock *quote.Quote) *Candles {
	candles := Candles{}
	for i := 0; i < len(stock.Date); i++ {
		candles = append(candles, Candle{
			Date:   stock.Date[i],
			Open:   stock.Open[i],
			High:   stock.High[i],
			Low:    stock.Low[i],
			Close:  stock.Close[i],
			Volume: stock.Volume[i],
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

// GetCandles gets candle data for limit by descending
// After get data, return DataFrame stored in data
func GetCandles(limit int) *CandleFrame {
	var candles Candles
	DB.Order("date desc").Limit(limit).Order("date asc").Find(&candles)

	cframe := CandleFrame{}
	for _, candle := range candles {
		cframe.Candles = append(cframe.Candles, candle)
	}

	return &cframe
}
