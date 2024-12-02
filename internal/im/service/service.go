package service

import (
	"context"
	"fmt"
	"github.com/DaHuangQwQ/webook/im/domain"
	"github.com/ecodeclub/ekit/net/httpx"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type UserService interface {
	Sync(ctx context.Context, user domain.User) error
}

type RESTUserService struct {
	// 部署 IM 时候配置的 IM Secret，默认是 openIM123
	secret string
	base   string
	client *http.Client
}

func NewRESTUserService(secret string, base string) *RESTUserService {
	return &RESTUserService{
		secret: secret,
		base:   base}
}

func (svc *RESTUserService) Sync(ctx context.Context, user domain.User) error {
	spanCtx := trace.SpanContextFromContext(ctx)
	var traceId string
	if spanCtx.HasSpanID() {
		traceId = spanCtx.TraceID().String()
	} else {
		// 随便生成一个，但是这样链路就拼接不起来了
		traceId = uuid.New().String()
	}
	var resp response
	err := httpx.NewRequest(ctx, http.MethodPost,
		svc.base+"/user/user_register").JSONBody(syncUserRequest{
		Secret: svc.secret,
		Users:  []domain.User{user},
	}).Client(svc.client).
		AddHeader("operationID", traceId).
		Do().JSONScan(&resp)
	if err != nil {
		return err
	}
	if resp.ErrCode != 0 {
		return fmt.Errorf("同步数据失败 %d, %s, %s", resp.ErrCode, resp.ErrMsg, resp.ErrDlt)
	}
	return nil
}

type syncUserRequest struct {
	Secret string        `json:"secret"`
	Users  []domain.User `json:"users"`
}

type response struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
}
