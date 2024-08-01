package web

import (
	"net/http"
	"regexp"
	"webook/internal/domain"
	"webook/internal/service"
)
import "github.com/gin-gonic/gin"

type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern  = "^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$"
		passwordExpPattern = "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[!@#$%^&*])[A-Za-z\\d!@#$%^&*]{8,}$\n"
	)
	emailExp := regexp.MustCompile(emailRegexPattern)
	passwordExp := regexp.MustCompile(passwordExpPattern)

	return &UserHandler{
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (u *UserHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/signup", u.SignUp)
}

func (u *UserHandler) Login(ctx *gin.Context) {

}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	err := u.svc.SignUp(ctx, domain.User{
		ID:        0,
		FirstName: "123",
		LastName:  "123",
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}
