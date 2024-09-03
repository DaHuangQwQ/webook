package article

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

type Producer interface {
	ProduceReadEvent(ctx context.Context, readEvent ReadEvent) error
}

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(producer sarama.SyncProducer) Producer {
	return &KafkaProducer{
		producer: producer,
	}
}

func (k *KafkaProducer) ProduceReadEvent(ctx context.Context, readEvent ReadEvent) error {
	data, err := json.Marshal(readEvent)
	if err != nil {
		return err
	}
	_, _, err = k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "read_article",
		Value: sarama.ByteEncoder(data),
	})
	return err
}

type ReadEvent struct {
	Uid int64
	Aid int64
}
