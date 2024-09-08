package api

import "webook/internal/domain"

type RecruitmentAddReq struct {
	Meta `path:"/recruit/add" method:"post"`
	domain.Recruitment
}
