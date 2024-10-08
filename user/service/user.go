package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"mime/multipart"
	domain2 "webook/user/domain"
	"webook/user/repository"
)

type UserService interface {
	Signup(ctx context.Context, u domain2.User) error
	Login(ctx context.Context, email string, password string) (domain2.User, error)
	Profile(ctx context.Context, id int64) (domain2.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain2.User, error)
	FindOrCreateByWechat(ctx context.Context, wechatInfo domain2.WechatInfo) (domain2.User, error)
	UpdateByID(ctx context.Context, user domain2.User) error
	FindByID(ctx context.Context, id int64) (domain2.User, error)
	AvatarUpdate(ctx context.Context, id int64, file multipart.File, fileType string) (string, error)
	GetAvatar(ctx context.Context, id int64) (string, error)
}

type userService struct {
	repo repository.UserRepository
}

func (svc *userService) GetAvatar(ctx context.Context, id int64) (string, error) {
	return svc.repo.GetAvatar(ctx, id)
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

func (svc *userService) Signup(ctx context.Context, u domain2.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *userService) Login(ctx context.Context, email string, password string) (domain2.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)

	if errors.Is(err, repository.ErrUserNotFound) {
		return domain2.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain2.User{}, err
	}
	// 检查密码对不对
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain2.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *userService) Profile(ctx context.Context, id int64) (domain2.User, error) {
	u, err := svc.repo.FindByID(ctx, id)
	return u, err
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain2.User, error) {

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
	err = svc.repo.Create(ctx, domain2.User{
		Phone: phone,
	})
	if err != nil && err != ErrDuplicate {
		return domain2.User{}, err
	}
	// 这里会遇到 主从延迟 的问题
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *userService) FindOrCreateByWechat(ctx context.Context, wechatInfo domain2.WechatInfo) (domain2.User, error) {
	// 快路径
	user, err := svc.repo.FindByWechat(ctx, wechatInfo.OpenId)
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
	err = svc.repo.Create(ctx, domain2.User{
		WechatInfo: wechatInfo,
	})
	if err != nil && err != ErrDuplicate {
		return domain2.User{}, err
	}
	// 这里会遇到 主从延迟 的问题
	return svc.repo.FindByWechat(ctx, wechatInfo.OpenId)
}

func (svc *userService) UpdateByID(ctx context.Context, user domain2.User) error {
	return svc.repo.UpdateByID(ctx, user)
}

func (svc *userService) FindByID(ctx context.Context, id int64) (domain2.User, error) {
	return svc.repo.FindByID(ctx, id)
}

func (svc *userService) AvatarUpdate(ctx context.Context, id int64, file multipart.File, fileType string) (string, error) {
	fileBytes, err := svc.fileToBytes(file)
	if err != nil {
		return " ", fmt.Errorf("%w", err)
	}
	return svc.repo.AvatarUpdate(ctx, id, fileBytes, fileType)
}

func (svc *userService) fileToBytes(file multipart.File) ([]byte, error) {
	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, file); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
