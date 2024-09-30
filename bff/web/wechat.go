package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	ijwt "webook/internal/web/jwt"
	"webook/pkg/logger"
	"webook/user/service"
	"webook/user/service/oauth2/wechat"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	ijwt.Handler
	key             []byte
	stateCookieName string

	l logger.LoggerV1
}

func NewOAuth2WechatHandler(svc wechat.Service,
	hdl ijwt.Handler,
	userSvc service.UserService,
	l logger.LoggerV1,
) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:             svc,
		userSvc:         userSvc,
		key:             []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgB"),
		stateCookieName: "jwt-state",
		Handler:         hdl,
		l:               l,
	}
}

func (o *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", o.Auth2URL)
	g.Any("/callback", o.Callback)
	g.Any("/code2session", o.Code2Session)
}

func (o *OAuth2WechatHandler) Auth2URL(ctx *gin.Context) {
	state := uuid.New()
	val, err := o.svc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "构造跳转URL失败",
			Code: 5,
		})
		return
	}
	println(val)
	err = o.setStateCookie(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "服务器异常",
			Code: 5,
		})
	}

	//token := jwt.NewWithClaims(jwt.SigningMethodHS512)
	//tokenStr, err := token.SignedString()
	//if err != nil {
	//	ctx.JSON(http.StatusOK, Result{})
	//}
	//ctx.SetCookie("jwt-state", tokenStr, 600, "/", "", false, true)

	ctx.JSON(http.StatusOK, Result{
		Data: val,
	})
}

func (o *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	err := o.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "非法请求",
			Code: 4,
		})
		o.l.Info("微信登入验证cookie", logger.Field{Key: "err", Val: err})
		return
	}
	// 你校验不校验都可以
	code := ctx.Query("code")
	// state := ctx.Query("state")
	wechatInfo, err := o.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "授权码有误",
			Code: 4,
		})
		o.l.Info("微信accessToken申请失败", logger.Field{Key: "err", Val: err})
		return
	}
	u, err := o.userSvc.FindOrCreateByWechat(ctx, wechatInfo)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		return
	}
	err = o.SetLoginToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
	return
}

func (o *OAuth2WechatHandler) Code2Session(ctx *gin.Context) {
	code := ctx.Query("code")
	wechatInfo, err := o.svc.Code2Session(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		o.l.Info("Code2Session服务出错", logger.Field{Key: "err", Val: err})
		return
	}
	u, err := o.userSvc.FindOrCreateByWechat(ctx, wechatInfo)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		return
	}
	err = o.SetLoginToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
	return
}

func (o *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	//state := ctx.Query("state")
	//ck, err := ctx.Cookie(o.stateCookieName)
	//if err != nil {
	//	return fmt.Errorf("无法获得 cookie %w", err)
	//}
	//var sc StateClaims
	//_, err = jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
	//	return o.key, nil
	//})
	//if err != nil {
	//	return fmt.Errorf("解析 token 失败 %w", err)
	//}
	//if state != sc.State {
	//	// state 不匹配，有人搞你
	//	return fmt.Errorf("state 不匹配")
	//}
	return nil
}

func (o *OAuth2WechatHandler) setStateCookie(ctx *gin.Context,
	state string) error {
	claims := StateClaims{
		State: state,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(o.key)
	if err != nil {

		return err
	}
	ctx.SetCookie(o.stateCookieName, tokenStr,
		600, "/oauth2/wechat/callback",
		"", false, true)
	return nil
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string
}
