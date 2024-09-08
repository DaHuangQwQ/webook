package system

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"webook/internal/api"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

type UserRepository interface {
	GetMenuIds(ctx context.Context, roleIds []uint) map[int64]int64
	GetAdminRoleIds(ctx context.Context, userId int64) (roleIds []uint, err error)
	List(ctx context.Context, req api.UserSearchReq) (total int, userList []domain.User, err error)
	Add(ctx context.Context, req api.UserAddReq) error
	DeleteByIds(ctx *gin.Context, ids []int) error
	GetUserInfoById(ctx *gin.Context, id uint64) (domain.User, error)
}

type CachedUserRepository struct {
	casbin           casbin.IEnforcer
	dao              dao.UserDao
	casBinUserPrefix string
}

func NewCachedUserRepository(casbin casbin.IEnforcer, dao dao.UserDao) UserRepository {
	return &CachedUserRepository{
		casbin:           casbin,
		casBinUserPrefix: "u_",
		dao:              dao,
	}
}
func (repo *CachedUserRepository) GetUserInfoById(ctx *gin.Context, id uint64) (domain.User, error) {
	user, err := repo.dao.FindById(ctx, int64(id))
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(user), nil
}

func (repo *CachedUserRepository) DeleteByIds(ctx *gin.Context, ids []int) error {
	return repo.dao.DeleteByIds(ctx, ids)
}

func (repo *CachedUserRepository) Add(ctx context.Context, req api.UserAddReq) error {
	userId, err := repo.dao.InsertAndGetId(ctx, dao.User{
		DeptID: req.DeptId,
		Email: sql.NullString{
			String: req.Email,
			Valid:  req.Email != "",
		},
		Nickname: req.NickName,
		Phone: sql.NullString{
			String: req.Mobile,
			Valid:  req.Mobile != "",
		},
		Remark:  req.Remark,
		Gender:  req.Sex,
		Status:  req.Status,
		IsAdmin: uint8(req.IsAdmin),
	})
	if err != nil {
		return err
	}
	return repo.addUserRole(ctx, userId, req.RoleIds)
}

func (repo *CachedUserRepository) List(ctx context.Context, req api.UserSearchReq) (total int, userList []domain.User, err error) {
	total, resList, err := repo.dao.FindAll(ctx, req)
	userList = make([]domain.User, len(resList))
	for i, res := range resList {
		userList[i] = repo.toDomain(res)
	}
	return
}

func (repo *CachedUserRepository) GetAdminRoleIds(ctx context.Context, userId int64) (roleIds []uint, err error) {
	groupPolicy, err := repo.casbin.GetFilteredGroupingPolicy(0, fmt.Sprintf("%s%d", repo.casBinUserPrefix, userId))
	if len(groupPolicy) > 0 {
		roleIds = make([]uint, len(groupPolicy))
		//得到角色id的切片

		for k, v := range groupPolicy {
			num, _ := strconv.ParseUint(v[1], 10, 64)
			roleIds[k] = uint(num)
		}
	}
	return
}

func (repo *CachedUserRepository) GetMenuIds(ctx context.Context, roleIds []uint) map[int64]int64 {
	menuIds := map[int64]int64{}
	for _, roleId := range roleIds {
		//查询当前权限
		gp, _ := repo.casbin.GetFilteredPolicy(0, strconv.Itoa(int(roleId)))
		for _, p := range gp {
			mid, _ := strconv.ParseInt(p[1], 10, 64)
			menuIds[mid] = mid
		}
	}
	return menuIds
}

func (repo *CachedUserRepository) toDomain(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email.String,
		Password: user.Password,
		Phone:    user.Phone.String,
		Nickname: user.Nickname,
		Grade:    user.Grade,
		Gender:   user.Gender,
		CTime:    time.UnixMilli(user.CTime),
		Avatar:   user.AvatarUrl,

		WechatInfo: domain.WechatInfo{
			OpenId:  user.WechatOpenId.String,
			UnionId: user.WechatUnionId.String,
		},

		Birthday:    user.Birthday,
		UserStatus:  uint(user.UserStatus),
		DeptId:      user.DeptID,
		Remark:      user.Remark,
		IsAdmin:     int(user.IsAdmin),
		Address:     user.Address,
		Describe:    user.Describe,
		LastLoginIp: user.LastLoginIP,
	}
}

func (repo *CachedUserRepository) addUserRole(ctx context.Context, userId int64, roleIds []int64) error {
	for _, v := range roleIds {
		_, err := repo.casbin.AddGroupingPolicy(fmt.Sprintf("%s%d", repo.casBinUserPrefix, userId), strconv.FormatInt(v, 10))
		return err
	}
	return nil
}
