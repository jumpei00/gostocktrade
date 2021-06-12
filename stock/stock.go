package stock

import (
	"github.com/sirupsen/logrus"
	"time"

	"github.com/markcheno/go-quote"
)

const timeFormat = "2006-01-02"

// GetStockData dawnloads daily stockdata for symbol(GOOGL, FB...etc) during today ~ before dayPeriod.
// dayPeriod must be day(1day, 30days...etc).
// If stock data is not dawnloaded due to bad symbol, output panic.
func GetStockData(symbol string, dayPeriod int) (*quote.Quote, error) {
	endDay := time.Now()
	startDay := endDay.AddDate(0, 0, -dayPeriod)

	logrus.Info("get stock data from yahoo")
	stock, err := quote.NewQuoteFromYahoo(
		symbol, startDay.Format(timeFormat), endDay.Format(timeFormat), quote.Daily, true)
	if err != nil {
		return nil, err
	}

	return &stock, nil
}
