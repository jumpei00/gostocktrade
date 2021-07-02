package indicator_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/jumpei00/gostocktrade/app/models/indicator"
	"testing"
)

func TestRsiBuyAndSell(t *testing.T) {
	assert := assert.New(t)

	signals := indicator.RsiSignals{}
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

func TestProfit(t *testing.T) {
	assert := assert.New(t)

	signals := indicator.RsiSignals{
		RsiSignals: []indicator.RsiSignal{
			indicator.RsiSignal{
				Symbol: "VOO",
				Time: 0,
				Price: 100,
				Action: indicator.BUY,
			},
			indicator.RsiSignal{
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

	signals.RsiSignals = append(signals.RsiSignals, indicator.RsiSignal{
		Symbol: "VOO", Time: 2, Price: 100, Action: indicator.BUY,
	})
	
	// when buy at 100, sell at 150, buy at 100,
	// expected profit is 50
	assert.Equal(50.0, signals.Profit())
}