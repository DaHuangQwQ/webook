package ioc

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/spf13/viper"
)

func InitCasbinService() casbin.IEnforcer {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var config = Config{
		DSN: "default",
	}
	err := viper.UnmarshalKey("db", &config)
	if err != nil {
		panic(err)
	}
	a, _ := gormadapter.NewAdapter("mysql", config.DSN, true)
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`)
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		panic(err)
	}
	e.LoadPolicy()
	return e
}
