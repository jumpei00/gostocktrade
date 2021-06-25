package indicator

// WillrBacktestParam represents some parameters used for backtest
type WillrBacktestParam struct {
	WillrPeriodLow      int     `json:"period_low"`
	WillrPeriodHigh     int     `json:"period_high"`
	WillrBuyThreadLow   float64 `json:"buy_low"`
	WillrBuyThreadHigh  float64 `json:"buy_high"`
	WillrSellThreadLow  float64 `json:"sell_low"`
	WillrSellThreadHigh float64 `json:"sell_high"`
}

// WillrSignals stores WillrSignal
type WillrSignals struct {
	WillrSignals []WillrSignal
}

// WillrSignal is signal results of backtest
type WillrSignal struct {
	ID     int     `gorm:"primary_key" json:"-"`
	Symbol string  `json:"-"`
	Time   int64   `json:"time"`
	Price  float64 `json:"-"`
	Action string  `json:"action"`
}

// Buy appends buy-signal to Signals, if can not buy, return false
func (wi *WillrSignals) Buy(symbol string, time int64, price float64) bool {
	if !(wi.CanBuy()) {
		return false
	}
	wi.WillrSignals = append(wi.WillrSignals, WillrSignal{Symbol: symbol, Time: time, Price: price, Action: BUY})
	return true
}

// CanBuy judges whether buy or not
func (wi *WillrSignals) CanBuy() bool {
	lenSignals := len(wi.WillrSignals)
	// not buy or sell
	if lenSignals == 0 {
		return true
	}

	if wi.WillrSignals[lenSignals-1].Action == SELL {
		return true
	}

	return false
}

// Sell appends sell-signal to Signals, if can not sell, return false
func (wi *WillrSignals) Sell(symbol string, time int64, price float64) bool {
	if !(wi.CanSell()) {
		return false
	}
	wi.WillrSignals = append(wi.WillrSignals, WillrSignal{Symbol: symbol, Time: time, Price: price, Action: SELL})
	return true
}

// CanSell judges whether sell or not
func (wi *WillrSignals) CanSell() bool {
	lenSignals := len(wi.WillrSignals)
	// not buy or sell
	if lenSignals == 0 {
		return false
	}

	if wi.WillrSignals[lenSignals-1].Action == BUY {
		return true
	}

	return false
}

// Profit calculates profit for backtest
func (wi *WillrSignals) Profit() float64 {
	profit := 0.0
	afterSell := 0.0
	isHolding := false

	for _, signal := range wi.WillrSignals {
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
