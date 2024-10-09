package events

import (
	"context"
	"github.com/IBM/sarama"
	"webook/article/repository"
	"webook/article/repository/dao"
	"webook/pkg/canalx"
	"webook/pkg/logger"
	"webook/pkg/saramax"
)

type MySQLBinlogConsumer struct {
	client sarama.Client
	l      logger.LoggerV1
	repo   *repository.CachedArticleRepository
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
			saramax.NewHandler[canalx.Message[dao.PublishedArticle]](r.l, r.Consume))
		if err != nil {
			r.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (r *MySQLBinlogConsumer) Consume(msg *sarama.ConsumerMessage,
	val canalx.Message[dao.PublishedArticle]) error {
	// 因为共用了一个 topic，所以会有很多表的数据，不是自己的就不用管了
	if val.Table != "published_articles" {
		return nil
	}
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	for _, data := range val.Data {
		var err error
		switch data.Status {
		//case domain.ArticleStatusPublished.ToUInt8():
		//	err = r.repo.Cache().SetPub(ctx, r.repo.ToDomain(dao.Article(data)))
		//case domain.ArticleStatusPrivate.ToUint8():
		//	err = r.repo.Cache().DelPub(ctx, data.Id)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
