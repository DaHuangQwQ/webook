package grpc

import (
	"context"
	"google.golang.org/grpc"
	"webook/api/proto/gen/code/v1"
	"webook/code/service"
)

type CodeServiceServer struct {
	codev1.UnimplementedCodeServiceServer
	svc service.CodeService
}

func NewCodeServiceServer(svc service.CodeService) *CodeServiceServer {
	return &CodeServiceServer{svc: svc}
}

func (s *CodeServiceServer) Register(server *grpc.Server) {
	codev1.RegisterCodeServiceServer(server, s)
}

func (s *CodeServiceServer) Send(ctx context.Context, request *codev1.CodeSendRequest) (*codev1.CodeSendResponse, error) {
	err := s.svc.Send(ctx, request.Biz, request.Phone)
	return &codev1.CodeSendResponse{}, err
}

func (s *CodeServiceServer) Verify(ctx context.Context, request *codev1.VerifyRequest) (*codev1.VerifyResponse, error) {
	verify, err := s.svc.Verify(ctx, request.Biz, request.Phone, request.InputCode)
	return &codev1.VerifyResponse{Answer: verify}, err
}
