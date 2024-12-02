package main

import (
	"github.com/DaHuangQwQ/webook/pkg/logger"
)

func main() {
	logger.InitViper()
	app := InitAPP()
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	go func() {
		err := app.webAdmin.Start()
		if err != nil {
			return
		}
	}()
	err := app.server.Serve()
	if err != nil {
		panic(err)
	}
}
