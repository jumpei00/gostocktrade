package main

import (
	"encoding/json"
	"fmt"
	"github.com/jumpei00/gostocktrade/app/models"
)

func main() {
	// log.SetLogging()
	// server.Run()
	// models.BackTest("VOO", 200).CreateBacktestResult()
	// models.DeleteBacktestResult("VOO")
	f := models.GetOptimizedParamFrame("VOO")
	js, _ := json.Marshal(f)
	fmt.Println(string(js))
}
