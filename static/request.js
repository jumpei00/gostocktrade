const backtest_params = {
    symbol: "",
    period: "",
    ema: {
        short_low: "", short_high: "",
        long_low: "", long_high: "",
    },
    bb: {
        n_low: "", n_high: "",
        k_low: "", k_high: "",
    },
    macd: {
        fast_low: "", fast_high: "",
        slow_low: "", slow_high: "",
        signal_low: "", signal_high: "",
    },
    rsi: {
        period_low: "", period_high: "",
        buy_low: "", buy_high: "",
        sell_low: "", sell_high: "",
    },
    willr: {
        period_low: "", period_high: "",
        buy_low: "", buy_high: "",
        sell_low: "", sell_high: "",
    }
}

export function mappingParams(params) {
    backtest_params.ema.short_low = +params.querySelector("#ema_short_low").value
    backtest_params.ema.short_high = +params.querySelector("#ema_short_high").value
    backtest_params.ema.long_low = +params.querySelector("#ema_long_low").value
    backtest_params.ema.long_high = +params.querySelector("#ema_long_high").value

    backtest_params.bb.n_low = +params.querySelector("#bb_n_low").value
    backtest_params.bb.n_high = +params.querySelector("#bb_n_high").value
    backtest_params.bb.k_low = +params.querySelector("#bb_k_low").value
    backtest_params.bb.k_high = +params.querySelector("#bb_k_high").value

    backtest_params.macd.fast_low = +params.querySelector("#macd_fast_low").value
    backtest_params.macd.fast_high = +params.querySelector("#macd_fast_high").value
    backtest_params.macd.slow_low = +params.querySelector("#macd_slow_low").value
    backtest_params.macd.slow_high = +params.querySelector("#macd_slow_high").value
    backtest_params.macd.signal_low = +params.querySelector("#macd_signal_low").value
    backtest_params.macd.signal_high = +params.querySelector("#macd_signal_high").value

    backtest_params.rsi.period_low = +params.querySelector("#rsi_period_low").value
    backtest_params.rsi.period_high = +params.querySelector("#rsi_period_high").value
    backtest_params.rsi.buy_low = +params.querySelector("#rsi_buy_low").value
    backtest_params.rsi.buy_high = +params.querySelector("#rsi_buy_high").value
    backtest_params.rsi.sell_low = +params.querySelector("#rsi_sell_low").value
    backtest_params.rsi.sell_high = +params.querySelector("#rsi_sell_high").value

    backtest_params.willr.period_low = +params.querySelector("#willr_period_low").value
    backtest_params.willr.period_high = +params.querySelector("#willr_period_high").value
    backtest_params.willr.buy_low = -params.querySelector("#willr_buy_low").value
    backtest_params.willr.buy_high = -params.querySelector("#willr_buy_high").value
    backtest_params.willr.sell_low = -params.querySelector("#willr_sell_low").value
    backtest_params.willr.sell_high = -params.querySelector("#willr_sell_high").value

    return backtest_params
}

// candleGetRequest fetches any data from server, return json
// this method is only used to get candle data 
export async function candleGetRequest(uri, query) {
    let response = await fetch(uri + "?" + query)
    if (!response.ok) {
        throw new Error(`HTTP error, starus: ${response.status}`)
    }
    return response.json()
}

// backtestRequest is only used to execute backtest
export async function backtestRequest(uri, params) {
    let response = await fetch(uri, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(params)
    })

    if (!response.ok) {
        throw new Error(`HTTP error, status: ${response.status}`)
    }

    return response.json()
}