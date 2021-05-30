package models

import (
	"github.com/markcheno/go-quote"
)

// Candles is slice of Candle
// Using this, create candle data in database
type Candles []Candle

// NewCandles converts Quote to slice of Candle due to creating in database,
// ex) [Date[1, 2, 3...], Open[1, 2, 3...]...] → [[Date[1], Open[1]...], [Date[2], Open[2]...]...]
// and return pointer of Candles(used as constructor)
func NewCandles(stock *quote.Quote) *Candles {
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

// Create creates candle data
func (cds *Candles) Create() {
	DB.Create(cds)
}
