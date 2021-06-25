package indicator

// BBBacktestParam represents some parameters used for backtest
type BBBacktestParam struct {
	BBnLow  int     `json:"n_low"`
	BBnHigh int     `json:"n_high"`
	BBkLow  float64 `json:"k_low"`
	BBkHigh float64 `json:"k_high"`
}

// BBSignals stores EmaSignal
type BBSignals struct {
	BBSignals []BBSignal
}

// BBSignal is signal results of backtest
type BBSignal struct {
	ID     int     `gorm:"primary_key" json:"-"`
	Symbol string  `json:"-"`
	Time   int64   `json:"time"`
	Price  float64 `json:"-"`
	Action string  `json:"action"`
}

// Buy appends buy-signal to Signals, if can not buy, return false
func (bb *BBSignals) Buy(symbol string, time int64, price float64) bool {
	if !(bb.CanBuy()) {
		return false
	}
	bb.BBSignals = append(bb.BBSignals, BBSignal{Symbol: symbol, Time: time, Price: price, Action: BUY})
	return true
}

// CanBuy judges whether buy or not
func (bb *BBSignals) CanBuy() bool {
	lenSignals := len(bb.BBSignals)
	// not buy or sell
	if lenSignals == 0 {
		return true
	}

	if bb.BBSignals[lenSignals-1].Action == SELL {
		return true
	}

	return false
}

// Sell appends sell-signal to Signals, if can not sell, return false
func (bb *BBSignals) Sell(symbol string, time int64, price float64) bool {
	if !(bb.CanSell()) {
		return false
	}
	bb.BBSignals = append(bb.BBSignals, BBSignal{Symbol: symbol, Time: time, Price: price, Action: SELL})
	return true
}

// CanSell judges whether sell or not
func (bb *BBSignals) CanSell() bool {
	lenSignals := len(bb.BBSignals)
	// not buy or sell
	if lenSignals == 0 {
		return false
	}

	if bb.BBSignals[lenSignals-1].Action == BUY {
		return true
	}

	return false
}

// Profit calculates profit for backtest
func (bb *BBSignals) Profit() float64 {
	profit := 0.0
	afterSell := 0.0
	isHolding := false

	for _, signal := range bb.BBSignals {
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
