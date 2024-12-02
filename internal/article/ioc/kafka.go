package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func InitSyncProducer() sarama.SyncProducer {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	scfg := sarama.NewConfig()
	scfg.Producer.Return.Successes = true
	client, err := sarama.NewClient(cfg.Addrs, scfg)
	if err != nil {
		panic(err)
	}
	p, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return p
}
