package indicator

// Macd is macd indicator
type Macd struct {
	Fast       int       `json:"fast_period,omitempty"`
	Slow       int       `json:"slow_period,omitempty"`
	Signal     int       `json:"signal_period,omitempty"`
	Macd       []float64 `json:"macd,omitempty"`
	MacdSignal []float64 `json:"macd_signal,omitempty"`
	MacdHist   []float64 `json:"macd_hist,omitempty"`
}

// MacdSignals stores EmaSignal
type MacdSignals struct {
	MacdSignals []MacdSignal
}

// MacdSignal is signal results of backtest
type MacdSignal struct {
	ID     int     `gorm:"primary_key" json:"-"`
	Symbol string  `json:"-"`
	Time   int64   `json:"time"`
	Price  float64 `json:"-"`
	Action string  `json:"action"`
}

// Buy appends buy-signal to Signals, if can not buy, return false
func (md *MacdSignals) Buy(symbol string, time int64, price float64) bool {
	if !(md.CanBuy()) {
		return false
	}
	md.MacdSignals = append(md.MacdSignals, MacdSignal{Symbol:symbol, Time: time, Price: price, Action: BUY})
	return true
}

// CanBuy judges whether buy or not
func (md *MacdSignals) CanBuy() bool {
	lenSignals := len(md.MacdSignals)
	// not buy or sell
	if lenSignals == 0 {
		return true
	}

	if md.MacdSignals[lenSignals-1].Action == SELL {
		return true
	}

	return false
}

// Sell appends sell-signal to Signals, if can not sell, return false
func (md *MacdSignals) Sell(symbol string, time int64, price float64) bool {
	if !(md.CanSell()) {
		return false
	}
	md.MacdSignals = append(md.MacdSignals, MacdSignal{Symbol:symbol, Time: time, Price: price, Action: SELL})
	return true
}

// CanSell judges whether sell or not
func (md *MacdSignals) CanSell() bool {
	lenSignals := len(md.MacdSignals)
	// not buy or sell
	if lenSignals == 0 {
		return false
	}

	if md.MacdSignals[lenSignals-1].Action == BUY {
		return true
	}

	return false
}

// Profit calculates profit for backtest
func (md *MacdSignals) Profit() float64 {
	profit := 0.0
	afterSell := 0.0
	isHolding := false

	for _, signal := range md.MacdSignals {
		if signal.Action == BUY {
			profit -= signal.Price
			isHolding = true
		} else if signal.Action == SELL {
			profit += signal.Price
			afterSell = profit
			isHolding = false
		}
	}

	if isHolding {
		return afterSell
	}

	return profit
}
