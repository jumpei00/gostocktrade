package indicator

// EmaBacktestParam represents some parameters used for backtest
type EmaBacktestParam struct {
	EmaShortLow  int `json:"short_low"`
	EmaShortHigh int `json:"short_high"`
	EmaLongLow   int `json:"long_low"`
	EmaLongHigh  int `json:"long_high"`
}

// EmaSignals stores EmaSignal
type EmaSignals struct {
	EmaSignals []EmaSignal
}

// EmaSignal is signal results of backtest
type EmaSignal struct {
	ID     int     `gorm:"primary_key" json:"-"`
	Symbol string  `json:"-"`
	Time   int64   `json:"time"`
	Price  float64 `json:"-"`
	Action string  `json:"action"`
}

// Buy appends buy-signal to Signals, if can not buy, return false
func (ema *EmaSignals) Buy(symbol string, time int64, price float64) bool {
	if !(ema.CanBuy()) {
		return false
	}
	ema.EmaSignals = append(ema.EmaSignals, EmaSignal{Symbol: symbol, Time: time, Price: price, Action: BUY})
	return true
}

// CanBuy judges whether buy or not
func (ema *EmaSignals) CanBuy() bool {
	lenSignals := len(ema.EmaSignals)
	// not buy or sell
	if lenSignals == 0 {
		return true
	}

	if ema.EmaSignals[lenSignals-1].Action == SELL {
		return true
	}

	return false
}

// Sell appends sell-signal to Signals, if can not sell, return false
func (ema *EmaSignals) Sell(symbol string, time int64, price float64) bool {
	if !(ema.CanSell()) {
		return false
	}
	ema.EmaSignals = append(ema.EmaSignals, EmaSignal{Symbol: symbol, Time: time, Price: price, Action: SELL})
	return true
}

// CanSell judges whether sell or not
func (ema *EmaSignals) CanSell() bool {
	lenSignals := len(ema.EmaSignals)
	// not buy or sell
	if lenSignals == 0 {
		return false
	}

	if ema.EmaSignals[lenSignals-1].Action == BUY {
		return true
	}

	return false
}

// Profit calculates profit for backtest
func (ema *EmaSignals) Profit() float64 {
	profit := 0.0
	afterSell := 0.0
	isHolding := false

	for _, signal := range ema.EmaSignals {
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
