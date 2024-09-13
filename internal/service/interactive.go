package service

import (
	"context"
	"golang.org/x/sync/errgroup"
	"webook/internal/domain"
	"webook/internal/repository"
)

//go:generate mockgen -source=interactive.go -package=svcmocks -destination=mocks/interactive.mock.go InteractiveService
type InteractiveService interface {
	Read(ctx context.Context, biz string, bizId int64) error
	Like(ctx context.Context, biz string, id int64, uid int64) error
	CancelLike(ctx context.Context, biz string, id int64, uid int64) error
	// Collect cid 是收藏夹
	Collect(ctx context.Context, biz string, bizId int64, cid int64, uid int64) error
	Get(ctx context.Context, biz string, bizId int64, uid int64) (domain.Interactive, error)
	GetByIds(ctx context.Context, biz string, bizIds []int64) (map[int64]domain.Interactive, error)
}

type interactiveService struct {
	repo repository.InteractiveRepository
}

func NewInteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &interactiveService{
		repo: repo,
	}
}

func (svc *interactiveService) GetByIds(ctx context.Context, biz string, bizIds []int64) (map[int64]domain.Interactive, error) {
	intrs, err := svc.repo.GetByIds(ctx, biz, bizIds)
	if err != nil {
		return nil, err
	}
	res := make(map[int64]domain.Interactive, len(intrs))
	for _, intr := range intrs {
		res[intr.BizId] = intr
	}
	return res, nil
}

func (svc *interactiveService) Get(ctx context.Context, biz string, bizId int64, uid int64) (domain.Interactive, error) {
	var (
		res       domain.Interactive
		err       error
		Liked     bool
		Collected bool
	)
	var eg errgroup.Group
	eg.Go(func() error {
		res, err = svc.repo.Get(ctx, biz, bizId)
		return err
	})
	eg.Go(func() error {
		var er error
		Liked, er = svc.repo.Liked(ctx, biz, bizId, uid)
		return er
	})
	eg.Go(func() error {
		var er error
		Collected, er = svc.repo.Collected(ctx, biz, bizId, uid)
		return er
	})
	err = eg.Wait()
	if err != nil {
		return domain.Interactive{}, err
	}
	res.Collected = Collected
	res.Liked = Liked
	return res, nil
}

func (svc *interactiveService) Read(ctx context.Context, biz string, bizId int64) error {
	return svc.repo.IncrReadCnt(ctx, biz, bizId)
}

func (svc *interactiveService) Collect(ctx context.Context, biz string, bizId int64, cid int64, uid int64) error {
	return svc.repo.AddCollectionItem(ctx, biz, bizId, cid, uid)
}

func (svc *interactiveService) Like(ctx context.Context, biz string, id int64, uid int64) error {
	return svc.repo.IncrLike(ctx, biz, id, uid)
}

func (svc *interactiveService) CancelLike(ctx context.Context, biz string, id int64, uid int64) error {
	return svc.repo.DecrLike(ctx, biz, id, uid)
}
