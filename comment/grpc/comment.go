package grpc

import (
	"context"
	"google.golang.org/grpc"
	"webook/api/proto/gen/comment/v1"
	"webook/comment/service"
)

type CommentServiceServer struct {
	commentv1.UnimplementedCommentServiceServer
	svc service.CommentService
}

func NewCommentServiceServer(svc service.CommentService) *CommentServiceServer {
	return &CommentServiceServer{svc: svc}
}

func (s *CommentServiceServer) Register(server *grpc.Server) {
	commentv1.RegisterCommentServiceServer(server, s)
}

func (s *CommentServiceServer) GetCommentList(ctx context.Context, request *commentv1.CommentListRequest) (*commentv1.CommentListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *CommentServiceServer) DeleteComment(ctx context.Context, request *commentv1.DeleteCommentRequest) (*commentv1.DeleteCommentResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *CommentServiceServer) CreateComment(ctx context.Context, request *commentv1.CreateCommentRequest) (*commentv1.CreateCommentResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *CommentServiceServer) GetMoreReplies(ctx context.Context, request *commentv1.GetMoreRepliesRequest) (*commentv1.GetMoreRepliesResponse, error) {
	//TODO implement me
	panic("implement me")
}
