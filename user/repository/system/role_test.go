package system

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
	"testing"
	"webook/user/domain"
)

func TestCachedRoleRepository_Save(t *testing.T) {
	a, _ := gormadapter.NewAdapter("mysql", "root:root@tcp(localhost:3307)/dahuang", true)
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
`)
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		panic(err)
	}
	e.LoadPolicy()
	roleRepo := NewCachedRoleRepository(e)
	err = roleRepo.Save(context.Background(), domain.Role{
		Id:  1,
		Ids: 2,
	})
	if err != nil {
		panic(err)
	}
}
