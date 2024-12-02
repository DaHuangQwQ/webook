package ioc

import (
	"github.com/DaHuangQwQ/webook/pkg/logger"
	"github.com/DaHuangQwQ/webook/user/service/oauth2/wechat"
	"github.com/spf13/viper"
)

func InitWechat(logger logger.LoggerV1) wechat.Service {
	type APP struct {
		ID     string `yaml:"ID"`
		SECRET string `yaml:"SECRET"`
	}
	var app APP
	err := viper.UnmarshalKey("APP", &app)
	if err != nil {
		panic("微信设置错误")
	}
	return wechat.NewService(app.ID, app.SECRET, logger)
}
