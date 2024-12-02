package domain

type User struct {
	Nickname string `json:"nickname"`
	UserID   string `json:"userID"`
	// 这个实际上是头像
	FaceURL string `json:"faceURL"`
}
