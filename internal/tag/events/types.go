package events

import (
	"github.com/DaHuangQwQ/gpkg/logger"
	"github.com/DaHuangQwQ/webook/internal/search/service"
	"github.com/IBM/sarama"
)

type SyncDataEvent struct {
	IndexName string
	DocID     string
	Data      string
}

type SyncDataEventConsumer struct {
	svc    service.SyncService
	client sarama.Client
	l      logger.LoggerV1
}
