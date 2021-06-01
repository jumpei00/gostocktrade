package server

import (
	"encoding/json"
	"github.com/jumpei00/gostocktrade/app/models"
	"github.com/jumpei00/gostocktrade/config"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)


func candlesAPIHandler(w http.ResponseWriter, req *http.Request) {
	// TODO
	df := models.GetCandles(50)

	js, err := json.Marshal(df)
	if err != nil {
		logrus.Warnf("json marshal error: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func helloWorld(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello World"))
}

// Run starts webserver
func Run() {
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/candles", candlesAPIHandler)
	logrus.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil))
}