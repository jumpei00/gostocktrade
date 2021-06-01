package models

// Sma is sma indicator
type Sma struct {
	Period int       `json:"period,omitempty"`
	Values []float64 `json:"values,omitempty"`
}

// Ema is ema indicator
type Ema struct {
	Period int       `json:"period,omitempty"`
	Values []float64 `json:"values,omitempty"`
}

// BBands is bollinger bands indicator
type BBands struct {
	N   int       `json:"n,omitempty"`
	K   float64   `json:"k,omitempty"`
	Up  []float64 `json:"up,omitempty"`
	Mid []float64 `json:"mid,omitempty"`
	Low []float64 `json:"low,omitempty"`
}

// Macd is macd indicator
type Macd struct {
	Fast       int       `json:"fast_period,omitempty"`
	Slow       int       `json:"slow_period,omitempty"`
	Signal     int       `json:"signal_period,omitempty"`
	Macd       []float64 `json:"macd,omitempty"`
	MacdSignal []float64 `json:"macd_signal,omitempty"`
	MacdHist   []float64 `json:"macd_hist,omitempty"`
}

// Rsi is rsi indicator
type Rsi struct {
	Period int       `json:"period,omitempty"`
	Values []float64 `json:"values,omitempty"`
}

// WillR is williamR indicator
type WillR struct {
	Period int       `json:"period,omitempty"`
	Values []float64 `json:"values,omitempty"`
}
