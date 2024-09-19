package events

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

type Producer interface {
	ProduceInconsistentEvent(ctx context.Context, evt InconsistentEvent) error
}

type SaramaProducer struct {
	p     sarama.SyncProducer
	topic string
}

func NewSaramaProducer(topic string, p sarama.SyncProducer) *SaramaProducer {
	return &SaramaProducer{
		topic: topic,
		p:     p,
	}
}

func (s *SaramaProducer) ProduceInconsistentEvent(ctx context.Context, evt InconsistentEvent) error {
	val, _ := json.Marshal(evt)
	_, _, err := s.p.SendMessage(&sarama.ProducerMessage{
		Topic: s.topic,
		Value: sarama.ByteEncoder(val),
	})
	return err
}
