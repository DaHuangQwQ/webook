package domain

import "time"

type Comment struct {
	Id int64 `json:"id"`
	// 评论者
	Commentator User `json:"user"`
	// 评论对象
	// 数据里面
	Biz   string `json:"biz"`
	BizID int64  `json:"bizid"`
	// 评论对象
	Content string `json:"content"`
	// 根评论
	RootComment *Comment `json:"rootComment"`
	// 父评论
	ParentComment *Comment  `json:"parentComment"`
	Children      []Comment `json:"children"`
	CTime         time.Time `json:"ctime"`
	UTime         time.Time `json:"utime"`
}

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
