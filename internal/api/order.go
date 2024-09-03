package api

type OrderListReq struct {
}

type OrderAddReq struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Author   string `json:"author"`
	Salary   string `json:"salary"`
	Duration int64  `json:"duration"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
}

type OrderDeleteReq struct {
	Ids []int64 `json:"ids"`
}

type OrderEditReq struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Author   string `json:"author"`
	Salary   string `json:"salary"`
	Duration int64  `json:"duration"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
}
