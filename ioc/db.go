package ioc

import (
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
	"gorm.io/plugin/prometheus"
	"time"
	"webook/internal/repository/dao"
	"webook/pkg/gormx"
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
	err = db.Use(prometheus.New(prometheus.Config{
		DBName:          "webook",
		RefreshInterval: 15,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"thread_running"},
			},
		},
	}))
	if err != nil {
		panic(err)
	}
	cb := gormx.NewCallbacks(prometheus2.SummaryOpts{
		Namespace: "DaHuang",
		Subsystem: "webook",
		Name:      "gorm_db",
		Help:      "统计 GORM 的数据库查询",
		ConstLabels: map[string]string{
			"instance_id": "my_instance",
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})
	err = db.Use(cb)
	if err != nil {
		panic(err)
	}
	err = db.Use(tracing.NewPlugin(tracing.WithoutMetrics(),
		tracing.WithDBName("webook")))
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
