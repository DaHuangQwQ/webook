package service

import "context"

type RankingService interface {
	TopN(ctx context.Context) error
}
