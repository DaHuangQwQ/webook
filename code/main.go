package main

import "webook/pkg/logger"

func main() {
	logger.InitViper()
	app := initApp()
	err := app.server.Serve()
	if err != nil {
		panic(err)
	}
}
