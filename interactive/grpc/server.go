package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"webook/api/proto/gen/interactive/v1"
	"webook/interactive/domain"
	"webook/interactive/service"
)

type InteractiveServiceServer struct {
	interactivev1.UnimplementedInteractiveServiceServer
	svc service.InteractiveService
}

func NewInteractiveServiceServer(svc service.InteractiveService) *InteractiveServiceServer {
	return &InteractiveServiceServer{svc: svc}
}

func (i *InteractiveServiceServer) Register(s *grpc.Server) {
	interactivev1.RegisterInteractiveServiceServer(s, i)
}

func (i *InteractiveServiceServer) IncrReadCnt(ctx context.Context, request *interactivev1.IncrReadCntRequest) (*interactivev1.IncrReadCntResponse, error) {
	err := i.svc.IncrReadCnt(ctx, request.GetBiz(), request.GetBizId())
	return &interactivev1.IncrReadCntResponse{}, err
}

func (i *InteractiveServiceServer) Like(ctx context.Context, request *interactivev1.LikeRequest) (*interactivev1.LikeResponse, error) {
	err := i.svc.Like(ctx, request.GetBiz(), request.GetBizId(), request.Uid)
	return &interactivev1.LikeResponse{}, err
}

func (i *InteractiveServiceServer) CancelLike(ctx context.Context, request *interactivev1.CancelLikeRequest) (*interactivev1.CancelLikeResponse, error) {
	if request.Uid < 0 {
		return &interactivev1.CancelLikeResponse{}, status.Error(codes.InvalidArgument, "Uid must be greater than zero")
	}
	err := i.svc.CancelLike(ctx, request.GetBiz(), request.GetBizId(), request.Uid)
	return &interactivev1.CancelLikeResponse{}, err
}

func (i *InteractiveServiceServer) Collect(ctx context.Context, request *interactivev1.CollectRequest) (*interactivev1.CollectResponse, error) {
	err := i.svc.Collect(ctx, request.GetBiz(), request.GetBizId(), request.Cid, request.Uid)
	return &interactivev1.CollectResponse{}, err
}

func (i *InteractiveServiceServer) Get(ctx context.Context, request *interactivev1.GetRequest) (*interactivev1.GetResponse, error) {
	intr, err := i.svc.Get(ctx, request.GetBiz(), request.GetBizId(), request.Uid)
	return &interactivev1.GetResponse{
		Intr: i.toDTO(intr),
	}, err
}

func (i *InteractiveServiceServer) GetByIds(ctx context.Context, request *interactivev1.GetByIdsRequest) (*interactivev1.GetByIdsResponse, error) {
	intrs, err := i.svc.GetByIds(ctx, request.GetBiz(), request.GetBizIds())
	m := make(map[int64]*interactivev1.Interactive, len(intrs))
	for k, v := range intrs {
		m[k] = i.toDTO(v)
	}
	return &interactivev1.GetByIdsResponse{
		Intrs: m,
	}, err
}

func (i *InteractiveServiceServer) toDTO(intr domain.Interactive) *interactivev1.Interactive {
	return &interactivev1.Interactive{
		BizId:      intr.BizId,
		ReadCnt:    intr.ReadCnt,
		LikeCnt:    intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		Liked:      intr.Liked,
		Collected:  intr.Collected,
	}
}
