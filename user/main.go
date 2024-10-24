package main

import (
	"github.com/spf13/viper"
	"webook/pkg/logger"
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
