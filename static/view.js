const chart = Highcharts.stockChart("container", {
    rangeSelector: {
        selected: 5,
    },

    title: {
        text: `Ticker: `
    },

    yAxis: [
        { height: "60%" },
        { top: "60%", height: "35%", offset: 0 }
    ],

    series: []
});

// viewRealTime views current datetime on a window
export function viewRealTime() {
    const realtimeID = document.querySelector("#realtime");
    setInterval(() => {
        realtimeID.innerHTML = new Date();
    }, 100);
}

// viewChart views HighstockChart
// also, brefore views, all current graphs are deleted
export function viewChart(symbol, json) {
    let candles = [];
    let volume = [];

    for (let candle of json) {
        candles.push([
            candle.time,
            candle.open,
            candle.high,
            candle.low,
            candle.close
        ])
        volume.push([
            candle.time,
            candle.volume
        ])
    }

    // all delete
    while (chart.series.length) {
        chart.series[0].remove();
    }

    chart.update({
        title: {
            text: `Ticker: ${symbol}`
        }
    })

    chart.addSeries(
        {
            type: "candlestick",
            id: `${symbol} chart`,
            name: `${symbol} Stock Price`,
            data: candles,
        }
    )

    chart.addSeries(
        {
            type: "column",
            id: `${symbol} volume`,
            name: `${symbol} Volume`,
            data: volume,
            yAxis: 1
        }
    )
}

// viewBacktestResults views optimized params,
// when optimized params exists, it's viewed, when no, it's cleared
export function viewBacktestResults(results_element, results, onchangeFunc) {
    results_element.innerHTML = "";

    // no data
    if (results == undefined) {
        return
    }

    const time = new Date(results.timestamp)

    results_element.innerHTML = `
        <p>Symbol: ${results.symbol} Latest Time: ${time.toString()}</p>
        <input type="checkbox" id="signal" value="ema">
        [EMA] Performance: ${results.ema_performance} Short: ${results.ema_short} Long: ${results.ema_long}
        <input type="checkbox" id="signal" value="bb">
        [BB] Performance: ${results.bb_performance} N: ${results.bb_n} K: ${results.bb_k}
        <input type="checkbox" id="signal" value="macd">
        [MACD] Performance: ${results.macd_performance} Fast: ${results.macd_fast} Slow: ${results.macd_slow} Signal: ${results.macd_signal}
        <input type="checkbox" id="signal" value="rsi">
        [RSI] Performance: ${results.rsi_performance} Period: ${results.rsi_period} Buy: ${results.rsi_buythread} Sell: ${results.rsi_sellthread}
        <input type="checkbox" id="signal" value="willr">
        [WILLr] Performance: ${results.willr_performance} Period: ${results.willr_period} Buy: ${results.willr_buythread} Sell: ${results.willr_sellthread}
    `

    // setting eventListener function for a part of signal
    const signals = results_element.querySelectorAll("#signal");
    for (let i = 0; i < signals.length; i++) {
        signals[i].addEventListener("change", () => {
            onchangeFunc(signals[i]);
        })
    }
}

export function viewTrade(trade_element, results) {
    trade_element.innerHTML = ""

    // no data
    if (results == undefined) {
        return
    }

    trade_element.innerHTML = `
        [EMA] <span style=${styleSet(results.last_ema, results.today_ema)}>${results.last_ema}</span>
        [BB] <span style=${styleSet(results.last_bb, results.today_bb)}>${results.last_bb}</span>
        [MACD] <span style=${styleSet(results.last_macd, results.today_macd)}>${results.last_macd}</span>
        [RSI] <span style=${styleSet(results.last_rsi, results.today_rsi)}>${results.last_rsi}</span>
        [WILLr] <span style=${styleSet(results.last_willr, results.today_willr)}>${results.last_willr}</span>
    `
}

function styleSet(signal, today_trade) {
    let style = ""

    switch (signal) {
        case "BUY":
            style = "color:red;"
            break
        case "SELL":
            style = "color:blue;"
            break
    }

    if (today_trade) {
        style += "font-weight:bold;"
    }
    
    return style
}

// viewSignal views signal(BUY or SELL) for some indicators, when checkbox is checked
export function viewSignal(symbol, signalName, signals) {
    let data = []
    for (let signal of signals[`${signalName}_signals`]) {
        data.push(
            {
                x: signal.time,
                title: signal.action,
                text: signalName
            }
        )
    }

    chart.addSeries(
        {
            type: "flags",
            onSeries: `${symbol} chart`,
            name: signalName,
            shape: "squarepin",
            width: 20,
            data: data
        }
    )
}

// viewSignal unviews signal(BUY or SELL) for some indicators, when checkbox is unchecked
export function removeSignal(signalName) {
    for (let i = 0; i < chart.series.length; i++) {
        if (chart.series[i].name == signalName) {
            chart.series[i].remove();
            return
        }
    }
}