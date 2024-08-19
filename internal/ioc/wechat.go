package ioc

import (
	"webook/internal/service/oauth2/wechat"
	"webook/pkg/logger"
)

func InitWechat(logger logger.LoggerV1) wechat.Service {
	return wechat.NewService("123", "123", logger)
}
