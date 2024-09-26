package ioc

import (
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"webook/interactive/repository/dao"
	"webook/pkg/ginx"
	"webook/pkg/gormx/connpool"
	"webook/pkg/logger"
	"webook/pkg/migrator/events"
	"webook/pkg/migrator/events/fixer"
	"webook/pkg/migrator/scheduler"
)

const topic = "migrator_interactives"

func InitFixDataConsumer(l logger.LoggerV1,
	src SrcDB,
	dst DstDB,
	client sarama.Client,
) *fixer.Consumer[dao.Interactive] {
	res, err := fixer.NewConsumer[dao.Interactive](client, l, topic, src, dst)
	if err != nil {
		panic(err)
	}
	return res
}

func InitMigradatorProducer(p sarama.SyncProducer) events.Producer {
	return events.NewSaramaProducer(topic, p)
}

func InitMigratorWeb(
	l logger.LoggerV1,
	src SrcDB,
	dst DstDB,
	pool *connpool.DoubleWritePool,
	producer events.Producer,
) *ginx.Server {
	web := gin.Default()
	ginx.InitCounter(prometheus2.CounterOpts{
		Namespace: "DaHuang",
		Subsystem: "webook_intr",
		Name:      "http_biz_code",
		Help:      "GIN 中 HTTP 请求",
		ConstLabels: map[string]string{
			"instance_id": "my-instance-1",
		},
	})
	intrs := scheduler.NewScheduler[dao.Interactive](l, src, dst, pool, producer)
	intrs.RegisterRoutes(web.Group("/intr"))
	addr := viper.GetString("migrator.http.addr")
	// 你在这里加别的
	return &ginx.Server{
		Engine: web,
		Addr:   addr,
	}
}

func InitDoubleWritePool(src SrcDB, dst DstDB, l logger.LoggerV1) *connpool.DoubleWritePool {
	return connpool.NewDoubleWritePool(src.ConnPool, dst.ConnPool, l)
}
