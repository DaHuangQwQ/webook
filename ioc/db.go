package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"time"
	"webook/internal/repository/dao"
	"webook/pkg/logger"
)

func InitDB(logger logger.LoggerV1) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var config = Config{
		DSN: "default",
	}
	err := viper.UnmarshalKey("db", &config)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(mysql.Open(config.DSN), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(logger.Debug), glogger.Config{
			// 慢查询阈值, 时间超过50毫秒才会输出日志
			// 一次磁盘io约10ms， mysql 命中索引
			SlowThreshold:        time.Millisecond * 50,
			ParameterizedQueries: true,
			LogLevel:             glogger.Info,
		}),
	})
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{Key: "args", Val: args})
}
