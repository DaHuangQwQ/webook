package ioc

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"gorm.io/plugin/prometheus"
	"webook/interactive/repository/dao"
	prometheus2 "webook/pkg/gormx/callbacks/prometheus"
	"webook/pkg/gormx/connpool"
)

type SrcDB *gorm.DB
type DstDB *gorm.DB

func InitSRC() SrcDB {
	return initDB("db.src", "webook")
}

func InitDST() DstDB {
	return initDB("db.dst", "webook_intr")
}

func InitBizDB(pool *connpool.DoubleWritePool) *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: pool,
	}))
	if err != nil {
		panic(err)
	}
	return db
}

func initDB(key string, name string) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	c := Config{
		DSN: "root:root@tcp(localhost:3307)/mysql",
	}
	err := viper.UnmarshalKey(key, &c)
	if err != nil {
		panic(fmt.Errorf("初始化配置失败 %v, 原因 %w", c, err))
	}
	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{
		// 使用 DEBUG 来打印
		//Logger: glogger.New(gormLoggerFunc(l.Debug),
		//	glogger.Config{
		//		SlowThreshold: 0,
		//		LogLevel:      glogger.Info,
		//	}),
	})
	if err != nil {
		panic(err)
	}

	// 接入 prometheus
	err = db.Use(prometheus.New(prometheus.Config{
		DBName: name,
		// 每 15 秒采集一些数据
		RefreshInterval: 15,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"Threads_running"},
			},
		}, // user defined metrics
	}))
	if err != nil {
		panic(err)
	}
	err = db.Use(tracing.NewPlugin(tracing.WithoutMetrics()))
	if err != nil {
		panic(err)
	}

	prom := prometheus2.Callbacks{
		Namespace:  "DaHuang",
		Subsystem:  "webook",
		Name:       "gorm_" + name,
		InstanceID: "my-instance-1",
		Help:       "gorm DB 查询",
	}
	err = prom.Register(db)
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
