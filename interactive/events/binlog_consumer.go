package events

import (
	"context"
	"github.com/IBM/sarama"
	"gorm.io/gorm"
	"sync/atomic"
	"time"
	"webook/article/repository"
	"webook/pkg/canalx"
	"webook/pkg/logger"
	"webook/pkg/migrator"
	"webook/pkg/migrator/events"
	"webook/pkg/migrator/validator"
	"webook/pkg/saramax"
)

type MySQLBinlogConsumer[T migrator.Entity] struct {
	client   sarama.Client
	l        logger.LoggerV1
	table    string
	repo     *repository.CachedArticleRepository
	srcToDst *validator.CanalIncrValidator[T]
	dstToSrc *validator.CanalIncrValidator[T]
	dstFirst *atomic.Bool
}

func NewMySQLBinlogConsumer[T migrator.Entity](
	client sarama.Client,
	l logger.LoggerV1,
	table string,
	src *gorm.DB,
	dst *gorm.DB,
	p events.Producer,
	repo *repository.CachedArticleRepository) *MySQLBinlogConsumer[T] {
	srcToDst := validator.NewCanalIncrValidator[T](src, dst, "SRC", l, p)
	dstToSrc := validator.NewCanalIncrValidator[T](src, dst, "DST", l, p)
	return &MySQLBinlogConsumer[T]{
		client: client, l: l,
		dstFirst: &atomic.Bool{},
		srcToDst: srcToDst,
		dstToSrc: dstToSrc,
		table:    table,
		repo:     repo}
}

func (r *MySQLBinlogConsumer[T]) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("migrator_incr",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{"webook_binlog"},
			saramax.NewHandler[canalx.Message[T]](r.l, r.Consume))
		if er != nil {
			r.l.Error("退出了消费循环异常", logger.Error(er))
		}
	}()
	return err
}

func (r *MySQLBinlogConsumer[T]) Consume(msg *sarama.ConsumerMessage,
	val canalx.Message[T]) error {
	dstFirst := r.dstFirst.Load()
	var v *validator.CanalIncrValidator[T]
	// db:
	//  src:
	//    dsn: "root:root@tcp(localhost:13316)/webook"
	//  dst:
	//    dsn: "root:root@tcp(localhost:13316)/webook_intr"
	if dstFirst && val.Database == "webook_intr" {
		// 校验，用 dst 的来校验
		v = r.dstToSrc
	} else if !dstFirst && val.Database == "webook" {
		v = r.srcToDst
	}
	if v != nil {
		for _, data := range val.Data {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			err := v.Validate(ctx, data.ID())
			cancel()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *MySQLBinlogConsumer[T]) DstFirst() {
	r.dstFirst.Store(true)
}
