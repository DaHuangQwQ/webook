package startup

import "github.com/DaHuangQwQ/webook/pkg/logger"

func InitLog() logger.LoggerV1 {
	return logger.NewNoOpLogger()
}
