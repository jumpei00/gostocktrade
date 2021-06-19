package models

import (
	"reflect"

	"github.com/sirupsen/logrus"

	"github.com/jumpei00/gostocktrade/app/models/indicator"
)

// Trade represents whether today is "buy" or "sell" or "no trade"
type Trade struct {
	EmaTrade   string `json:"ema_trade"`
	BBTrade    string `json:"bb_trade"`
	MacdTrade  string `json:"macd_trade"`
	RsiTrade   string `json:"rsi_trade"`
	WillrTrade string `json:"willr_trade"`
}

// TodayTradeState returns Trade, after examining today trading
func TodayTradeState(symbol string) *Trade {
	signalEvents := GetSignalFrame(symbol, true, true, true, true, true).Signals
	lastCandleTime, err := LastCandleTime()
	if err != nil {
		logrus.Warnf("last candle get error: %v", err)
		return nil
	}

	trade := Trade{
		EmaTrade:   indicator.NOTRADE,
		BBTrade:    indicator.NOTRADE,
		MacdTrade:  indicator.NOTRADE,
		RsiTrade:   indicator.NOTRADE,
		WillrTrade: indicator.NOTRADE,
	}

	for k, v := range signalEvents.LastSignalTimes() {
		if lastCandleTime != v {
			continue
		}
		switch k {
		case "emaTime":
			trade.EmaTrade = signalEvents.EmaSignals[len(signalEvents.EmaSignals)-1].Action
		case "bbTime":
			trade.BBTrade = signalEvents.BBSignals[len(signalEvents.BBSignals)-1].Action
		case "macdTime":
			trade.MacdTrade = signalEvents.MacdSignals[len(signalEvents.MacdSignals)-1].Action
		case "rsiTime":
			trade.RsiTrade = signalEvents.RsiSignals[len(signalEvents.RsiSignals)-1].Action
		case "willrTime":
			trade.WillrTrade = signalEvents.WillrSignals[len(signalEvents.WillrSignals)-1].Action
		}
	}

	return &trade
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

// SignalTest execute backtest from last signal day
func SignalTest(symbol string, period int) {
	cframe := GetCandleFrame(symbol, period)
	opParam := GetOptimizedParamFrame(symbol).Param
	signalEvents := GetSignalFrame(symbol, true, true, true, true, true).Signals

	if opParam == nil || signalEvents == nil {
		return
	}

	firstTime := cframe.Candles[0].ID
	for k, v := range signalEvents.LastSignalTimes() {
		machID, err := MatchTime(v)
		if err != nil {
			continue
		}

		startDay := machID - firstTime + 1
		switch k {
		case "emaTime":
			emaSignals := cframe.backtestEma(startDay, opParam.EmaShort, opParam.EmaLong).EmaSignals
			DB.Model(opParam).Association("EmaSignals").Append(emaSignals)
		case "bbTime":
			bbSignals := cframe.backtestBB(startDay, opParam.BBn, opParam.BBk).BBSignals
			DB.Model(opParam).Association("BBSignals").Append(bbSignals)
		case "macdTime":
			macdSignals := cframe.backtestMacd(startDay, opParam.MacdFast, opParam.MacdSlow, opParam.MacdSignal).MacdSignals
			DB.Model(opParam).Association("MacdSignals").Append(macdSignals)
		case "rsiTime":
			rsiSignals := cframe.backtestRsi(startDay, opParam.RsiPeriod, opParam.RsiBuyThread, opParam.RsiSellThread).RsiSignals
			DB.Model(opParam).Association("RsiSignals").Append(rsiSignals)
		case "willrTime":
			willrSignals := cframe.backtestWillr(startDay, opParam.WillrPeriod, opParam.WillrBuyThread, opParam.WillrSellThread).WillrSignals
			DB.Model(opParam).Association("WillrSignals").Append(willrSignals)

		}
	}
}
