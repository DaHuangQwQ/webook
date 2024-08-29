package domain

import "time"

type User struct {
	Id    int64
	Email string

	Phone    string
	Password string

	Nickname string
	Grade    int

	Gender int
	Avatar string

	CTime time.Time

	WechatInfo WechatInfo
}

type UserInfo struct {
	Nickname string `json:"nickname"`
	Grade    int    `json:"grade"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	// 0 未知， 1 男， 2 女
	Gender int `json:"gender"`
}
