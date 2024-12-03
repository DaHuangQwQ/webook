package main

import "github.com/DaHuangQwQ/gpkg/logger"

func main() {
	logger.InitViper()
	app := InitApp()
	go func() {
		err := app.web.Start()
		if err != nil {
			panic(err)
		}
	}()
	err := app.server.Serve()
	if err != nil {
		panic(err)
	}
}
