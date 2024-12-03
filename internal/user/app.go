package user

import (
	"github.com/DaHuangQwQ/webook/internal/user/service"
)

type App struct {
	Server service.UserService
	//webAdmin *ginx.Server
}
