import { viewRealTime, viewChart, viewBacktestResults } from "./view.js"
import { candleGetRequest, backtestRequest, mappingParams } from "./request.js"


const candle = document.querySelector("#candle")
const backtest = document.querySelector("#backtest")

// candleGet gets candle data from server, view graph
function candlesGet() {
    const symbol = candle.querySelector("#symbol").value;
    const period = candle.querySelector("#period").value;
    const query = new URLSearchParams({ symbol: symbol, period: period, get: true });

    candleGetRequest("/candles", query).then(function (json) {
        viewChart(symbol, json["candles"]);
    })
}

// getButtonAction is executed when GET button is pushed
function getButtonAction() {
    const getButton = candle.querySelector("#get");
    getButton.addEventListener("click", () => {
        candlesGet();
    })
}

function executeBacktest() {
    const params = backtest.querySelector("#params")
    let backtest_params = mappingParams(params)

    backtest_params.symbol = candle.querySelector("#symbol").value
    backtest_params.period = +backtest.querySelector("#period").value

    backtestRequest("/backtest", backtest_params).then(function (json) {
        const tag = backtest.querySelector("#results")
        viewBacktestResults(tag, json["optimized_params"]);
    })
}

function testButtonAction() {
    const testButton = backtest.querySelector("#test");
    testButton.addEventListener("click", () => {
        executeBacktest();
    })
}

window.addEventListener("load", () => {
    viewRealTime();
    candlesGet();
    getButtonAction();
    testButtonAction();
}, false)