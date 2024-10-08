package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	userv1 "webook/api/proto/gen/user/v1"
	"webook/user/domain"
	"webook/user/service"
)

type UserServiceServer struct {
	userv1.UnimplementedUserServiceServer
	svc service.UserService
}

func NewUserServiceServer(svc service.UserService) *UserServiceServer {
	return &UserServiceServer{svc: svc}
}

func (s *UserServiceServer) Register(grpcServer *grpc.Server) {
	userv1.RegisterUserServiceServer(grpcServer, s)
}

func (s *UserServiceServer) Signup(ctx context.Context, request *userv1.SignupRequest) (*userv1.SignupResponse, error) {
	err := s.svc.Signup(ctx, s.toDomain(request.User))
	return &userv1.SignupResponse{}, err
}

func (s *UserServiceServer) FindOrCreate(ctx context.Context, request *userv1.FindOrCreateRequest) (*userv1.FindOrCreateResponse, error) {
	user, err := s.svc.FindOrCreate(ctx, request.Phone)
	if err != nil {
		return nil, err
	}
	return &userv1.FindOrCreateResponse{User: s.toRpc(user)}, err
}

func (s *UserServiceServer) Login(ctx context.Context, request *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	user, err := s.svc.Login(ctx, request.Email, request.Password)
	return &userv1.LoginResponse{User: s.toRpc(user)}, err
}

func (s *UserServiceServer) Profile(ctx context.Context, request *userv1.ProfileRequest) (*userv1.ProfileResponse, error) {
	user, err := s.svc.Profile(ctx, request.Id)
	return &userv1.ProfileResponse{User: s.toRpc(user)}, err
}

func (s *UserServiceServer) UpdateNonSensitiveInfo(ctx context.Context, request *userv1.UpdateNonSensitiveInfoRequest) (*userv1.UpdateNonSensitiveInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *UserServiceServer) FindOrCreateByWechat(ctx context.Context, request *userv1.FindOrCreateByWechatRequest) (*userv1.FindOrCreateByWechatResponse, error) {
	user, err := s.svc.FindOrCreateByWechat(ctx, domain.WechatInfo{
		UnionId: request.Info.UnionId,
		OpenId:  request.Info.OpenId,
	})
	return &userv1.FindOrCreateByWechatResponse{User: s.toRpc(user)}, err
}

func (s *UserServiceServer) toDomain(user *userv1.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: user.Password,
		Nickname: user.Nickname,
		Grade:    0,
		Gender:   0,
		Avatar:   "",
		CTime:    user.Ctime.AsTime(),
		WechatInfo: domain.WechatInfo{
			UnionId: user.WechatInfo.UnionId,
			OpenId:  user.WechatInfo.OpenId,
		},
		Birthday:    user.Birthday.AsTime(),
		UserStatus:  0,
		DeptId:      0,
		Remark:      "",
		IsAdmin:     0,
		Address:     "",
		AboutMe:     user.AboutMe,
		LastLoginIp: "",
	}
}

func (s *UserServiceServer) toRpc(user domain.User) *userv1.User {
	return &userv1.User{
		Id:       user.Id,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: user.Password,
		Nickname: user.Nickname,
		WechatInfo: &userv1.WechatInfo{
			UnionId: user.WechatInfo.UnionId,
			OpenId:  user.WechatInfo.OpenId,
		},
		Birthday: timestamppb.New(user.Birthday),
		AboutMe:  user.AboutMe,
	}
}
