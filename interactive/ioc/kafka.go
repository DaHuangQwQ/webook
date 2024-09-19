package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"webook/interactive/events"
	"webook/interactive/repository/dao"
	"webook/pkg/migrator/events/fixer"
	"webook/pkg/saramax"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	client, err := sarama.NewClient(cfg.Addrs, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitSyncProducer(client sarama.Client) sarama.SyncProducer {
	res, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return res
}

// NewConsumers 面临的问题依旧是所有的 Consumer 在这里注册一下
func NewConsumers(intr *events.InteractiveReadEventConsumer,
	fix *fixer.Consumer[dao.Interactive],
) []saramax.Consumer {
	return []saramax.Consumer{
		intr,
		fix,
	}
}
