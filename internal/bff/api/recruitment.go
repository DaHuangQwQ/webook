package api

import "github.com/DaHuangQwQ/webook/internal_temp/domain"

type RecruitmentAddReq struct {
	Meta `path:"/recruit/add" method:"post"`
	domain.Recruitment
}
