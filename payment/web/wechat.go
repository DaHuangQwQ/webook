package web

import (
	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"net/http"
	"webook/payment/service/wechat"
	"webook/pkg/ginx"
	"webook/pkg/logger"
)

type WechatHandler struct {
	handler   *notify.Handler
	l         logger.LoggerV1
	nativeSvc wechat.NativePaymentService
}

func NewWechatHandler(handler *notify.Handler, l logger.LoggerV1, nativeSvc wechat.NativePaymentService) *WechatHandler {
	return &WechatHandler{handler: handler, l: l, nativeSvc: nativeSvc}
}

func (h *WechatHandler) RegisterRoutes(router *gin.Engine) {
	router.Any("/pay/callback", h.HandlerNative)
}

func (h *WechatHandler) HandlerNative(ctx *gin.Context) {
	txn := new(payments.Transaction)
	_, err := h.handler.ParseNotifyRequest(ctx.Request.Context(), ctx.Request, &txn)
	if err != nil {
		h.l.Warn("有人伪造微信支付请求", logger.Error(err))
		return
	}
	// 交易信息已经被解密
	err = h.nativeSvc.HandlerCallback(ctx, txn)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ginx.Result{})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "ok",
	})
}
