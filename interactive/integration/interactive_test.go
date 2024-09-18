package integration

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
	intrv1 "webook/api/proto/gen/interactive/v1"
	"webook/interactive/grpc"
	"webook/interactive/integration/startup"
	"webook/interactive/repository/dao"
)

type InteractiveTestSuite struct {
	suite.Suite
	db     *gorm.DB
	rdb    redis.Cmdable
	server *grpc.InteractiveServiceServer
}

func (s *InteractiveTestSuite) SetupSuite() {
	s.db = startup.InitTestDB()
	s.rdb = startup.InitRedis()
	s.server = startup.InitGRPCServer()
}

func (s *InteractiveTestSuite) TearDownTest() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err := s.db.Exec("TRUNCATE TABLE `interactives`").Error
	assert.NoError(s.T(), err)
	err = s.db.Exec("TRUNCATE TABLE `user_like_bizs`").Error
	assert.NoError(s.T(), err)
	err = s.db.Exec("TRUNCATE TABLE `user_collection_bizs`").Error
	assert.NoError(s.T(), err)
	// 清空 Redis
	err = s.rdb.FlushDB(ctx).Err()
	assert.NoError(s.T(), err)
}

func (s *InteractiveTestSuite) TestIncrReadCnt() {
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		biz   string
		bizId int64

		wantErr  error
		wantResp *intrv1.IncrReadCntResponse
	}{
		{
			// DB 和缓存都有数据
			name: "增加成功,db和redis",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				err := s.db.Create(dao.Interactive{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    5,
					CTime:      6,
					UTime:      7,
				}).Error
				assert.NoError(t, err)
				err = s.rdb.HSet(ctx, "interactive:test:2",
					"read_cnt", 3).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data dao.Interactive
				err := s.db.Where("id = ?", 1).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.UTime > 7)
				data.UTime = 0
				assert.Equal(t, dao.Interactive{
					Id:    1,
					Biz:   "test",
					BizId: 2,
					// +1 之后
					ReadCnt:    4,
					CollectCnt: 4,
					LikeCnt:    5,
					CTime:      6,
				}, data)
				cnt, err := s.rdb.HGet(ctx, "interactive:test:2", "read_cnt").Int()
				assert.NoError(t, err)
				assert.Equal(t, 4, cnt)
				err = s.rdb.Del(ctx, "interactive:test:2").Err()
				assert.NoError(t, err)
			},
			biz:      "test",
			bizId:    2,
			wantResp: &intrv1.IncrReadCntResponse{},
		},
		{
			// DB 有数据，缓存没有数据
			name: "增加成功,db有",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				err := s.db.WithContext(ctx).Create(dao.Interactive{
					Id:         3,
					Biz:        "test",
					BizId:      3,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    5,
					CTime:      6,
					UTime:      7,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data dao.Interactive
				err := s.db.Where("id = ?", 3).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.UTime > 7)
				data.UTime = 0
				assert.Equal(t, dao.Interactive{
					Id:    3,
					Biz:   "test",
					BizId: 3,
					// +1 之后
					ReadCnt:    4,
					CollectCnt: 4,
					LikeCnt:    5,
					CTime:      6,
				}, data)
				cnt, err := s.rdb.Exists(ctx, "interactive:test:3").Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), cnt)
			},
			biz:      "test",
			bizId:    3,
			wantResp: &intrv1.IncrReadCntResponse{},
		},
		{
			name:   "增加成功-都没有",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data dao.Interactive
				err := s.db.Where("biz = ? AND biz_id = ?", "test", 4).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.UTime > 0)
				assert.True(t, data.CTime > 0)
				assert.True(t, data.Id > 0)
				data.Id = 0
				data.UTime = 0
				data.CTime = 0
				assert.Equal(t, dao.Interactive{
					Biz:     "test",
					BizId:   4,
					ReadCnt: 1,
				}, data)
				cnt, err := s.rdb.Exists(ctx, "interactive:test:4").Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), cnt)
			},
			biz:      "test",
			bizId:    4,
			wantResp: &intrv1.IncrReadCntResponse{},
		},
	}

	// 不同于 AsyncSms 服务，我们不需要 mock，所以创建一个就可以
	// 不需要每个测试都创建
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			resp, err := s.server.IncrReadCnt(context.Background(), &intrv1.IncrReadCntRequest{
				Biz: tc.biz, BizId: tc.bizId,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantResp, resp)
			tc.after(t)
		})
	}
}

func (s *InteractiveTestSuite) TestLike() {
	t := s.T()
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		biz   string
		bizId int64
		uid   int64

		wantErr  error
		wantResp *intrv1.LikeResponse
	}{
		{
			name: "点赞-DB和cache都有",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				err := s.db.Create(dao.Interactive{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    5,
					CTime:      6,
					UTime:      7,
				}).Error
				assert.NoError(t, err)
				err = s.rdb.HSet(ctx, "interactive:test:2",
					"like_cnt", 3).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data dao.Interactive
				err := s.db.Where("id = ?", 1).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.UTime > 7)
				data.UTime = 0
				assert.Equal(t, dao.Interactive{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    6,
					CTime:      6,
				}, data)

				var likeBiz dao.UserLikeBiz
				err = s.db.Where("biz = ? AND biz_id = ? AND uid = ?",
					"test", 2, 123).First(&likeBiz).Error
				assert.NoError(t, err)
				assert.True(t, likeBiz.Id > 0)
				assert.True(t, likeBiz.CTime > 0)
				assert.True(t, likeBiz.UTime > 0)
				likeBiz.Id = 0
				likeBiz.CTime = 0
				likeBiz.UTime = 0
				assert.Equal(t, dao.UserLikeBiz{
					Biz:    "test",
					BizId:  2,
					Uid:    123,
					Status: 1,
				}, likeBiz)

				cnt, err := s.rdb.HGet(ctx, "interactive:test:2", "like_cnt").Int()
				assert.NoError(t, err)
				assert.Equal(t, 4, cnt)
				err = s.rdb.Del(ctx, "interactive:test:2").Err()
				assert.NoError(t, err)
			},
			biz:      "test",
			bizId:    2,
			uid:      123,
			wantResp: &intrv1.LikeResponse{},
		},
		{
			name:   "点赞-都没有",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data dao.Interactive
				err := s.db.Where("biz = ? AND biz_id = ?", "test", 3).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.UTime > 0)
				assert.True(t, data.CTime > 0)
				assert.True(t, data.Id > 0)
				data.UTime = 0
				data.CTime = 0
				data.Id = 0
				assert.Equal(t, dao.Interactive{
					Biz:     "test",
					BizId:   3,
					LikeCnt: 1,
				}, data)

				var likeBiz dao.UserLikeBiz
				err = s.db.Where("biz = ? AND biz_id = ? AND uid = ?",
					"test", 3, 123).First(&likeBiz).Error
				assert.NoError(t, err)
				assert.True(t, likeBiz.Id > 0)
				assert.True(t, likeBiz.CTime > 0)
				assert.True(t, likeBiz.UTime > 0)
				likeBiz.Id = 0
				likeBiz.CTime = 0
				likeBiz.UTime = 0
				assert.Equal(t, dao.UserLikeBiz{
					Biz:    "test",
					BizId:  3,
					Uid:    123,
					Status: 1,
				}, likeBiz)

				cnt, err := s.rdb.Exists(ctx, "interactive:test:2").Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), cnt)
			},
			biz:      "test",
			bizId:    3,
			uid:      123,
			wantResp: &intrv1.LikeResponse{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			resp, err := s.server.Like(context.Background(), &intrv1.LikeRequest{
				Biz: tc.biz, BizId: tc.bizId, Uid: tc.uid,
			})
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResp, resp)
			tc.after(t)
		})
	}
}

func (s *InteractiveTestSuite) TestDislike() {
	t := s.T()
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		biz   string
		bizId int64
		uid   int64

		wantErr  error
		wantResp *intrv1.CancelLikeResponse
	}{
		{
			name: "取消点赞-DB和cache都有",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				err := s.db.Create(dao.Interactive{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    5,
					CTime:      6,
					UTime:      7,
				}).Error
				assert.NoError(t, err)
				err = s.db.Create(dao.UserLikeBiz{
					Id:     1,
					Biz:    "test",
					BizId:  2,
					Uid:    123,
					CTime:  6,
					UTime:  7,
					Status: 1,
				}).Error
				assert.NoError(t, err)
				err = s.rdb.HSet(ctx, "interactive:test:2",
					"like_cnt", 3).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var data dao.Interactive
				err := s.db.Where("id = ?", 1).First(&data).Error
				assert.NoError(t, err)
				assert.True(t, data.UTime > 7)
				data.UTime = 0
				assert.Equal(t, dao.Interactive{
					Id:         1,
					Biz:        "test",
					BizId:      2,
					ReadCnt:    3,
					CollectCnt: 4,
					LikeCnt:    4,
					CTime:      6,
				}, data)

				var likeBiz dao.UserLikeBiz
				err = s.db.Where("id = ?", 1).First(&likeBiz).Error
				assert.NoError(t, err)
				assert.True(t, likeBiz.UTime > 7)
				likeBiz.UTime = 0
				assert.Equal(t, dao.UserLikeBiz{
					Id:     1,
					Biz:    "test",
					BizId:  2,
					Uid:    123,
					CTime:  6,
					Status: 0,
				}, likeBiz)

				cnt, err := s.rdb.HGet(ctx, "interactive:test:2", "like_cnt").Int()
				assert.NoError(t, err)
				assert.Equal(t, 2, cnt)
				err = s.rdb.Del(ctx, "interactive:test:2").Err()
				assert.NoError(t, err)
			},
			biz:      "test",
			bizId:    2,
			uid:      123,
			wantResp: &intrv1.CancelLikeResponse{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			resp, err := s.server.CancelLike(context.Background(), &intrv1.CancelLikeRequest{
				Biz: tc.biz, BizId: tc.bizId, Uid: tc.uid,
			})
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResp, resp)
			tc.after(t)
		})
	}
}

func (s *InteractiveTestSuite) TestCollect() {
	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		bizId int64
		biz   string
		cid   int64
		uid   int64

		wantErr  error
		wantResp *intrv1.CollectResponse
	}{
		{
			name:   "收藏成功,db和缓存都没有",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				var intr dao.Interactive
				err := s.db.Where("biz = ? AND biz_id = ?", "test", 1).First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.CTime > 0)
				intr.CTime = 0
				assert.True(t, intr.UTime > 0)
				intr.UTime = 0
				assert.True(t, intr.Id > 0)
				intr.Id = 0
				assert.Equal(t, dao.Interactive{
					Biz:        "test",
					BizId:      1,
					CollectCnt: 1,
				}, intr)
				cnt, err := s.rdb.Exists(ctx, "interactive:test:1").Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), cnt)
				// 收藏记录
				var cbiz dao.UserCollectionBiz
				err = s.db.WithContext(ctx).
					Where("uid = ? AND biz = ? AND biz_id = ?", 1, "test", 1).
					First(&cbiz).Error
				assert.NoError(t, err)
				assert.True(t, cbiz.CTime > 0)
				cbiz.CTime = 0
				assert.True(t, cbiz.UTime > 0)
				cbiz.UTime = 0
				assert.True(t, cbiz.Id > 0)
				cbiz.Id = 0
				assert.Equal(t, dao.UserCollectionBiz{
					Biz:   "test",
					BizId: 1,
					Cid:   1,
					Uid:   1,
				}, cbiz)
			},
			bizId:    1,
			biz:      "test",
			cid:      1,
			uid:      1,
			wantResp: &intrv1.CollectResponse{},
		},
		{
			name: "收藏成功,db有缓存没有",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := s.db.WithContext(ctx).Create(&dao.Interactive{
					Biz:        "test",
					BizId:      2,
					CollectCnt: 10,
					CTime:      123,
					UTime:      234,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				var intr dao.Interactive
				err := s.db.WithContext(ctx).
					Where("biz = ? AND biz_id = ?", "test", 2).First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.CTime > 0)
				intr.CTime = 0
				assert.True(t, intr.UTime > 0)
				intr.UTime = 0
				assert.True(t, intr.Id > 0)
				intr.Id = 0
				assert.Equal(t, dao.Interactive{
					Biz:        "test",
					BizId:      2,
					CollectCnt: 11,
				}, intr)
				cnt, err := s.rdb.Exists(ctx, "interactive:test:2").Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), cnt)

				var cbiz dao.UserCollectionBiz
				err = s.db.WithContext(ctx).
					Where("uid = ? AND biz = ? AND biz_id = ?", 1, "test", 2).
					First(&cbiz).Error
				assert.NoError(t, err)
				assert.True(t, cbiz.CTime > 0)
				cbiz.CTime = 0
				assert.True(t, cbiz.UTime > 0)
				cbiz.UTime = 0
				assert.True(t, cbiz.Id > 0)
				cbiz.Id = 0
				assert.Equal(t, dao.UserCollectionBiz{
					Biz:   "test",
					BizId: 2,
					Cid:   1,
					Uid:   1,
				}, cbiz)
			},
			bizId:    2,
			biz:      "test",
			cid:      1,
			uid:      1,
			wantResp: &intrv1.CollectResponse{},
		},
		{
			name: "收藏成功,db和缓存都有",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := s.db.WithContext(ctx).Create(&dao.Interactive{
					Biz:        "test",
					BizId:      3,
					CollectCnt: 10,
					CTime:      123,
					UTime:      234,
				}).Error
				assert.NoError(t, err)
				err = s.rdb.HSet(ctx, "interactive:test:3", "collect_cnt", 10).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				var intr dao.Interactive
				err := s.db.WithContext(ctx).
					Where("biz = ? AND biz_id = ?", "test", 3).First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.CTime > 0)
				intr.CTime = 0
				assert.True(t, intr.UTime > 0)
				intr.UTime = 0
				assert.True(t, intr.Id > 0)
				intr.Id = 0
				assert.Equal(t, dao.Interactive{
					Biz:        "test",
					BizId:      3,
					CollectCnt: 11,
				}, intr)
				cnt, err := s.rdb.HGet(ctx, "interactive:test:3", "collect_cnt").Int()
				assert.NoError(t, err)
				assert.Equal(t, 11, cnt)

				var cbiz dao.UserCollectionBiz
				err = s.db.WithContext(ctx).
					Where("uid = ? AND biz = ? AND biz_id = ?", 1, "test", 3).
					First(&cbiz).Error
				assert.NoError(t, err)
				assert.True(t, cbiz.CTime > 0)
				cbiz.CTime = 0
				assert.True(t, cbiz.UTime > 0)
				cbiz.UTime = 0
				assert.True(t, cbiz.Id > 0)
				cbiz.Id = 0
				assert.Equal(t, dao.UserCollectionBiz{
					Biz:   "test",
					BizId: 3,
					Cid:   1,
					Uid:   1,
				}, cbiz)
			},
			bizId:    3,
			biz:      "test",
			cid:      1,
			uid:      1,
			wantResp: &intrv1.CollectResponse{},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			resp, err := s.server.Collect(context.Background(), &intrv1.CollectRequest{
				Biz: tc.biz, BizId: tc.bizId, Cid: tc.cid, Uid: tc.uid,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantResp, resp)
			tc.after(t)
		})
	}
}

func (s *InteractiveTestSuite) TestGet() {
	testCases := []struct {
		name string

		before func(t *testing.T)

		bizId int64
		biz   string
		uid   int64

		wantErr error
		wantRes *intrv1.GetResponse
	}{
		{
			name:  "全部取出来了-无缓存",
			biz:   "test",
			bizId: 12,
			uid:   123,
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				err := s.db.WithContext(ctx).Create(&dao.Interactive{
					Biz:        "test",
					BizId:      12,
					ReadCnt:    100,
					CollectCnt: 200,
					LikeCnt:    300,
					CTime:      123,
					UTime:      234,
				}).Error
				assert.NoError(t, err)
			},
			wantRes: &intrv1.GetResponse{
				Intr: &intrv1.Interactive{
					Biz:        "test",
					BizId:      12,
					ReadCnt:    100,
					CollectCnt: 200,
					LikeCnt:    300,
				},
			},
		},
		{
			name:  "全部取出来了-命中缓存-用户已点赞收藏",
			biz:   "test",
			bizId: 3,
			uid:   123,
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				err := s.db.WithContext(ctx).
					Create(&dao.UserCollectionBiz{
						Cid:   1,
						Biz:   "test",
						BizId: 3,
						Uid:   123,
						CTime: 123,
						UTime: 124,
					}).Error
				assert.NoError(t, err)
				err = s.db.WithContext(ctx).
					Create(&dao.UserLikeBiz{
						Biz:    "test",
						BizId:  3,
						Uid:    123,
						CTime:  123,
						UTime:  124,
						Status: 1,
					}).Error
				assert.NoError(t, err)
				err = s.rdb.HSet(ctx, "interactive:test:3",
					"read_cnt", 0, "collect_cnt", 1).Err()
				assert.NoError(t, err)
			},
			wantRes: &intrv1.GetResponse{
				Intr: &intrv1.Interactive{
					BizId:      3,
					CollectCnt: 1,
					Collected:  true,
					Liked:      true,
				},
			},
		},
	}
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			res, err := s.server.Get(context.Background(), &intrv1.GetRequest{
				Biz: tc.biz, BizId: tc.bizId, Uid: tc.uid,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func (s *InteractiveTestSuite) TestGetByIds() {
	preCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// 准备数据
	for i := 1; i < 5; i++ {
		i := int64(i)
		err := s.db.WithContext(preCtx).
			Create(&dao.Interactive{
				Id:         i,
				Biz:        "test",
				BizId:      i,
				ReadCnt:    i,
				CollectCnt: i + 1,
				LikeCnt:    i + 2,
			}).Error
		assert.NoError(s.T(), err)
	}

	testCases := []struct {
		name string

		before func(t *testing.T)
		biz    string
		ids    []int64

		wantErr error
		wantRes *intrv1.GetByIdsResponse
	}{
		{
			name: "查找成功",
			biz:  "test",
			ids:  []int64{1, 2},
			wantRes: &intrv1.GetByIdsResponse{
				Intrs: map[int64]*intrv1.Interactive{
					1: {
						Biz:        "test",
						BizId:      1,
						ReadCnt:    1,
						CollectCnt: 2,
						LikeCnt:    3,
					},
					2: {
						Biz:        "test",
						BizId:      2,
						ReadCnt:    2,
						CollectCnt: 3,
						LikeCnt:    4,
					},
				},
			},
		},
		{
			name: "没有对应的数据",
			biz:  "test",
			ids:  []int64{100, 200},
			wantRes: &intrv1.GetByIdsResponse{
				Intrs: map[int64]*intrv1.Interactive{},
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			res, err := s.server.GetByIds(context.Background(), &intrv1.GetByIdsRequest{
				Biz: tc.biz, BizIds: tc.ids,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestInteractiveService(t *testing.T) {
	suite.Run(t, &InteractiveTestSuite{})
}
