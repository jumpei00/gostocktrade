package indicator

// RsiBacktestParam represents some parameters used for backtest
type RsiBacktestParam struct {
	RsiPeriodLow      int     `json:"period_low"`
	RsiPeriodHigh     int     `json:"period_high"`
	RsiBuyThreadLow   float64 `json:"buy_low"`
	RsiBuyThreadHigh  float64 `json:"buy_high"`
	RsiSellThreadLow  float64 `json:"sell_low"`
	RsiSellThreadHigh float64 `json:"sell_high"`
}

// RsiSignals stores EmaSignal
type RsiSignals struct {
	RsiSignals []RsiSignal
}

// RsiSignal is signal results of backtest
type RsiSignal struct {
	ID     int     `gorm:"primary_key" json:"-"`
	Symbol string  `json:"-"`
	Time   int64   `json:"time"`
	Price  float64 `json:"-"`
	Action string  `json:"action"`
}

// Buy appends buy-signal to Signals, if can not buy, return false
func (rsi *RsiSignals) Buy(symbol string, time int64, price float64) bool {
	if !(rsi.CanBuy()) {
		return false
	}
	rsi.RsiSignals = append(rsi.RsiSignals, RsiSignal{Symbol: symbol, Time: time, Price: price, Action: BUY})
	return true
}

// CanBuy judges whether buy or not
func (rsi *RsiSignals) CanBuy() bool {
	lenSignals := len(rsi.RsiSignals)
	// not buy or sell
	if lenSignals == 0 {
		return true
	}

	if rsi.RsiSignals[lenSignals-1].Action == SELL {
		return true
	}

	return false
}

// Sell appends sell-signal to Signals, if can not sell, return false
func (rsi *RsiSignals) Sell(symbol string, time int64, price float64) bool {
	if !(rsi.CanSell()) {
		return false
	}
	rsi.RsiSignals = append(rsi.RsiSignals, RsiSignal{Symbol: symbol, Time: time, Price: price, Action: SELL})
	return true
}

// CanSell judges whether sell or not
func (rsi *RsiSignals) CanSell() bool {
	lenSignals := len(rsi.RsiSignals)
	// not buy or sell
	if lenSignals == 0 {
		return false
	}

	if rsi.RsiSignals[lenSignals-1].Action == BUY {
		return true
	}

	return false
}

// Profit calculates profit for backtest
func (rsi *RsiSignals) Profit() float64 {
	profit := 0.0
	afterSell := 0.0
	isHolding := false

	for _, signal := range rsi.RsiSignals {
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
