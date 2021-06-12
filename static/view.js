// viewRealTime views current datetime on a window
export function viewRealTime() {
    const realtimeID = document.querySelector("#realtime");
    setInterval(() => {
        realtimeID.innerHTML = new Date();
    }, 100);
}

// viewChart views HighstockChart, also time parameter is converted due to highchart specification
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

    Highcharts.stockChart("container", {
        rangeSelector: {
            selected: 5,
        },

        title: {
            text: `Ticker: ${symbol}`
        },

        yAxis: [
            { height: "60%" },
            { top: "60%", height: "35%", offset: 0 }
        ],

        series: [{
            type: "candlestick",
            id: `${symbol} chart`,
            name: `${symbol} Stock Price`,
            data: candles,
        },
        {
            type: "column",
            id: `${symbol} volume`,
            name: `${symbol} Volume`,
            data: volume,
            yAxis: 1
        }]
    })
}

export function viewBacktestResults(tag, results) {
    tag.innerHTML = "";
    const time = new Date(results.timestamp)

    tag.innerHTML = `
        <p>Symbol: ${results.symbol} Latest Time: ${time.toString()}</p>
        <input type="checkbox" id="ema">
        [EMA] Performance: ${results.ema_performance} Short: ${results.ema_short} Long: ${results.ema_long}
        <input type="checkbox" id="bb">
        [BB] Performance: ${results.bb_performance} N: ${results.bb_n} K: ${results.bb_k}
        <input type="checkbox" id="bb">
        [MACD] Performance: ${results.macd_performance} Fast: ${results.macd_fast} Slow: ${results.macd_slow} Signal: ${results.macd_signal}
        <input type="checkbox" id="bb">
        [RSI] Performance: ${results.rsi_performance} Period: ${results.rsi_period} Buy: ${results.rsi_buythread} Sell: ${results.rsi_sellthread}
        <input type="checkbox" id="bb">
        [Willr] Performance: ${results.willr_performance} Period: ${results.willr_period} Buy: ${results.willr_buythread} Sell: ${results.willr_sellthread}
    `
}