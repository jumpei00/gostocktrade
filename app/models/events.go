package models

import "github.com/jumpei00/gostocktrade/app/models/indicator"

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
