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

	CTime time.Time

	WechatInfo WechatInfo
}

type UserInfo struct {
	Nickname string
	Grade    int
	Phone    string
	// 0 未知， 1 男， 2 女
	Gender int
}
