package web

import (
	"github.com/gin-gonic/gin"
	articlev1 "webook/api/proto/gen/article/v1"
	rewardv1 "webook/api/proto/gen/reward/v1"
	"webook/bff/api"
	"webook/pkg/ginx"
)

type RewardHandler struct {
	client    rewardv1.RewardServiceClient
	artClient articlev1.ArticleServiceClient
}

func NewRewardHandler(client rewardv1.RewardServiceClient, artClient articlev1.ArticleServiceClient) *RewardHandler {
	return &RewardHandler{client: client, artClient: artClient}
}

func (h *RewardHandler) RegisterRoutes(server *gin.Engine) {
	//rg := server.Group("/reward")
	server.POST(ginx.WarpWithToken[api.GetRewardReq](h.GetReward))
}

func (h *RewardHandler) GetReward(
	ctx *gin.Context,
	req api.GetRewardReq,
	claims ginx.UserClaims) (ginx.Result, error) {
	resp, err := h.client.GetReward(ctx, &rewardv1.GetRewardRequest{
		Rid: req.Rid,
		Uid: claims.Id,
	})
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return ginx.Result{
		// 暂时也就是只需要状态
		Data: resp.Status.String(),
	}, nil
}

type RewardArticleReq struct {
	Aid int64 `json:"aid"`
	Amt int64 `json:"amt"`
}
