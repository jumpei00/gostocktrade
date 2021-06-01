package server

import (
	"github.com/jumpei00/gostocktrade/stock"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/jumpei00/gostocktrade/app/models"
	"github.com/jumpei00/gostocktrade/config"
	"github.com/sirupsen/logrus"
)

var temp = template.Must(template.ParseFiles("templates/index.html"))

func indexAPIHandler(w http.ResponseWriter, req *http.Request) {
	temp.ExecuteTemplate(w, "index.html", nil)
}

func candleGetAPIHandler(w http.ResponseWriter, req *http.Request) {
	models.AllDeleteCandles()
	stockData, err := stock.GetStockData("VOO", 365)
	if err != nil {
		logrus.Warnf("stock get error: %v", err)
	}
	models.NewCandlesFromQuote(stockData).CreateCandles()

	cframe := models.GetCandles(365)

	js, err := json.Marshal(cframe)
	if err != nil {
		logrus.Warnf("candle json error: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func indicatorAPIHandler(w http.ResponseWriter, req *http.Request) {
	iframe := models.NewIndicator(365)
	iframe.AddEma(8)

	js, err := json.Marshal(iframe)
	if err != nil {
		logrus.Warnf("indicator json error: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// Run starts webserver
func Run() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", indexAPIHandler)
	http.HandleFunc("/candles", candleGetAPIHandler)
	http.HandleFunc("/indicator", indicatorAPIHandler)
	logrus.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil))
}
