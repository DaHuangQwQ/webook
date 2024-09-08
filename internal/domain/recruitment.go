package domain

type Recruitment struct {
	Id          uint64 `json:"id"`          // id (bigint)
	Name        string `json:"name"`        // 姓名 (varchar(10))
	StudentID   string `json:"studentId"`   // 学号 (varchar(10))
	Major       uint8  `json:"major"`       // 专业 (tinyint)
	Situation   string `json:"situation"`   // 基础情况 (text)
	Expectation string `json:"expectation"` // 未来期望 (text)
	Selfie      string `json:"selfie"`      // 自拍图片 (varchar(255))
	ErrorNum    int    `json:"errorNum"`    // 错误次数 (int)
}
