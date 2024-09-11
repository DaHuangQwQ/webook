package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"net/http"
	"os"
	"webook/internal/api"
	"webook/ioc"
	"webook/pkg/ginx"
)

func main() {
	initViper()
	initPrometheus()
	app := InitWebServer()
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	server := app.server
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello，启动成功了！")
	})
	initDoc()
	ginx.NewWarpLogger(ioc.InitLogger())
	_ = server.Run(":8090")
}

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			return
		}
	}()
}

func initDoc() {
	spec := api.OpenAPISpec{
		OpenAPI: "3.0.0",
		Info: api.Info{
			Title:       "Sample API",
			Description: "This is a sample API to demonstrate OpenAPI generation in Go",
			Version:     "1.0.0",
		},
		Paths: ginx.Paths,
	}
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		panic(err)
	}
	if _, err := os.Stat("./docs"); os.IsNotExist(err) {
		err = os.Mkdir("./docs", 0755)
		if err != nil {
			panic(err)
		}
	}
	file, err := os.Create("./docs/doc.json")
	if err != nil {
		// 如果打开（或创建）文件时发生错误，则 panic
		panic(err)
	}
	_, err = file.Write(data)
	if err != nil {
		// 如果写入文件时发生错误，则 panic
		panic(err)
	}
	defer file.Close() // 确保在函数结束时关闭文件
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
