package service

import (
	"context"
	"github.com/DaHuangQwQ/gutil/slice"
	"github.com/ecodeclub/ekit/queue"
	"math"
	"time"
	"webook/internal/domain"
	"webook/internal/repository"
)

type RankingService interface {
	TopN(ctx context.Context) error
}

type BatchRankingService struct {
	repo      repository.RankingRepository
	artSvc    ArticleService
	intrSvc   InteractiveService
	batchSize int
	n         int
	scoreFunc func(t time.Time, likeCnt int64) float64
}

func NewBatchRankingService(artSvc ArticleService, intrSvc InteractiveService, repo repository.RankingRepository) RankingService {
	return &BatchRankingService{
		repo:      repo,
		artSvc:    artSvc,
		intrSvc:   intrSvc,
		batchSize: 100,
		n:         100,
		scoreFunc: func(t time.Time, likeCnt int64) float64 {
			dur := time.Since(t).Seconds()
			return float64(likeCnt-1) / math.Pow(dur+2, 1.5)
		},
	}
}

func (svc *BatchRankingService) TopN(ctx context.Context) error {
	arts, err := svc.topN(ctx)
	if err != nil {
		return err
	}
	// redis 缓存
	return svc.repo.ReplaceTopN(ctx, arts)
}

func (svc *BatchRankingService) topN(ctx context.Context) ([]domain.Article, error) {
	// 拿一批数据
	offset := 0
	start := time.Now()
	ddl := start.Add(-30 * 24 * time.Hour)

	type Score struct {
		art   domain.Article
		score float64
	}
	topN := queue.NewConcurrentPriorityQueue[Score](svc.n, func(src Score, dst Score) int {
		if src.score > dst.score {
			return 1
		} else if src.score == dst.score {
			return 0
		} else {
			return -1
		}
	})
	for {
		// 取数据
		arts, err := svc.artSvc.ListPub(ctx, start, offset, svc.batchSize)
		if err != nil {
			return nil, err
		}
		//if len(arts) == 0 {
		//	break
		//}
		ids := slice.Map(arts, func(idx int, art domain.Article) int64 {
			return art.Id
		})
		// 取点赞数
		intrMap, err := svc.intrSvc.GetByIds(ctx, "article", ids)
		if err != nil {
			return nil, err
		}
		for _, art := range arts {
			intr := intrMap[art.Id]
			score := svc.scoreFunc(art.UTime, intr.LikeCnt)
			ele := Score{
				score: score,
				art:   art,
			}
			err = topN.Enqueue(ele)
			if err == queue.ErrOutOfCapacity {
				// 这个也是满了
				// 拿出最小的元素
				minEle, _ := topN.Dequeue()
				if minEle.score < score {
					_ = topN.Enqueue(ele)
				} else {
					_ = topN.Enqueue(minEle)
				}
			}
		}
		offset = offset + len(arts)
		// 没有取够一批，我们就直接中断执行
		// 没有下一批了
		if len(arts) < svc.batchSize ||
			// 这个是一个优化
			arts[len(arts)-1].UTime.Before(ddl) {
			break
		}
	}

	res := make([]domain.Article, svc.n)
	for i := topN.Len() - 1; i >= 0; i-- {
		val, err := topN.Dequeue()
		if err != nil {
			break
		}
		res[i] = val.art
	}
	return res, nil
}
