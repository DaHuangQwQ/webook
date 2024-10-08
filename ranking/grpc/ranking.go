package grpc

import (
	"context"
	"google.golang.org/grpc"
	"webook/api/proto/gen/ranking/v1"
	"webook/ranking/service"
)

type RankingServiceServer struct {
	rankingv1.UnimplementedRankingServiceServer
	svc service.RankingService
}

func NewRankingServiceServer(svc service.RankingService) *RankingServiceServer {
	return &RankingServiceServer{svc: svc}
}

func (s *RankingServiceServer) Register(server *grpc.Server) {
	rankingv1.RegisterRankingServiceServer(server, s)
}

func (s *RankingServiceServer) RankTopN(ctx context.Context, request *rankingv1.RankTopNRequest) (*rankingv1.RankTopNResponse, error) {
	err := s.svc.TopN(ctx)
	return &rankingv1.RankTopNResponse{}, err
}

func (s *RankingServiceServer) TopN(ctx context.Context, request *rankingv1.TopNRequest) (*rankingv1.TopNResponse, error) {
	//TODO implement me
	panic("implement me")
}
