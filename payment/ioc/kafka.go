package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"webook/payment/events"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	// 哈希类算法 保证 有序性
	// 新加分区要保证有序性， 让新加的分区阻塞一会，等到其他分区都消费完
	saramaCfg.Producer.Partitioner = sarama.NewConsistentCRCHashPartitioner
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

func InitProducer(client sarama.Client) events.Producer {
	res, err := events.NewSaramaProducer(client)
	if err != nil {
		panic(err)
	}
	return res
}
