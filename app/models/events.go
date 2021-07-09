package models

import (
	"reflect"

	"github.com/sirupsen/logrus"

	"github.com/jumpei00/gostocktrade/app/models/indicator"
)

// Trade represents whether today is "buy" or "sell" or "no trade"
type Trade struct {
	LastEmaTrade   string `json:"last_ema"`
	IsEmaToday     bool   `json:"today_ema"`
	LastBBTrade    string `json:"last_bb"`
	IsBBToday      bool   `json:"today_bb"`
	LastMacdTrade  string `json:"last_macd"`
	IsMacdToday    bool   `json:"today_macd"`
	LastRsiTrade   string `json:"last_rsi"`
	IsRsiToday     bool   `json:"today_rsi"`
	LastWillrTrade string `json:"last_willr"`
	IsWillrToday   bool   `json:"today_willr"`
}

// GetTradeState returns Trade, after examining today trading,
// the symbol argument is certainly the same to the candle symbol
func GetTradeState(symbol string) *TradeFrame {
	signalEvents := GetSignalFrame(symbol, true, true, true, true, true).Signals
	lastCandleTime, err := LastCandleTime()
	if err != nil {
		logrus.Warnf("last candle get error: %v", err)
		return &TradeFrame{Trade: nil}
	}

	trade := Trade{
		LastEmaTrade:   indicator.NOTRADE,
		IsEmaToday:     false,
		LastBBTrade:    indicator.NOTRADE,
		IsBBToday:      false,
		LastMacdTrade:  indicator.NOTRADE,
		IsMacdToday:    false,
		LastRsiTrade:   indicator.NOTRADE,
		IsRsiToday:     false,
		LastWillrTrade: indicator.NOTRADE,
		IsWillrToday:   false,
	}

	for signal, time := range signalEvents.LastSignalTimes() {
		// no signal
		if time == 0 {
			continue
		}
		switch signal {
		case "emaTime":
			lastEma := signalEvents.EmaSignals[len(signalEvents.EmaSignals)-1]
			trade.LastEmaTrade = lastEma.Action
			trade.IsEmaToday = (lastEma.Time == lastCandleTime)
		case "bbTime":
			lastBB := signalEvents.BBSignals[len(signalEvents.BBSignals)-1]
			trade.LastBBTrade = lastBB.Action
			trade.IsBBToday = (lastBB.Time == lastCandleTime)
		case "macdTime":
			lastMacd := signalEvents.MacdSignals[len(signalEvents.MacdSignals)-1]
			trade.LastMacdTrade = lastMacd.Action
			trade.IsMacdToday = (lastMacd.Time == lastCandleTime)
		case "rsiTime":
			lastRsi := signalEvents.RsiSignals[len(signalEvents.RsiSignals)-1]
			trade.LastRsiTrade = lastRsi.Action
			trade.IsRsiToday = (lastRsi.Time == lastCandleTime)
		case "willrTime":
			lastWillr := signalEvents.WillrSignals[len(signalEvents.WillrSignals)-1]
			trade.LastWillrTrade = lastWillr.Action
			trade.IsWillrToday = (lastWillr.Time == lastCandleTime)
		}
	}

	return &TradeFrame{Trade: &trade}
}

// SignalEvents stores a part of signal
type SignalEvents struct {
	EmaSignals   []indicator.EmaSignal   `json:"ema_signals,omitempty"`
	BBSignals    []indicator.BBSignal    `json:"bb_signals,omitempty"`
	MacdSignals  []indicator.MacdSignal  `json:"macd_signals,omitempty"`
	RsiSignals   []indicator.RsiSignal   `json:"rsi_signals,omitempty"`
	WillrSignals []indicator.WillrSignal `json:"willr_signals,omitempty"`
}

// GetSignalFrame returns SignalFrame including a part of signal events
func GetSignalFrame(symbol string, ema, bb, macd, rsi, willr bool) *SignalFrame {
	if !(ema || bb || macd || rsi || willr) {
		return &SignalFrame{Signals: nil}
	}

	signalEvents := &SignalEvents{}

	if ema {
		emaSignals := []indicator.EmaSignal{}
		DB.Where("Symbol = ?", symbol).Find(&emaSignals)
		signalEvents.EmaSignals = emaSignals
	}

	if bb {
		bbSignals := []indicator.BBSignal{}
		DB.Where("Symbol = ?", symbol).Find(&bbSignals)
		signalEvents.BBSignals = bbSignals
	}

	if macd {
		macdSignals := []indicator.MacdSignal{}
		DB.Where("Symbol = ?", symbol).Find(&macdSignals)
		signalEvents.MacdSignals = macdSignals
	}

	if rsi {
		rsiSignals := []indicator.RsiSignal{}
		DB.Where("Symbol = ?", symbol).Find(&rsiSignals)
		signalEvents.RsiSignals = rsiSignals
	}

	if willr {
		willrSignals := []indicator.WillrSignal{}
		DB.Where("Symbol = ?", symbol).Find(&willrSignals)
		signalEvents.WillrSignals = willrSignals
	}

	return &SignalFrame{Signals: signalEvents}
}

// SignalTest execute backtest from last signal day, in other words, update each signal event
func SignalTest(symbol string, period int) bool {
	opParam := GetOptimizedParamFrame(symbol).Param
	if opParam == nil {
		return false
	}

	cframe := GetCandleFrame(symbol, period)
	signalEvents := GetSignalFrame(symbol, true, true, true, true, true).Signals

	firstTime := cframe.Candles[0].ID
	for k, v := range signalEvents.LastSignalTimes() {
		machID, err := MatchTime(v)
		if err != nil {
			continue
		}

		startDay := machID - firstTime + 1
		switch k {
		case "emaTime":
			emaSignals := cframe.backtestEma(
				startDay, opParam.EmaShort, opParam.EmaLong, &signalEvents.EmaSignals[len(signalEvents.EmaSignals)-1]).EmaSignals
			DB.Model(opParam).Association("EmaSignals").Append(emaSignals)
		case "bbTime":
			bbSignals := cframe.backtestBB(
				startDay, opParam.BBn, opParam.BBk, &signalEvents.BBSignals[len(signalEvents.BBSignals)-1]).BBSignals
			DB.Model(opParam).Association("BBSignals").Append(bbSignals)
		case "macdTime":
			macdSignals := cframe.backtestMacd(
				startDay, opParam.MacdFast, opParam.MacdSlow, opParam.MacdSignal, &signalEvents.MacdSignals[len(signalEvents.MacdSignals)-1]).MacdSignals
			DB.Model(opParam).Association("MacdSignals").Append(macdSignals)
		case "rsiTime":
			rsiSignals := cframe.backtestRsi(
				startDay, opParam.RsiPeriod, opParam.RsiBuyThread, opParam.RsiSellThread, &signalEvents.RsiSignals[len(signalEvents.RsiSignals)-1]).RsiSignals
			DB.Model(opParam).Association("RsiSignals").Append(rsiSignals)
		case "willrTime":
			willrSignals := cframe.backtestWillr(
				startDay, opParam.WillrPeriod, opParam.WillrBuyThread, opParam.WillrSellThread, &signalEvents.WillrSignals[len(signalEvents.WillrSignals)-1]).WillrSignals
			DB.Model(opParam).Association("WillrSignals").Append(willrSignals)

		}
	}

	return true
}

// LastSignalTimes returns a slice including Time for a last element of Signals
func (sg *SignalEvents) LastSignalTimes() map[string]int64 {
	lastTimes := []int64{}

	rv := reflect.ValueOf(*sg)
	for i := 0; i < rv.NumField(); i++ {
		signals := rv.Field(i)
		if signals.Len() != 0 {
			lastTimes = append(lastTimes, signals.Index(signals.Len()-1).FieldByName("Time").Int())
		} else {
			lastTimes = append(lastTimes, 0)
		}
	}

	return map[string]int64{
		"emaTime":   lastTimes[0],
		"bbTime":    lastTimes[1],
		"macdTime":  lastTimes[2],
		"rsiTime":   lastTimes[3],
		"willrTime": lastTimes[4],
	}

}
