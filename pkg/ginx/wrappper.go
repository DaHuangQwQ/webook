package ginx

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webook/pkg/logger"
)

var L logger.LoggerV1

//func init() {
//	L = ioc.InitLogger()
//}

func Warp[Req any](fn func(ctx *gin.Context, req Req) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "参数错误" + err.Error(),
			})
			return
		}

		result, err := fn(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusOK, result)
			//L.Info("系统错误", logger.Field{Key: "err", Val: err})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
