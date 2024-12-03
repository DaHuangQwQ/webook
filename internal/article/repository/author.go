package repository

import (
	"context"
	"github.com/DaHuangQwQ/webook/internal/user"

	//userv1 "github.com/DaHuangQwQ/webook/internal/api/proto/gen/user/v1"
	"github.com/DaHuangQwQ/webook/internal/article/domain"
	"github.com/DaHuangQwQ/webook/internal/article/repository/dao"
)

// AuthorRepository 封装user的client用于获取用户信息
type AuthorRepository interface {
	// FindAuthor id为文章id
	FindAuthor(ctx context.Context, id int64) (domain.Author, error)
}

type GrpcAuthorRepository struct {
	client user.App
	dao    dao.ArticleDao
}

func NewGrpcAuthorRepository(articleDao dao.ArticleDao, client user.App) AuthorRepository {
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
	u, err := g.client.Server.Profile(ctx, art.AuthorId)
	if err != nil {
		return domain.Author{}, err
	}
	return domain.Author{
		Id:   u.Id,
		Name: u.Nickname,
	}, nil
}
