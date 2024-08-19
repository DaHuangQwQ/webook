package ioc

import (
	"go.uber.org/zap"
)

func InitLogger() *zap.Logger {
	development, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return development
}
