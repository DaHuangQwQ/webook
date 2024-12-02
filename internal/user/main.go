package main

import (
	"github.com/DaHuangQwQ/webook/pkg/logger"
	"github.com/spf13/viper"
)

func main() {
	logger.InitViper()
	a := viper.AllKeys()
	println(a)
	app := initApp()
	err := app.server.Serve()
	if err != nil {
		panic(err)
	}
}
