package web

import (
	"github.com/DaHuangQwQ/webook/bff/api"
	"github.com/DaHuangQwQ/webook/internal_temp/service"
	"github.com/DaHuangQwQ/webook/pkg/ginx"
	"github.com/gin-gonic/gin"
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
	//err := h.articleSvc.Add(ctx, req.Recruitment)
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
