import { viewRealTime, viewChart } from "./view.js"
import { serverRequest } from "./request.js"


const body = document.querySelector("body");

// candleGet gets candle data from server, view graph
function candlesGet() {
    const symbol = body.querySelector("#symbol").value;
    const period = body.querySelector("#period").value;
    const query = new URLSearchParams({ symbol: symbol, period: period });

    serverRequest("/candles", query).then(function (json) {
        viewChart(symbol, json["candles"]);
    })
}

// getButtonAction is executed when GET button is pushed
function getButtonAction() {
    const getButton = body.querySelector("#candlesGet");
    getButton.addEventListener("click", () => {
        candlesGet();
    })
}

window.addEventListener("load", () => {
    viewRealTime();
    candlesGet();
    getButtonAction();
}, false)