package grpc

import (
	"context"
	"google.golang.org/grpc"
	"webook/api/proto/gen/reward/v1"
	"webook/reward/domain"
	"webook/reward/service"
)

type RewardServiceServer struct {
	rewardv1.UnimplementedRewardServiceServer
	svc service.RewardService
}

func NewRewardServiceServer(svc service.RewardService) *RewardServiceServer {
	return &RewardServiceServer{svc: svc}
}

func (s *RewardServiceServer) Register(server *grpc.Server) {
	rewardv1.RegisterRewardServiceServer(server, s)
}

func (s *RewardServiceServer) PreReward(ctx context.Context, request *rewardv1.PreRewardRequest) (*rewardv1.PreRewardResponse, error) {
	reward, err := s.svc.PreReward(ctx, domain.Reward{
		Uid: request.GetUid(),
		Target: domain.Target{
			BizId:   request.GetBizId(),
			Biz:     request.GetBiz(),
			BizName: request.GetBizName(),
			Uid:     request.GetTargetUid(),
		},
		Amt: request.GetAmt(),
	})
	if err != nil {
		return nil, err
	}
	return &rewardv1.PreRewardResponse{
		CodeUrl: reward.URL,
		Rid:     reward.Rid,
	}, nil
}

func (s *RewardServiceServer) GetReward(ctx context.Context, request *rewardv1.GetRewardRequest) (*rewardv1.GetRewardResponse, error) {
	reward, err := s.svc.GetReward(ctx, request.GetRid(), request.GetUid())
	if err != nil {
		return nil, err
	}
	return &rewardv1.GetRewardResponse{
		Status: rewardv1.RewardStatus(reward.Status),
	}, nil
}
