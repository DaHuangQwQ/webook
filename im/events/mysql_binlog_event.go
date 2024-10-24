package events

import (
	"context"
	"github.com/IBM/sarama"
	"strconv"
	"time"
	"webook/im/domain"
	"webook/im/service"
	"webook/pkg/canalx"
	"webook/pkg/logger"
	"webook/pkg/saramax"
)

type MySQLBinlogConsumer struct {
	client sarama.Client
	l      logger.LoggerV1
	svc    service.UserService
}

func (r *MySQLBinlogConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("pub_articles_cache",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{"webook_binlog"},
			saramax.NewHandler[canalx.Message[User]](r.l, r.Consume))
		if err != nil {
			r.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (r *MySQLBinlogConsumer) Consume(msg *sarama.ConsumerMessage,
	val canalx.Message[User]) error {
	// 因为共用了一个 topic，所以会有很多表的数据，不是自己的就不用管了
	// 只处理用户表的
	if val.Table != "users" {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for _, data := range val.Data {
		err := r.svc.Sync(ctx, domain.User{
			Nickname: data.Nickname,
			UserID:   strconv.FormatInt(data.Id, 10),
		})
		if err != nil {
			// 记录日志。
			continue
		}
	}
	return nil
}

type User struct {
	Id            int64
	Email         string
	Password      string
	Phone         string
	Birthday      string
	Nickname      string
	AboutMe       string
	WechatOpenId  string
	WechatUnionId string

	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
}
