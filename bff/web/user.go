package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"webook/code/repository/cache"
	service2 "webook/code/service"
	ijwt "webook/internal/web/jwt"
	"webook/pkg/ginx"
	"webook/pkg/logger"
	"webook/user/domain"
	"webook/user/service"
)
import "github.com/gin-gonic/gin"

// 确保实现了 Handler 接口
var _ Handler = (*UserHandler)(nil)

type UserHandler struct {
	ijwt.Handler
	svc         service.UserService
	codeSvc     service2.CodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	phoneExp    *regexp.Regexp

	l logger.LoggerV1
}

const biz = "login"

func NewUserHandler(svc service.UserService, codeSvc service2.CodeService, hdl ijwt.Handler, l logger.LoggerV1) *UserHandler {
	const (
		emailRegexPattern  = "^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$"
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
		l:           l,
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
	ug.POST("/phone_update", h.PhoneUpdate)

	ug.POST("/info_update", h.InfoUpdate)
	ug.GET("/get_info", h.GetInfo)

	ug.POST("/avatar_update", h.AvatarUpdate)
	ug.GET("/get_avatar", h.GetAvatar)

	ug.GET("/demo", h.Demo)
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

	err = h.svc.Signup(ctx.Request.Context(), domain.User{
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
	u, err := h.svc.Login(ctx.Request.Context(), req.Email, req.Password)
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
	u, err := h.svc.Login(ctx.Request.Context(), req.Email, req.Password)
	switch err {
	case nil:
		err = h.SetLoginToken(ctx, u.Id)
		if err != nil {
			ctx.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "系统错误",
			})
			return
		}
		ctx.JSON(http.StatusOK, Result{
			Code: 0,
			Msg:  "登录成功",
		})
	case service.ErrInvalidUserOrPassword:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "用户名或者密码不对",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
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

	err = h.codeSvc.Send(ctx.Request.Context(), biz, req.Phone)
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
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	// 正则表达式
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "请输入手机号",
		})
		return
	}
	verify, err := h.codeSvc.Verify(ctx.Request.Context(), biz, req.Phone, req.Code)
	switch err {
	case nil:
	case cache.ErrCodeVerifyTooMany:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "验证太频繁",
		})
		return
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !verify {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "验证码不正确",
		})
		return
	}
	user, err := h.svc.FindOrCreate(ctx.Request.Context(), req.Phone)
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

func (h *UserHandler) InfoUpdate(ctx *gin.Context) {
	type Req struct {
		Name   string `json:"username"`
		Email  string `json:"email"`
		Grade  int    `json:"grade"`
		Gender int    `json:"gender"`
		Phone  string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "参数错误",
		})
		return
	}
	reqNil := Req{}
	if req == reqNil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "参数为空",
		})
		return
	}
	if req.Phone != "" {
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
				Code: 5,
				Msg:  "非法手机格式",
			})
			return
		}
	}
	if req.Email != "" {
		isEmail, err := h.emailExp.MatchString(req.Email)
		if err != nil {
			ctx.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "系统错误",
			})
			return
		}
		if !isEmail {
			ctx.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "非法邮箱格式",
			})
			return
		}
	}

	h.l.Info("信息更新", logger.Field{Key: "info", Val: req})
	userId := ctx.MustGet("claims").(ijwt.UserClaims)
	err := h.svc.UpdateByID(ctx.Request.Context(), domain.User{
		Id:       userId.Uid,
		Nickname: req.Name,
		Grade:    req.Grade,
		Gender:   req.Gender,
		Email:    req.Email,
		Phone:    req.Phone,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "ok",
	})
}

func (h *UserHandler) GetInfo(ctx *gin.Context) {
	UserId := ctx.MustGet("claims").(ijwt.UserClaims)
	UserInfo, err := h.svc.FindByID(ctx.Request.Context(), UserId.Uid)
	h.l.Info("UserInfo", logger.Field{Key: "err", Val: UserInfo})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	//phoneRes := ""
	//if UserInfo.Phone != "" {
	//	phoneRes = UserInfo.Phone[:3] + "****" + UserInfo.Phone[7:]
	//}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "ok",
		Data: domain.UserInfo{
			Nickname: UserInfo.Nickname,
			Phone:    UserInfo.Phone,
			Grade:    UserInfo.Grade,
			Gender:   UserInfo.Gender,
			Avatar:   UserInfo.Avatar,
			Email:    UserInfo.Email,
		},
	})
}

func (h *UserHandler) AvatarUpdate(ctx *gin.Context) {
	file, fileTypes, err := ctx.Request.FormFile("file")
	if file == nil || err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "图片上传失败",
		})
		return
	}
	if fileTypes.Size > 2*1024*1024 {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "文件不应该超过2MB",
		})
		return
	}
	fileType := strings.Split(fileTypes.Header.Get("Content-Type"), "/")[1]
	user := ctx.MustGet("claims").(ijwt.UserClaims)
	ossAdress, err := h.svc.AvatarUpdate(ctx.Request.Context(), user.Uid, file, fileType)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Info("AvatarUpdate错误", logger.Field{Key: "err", Val: err})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "ok",
		Data: ossAdress,
	})
	file.Close()
}

func (h *UserHandler) GetAvatar(ctx *gin.Context) {
	userId := ctx.MustGet("claims").(ijwt.UserClaims)
	avatar, err := h.svc.GetAvatar(ctx.Request.Context(), userId.Uid)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "ok",
		Data: avatar,
	})
}

func (h *UserHandler) PhoneUpdate(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	// 正则表达式
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "请输入手机号",
		})
		return
	}
	verify, err := h.codeSvc.Verify(ctx.Request.Context(), biz, req.Phone, req.Code)
	switch err {
	case nil:
	case cache.ErrCodeVerifyTooMany:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "验证太频繁",
		})
		return
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Info("短信错误", logger.Field{Key: "err", Val: err})
		return
	}
	if !verify {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "验证码不正确",
		})
		return
	}
	UserId := ctx.MustGet("claims").(ijwt.UserClaims)
	err = h.svc.UpdateByID(ctx.Request.Context(), domain.User{
		Id:    UserId.Uid,
		Phone: req.Phone,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Info("信息更新错误", logger.Field{Key: "err", Val: err})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "ok",
	})
}

func (h *UserHandler) Demo(ctx *gin.Context) {
	str := `
{
  "userList": [
    {
    "id": 31,
      "email": "2195921486@qq.com",
      "mobile": "13029375111",
       "userName": "黄俊烨",
      "grade": 2022,
      "gender": 1,
      "avatar": "https://ceit.oss-cn-wulanchabu.aliyuncs.com/avatar/1725809043994.jpeg"
    },
    {
    "id": 27,
      "email": "2195921486@qq.com",
      "mobile": "13029375111",
       "userName": "肖岳阳",
      "grade": 2022,
      "gender": 1,
      "avatar": "https://ceit.oss-cn-wulanchabu.aliyuncs.com/avatar/1725871599659.jpeg"
    },
    {
    "id": 28,
      "email": "2195921486@qq.com",
      "mobile": "13029375111",
       "userName": "王博涵",
      "grade": 2022,
      "gender": 1,
      "avatar": "https://ceit.oss-cn-wulanchabu.aliyuncs.com/avatar/1725871559060.jpeg"
    },
    {
    "id": 29,
      "email": "2195921486@qq.com",
      "mobile": "13029375111",
       "userName": "文博熙",
      "grade": 2022,
      "gender": 1,
      "avatar": "https://ceit.oss-cn-wulanchabu.aliyuncs.com/avatar/1725871578342.jpeg"
    },
    {
      "id": 25,
      "email": "2195921486@qq.com",
      "mobile": "13029375111",
      "password": "",
      "userName": "尹欣雨",
      "grade": 2023,
      "gender": 2,
      "avatar": "https://ceit.oss-cn-wulanchabu.aliyuncs.com/avatar/1725871637021.jpeg",
      "ctime": "2024-09-08T23:23:30.374+08:00",
      "WechatInfo": {
        "UnionId": "",
        "OpenId": "oZ7s043gWIyjR_W0f5qhZ6avje9k"
      },
      "birthday": 0,
      "userStatus": 1,
      "deptId": 0,
      "remark": "",
      "isAdmin": 1,
      "address": "",
      "describe": "",
      "lastLoginIp": "",
      "dept": {
        "deptId": 0,
        "parentId": 0,
        "ancestors": "",
        "deptName": "",
        "orderNum": 0,
        "leader": "",
        "phone": "",
        "email": "",
        "status": 0
      },
      "roleInfo": [
        {
          "roleId": 1,
          "name": "root"
        },
        {
          "roleId": 2,
          "name": "root"
        }
      ]
    },
    {
      "id": 26,
      "email": "1623929456@qq.com",
      "mobile": "15655440534",
      "password": "",
      "userName": "城猫",
      "grade": 2022,
      "gender": 1,
      "avatar": "https://ceit.oss-cn-wulanchabu.aliyuncs.com/avatar/1725810736043.jpeg",
      "ctime": "2024-09-08T23:29:14.809+08:00",
      "WechatInfo": {
        "UnionId": "",
        "OpenId": "oZ7s044EqiPnrbbHeILelQL00Oss"
      },
      "birthday": 0,
      "userStatus": 1,
      "deptId": 0,
      "remark": "",
      "isAdmin": 1,
      "address": "",
      "describe": "",
      "lastLoginIp": "",
      "dept": {
        "deptId": 0,
        "parentId": 0,
        "ancestors": "",
        "deptName": "",
        "orderNum": 0,
        "leader": "",
        "phone": "",
        "email": "",
        "status": 0
      },
      "roleInfo": [
        {
          "roleId": 1,
          "name": "root"
        },
        {
          "roleId": 2,
          "name": "root"
        }
      ]
    },
    {
      "id": 30,
      "email": "2195921481@qq.com",
      "mobile": "",
      "password": "$2a$10$ItQPeBXtoTA2BMg.41drIO9vbslG2jRr.coChOvpEhwuuCv9/Pct.",
      "userName": "",
      "grade": 0,
      "gender": 0,
      "avatar": "",
      "ctime": "2024-09-08T23:39:48.937+08:00",
      "WechatInfo": {
        "UnionId": "",
        "OpenId": ""
      },
      "birthday": 0,
      "userStatus": 1,
      "deptId": 0,
      "remark": "",
      "isAdmin": 1,
      "address": "",
      "describe": "",
      "lastLoginIp": "",
      "dept": {
        "deptId": 0,
        "parentId": 0,
        "ancestors": "",
        "deptName": "",
        "orderNum": 0,
        "leader": "",
        "phone": "",
        "email": "",
        "status": 0
      },
      "roleInfo": [
        {
          "roleId": 1,
          "name": "root"
        },
        {
          "roleId": 2,
          "name": "root"
        }
      ]
    },
    {
      "id": 31,
      "email": "",
      "mobile": "",
      "password": "",
      "userName": "",
      "grade": 0,
      "gender": 0,
      "avatar": "",
      "ctime": "2024-09-08T23:46:56.975+08:00",
      "WechatInfo": {
        "UnionId": "",
        "OpenId": "oZ7s048EKnAjql6Zr9NklXUOxmo0"
      },
      "birthday": 0,
      "userStatus": 1,
      "deptId": 0,
      "remark": "",
      "isAdmin": 1,
      "address": "",
      "describe": "",
      "lastLoginIp": "",
      "dept": {
        "deptId": 0,
        "parentId": 0,
        "ancestors": "",
        "deptName": "",
        "orderNum": 0,
        "leader": "",
        "phone": "",
        "email": "",
        "status": 0
      },
      "roleInfo": [
        {
          "roleId": 1,
          "name": "root"
        },
        {
          "roleId": 2,
          "name": "root"
        }
      ]
    }
  ],
  "currentPage": 0,
  "total": 4,
  "list": null
}

`
	//strs := strings.Replace(str, "\n", "", -1)
	//strs = strings.Replace(strs, "\\", "", -1)
	ctx.JSON(http.StatusOK, ginx.Result{
		Code: 0,
		Msg:  "ok",
		Data: str,
	})
}
