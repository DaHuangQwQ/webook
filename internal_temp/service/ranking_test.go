package service

//func TestBatchRankingService_TopN(t *testing.T) {
//	const batchSize = 2
//	now := time.Now()
//	testCases := []struct {
//		name string
//
//		mock func(ctrl *gomock.Controller) (service.InteractiveService, ArticleService)
//
//		wantArts []domain.Article
//		wantErr  error
//	}{
//		{
//			name: "成功获取",
//			mock: func(ctrl *gomock.Controller) (service.InteractiveService, ArticleService) {
//				intrSvc := svcmocks.NewMockInteractiveService(ctrl)
//				artSvc := svcmocks.NewMockArticleService(ctrl)
//				// 先模拟批量获取数据
//				// 先模拟第一批
//				artSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), 0, 2).
//					Return([]domain.Article{
//						{Id: 1, UTime: now},
//						{Id: 2, UTime: now},
//					}, nil)
//				// 模拟第二批
//				artSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), 2, 2).
//					Return([]domain.Article{
//						{Id: 3, UTime: now},
//						{Id: 4, UTime: now},
//					}, nil)
//				// 模拟第三批
//				artSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), 4, 2).
//					// 没数据了
//					Return([]domain.Article{}, nil)
//
//				// 第一批的点赞数据
//				intrSvc.EXPECT().GetByIds(gomock.Any(), "article", []int64{1, 2}).
//					Return(map[int64]domain2.Interactive{
//						1: {LikeCnt: 1},
//						2: {LikeCnt: 2},
//					}, nil)
//
//				// 第二批的点赞数据
//				intrSvc.EXPECT().GetByIds(gomock.Any(), "article", []int64{3, 4}).
//					Return(map[int64]domain2.Interactive{
//						3: {LikeCnt: 3},
//						4: {LikeCnt: 4},
//					}, nil)
//
//				// 第三批的点赞数据
//				intrSvc.EXPECT().GetByIds(gomock.Any(), "article", []int64{}).
//					Return(map[int64]domain2.Interactive{}, nil)
//
//				return intrSvc, artSvc
//			},
//
//			wantErr: nil,
//			wantArts: []domain.Article{
//				{Id: 4, UTime: now},
//				{Id: 3, UTime: now},
//				{Id: 2, UTime: now},
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//			intrSvc, artSvc := tc.mock(ctrl)
//			svc := &BatchRankingService{
//				intrSvc:   intrSvc,
//				artSvc:    artSvc,
//				batchSize: batchSize,
//				n:         3,
//				scoreFunc: func(UTime time.Time, likeCnt int64) float64 {
//					return float64(likeCnt)
//				},
//			}
//			arts, err := svc.topN(context.Background())
//			assert.Equal(t, tc.wantErr, err)
//			assert.Equal(t, tc.wantArts, arts)
//		})
//	}
//}
