package stock_test

import (
	"testing"

	"github.com/jumpei00/gostocktrade/stock"
	"github.com/stretchr/testify/assert"
)

func TestGetStockDataa(t *testing.T) {
	assert := assert.New(t)
	stock1, err1 := stock.GetStockData("VOO", 10, true)
	stock2, err2 := stock.GetStockData("TEST", 10, true)

	assert.Nil(err1)
	assert.Equal("VOO", stock1.Symbol)
	assert.NotEmpty(stock1.Date)
	// wrong symbol
	// err is nil, even if symbol is wrong
	assert.Nil(err2)
	assert.Len(stock2.Date, 0)
}
