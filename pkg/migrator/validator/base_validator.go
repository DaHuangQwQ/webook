package validator

import (
	"context"
	"gorm.io/gorm"
	"time"
	"webook/pkg/logger"
	events2 "webook/pkg/migrator/events"
)

type baseValidator struct {
	base   *gorm.DB
	target *gorm.DB

	// 这边需要告知，是以 SRC 为准，还是以 DST 为准
	// 修复数据需要知道
	direction string

	l        logger.LoggerV1
	producer events2.Producer
}

// 上报不一致的数据
func (v *baseValidator) notify(id int64, typ string) {
	// 这里我们要单独控制超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	evt := events2.InconsistentEvent{
		Direction: v.direction,
		ID:        id,
		Type:      typ,
	}

	err := v.producer.ProduceInconsistentEvent(ctx, evt)
	if err != nil {
		v.l.Error("发送消息失败", logger.Error(err),
			logger.Field{Key: "event", Val: evt})
	}
}
