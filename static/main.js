import { viewRealTime, viewChart, viewBacktestResults, viewTrade, viewSignal, removeSignal } from "./view.js"
import { candleGetRequest, backtestRequest, signalRequest, mappingParams } from "./request.js"

const candle = document.querySelector("#candle")
const backtest = document.querySelector("#backtest")

// getButtonAction is executed when GET button is pushed
function getButtonAction() {
    const getButton = candle.querySelector("#get");
    getButton.addEventListener("click", () => {
        candlesGet();
    })
}

// candleGet gets candle data from server, view graph
function candlesGet() {
    const symbol = candle.querySelector("#symbol").value;
    const period = candle.querySelector("#period").value;
    const query = new URLSearchParams({ symbol: symbol, period: period, get: true });

    candleGetRequest("/candles", query).then(function (json) {
        const result_tag = backtest.querySelector("#results");
        const trade_tag = backtest.querySelector("#trade");

        viewChart(symbol, json["candles"]);
        viewBacktestResults(result_tag, json["optimized_params"], signalButtonAction);
        viewTrade(trade_tag, json["trade"]);
    })
}

// testButtonAction is executed when TEST button is pushed
function testButtonAction() {
    const testButton = backtest.querySelector("#test");
    testButton.addEventListener("click", () => {
        executeBacktest();
    })
}

// executeBacktest executes backtest
function executeBacktest() {
    const params = backtest.querySelector("#params")
    let backtest_params = mappingParams(params)

    backtest_params.symbol = candle.querySelector("#symbol").value
    backtest_params.period = +backtest.querySelector("#period").value

    backtestRequest("/backtest", backtest_params).then(function (json) {
        const result_tag = backtest.querySelector("#results");
        const trade_tag = backtest.querySelector("#trade");

        viewBacktestResults(result_tag, json["optimized_params"], signalButtonAction);
        viewTrade(trade_tag, json["trade"]);
    })
}

// signalButtonAction is executed when checkbox state changes
function signalButtonAction(signal) {
    if (signal.checked) {
        signalGet(signal.value);
    } else {
        removeSignal(signal.value);
    }
}

// signalGet gets signals(BUY or SELL) for some indicators from server 
function signalGet(indicator) {
    const symbol = candle.querySelector("#symbol").value;
    const query = new URLSearchParams({ symbol: symbol, [indicator]: true })
    signalRequest("/candles", query).then(function (json) {
        viewSignal(symbol, indicator, json["signals"]);
    })
}

window.addEventListener("load", () => {
    viewRealTime();
    candlesGet();
    getButtonAction();
    testButtonAction();
}, false)