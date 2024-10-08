package grpc

import (
	"context"
	"google.golang.org/grpc"
	"webook/api/proto/gen/search/v1"
	"webook/search/service"
)

type SearchSyncServiceServer struct {
	searchv1.UnimplementedSearchServiceServer
	svc service.SyncService
}

func NewSearchSyncServiceServer(svc service.SyncService) *SearchSyncServiceServer {
	return &SearchSyncServiceServer{svc: svc}
}

func (s *SearchSyncServiceServer) Register(server *grpc.Server) {
	searchv1.RegisterSearchServiceServer(server, s)
}

func (s *SearchSyncServiceServer) InputUser(ctx context.Context, request *searchv1.InputUserRequest) (*searchv1.InputUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SearchSyncServiceServer) InputArticle(ctx context.Context, request *searchv1.InputArticleRequest) (*searchv1.InputArticleResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SearchSyncServiceServer) InputAny(ctx context.Context, request *searchv1.InputAnyRequest) (*searchv1.InputAnyResponse, error) {
	//TODO implement me
	panic("implement me")
}
