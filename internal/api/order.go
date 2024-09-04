package api

type OrderListReq struct {
	Meta `path:"/orders/list" method:"get"`
}

type OrderAddReq struct {
	Meta     `path:"/orders/add" method:"post"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Author   string `json:"author"`
	Salary   string `json:"salary"`
	Duration int64  `json:"duration"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
}

type OrderDeleteReq struct {
	Meta `path:"/orders/delete" method:"post"`
	Ids  []int64 `json:"ids"`
}

type OrderEditReq struct {
	Meta     `path:"/orders/edit" method:"post"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Author   string `json:"author"`
	Salary   string `json:"salary"`
	Duration int64  `json:"duration"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
}
