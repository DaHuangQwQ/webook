package main

import (
	"webook/pkg/logger"
)

func main() {
	logger.InitViper()
	app := InitApp()
	if err := app.WebServer.Start(); err != nil {
		panic(err)
	}
}
