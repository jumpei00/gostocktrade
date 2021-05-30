package main

import (
	"github.com/jumpei00/gostocktrade/app/models"
	"github.com/jumpei00/gostocktrade/log"
	"github.com/jumpei00/gostocktrade/stock"
)

func main() {
	log.SetLogging()

	q := stock.GetStockData("VOO", 365)
	candles := models.NewCandles(q)
	candles.Create()
}
