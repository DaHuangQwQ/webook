package repository

import (
	"context"
	"webook/internal/domain"
)

type HistoryRecordRepository interface {
	AddRecord(ctx context.Context, record domain.HistoryRecord) error
}

type CachedHistoryRecordRepository struct {
}

func NewCachedHistoryRecordRepository() HistoryRecordRepository {
	return &CachedHistoryRecordRepository{}
}

func (c CachedHistoryRecordRepository) AddRecord(ctx context.Context, record domain.HistoryRecord) error {
	//TODO implement me
	panic("implement me")
}
