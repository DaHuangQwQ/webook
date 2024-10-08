package main

import "webook/pkg/logger"

func main() {
	logger.InitViper()
	app := InitApp()
	err := app.server.Serve()
	if err != nil {
		panic(err)
	}
}
