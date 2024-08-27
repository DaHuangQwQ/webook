package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/service"
	ijwt "webook/internal/web/jwt"
)
import "github.com/gin-gonic/gin"

// 确保实现了 handler 接口
var _ handler = (*UserHandler)(nil)

type UserHandler struct {
	ijwt.Handler
	svc         service.UserService
	codeSvc     service.CodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	phoneExp    *regexp.Regexp
}

const biz = "login"

func NewUserHandler(svc service.UserService, codeSvc service.CodeService, hdl ijwt.Handler) *UserHandler {
	const (
		emailRegexPattern  = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordExpPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
		phoneExpPattern    = "^1(3[0-9]|5[0-3,5-9]|7[1-3,5-8]|8[0-9])\\d{8}$"
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordExpPattern, regexp.None)
	phoneExp := regexp.MustCompile(phoneExpPattern, regexp.None)

	return &UserHandler{
		svc:         svc,
		codeSvc:     codeSvc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		phoneExp:    phoneExp,
		Handler:     hdl,
	}
}

func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	ug := router.Group("/users")
	ug.POST("/signup", h.SignUp)
	ug.POST("/login", h.LoginJWT)
	ug.POST("/logout", h.LogoutJWT)

	ug.POST("/edit", h.Edit)
	ug.GET("/profile", h.Profile)

	ug.GET("/refresh_token", h.RefreshToken)

	// 手机验证码登录相关功能
	ug.POST("/login_sms/code/send", h.SendLoginSmsCode)
	ug.POST("/login_sms", h.LoginSms)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	isEmail, err := h.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "非法邮箱格式")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入密码不对")
		return
	}

	isPassword, err := h.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含字母、数字、特殊字符，并且不少于八位")
		return
	}

	err = h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		ctx.String(http.StatusOK, "注册成功")
	case service.ErrDuplicate:
		ctx.String(http.StatusOK, "邮箱冲突，请换一个")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) LoginSession(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		sess := sessions.Default(ctx)
		sess.Set("userId", u.Id)
		sess.Options(sessions.Options{
			// 十五分钟
			MaxAge: 900,
		})
		err = sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) LoginJWT(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		err = h.SetLoginToken(ctx, u.Id)
		if err != nil {
			ctx.JSON(http.StatusOK, "系统错误")
			return
		}
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
	return
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "profile")
}

func (h *UserHandler) Edit(ctx *gin.Context) {
	type Req struct{}
}

func (h *UserHandler) SendLoginSmsCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	isPhone, err := h.phoneExp.MatchString(req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !isPhone {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "请输入手机号码",
		})
		return
	}
	//if req.Phone == "" {
	//	ctx.JSON(http.StatusOK, Result{
	//		Code: 4,
	//		Msg:  "请输入手机号码",
	//	})
	//	return
	//}

	err = h.codeSvc.Send(ctx, biz, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case cache.ErrSendCodeTooMany:
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "短信发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
}

func (h *UserHandler) LoginSms(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	// 正则表达式
	if req.Phone == "" {
		ctx.String(http.StatusOK, "请输入手机号")
		return
	}
	verify, err := h.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	switch err {
	case nil:
	case cache.ErrCodeVerifyTooMany:
		ctx.String(http.StatusOK, "验证太频繁")
		return
	default:
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !verify {
		ctx.String(http.StatusOK, "验证码不正确")
		return
	}
	user, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}
	err = h.SetLoginToken(ctx, user.Id)
	if err != nil {
		return
	}
	ctx.String(http.StatusOK, "验证码正确")
}

func (h *UserHandler) RefreshToken(ctx *gin.Context) {
	// 约定，前端在 Authorization 里面带上这个 refresh_token
	tokenStr := h.ExtractToken(ctx)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RCJWTKey, nil
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if token == nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.CheckSession(ctx, rc.Ssid)
	if err != nil {
		// token 无效或者 redis 有问题
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.SetJWTToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

func (h *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := h.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Msg: "退出登录成功"})
}
