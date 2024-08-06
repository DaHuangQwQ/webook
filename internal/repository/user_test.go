package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	cachemocks "webook/internal/repository/cache/mocks"
	"webook/internal/repository/dao"
	daomocks "webook/internal/repository/dao/mocks"
)

func TestCachedUserRepository_FindById(t *testing.T) {
	nowMs := time.Now().UnixMilli()
	now := time.UnixMilli(nowMs)
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao)

		ctx context.Context
		uid int64

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "查找成功，缓存未命中",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				uid := int64(123)
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).
					Return(domain.User{}, cache.ErrUserKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), uid).
					Return(dao.User{
						Id: uid,
						Email: sql.NullString{
							String: "123@qq.com",
							Valid:  true,
						},
						Password: "123456",
						Phone: sql.NullString{
							String: "15212345678",
							Valid:  true,
						},
						CTime: nowMs,
						UTime: 102,
					}, nil)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "123456",
					Phone:    "15212345678",
					CTime:    now,
				}).Return(nil)
				return c, d
			},
			uid: 123,
			ctx: context.Background(),
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123456",
				Phone:    "15212345678",
				CTime:    now,
			},
			wantErr: nil,
		},

		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				uid := int64(123)
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).
					Return(domain.User{
						Id:       123,
						Email:    "123@qq.com",
						Password: "123456",
						Phone:    "15212345678",
						CTime:    time.UnixMilli(101),
					}, nil)
				return c, d
			},
			uid: 123,
			ctx: context.Background(),
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123456",
				Phone:    "15212345678",
				CTime:    time.UnixMilli(101),
			},
			wantErr: nil,
		},

		{
			name: "未找到用户",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				uid := int64(123)
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).
					Return(domain.User{}, cache.ErrUserKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), uid).
					Return(dao.User{}, dao.ErrRecordNotFound)
				return c, d
			},
			uid:      123,
			ctx:      context.Background(),
			wantUser: domain.User{},
			wantErr:  dao.ErrRecordNotFound,
		},

		{
			name: "回写缓存失败",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDao) {
				uid := int64(123)
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).
					Return(domain.User{}, cache.ErrUserKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), uid).
					Return(dao.User{
						Id: uid,
						Email: sql.NullString{
							String: "123@qq.com",
							Valid:  true,
						},
						Password: "123456",
						Phone: sql.NullString{
							String: "15212345678",
							Valid:  true,
						},
						CTime: 101,
						UTime: 102,
					}, nil)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "123456",
					Phone:    "15212345678",
					CTime:    time.UnixMilli(101),
				}).Return(errors.New("redis错误"))
				return c, d
			},
			uid: 123,
			ctx: context.Background(),
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123456",
				Phone:    "15212345678",
				CTime:    time.UnixMilli(101),
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc, ud := tc.mock(ctrl)
			svc := NewUserRepository(ud, uc)
			user, err := svc.FindByID(tc.ctx, tc.uid)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
			time.Sleep(time.Second)
		})
	}
}
