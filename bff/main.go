package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	InitViper()
	app := InitApp()
	if err := app.WebServer.Start(); err != nil {
		panic(err)
	}
}

func InitViper() {
	file := pflag.String("config", "config/dev.yaml", "指定文件路径")
	pflag.Parse()
	viper.SetConfigFile(*file)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
