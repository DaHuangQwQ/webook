package api

import (
	articlev1 "webook/api/proto/gen/article/v1"
)

type GetArticleListReq struct {
	Meta `path:"/list" method:"post"`
	ListReq
}

type GetArticlePubDetailReq struct {
	Meta `path:"/:id" method:"get"`
}

type RewardReq struct {
	Meta `path:"/reward" method:"post"`
	Id   int64 `json:"id"`
	Amt  int64 `json:"amt"`
}

type ArticleLikeReq struct {
	Meta `path:"/like" method:"post"`
	Id   int64 `json:"id"`
	Like bool  `json:"like"`
}

type CollectReq struct {
	Meta `path:"/collect" method:"post"`
	Id   int64 `json:"id"`
	Cid  int64 `json:"cid"`
}

type ArticleVo struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	// 摘要
	Abstract string `json:"abstract"`
	// 内容
	Content string `json:"content"`
	Status  int32  `json:"status"`
	Author  string `json:"author"`
	Ctime   string `json:"ctime"`
	Utime   string `json:"utime"`

	// 点赞之类的信息
	LikeCnt    int64 `json:"likeCnt"`
	CollectCnt int64 `json:"collectCnt"`
	ReadCnt    int64 `json:"readCnt"`

	// 个人是否点赞的信息
	Liked     bool `json:"liked"`
	Collected bool `json:"collected"`
}

type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (req ArticleReq) ToDTO(uid int64) *articlev1.Article {
	return &articlev1.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: &articlev1.Author{
			Id: uid,
		},
	}
}
