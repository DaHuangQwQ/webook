package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"webook/payment/web"
	"webook/pkg/ginx"
)

func NewWebServer(wechatHdl *web.WechatHandler) *ginx.Server {
	server := gin.Default()
	wechatHdl.RegisterRoutes(server)
	addr := viper.GetString("web.http.addr")
	return &ginx.Server{
		Engine: server,
		Addr:   addr,
	}
}
