package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"time"
	codev1 "webook/api/proto/gen/code/v1"
	userv1 "webook/api/proto/gen/user/v1"
	"webook/bff/api"
	ijwt "webook/bff/web/jwt"
	"webook/pkg/ginx"
	"webook/user/errs"
)
import "github.com/gin-gonic/gin"

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`

	userIdKey = "userId"
	bizLogin  = "login"
)

var _ Handler = &UserHandler{}

type UserHandler struct {
	svc              userv1.UserServiceClient
	codeSvc          codev1.CodeServiceClient
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
	ijwt.Handler
}

func NewUserHandler(svc userv1.UserServiceClient,
	codeSvc codev1.CodeServiceClient, jwthdl ijwt.Handler) *UserHandler {
	return &UserHandler{
		svc:              svc,
		codeSvc:          codeSvc,
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		Handler:          jwthdl,
	}
}

func (c *UserHandler) RegisterRoutes(server *gin.Engine) {

	// 分组注册
	ug := server.Group("/users")
	server.POST(ginx.Warp[api.UserSignUpReq](c.SignUp))
	// session 机制
	//ug.POST("/login", c.Login)
	// JWT 机制
	server.POST(ginx.Warp[api.UserLoginReq](c.LoginJWT))
	ug.POST("/logout", c.Logout)
	ug.POST("/edit", c.Edit)
	//ug.GET("/profile", c.Profile)
	ug.GET("/profile", c.ProfileJWT)
	ug.POST("/login_sms/code/send", c.SendSMSLoginCode)
	ug.POST("/login_sms", c.LoginSMS)
	ug.POST("/refresh_token", c.RefreshToken)
	//ug.POST("/avatar_update", c.AvatarUpdate)
}

func (c *UserHandler) RefreshToken(ctx *gin.Context) {
	// 假定长 token 也放在这里
	tokenStr := c.ExtractTokenString(ctx)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RefreshTokenKey, nil
	})
	// 这边要保持和登录校验一直的逻辑，即返回 401 响应
	if err != nil || token == nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, Result{Code: 4, Msg: "请登录"})
		return
	}

	// 校验 ssid
	err = c.CheckSession(ctx, rc.Ssid)
	if err != nil {
		// 系统错误或者用户已经主动退出登录了
		// 这里也可以考虑说，如果在 Redis 已经崩溃的时候，
		// 就不要去校验是不是已经主动退出登录了。
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = c.SetJWTToken(ctx, rc.Ssid, rc.Id)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, Result{Code: 4, Msg: "请登录"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Msg: "刷新成功"})
}

func (c *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	resp, err := c.codeSvc.Verify(ctx, &codev1.VerifyRequest{
		Biz: bizLogin, Phone: req.Phone, InputCode: req.Code,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统异常"})
		zap.L().Error("用户手机号码登录失败", zap.Error(err))
		return
	}
	if resp.Answer {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "验证码错误"})
		return
	}

	// 验证码是对的
	// 登录或者注册用户
	u, err := c.svc.FindOrCreate(ctx, &userv1.FindOrCreateRequest{
		Phone: req.Phone,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "系统错误"})
		return
	}
	// 用 uuid 来标识这一次会话
	ssid := uuid.New().String()
	err = c.SetJWTToken(ctx, ssid, u.User.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Msg: "登录成功"})
}

// SendSMSLoginCode 发送短信验证码
func (c *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 你也可以用正则表达式校验是不是合法的手机号
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "请输入手机号码"})
		return
	}
	_, err := c.codeSvc.Send(ctx, &codev1.CodeSendRequest{
		Biz: bizLogin, Phone: req.Phone,
	})
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{Msg: "发送成功"})
	//case .ErrCodeSendTooMany:
	// TODO 利用 grpc 来传递错误码
	//	ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "短信发送太频繁，请稍后再试"})
	default:
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		// 要打印日志
		return
	}
}

// SignUp 用户注册接口
func (c *UserHandler) SignUp(ctx *gin.Context, req api.UserSignUpReq) (ginx.Result, error) {

	isEmail, err := c.emailRegexExp.MatchString(req.Email)
	if err != nil {
		return Result{
			Code: errs.UserInternalServerError,
			Msg:  "系统错误",
		}, err
	}
	if !isEmail {
		return Result{
			Code: errs.UserInvalidInput,
			Msg:  "邮箱输入错误",
		}, nil
	}

	if req.Password != req.ConfirmPassword {
		return Result{
			Code: errs.UserInvalidInput,
			Msg:  "两次输入密码不对",
		}, nil
	}

	isPassword, err := c.passwordRegexExp.MatchString(req.Password)
	if err != nil {
		return Result{
			Code: errs.UserInvalidInput,
			Msg:  "系统错误",
		}, err
	}
	if !isPassword {
		return Result{
			Code: errs.UserInvalidInput,
			Msg:  "密码必须包含数字、特殊字符，并且长度不能小于 8 位",
		}, nil
	}

	_, err = c.svc.Signup(ctx.Request.Context(), &userv1.SignupRequest{User: &userv1.User{Email: req.Email, Password: req.ConfirmPassword}})
	if err != nil {
		return Result{
			Code: errs.UserInternalServerError,
			Msg:  "系统错误",
		}, err
	}
	return Result{
		Msg: "OK",
	}, nil
}

// LoginJWT 用户登录接口，使用的是 JWT，如果你想要测试 JWT，就启用这个
func (c *UserHandler) LoginJWT(ctx *gin.Context, req api.UserLoginReq) (ginx.Result, error) {
	u, err := c.svc.Login(ctx.Request.Context(), &userv1.LoginRequest{
		Email: req.Email, Password: req.Password,
	})

	if err != nil {
		return ginx.Result{}, err
	}
	err = c.SetLoginToken(ctx, u.User.Id)
	if err != nil {
		return ginx.Result{}, err
	}
	return ginx.Result{Msg: "登录成功"}, nil
}

func (c *UserHandler) Logout(ctx *gin.Context) {
	err := c.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

// Login 用户登录接口
func (c *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	// 当我们调用 Bind 方法的时候，如果有问题，Bind 方法已经直接写响应回去了
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := c.svc.Login(ctx.Request.Context(), &userv1.LoginRequest{
		Email: req.Email, Password: req.Password})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	// TODO 利用 grpc 来传递错误码
	//if err == service.ErrInvalidUserOrPassword {
	//	ctx.String(http.StatusOK, "用户名或者密码不正确，请重试")
	//	return
	//}
	sess := sessions.Default(ctx)
	sess.Set(userIdKey, u.User.Id)
	sess.Options(sessions.Options{
		// 60 秒过期
		MaxAge: 60,
	})
	err = sess.Save()
	if err != nil {
		ctx.String(http.StatusOK, "服务器异常")
		return
	}
	ctx.String(http.StatusOK, "登录成功")
}

// Edit 用户编译信息
func (c *UserHandler) Edit(ctx *gin.Context) {
	type Req struct {
		// 注意，其它字段，尤其是密码、邮箱和手机，
		// 修改都要通过别的手段
		// 邮箱和手机都要验证
		// 密码更加不用多说了
		Nickname string `json:"nickname"`
		// 2023-01-01
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 你可以尝试在这里校验。
	// 比如说你可以要求 Nickname 必须不为空
	// 校验规则取决于产品经理
	if req.Nickname == "" {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "昵称不能为空"})
		return
	}

	if len(req.AboutMe) > 1024 {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "关于我过长"})
		return
	}
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		// 也就是说，我们其实并没有直接校验具体的格式
		// 而是如果你能转化过来，那就说明没问题
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "日期格式不对"})
		return
	}

	uc := ctx.MustGet("user").(ijwt.UserClaims)
	_, err = c.svc.UpdateNonSensitiveInfo(ctx,
		&userv1.UpdateNonSensitiveInfoRequest{
			User: &userv1.User{
				Id:       uc.Id,
				Nickname: req.Nickname,
				AboutMe:  req.AboutMe,
				Birthday: timestamppb.New(birthday),
			},
		})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Msg: "OK"})
}

// ProfileJWT 用户详情, JWT 版本
func (c *UserHandler) ProfileJWT(ctx *gin.Context) {
	type Profile struct {
		Email    string
		Phone    string
		Nickname string
		Birthday string
		AboutMe  string
	}
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	resp, err := c.svc.Profile(ctx, &userv1.ProfileRequest{Id: uc.Id})
	if err != nil {
		// 按照道理来说，这边 id 对应的数据肯定存在，所以要是没找到，
		// 那就说明是系统出了问题。
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	u := resp.User
	ctx.JSON(http.StatusOK, Profile{
		Email:    u.Email,
		Phone:    u.Phone,
		Nickname: u.Nickname,
		Birthday: u.Birthday.AsTime().Format(time.DateOnly),
		AboutMe:  u.AboutMe,
	})
}

// Profile 用户详情
func (c *UserHandler) Profile(ctx *gin.Context) {
	type Profile struct {
		Email string
	}
	sess := sessions.Default(ctx)
	id := sess.Get(userIdKey).(int64)
	u, err := c.svc.Profile(ctx, &userv1.ProfileRequest{
		Id: id,
	})
	if err != nil {
		// 按照道理来说，这边 id 对应的数据肯定存在，所以要是没找到，
		// 那就说明是系统出了问题。
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, Profile{
		Email: u.User.Email,
	})
}

//func (h *UserHandler) AvatarUpdate(ctx *gin.Context) {
//	file, fileTypes, err := ctx.Request.FormFile("file")
//	if file == nil || err != nil {
//		ctx.JSON(http.StatusOK, Result{
//			Code: 5,
//			Msg:  "图片上传失败",
//		})
//		return
//	}
//	if fileTypes.Size > 2*1024*1024 {
//		ctx.JSON(http.StatusOK, Result{
//			Code: 5,
//			Msg:  "文件不应该超过2MB",
//		})
//		return
//	}
//	fileType := strings.Split(fileTypes.Header.Get("Content-Type"), "/")[1]
//	user := ctx.MustGet("claims").(ijwt.UserClaims)
//	ossAdress, err := h.svc.AvatarUpdate(ctx.Request.Context(), user.Uid, file, fileType)
//	if err != nil {
//		ctx.JSON(http.StatusOK, Result{
//			Code: 5,
//			Msg:  "系统错误",
//		})
//		h.l.Info("AvatarUpdate错误", logger.Field{Key: "err", Val: err})
//		return
//	}
//	ctx.JSON(http.StatusOK, Result{
//		Code: 0,
//		Msg:  "ok",
//		Data: ossAdress,
//	})
//	file.Close()
//}
