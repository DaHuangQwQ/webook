package api

type GetListReq struct {
	PageReq
}

type GetListRes struct {
	ListRes
}

type LikeReq struct {
	ArticleID int64 `json:"article_id"`
	Like      bool  `json:"like"`
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
