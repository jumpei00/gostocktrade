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