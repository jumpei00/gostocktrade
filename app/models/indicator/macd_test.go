package indicator_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/jumpei00/gostocktrade/app/models/indicator"
	"testing"
)

func TestMacdBuyAndSell(t *testing.T) {
	assert := assert.New(t)

	signals := indicator.MacdSignals{}
	// when empty
	assert.False(signals.Sell("VOO", 0, 100))
	assert.True(signals.Buy("VOO", 0, 100))

	// when last is BUY
	assert.False(signals.Buy("VOO", 1, 100))
	assert.True(signals.Sell("VOO", 1, 100))

	// when last is SELL
	assert.False(signals.Sell("VOO", 2, 100))
	assert.True(signals.Buy("VOO", 2, 100))
}

func TestMacdProfit(t *testing.T) {
	assert := assert.New(t)

	signals := indicator.MacdSignals{
		MacdSignals: []indicator.MacdSignal{
			indicator.MacdSignal{
				Symbol: "VOO",
				Time: 0,
				Price: 100,
				Action: indicator.BUY,
			},
			indicator.MacdSignal{
				Symbol: "VOO",
				Time: 1,
				Price: 150,
				Action: indicator.SELL,
			},
		},
	}

	// when buy at 100, sell at 150,
	// expected profit is 50
	assert.Equal(50.0, signals.Profit())

	signals.MacdSignals = append(signals.MacdSignals, indicator.MacdSignal{
		Symbol: "VOO", Time: 2, Price: 100, Action: indicator.BUY,
	})
	
	// when buy at 100, sell at 150, buy at 100,
	// expected profit is 50
	assert.Equal(50.0, signals.Profit())
}