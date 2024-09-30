package repository

import (
	"context"
	userv1 "webook/api/proto/gen/user/v1"
	"webook/article/domain"
	"webook/article/repository/dao"
)

// AuthorRepository 封装user的client用于获取用户信息
type AuthorRepository interface {
	// FindAuthor id为文章id
	FindAuthor(ctx context.Context, id int64) (domain.Author, error)
}

type GrpcAuthorRepository struct {
	client userv1.UserServiceClient
	dao    dao.ArticleDao
}

func NewGrpcAuthorRepository(articleDao dao.ArticleDao, client userv1.UserServiceClient) AuthorRepository {
	return &GrpcAuthorRepository{
		client: client,
		dao:    articleDao,
	}
}

func (g *GrpcAuthorRepository) FindAuthor(ctx context.Context, id int64) (domain.Author, error) {
	art, err := g.dao.FindById(ctx, id)
	if err != nil {
		return domain.Author{}, nil
	}
	u, err := g.client.Profile(ctx, &userv1.ProfileRequest{
		Id: art.AuthorId,
	})
	if err != nil {
		return domain.Author{}, err
	}
	return domain.Author{
		Id:   u.User.Id,
		Name: u.User.Nickname,
	}, nil
}
