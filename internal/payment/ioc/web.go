package ioc

import (
	"github.com/DaHuangQwQ/webook/payment/web"
	"github.com/DaHuangQwQ/webook/pkg/ginx"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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
