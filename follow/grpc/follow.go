package grpc

import (
	"context"
	"google.golang.org/grpc"
	"webook/api/proto/gen/follow/v1"
	"webook/follow/service"
)

type FollowServiceServer struct {
	followv1.UnimplementedFollowServiceServer
	svc service.FollowRelationService
}

func NewFollowServiceServer(svc service.FollowRelationService) *FollowServiceServer {
	return &FollowServiceServer{svc: svc}
}

func (s *FollowServiceServer) Register(server *grpc.Server) {
	followv1.RegisterFollowServiceServer(server, s)
}

func (s *FollowServiceServer) Follow(ctx context.Context, request *followv1.FollowRequest) (*followv1.FollowResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *FollowServiceServer) CancelFollow(ctx context.Context, request *followv1.CancelFollowRequest) (*followv1.CancelFollowResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *FollowServiceServer) GetFollowee(ctx context.Context, request *followv1.GetFolloweeRequest) (*followv1.GetFolloweeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *FollowServiceServer) FollowInfo(ctx context.Context, request *followv1.FollowInfoRequest) (*followv1.FollowInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *FollowServiceServer) GetFollower(ctx context.Context, request *followv1.GetFollowerRequest) (*followv1.GetFollowerResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *FollowServiceServer) GetFollowStatic(ctx context.Context, request *followv1.GetFollowStaticRequest) (*followv1.GetFollowStaticResponse, error) {
	//TODO implement me
	panic("implement me")
}
