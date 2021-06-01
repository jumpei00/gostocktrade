package main

import (
	"github.com/jumpei00/gostocktrade/app/server"

	"github.com/jumpei00/gostocktrade/log"
)

func main() {
	log.SetLogging()
	server.Run()
}
