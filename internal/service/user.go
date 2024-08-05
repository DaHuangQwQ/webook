package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"webook/internal/domain"
	"webook/internal/repository"
)

type UserService interface {
	Signup(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email string, password string) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

var (
	ErrDuplicate             = repository.ErrDuplicate
	ErrInvalidUserOrPassword = errors.New("用户不存在或者密码不对")
)

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)

	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 检查密码对不对
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	u, err := svc.repo.FindByID(ctx, id)
	return u, err
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {

	// 快路径
	user, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		// err 为 nil 进入这里
		// err 未找到 进入这里
		return user, err
	}
	// 降级策略
	//if ctx.Value("降级") == "true" {
	//	return
	//}
	// 慢路径
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	if err != nil && err != ErrDuplicate {
		return domain.User{}, err
	}
	// 这里会遇到 主从延迟 的问题
	return svc.repo.FindByPhone(ctx, phone)
}
