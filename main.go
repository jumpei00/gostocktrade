package main

import (
	"github.com/jumpei00/gostocktrade/app/models"
	"github.com/jumpei00/gostocktrade/app/server"
	"github.com/jumpei00/gostocktrade/config"
	"github.com/jumpei00/gostocktrade/log"
)

func main() {
	config.InitConfig()
	log.SetLogging()
	models.InitDB()
	server.Run()
}
