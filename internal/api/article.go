package api

import "webook/internal/domain"

type GetListReq struct {
	Meta `path:"/articles/list" method:"get"`
	PageReq
}

type GetListRes struct {
	ListRes
}

type LikeReq struct {
	Meta      `path:"/articles/like" method:"post"`
	ArticleID int64 `json:"article_id"`
	Like      bool  `json:"like"`
}

type ArticleEditReq struct {
	Meta `path:"/articles/edit" method:"post"`
	ArticleReq
}

type ArticlePublishReq struct {
	Meta `path:"/articles/publish" method:"post"`
	ArticleReq
}

type ArticleWithdrawReq struct {
	Meta `path:"/articles/withdraw" method:"post"`
	Id   int64 `json:"id"`
}

type ArticleGetListReq struct {
	Meta `path:"/articles/getlist" method:"get"`
}

type Article struct {
	Id        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`

	Status     uint8  `json:"status"`
	Author     string `json:"author"`
	ReadCnt    int64  `json:"read_cnt"`
	LikeCnt    int64  `json:"like_cnt"`
	CollectCnt int64  `json:"collect_cnt"`

	Liked     bool `json:"liked"`
	Collected bool `json:"collected"`
}

type ArticleReq struct {
	Id      int64  `json:"id"`
	ImgUrl  string `json:"imgUrl"`
	Title   string `json:"title"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

func (art *ArticleReq) ToDomain(uid int64) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: uid,
		},
		ImgUrl: art.ImgUrl,
		Type:   art.Type,
	}
}
