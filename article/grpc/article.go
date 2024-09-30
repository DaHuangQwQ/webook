package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	articlev1 "webook/api/proto/gen/article/v1"
	"webook/article/domain"
	"webook/article/service"
)

type ArticleServiceServer struct {
	articlev1.UnimplementedArticleServiceServer
	svc service.ArticleService
}

func NewArticleServiceServer(svc service.ArticleService) *ArticleServiceServer {
	return &ArticleServiceServer{svc: svc}
}

func (s *ArticleServiceServer) Register(server *grpc.Server) {
	articlev1.RegisterArticleServiceServer(server, s)
}

func (s *ArticleServiceServer) Save(ctx context.Context, request *articlev1.SaveRequest) (*articlev1.SaveResponse, error) {
	id, err := s.svc.Save(ctx, domain.Article{
		Id:      request.Article.Id,
		Title:   request.Article.Title,
		Content: request.Article.Content,
		Author: domain.Author{
			Id:   request.Article.Author.Id,
			Name: request.Article.Author.Name,
		},
		//ImgUrl: "",
		//Type:   request.Article.,
		Status: domain.ArticleStatus(request.Article.Status),
		UTime:  request.Article.Utime.AsTime(),
	})
	return &articlev1.SaveResponse{
		Id: id,
	}, err
}

func (s *ArticleServiceServer) Publish(ctx context.Context, request *articlev1.PublishRequest) (*articlev1.PublishResponse, error) {
	id, err := s.svc.Publish(ctx, domain.Article{
		Id:      request.Article.Id,
		Title:   request.Article.Title,
		Content: request.Article.Content,
		Author: domain.Author{
			Id:   request.Article.Author.Id,
			Name: request.Article.Author.Name,
		},
		//ImgUrl: "",
		//Type:   request.Article.,
		Status: domain.ArticleStatus(request.Article.Status),
		UTime:  request.Article.Utime.AsTime(),
	})
	return &articlev1.PublishResponse{
		Id: id,
	}, err
}

func (s *ArticleServiceServer) Withdraw(ctx context.Context, request *articlev1.WithdrawRequest) (*articlev1.WithdrawResponse, error) {
	err := s.svc.Withdraw(ctx, domain.Article{
		Id: request.Id,
		Author: domain.Author{
			Id: request.Uid,
		},
	})
	return &articlev1.WithdrawResponse{}, err
}

func (s *ArticleServiceServer) List(ctx context.Context, request *articlev1.ListRequest) (*articlev1.ListResponse, error) {
	panic("implement me")
}

func (s *ArticleServiceServer) GetById(ctx context.Context, request *articlev1.GetByIdRequest) (*articlev1.GetByIdResponse, error) {
	panic("implement me")
}

func (s *ArticleServiceServer) GetPublishedById(ctx context.Context, request *articlev1.GetPublishedByIdRequest) (*articlev1.GetPublishedByIdResponse, error) {
	art, err := s.svc.GetPublishedById(ctx, request.GetUid(), request.GetId())
	return &articlev1.GetPublishedByIdResponse{
		Article: &articlev1.Article{
			Id:      art.Id,
			Title:   art.Title,
			Status:  int32(art.Status),
			Content: art.Content,
			Author: &articlev1.Author{
				Id:   art.Author.Id,
				Name: art.Author.Name,
			},
			//Ctime: nil,
			Utime: timestamppb.New(art.UTime),
			//Abstract: "",
		},
	}, err
}

func (s *ArticleServiceServer) ListPub(ctx context.Context, request *articlev1.ListPubRequest) (*articlev1.ListPubResponse, error) {
	//TODO implement me
	panic("implement me")
}
