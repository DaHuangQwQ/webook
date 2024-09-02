package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	initViper()
	server := InitWebServer()
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello world")
	})
	err := server.Run(":8090")
	if err != nil {
		return
	}
}

func initViper() {
	file := pflag.String("config", "config/dev.yaml", "指定文件路径")
	pflag.Parse()
	viper.SetConfigFile(*file)
	// 实时监听配置更新
	//viper.WatchConfig()
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	println(in.Name, in.Op)
	//})
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func initViperRemote() {
	viper.SetConfigType("yaml")
	err := viper.AddRemoteProvider("etcd3", "127.0.0.1:12379", "/webook")
	if err != nil {
		panic(err)
	}
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

func initLogger() {
	development, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(development)
	zap.L().Info("ni")
}
