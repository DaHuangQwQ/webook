package api

type Meta struct{}

// PageReq 公共请求参数
type PageReq struct {
	DateRange []string `json:"dateRange"` //日期范围
	PageNum   int      `json:"pageNum"`   //当前页码
	PageSize  int      `json:"pageSize"`  //每页数
	OrderBy   string   //排序方式
}

// ListRes 列表公共返回
type ListRes struct {
	CurrentPage int         `json:"currentPage"`
	Total       interface{} `json:"total"`
	List        []any       `json:"list"`
}

type ListReq struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}
