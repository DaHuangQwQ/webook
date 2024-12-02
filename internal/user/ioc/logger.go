package ioc

import (
	"github.com/DaHuangQwQ/webook/pkg/logger"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() logger.LoggerV1 {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "/var/log/user.log",
		MaxSize:    50,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   false,
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(lumberJackLogger),
		zapcore.DebugLevel,
	)
	l := zap.New(core)
	return logger.NewZapLogger(l)
}
