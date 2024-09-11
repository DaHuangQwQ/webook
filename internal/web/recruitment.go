package web

import (
	"github.com/gin-gonic/gin"
	"webook/internal/api"
	"webook/internal/service"
	"webook/pkg/ginx"
)

var _ Handler = (*RecruitmentHandler)(nil)

type RecruitmentHandler struct {
	svc service.RecruitmentService
}

func NewRecruitmentHandler(svc service.RecruitmentService) *RecruitmentHandler {
	return &RecruitmentHandler{
		svc: svc,
	}
}

func (h *RecruitmentHandler) RegisterRoutes(router *gin.Engine) {
	router.POST(ginx.Warp[api.RecruitmentAddReq](h.Add))
}

func (h *RecruitmentHandler) Add(ctx *gin.Context, req api.RecruitmentAddReq) (ginx.Result, error) {
	//err := h.svc.Add(ctx, req.Recruitment)
	//if err != nil {
	//	return ginx.Result{
	//		Code: 5,
	//		Msg:  "系统错误" + err.Error(),
	//	}, err
	//}
	return ginx.Result{
		Code: 5,
		Msg:  "纳新截止",
	}, nil
}
